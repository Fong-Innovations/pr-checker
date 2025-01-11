package main

import (
	router "ai-api/routers" // Import the router package
)

func main() {
	// Setup the router from the external package
	r := router.SetupRouter()

	// Start the Gin server on port 8080
	r.Run(":8080")
}
