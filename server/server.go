package router

import (
	config "ai-api/config"
	handler "ai-api/handlers" // Import the handler package
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	Config *config.Config
	Router *gin.Engine
}

// SetupRouter sets up all routes for the application
func NewServer(cfg *config.Config) Server {

	logger := setupLogger()
	r := gin.Default()
	r.Use(ZlogMiddleware(logger))
	r.SetTrustedProxies([]string{})

	r = addRoutes(r)
	return Server{
		Config: cfg,
		Router: r,
	}
}

func setupLogger() zerolog.Logger {
	// Create a new logger that writes to standard output
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func ZlogMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add the logger to the context
		c.Set("zlog", logger)
		// logger.Info().Msg("Hello from zlog!")

		// Proceed to the next middleware or handler
		c.Next()
	}
}

func addRoutes(r *gin.Engine) *gin.Engine {

	api := r.Group("/v1/api")
	{
		// testing
		api.GET("/hello", handler.HelloHandler)

		// PULL REQUEST ROUTES
		pr := api.Group("/pull-request")
		{
			pr.GET("/:id", handler.GetPR)

		}
	}

	return r
}
