package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/internal/service"
)

// AnalyzerHandler provides HTTP handlers for web analysis
type AnalyzerHandler struct {
	analyzerService services.AnalyzerService
}

// NewAnalyzerHandler creates a new instance of AnalyzerHandler
func NewAnalyzerHandler(service services.AnalyzerService) *AnalyzerHandler {
	return &AnalyzerHandler{analyzerService: service}
}

// AnalyzePage handles requests for analyzing web pages
func (h *AnalyzerHandler) AnalyzePage(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	// Execute analysis in a separate goroutine
	resultChan := make(chan services.AnalysisResult)
	go func() {
		resultChan <- h.analyzerService.Analyze(url)
	}()
	
	// Retrieve results
	result := <-resultChan
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"title":       result.Title,
		"htmlVersion": result.HTMLVersion,
	})
}