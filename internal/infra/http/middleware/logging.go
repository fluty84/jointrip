package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware provides request logging middleware
type LoggingMiddleware struct {
	logger *logrus.Logger
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger *logrus.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// RequestLogger middleware that logs HTTP requests
func (m *LoggingMiddleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get user agent
		userAgent := c.GetHeader("User-Agent")

		// Get user ID if available
		userID := ""
		if userIDInterface, exists := c.Get("user_id"); exists {
			if uid, ok := userIDInterface.(uuid.UUID); ok {
				userID = uid.String()
			}
		}

		// Create log entry
		logEntry := m.logger.WithFields(logrus.Fields{
			"request_id":   requestID,
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"status_code":  statusCode,
			"latency_ms":   latency.Milliseconds(),
			"client_ip":    clientIP,
			"user_agent":   userAgent,
			"content_type": c.GetHeader("Content-Type"),
		})

		// Add user ID if available
		if userID != "" {
			logEntry = logEntry.WithField("user_id", userID)
		}

		// Log based on status code
		switch {
		case statusCode >= 500:
			logEntry.Error("HTTP request completed with server error")
		case statusCode >= 400:
			logEntry.Warn("HTTP request completed with client error")
		case statusCode >= 300:
			logEntry.Info("HTTP request completed with redirect")
		default:
			logEntry.Info("HTTP request completed successfully")
		}
	}
}

// ErrorLogger middleware that logs errors
func (m *LoggingMiddleware) ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log any errors that occurred during request processing
		for _, err := range c.Errors {
			requestID := ""
			if reqID, exists := c.Get("request_id"); exists {
				if id, ok := reqID.(string); ok {
					requestID = id
				}
			}

			logEntry := m.logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"error_type": err.Type,
			})

			switch err.Type {
			case gin.ErrorTypePublic:
				logEntry.WithError(err.Err).Warn("Public error occurred")
			case gin.ErrorTypeBind:
				logEntry.WithError(err.Err).Warn("Binding error occurred")
			case gin.ErrorTypeRender:
				logEntry.WithError(err.Err).Error("Render error occurred")
			default:
				logEntry.WithError(err.Err).Error("Internal error occurred")
			}
		}
	}
}
