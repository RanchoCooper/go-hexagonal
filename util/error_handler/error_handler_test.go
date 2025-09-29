package error_handler

import (
	"context"
	"fmt"
	"testing"

	"go-hexagonal/api/error_code"
	util_errors "go-hexagonal/util/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock implementation of the logger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Errorf(format string, args ...any) {
	m.Called(format, args)
}

func (m *MockLogger) Warnf(format string, args ...any) {
	m.Called(format, args)
}

func (m *MockLogger) Infof(format string, args ...any) {
	m.Called(format, args)
}

func TestNew(t *testing.T) {
	handler := New()
	assert.NotNil(t, handler)
	assert.IsType(t, &ErrorHandler{}, handler)
}

func TestHandleError_NilError(t *testing.T) {
	handler := New()
	ctx := context.Background()

	err := handler.HandleError(ctx, nil, "test-operation")
	assert.Nil(t, err)
}

func TestHandleError_WithError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	expectedErr := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	err := handler.HandleError(ctx, expectedErr, "test-operation")
	assert.Equal(t, expectedErr, err)
}

func TestHandleAndWrapError_NilError(t *testing.T) {
	handler := New()
	ctx := context.Background()

	err := handler.HandleAndWrapError(ctx, nil, "test-operation", "test message")
	assert.Nil(t, err)
}

func TestHandleAndWrapError_WithError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	expectedErr := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	err := handler.HandleAndWrapError(ctx, expectedErr, "test-operation", "wrapped")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wrapped")
	assert.Contains(t, err.Error(), "test error")
}

func TestHandleAndConvertError_NilError(t *testing.T) {
	handler := New()
	ctx := context.Background()

	err := handler.HandleAndConvertError(ctx, nil, "test-operation", util_errors.ErrorTypeValidation)
	assert.Nil(t, err)
}

func TestHandleAndConvertError_ValidationError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	expectedErr := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	err := handler.HandleAndConvertError(ctx, expectedErr, "test-operation", util_errors.ErrorTypeValidation)
	assert.NotNil(t, err)
	assert.True(t, util_errors.IsValidationError(err))
	assert.Contains(t, err.Error(), "test-operation failed")
}

func TestHandleAndConvertError_NotFoundError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	expectedErr := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	err := handler.HandleAndConvertError(ctx, expectedErr, "test-operation", util_errors.ErrorTypeNotFound)
	assert.NotNil(t, err)
	assert.True(t, util_errors.IsNotFoundError(err))
	assert.Contains(t, err.Error(), "test-operation failed")
}

func TestHandleAndConvertError_PersistenceError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	expectedErr := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	err := handler.HandleAndConvertError(ctx, expectedErr, "test-operation", util_errors.ErrorTypePersistence)
	assert.NotNil(t, err)
	assert.True(t, util_errors.IsPersistenceError(err))
	assert.Contains(t, err.Error(), "test-operation failed")
}

func TestHandleAndConvertError_SystemError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	expectedErr := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	err := handler.HandleAndConvertError(ctx, expectedErr, "test-operation", util_errors.ErrorTypeSystem)
	assert.NotNil(t, err)
	assert.True(t, util_errors.IsSystemError(err))
	assert.Contains(t, err.Error(), "test-operation failed")
}

func TestHandleAndConvertError_BusinessError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	expectedErr := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	err := handler.HandleAndConvertError(ctx, expectedErr, "test-operation", util_errors.ErrorTypeBusiness)
	assert.NotNil(t, err)
	assert.True(t, util_errors.IsBusinessError(err))
	assert.Contains(t, err.Error(), "test-operation failed")
}

func TestHandleDomainError_NilError(t *testing.T) {
	handler := New()
	ctx := context.Background()

	err := handler.HandleDomainError(ctx, nil, "test-operation", "user")
	assert.Nil(t, err)
}

