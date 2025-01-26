package services_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uikee/web-analyzer-service/internal/services"
)

// MockUtils is a mock implementation of utility functions
type MockUtils struct {
	mock.Mock
}

func (m *MockUtils) CountHeadings(htmlContent string) map[string]int {
	args := m.Called(htmlContent)
	return args.Get(0).(map[string]int)
}

func (m *MockUtils) ContainsLoginForm(htmlContent string) bool {
	args := m.Called(htmlContent)
	return args.Bool(0)
}

func (m *MockUtils) CountLinksConcurrently(baseURL, htmlContent string) (int, int, int, error) {
	args := m.Called(baseURL, htmlContent)
	return args.Int(0), args.Int(1), args.Int(2), args.Error(3)
}

func (m *MockUtils) ExtractTitle(htmlContent string) string {
	args := m.Called(htmlContent)
	return args.String(0)
}

func (m *MockUtils) DetectHTMLVersion(htmlContent string) string {
	args := m.Called(htmlContent)
	return args.String(0)
}

func TestAnalyze_Success(t *testing.T) {
	mockHTMLContent := "<html><head><title>Test Page</title></head><body><h1>Heading 1</h1><a href=\"http://example.com\">Link</a></body></html>"

	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockHTMLContent))
	}))
	defer server.Close()

	// Mock utility functions
	mockUtils := services.UtilityFunctions{
		CountHeadings: func(htmlContent string) map[string]int {
			return map[string]int{"h1": 1}
		},
		ContainsLoginForm: func(htmlContent string) bool {
			return false
		},
		CountLinksConcurrently: func(baseURL, htmlContent string) (int, int, int, error) {
			return 1, 0, 0, nil
		},
		ExtractTitle: func(htmlContent string) string {
			return "Test Page"
		},
		DetectHTMLVersion: func(htmlContent string) string {
			return "HTML5"
		},
	}

	// Create the service with mock utilities
	service := services.NewAnalyzerServiceWithUtils(mockUtils)

	// Call the Analyze method
	result, err := service.Analyze(server.URL)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Test Page", result.Title)
	assert.Equal(t, "HTML5", result.HTMLVersion)
	assert.Equal(t, map[string]int{"h1": 1}, result.Headings)
	assert.Equal(t, 1, result.InternalLinks)
	assert.Equal(t, 0, result.ExternalLinks)
	assert.Equal(t, 0, result.InaccessibleLinks)
	assert.False(t, result.HasLoginForm)
}

func TestAnalyze_FetchFailed(t *testing.T) {
	service := services.NewAnalyzerService()

	// Call the Analyze method with an invalid URL
	_, err := service.Analyze("http://invalid-url")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, services.ErrFetchFailed, err)
}

func TestAnalyze_ReadBodyFailed(t *testing.T) {
	// Mock HTTP server that returns an error for reading the body
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1") // Simulate a mismatch to cause an error
	}))
	defer server.Close()

	service := services.NewAnalyzerService()

	// Call the Analyze method
	_, err := service.Analyze(server.URL)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, services.ErrReadBodyFailed, err)
}

// func TestAnalyze_UtilityError(t *testing.T) {
// 	// Mock HTTP server
// 	mockHTMLContent := "<html><head><title>Test Page</title></head><body><h1>Heading 1</h1><a href=\"http://example.com\">Link</a></body></html>"
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte(mockHTMLContent))
// 	}))
// 	defer server.Close()

// 	// Mock utility functions
// 	mockUtils := &MockUtils{}
// 	mockUtils.On("CountHeadings", mockHTMLContent).Return(map[string]int{"h1": 1})
// 	mockUtils.On("ContainsLoginForm", mockHTMLContent).Return(false)
// 	mockUtils.On("CountLinksConcurrently", server.URL, mockHTMLContent).Return(0, 0, 0, errors.New("utility error"))

// 	// Replace utils with mockUtils
// 	services.SetUtils(mockUtils)

// 	// Create the service
// 	service := services.NewAnalyzerService()

// 	// Call the Analyze method
// 	_, err := service.Analyze(server.URL)

// 	// Assertions
// 	assert.Error(t, err)
// 	assert.Equal(t, "utility error", err.Error())
// }
