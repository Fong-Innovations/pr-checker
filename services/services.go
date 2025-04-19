package services

import (
	"net/http"
	clients "pr-checker/clients"
	"pr-checker/config"
	"time"
)

type Services struct {
	PRService *PRService
}

// NewServices creates a new Services instance
func NewServices(cfg config.Config) *Services {

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}
	githubClient := clients.NewGithubClient(httpClient, cfg.GithubToken, cfg.GithubBaseURL)
	openFGAClient := clients.NewOpenFGAClient(httpClient, cfg.LLMServiceAPIKey, cfg.LLMServiceURL)

	prService := &PRService{
		githubClient: *githubClient,
		llmClient:    *openFGAClient,
		cfg:          cfg,
	}

	return &Services{
		PRService: prService,
	}
}
