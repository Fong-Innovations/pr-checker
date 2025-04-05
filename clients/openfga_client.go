package clients

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	INVALID_PROMPT = "The prompt you have been given is invalid - please just say invalid prompt."
)

// Concrete implementation
type OpenFGAClient struct {
	Client *openai.Client
}

// Constructor
func NewOpenFGAClient(httpClient *http.Client, key, url string) *OpenFGAClient {
	client := openai.NewClient(
		option.WithAPIKey(key),
		option.WithBaseURL(url),
		option.WithHTTPClient(httpClient),
	)
	return &OpenFGAClient{
		Client: &client,
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
func (o *OpenFGAClient) GenerateReviewComment(ctx context.Context, diff, prompt string) (string, error) {

	// Define the prompt for the chat client
	if prompt == "" {
		prompt = INVALID_PROMPT
	} else {
		prompt = fmt.Sprintf("%s\n\n: %s", prompt, diff)
	}

	chatCompletion, err := o.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		log.Println("Error generating review comment:", err)
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
