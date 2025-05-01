package router

import (
	"os"
	config "pr-checker/config"
	"pr-checker/handlers"
	"pr-checker/services"
	zlog "pr-checker/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	Config    *config.Config
	PRHandler *handlers.PRHandler
	Router    *gin.Engine
}

// SetupRouter sets up all routes for the application
func NewServer(cfg *config.Config, services *services.Services) Server {

	zlog.InitLogger()
	r := gin.Default()

	// create handlers
	prHandler := handlers.NewPRHandler(services.PRService)
	r.SetTrustedProxies([]string{})

	// Register routes
	server := &Server{
		Config:    cfg,
		Router:    r,
		PRHandler: prHandler,
	}

	server.routes()

	return *server
}

func setupLogger() zerolog.Logger {
	// Set global log level (e.g., info, debug, error)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create a new logger that writes to standard output
	return zerolog.New(os.Stdout).With().
		Timestamp().                  // Add timestamps to logs
		Str("service", "pr-checker"). // Add a service name field
		Logger()
}

func ZlogMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Attach the logger to the context
		c.Set("zlog", logger)

		start := time.Now()

		// Process the request
		c.Next()

		// Log request details after processing
		logger.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Msg("HTTP request")
	}
}
func (s *Server) routes() {

	api := s.Router.Group("/v1/api")
	{
		// PULL REQUEST ROUTES
		pr := api.Group("/pr")
		{
			pr.PUT("store/:owner/:repo/:id", s.PRHandler.StorePRData)
			pr.GET("changes/:owner/:repo/:id", s.PRHandler.AnalyzePR)
			pr.GET("test", s.PRHandler.TestFunction)
		}
	}
}
