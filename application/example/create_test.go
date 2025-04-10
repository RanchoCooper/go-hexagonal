package example

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// MockExampleService mocks the example service for testing
type MockExampleService struct {
	mock.Mock
}

// Create implements the Create method
func (m *MockExampleService) Create(ctx context.Context, name, alias string) (*model.Example, error) {
	args := m.Called(ctx, name, alias)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// Get implements the Get method
func (m *MockExampleService) Get(ctx context.Context, id int) (*model.Example, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// FindByName implements the FindByName method
func (m *MockExampleService) FindByName(ctx context.Context, name string) (*model.Example, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// Delete implements the Delete method
func (m *MockExampleService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Update implements the Update method
func (m *MockExampleService) Update(ctx context.Context, id int, name, alias string) error {
	args := m.Called(ctx, id, name, alias)
	return args.Error(0)
}

// TestablCreateUseCase modifies CreateUseCase for testing purposes
type TestablCreateUseCase struct {
	CreateUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

// NewTestablCreateUseCase creates a testable create use case
func NewTestablCreateUseCase(svc *MockExampleService) *TestablCreateUseCase {
	return &TestablCreateUseCase{
		CreateUseCase: CreateUseCase{
			exampleService: svc,
		},
		txProvider: CreateTestTransaction,
	}
}

// Execute overrides the Execute method to replace transaction handling logic
func (uc *TestablCreateUseCase) Execute(ctx context.Context, in CreateInput) (*ExampleOutput, error) {
	// Validate input
	if err := in.Validate(); err != nil {
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// Use test transaction
	tx, err := uc.txProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Call domain service
	example, err := uc.exampleService.Create(ctx, in.Name, in.Alias)
	if err != nil {
		return nil, fmt.Errorf("failed to create example: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Build output
	out := NewExampleOutput(example)
	return out, nil
}

// TestCreateInput_Validate tests the validation logic for create input
func TestCreateInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: CreateInput{
				Name:  "Valid Example",
				Alias: "valid",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: CreateInput{
				Name:  "",
				Alias: "valid",
			},
			wantErr: true,
		},
		{
			name: "empty alias",
			input: CreateInput{
				Name:  "Valid Example",
				Alias: "",
			},
			wantErr: false, // alias can be empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCreateUseCase_Success tests the successful case of creating an example
func TestCreateUseCase_Success(t *testing.T) {
	// Prepare test data
	now := time.Now()
	example := &model.Example{
		Id:        1,
		Name:      "Test Example",
		Alias:     "test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create mock service
	mockSvc := new(MockExampleService)
	mockSvc.On("Create", mock.Anything, "Test Example", "test").Return(example, nil)

	// Create use case and execute
	uc := NewTestablCreateUseCase(mockSvc)
	result, err := uc.Execute(context.Background(), CreateInput{
		Name:  "Test Example",
		Alias: "test",
	})

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Test Example", result.Name)
	assert.Equal(t, "test", result.Alias)
	mockSvc.AssertExpectations(t)
}

// TestCreateUseCase_Error tests the error cases of creating an example
func TestCreateUseCase_Error(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateInput
		setupMock    func(*MockExampleService)
		txProvider   func(context.Context) (repo.Transaction, error)
		wantErr      bool
		errorMessage string
	}{
		{
			name: "input validation failed",
			input: CreateInput{
				Name:  "",
				Alias: "test",
			},
			setupMock:    func(mockSvc *MockExampleService) {},
			txProvider:   CreateTestTransaction,
			wantErr:      true,
			errorMessage: "input validation failed",
		},
		{
			name: "transaction creation failed",
			input: CreateInput{
				Name:  "Test Example",
				Alias: "test",
			},
			setupMock:    func(mockSvc *MockExampleService) {},
			txProvider:   ErrorTestTransaction,
			wantErr:      true,
			errorMessage: "failed to create transaction",
		},
		{
			name: "service error",
			input: CreateInput{
				Name:  "Test Example",
				Alias: "test",
			},
			setupMock: func(mockSvc *MockExampleService) {
				mockSvc.On("Create", mock.Anything, "Test Example", "test").
					Return(nil, errors.New("service error"))
			},
			txProvider:   CreateTestTransaction,
			wantErr:      true,
			errorMessage: "failed to create example",
		},
		{
			name: "transaction commit failed",
			input: CreateInput{
				Name:  "Test Example",
				Alias: "test",
			},
			setupMock: func(mockSvc *MockExampleService) {
				now := time.Now()
				example := &model.Example{
					Id:        1,
					Name:      "Test Example",
					Alias:     "test",
					CreatedAt: now,
					UpdatedAt: now,
				}
				mockSvc.On("Create", mock.Anything, "Test Example", "test").Return(example, nil)
			},
			txProvider:   CommitErrorTestTransaction,
			wantErr:      true,
			errorMessage: "failed to commit transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service and use case
			mockSvc := new(MockExampleService)
			tt.setupMock(mockSvc)

			uc := NewTestablCreateUseCase(mockSvc)
			uc.txProvider = tt.txProvider

			// Execute use case
			result, err := uc.Execute(context.Background(), tt.input)

			// Verify results
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

// TestNewExampleOutput tests the creation of output object
func TestNewExampleOutput(t *testing.T) {
	now := time.Now()
	example := &model.Example{
		Id:        1,
		Name:      "Test Example",
		Alias:     "test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	output := NewExampleOutput(example)

	assert.NotNil(t, output)
	assert.Equal(t, example.Id, output.ID)
	assert.Equal(t, example.Name, output.Name)
	assert.Equal(t, example.Alias, output.Alias)
	assert.Equal(t, example.CreatedAt, output.CreatedAt)
	assert.Equal(t, example.UpdatedAt, output.UpdatedAt)
}
