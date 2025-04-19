package handlers

import (
	"fmt"
	"net/http"
	"pr-checker/models"
	"pr-checker/services"

	"github.com/gin-gonic/gin"
)

// PRHandler represents the handler for handling PR-related requests
type PRHandler struct {
	Service *services.PRService
}

type PRHandlerInterface interface {
	GetPR(c *gin.Context)
	AnalyzePR(c *gin.Context)
}

// NewPRHandler creates a new PR handler
func NewPRHandler(service *services.PRService) *PRHandler {
	return &PRHandler{
		Service: service,
	}
}

// GetPR handles GET requests to fetch a single PR by ID
// GetPR handles the HTTP request to fetch pull request details from GitHub.
// It parses the request body to extract the necessary information, calls the
// service to retrieve the pull request details, and returns the details in the
// response.
//
// @Summary Fetch pull request details
// @Description Fetches pull request details from GitHub based on the provided request body.
// @Tags pull requests
// @Accept json
// @Produce json
// @Param request body FetchPullRequestBody true "Fetch Pull Request Body"
// @Success 200 {object} gin.H{"message": string, "pr": interface{}}
// @Failure 400 {object} gin.H{"error": string}
// @Router /pullrequest [post]
func (h *PRHandler) AnalyzePR(ctx *gin.Context) {

	// parse pr request data
	prRequestBody, err := parseFetchPullRequestBody(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error parsing request body. error:": err.Error()})
		return
	}

	// fetch changes from github for requested pr
	pr, err := h.Service.GetPRChangeFilesFromGitHub(ctx, *prRequestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error fetching pr changes", "error:": err.Error()})
		return
	}

	// analyze the change files and generate a list of comments
	codeReviews, err := h.Service.ReviewChanges(ctx, pr, prRequestBody.OwnerID, prRequestBody.RepoID, prRequestBody.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error posting PR comments", "error: ": err.Error()})
		return
	}

	status, err := h.Service.PostPRComments(ctx, codeReviews)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error posting PR comments", "error: ": err.Error()})
		return
	}

	// return status
	ctx.JSON(http.StatusOK, gin.H{"message": "PR Analyzed", "commented_files_count": len(codeReviews), "status": status})
}

func (h *PRHandler) StorePRData(ctx *gin.Context) {
	// parse pr request data
	prRequestBody, err := parseFetchPullRequestBody(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error parsing request body. error:": err.Error()})
		return
	}
	// fetch changes from github for requested pr
	pr, err := h.Service.GetPRDataFromGithub(ctx, *prRequestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error fetching pr changes", "error:": err.Error()})
		return
	}
	smaller := models.PullRequestDBEntry{
		Repo:         pr.Base.Repo.Name,
		TargetBranch: pr.Base.Ref,
		Merged:       pr.Merged,
		SourceBranch: pr.Head.Ref,
		Comments:     pr.Comments,
		ChangedFiles: pr.ChangedFiles,
		OpenedAt:     pr.CreatedAt,
		MergedAt:     *pr.MergedAt,
		ClosedAt:     *pr.ClosedAt,
		IssueUrl:     pr.IssueURL,
		User:         pr.User.Login,
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "PR Fetched", "pr": smaller, "status": http.StatusOK})

}

func parseFetchPullRequestBody(c *gin.Context) (*models.PullRequestRequest, error) {
	req := &models.PullRequestRequest{}
	req.OwnerID = c.Param("owner")
	req.RepoID = c.Param("repo")
	req.ID = c.Param("id")
	if req.ID == "" || req.OwnerID == "" || req.RepoID == "" {
		return nil, fmt.Errorf("missing field in FetchPullRequest body")
	}
	return req, nil
}
