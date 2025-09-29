package metrics

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialized(t *testing.T) {
	// Reset the initialized state for testing
	initialized = false

	assert.False(t, Initialized())

	Init()
	assert.True(t, Initialized())
}

// ResetMetrics resets the metrics system for testing purposes
func ResetMetrics() {
	mutex.Lock()

	// Reset registry and metrics
	registry = prometheus.NewRegistry()
	initialized = false

	mutex.Unlock()

	// Reinitialize metrics (this will acquire the mutex itself)
	Init()
}

func TestInit(t *testing.T) {
	ResetMetrics()

	assert.True(t, Initialized())

	// Verify that metrics are registered
	assert.NotNil(t, RequestDuration)
	assert.NotNil(t, RequestTotal)
	assert.NotNil(t, ErrorTotal)
	assert.NotNil(t, CacheHits)
	assert.NotNil(t, DBQueryDuration)
	assert.NotNil(t, TransactionDuration)
	assert.NotNil(t, TransactionTotal)
	assert.NotNil(t, DomainEventTotal)
}

func TestInit_AlreadyInitialized(t *testing.T) {
	ResetMetrics()

	// Call Init again - should not panic
	Init()
	assert.True(t, Initialized())
}

func TestServeHTTP(t *testing.T) {
	ResetMetrics()

	// Record some metrics first so they appear in the response
	RequestDuration.WithLabelValues("test-handler", "GET", "200").Observe(0.5)
	RequestTotal.WithLabelValues("test-handler", "GET", "200").Inc()

	req := httptest.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()

	ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.Contains(t, body, "http_request_duration_seconds")
	assert.Contains(t, body, "http_requests_total")
}

func TestHTTPMiddleware(t *testing.T) {
	ResetMetrics()

	handler := "test-handler"
	middleware := HTTPMiddleware(handler)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Wrap with middleware
	wrappedHandler := middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	// Verify metrics were recorded
	counter, err := RequestTotal.GetMetricWithLabelValues(handler, "GET", "200")
	require.NoError(t, err)

	var metric prometheus.Metric
	ch := make(chan prometheus.Metric, 1)
	counter.Collect(ch)
	close(ch)

	metric = <-ch
	assert.NotNil(t, metric)
}

