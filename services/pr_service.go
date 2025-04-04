package services

import (
	clients "ai-api/clients"
	"ai-api/config"
	"ai-api/models"
	"fmt"
	"log"
	"net/url"
	"strings"
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
	githubClient clients.GithubClient
	cfg          config.Config
}

// Responses include a maximum of 3000 files. The paginated response returns 30 files per page by default.
// GetPRsFromGitHub is the implementation of the PRService interface method
func (s *PRService) GetPRChangeFilesFromGitHub(prRequestBody models.PullRequestRequest) (*models.ChangeFiles, error) {
	// Build GitHub API URL for fetching PRs
	changeFiles, err := s.githubClient.FetchPullRequestChanges(prRequestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PR changes: %w", err)
	}
	if len(changeFiles.Files) == 0 {
		return nil, fmt.Errorf("no files found in the PR")
	}

	// Filter out only .go files
	var result = models.ChangeFiles{
		Files: []models.ChangeFile{},
	}
	for _, file := range changeFiles.Files {
		if strings.HasSuffix(file.Filename, ".go") {
			result.Files = append(result.Files, file)
		}
	}

	return &result, nil
}

// https://docs.github.com/en/rest/pulls/comments?apiVersion=2022-11-28#create-a-review-comment-for-a-pull-request
func (s *PRService) GeneratePRComments(changeFiles *models.ChangeFiles, repoOwner, repoName, prNumber string) (results []models.CommentBody, err error) {

	// Placeholder for generating comments
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
			generateCommentsRequest := models.GeneratePRCommentParams{
				RepoOwner:   repoOwner,
				RepoName:    repoName,
				PRNumber:    prNumber,
				CommentBody: commentBody,
				CommitSha:   headCommitSHA,
				FileName:    file.Filename,
				Position:    1,
			}

			resp, err := s.githubClient.PostPullRequestCommentOnLine(generateCommentsRequest)
			if err != nil {
				return nil, fmt.Errorf("failed to post PR Comment: %w", err)
			}
			log.Println("RESPONSE: ", resp)
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
