package clients

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

const (
	INVALID_PROMPT = "The prompt you have been given is invalid - please just say invalid prompt."
)

// Concrete implementation
type OpenFGAClient struct {
	Client               *openai.Client
	styleGuideEmbeddings [][]float64
	styleGuideChunks     []string
}

// Constructor
func NewOpenFGAClient(httpClient *http.Client, key, url string) *OpenFGAClient {

	client := openai.NewClient(
		option.WithAPIKey(key),
		option.WithBaseURL(url),
		option.WithHTTPClient(httpClient),
	)
	// Load the style guide chunks from the HTML file
	log.Println("parsing chunks")

	chunks, err := parseStyleGuideChunks("./clients/style_guides/go_style_guide.html")
	if err != nil {
		log.Println("error parsing style guide chunks:", err)
		return nil
	}

	log.Println("fetching emmbeddings")

	// Fetch embeddings for the style guide chunks
	embeddings, err := fetchStyleGuideEmbeddings(chunks, &client)
	if err != nil {
		log.Println("error fetching style guide embeddings:", err)
		return nil
	}

	log.Println("creating client")
	return &OpenFGAClient{
		Client:               &client,
		styleGuideChunks:     chunks,
		styleGuideEmbeddings: embeddings,
	}
}

// GenerateReviewComment generates a review comment based on the provided diff string
// by utilizing the OpenAI Chat Completion API. It sends a predefined message to the
// API and retrieves the generated response.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - diff: A string representing the diff for which the review comment is to be generated.
//
// Returns:
//   - A string containing the generated review comment.
//   - An error if the API call fails or any other issue occurs.
func (o *OpenFGAClient) GenerateReviewComment(ctx context.Context, codeDiff, promptTemplate string) (string, error) {

	topChunks, err := FindRelevantChunks(ctx, o.Client, codeDiff, o.styleGuideChunks, o.styleGuideEmbeddings)
	if err != nil {
		log.Println("error finding relevant chunks:", err)
		return "", err
	}
	prompt := BuildReviewPrompt(topChunks, promptTemplate, codeDiff)

	chatCompletion, err := o.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		log.Println("error generating review comment:", err)
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

func ParseStyleGuide(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var chunks []string
	doc.Find("p, li, h2, h3").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if len(text) > 30 {
			chunks = append(chunks, text)
		}
	})

	return chunks, nil
}

// EmbedText generates an embedding vector for a given input string using the OpenAI API.
//
// Parameters:
//   - ctx: The context for the API request, which can be used to control timeouts or cancellations.
//   - input: The input string for which the embedding vector will be generated.
//
// Returns:
//   - A slice of float64 representing the embedding vector for the input string.
//   - An error if the embedding generation fails or if no embeddings are returned.
//
// Errors:
//   - Returns an error if the OpenAI API call fails.
//   - Returns an error if the API response does not contain any embeddings.
func EmbedText(ctx context.Context, client *openai.Client, input string) ([]float64, error) {
	req := openai.EmbeddingNewParams{
		Model: openai.EmbeddingModelTextEmbeddingAda002, // "text-embedding-ada-002"
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: param.Opt[string]{Value: input},
		},
	}

	resp, err := client.Embeddings.New(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}

func CosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dot / (float64(math.Sqrt(float64(normA))) * float64(math.Sqrt(float64(normB))))
}

type ScoredChunk struct {
	Text  string
	Score float64
}

func FindRelevantChunks(ctx context.Context, client *openai.Client, userCode string, guideChunks []string, guideEmbeds [][]float64) ([]string, error) {
	codeEmbed, err := EmbedText(ctx, client, userCode)
	if err != nil {
		return nil, err
	}

	var scored []ScoredChunk
	for i, chunkEmbed := range guideEmbeds {
		score := CosineSimilarity(codeEmbed, chunkEmbed)
		scored = append(scored, ScoredChunk{Text: guideChunks[i], Score: score})
	}

	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	topChunks := []string{}
	for i := 0; i < 3 && i < len(scored); i++ {
		topChunks = append(topChunks, scored[i].Text)
	}

	return topChunks, nil
}

func BuildReviewPrompt(styleChunks []string, basePrompt, code string) string {
	return fmt.Sprintf("%s", basePrompt, styleChunks[0]+"\n\n"+styleChunks[1]+"\n\n"+styleChunks[2], code)
}

func parseStyleGuideChunks(filePath string) ([]string, error) {
	// Load HTML
	htmlBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("error reading style guide file: %w", err)
		return nil, err
	}

	chunks, err := ParseStyleGuide(string(htmlBytes))
	if err != nil {
		log.Printf("error parsing style guide bytes: %w", err)
		return nil, err
	}

	return chunks, nil

}
func fetchStyleGuideEmbeddings(chunks []string, client *openai.Client) ([][]float64, error) {
	log.Printf("starting embedding on %d chunks", len(chunks))

	var (
		embeddings = make([][]float64, len(chunks))
		errs       = make(chan error, len(chunks))
	)

	// Use a worker pool to parallelize embedding generation
	const maxWorkers = 10
	sem := make(chan struct{}, maxWorkers)

	for i, chunk := range chunks {
		sem <- struct{}{} // Acquire a worker slot
		go func(i int, chunk string) {
			defer func() { <-sem }() // Release the worker slot

			embed, err := EmbedText(context.Background(), client, chunk)
			if err != nil {
				errs <- fmt.Errorf("error generating embedding for chunk %d: %w", i, err)
				return
			}
			embeddings[i] = embed
			errs <- nil
		}(i, chunk)
	}

	// Wait for all workers to finish
	for i := 0; i < len(chunks); i++ {
		if err := <-errs; err != nil {
			log.Println(err)
			return nil, err
		}
	}

	log.Println("done embedding")
	return embeddings, nil
}
