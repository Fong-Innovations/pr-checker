package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Interface for mocking and extensibility
type LLMClient interface {
	GenerateReviewComment(diff string) (string, error)
}

// Concrete implementation
type OpenAIClient struct {
	APIKey string
	Model  string
	URL    string
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

type openAIResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

// Constructor
func NewOpenAIClient() *OpenAIClient {
	return &OpenAIClient{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  "gpt-4", // or gpt-3.5-turbo
		URL:    "https://api.openai.com/v1/chat/completions",
	}
}

// Implementation of GenerateReviewComment
func (c *OpenAIClient) GenerateReviewComment(diff string) (string, error) {
	reqBody := openAIRequest{
		Model: c.Model,
		Messages: []chatMessage{
			{Role: "system", Content: "You're a senior software engineer reviewing pull requests."},
			{Role: "user", Content: fmt.Sprintf("Review this code diff and generate feedback:\n\n%s", diff)},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error: %s", string(bodyBytes))
	}

	var res openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return res.Choices[0].Message.Content, nil
}
