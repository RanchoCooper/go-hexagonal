package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/domain/repo"
)

// testableDeleteUseCase is a testable implementation of the delete use case that replaces actual transaction handling
type testableDeleteUseCase struct {
	DeleteUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

// newTestableDeleteUseCase creates a testable delete use case
func newTestableDeleteUseCase(svc *MockExampleService) *testableDeleteUseCase {
	return &testableDeleteUseCase{
		DeleteUseCase: DeleteUseCase{
			exampleService: svc,
		},
		txProvider: CreateTestTransaction,
	}
}

// Execute overrides the Execute method to replace transaction handling logic
func (uc *testableDeleteUseCase) Execute(ctx context.Context, id int) error {
	// Use mock transaction
	tx, err := uc.txProvider(ctx)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

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

// TestDeleteUseCase_Execute_Success tests the successful case of deleting an example
func TestDeleteUseCase_Execute_Success(t *testing.T) {
	// Create mock service
	mockService := new(MockExampleService)

	// Setup mock behavior
	mockService.On("Delete", mock.Anything, 1).Return(nil)

	// Create testable use case
	useCase := newTestableDeleteUseCase(mockService)

	// Test data
	ctx := context.Background()
	exampleId := 1

	// Execute use case
	err := useCase.Execute(ctx, exampleId)

	// Verify results
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

// TestDeleteUseCase_Execute_Error tests the error case when deleting an example
func TestDeleteUseCase_Execute_Error(t *testing.T) {
	// Create mock service
	mockService := new(MockExampleService)

	// Setup mock behavior - simulate error
	expectedError := assert.AnError
	mockService.On("Delete", mock.Anything, 999).Return(expectedError)

	// Create testable use case
	useCase := newTestableDeleteUseCase(mockService)

	// Test data
	ctx := context.Background()
	exampleId := 999 // Non-existent ID

	// Execute use case
	err := useCase.Execute(ctx, exampleId)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete example")
	mockService.AssertExpectations(t)
}
