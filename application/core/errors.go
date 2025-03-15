package core

import (
	"errors"
	"fmt"
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

// Application layer error definitions
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("resource conflict")
	ErrInternal     = errors.New("internal error")
)

// Error represents an application error
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
func NewError(errorType string, code int, message string, err error) *Error {
	return &Error{
		Type:    errorType,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewValidationError creates a validation error
func NewValidationError(code int, message string, err error) *Error {
	return NewError(ErrorTypeValidation, code, message, err)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(code int, message string, err error) *Error {
	return NewError(ErrorTypeNotFound, code, message, err)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(code int, message string, err error) *Error {
	return NewError(ErrorTypeUnauthorized, code, message, err)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(code int, message string, err error) *Error {
	return NewError(ErrorTypeForbidden, code, message, err)
}

// NewConflictError creates a conflict error
func NewConflictError(code int, message string, err error) *Error {
	return NewError(ErrorTypeConflict, code, message, err)
}

// NewInternalError creates an internal error
func NewInternalError(code int, message string, err error) *Error {
	return NewError(ErrorTypeInternal, code, message, err)
}
