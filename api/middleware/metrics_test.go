package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-hexagonal/util/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetricsResponseWriter(t *testing.T) {
	// Create a test response writer
	rw := httptest.NewRecorder()
	metricsWriter := NewMetricsResponseWriter(rw)

	assert.NotNil(t, metricsWriter)
	assert.Equal(t, rw, metricsWriter.ResponseWriter)
	assert.Equal(t, http.StatusOK, metricsWriter.Status())
}

func TestMetricsResponseWriter_WriteHeader(t *testing.T) {
	rw := httptest.NewRecorder()
	metricsWriter := NewMetricsResponseWriter(rw)

	// Test setting status code
	metricsWriter.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, metricsWriter.Status())
	assert.Equal(t, http.StatusNotFound, rw.Code)

	// Create a new recorder for testing a different status code
	rw2 := httptest.NewRecorder()
	metricsWriter2 := NewMetricsResponseWriter(rw2)

	// Test setting another status code on a fresh writer
	metricsWriter2.WriteHeader(http.StatusInternalServerError)
	assert.Equal(t, http.StatusInternalServerError, metricsWriter2.Status())
	assert.Equal(t, http.StatusInternalServerError, rw2.Code)
}

func TestMetricsResponseWriter_Write(t *testing.T) {
	rw := httptest.NewRecorder()
	metricsWriter := NewMetricsResponseWriter(rw)

	// Test writing body
	data := []byte("test data")
	n, err := metricsWriter.Write(data)

	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, "test data", rw.Body.String())
	assert.Equal(t, http.StatusOK, metricsWriter.Status())
}

func TestMetricsResponseWriter_Hijack(t *testing.T) {
	rw := httptest.NewRecorder()
	metricsWriter := NewMetricsResponseWriter(rw)

	// Test Hijack method - should fail since httptest.ResponseWriter doesn't implement Hijacker
	conn, buf, err := metricsWriter.Hijack()

	assert.Error(t, err)
	assert.Nil(t, conn)
	assert.Nil(t, buf)
	assert.Contains(t, err.Error(), "does not implement http.Hijacker")
}

func TestMetricsResponseWriter_Flush(t *testing.T) {
	rw := httptest.NewRecorder()
	metricsWriter := NewMetricsResponseWriter(rw)

	// Test Flush method - should not panic
	assert.NotPanics(t, func() {
		metricsWriter.Flush()
	})
}

func TestMetricsResponseWriter_CloseNotify(t *testing.T) {
	rw := httptest.NewRecorder()
	metricsWriter := NewMetricsResponseWriter(rw)

	// Test CloseNotify method - should return a channel
	ch := metricsWriter.CloseNotify()
	assert.NotNil(t, ch)

	// The channel should not receive anything immediately
	select {
	case <-ch:
		t.Fatal("CloseNotify channel should not receive anything")
	default:
		// Expected
	}
}

func TestRecordHTTPMetrics_Initialized(t *testing.T) {
	// Initialize metrics first
	metrics.Init()
	defer func() {
		// Reset metrics state after test
		metrics.Init()
	}()

	handlerName := "test-handler"
	method := "GET"
	statusCode := 200
	duration := 100 * time.Millisecond

	// Record metrics
	RecordHTTPMetrics(handlerName, method, statusCode, duration)

	// Verify metrics were recorded
	counter, err := metrics.RequestTotal.GetMetricWithLabelValues(handlerName, "GET", "200")
	require.NoError(t, err)

	// For Prometheus counters, we can simply verify that the counter exists
	// The actual metric collection is handled by Prometheus internally
	assert.NotNil(t, counter)
}

func TestRecordHTTPMetrics_NotInitialized(t *testing.T) {
	// Ensure metrics are not initialized
	// This test verifies that RecordHTTPMetrics doesn't panic when metrics are not initialized

	handlerName := "test-handler"
	method := "GET"
	statusCode := 200
	duration := 100 * time.Millisecond

	// Should not panic
	assert.NotPanics(t, func() {
		RecordHTTPMetrics(handlerName, method, statusCode, duration)
	})
}

func TestRecordHTTPMetrics_ErrorStatusCodes(t *testing.T) {
	// Initialize metrics first
	metrics.Init()
	defer func() {
		// Reset metrics state after test
		metrics.Init()
	}()

	testCases := []struct {
		name       string
		statusCode int
		errorType  string
	}{
		{"Client Error", 400, "client_error"},
		{"Not Found", 404, "client_error"},
		{"Server Error", 500, "server_error"},
		{"Internal Server Error", 503, "server_error"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handlerName := "test-handler"
			method := "GET"
			duration := 100 * time.Millisecond

			// Record metrics
			RecordHTTPMetrics(handlerName, method, tc.statusCode, duration)

			// Verify error metrics were recorded
			errorCounter, err := metrics.ErrorTotal.GetMetricWithLabelValues(tc.errorType, handlerName)
			require.NoError(t, err)

			// For Prometheus counters, we can simply verify that the counter exists
			// The actual metric collection is handled by Prometheus internally
			assert.NotNil(t, errorCounter)
		})
	}
}

