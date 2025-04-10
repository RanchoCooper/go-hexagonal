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

// TestableFindByNameUseCase is a testable implementation of the query use case that replaces actual transaction handling
type TestableFindByNameUseCase struct {
	FindByNameUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

// NewTestableFindByNameUseCase creates a testable query use case
func NewTestableFindByNameUseCase(svc *MockExampleService) *TestableFindByNameUseCase {
	return &TestableFindByNameUseCase{
		FindByNameUseCase: FindByNameUseCase{
			exampleService: svc,
		},
		txProvider: CreateTestTransaction,
	}
}

// Execute overrides the method to replace transaction handling logic
func (uc *TestableFindByNameUseCase) Execute(ctx context.Context, name string) (*dto.GetExampleResponse, error) {
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
		return nil, fmt.Errorf("failed to find example: %w", err)
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

// TestFindByNameUseCase_Success tests the successful case of finding an example by name
func TestFindByNameUseCase_Success(t *testing.T) {
	// Create mock service
	mockService := new(MockExampleService)

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
	useCase := NewTestableFindByNameUseCase(mockService)

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

// TestFindByNameUseCase_NotFound tests the case when an example is not found
func TestFindByNameUseCase_NotFound(t *testing.T) {
	// Create mock service
	mockService := new(MockExampleService)

	// Setup mock behavior for not found case
	mockService.On("FindByName", mock.Anything, "non-existent").Return(nil, repo.ErrNotFound)

	// Create use case with testable version
	useCase := NewTestableFindByNameUseCase(mockService)

	// Test data
	ctx := context.Background()
	name := "non-existent"

	// Execute use case
	result, err := useCase.Execute(ctx, name)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to find example")
	mockService.AssertExpectations(t)
}
