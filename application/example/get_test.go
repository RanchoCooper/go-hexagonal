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

// MockExampleService is defined in create_test.go

// TestablGetUseCase modifies GetUseCase for testing purposes
type TestablGetUseCase struct {
	GetUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

func NewTestablGetUseCase(svc *MockExampleService) *TestablGetUseCase {
	return &TestablGetUseCase{
		GetUseCase: GetUseCase{
			exampleService: svc,
		},
		txProvider: CreateTestTransaction,
	}
}

// Execute overrides Execute method to replace transaction handling logic
func (uc *TestablGetUseCase) Execute(ctx context.Context, id int) (*dto.GetExampleResponse, error) {
	// Use test transaction
	tx, err := uc.txProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

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

// TestGetUseCase_Success tests successful retrieval of example by ID
func TestGetUseCase_Success(t *testing.T) {
	// Create mock service
	mockService := new(MockExampleService)

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

	// Set mock behavior
	mockService.On("Get", mock.Anything, exampleId).Return(expectedExample, nil)

	// Create use case with testable version
	useCase := NewTestablGetUseCase(mockService)

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

// TestGetUseCase_Error tests error scenario when retrieving example by ID
func TestGetUseCase_Error(t *testing.T) {
	// Create mock service
	mockService := new(MockExampleService)

	// Test data
	exampleId := 999 // Non-existent ID

	// Set mock behavior - simulate error
	expectedError := assert.AnError
	mockService.On("Get", mock.Anything, exampleId).Return(nil, expectedError)

	// Create use case with testable version
	useCase := NewTestablGetUseCase(mockService)

	// Execute use case
	ctx := context.Background()
	result, err := useCase.Execute(ctx, exampleId)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get example")

	mockService.AssertExpectations(t)
}
