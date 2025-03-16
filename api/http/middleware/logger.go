package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-hexagonal/util/log"
)

// bodyLogWriter is a custom response writer that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLogger is a middleware that logs request and response details
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a custom response writer to capture the response
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		log.SugaredLogger.With(
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("request_id", c.GetHeader("X-Request-ID")),
			zap.String("user_agent", c.Request.UserAgent()),
		).Infof("Request completed")

		// Log request and response bodies in debug mode
		if log.SugaredLogger.Desugar().Core().Enabled(zap.DebugLevel) {
			// Only log bodies for non-binary content types
			contentType := c.Writer.Header().Get("Content-Type")
			if isTextContentType(contentType) {
				log.SugaredLogger.Debugf("Request body: %s", string(requestBody))
				log.SugaredLogger.Debugf("Response body: %s", blw.body.String())
			} else {
				log.SugaredLogger.Debug("Binary content type - body not logged")
			}
		}
	}
}

// isTextContentType determines if the content type is text-based
func isTextContentType(contentType string) bool {
	switch {
	case contentType == "":
		return true // Default to text if not specified
	case bytes.Contains([]byte(contentType), []byte("text/")):
		return true
	case bytes.Contains([]byte(contentType), []byte("application/json")):
		return true
	case bytes.Contains([]byte(contentType), []byte("application/xml")):
		return true
	case bytes.Contains([]byte(contentType), []byte("application/x-www-form-urlencoded")):
		return true
	default:
		return false
	}
}
