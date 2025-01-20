package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/internal/handler"
	"github.com/uikee/web-analyzer-service/internal/service"
)

// RegisterRoutes sets up API endpoints
func RegisterRoutes(router *gin.Engine) {
	analyzerService := services.NewAnalyzerService()
	analyzerHandler := handler.NewAnalyzerHandler(analyzerService)

	router.GET("/analyze", analyzerHandler.AnalyzePage)
}