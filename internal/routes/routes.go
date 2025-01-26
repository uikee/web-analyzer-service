package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/config"
	"github.com/uikee/web-analyzer-service/internal/handler"
	"github.com/uikee/web-analyzer-service/internal/services"
	"github.com/uikee/web-analyzer-service/internal/validators"
)

// RegisterRoutes sets up API endpoints
func RegisterRoutes(router *gin.Engine) {
	// Attempt to initialize the analyzer service
	analyzerService := services.NewAnalyzerService()

	// Initialize the URL validator
	urlValidator := validators.NewURLValidator()

	// Log successful initialization of the service and validator
	config.Logger.Info().Msg("Analyzer service and URL validator initialized successfully")

	// Create the handler instance with both the service and validator
	analyzerHandler := handler.NewAnalyzerHandler(analyzerService, urlValidator)

	// Register the /analyze route and log the registration
	router.GET("/analyze", func(c *gin.Context) {
		// Log request for analysis
		config.Logger.Info().Msg("Received request for /analyze endpoint")
		analyzerHandler.AnalyzePage(c)
	})

	// Log successful route registration
	config.Logger.Info().Msg("Routes registered successfully")
}
