package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/internal/routes"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Load API routes
	routes.SetupRoutes(router)

	// Start the server
	port := "8080"
	log.Printf("Server running on port %s\n", port)
	log.Fatal(router.Run(":" + port))
}
