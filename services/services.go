package services

import (
	"ai-api/config"
	"net/http"
)

type Services struct {
	PRService PRService
}

// NewServices creates a new Services instance
func NewServices(cfg config.Config) *Services {
	prService := &PRServiceImpl{
		httpClient: &http.Client{},
		cfg:        cfg,
	}
	return &Services{
		PRService: prService,
	}
}
