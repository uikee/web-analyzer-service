package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/config"
	"github.com/uikee/web-analyzer-service/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create a new Gin router
	router := gin.Default()

	// Load API routes
	routes.RegisterRoutes(router)

	// Start the server
	log.Printf("Server running on port %s\n", cfg.ServerPort)
	log.Fatal(router.Run(":" + cfg.ServerPort))
}