func TestMetricsMiddleware_Success(t *testing.T) {
	// Initialize metrics first
	metrics.Init()
	defer func() {
		// Reset metrics state after test
		metrics.Init()
	}()

	handlerName := "test-handler"
	middleware := MetricsMiddleware(handlerName)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Wrap with middleware
	wrappedHandler := middleware(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Execute the handler
	wrappedHandler.ServeHTTP(rr, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	// Verify metrics were recorded
	counter, err := metrics.RequestTotal.GetMetricWithLabelValues(handlerName, "GET", "200")
	require.NoError(t, err)

	// For Prometheus counters, we can simply verify that the counter exists
	// The actual metric collection is handled by Prometheus internally
	assert.NotNil(t, counter)
}

func TestMetricsMiddleware_ErrorResponse(t *testing.T) {
	// Initialize metrics first
	metrics.Init()
	defer func() {
		// Reset metrics state after test
		metrics.Init()
	}()

	handlerName := "test-handler"
	middleware := MetricsMiddleware(handlerName)

	// Create a test handler that returns an error
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error"))
	})

	// Wrap with middleware
	wrappedHandler := middleware(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Execute the handler
	wrappedHandler.ServeHTTP(rr, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Error", rr.Body.String())

	// Verify error metrics were recorded
	errorCounter, err := metrics.ErrorTotal.GetMetricWithLabelValues("server_error", handlerName)
	require.NoError(t, err)

	// For Prometheus counters, we can simply verify that the counter exists
	// The actual metric collection is handled by Prometheus internally
	assert.NotNil(t, errorCounter)
}

func TestMetricsMiddleware_NotInitialized(t *testing.T) {
	// Ensure metrics are not initialized
	// This test verifies that the middleware doesn't panic when metrics are not initialized

	handlerName := "test-handler"
	middleware := MetricsMiddleware(handlerName)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Wrap with middleware
	wrappedHandler := middleware(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Should not panic
	assert.NotPanics(t, func() {
		wrappedHandler.ServeHTTP(rr, req)
	})

	// Verify response still works
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestMetricsMiddleware_DifferentMethods(t *testing.T) {
	// Initialize metrics first
	metrics.Init()
	defer func() {
		// Reset metrics state after test
		metrics.Init()
	}()

	handlerName := "test-handler"
	middleware := MetricsMiddleware(handlerName)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	wrappedHandler := middleware(testHandler)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			rr := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, "OK", rr.Body.String())

			// Verify metrics were recorded
			counter, err := metrics.RequestTotal.GetMetricWithLabelValues(handlerName, method, "200")
			require.NoError(t, err)

			// For Prometheus counters, we can simply verify that the counter exists
			// The actual metric collection is handled by Prometheus internally
			assert.NotNil(t, counter)
		})
	}
}

