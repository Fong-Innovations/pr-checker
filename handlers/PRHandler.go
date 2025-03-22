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
	Service services.PRService
}

// NewPRHandler creates a new PR handler
func NewPRHandler(service services.PRService) *PRHandler {
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
func (h *PRHandler) GetPR(c *gin.Context) {
	prRequestBody, err := parseFetchPullRequestBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr, err := h.Service.GetPRsFromGitHub(*prRequestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PR details", "pr": pr})
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
