package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		error    *AppError
		expected string
	}{
		{
			name: "error without cause",
			error: &AppError{
				Type:    ErrorTypeValidation,
				Message: "invalid input",
			},
			expected: "VALIDATION: invalid input",
		},
		{
			name: "error with cause",
			error: &AppError{
				Type:    ErrorTypeValidation,
				Message: "invalid input",
				Cause:   errors.New("underlying error"),
			},
			expected: "VALIDATION: invalid input (caused by: underlying error)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.error.Error())
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &AppError{
		Type:    ErrorTypeValidation,
		Message: "invalid input",
		Cause:   cause,
	}

	assert.Equal(t, cause, err.Unwrap())
}

func TestAppError_WithDetails(t *testing.T) {
	err := NewValidationError("invalid input", nil)
	details := map[string]any{
		"field": "name",
		"value": "",
	}

	errWithDetails := err.WithDetails(details)
	assert.Equal(t, details, errWithDetails.Details)
	assert.Equal(t, err.Type, errWithDetails.Type)
	assert.Equal(t, err.Message, errWithDetails.Message)
}

func TestAppError_WithCode(t *testing.T) {
	err := NewValidationError("invalid input", nil)
	errWithCode := err.WithCode(400)

	assert.Equal(t, 400, errWithCode.Code)
	assert.Equal(t, err.Type, errWithCode.Type)
	assert.Equal(t, err.Message, errWithCode.Message)
}

func TestNewValidationError(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewValidationError("invalid input", cause)

	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "invalid input", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("resource not found", nil)

	assert.Equal(t, ErrorTypeNotFound, err.Type)
	assert.Equal(t, "resource not found", err.Message)
}

func TestNewPersistenceError(t *testing.T) {
	err := NewPersistenceError("database error", nil)

	assert.Equal(t, ErrorTypePersistence, err.Type)
	assert.Equal(t, "database error", err.Message)
}

func TestNewSystemError(t *testing.T) {
	err := NewSystemError("system error", nil)

	assert.Equal(t, ErrorTypeSystem, err.Type)
	assert.Equal(t, "system error", err.Message)
}

func TestNewBusinessError(t *testing.T) {
	err := NewBusinessError("business error", nil)

	assert.Equal(t, ErrorTypeBusiness, err.Type)
	assert.Equal(t, "business error", err.Message)
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "validation error",
			error:    NewValidationError("invalid", nil),
			expected: true,
		},
		{
			name:     "not found error",
			error:    NewNotFoundError("not found", nil),
			expected: false,
		},
		{
			name:     "standard error",
			error:    errors.New("standard error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsValidationError(tt.error))
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "not found error",
			error:    NewNotFoundError("not found", nil),
			expected: true,
		},
		{
			name:     "validation error",
			error:    NewValidationError("invalid", nil),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsNotFoundError(tt.error))
		})
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := Wrap(cause, ErrorTypeValidation, "wrapped error")

	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "wrapped error", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestWrapf(t *testing.T) {
	cause := errors.New("underlying error")
	err := Wrapf(cause, ErrorTypeValidation, "wrapped %s", "error")

	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "wrapped error", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestNew(t *testing.T) {
	err := New(ErrorTypeValidation, "new error")

	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "new error", err.Message)
	assert.Nil(t, err.Cause)
}

func TestNewf(t *testing.T) {
	err := Newf(ErrorTypeValidation, "new %s", "error")

	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "new error", err.Message)
	assert.Nil(t, err.Cause)
}

func TestStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		error    *AppError
		expected int
	}{
		{
			name:     "validation error",
			error:    NewValidationError("invalid", nil),
			expected: 400,
		},
		{
			name:     "not found error",
			error:    NewNotFoundError("not found", nil),
			expected: 404,
		},
		{
			name:     "unauthorized error",
			error:    &AppError{Type: ErrorTypeUnauthorized},
			expected: 401,
		},
		{
			name:     "forbidden error",
			error:    &AppError{Type: ErrorTypeForbidden},
			expected: 403,
		},
		{
			name:     "conflict error",
			error:    &AppError{Type: ErrorTypeConflict},
			expected: 409,
		},
		{
			name:     "default error",
			error:    &AppError{Type: ErrorTypeSystem},
			expected: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.error.StatusCode())
		})
	}
}

func TestIsPersistenceError(t *testing.T) {
	err := NewPersistenceError("db error", nil)
	assert.True(t, IsPersistenceError(err))
	assert.False(t, IsPersistenceError(NewValidationError("invalid", nil)))
}

func TestIsSystemError(t *testing.T) {
	err := NewSystemError("system error", nil)
	assert.True(t, IsSystemError(err))
	assert.False(t, IsSystemError(NewValidationError("invalid", nil)))
}

func TestIsBusinessError(t *testing.T) {
	err := NewBusinessError("business error", nil)
	assert.True(t, IsBusinessError(err))
	assert.False(t, IsBusinessError(NewValidationError("invalid", nil)))
}
