package main

import (
	"ai-api/config"
	router "ai-api/server"
)

func main() {

	// Setup the router from the external package
	// Load configuration from the .env file
	cfg := config.LoadConfig(".env")
	server := router.NewServer(cfg)

	server.Router.Run("localhost:8080")

}
