package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uikee/web-analyzer-service/internal/services"
)

// MockAnalyzerService mocks the AnalyzerService interface
type MockAnalyzerService struct {
	mock.Mock
}

func (m *MockAnalyzerService) Analyze(url string) (services.AnalysisResult, error) {
	args := m.Called(url)
	return args.Get(0).(services.AnalysisResult), args.Error(1)
}

// MockURLValidator mocks the URLValidator interface
type MockURLValidator struct {
	mock.Mock
}

func (m *MockURLValidator) Validate(url string) error {
	args := m.Called(url)
	return args.Error(0)
}

func TestAnalyzePage_Success(t *testing.T) {
	// Create a gin context
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Create mock services
	mockAnalyzerService := new(MockAnalyzerService)
	mockValidator := new(MockURLValidator)

	// Create handler
	handler := NewAnalyzerHandler(mockAnalyzerService, mockValidator)

	// Define mock behavior
	mockValidator.On("Validate", "http://example.com").Return(nil)
	mockAnalyzerService.On("Analyze", "http://example.com").Return(services.AnalysisResult{}, nil)

	// Test case: Valid URL
	r.GET("/analyze", handler.AnalyzePage)
	w := performRequest(r, "GET", "/analyze?url=http://example.com")

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify that the methods were called
	mockValidator.AssertExpectations(t)
	mockAnalyzerService.AssertExpectations(t)
}

func TestAnalyzePage_MissingURL(t *testing.T) {
	// Create mock services
	mockAnalyzerService := new(MockAnalyzerService)
	mockValidator := new(MockURLValidator)

	// Create handler
	handler := NewAnalyzerHandler(mockAnalyzerService, mockValidator)

	// Test case: Missing URL
	r := gin.Default()
	r.GET("/analyze", handler.AnalyzePage)
	w := performRequest(r, "GET", "/analyze")

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "URL parameter is required")
}

func TestAnalyzePage_InvalidURL(t *testing.T) {
	// Create mock services
	mockAnalyzerService := new(MockAnalyzerService)
	mockValidator := new(MockURLValidator)

	// Create handler
	handler := NewAnalyzerHandler(mockAnalyzerService, mockValidator)

	// Define mock behavior
	mockValidator.On("Validate", "invalid-url").Return(errors.New("Invalid URL"))

	// Test case: Invalid URL format
	r := gin.Default()
	r.GET("/analyze", handler.AnalyzePage)
	w := performRequest(r, "GET", "/analyze?url=invalid-url")

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid URL")
}

func TestAnalyzePage_AnalyzerError(t *testing.T) {
	// Create mock services
	mockAnalyzerService := new(MockAnalyzerService)
	mockValidator := new(MockURLValidator)

	// Create handler
	handler := NewAnalyzerHandler(mockAnalyzerService, mockValidator)

	// Define mock behavior
	mockValidator.On("Validate", "http://example.com").Return(nil)
	mockAnalyzerService.On("Analyze", "http://example.com").Return(services.AnalysisResult{}, errors.New("Error during page analysis"))

	// Test case: Error during analysis
	r := gin.Default()
	r.GET("/analyze", handler.AnalyzePage)
	w := performRequest(r, "GET", "/analyze?url=http://example.com")

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error during page analysis")
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