func TestHandleDomainError_WithError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	expectedErr := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	err := handler.HandleDomainError(ctx, expectedErr, "test-operation", "user")
	assert.NotNil(t, err)
	assert.True(t, util_errors.IsBusinessError(err))
	assert.Contains(t, err.Error(), "test-operation error in user domain")
}

func TestHandleAPIError_NilError(t *testing.T) {
	handler := New()
	ctx := context.Background()

	err := handler.HandleAPIError(ctx, nil, "test-operation")
	assert.Nil(t, err)
}

func TestHandleAPIError_ValidationError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	err := util_errors.NewValidationError("invalid input", nil)

	apiErr := handler.HandleAPIError(ctx, err, "test-operation")
	assert.NotNil(t, apiErr)
	assert.Equal(t, error_code.InvalidParams.Code, apiErr.Code)
}

func TestHandleAPIError_NotFoundError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	err := util_errors.NewNotFoundError("not found", nil)

	apiErr := handler.HandleAPIError(ctx, err, "test-operation")
	assert.NotNil(t, apiErr)
	assert.Equal(t, error_code.NotFound.Code, apiErr.Code)
}

func TestHandleAPIError_PersistenceError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	err := util_errors.NewPersistenceError("db error", nil)

	apiErr := handler.HandleAPIError(ctx, err, "test-operation")
	assert.NotNil(t, apiErr)
	assert.Equal(t, error_code.ServerError.Code, apiErr.Code)
}

func TestHandleAPIError_SystemError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	err := util_errors.NewSystemError("system error", nil)

	apiErr := handler.HandleAPIError(ctx, err, "test-operation")
	assert.NotNil(t, apiErr)
	assert.Equal(t, error_code.ServerError.Code, apiErr.Code)
}

func TestHandleAPIError_BusinessError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	err := util_errors.NewBusinessError("business error", nil)

	apiErr := handler.HandleAPIError(ctx, err, "test-operation")
	assert.NotNil(t, apiErr)
	assert.Equal(t, error_code.ServerError.Code, apiErr.Code)
}

func TestHandleAPIError_DefaultError(t *testing.T) {
	handler := New()
	ctx := context.Background()
	err := util_errors.New(util_errors.ErrorTypeSystem, "unknown error")

	apiErr := handler.HandleAPIError(ctx, err, "test-operation")
	assert.NotNil(t, apiErr)
	assert.Equal(t, error_code.ServerError.Code, apiErr.Code)
}

func TestErrorTypeDetection(t *testing.T) {
	handler := New()
	ctx := context.Background()

	testCases := []struct {
		name        string
		error       error
		expectedAPI bool
	}{
		{"Validation Error", util_errors.NewValidationError("validation", nil), true},
		{"Not Found Error", util_errors.NewNotFoundError("not found", nil), true},
		{"Persistence Error", util_errors.NewPersistenceError("persistence", nil), true},
		{"System Error", util_errors.NewSystemError("system", nil), true},
		{"Business Error", util_errors.NewBusinessError("business", nil), true},
		{"Standard Error", util_errors.New(util_errors.ErrorTypeSystem, "standard"), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiErr := handler.HandleAPIError(ctx, tc.error, "test-operation")
			if tc.expectedAPI {
				assert.NotNil(t, apiErr)
				assert.NotEmpty(t, apiErr.Code)
				assert.NotEmpty(t, apiErr.Msg)
			} else {
				assert.Nil(t, apiErr)
			}
		})
	}
}

func TestErrorWrappingPreservation(t *testing.T) {
	handler := New()
	ctx := context.Background()

	// Create a wrapped error
	originalErr := util_errors.New(util_errors.ErrorTypeSystem, "original error")
	wrappedErr := util_errors.Wrapf(originalErr, util_errors.ErrorTypeValidation, "wrapped message")

	// Handle the wrapped error
	handledErr := handler.HandleAndWrapError(ctx, wrappedErr, "test-operation", "additional wrap")

	assert.Error(t, handledErr)
	assert.Contains(t, handledErr.Error(), "additional wrap")
	assert.Contains(t, handledErr.Error(), "wrapped message")
	assert.Contains(t, handledErr.Error(), "original error")
}

