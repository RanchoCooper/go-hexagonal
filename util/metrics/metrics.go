// Package metrics provides functionality for collecting and exposing application metrics
package metrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registry    = prometheus.NewRegistry()
	initialized = false
	mutex       sync.Mutex

	// RequestDuration measures the duration of HTTP requests
	RequestDuration *prometheus.HistogramVec

	// RequestTotal counts the total number of HTTP requests
	RequestTotal *prometheus.CounterVec

	// ErrorTotal counts the total number of errors
	ErrorTotal *prometheus.CounterVec

	// CacheHits tracks cache hits and misses
	CacheHits *prometheus.CounterVec

	// DBQueryDuration measures the duration of database queries
	DBQueryDuration *prometheus.HistogramVec

	// TransactionDuration measures the duration of database transactions
	TransactionDuration *prometheus.HistogramVec

	// TransactionTotal counts the total number of transaction operations
	TransactionTotal *prometheus.CounterVec

	// DomainEventTotal counts the total number of domain events
	DomainEventTotal *prometheus.CounterVec
)

// Initialized returns whether metrics has been initialized
func Initialized() bool {
	return initialized
}

// Init initializes the metrics collection system
func Init() {
	mutex.Lock()
	defer mutex.Unlock()

	if initialized {
		return
	}

	// HTTP metrics
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "method", "status"},
	)

	RequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"handler", "method", "status"},
	)

	// Error metrics
	ErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type", "source"},
	)

	// Cache metrics
	CacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache", "operation"},
	)

	// Database metrics
	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"db", "operation"},
	)

	TransactionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_duration_seconds",
			Help:    "Duration of transactions",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"store_type"},
	)

	TransactionTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_operations_total",
			Help: "Total number of transaction operations",
		},
		[]string{"operation", "store_type"},
	)

	// Domain event metrics
	DomainEventTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "domain_events_total",
			Help: "Total number of domain events",
		},
		[]string{"event_type", "source"},
	)

	// Register all metrics
	registry.MustRegister(
		RequestDuration,
		RequestTotal,
		ErrorTotal,
		CacheHits,
		DBQueryDuration,
		TransactionDuration,
		TransactionTotal,
		DomainEventTotal,
	)

	initialized = true
}

// ServeHTTP serves the metrics endpoint for Prometheus scraping
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

// StartServer starts a metrics server on the given address
func StartServer(ctx context.Context, addr string) error {
	if !initialized {
		Init()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", ServeHTTP)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	fmt.Printf("Metrics server started on %s\n", addr)
	return server.ListenAndServe()
}

// HTTPMiddleware creates a middleware for measuring HTTP request metrics
func HTTPMiddleware(handler string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !initialized {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			rw := NewResponseWriter(w)
			next.ServeHTTP(rw, r)

			duration := time.Since(start).Seconds()
			status := fmt.Sprintf("%d", rw.Status())
			RequestDuration.WithLabelValues(handler, r.Method, status).Observe(duration)
			RequestTotal.WithLabelValues(handler, r.Method, status).Inc()
		})
	}
}

// MeasureDBQuery measures the duration of a database query
func MeasureDBQuery(db, operation string, f func() error) error {
	if !initialized {
		return f()
	}

	start := time.Now()
	err := f()
	duration := time.Since(start).Seconds()

	DBQueryDuration.WithLabelValues(db, operation).Observe(duration)
	if err != nil {
		ErrorTotal.WithLabelValues("db", db).Inc()
	}

	return err
}

// MeasureTransaction measures the duration of a database transaction
func MeasureTransaction(storeType string, f func() error) error {
	if !initialized {
		return f()
	}

	start := time.Now()
	err := f()
	duration := time.Since(start).Seconds()

	TransactionDuration.WithLabelValues(storeType).Observe(duration)
	if err != nil {
		ErrorTotal.WithLabelValues("transaction", storeType).Inc()
	}

	return err
}

// RecordTransactionOperation records a transaction operation
func RecordTransactionOperation(operation, storeType string) {
	if !initialized {
		return
	}
	TransactionTotal.WithLabelValues(operation, storeType).Inc()
}

// RecordCacheHit records a cache hit or miss
func RecordCacheHit(cache, operation string) {
	if !initialized {
		return
	}
	CacheHits.WithLabelValues(cache, operation).Inc()
}

// RecordDomainEvent records a domain event
func RecordDomainEvent(eventType, source string) {
	if !initialized {
		return
	}
	DomainEventTotal.WithLabelValues(eventType, source).Inc()
}

// RecordError records an error
func RecordError(errorType, source string) {
	if !initialized {
		return
	}
	ErrorTotal.WithLabelValues(errorType, source).Inc()
}

// ResponseWriter is a wrapper around http.ResponseWriter that captures the status code
type ResponseWriter struct {
	http.ResponseWriter
	status int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Status returns the status code
func (rw *ResponseWriter) Status() int {
	return rw.status
}
