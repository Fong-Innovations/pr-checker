package services

import (
	"ai-api/config"
	"ai-api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	githubFetchPRChangesURL = "https://api.github.com/repos/%s/%s/pulls/%s/files"
	githubPostPRCommentURL  = "https://api.github.com/repos/%s/%s/pulls/%s/comments" // github treats prs as issues for comments!
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

// https://docs.github.com/en/rest/pulls/comments?apiVersion=2022-11-28#create-a-review-comment-for-a-pull-request
func (s *PRService) GeneratePRComments(changeFiles *models.ChangeFiles, repoOwner, repoName, prNumber string) (results []models.CommentBody, err error) {

	// Placeholder for generating comments
	url := fmt.Sprintf(githubPostPRCommentURL, repoOwner, repoName, prNumber)
	for _, file := range changeFiles.Files {
		if file.Filename == "config/config.go" {
			// get the sha from the contents url
			headCommitSHA, err := parseRefForHeadCommitSHA(file.Contents_url)
			if err != nil {
				return nil, fmt.Errorf("failed to extra head commit sha: %w", err)
			}

			commentBody, err := generateCommentBody(file.Patch)
			if err != nil {
				return nil, fmt.Errorf("failed to generate comment body: %w", err)
			}

			prReviewCommentRequestBody := models.CommentBody{
				Body:     commentBody,
				CommitID: headCommitSHA,
				Path:     file.Filename,
				Position: 1,
			}
			jsonData, err := json.Marshal(prReviewCommentRequestBody)

			if err != nil {
				return nil, fmt.Errorf("failed to marshal PR Comment body: %w", err)
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				return nil, fmt.Errorf("error making PR Comment request: %w", err)
			}
			req.Header.Set("Authorization", "token "+s.cfg.GithubToken)
			req.Header.Set("Accept", "application/vnd.github+json")
			req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
			log.Print(req)
			resp, err := s.httpClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				return nil, fmt.Errorf(resp.Status)
			}

			results = append(results, prReviewCommentRequestBody)

			fmt.Println("Comment posted for: ", file.Filename)
		}
	}
	return results, nil
}

// parseRefForHeadCommitSHA parses the rawURL string to get the head commit SHA for a PR
func parseRefForHeadCommitSHA(rawURL string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Extract query parameters
	queryParams := parsedURL.Query()

	// Get the value of the 'ref' parameter
	sha := queryParams.Get("ref")
	return sha, nil
}

func generateCommentBody(changePatch string) (string, error) {
	// Placeholder for generating a comment body
	log.Println(changePatch)
	commentBody := "This is a sample comment body from the API"
	return commentBody, nil
}
