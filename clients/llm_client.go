package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Interface for mocking and extensibility
type AIClientInterface interface {
	GenerateReviewComment(diff string) (string, error)
}

// Concrete implementation
type AIClient struct {
	HttpClient *http.Client
	APIKey     string
	Model      string
	URL        string
}

// Request/Response format (simplified for OpenAI Chat API, can adjust per your LLM)
type openAIRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type aiResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

// Constructor
func NewAIClient(httpClient *http.Client, key, model, url string) *AIClient {
	return &AIClient{
		APIKey: key,
		Model:  model, // or gpt-3.5-turbo
		URL:    url,
	}
}

// Implementation of GenerateReviewComment
func (c *AIClient) GenerateReviewComment(diff string) (string, error) {
	reqBody := openAIRequest{
		Model: c.Model,
		Messages: []chatMessage{
			{Role: "system", Content: "You're a senior software engineer reviewing pull requests."},
			{Role: "user", Content: fmt.Sprintf("Review this code diff and generate feedback:\n\n%s", diff)},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("error marshalling request body:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("error creating request for LLM comments:", err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %s", resp.Status)
	}

	var res aiResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode open ai response body: %w", err)
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return res.Choices[0].Message.Content, nil
}
