package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	ServerPort  string
	FrontendUrl string
	Env         string
}

// LoadConfig loads environment variables from .env file
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env variables")
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8081"),
		FrontendUrl: getEnv("FRONTEND_URL", "http://localhost:3000"),
		Env:        getEnv("ENV", "dev"),
	}
}

// getEnv fetches an env variable with a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}