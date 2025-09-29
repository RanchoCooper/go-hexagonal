package log

import (
	"context"

	"go.uber.org/zap"
)

// LogHelper provides helper methods for common logging patterns
type LogHelper struct {
	component string
}

// NewLogHelper creates a new log helper for a specific component
func NewLogHelper(component string) *LogHelper {
	return &LogHelper{
		component: component,
	}
}

// WithContext creates a log context with the given operation
func (h *LogHelper) WithContext(ctx context.Context, operation string) *LogContext {
	return NewLogContext().WithComponent(h.component).WithOperation(operation)
}

// Errorf logs an error with structured context
func (h *LogHelper) Errorf(ctx context.Context, operation, format string, args ...any) {
	h.WithContext(ctx, operation)
	SugaredLogger.Errorf(format, args...)
}

// Error logs an error with structured context
func (h *LogHelper) Error(ctx context.Context, operation string, err error) {
	if err == nil {
		return
	}
	h.WithContext(ctx, operation)
	SugaredLogger.Errorw("Operation failed",
		"component", h.component,
		"operation", operation,
		"error", err.Error(),
	)
}

// ErrorWithDetails logs an error with additional details
func (h *LogHelper) ErrorWithDetails(ctx context.Context, operation string, err error, details map[string]any) {
	if err == nil {
		return
	}

	fields := []zap.Field{
		zap.String("component", h.component),
		zap.String("operation", operation),
		zap.String("error", err.Error()),
	}

	for key, value := range details {
		fields = append(fields, zap.Any(key, value))
	}

	Logger.Error("Operation failed", fields...)
}

// TransactionError logs a transaction-related error
func (h *LogHelper) TransactionError(ctx context.Context, operation string, err error) {
	h.ErrorWithDetails(ctx, operation, err, map[string]any{
		"transaction_operation": operation,
	})
}

// DatabaseError logs a database-related error
func (h *LogHelper) DatabaseError(ctx context.Context, operation string, err error) {
	h.ErrorWithDetails(ctx, operation, err, map[string]any{
		"database_operation": operation,
	})
}

// CacheError logs a cache-related error
func (h *LogHelper) CacheError(ctx context.Context, operation string, err error) {
	h.ErrorWithDetails(ctx, operation, err, map[string]any{
		"cache_operation": operation,
	})
}

// ValidationError logs a validation-related error
func (h *LogHelper) ValidationError(ctx context.Context, operation string, err error) {
	h.ErrorWithDetails(ctx, operation, err, map[string]any{
		"validation_operation": operation,
	})
}

// Info logs an informational message
func (h *LogHelper) Info(ctx context.Context, operation, message string) {
	h.WithContext(ctx, operation)
	SugaredLogger.Infow(message,
		"component", h.component,
		"operation", operation,
	)
}

// Debug logs a debug message
func (h *LogHelper) Debug(ctx context.Context, operation, message string) {
	h.WithContext(ctx, operation)
	SugaredLogger.Debugw(message,
		"component", h.component,
		"operation", operation,
	)
}

// Warn logs a warning message
func (h *LogHelper) Warn(ctx context.Context, operation, message string) {
	h.WithContext(ctx, operation)
	SugaredLogger.Warnw(message,
		"component", h.component,
		"operation", operation,
	)
}

// Global helper functions for common use cases

// ExampleServiceLogger returns a log helper for the example service
func ExampleServiceLogger() *LogHelper {
	return NewLogHelper("example_service")
}

// TransactionLogger returns a log helper for transaction operations
func TransactionLogger() *LogHelper {
	return NewLogHelper("transaction")
}

// DatabaseLogger returns a log helper for database operations
func DatabaseLogger() *LogHelper {
	return NewLogHelper("database")
}

// CacheLogger returns a log helper for cache operations
func CacheLogger() *LogHelper {
	return NewLogHelper("cache")
}

// APILogger returns a log helper for API operations
func APILogger() *LogHelper {
	return NewLogHelper("api")
}
