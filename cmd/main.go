package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/internal/routes"
)

// Handler function for Vercel to call
func Handler(w http.ResponseWriter, r *http.Request) {
	// Create a Gin router instance
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://web-analyzer-frontend-omega.vercel.app", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Load API routes
	routes.RegisterRoutes(router)

	// Serve the request using Gin
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

func main() {
	// Log to indicate that the app is ready
	fmt.Println("App started!")
}