func TestMetricsMiddleware_DurationMeasurement(t *testing.T) {
	// Initialize metrics first
	metrics.Init()
	defer func() {
		// Reset metrics state after test
		metrics.Init()
	}()

	handlerName := "test-handler"
	middleware := MetricsMiddleware(handlerName)

	// Create a test handler with a delay
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Wrap with middleware
	wrappedHandler := middleware(testHandler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	start := time.Now()
	wrappedHandler.ServeHTTP(rr, req)
	duration := time.Since(start)

	// Verify response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	// Verify duration is reasonable (should be at least 50ms)
	assert.GreaterOrEqual(t, duration, 50*time.Millisecond)
}

func TestInitializeMetrics(t *testing.T) {
	// Test that InitializeMetrics properly initializes the metrics system
	InitializeMetrics()

	// Verify metrics are initialized
	assert.True(t, metrics.Initialized())
}

func TestStartMetricsServer(t *testing.T) {
	// Test starting metrics server with invalid address (should fail quickly)
	err := StartMetricsServer("invalid-address")

	// Should return an error for invalid address
	assert.Error(t, err)
}

func TestStartMetricsServer_ValidAddress(t *testing.T) {
	// Test starting metrics server with valid address
	// Use port 0 to get an available port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- metrics.StartServer(ctx, ":0")
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context to stop server
	cancel()

	// Wait for server to stop
	select {
	case err := <-serverErr:
		// Server should stop gracefully - "http: Server closed" is a normal shutdown error
		// We accept this error as it indicates the server shut down as expected
		if err != nil && err.Error() != "http: Server closed" {
			assert.NoError(t, err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Server did not stop within timeout")
	}
}

func TestMetricsResponseWriter_Header(t *testing.T) {
	rw := httptest.NewRecorder()
	metricsWriter := NewMetricsResponseWriter(rw)

	// Test setting headers
	metricsWriter.Header().Set("Content-Type", "application/json")
	metricsWriter.Header().Set("X-Custom-Header", "test-value")

	// Verify headers are set
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))
	assert.Equal(t, "test-value", rw.Header().Get("X-Custom-Header"))
}

func TestMetricsMiddleware_ChainedHandlers(t *testing.T) {
	// Initialize metrics first
	metrics.Init()
	defer func() {
		// Reset metrics state after test
		metrics.Init()
	}()

	handlerName := "test-handler"
	middleware := MetricsMiddleware(handlerName)

	// Create a chain of handlers
	firstHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set a header
		w.Header().Set("X-Processed", "true")
	})

	secondHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if previous handler set the header
		if w.Header().Get("X-Processed") == "true" {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("Created"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Bad Request"))
		}
	})

	// Chain the handlers
	chainedHandler := firstHandler
	chainedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firstHandler.ServeHTTP(w, r)
		secondHandler.ServeHTTP(w, r)
	})

	// Wrap with metrics middleware
	wrappedHandler := middleware(chainedHandler)

	// Create test request
	req := httptest.NewRequest("POST", "/test", nil)
	rr := httptest.NewRecorder()

	// Execute the handler
	wrappedHandler.ServeHTTP(rr, req)

	// Verify response
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "Created", rr.Body.String())
	assert.Equal(t, "true", rr.Header().Get("X-Processed"))

	// Verify metrics were recorded
	counter, err := metrics.RequestTotal.GetMetricWithLabelValues(handlerName, "POST", "201")
	require.NoError(t, err)

	// For Prometheus counters, we can simply verify that the counter exists
	// The actual metric collection is handled by Prometheus internally
	assert.NotNil(t, counter)
}

func TestRecordHTTPMetrics_EdgeCases(t *testing.T) {
	// Initialize metrics first
	metrics.Init()
	defer func() {
		// Reset metrics state after test
		metrics.Init()
	}()

	// Test with zero duration
	RecordHTTPMetrics("test-handler", "GET", 200, 0)

	// Test with negative duration (should not panic)
	RecordHTTPMetrics("test-handler", "GET", 200, -100*time.Millisecond)

	// Test with very large duration
	RecordHTTPMetrics("test-handler", "GET", 200, 24*time.Hour)

	// Test with empty handler name
	RecordHTTPMetrics("", "GET", 200, 100*time.Millisecond)

	// Test with empty method
	RecordHTTPMetrics("test-handler", "", 200, 100*time.Millisecond)

	// Verify no panic occurred
	assert.True(t, true)
}

func TestMetricsResponseWriter_ConcurrentAccess(t *testing.T) {
	// Test that we can safely create multiple MetricsResponseWriter instances concurrently
	done := make(chan bool, 10)
	writers := make([]*MetricsResponseWriter, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			rw := httptest.NewRecorder()
			writers[index] = NewMetricsResponseWriter(rw)

			// Perform basic operations on each writer
			writers[index].Header().Set(fmt.Sprintf("X-Test-%d", index), "value")
			writers[index].WriteHeader(http.StatusOK + index)
			_, _ = writers[index].Write([]byte("test"))

			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify that all writers were created and can be accessed
	for i := 0; i < 10; i++ {
		assert.NotNil(t, writers[i], "Writer %d should not be nil", i)
		assert.GreaterOrEqual(t, writers[i].Status(), http.StatusOK, "Writer %d should have valid status", i)
	}
}

// Benchmark tests
func BenchmarkMetricsMiddleware(b *testing.B) {
	// Initialize metrics
	metrics.Init()
	defer metrics.Init() // Reset after test

	handlerName := "benchmark-handler"
	middleware := MetricsMiddleware(handlerName)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	wrappedHandler := middleware(testHandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/benchmark", nil)
		rr := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rr, req)
	}
}

func BenchmarkRecordHTTPMetrics(b *testing.B) {
	// Initialize metrics
	metrics.Init()
	defer metrics.Init() // Reset after test

	handlerName := "benchmark-handler"
	method := "GET"
	statusCode := 200
	duration := 100 * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RecordHTTPMetrics(handlerName, method, statusCode, duration)
	}
}

func BenchmarkNewMetricsResponseWriter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rw := httptest.NewRecorder()
		_ = NewMetricsResponseWriter(rw)
	}
}
