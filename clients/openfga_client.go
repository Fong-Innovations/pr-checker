package clients

import (
	"context"
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

// OpenFGAClient is a struct that represents a client for interacting with the OpenAI API.
// It contains a pointer to the OpenAI client and slices for storing
// style guide embeddings and chunks. The struct is used to generate review comments
// based on code diffs and style guides.
type OpenFGAClient struct {
	Client               *openai.Client
	styleGuideEmbeddings [][]float64
	styleGuideChunks     []string
}

// ScoredChunk represents a chunk of text with its associated score.
// It is used to store the text and its similarity score when comparing
// it to a user-provided code snippet.
// The score is typically a floating-point number representing the
// similarity between the chunk and the user code, calculated using
// cosine similarity or another metric.
type ScoredChunk struct {
	Text  string
	Score float64
}

// OpenFGAClientInterface defines the methods for interacting with the OpenAI API
type OpenFGAClientInterface interface {
	GenerateReviewComment(ctx context.Context, codeDiff, promptTemplate string) (string, error)
}

// NewOpenFGAClient creates a new instance of OpenFGAClient with the provided HTTP client, API key, and base URL.
// It initializes the OpenAI client, parses style guide chunks from an HTML file, and fetches embeddings for the
// style guide chunks. If any error occurs during the initialization process, it logs the error and returns nil.
//
// Parameters:
//   - httpClient: The HTTP client to be used for making requests.
//   - key: The API key for authenticating with the OpenAI service.
//   - url: The base URL for the OpenAI service.
//
// Returns:
//   - A pointer to an OpenFGAClient instance if successful, or nil if an error occurs.
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
func NewOpenFGAClient(httpClient *http.Client, key, url string) *OpenFGAClient {

	client := openai.NewClient(
		option.WithAPIKey(key),
		option.WithBaseURL(url),
		option.WithHTTPClient(httpClient),
	)
	// Load the style guide chunks from the HTML file
	log.Println("loading style guide embeddings")
	chunks, err := parseStyleGuideChunks("./clients/style_guides/go_style_guide.html")
	if err != nil {
		log.Println("error parsing style guide chunks:", err)
		return nil
	}

	// Fetch embeddings for the style guide chunks
	embeddings, err := fetchStyleGuideEmbeddings(chunks, &client)
	if err != nil {
		log.Println("error fetching style guide embeddings:", err)
		return nil
	}

	log.Println("creating openfga client")
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
		return "", fmt.Errorf("error finding relevant chunks: %w", err)
	}
	prompt := buildReviewPrompt(topChunks, promptTemplate, codeDiff)

	chatCompletion, err := o.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		return "", fmt.Errorf("error generating review comment: %w", err)
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

// parseStyleGuide extracts and returns textual content from the provided HTML string.
// It parses the HTML using goquery and selects text from <p>, <li>, <h2>, and <h3> elements.
// Only text content with a length greater than 30 characters is included in the result.
//
// Parameters:
//   - html: A string containing the HTML content to parse.
//
// Returns:
//   - []string: A slice of strings containing the extracted text chunks.
//   - error: An error if the HTML parsing fails.
func parseStyleGuide(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML style guide: %w", err)
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
		return nil, fmt.Errorf("error generating embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}

// CosineSimilarity calculates the cosine similarity between two vectors a and b.
// The cosine similarity is a measure of similarity between two non-zero vectors
// of an inner product space that measures the cosine of the angle between them.
// It is defined as the dot product of the vectors divided by the product of their magnitudes.
//
// Parameters:
//   - a: A slice of float64 representing the first vector.
//   - b: A slice of float64 representing the second vector.
//
// Returns:
//   - A float64 value representing the cosine similarity between the two vectors.
//     The result ranges from -1 (exactly opposite) to 1 (exactly the same), with 0
//     indicating orthogonality (no similarity).
//
// Note:
//   - The input slices a and b must have the same length.
//   - If either vector has a magnitude of zero, the function may produce a NaN result
//     due to division by zero.
func CosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dot / (float64(math.Sqrt(float64(normA))) * float64(math.Sqrt(float64(normB))))
}

// FindRelevantChunks identifies the most relevant guide chunks based on their
// semantic similarity to the provided user code.
//
// This function computes an embedding for the user-provided code and compares
// it to precomputed embeddings of guide chunks using cosine similarity. It
// then selects the top 3 most relevant chunks based on their similarity scores.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - client: An OpenAI client used for generating embeddings.
//   - userCode: The code snippet provided by the user.
//   - guideChunks: A slice of guide text chunks to compare against.
//   - guideEmbeds: A slice of precomputed embeddings corresponding to the guide chunks.
//
// Returns:
//   - A slice of the top 3 most relevant guide chunks, sorted by similarity score.
//   - An error if embedding generation or any other operation fails.
func FindRelevantChunks(ctx context.Context, client *openai.Client, userCode string, guideChunks []string, guideEmbeds [][]float64) ([]string, error) {
	codeEmbed, err := EmbedText(ctx, client, userCode)
	if err != nil {
		return nil, fmt.Errorf("error generating embedding for user code: %w", err)
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

// buildReviewPrompt constructs a review prompt by combining a base prompt,
// style guide chunks, and the code to be reviewed. It formats the output
// as a single string.
//
// Parameters:
//   - styleChunks: A slice of strings representing chunks of the style guide.
//     It is expected to have at least three elements.
//   - basePrompt: A string containing the base prompt or introductory text.
//   - code: A string containing the code that needs to be reviewed.
//
// Returns:
//
//	A formatted string that includes the base prompt, the style guide,
//	and the code to be reviewed.
func buildReviewPrompt(styleChunks []string, basePrompt, code string) string {
	return fmt.Sprintf("%s. Here is the style guide: %s\n\n Here is the code to review: \n\n%s", basePrompt, styleChunks[0]+"\n\n"+styleChunks[1]+"\n\n"+styleChunks[2], code)
}

// parseStyleGuideChunks reads an HTML file from the specified file path,
// parses its content into chunks using the ParseStyleGuide function, and
// returns the resulting chunks as a slice of strings. If an error occurs
// during file reading or parsing, it returns an error.
//
// Parameters:
//   - filePath: The path to the HTML file to be read and parsed.
//
// Returns:
//   - []string: A slice of strings representing the parsed chunks of the style guide.
//   - error: An error if the file cannot be read or parsed.
func parseStyleGuideChunks(filePath string) ([]string, error) {
	// Load HTML
	htmlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading style guide file: %w", err)
	}

	chunks, err := parseStyleGuide(string(htmlBytes))
	if err != nil {
		return nil, fmt.Errorf("error parsing style guide bytes: %w", err)
	}

	return chunks, nil

}

// fetchStyleGuideEmbeddings generates embeddings for a list of text chunks using the OpenAI client.
// It processes the chunks in parallel using a worker pool to improve performance.
//
// Parameters:
//   - chunks: A slice of strings, where each string represents a chunk of text to be embedded.
//   - client: An instance of the OpenAI client used to generate embeddings.
//
// Returns:
//   - A 2D slice of float64 values, where each inner slice represents the embedding for a corresponding chunk.
//   - An error if any embedding generation fails.
//
// The function uses a maximum of 10 concurrent workers to process the chunks. If an error occurs
// while generating an embedding for a chunk, the function logs the error and returns it immediately.
func fetchStyleGuideEmbeddings(chunks []string, client *openai.Client) ([][]float64, error) {
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
	for range chunks {
		if err := <-errs; err != nil {
			return nil, fmt.Errorf("error embedding chunk: %w", err)
		}
	}

	return embeddings, nil
}
