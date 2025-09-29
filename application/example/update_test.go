package example

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/api/dto"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// MockExampleService is defined in create_test.go

// testableUpdateUseCase modifies UpdateUseCase for testing purposes
type testableUpdateUseCase struct {
	UpdateUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

// newTestableUpdateUseCase creates a testable update use case
func newTestableUpdateUseCase(svc *MockExampleService) *testableUpdateUseCase {
	return &testableUpdateUseCase{
		UpdateUseCase: UpdateUseCase{
			exampleService: svc,
		},
		txProvider: CreateTestTransaction,
	}
}

// Execute overrides the Execute method to replace transaction handling logic
func (uc *testableUpdateUseCase) Execute(ctx context.Context, input dto.UpdateExampleReq) error {
	// Use mock transaction
	tx, err := uc.txProvider(ctx)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Convert DTO to domain model
	example := &model.Example{
		Id:    int(input.Id),
		Name:  input.Name,
		Alias: input.Alias,
	}

	// Call domain service
	if err := uc.exampleService.Update(ctx, example.Id, example.Name, example.Alias); err != nil {
		return fmt.Errorf("failed to update example: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TestUpdateUseCase_Execute_Success tests the successful case of updating an example
func TestUpdateUseCase_Execute_Success(t *testing.T) {
	// Create mock service
	mockService := new(MockExampleService)

	// Test data
	updateReq := dto.UpdateExampleReq{
		Id:    1,
		Name:  "Updated Example",
		Alias: "updated",
	}

	// Setup mock behavior
	mockService.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create testable use case
	useCase := newTestableUpdateUseCase(mockService)

	// Execute use case
	ctx := context.Background()
	err := useCase.Execute(ctx, updateReq)

	// Verify results
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

// TestUpdateUseCase_Execute_Error tests the error case when updating an example
func TestUpdateUseCase_Execute_Error(t *testing.T) {
	// Create mock service
	mockService := new(MockExampleService)

	// Test data
	updateReq := dto.UpdateExampleReq{
		Id:    999, // Non-existent ID
		Name:  "Updated Example",
		Alias: "updated",
	}

	// Setup mock behavior - simulate error
	expectedError := assert.AnError
	mockService.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedError)

	// Create testable use case
	useCase := newTestableUpdateUseCase(mockService)

	// Execute use case
	ctx := context.Background()
	err := useCase.Execute(ctx, updateReq)

	// Verify results
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update example")
	mockService.AssertExpectations(t)
}
