package router

import (
	"os"
	config "pr-checker/config"
	"pr-checker/handlers"
	handler "pr-checker/handlers" // Import the handler package
	"pr-checker/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	Config    *config.Config
	PRHandler *handler.PRHandler
	Router    *gin.Engine
}

// SetupRouter sets up all routes for the application
func NewServer(cfg *config.Config, services *services.Services) Server {

	logger := setupLogger()

	r := gin.Default()

	// create handlers
	prHandler := handlers.NewPRHandler(services.PRService)

	r.Use(ZlogMiddleware(logger))
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
	// Create a new logger that writes to standard output
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func ZlogMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("zlog", logger)
		c.Next()
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
		}
	}
}
