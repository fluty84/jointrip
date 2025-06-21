package logger

import (
	"os"

	"jointrip/internal/infra/config"

	"github.com/sirupsen/logrus"
)

// NewLogger creates a new configured logger instance
func NewLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		logger.Warn("Invalid log level, defaulting to info")
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set log format
	switch cfg.Log.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
			},
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		// Default to JSON in production, text in development
		if cfg.IsProduction() {
			logger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			})
		} else {
			logger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
				ForceColors:     true,
			})
		}
	}

	// Set output
	logger.SetOutput(os.Stdout)

	// Add common fields
	logger = logger.WithFields(logrus.Fields{
		"service": "jointrip-api",
		"version": "1.0.0",
		"env":     cfg.Server.Env,
	}).Logger

	return logger
}

// WithRequestID adds a request ID to the logger context
func WithRequestID(logger *logrus.Logger, requestID string) *logrus.Entry {
	return logger.WithField("request_id", requestID)
}

// WithUserID adds a user ID to the logger context
func WithUserID(logger *logrus.Logger, userID string) *logrus.Entry {
	return logger.WithField("user_id", userID)
}

// WithError adds an error to the logger context
func WithError(logger *logrus.Logger, err error) *logrus.Entry {
	return logger.WithError(err)
}
