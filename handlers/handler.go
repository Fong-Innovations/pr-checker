package handler

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
