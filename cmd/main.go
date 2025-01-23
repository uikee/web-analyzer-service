package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/config"
	"github.com/uikee/web-analyzer-service/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create a new Gin router
	router := gin.Default()

	// Configure CORS
    allowedOrigins := []string{"http://localhost:3000"}
    if envOrigin := cfg.FrontendUrl; envOrigin != "" {
        allowedOrigins = append(allowedOrigins, envOrigin)
    }

    router.Use(cors.New(cors.Config{
        AllowOrigins:     allowedOrigins,
        AllowMethods:     []string{"GET"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
    }))

	// Load API routes
	routes.RegisterRoutes(router)

	// Start the server
	log.Printf("Server running on port %s\n", cfg.ServerPort)
	log.Fatal(router.Run(":" + cfg.ServerPort))
}
