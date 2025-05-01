package services

import (
	"net/http"
	clients "pr-checker/clients"
	"pr-checker/config"
	"pr-checker/repository"
	zlog "pr-checker/utils"
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

	mariadDBClient, err := clients.GetMariaDBClient(cfg.DBUser, cfg.DBPass, cfg.DBUrl, cfg.DBPort, cfg.DBName)
	if err != nil {
		zlog.Log.Error().Err(err).Msgf("error creating database client")
		return nil
	}

	prDataRepository := repository.NewPRDataRepository(mariadDBClient.DB)
	githubClient := clients.NewGithubClient(httpClient, cfg.GithubToken, cfg.GithubAPIVersion, cfg.GithubBaseURL)
	openAIClient := clients.NewOpenAIClient(httpClient, cfg.LLMServiceAPIKey, cfg.LLMServiceURL)

	prService := &PRService{
		githubClient:          *githubClient,
		llmClient:             *openAIClient,
		PullRequestRepository: *prDataRepository,
		cfg:                   cfg,
	}

	return &Services{
		PRService: prService,
	}
}
