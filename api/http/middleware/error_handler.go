// Package middleware provides HTTP request processing middleware
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-hexagonal/api/error_code"
	"go-hexagonal/api/http/handle"
	"go-hexagonal/util/errors"
	"go-hexagonal/util/log"
)

// ErrorHandlerMiddleware handles API layer error responses uniformly
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handle.Error(c, err)
			return
		}

		// Handle errors that might be set in the context but not in c.Errors
		if c.Writer.Status() >= 400 {
			// If status code indicates error but no error was explicitly set
			// Create an appropriate error based on the status code
			switch c.Writer.Status() {
			case http.StatusBadRequest:
				handle.Error(c, error_code.InvalidParams)
			case http.StatusNotFound:
				handle.Error(c, error_code.NotFound)
			case http.StatusUnauthorized:
				handle.Error(c, error_code.UnauthorizedTokenError)
			case http.StatusForbidden:
				handle.Error(c, error_code.UnauthorizedTokenError.WithMessage("Access forbidden"))
			case http.StatusTooManyRequests:
				handle.Error(c, error_code.TooManyRequests)
			default:
				handle.Error(c, error_code.ServerError)
			}
		}
	}
}

// EnhancedErrorHandlerMiddleware provides comprehensive error handling with logging and metrics
func EnhancedErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store start time for performance metrics
		// startTime := time.Now()

		// Process the request
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Log the error with request context
			logErrorWithContext(c, err)

			// Handle the error
			handleEnhancedError(c, err)
			return
		}

		// Handle HTTP status code errors
		if c.Writer.Status() >= 400 {
			handleHTTPStatusError(c)
		}
	}
}

// logErrorWithContext logs errors with request context information
func logErrorWithContext(c *gin.Context, err error) {
	requestID := c.GetString("X-Request-ID")
	method := c.Request.Method
	path := c.Request.URL.Path
	statusCode := c.Writer.Status()

	// Log based on error type
	switch {
	case errors.IsValidationError(err):
		log.SugaredLogger.Warnf("Validation error [%s] %s %s - Status: %d - Error: %v",
			requestID, method, path, statusCode, err)
	case errors.IsNotFoundError(err):
		log.SugaredLogger.Infof("Resource not found [%s] %s %s - Status: %d - Error: %v",
			requestID, method, path, statusCode, err)
	case errors.IsPersistenceError(err):
		log.SugaredLogger.Errorf("Persistence error [%s] %s %s - Status: %d - Error: %v",
			requestID, method, path, statusCode, err)
	case errors.IsSystemError(err):
		log.SugaredLogger.Errorf("System error [%s] %s %s - Status: %d - Error: %v",
			requestID, method, path, statusCode, err)
	default:
		log.SugaredLogger.Errorf("Unexpected error [%s] %s %s - Status: %d - Error: %v",
			requestID, method, path, statusCode, err)
	}
}

// handleEnhancedError provides enhanced error handling with better error mapping
func handleEnhancedError(c *gin.Context, err error) {
	// Try to handle as api/error_code.Error first
	if apiErr, ok := err.(*error_code.Error); ok {
		handle.Error(c, apiErr)
		return
	}

	// Try to handle as util/errors.AppError
	if appErr, ok := err.(*errors.AppError); ok {
		handleAppError(c, appErr)
		return
	}

	// Handle standard error types
	switch {
	case errors.IsValidationError(err):
		handle.Error(c, error_code.InvalidParams.WithMessage("Validation error: %s", err.Error()))
	case errors.IsNotFoundError(err):
		handle.Error(c, error_code.NotFound.WithMessage("Resource not found: %s", err.Error()))
	case errors.IsPersistenceError(err):
		handle.Error(c, error_code.ServerError.WithMessage("Database operation failed"))
	case errors.IsSystemError(err):
		handle.Error(c, error_code.ServerError.WithMessage("Internal system error"))
	case errors.IsBusinessError(err):
		handle.Error(c, error_code.ServerError.WithMessage("Business error: %s", err.Error()))
	default:
		handle.Error(c, error_code.ServerError.WithMessage("Internal server error"))
	}
}

// handleAppError converts util/errors.AppError to API response
func handleAppError(c *gin.Context, appErr *errors.AppError) {
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

	handle.Error(c, apiErr)
}

// handleHTTPStatusError handles errors based on HTTP status code
func handleHTTPStatusError(c *gin.Context) {
	switch c.Writer.Status() {
	case http.StatusBadRequest:
		handle.Error(c, error_code.InvalidParams)
	case http.StatusUnauthorized:
		handle.Error(c, error_code.UnauthorizedTokenError)
	case http.StatusForbidden:
		handle.Error(c, error_code.UnauthorizedTokenError.WithMessage("Access forbidden"))
	case http.StatusNotFound:
		handle.Error(c, error_code.NotFound)
	case http.StatusTooManyRequests:
		handle.Error(c, error_code.TooManyRequests)
	case http.StatusConflict:
		handle.Error(c, error_code.AccountExist.WithMessage("Resource conflict"))
	default:
		handle.Error(c, error_code.ServerError)
	}
}
