// Package errors provides unified error creation and handling methods
package errors

import (
	"fmt"

	stderrors "errors"
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
)

// AppError defines the application error structure
type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
	Details map[string]interface{}
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
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
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

// Wrap wraps a standard error as an application error
func Wrap(err error, errType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
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
