package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/config"
	"github.com/uikee/web-analyzer-service/internal/services"
	"github.com/uikee/web-analyzer-service/internal/validators"
)

// ErrorResponse represents a standardized error response with HTTP status
type ErrorResponse struct {
	HTTPStatus int    `json:"status"`
	Message    string `json:"error"`
}

// AnalyzerHandler provides HTTP handlers for web analysis
type AnalyzerHandler struct {
	analyzerService services.AnalyzerService
	validator       validators.URLValidator
}

// NewAnalyzerHandler creates a new instance of AnalyzerHandler
func NewAnalyzerHandler(service services.AnalyzerService, validator validators.URLValidator) *AnalyzerHandler {
	return &AnalyzerHandler{analyzerService: service, validator: validator}
}

// handleError sends an appropriate JSON error response and logs it
func (h *AnalyzerHandler) handleError(c *gin.Context, statusCode int, err error, context string) {
	config.Logger.Error().
		Err(err).
		Int("status", statusCode).
		Str("context", context).
		Msg("API error occurred")

	c.JSON(statusCode, ErrorResponse{
		HTTPStatus: statusCode,
		Message:    err.Error(),
	})
}

// AnalyzePage handles requests for analyzing web pages
func (h *AnalyzerHandler) AnalyzePage(c *gin.Context) {
	urlParam := c.Query("url")
	if urlParam == "" {
		h.handleError(c, http.StatusBadRequest, validators.ErrMissingURL, "Missing URL parameter")
		return
	}

	// Validate URL
	if err := h.validator.Validate(urlParam); err != nil {
		h.handleError(c, http.StatusBadRequest, err, "Invalid URL format")
		return
	}

	config.Logger.Info().Str("url", urlParam).Msg("Start analyzing web page")

	// Perform the web page analysis
	result, err := h.analyzerService.Analyze(urlParam)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err, "Error during page analysis")
		return
	}

	config.Logger.Info().Str("url", urlParam).Msg("Web page analysis completed successfully")

	c.JSON(http.StatusOK, result)
}