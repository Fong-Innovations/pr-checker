// File: config/config.go
// Loads config variables from env file
package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

// Config struct to hold application configuration
type Config struct {
	GithubToken      string `koanf:"github_token"`
	GithubBaseURL    string `koanf:"github_base_url"`
	LLMServiceURL    string `koanf:"llm_base_url"`
	LLMServiceAPIKey string `koanf:"llm_api_key"`
	LLMModel         string `koanf:"llm_model"`
	LLMAnalyzePrompt string `koanf:"llm_analyze_pr_prompt"`
}

// LoadConfig reads configuration from a .env file and environment variables
func LoadConfig(envFile string) (*Config, error) {
	// Load the .env file into environment variables
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("failed to load env file: %v", err)
		return nil, err
	}

	var k = koanf.New(".")

	// Load environment variables with the prefix "AICHECKER_"
	err = k.Load(env.Provider("AI_CHECKER_", ".", func(s string) string {
		// Transform environment variable names to match struct field names
		return strings.ToLower(strings.TrimPrefix(s, "AI_CHECKER_"))
	}), nil)
	if err != nil {
		log.Printf("failed to load environment variables: %v", err)
		return nil, err
	}
	// Create and populate a Config struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Printf("error unmarshaling config: %v", err)
		return nil, err
	}
	return &cfg, nil
}
