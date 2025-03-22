package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloHandler is the handler for the /hello route
func HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from external router!",
	})
}

// type AnalyzePRRequest struct {
// 	URL string `json:"url"`
// }

// // AnalyzePR analyzes the pull request and posts a comment if needed
// func AnalyzePR(c *gin.Context) {
// 	// Define a struct to capture the PR link from the request body

// 	req := AnalyzePRRequest{}
// 	// Bind the JSON body to the struct
// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	// Call the analyzePR function (this can be a security check or analysis)
// 	message := analyzePR(req.URL)

// 	// Post a comment on the PR (this uses GitHub API or any other platform)
// 	err = postCommentToPR(req.URL, message)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post comment"})
// 		return
// 	}

// 	// Return a success message with the result
// 	c.JSON(http.StatusOK, gin.H{"message": message})
// }

// // analyzePR checks for any issues in the PR link and returns a message
// func analyzePR(prLink string) string {
// 	// Simulate analysis; you would replace this with actual logic
// 	// For now, we'll just check if the PR link contains the word "security_issue"

// 	if prLink == "https://github.com/some/repo/pull/security_issue" {
// 		return "Security issue found! Please review the code."
// 	}
// 	return "All good! No security issues found."
// }

// // postCommentToPR posts a comment on the PR
// func postCommentToPR(prLink, message string) error {
// 	// You would integrate the GitHub API (or other platform) to post the comment
// 	// For now, we will simulate a successful comment posting
// 	return nil
// }
