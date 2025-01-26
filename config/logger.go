package config

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger is a globally accessible structured logger
var Logger zerolog.Logger

func init() {
	// Create log file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	// Configure logger
	Logger = zerolog.New(logFile).
		With().
		Timestamp().
		Logger()

	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}