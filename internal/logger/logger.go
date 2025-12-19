package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Log zerolog.Logger

func Init() {
	// Enable stack traces for errors
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339

	// Determine environment
	isDevelopment := os.Getenv("APP_ENV") == "development"

	// Configure output format
	if isDevelopment {
		// Pretty console output for development
		Log = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).With().Timestamp().Caller().Logger()
	} else {
		// JSON output for production
		Log = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	}

	// Set log level (default: info)
	logLevel := zerolog.InfoLevel
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		parsedLevel, err := zerolog.ParseLevel(level)
		if err == nil {
			logLevel = parsedLevel
		}
	}

	Log = Log.Level(logLevel)
}
