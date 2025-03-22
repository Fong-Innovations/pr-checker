package services

import (
	"ai-api/config"
	"ai-api/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	githubPullRequestURL = "https://api.github.com/repos/%s/%s/pulls/%s"
)

// PRService is an interface that defines methods for PR-related operations
type PRService interface {
	GetPRsFromGitHub(prRequestBody models.PullRequestRequest) ([]models.PullRequestResponse, error)
}

// PRServiceImpl is a concrete implementation of the PRService interface
type PRServiceImpl struct {
	httpClient *http.Client
	cfg        config.Config
	// Add any dependencies or configurations needed
}

// GetPRsFromGitHub is the implementation of the PRService interface method
func (s *PRServiceImpl) GetPRsFromGitHub(prRequestBody models.PullRequestRequest) (*models.PullRequestResponse, error) {
	// Build GitHub API URL for fetching PRs
	url := fmt.Sprintf(githubPullRequestURL, prRequestBody.OwnerID, prRequestBody.RepoID, prRequestBody.ID)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers to the request if needed
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	// req.Header.Set("Authorization: Bearer", s.cfg.GithubToken)
	log.Println("s.cfg.GithubToken")
	log.Println(s.cfg.GithubToken)
	req.Header.Set("Authorization", "Bearer "+s.cfg.GithubToken)

	// Add more headers if necessary

	// Send the HTTP request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PRs from GitHub: %w", err)
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response from GitHub: %s", resp.Status)
	}

	// Decode the JSON response body into a slice of PRs
	var pr *models.PullRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(pr.PR); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub response: %w", err)
	}

	return pr, nil
}
