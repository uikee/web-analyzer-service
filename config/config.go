package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	ServerPort string
	Env        string
}

// LoadConfig loads environment variables from .env file
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env variables")
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
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