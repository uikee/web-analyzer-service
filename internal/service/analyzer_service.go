package services

import (
	"errors"
	"io"
	"net/http"

	// "sync"

	"github.com/uikee/web-analyzer-service/internal/utils"
)

// AnalysisResult represents the result of a web analysis
type AnalysisResult struct {
	Title             string
	HTMLVersion       string
	Headings          map[string]int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	HasLoginForm      bool
	Err               error
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

	// Channels for concurrent tasks
	headingsChan := make(chan map[string]int)
	loginFormChan := make(chan bool)
	linkCountsChan := make(chan [3]int)
	errorChan := make(chan error)

	// Goroutines for concurrent tasks
	go func() { headingsChan <- utils.CountHeadings(htmlContent) }()
	go func() { loginFormChan <- utils.ContainsLoginForm(htmlContent) }()
	go func() {
		internal, external, inaccessible, err := utils.CountLinksConcurrently(url, htmlContent)
		if err != nil {
			errorChan <- err
		} else {
			linkCountsChan <- [3]int{internal, external, inaccessible}
		}
	}()

	// Detect title and HTML version
	title := utils.ExtractTitle(htmlContent)
	htmlVersion := utils.DetectHTMLVersion(htmlContent)

	// Wait for results
	headings := <-headingsChan
	hasLoginForm := <-loginFormChan

	var linkCounts [3]int
	select {
	case linkCounts = <-linkCountsChan:
	case err = <-errorChan:
		return AnalysisResult{Err: err}
	}

	if title == "" {
		title = "Title not found"
	}
	if htmlVersion == "" {
		htmlVersion = "Unknown HTML version"
	}

	return AnalysisResult{
		Title:             title,
		HTMLVersion:       htmlVersion,
		Headings:          headings,
		InternalLinks:     linkCounts[0],
		ExternalLinks:     linkCounts[1],
		InaccessibleLinks: linkCounts[2],
		HasLoginForm:      hasLoginForm,
		Err:               nil,
	}
}
