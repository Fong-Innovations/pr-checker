package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/gofiber/fiber/v2/log"
)

type PullRequest struct {
	ID       string
	Owner    string
	Repo     string
	Comments []Comment
}

type PullRequestRequest struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	RepoID  string `json:"repo_id"`
}

type PullRequestResponse struct {
	ID          string
	PullRequest PullRequest
}

type Comment struct {
	ID       string
	Category string
	Content  string
}

// GetPR handles GET requests to fetch a single PR by ID
func GetPR(c *gin.Context) {

	requestBody, err := parseFetchPullRequestBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Logic to fetch a PR by ID from the database or service
	c.JSON(http.StatusOK, gin.H{"message": "PR details", "pr": &requestBody})
}

func parseFetchPullRequestBody(c *gin.Context) (req *PullRequestRequest, err error) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading FetchPullRequest body")
		return nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("Error unmarshalling FetchPullRequest body")
		return nil, err
	}

	if req.ID == "" || req.OwnerID == "" || req.RepoID == "" {
		return nil, fmt.Errorf("missing field in FetchPullRequest body")
	}

	return req, nil
}