func TestContextUsage(t *testing.T) {
	handler := New()

	// Create a context with values
	type contextKey string
	const testKey contextKey = "test-key"

	ctx := context.WithValue(context.Background(), testKey, "test-value")
	err := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	// Test that context is passed through
	handledErr := handler.HandleError(ctx, err, "test-operation")
	assert.Equal(t, err, handledErr)

	// Test with wrapped error
	wrappedErr := handler.HandleAndWrapError(ctx, err, "test-operation", "wrapped")
	assert.Error(t, wrappedErr)

	// Test with converted error
	convertedErr := handler.HandleAndConvertError(ctx, err, "test-operation", util_errors.ErrorTypeValidation)
	assert.Error(t, convertedErr)
}

func TestNilHandlerSafety(t *testing.T) {
	// Test that methods handle nil receiver gracefully
	var handler *ErrorHandler = nil
	ctx := context.Background()
	err := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	// These should not panic
	assert.NotPanics(t, func() {
		_ = handler.HandleError(ctx, err, "test-operation")
	})

	assert.NotPanics(t, func() {
		_ = handler.HandleAndWrapError(ctx, err, "test-operation", "wrapped")
	})

	assert.NotPanics(t, func() {
		_ = handler.HandleAndConvertError(ctx, err, "test-operation", util_errors.ErrorTypeValidation)
	})

	assert.NotPanics(t, func() {
		_ = handler.HandleDomainError(ctx, err, "test-operation", "user")
	})

	assert.NotPanics(t, func() {
		_ = handler.HandleAPIError(ctx, err, "test-operation")
	})
}

func TestEmptyOperationName(t *testing.T) {
	handler := New()
	ctx := context.Background()
	err := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	// Test with empty operation name
	handledErr := handler.HandleError(ctx, err, "")
	assert.Equal(t, err, handledErr)

	wrappedErr := handler.HandleAndWrapError(ctx, err, "", "wrapped")
	assert.Error(t, wrappedErr)

	convertedErr := handler.HandleAndConvertError(ctx, err, "", util_errors.ErrorTypeValidation)
	assert.Error(t, convertedErr)
}

func TestConcurrentAccess(t *testing.T) {
	handler := New()
	ctx := context.Background()
	err := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	// Test concurrent access to handler methods
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			operation := fmt.Sprintf("operation-%d", index)

			_ = handler.HandleError(ctx, err, operation)
			_ = handler.HandleAndWrapError(ctx, err, operation, "wrapped")
			_ = handler.HandleAndConvertError(ctx, err, operation, util_errors.ErrorTypeValidation)
			_ = handler.HandleDomainError(ctx, err, operation, "user")
			_ = handler.HandleAPIError(ctx, err, operation)

			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// Benchmark tests
func BenchmarkHandleError(b *testing.B) {
	handler := New()
	ctx := context.Background()
	err := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.HandleError(ctx, err, "benchmark-operation")
	}
}

func BenchmarkHandleAndWrapError(b *testing.B) {
	handler := New()
	ctx := context.Background()
	err := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.HandleAndWrapError(ctx, err, "benchmark-operation", "wrapped message")
	}
}

func BenchmarkHandleAndConvertError(b *testing.B) {
	handler := New()
	ctx := context.Background()
	err := util_errors.New(util_errors.ErrorTypeSystem, "test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.HandleAndConvertError(ctx, err, "benchmark-operation", util_errors.ErrorTypeValidation)
	}
}

func BenchmarkHandleAPIError(b *testing.B) {
	handler := New()
	ctx := context.Background()
	err := util_errors.NewValidationError("validation error", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.HandleAPIError(ctx, err, "benchmark-operation")
	}
}
