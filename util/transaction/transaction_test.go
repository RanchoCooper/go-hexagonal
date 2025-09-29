package transaction

import (
	"context"
	stdErrors "errors"
	"testing"

	"go-hexagonal/domain/repo"
	"go-hexagonal/util/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionFactory is a mock implementation of repo.TransactionFactory
type MockTransactionFactory struct {
	mock.Mock
}

func (m *MockTransactionFactory) NewTransaction(ctx context.Context, storeType repo.StoreType, opts any) (repo.Transaction, error) {
	args := m.Called(ctx, storeType, opts)
	return args.Get(0).(repo.Transaction), args.Error(1)
}

// MockTransaction is a mock implementation of repo.Transaction
type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) GetContext() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *MockTransaction) Begin() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) WithContext(ctx context.Context) repo.Transaction {
	args := m.Called(ctx)
	return args.Get(0).(repo.Transaction)
}

func (m *MockTransaction) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *MockTransaction) StoreType() repo.StoreType {
	args := m.Called()
	return args.Get(0).(repo.StoreType)
}

func (m *MockTransaction) Options() *repo.TransactionOptions {
	args := m.Called()
	return args.Get(0).(*repo.TransactionOptions)
}

func TestExecuteWithRollback_Success(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)
	mockTx.On("Commit").Return(nil)

	// Test function
	expectedResult := "test-result"
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return expectedResult, nil
	}

	// Execute
	result, err := ExecuteWithRollback(ctx, mockTxFactory, storeType, fn)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestExecuteWithRollback_TransactionCreationFailed(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)

	ctx := context.Background()
	storeType := repo.MySQLStore
	expectedErr := stdErrors.New("connection failed")

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return((*MockTransaction)(nil), expectedErr)

	// Test function
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return nil, nil
	}

	// Execute
	result, err := ExecuteWithRollback(ctx, mockTxFactory, storeType, fn)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create transaction")

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
}

func TestExecuteWithRollback_FunctionExecutionFailed(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore
	expectedErr := stdErrors.New("function failed")

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)

	// Test function
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return nil, expectedErr
	}

	// Execute
	result, err := ExecuteWithRollback(ctx, mockTxFactory, storeType, fn)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestExecuteWithRollback_CommitFailed(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore
	commitErr := stdErrors.New("commit failed")

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)
	mockTx.On("Commit").Return(commitErr)

	// Test function
	expectedResult := "test-result"
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return expectedResult, nil
	}

	// Execute
	result, err := ExecuteWithRollback(ctx, mockTxFactory, storeType, fn)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to commit transaction")

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestExecuteWithRollback_RollbackErrorIgnored(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore
	rollbackErr := stdErrors.New("rollback failed")

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(rollbackErr)

	// Test function - this will trigger rollback
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return nil, stdErrors.New("function failed")
	}

	// Execute
	result, err := ExecuteWithRollback(ctx, mockTxFactory, storeType, fn)

	// Verify - rollback error should be ignored, only function error should be returned
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "function failed", err.Error())

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestExecuteWithMetrics_Success(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore
	useCaseName := "test-use-case"

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)
	mockTx.On("Commit").Return(nil)

	// Test function
	expectedResult := "test-result"
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return expectedResult, nil
	}

	// Execute
	result, err := ExecuteWithMetrics(ctx, mockTxFactory, storeType, useCaseName, fn)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestExecuteWithMetrics_TransactionCreationFailed(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)

	ctx := context.Background()
	storeType := repo.MySQLStore
	useCaseName := "test-use-case"
	expectedErr := stdErrors.New("connection failed")

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return((*MockTransaction)(nil), expectedErr)

	// Test function
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return nil, nil
	}

	// Execute
	result, err := ExecuteWithMetrics(ctx, mockTxFactory, storeType, useCaseName, fn)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create transaction")

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
}

func TestExecuteWithMetrics_FunctionExecutionFailed(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore
	useCaseName := "test-use-case"
	expectedErr := stdErrors.New("function failed")

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)

	// Test function
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return nil, expectedErr
	}

	// Execute
	result, err := ExecuteWithMetrics(ctx, mockTxFactory, storeType, useCaseName, fn)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestExecuteWithMetrics_CommitFailed(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore
	useCaseName := "test-use-case"
	commitErr := stdErrors.New("commit failed")

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)
	mockTx.On("Commit").Return(commitErr)

	// Test function
	expectedResult := "test-result"
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return expectedResult, nil
	}

	// Execute
	result, err := ExecuteWithMetrics(ctx, mockTxFactory, storeType, useCaseName, fn)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to commit transaction")

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestWithBackgroundContext(t *testing.T) {
	ctx := WithBackgroundContext()

	assert.NotNil(t, ctx)
	// Should not be canceled
	select {
	case <-ctx.Done():
		t.Fatal("Background context should not be done")
	default:
		// Expected
	}
}

