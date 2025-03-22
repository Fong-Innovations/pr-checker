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
	APIKey      string `koanf:"api_key"`
	GithubToken string `koanf:"github_token"`
}

var k = koanf.New(".")

// LoadConfig reads configuration from a .env file and environment variables
func LoadConfig(envFile string) *Config {
	// Load the .env file into environment variables
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No .env file found or failed to load: %v", err)
	}

	// Load environment variables with the prefix "AICHECKER_"
	err := k.Load(env.Provider("AICHECKER_", ".", func(s string) string {
		// Transform environment variable names to match struct field names
		return strings.ToLower(strings.TrimPrefix(s, "AICHECKER_"))
	}), nil)
	if err != nil {
		log.Fatalf("error loading environment variables: %v", err)
	}
	// Create and populate a Config struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Printf("error unmarshaling config: %v", err)
	}

	return &cfg
}
