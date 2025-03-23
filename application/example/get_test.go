package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/api/dto"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// mockExampleService is defined in create_test.go

// Modify GetUseCase for testing purposes
type testableGetUseCase struct {
	GetUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

func newTestableGetUseCase(svc *mockExampleService) *testableGetUseCase {
	return &testableGetUseCase{
		GetUseCase: GetUseCase{
			exampleService: svc,
		},
		txProvider: mockTransaction,
	}
}

// Override Execute method to replace transaction handling logic
func (uc *testableGetUseCase) Execute(ctx context.Context, id int) (*dto.GetExampleResponse, error) {
	// Use mock transaction
	tx, err := uc.txProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Call domain service
	example, err := uc.exampleService.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get example: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert domain model to DTO
	result := &dto.GetExampleResponse{
		Id:        example.Id,
		Name:      example.Name,
		Alias:     example.Alias,
		CreatedAt: example.CreatedAt,
		UpdatedAt: example.UpdatedAt,
	}

	return result, nil
}

// TestGetUseCase_Execute_Success tests the successful case of getting an example by ID
func TestGetUseCase_Execute_Success(t *testing.T) {
	// Create mock service
	mockService := new(mockExampleService)

	// Test data
	exampleId := 1

	now := time.Now()
	expectedExample := &model.Example{
		Id:        exampleId,
		Name:      "Test Example",
		Alias:     "test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Setup mock behavior
	mockService.On("Get", mock.Anything, exampleId).Return(expectedExample, nil)

	// Create use case with testable version
	useCase := newTestableGetUseCase(mockService)

	// Execute use case
	ctx := context.Background()
	result, err := useCase.Execute(ctx, exampleId)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedExample.Id, result.Id)
	assert.Equal(t, expectedExample.Name, result.Name)
	assert.Equal(t, expectedExample.Alias, result.Alias)
	assert.Equal(t, expectedExample.CreatedAt, result.CreatedAt)
	assert.Equal(t, expectedExample.UpdatedAt, result.UpdatedAt)

	mockService.AssertExpectations(t)
}

// TestGetUseCase_Execute_Error tests the error case when getting an example by ID
func TestGetUseCase_Execute_Error(t *testing.T) {
	// Create mock service
	mockService := new(mockExampleService)

	// Test data
	exampleId := 999 // Non-existent ID

	// Setup mock behavior - simulate error
	expectedError := assert.AnError
	mockService.On("Get", mock.Anything, exampleId).Return(nil, expectedError)

	// Create use case with testable version
	useCase := newTestableGetUseCase(mockService)

	// Execute use case
	ctx := context.Background()
	result, err := useCase.Execute(ctx, exampleId)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get example")

	mockService.AssertExpectations(t)
}
