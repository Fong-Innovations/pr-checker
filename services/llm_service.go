package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LLMService struct {
	httpClient *http.Client
	apiURL     string
	apiKey     string
}

func NewLLMService(httpClient *http.Client, apiURL, apiKey string) *LLMService {
	return &LLMService{
		httpClient: httpClient,
		apiURL:     apiURL,
		apiKey:     apiKey,
	}
}

func (s *LLMService) GenerateComment(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"prompt":      prompt,
		"max_tokens":  100,
		"temperature": 0.7,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal LLM request body: %w", err)
	}

	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create LLM request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send LLM request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response from LLM: %s", resp.Status)
	}

	var responseBody struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return "", fmt.Errorf("failed to decode LLM response: %w", err)
	}

	if len(responseBody.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from LLM")
	}

	return responseBody.Choices[0].Text, nil
}
