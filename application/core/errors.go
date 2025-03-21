package core

import (
	"errors"
	"fmt"

	apperrors "go-hexagonal/util/errors"
)

// HTTP status codes as registered with IANA
const (
	StatusBadRequest          = 400 // RFC 7231, 6.5.1
	StatusNotFound            = 404 // RFC 7231, 6.5.4
	StatusInternalServerError = 500 // RFC 7231, 6.6.1
)

// Error types
const (
	ErrorTypeValidation   = "VALIDATION_ERROR"
	ErrorTypeNotFound     = "NOT_FOUND"
	ErrorTypeUnauthorized = "UNAUTHORIZED"
	ErrorTypeForbidden    = "FORBIDDEN"
	ErrorTypeConflict     = "CONFLICT"
	ErrorTypeInternal     = "INTERNAL_ERROR"
)

// Compatibility errors for backward compatibility
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("resource conflict")
	ErrInternal     = errors.New("internal error")
)

// Error represents an application error (will be deprecated in favor of util/errors)
type Error struct {
	Type    string // Error type
	Code    int    // Error code
	Message string // Error message
	Err     error  // Original error
}

// Error returns the error message
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the original error
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new application error
// Deprecated: Use util/errors package instead
func NewError(errorType string, code int, message string, err error) *Error {
	return &Error{
		Type:    errorType,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewValidationError creates a validation error
// Deprecated: Use util/errors.NewValidationError instead
func NewValidationError(code int, message string, err error) *Error {
	return NewError("VALIDATION_ERROR", code, message, err)
}

// NewNotFoundError creates a not found error
// Deprecated: Use util/errors.NewNotFoundError instead
func NewNotFoundError(code int, message string, err error) *Error {
	return NewError("NOT_FOUND", code, message, err)
}

// NewUnauthorizedError creates an unauthorized error
// Deprecated: Use util/errors package instead
func NewUnauthorizedError(code int, message string, err error) *Error {
	return NewError("UNAUTHORIZED", code, message, err)
}

// NewForbiddenError creates a forbidden error
// Deprecated: Use util/errors package instead
func NewForbiddenError(code int, message string, err error) *Error {
	return NewError("FORBIDDEN", code, message, err)
}

// NewConflictError creates a conflict error
// Deprecated: Use util/errors package instead
func NewConflictError(code int, message string, err error) *Error {
	return NewError("CONFLICT", code, message, err)
}

// NewInternalError creates an internal error
// Deprecated: Use util/errors.NewSystemError instead
func NewInternalError(code int, message string, err error) *Error {
	return NewError("INTERNAL_ERROR", code, message, err)
}

// ToAppError converts a legacy Error to the new AppError type
func ToAppError(err *Error) *apperrors.AppError {
	switch err.Type {
	case "VALIDATION_ERROR":
		return apperrors.NewValidationError(err.Message, err.Err)
	case "NOT_FOUND":
		return apperrors.NewNotFoundError(err.Message, err.Err)
	case "INTERNAL_ERROR":
		return apperrors.NewSystemError(err.Message, err.Err)
	default:
		return apperrors.New(apperrors.ErrorTypeSystem, err.Message)
	}
}
