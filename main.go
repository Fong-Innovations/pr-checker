package main

import (
	"pr-checker/config"
	router "pr-checker/server"
	"pr-checker/services"

	"github.com/gofiber/fiber/v2/log"
)

func main() {

	// Setup the router from the external package
	// Load configuration from the .env file
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		// Handle error
		log.Fatal(err)
		return
	}
	services := services.NewServices(*cfg)
	server := router.NewServer(cfg, services)

	server.Router.Run("localhost:8080")

}