func TestWithTestContext(t *testing.T) {
	ctx := WithTestContext()

	assert.NotNil(t, ctx)
	// Should not be canceled
	select {
	case <-ctx.Done():
		t.Fatal("Test context should not be done")
	default:
		// Expected
	}
}

func TestTransactionErrorWrapping(t *testing.T) {
	// Test that system errors are properly wrapped
	originalErr := stdErrors.New("database error")

	// Simulate the error wrapping that happens in ExecuteWithRollback
	wrappedErr := errors.NewSystemError("failed to create transaction", originalErr)

	assert.Error(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "failed to create transaction")
	assert.Contains(t, wrappedErr.Error(), "database error")
}

func TestDifferentStoreTypes(t *testing.T) {
	testCases := []struct {
		name      string
		storeType repo.StoreType
	}{
		{"MySQL", repo.MySQLStore},
		{"PostgreSQL", repo.PostgresStore},
		{"SQLite", repo.MongoStore},
		{"Memory", repo.RedisStore},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			mockTxFactory := new(MockTransactionFactory)
			mockTx := new(MockTransaction)

			ctx := context.Background()

			// Mock expectations
			mockTxFactory.On("NewTransaction", ctx, tc.storeType, mock.Anything).Return(mockTx, nil)
			mockTx.On("Rollback").Return(nil)
			mockTx.On("Commit").Return(nil)

			// Test function
			fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
				return "success", nil
			}

			// Execute
			result, err := ExecuteWithRollback(ctx, mockTxFactory, tc.storeType, fn)

			// Verify
			assert.NoError(t, err)
			assert.Equal(t, "success", result)

			// Verify mock expectations
			mockTxFactory.AssertExpectations(t)
			mockTx.AssertExpectations(t)
		})
	}
}

func TestNilFunctionHandling(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)

	// This should not panic
	var fn func(context.Context, repo.Transaction) (any, error) = nil

	// Execute
	result, err := ExecuteWithRollback(ctx, mockTxFactory, storeType, fn)

	// Verify - should return an error
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "function parameter cannot be nil")

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestContextPropagation(t *testing.T) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	// Create a context with a value
	type contextKey string
	const testKey contextKey = "test-key"

	ctx := context.WithValue(context.Background(), testKey, "test-value")
	storeType := repo.MySQLStore

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)
	mockTx.On("Commit").Return(nil)

	// Test function that checks context propagation
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		value := ctx.Value(testKey)
		assert.Equal(t, "test-value", value)
		return value, nil
	}

	// Execute
	result, err := ExecuteWithRollback(ctx, mockTxFactory, storeType, fn)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "test-value", result)

	// Verify mock expectations
	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestTransactionOptions(t *testing.T) {
	// Test that transaction options are properly handled
	// Note: The current implementation doesn't use options, but we test the interface

	ctx := context.Background()
	storeType := repo.MySQLStore

	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	// Mock expectations with options - but note current implementation doesn't pass options
	// So we use mock.Anything to match any parameters
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)
	mockTx.On("Commit").Return(nil)

	// Test function
	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return "success", nil
	}

	// Note: Our current implementation doesn't support passing options
	// This test is for future enhancement

	result, err := ExecuteWithRollback(ctx, mockTxFactory, storeType, fn)

	assert.NoError(t, err)
	assert.Equal(t, "success", result)

	mockTxFactory.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkExecuteWithRollback(b *testing.B) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)
	mockTx.On("Commit").Return(nil)

	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return "benchmark-result", nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ExecuteWithRollback(ctx, mockTxFactory, storeType, fn)
	}
}

func BenchmarkExecuteWithMetrics(b *testing.B) {
	// Setup mocks
	mockTxFactory := new(MockTransactionFactory)
	mockTx := new(MockTransaction)

	ctx := context.Background()
	storeType := repo.MySQLStore
	useCaseName := "benchmark-use-case"

	// Mock expectations
	mockTxFactory.On("NewTransaction", ctx, storeType, mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)
	mockTx.On("Commit").Return(nil)

	fn := func(ctx context.Context, tx repo.Transaction) (any, error) {
		return "benchmark-result", nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ExecuteWithMetrics(ctx, mockTxFactory, storeType, useCaseName, fn)
	}
}
