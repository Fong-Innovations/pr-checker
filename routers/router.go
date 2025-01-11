package router

import (
	handler "ai-api/handlers" // Import the handler package

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up all routes for the application
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Define the /hello route
	// Simple group: v1
	api := r.Group("/v1/api")
	{
		api.GET("/hello", handler.HelloHandler)

	}

	return r
}
