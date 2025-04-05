package handlers

import (
	"ai-api/models"
	"ai-api/services"
	"fmt"
	"net/http"

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

	// get pr head commit sha

	// fetch changes from github for requested pr
	pr, err := h.Service.GetPRChangeFilesFromGitHub(ctx, *prRequestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error fetching pr changes", "error:": err.Error()})
		return
	}

	// analyze the change files and generate a list of comments
	filesCommented, err := h.Service.GeneratePRComments(ctx, pr, prRequestBody.OwnerID, prRequestBody.RepoID, prRequestBody.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error posting PR comments", "error: ": err.Error()})
		return
	}

	// return status
	ctx.JSON(http.StatusOK, gin.H{"message": "PR Analyzed", "commented_files_count": len(filesCommented), "files": filesCommented})
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
