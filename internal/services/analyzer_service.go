package services

import (
	"errors"
	"io"
	"net/http"

	"github.com/uikee/web-analyzer-service/config"
	"github.com/uikee/web-analyzer-service/internal/utils"
)

// AnalysisResult represents the result of a web analysis
type AnalysisResult struct {
	Title             string         `json:"title"`
	HTMLVersion       string         `json:"html_version"`
	Headings          map[string]int `json:"headings"`
	InternalLinks     int            `json:"internal_links"`
	ExternalLinks     int            `json:"external_links"`
	InaccessibleLinks int            `json:"inaccessible_links"`
	HasLoginForm      bool           `json:"has_login_form"`
}

// AnalyzerService provides functionality to analyze web pages
type AnalyzerService interface {
	Analyze(url string) (AnalysisResult, error)
}

// analyzerServiceImpl is the concrete implementation of AnalyzerService
type analyzerServiceImpl struct{}

// NewAnalyzerService creates a new instance of AnalyzerService
func NewAnalyzerService() AnalyzerService {
	return &analyzerServiceImpl{}
}

var (
	// ErrFetchFailed indicates that fetching the URL failed
	ErrFetchFailed = errors.New("failed to fetch the URL")

	// ErrReadBodyFailed indicates that reading the response body failed
	ErrReadBodyFailed = errors.New("failed to read response body")
)

// Analyze fetches the webpage and extracts analysis data
func (s *analyzerServiceImpl) Analyze(targetURL string) (AnalysisResult, error) {
	config.Logger.Info().Str("url", targetURL).Msg("Fetching the web page")

	resp, err := http.Get(targetURL)
	if err != nil {
		return AnalysisResult{}, ErrFetchFailed
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AnalysisResult{}, ErrReadBodyFailed
	}

	htmlContent := string(body)

	// Concurrent execution using channels
	headingsChan := make(chan map[string]int)
	loginFormChan := make(chan bool)
	linkCountsChan := make(chan [3]int)
	errorChan := make(chan error)

	go func() { headingsChan <- utils.CountHeadings(htmlContent) }()
	go func() { loginFormChan <- utils.ContainsLoginForm(htmlContent) }()
	go func() {
		internal, external, inaccessible, err := utils.CountLinksConcurrently(targetURL, htmlContent)
		if err != nil {
			errorChan <- err
		} else {
			linkCountsChan <- [3]int{internal, external, inaccessible}
		}
	}()

	// Get results
	headings := <-headingsChan
	hasLoginForm := <-loginFormChan

	var linkCounts [3]int
	select {
	case linkCounts = <-linkCountsChan:
	case err = <-errorChan:
		return AnalysisResult{}, err
	}

	return AnalysisResult{
		Title:             utils.ExtractTitle(htmlContent),
		HTMLVersion:       utils.DetectHTMLVersion(htmlContent),
		Headings:          headings,
		InternalLinks:     linkCounts[0],
		ExternalLinks:     linkCounts[1],
		InaccessibleLinks: linkCounts[2],
		HasLoginForm:      hasLoginForm,
	}, nil
}
