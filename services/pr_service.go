package services

import (
	"ai-api/config"
	"ai-api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	githubFetchPRChangesURL = "https://api.github.com/repos/%s/%s/pulls/%s/files"
	githubPostPRCommentURL  = "https://api.github.com/repos/%s/%s/pulls/%s/comments"
)

// DiffEntry represents a single entry in the diff response from GitHub
type DiffEntry struct {
	SHA              string `json:"sha"`
	Filename         string `json:"filename"`
	Status           string `json:"status"`
	Additions        int    `json:"additions"`
	Deletions        int    `json:"deletions"`
	Changes          int    `json:"changes"`
	BlobURL          string `json:"blob_url"`
	RawURL           string `json:"raw_url"`
	ContentsURL      string `json:"contents_url"`
	Patch            string `json:"patch,omitempty"`             // Optional field
	PreviousFilename string `json:"previous_filename,omitempty"` // Optional field
}

// PRService is a concrete implementation of the PRService interface
type PRService struct {
	httpClient *http.Client
	cfg        config.Config
}

// Responses include a maximum of 3000 files. The paginated response returns 30 files per page by default.
// GetPRsFromGitHub is the implementation of the PRService interface method
func (s *PRService) GetPRChangeFilesFromGitHub(prRequestBody models.PullRequestRequest) (*models.ChangeFiles, error) {
	// Build GitHub API URL for fetching PRs
	url := fmt.Sprintf(githubFetchPRChangesURL, prRequestBody.OwnerID, prRequestBody.RepoID, prRequestBody.ID)

	// Create a new HTTP request
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers to the request if needed
	req.Header.Set("Accept", "application/vnd.github.full+json")

	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", "Bearer "+s.cfg.GithubToken)

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

	// Parse the response body into a Go struct
	var prResponse models.ChangeFiles
	var result = models.ChangeFiles{
		Files: []models.ChangeFile{},
	}
	err = json.NewDecoder(resp.Body).Decode(&prResponse.Files)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PR response body: %w", err)
	}

	// RETURN ONLY .go FILES
	for _, file := range prResponse.Files {
		if strings.HasSuffix(file.Filename, ".go") {
			result.Files = append(result.Files, file)
		}
	}

	return &result, nil
}

func (s *PRService) GeneratePRComments(changeFiles *models.ChangeFiles, repoOwner, repoName, prNumber string) (results []models.Comment, err error) {
	// Placeholder for generating comments
	url := fmt.Sprintf(githubPostPRCommentURL, repoOwner, repoName, prNumber)
	for _, file := range changeFiles.Files {
		commentBody := models.Comment{
			Owner:      repoOwner,
			Repo:       repoName,
			PullNumber: prNumber,
			Body: models.CommentBody{
				Body:        "This is a test comment",
				CommitID:    file.Sha,
				Path:        file.Filename,
				Line:        1,
				Side:        "RIGHT",
				SubjectType: "line",
			},
		}

		jsonData, err := json.Marshal(commentBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal PR Comment body: %w", err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("error making PR Comment request: %w", err)
		}
		req.Header.Set("Authorization", "token "+s.cfg.GithubToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return nil, fmt.Errorf("posting comment received non-OK response from GitHub: %s", resp.Status)
		}

		results = append(results, commentBody)

		fmt.Println("Comment posted for: ", file.Filename)
	}
	return results, nil
}
