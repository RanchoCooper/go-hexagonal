// Package error_handler provides unified error handling utilities for the application
package error_handler

import (
	"context"
	"fmt"

	"go-hexagonal/api/error_code"
	"go-hexagonal/util/errors"
	"go-hexagonal/util/log"

	"go.uber.org/zap"
)

// ErrorHandler provides unified error handling capabilities
type ErrorHandler struct {
	// Additional configuration can be added here
}

// New creates a new ErrorHandler instance
func New() *ErrorHandler {
	// Ensure logger is initialized
	if log.SugaredLogger == nil {
		// Create a safe fallback logger for testing environments
		logger, _ := zap.NewDevelopment()
		log.SugaredLogger = logger.Sugar()
	}
	return &ErrorHandler{}
}

// HandleError handles errors with context and logging
func (h *ErrorHandler) HandleError(ctx context.Context, err error, operation string) error {
	if err == nil {
		return nil
	}

	// Log the error with context
	h.logErrorWithContext(ctx, err, operation)

	// Return the original error for further handling
	return err
}

// HandleAndWrapError handles errors and wraps them with additional context
func (h *ErrorHandler) HandleAndWrapError(ctx context.Context, err error, operation string, message string) error {
	if err == nil {
		return nil
	}

	// Log the error with context
	h.logErrorWithContext(ctx, err, operation)

	// Wrap the error with additional context
	return fmt.Errorf("%s: %w", message, err)
}

// HandleAndConvertError handles errors and converts them to appropriate error types
func (h *ErrorHandler) HandleAndConvertError(ctx context.Context, err error, operation string, errorType errors.ErrorType) error {
	if err == nil {
		return nil
	}

	// Log the error with context
	h.logErrorWithContext(ctx, err, operation)

	// Convert to appropriate error type
	switch errorType {
	case errors.ErrorTypeValidation:
		return errors.NewValidationError(operation+" failed", err)
	case errors.ErrorTypeNotFound:
		return errors.NewNotFoundError(operation+" failed", err)
	case errors.ErrorTypePersistence:
		return errors.NewPersistenceError(operation+" failed", err)
	case errors.ErrorTypeSystem:
		return errors.NewSystemError(operation+" failed", err)
	case errors.ErrorTypeBusiness:
		return errors.NewBusinessError(operation+" failed", err)
	default:
		return errors.Wrap(err, errors.ErrorTypeSystem, operation+" failed")
	}
}

// HandleDomainError handles domain-specific errors
func (h *ErrorHandler) HandleDomainError(ctx context.Context, err error, operation string, domain string) error {
	if err == nil {
		return nil
	}

	// Log the error with context
	h.logErrorWithContext(ctx, err, operation)

	// Create a domain-specific error
	return errors.Newf(errors.ErrorTypeBusiness, "%s error in %s domain: %v", operation, domain, err)
}

// HandleAPIError handles API-specific errors and converts them to error_code.Error
func (h *ErrorHandler) HandleAPIError(ctx context.Context, err error, operation string) *error_code.Error {
	if err == nil {
		return nil
	}

	// Log the error with context
	h.logErrorWithContext(ctx, err, operation)

	// Convert to appropriate API error
	if appErr, ok := err.(*errors.AppError); ok {
		return h.convertAppErrorToAPIError(appErr)
	}

	// Handle standard error types
	switch {
	case errors.IsValidationError(err):
		return error_code.InvalidParams.WithMessage("Validation error: %s", err.Error())
	case errors.IsNotFoundError(err):
		return error_code.NotFound.WithMessage("Resource not found: %s", err.Error())
	case errors.IsPersistenceError(err):
		return error_code.ServerError.WithMessage("Database operation failed")
	case errors.IsSystemError(err):
		return error_code.ServerError.WithMessage("Internal system error")
	default:
		return error_code.ServerError.WithMessage("Internal server error")
	}
}

