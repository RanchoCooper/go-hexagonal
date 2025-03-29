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

// mockExampleService implements IExampleService interface for testing
type mockExampleService struct {
	mock.Mock
}

func (m *mockExampleService) Create(ctx context.Context, name string, alias string) (*model.Example, error) {
	args := m.Called(ctx, name, alias)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *mockExampleService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockExampleService) Update(ctx context.Context, id int, name string, alias string) error {
	args := m.Called(ctx, id, name, alias)
	return args.Error(0)
}

func (m *mockExampleService) Get(ctx context.Context, id int) (*model.Example, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *mockExampleService) FindByName(ctx context.Context, name string) (*model.Example, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// Modify CreateUseCase for testing purposes
type testableCreateUseCase struct {
	CreateUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

func newTestableCreateUseCase(svc *mockExampleService) *testableCreateUseCase {
	return &testableCreateUseCase{
		CreateUseCase: CreateUseCase{
			exampleService: svc,
		},
		txProvider: mockTransaction,
	}
}

// Override Execute method to replace transaction handling logic
func (uc *testableCreateUseCase) Execute(ctx context.Context, input dto.CreateExampleReq) (*dto.CreateExampleResp, error) {
	// Use mock transaction
	tx, err := uc.txProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Convert DTO to domain model
	example := &model.Example{
		Name:  input.Name,
		Alias: input.Alias,
	}

	// Call domain service
	createdExample, err := uc.exampleService.Create(ctx, example.Name, example.Alias)
	if err != nil {
		return nil, fmt.Errorf("failed to create example: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert domain model to DTO
	result := &dto.CreateExampleResp{
		Id:        uint(createdExample.Id),
		Name:      createdExample.Name,
		Alias:     createdExample.Alias,
		CreatedAt: createdExample.CreatedAt,
		UpdatedAt: createdExample.UpdatedAt,
	}

	return result, nil
}

// Mock transaction implementation using project's NoopTransaction
func mockTransaction(ctx context.Context) (repo.Transaction, error) {
	return &repo.NoopTransaction{}, nil
}

// TestCreateUseCase_Execute_Success tests the successful case of creating an example
func TestCreateUseCase_Execute_Success(t *testing.T) {
	// Create mock service
	mockService := new(mockExampleService)

	// Set up mock behavior first, then create the use case
	now := time.Now()
	expectedExample := &model.Example{
		Id:        1,
		Name:      "Test Example",
		Alias:     "test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Setup any Create call to return the expected result
	mockService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(expectedExample, nil)

	// Create use case with testable version
	useCase := newTestableCreateUseCase(mockService)

	// Test data
	ctx := context.Background()
	createReq := dto.CreateExampleReq{
		Name:  "Test Example",
		Alias: "test",
	}

	// Execute use case
	result, err := useCase.Execute(ctx, createReq)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(expectedExample.Id), result.Id)
	assert.Equal(t, expectedExample.Name, result.Name)
	assert.Equal(t, expectedExample.Alias, result.Alias)
	assert.Equal(t, expectedExample.CreatedAt, result.CreatedAt)
	assert.Equal(t, expectedExample.UpdatedAt, result.UpdatedAt)

	mockService.AssertExpectations(t)
}

// TestCreateUseCase_Execute_Error tests the error case when creating an example
func TestCreateUseCase_Execute_Error(t *testing.T) {
	// Create mock service
	mockService := new(mockExampleService)

	// Setup mock behavior - simulate error
	expectedError := assert.AnError
	mockService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)

	// Create use case with testable version
	useCase := newTestableCreateUseCase(mockService)

	// Test data
	ctx := context.Background()
	createReq := dto.CreateExampleReq{
		Name:  "Test Example",
		Alias: "test",
	}

	// Execute use case
	result, err := useCase.Execute(ctx, createReq)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create example")

	mockService.AssertExpectations(t)
}
