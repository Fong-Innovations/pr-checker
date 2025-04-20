package services

import (
	"log"
	"net/http"
	clients "pr-checker/clients"
	"pr-checker/config"
	"pr-checker/repository"
	"time"
)

type Services struct {
	PRService             *PRService
	PullRequestRepository *repository.PRDataRepository
}

// NewServices creates a new Services instance
func NewServices(cfg config.Config) *Services {

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	mariadDBClient, err := clients.NewMariaDBClient(cfg.DBUser, cfg.DBPass, cfg.DBUrl, cfg.DBPort, cfg.DBName)
	if err != nil {
		log.Printf("Error creating database client:: %v", err)
	}

	prDataRepository := repository.NewPRDataRepository(mariadDBClient.DB)
	githubClient := clients.NewGithubClient(httpClient, cfg.GithubToken, cfg.GithubBaseURL)
	OpenAIClient := clients.NewOpenAIClient(httpClient, cfg.LLMServiceAPIKey, cfg.LLMServiceURL)

	prService := &PRService{
		githubClient:          *githubClient,
		llmClient:             *OpenAIClient,
		PullRequestRepository: *prDataRepository,
		cfg:                   cfg,
	}

	return &Services{
		PRService: prService,
	}
}
