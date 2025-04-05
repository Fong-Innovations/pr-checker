package clients

import (
	"ai-api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Concrete implementation
type GithubClient struct {
	HttpClient *http.Client
	APIKey     string
	BaseURL    string
}

type GithubClientInterface interface {
	FetchPullRequestChanges(prRequestBody models.PullRequestRequest) (*models.ChangeFiles, error)
	PostPullRequestCommentOnLine(params models.GeneratePRCommentParams) (results []models.CommentBody, err error)
}

func NewGithubClient(httpClient *http.Client, apiKey, baseUrl string) *GithubClient {
	return &GithubClient{
		HttpClient: httpClient,
		APIKey:     apiKey,
		BaseURL:    baseUrl,
	}
}

const (
	githubFetchPRChangesURL = "https://api.github.com/repos/%s/%s/pulls/%s/files"
	githubPostPRCommentURL  = "https://api.github.com/repos/%s/%s/pulls/%s/comments" // github treats prs as issues for comments!
)

func (g *GithubClient) FetchPullRequestChanges(prRequestBody models.PullRequestRequest) (*models.ChangeFiles, error) {
	// Create a new HTTP request
	url := fmt.Sprintf(githubFetchPRChangesURL, prRequestBody.OwnerID, prRequestBody.RepoID, prRequestBody.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.full+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", "Bearer "+g.APIKey)

	resp, err := g.HttpClient.Do(req)
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
	err = json.NewDecoder(resp.Body).Decode(&prResponse.Files)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PR response body: %w", err)
	}

	return &prResponse, nil
}

func (g *GithubClient) PostPullRequestCommentOnLine(params models.GeneratePRCommentParams) (results []models.CommentBody, err error) {

	url := fmt.Sprintf(githubPostPRCommentURL, params.RepoOwner, params.RepoName, params.PRNumber)
	prReviewCommentRequestBody := models.CommentBody{
		Body:     params.CommentBody,
		CommitID: params.CommitSha,
		Path:     params.FileName,
		Position: params.Position,
	}
	jsonData, err := json.Marshal(prReviewCommentRequestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PR Comment body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making PR Comment request: %w", err)
	}

	req.Header.Set("Authorization", "token "+g.APIKey)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := g.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making PR Comment request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error posting PR comment: %s", resp.Status)
	}

	results = append(results, prReviewCommentRequestBody)
	return results, nil
}
