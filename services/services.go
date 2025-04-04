package services

import (
	clients "ai-api/clients"
	"ai-api/config"
	"net/http"
)

type Services struct {
	PRService *PRService
}

// NewServices creates a new Services instance
func NewServices(cfg config.Config) *Services {

	httpClient := &http.Client{}
	githubClient := &clients.GithubClient{
		HttpClient: httpClient,
		APIKey:     cfg.GithubToken,
		BaseURL:    "https://api.github.com",
	}

	prService := &PRService{
		githubClient: *githubClient,
		cfg:          cfg,
	}

	return &Services{
		PRService: prService,
	}
}
