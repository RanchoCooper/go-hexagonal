// Package errors provides unified error creation and handling methods
package errors

import (
	stderrors "errors"
	"fmt"
)

// ErrorType defines the error type
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "VALIDATION"
	// ErrorTypeNotFound represents resource not found errors
	ErrorTypeNotFound ErrorType = "NOT_FOUND"
	// ErrorTypePersistence represents persistence layer errors
	ErrorTypePersistence ErrorType = "PERSISTENCE"
	// ErrorTypeSystem represents internal system errors
	ErrorTypeSystem ErrorType = "SYSTEM"
	// ErrorTypeBusiness represents business logic errors
	ErrorTypeBusiness ErrorType = "BUSINESS"
	// ErrorTypeUnauthorized represents authentication/authorization errors
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	// ErrorTypeForbidden represents permission errors
	ErrorTypeForbidden ErrorType = "FORBIDDEN"
	// ErrorTypeConflict represents resource conflict errors
	ErrorTypeConflict ErrorType = "CONFLICT"
)

// AppError defines the application error structure
type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
	Details map[string]any
	Code    int
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap supports Go 1.13+ error unwrapping
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithDetails adds error details
func (e *AppError) WithDetails(details map[string]any) *AppError {
	e.Details = details
	return e
}

// WithCode adds an error code
func (e *AppError) WithCode(code int) *AppError {
	e.Code = code
	return e
}

// NewValidationError creates a validation error
func NewValidationError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Cause:   cause,
	}
}

// NewNotFoundError creates a resource not found error
func NewNotFoundError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Cause:   cause,
	}
}

// NewPersistenceError creates a persistence error
func NewPersistenceError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypePersistence,
		Message: message,
		Cause:   cause,
	}
}

// NewSystemError creates a system error
func NewSystemError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeSystem,
		Message: message,
		Cause:   cause,
	}
}

// NewBusinessError creates a business logic error
func NewBusinessError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeBusiness,
		Message: message,
		Cause:   cause,
	}
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr.Type == ErrorTypeValidation
	}
	return false
}

// IsNotFoundError checks if the error is a resource not found error
func IsNotFoundError(err error) bool {
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsPersistenceError checks if the error is a persistence error
func IsPersistenceError(err error) bool {
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr.Type == ErrorTypePersistence
	}
	return false
}

// IsSystemError checks if the error is a system error
func IsSystemError(err error) bool {
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr.Type == ErrorTypeSystem
	}
	return false
}

// IsBusinessError checks if the error is a business logic error
func IsBusinessError(err error) bool {
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr.Type == ErrorTypeBusiness
	}
	return false
}

// Wrap wraps a standard error as an application error
func Wrap(err error, errType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Cause:   err,
	}
}

// Wrapf wraps an error with a formatted message
func Wrapf(err error, errType ErrorType, format string, args ...any) *AppError {
	return &AppError{
		Type:    errType,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

// New creates a new application error
func New(errType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
	}
}

// Newf creates a new application error with a formatted message
func Newf(errType ErrorType, format string, args ...any) *AppError {
	return &AppError{
		Type:    errType,
		Message: fmt.Sprintf(format, args...),
	}
}

// StatusCode returns the appropriate HTTP status code for the error type
func (e *AppError) StatusCode() int {
	switch e.Type {
	case ErrorTypeValidation:
		return 400 // Bad Request
	case ErrorTypeNotFound:
		return 404 // Not Found
	case ErrorTypeUnauthorized:
		return 401 // Unauthorized
	case ErrorTypeForbidden:
		return 403 // Forbidden
	case ErrorTypeConflict:
		return 409 // Conflict
	case ErrorTypePersistence, ErrorTypeSystem:
		return 500 // Internal Server Error
	default:
		return 500 // Internal Server Error
	}
}

// Is implements error comparison for errors.Is
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Type == t.Type
}