// convertAppErrorToAPIError converts util/errors.AppError to api/error_code.Error
func (h *ErrorHandler) convertAppErrorToAPIError(appErr *errors.AppError) *error_code.Error {
	// Map AppError type to appropriate API error
	var apiErr *error_code.Error

	switch appErr.Type {
	case errors.ErrorTypeValidation:
		apiErr = error_code.InvalidParams.WithMessage("Validation error: %s", appErr.Message)
	case errors.ErrorTypeNotFound:
		apiErr = error_code.NotFound.WithMessage("Resource not found: %s", appErr.Message)
	case errors.ErrorTypeUnauthorized:
		apiErr = error_code.UnauthorizedTokenError.WithMessage("Unauthorized: %s", appErr.Message)
	case errors.ErrorTypeForbidden:
		apiErr = error_code.UnauthorizedTokenError.WithMessage("Access forbidden")
	case errors.ErrorTypeConflict:
		apiErr = error_code.AccountExist.WithMessage("Conflict: %s", appErr.Message)
	case errors.ErrorTypePersistence, errors.ErrorTypeSystem:
		apiErr = error_code.ServerError.WithMessage("Server error: %s", appErr.Message)
	case errors.ErrorTypeBusiness:
		apiErr = error_code.ServerError.WithMessage("Business error: %s", appErr.Message)
	default:
		apiErr = error_code.ServerError.WithMessage("Server error: %s", appErr.Message)
	}

	// Add details if available
	if len(appErr.Details) > 0 {
		if details, ok := appErr.Details["details"].([]string); ok {
			apiErr = apiErr.WithDetails(details...)
		}
	}

	return apiErr
}

// logErrorWithContext logs errors with request context information
func (h *ErrorHandler) logErrorWithContext(ctx context.Context, err error, operation string) {
	// Extract context information
	requestID := ""
	if reqID, ok := ctx.Value("X-Request-ID").(string); ok {
		requestID = reqID
	}

	// Log based on error type
	switch {
	case errors.IsValidationError(err):
		log.SugaredLogger.Warnf("Validation error [%s] %s - Error: %v", requestID, operation, err)
	case errors.IsNotFoundError(err):
		log.SugaredLogger.Infof("Resource not found [%s] %s - Error: %v", requestID, operation, err)
	case errors.IsPersistenceError(err):
		log.SugaredLogger.Errorf("Persistence error [%s] %s - Error: %v", requestID, operation, err)
	case errors.IsSystemError(err):
		log.SugaredLogger.Errorf("System error [%s] %s - Error: %v", requestID, operation, err)
	default:
		log.SugaredLogger.Errorf("Unexpected error [%s] %s - Error: %v", requestID, operation, err)
	}
}

// IsRetryableError checks if an error is retryable
func (h *ErrorHandler) IsRetryableError(err error) bool {
	// Network errors, temporary failures, etc. are retryable
	// Business logic errors, validation errors are not retryable
	switch {
	case errors.IsValidationError(err):
		return false
	case errors.IsNotFoundError(err):
		return false
	case errors.IsBusinessError(err):
		return false
	case errors.IsPersistenceError(err):
		// Some persistence errors might be retryable (e.g., deadlock)
		return true
	case errors.IsSystemError(err):
		// System errors might be retryable if temporary
		return true
	default:
		return false
	}
}

// ShouldLogError determines if an error should be logged
func (h *ErrorHandler) ShouldLogError(err error) bool {
	// Don't log validation errors at error level
	if errors.IsValidationError(err) {
		return false
	}

	// Don't log not found errors at error level
	if errors.IsNotFoundError(err) {
		return false
	}

	// Log all other errors
	return true
}

// Global error handler instance
var globalErrorHandler = New()

// HandleError is a convenience function that uses the global error handler
func HandleError(ctx context.Context, err error, operation string) error {
	return globalErrorHandler.HandleError(ctx, err, operation)
}

// HandleAndWrapError is a convenience function that uses the global error handler
func HandleAndWrapError(ctx context.Context, err error, operation string, message string) error {
	return globalErrorHandler.HandleAndWrapError(ctx, err, operation, message)
}

// HandleAndConvertError is a convenience function that uses the global error handler
func HandleAndConvertError(ctx context.Context, err error, operation string, errorType errors.ErrorType) error {
	return globalErrorHandler.HandleAndConvertError(ctx, err, operation, errorType)
}

// HandleDomainError is a convenience function that uses the global error handler
func HandleDomainError(ctx context.Context, err error, operation string, domain string) error {
	return globalErrorHandler.HandleDomainError(ctx, err, operation, domain)
}

// HandleAPIError is a convenience function that uses the global error handler
func HandleAPIError(ctx context.Context, err error, operation string) *error_code.Error {
	return globalErrorHandler.HandleAPIError(ctx, err, operation)
}

// IsRetryableError is a convenience function that uses the global error handler
func IsRetryableError(err error) bool {
	return globalErrorHandler.IsRetryableError(err)
}

// ShouldLogError is a convenience function that uses the global error handler
func ShouldLogError(err error) bool {
	return globalErrorHandler.ShouldLogError(err)
}
