package services

import (
	clients "ai-api/clients"
	"ai-api/config"
	"ai-api/models"
	"context"
	"fmt"
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
	llmClient    clients.OpenFGAClient
	cfg          config.Config
}

// Responses include a maximum of 3000 files. The paginated response returns 30 files per page by default.
// GetPRsFromGitHub is the implementation of the PRService interface method
func (s *PRService) GetPRChangeFilesFromGitHub(ctx context.Context, prRequestBody models.PullRequestRequest) (*models.ChangeFiles, error) {
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

// ReviewChanges reviews the changes in a pull request by analyzing the provided change files
// and generating review comments for each file.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - changeFiles: A pointer to a models.ChangeFiles object containing the list of changed files.
//   - repoOwner: The owner of the repository where the pull request resides.
//   - repoName: The name of the repository where the pull request resides.
//   - prNumber: The pull request number.
//
// Returns:
//   - reviews: A slice of models.GeneratePRCommentParams containing the generated review comments.
//   - err: An error if any issue occurs during the review process.
//
// The function iterates over the list of changed files, extracts the head commit SHA from the file's
// contents URL, and uses an LLM client to generate a review comment body based on the file's patch.
// It then constructs a GeneratePRCommentParams object for each file and appends it to the reviews slice.
func (s *PRService) ReviewChanges(ctx context.Context, changeFiles *models.ChangeFiles, repoOwner, repoName, prNumber string) (reviews []models.GeneratePRCommentParams, err error) {
	for _, file := range changeFiles.Files {
		// get the sha from the contents url (find a better way to do this?)
		headCommitSHA, err := parseRefForHeadCommitSHA(file.Contents_url)
		if err != nil {
			return nil, fmt.Errorf("failed to extra head commit sha: %w", err)
		}

		// Generate the comment body using the LLM client
		commentBody, err := s.llmClient.GenerateReviewComment(ctx, file.Patch, s.cfg.LLMAnalyzePrompt)
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

		reviews = append(reviews, generateCommentsRequest)
	}
	return reviews, nil
}

func (s *PRService) PostPRComments(ctx context.Context, codeReviews []models.GeneratePRCommentParams) (status string, err error) {
	var failedComments []models.GeneratePRCommentParams

	for _, codeReview := range codeReviews {
		resp, err := s.githubClient.PostPullRequestCommentOnLine(codeReview)
		if err != nil {
			// Log the failed comment and continue with the next one
			fmt.Printf("failed to post comment for file %s: %v\n", codeReview.FileName, err)
			failedComments = append(failedComments, codeReview)
			continue
		}
		fmt.Println("Comment posted for: ", codeReview.FileName)
		status = resp[0].Body
	}

	if len(failedComments) > 0 {
		return "", fmt.Errorf("some comments failed to post: %v", failedComments)
	}
	// If all comments were posted successfully, return the status
	if status == "" {
		return "", fmt.Errorf("no comments posted")
	}
	return status, nil
}

// parseRefForHeadCommitSHA parses the rawURL string to get the head commit SHA for a PR
func parseRefForHeadCommitSHA(rawURL string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract query parameters
	queryParams := parsedURL.Query()

	// Get the value of the 'ref' parameter
	sha := queryParams.Get("ref")
	return sha, nil
}
