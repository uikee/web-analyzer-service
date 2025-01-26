package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the Logger to prevent actual logging during tests
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Warn() *MockLogger {
	m.Called()
	return m
}

func (m *MockLogger) Msg(msg string) {
	m.Called(msg)
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Set up the environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("FRONTEND_URL", "http://example.com")

	// Create a config using LoadConfig
	config := LoadConfig()

	// Assertions
	assert.Equal(t, "9090", config.ServerPort)
	assert.Equal(t, "http://example.com", config.FrontendURL)
}

func TestLoadConfig_WithFallbackValues(t *testing.T) {
	// Unset the environment variables to test the fallback
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("FRONTEND_URL")

	// Create a config using LoadConfig
	config := LoadConfig()

	// Assertions
	assert.Equal(t, "8081", config.ServerPort)       // default fallback
	assert.Equal(t, "http://localhost:3000", config.FrontendURL) // default fallback
}

func TestGetEnv_WithEnvVar(t *testing.T) {
	// Set environment variable for the test
	os.Setenv("MY_VAR", "test_value")

	// Test if getEnv returns the correct value
	value := getEnv("MY_VAR", "default_value")

	// Assertions
	assert.Equal(t, "test_value", value)
}

func TestGetEnv_WithFallback(t *testing.T) {
	// Unset the environment variable to test the fallback
	os.Unsetenv("MY_VAR")

	// Test if getEnv returns the fallback value
	value := getEnv("MY_VAR", "default_value")

	// Assertions
	assert.Equal(t, "default_value", value)
}
