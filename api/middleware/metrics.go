package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-hexagonal/util/log"
	"go-hexagonal/util/metrics"
)

// MetricsResponseWriter is a wrapper around http.ResponseWriter that tracks the status code
type MetricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewMetricsResponseWriter creates a new MetricsResponseWriter
func NewMetricsResponseWriter(w http.ResponseWriter) *MetricsResponseWriter {
	return &MetricsResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
	}
}

// WriteHeader captures the status code
func (w *MetricsResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Status returns the response status code
func (w *MetricsResponseWriter) Status() int {
	return w.statusCode
}

// Implements http.Hijacker if the underlying ResponseWriter does
func (w *MetricsResponseWriter) Hijack() (interface{}, interface{}, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("underlying ResponseWriter does not implement http.Hijacker")
}

// Flush implements http.Flusher if the underlying ResponseWriter does
func (w *MetricsResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// CloseNotify implements http.CloseNotifier if the underlying ResponseWriter does
func (w *MetricsResponseWriter) CloseNotify() <-chan bool {
	if cn, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return make(<-chan bool)
}

// RecordHTTPMetrics records HTTP metrics for a request
func RecordHTTPMetrics(handlerName, method string, statusCode int, duration time.Duration) {
	if !metrics.Initialized() {
		return
	}

	// Convert status code to string
	statusCodeStr := fmt.Sprintf("%d", statusCode)

	// Record request duration and count
	metrics.RequestDuration.WithLabelValues(handlerName, method, statusCodeStr).Observe(duration.Seconds())
	metrics.RequestTotal.WithLabelValues(handlerName, method, statusCodeStr).Inc()

	// Record errors if any
	if statusCode >= 400 {
		errorType := "client_error"
		if statusCode >= 500 {
			errorType = "server_error"
		}
		metrics.RecordError(errorType, handlerName)
		log.SugaredLogger.Debugf("HTTP %s error for %s %s (handler: %s): %d",
			errorType, method, handlerName, handlerName, statusCode)
	}
}

// MetricsMiddleware creates a middleware for collecting HTTP metrics
func MetricsMiddleware(handlerName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !metrics.Initialized() {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			metricsWriter := NewMetricsResponseWriter(w)

			// Process the request
			next.ServeHTTP(metricsWriter, r)

			// Record metrics after request is processed
			duration := time.Since(start)
			statusCode := metricsWriter.Status()

			// Use RecordHTTPMetrics to record the metrics
			RecordHTTPMetrics(handlerName, r.Method, statusCode, duration)
		})
	}
}

// InitializeMetrics initializes the metrics collection system
func InitializeMetrics() {
	metrics.Init()
	log.SugaredLogger.Info("Metrics collection system initialized")
}

// StartMetricsServer starts the metrics server
func StartMetricsServer(addr string) error {
	log.SugaredLogger.Infof("Starting metrics server on %s", addr)
	ctx := context.Background()
	return metrics.StartServer(ctx, addr)
}
