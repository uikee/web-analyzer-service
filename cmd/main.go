package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/config"
	"github.com/uikee/web-analyzer-service/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	logger := config.Logger

	// Create a new Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Load API routes
	routes.RegisterRoutes(router)

	// Graceful shutdown handling
	go func() {
		logger.Info().Msgf("Server running on port %s", cfg.ServerPort)
		if err := router.Run(":" + cfg.ServerPort); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")
}