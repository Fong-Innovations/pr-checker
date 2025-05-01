package utils

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger() {
	// Set global log level (e.g., info, debug, error)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Initialize the global logger
	Log = zerolog.New(os.Stdout).With().
		Timestamp().                  // Add timestamps to logs
		Str("service", "pr-checker"). // Add a service name field
		Logger()
}
