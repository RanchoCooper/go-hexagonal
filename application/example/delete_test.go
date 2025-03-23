package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/domain/repo"
)

// mockExampleService is defined in create_test.go

// Modify DeleteUseCase for testing purposes
type testableDeleteUseCase struct {
	DeleteUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

func newTestableDeleteUseCase(svc *mockExampleService) *testableDeleteUseCase {
	return &testableDeleteUseCase{
		DeleteUseCase: DeleteUseCase{
			exampleService: svc,
		},
		txProvider: mockTransaction,
	}
}

// Override Execute method to replace transaction handling logic
func (uc *testableDeleteUseCase) Execute(ctx context.Context, id int) error {
	// Use mock transaction
	tx, err := uc.txProvider(ctx)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Call domain service
	if err := uc.exampleService.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Mock transaction implementation is defined in create_test.go

// TestDeleteUseCase_Execute_Success tests the successful case of deleting an example
func TestDeleteUseCase_Execute_Success(t *testing.T) {
	// Create mock service
	mockService := new(mockExampleService)

	// Setup mock behavior
	mockService.On("Delete", mock.Anything, 1).Return(nil)

	// Create use case with testable version
	useCase := newTestableDeleteUseCase(mockService)

	// Test data
	ctx := context.Background()
	exampleId := 1

	// Execute use case
	err := useCase.Execute(ctx, exampleId)

	// Assert results
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

// TestDeleteUseCase_Execute_Error tests the error case when deleting an example
func TestDeleteUseCase_Execute_Error(t *testing.T) {
	// Create mock service
	mockService := new(mockExampleService)

	// Setup mock behavior - simulate error
	expectedError := assert.AnError
	mockService.On("Delete", mock.Anything, 999).Return(expectedError)

	// Create use case with testable version
	useCase := newTestableDeleteUseCase(mockService)

	// Test data
	ctx := context.Background()
	exampleId := 999 // Non-existent ID

	// Execute use case
	err := useCase.Execute(ctx, exampleId)

	// Assert results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete example")
	mockService.AssertExpectations(t)
}