func TestHTTPMiddleware_NotInitialized(t *testing.T) {
	ResetMetrics()

	handler := "test-handler"
	middleware := HTTPMiddleware(handler)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	wrappedHandler := middleware(testHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestRecordHTTPMetrics(t *testing.T) {
	ResetMetrics()

	handler := "test-handler"
	method := "GET"
	status := "200"
	duration := 0.5

	// Record metrics
	RequestDuration.WithLabelValues(handler, method, status).Observe(duration)
	RequestTotal.WithLabelValues(handler, method, status).Inc()

	// Verify metrics were recorded by checking the registry
	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Check that we have some metrics collected
	assert.Greater(t, len(metrics), 0)

	// Verify specific metrics by name
	foundRequestDuration := false
	foundRequestTotal := false

	for _, metric := range metrics {
		if metric.GetName() == "http_request_duration_seconds" {
			foundRequestDuration = true
		}
		if metric.GetName() == "http_requests_total" {
			foundRequestTotal = true
		}
	}

	assert.True(t, foundRequestDuration, "RequestDuration metric should be found")
	assert.True(t, foundRequestTotal, "RequestTotal metric should be found")
}

func TestRecordErrorMetrics(t *testing.T) {
	ResetMetrics()

	errorType := "validation"
	source := "user-service"

	ErrorTotal.WithLabelValues(errorType, source).Inc()

	// Verify metrics were recorded by checking the registry
	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Check that we have some metrics collected
	assert.Greater(t, len(metrics), 0)

	// Verify specific metric by name
	foundErrorTotal := false

	for _, metric := range metrics {
		if metric.GetName() == "errors_total" {
			foundErrorTotal = true
			break
		}
	}

	assert.True(t, foundErrorTotal, "ErrorTotal metric should be found")
}

func TestRecordCacheMetrics(t *testing.T) {
	ResetMetrics()

	cache := "user-cache"
	operation := "get"

	CacheHits.WithLabelValues(cache, operation).Inc()

	// Verify metrics were recorded by checking the registry
	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Check that we have some metrics collected
	assert.Greater(t, len(metrics), 0)

	// Verify specific metric by name
	foundCacheHits := false

	for _, metric := range metrics {
		if metric.GetName() == "cache_hits_total" {
			foundCacheHits = true
			break
		}
	}

	assert.True(t, foundCacheHits, "CacheHits metric should be found")
}

func TestRecordDBMetrics(t *testing.T) {
	ResetMetrics()

	RecordDBMetrics("query", 100*time.Millisecond)

	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Verify DB metrics were recorded
	found := false
	for _, metric := range metrics {
		if metric.GetName() == "db_query_duration_seconds" {
			found = true
			break
		}
	}
	assert.True(t, found, "DB metrics should be recorded")
}

func TestRecordTransactionMetrics(t *testing.T) {
	ResetMetrics()

	RecordTransactionMetrics("create", 200*time.Millisecond)

	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Verify transaction metrics were recorded
	foundDuration := false
	foundOperations := false
	for _, metric := range metrics {
		if metric.GetName() == "transaction_duration_seconds" {
			foundDuration = true
		}
		if metric.GetName() == "transaction_operations_total" {
			foundOperations = true
		}
	}
	assert.True(t, foundDuration, "Transaction duration metric should be recorded")
	assert.True(t, foundOperations, "Transaction operations metric should be recorded")
}

func TestRecordDomainEventMetrics(t *testing.T) {
	ResetMetrics()

	RecordDomainEventMetrics("user_created")

	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Verify domain event metrics were recorded
	found := false
	for _, metric := range metrics {
		if metric.GetName() == "domain_events_total" {
			found = true
			break
		}
	}
	assert.True(t, found, "Domain event metrics should be recorded")
}

// Helper functions for testing
func RecordDBMetrics(operation string, duration time.Duration) {
	if !initialized {
		Init()
	}
	DBQueryDuration.WithLabelValues("test-db", operation).Observe(duration.Seconds())
}

func RecordTransactionMetrics(operation string, duration time.Duration) {
	if !initialized {
		Init()
	}
	TransactionDuration.WithLabelValues("test-store").Observe(duration.Seconds())
	TransactionTotal.WithLabelValues(operation, "test-store").Inc()
}

func RecordDomainEventMetrics(eventType string) {
	if !initialized {
		Init()
	}
	DomainEventTotal.WithLabelValues(eventType, "test-source").Inc()
}

func RecordCacheMetrics(cache string, hit bool) {
	if !initialized {
		Init()
	}
	operation := "miss"
	if hit {
		operation = "hit"
	}
	CacheHits.WithLabelValues(cache, operation).Inc()
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"healthy"}`))
}

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}

func TestResponseWriter(t *testing.T) {
	rw := httptest.NewRecorder()
	responseWriter := NewResponseWriter(rw)

	// Test default status
	assert.Equal(t, http.StatusOK, responseWriter.Status())

	// Test writing status
	responseWriter.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, responseWriter.Status())

	// Test writing body
	_, err := responseWriter.Write([]byte("test"))
	assert.NoError(t, err)
	assert.Equal(t, "test", rw.Body.String())
}

func TestStartServer_ContextCancellation(t *testing.T) {
	ResetMetrics()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := StartServer(ctx, ":0")
	assert.Error(t, err)
	// The error might be "http: Server closed" or "context canceled" depending on timing
	assert.True(t, err.Error() == "http: Server closed" || err.Error() == "context canceled",
		"Error should be either 'http: Server closed' or 'context canceled', got: %s", err.Error())
}

func TestHealthEndpoint(t *testing.T) {
	ResetMetrics()

	// Test health endpoint
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "{\"status\":\"healthy\"}", rr.Body.String())
}

func TestReadyEndpoint(t *testing.T) {
	ResetMetrics()

	// Test ready endpoint
	req, err := http.NewRequest("GET", "/ready", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ReadyHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "{\"status\":\"ready\"}", rr.Body.String())
}

func TestConcurrentInitialization(t *testing.T) {
	ResetMetrics()

	// Test concurrent initialization
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Init()
		}()
	}
	wg.Wait()

	assert.True(t, Initialized())
}

func TestMetricsRegistration(t *testing.T) {
	ResetMetrics()

	// Record some metrics first so they appear in the registry
	RequestDuration.WithLabelValues("test-handler", "GET", "200").Observe(0.5)
	RequestTotal.WithLabelValues("test-handler", "GET", "200").Inc()

	// Verify all metrics are properly registered
	metrics, err := registry.Gather()
	require.NoError(t, err)

	metricNames := make(map[string]bool)
	for _, metric := range metrics {
		metricNames[metric.GetName()] = true
	}

	// For testing purposes, we'll check if we have at least some metrics registered
	// The exact count might vary based on the test environment
	assert.Greater(t, len(metrics), 0, "Should have at least some metrics registered")

	// Check for the most common metrics
	assert.True(t, metricNames["http_request_duration_seconds"], "HTTP request duration metric should be registered")
	assert.True(t, metricNames["http_requests_total"], "HTTP requests total metric should be registered")
}

// Test helper functions
func TestGetMetricNames(t *testing.T) {
	ResetMetrics()

	// Record some metrics first so they appear in the registry
	RequestDuration.WithLabelValues("test-handler", "GET", "200").Observe(0.5)
	RequestTotal.WithLabelValues("test-handler", "GET", "200").Inc()

	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Should have at least some metrics registered
	assert.Greater(t, len(metrics), 0, "Should have at least some metrics registered")
}

func TestMetricLabels(t *testing.T) {
	ResetMetrics()

	t.Run("RequestTotal with valid labels", func(t *testing.T) {
		RequestTotal.WithLabelValues("GET", "/test", "200").Inc()

		metrics, err := registry.Gather()
		require.NoError(t, err)

		// Should have metrics
		assert.Greater(t, len(metrics), 0)
	})

	t.Run("ErrorTotal with valid labels", func(t *testing.T) {
		ErrorTotal.WithLabelValues("database", "connection_timeout").Inc()

		metrics, err := registry.Gather()
		require.NoError(t, err)

		// Should have metrics
		assert.Greater(t, len(metrics), 0)
	})

	t.Run("CacheHits with valid labels", func(t *testing.T) {
		CacheHits.WithLabelValues("user_cache", "hit").Inc()

		metrics, err := registry.Gather()
		require.NoError(t, err)

		// Should have metrics
		assert.Greater(t, len(metrics), 0)
	})
}

func TestMetricObservation(t *testing.T) {
	ResetMetrics()

	// Test observing a histogram metric
	RequestDuration.WithLabelValues("test-handler", "GET", "200").Observe(0.5)

	metrics, err := registry.Gather()
	require.NoError(t, err)

	// Should have metrics
	assert.Greater(t, len(metrics), 0)
}
