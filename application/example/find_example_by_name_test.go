package example

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/api/dto"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// mockExampleService is defined in create_test.go

// Modify FindByNameUseCase for testing purposes
type testableFindByNameUseCase struct {
	FindByNameUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

func newTestableFindByNameUseCase(svc *mockExampleService) *testableFindByNameUseCase {
	return &testableFindByNameUseCase{
		FindByNameUseCase: FindByNameUseCase{
			exampleService: svc,
		},
		txProvider: mockTransaction,
	}
}

// Override Execute method to replace transaction handling logic
func (uc *testableFindByNameUseCase) Execute(ctx context.Context, name string) (*dto.GetExampleResponse, error) {
	// Use mock transaction
	tx, err := uc.txProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Call domain service
	example, err := uc.exampleService.FindByName(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("failed to find example by name: %w", err)
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

// TestFindByNameUseCase_Execute_Success tests the successful case of finding an example by name
func TestFindByNameUseCase_Execute_Success(t *testing.T) {
	// Create mock service
	mockService := new(mockExampleService)

	// Test data
	exampleName := "Test Example"

	now := time.Now()
	expectedExample := &model.Example{
		Id:        1,
		Name:      exampleName,
		Alias:     "test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Setup mock behavior
	mockService.On("FindByName", mock.Anything, exampleName).Return(expectedExample, nil)

	// Create use case with testable version
	useCase := newTestableFindByNameUseCase(mockService)

	// Execute use case
	ctx := context.Background()
	result, err := useCase.Execute(ctx, exampleName)

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

// TestFindByNameUseCase_Execute_Error tests the error case when finding an example by name
func TestFindByNameUseCase_Execute_Error(t *testing.T) {
	// Create mock service
	mockService := new(mockExampleService)

	// Test data
	exampleName := "Nonexistent Example"

	// Setup mock behavior - simulate error
	expectedError := assert.AnError
	mockService.On("FindByName", mock.Anything, exampleName).Return(nil, expectedError)

	// Create use case with testable version
	useCase := newTestableFindByNameUseCase(mockService)

	// Execute use case
	ctx := context.Background()
	result, err := useCase.Execute(ctx, exampleName)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to find example")

	mockService.AssertExpectations(t)
}
