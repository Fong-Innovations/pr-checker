package main

import (
	"ai-api/config"
	router "ai-api/server"
	"ai-api/services"
)

func main() {

	// Setup the router from the external package
	// Load configuration from the .env file
	cfg := config.LoadConfig(".env")
	services := services.NewServices(*cfg)
	server := router.NewServer(cfg, services)

	server.Router.Run("localhost:8080")

}
