package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-hexagonal/util/log"
)

// ResponseWriter is a wrapper for gin.ResponseWriter that captures the
// response status code and size
type ResponseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

// Write captures the response body and writes it to the underlying writer
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Status returns the status code
func (rw *ResponseWriter) Status() int {
	return rw.statusCode
}

// LoggingConfig holds configuration for the request logging middleware
type LoggingConfig struct {
	// Whether to log request body (disabled by default for privacy and size reasons)
	LogRequestBody bool
	// Whether to log response body (disabled by default for privacy and size reasons)
	LogResponseBody bool
	// Maximum size of request/response body to log
	MaxBodyLogSize int
	// Skip logging for specified paths
	SkipPaths []string
}

// DefaultLoggingConfig returns the default logging configuration
func DefaultLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		LogRequestBody:  false,
		LogResponseBody: false,
		MaxBodyLogSize:  1024, // 1 KB
		SkipPaths:       []string{"/ping", "/health"},
	}
}

// RequestLogger returns a middleware that logs incoming requests and outgoing responses
func RequestLogger() gin.HandlerFunc {
	return RequestLoggerWithConfig(DefaultLoggingConfig())
}

// RequestLoggerWithConfig returns a middleware that logs requests and responses with custom config
func RequestLoggerWithConfig(config *LoggingConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultLoggingConfig()
	}

	return func(c *gin.Context) {
		// Skip logging for certain paths
		for _, path := range config.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// Get request ID from context
		var requestID string
		if id, exists := c.Get(RequestIDHeader); exists {
			requestID = id.(string)
		}

		// Start timer
		start := time.Now()

		// Read request body if enabled
		var requestBody []byte
		if config.LogRequestBody && c.Request.Body != nil {
			var err error
			requestBody, err = io.ReadAll(c.Request.Body)
			if err != nil {
				log.Logger.Error("Failed to read request body", zap.Error(err))
			}

			// Restore request body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a response writer wrapper to capture the response
		responseWriter := &ResponseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			statusCode:     http.StatusOK, // Default status is 200
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Check if request path is an API path
		isAPIPath := len(c.Errors) == 0 && c.Writer.Status() < 500

		// Determine log level based on status code
		var logMethod func(string, ...zap.Field)
		if c.Writer.Status() >= 500 {
			logMethod = log.Logger.Error
		} else if c.Writer.Status() >= 400 {
			logMethod = log.Logger.Warn
		} else {
			logMethod = log.Logger.Info
		}

		// Create log fields
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Int("status", responseWriter.Status()),
			zap.String("latency", duration.String()),
			zap.Int64("latency_ms", duration.Milliseconds()),
			zap.String("request_id", requestID),
		}

		// Add request body if enabled and present
		if config.LogRequestBody && len(requestBody) > 0 {
			// Limit the size of the logged body
			if len(requestBody) > config.MaxBodyLogSize {
				requestBody = requestBody[:config.MaxBodyLogSize]
			}
			fields = append(fields, zap.ByteString("request_body", requestBody))
		}

		// Add response body if enabled
		if config.LogResponseBody && responseWriter.body.Len() > 0 {
			responseBody := responseWriter.body.Bytes()
			// Limit the size of the logged body
			if len(responseBody) > config.MaxBodyLogSize {
				responseBody = responseBody[:config.MaxBodyLogSize]
			}
			fields = append(fields, zap.ByteString("response_body", responseBody))
		}

		// Add error if present
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// Log the request
		message := "Request"
		if isAPIPath {
			message = "API Request"
		}
		logMethod(message, fields...)
	}
}
