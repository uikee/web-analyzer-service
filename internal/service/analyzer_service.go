package services

import (
	"errors"
	"io"
	"net/http"

	"github.com/uikee/web-analyzer-service/internal/utils"
)

// AnalysisResult represents the result of a web analysis
type AnalysisResult struct {
	Title       string
	HTMLVersion string
	Err         error
}

// AnalyzerService provides functionality to analyze web pages
type AnalyzerService interface {
	Analyze(url string) AnalysisResult
}

// analyzerServiceImpl is the concrete implementation of AnalyzerService
type analyzerServiceImpl struct{}

// NewAnalyzerService creates a new instance of AnalyzerService
func NewAnalyzerService() AnalyzerService {
	return &analyzerServiceImpl{}
}

// Analyze fetches the webpage and extracts title and HTML version
func (s *analyzerServiceImpl) Analyze(url string) AnalysisResult {
	resp, err := http.Get(url)
	if err != nil {
		return AnalysisResult{Err: errors.New("failed to fetch the URL")}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AnalysisResult{Err: errors.New("URL returned non-200 status")}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AnalysisResult{Err: errors.New("failed to read response body")}
	}

	htmlContent := string(body)

	// Use utility functions for title and HTML version extraction
	title := utils.ExtractTitle(htmlContent)
	htmlVersion := utils.DetectHTMLVersion(htmlContent)

	if title == "" {
		title = "Title not found"
	}
	if htmlVersion == "" {
		htmlVersion = "Unknown HTML version"
	}

	return AnalysisResult{
		Title:       title,
		HTMLVersion: htmlVersion,
		Err:         nil,
	}
}