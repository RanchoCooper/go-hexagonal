package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/domain/event"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// Create Mock repository
type MockExampleRepo struct {
	mock.Mock
}

func (m *MockExampleRepo) Create(ctx context.Context, tr repo.Transaction, example *model.Example) (*model.Example, error) {
	// Set Id for parameter example so it has correct Id for event generation
	if example.Id == 0 {
		example.Id = 1 // Set a default ID
	}

	args := m.Called(ctx, tr, example)

	// If mock is configured to return an Example, ensure using it as return value
	if e, ok := args.Get(0).(*model.Example); ok {
		return e, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockExampleRepo) Delete(ctx context.Context, tr repo.Transaction, id int) error {
	args := m.Called(ctx, tr, id)
	return args.Error(0)
}

func (m *MockExampleRepo) Update(ctx context.Context, tr repo.Transaction, entity *model.Example) error {
	args := m.Called(ctx, tr, entity)
	return args.Error(0)
}

func (m *MockExampleRepo) GetByID(ctx context.Context, tr repo.Transaction, Id int) (*model.Example, error) {
	args := m.Called(ctx, tr, Id)
	if e, ok := args.Get(0).(*model.Example); ok {
		return e, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExampleRepo) FindByName(ctx context.Context, tr repo.Transaction, name string) (*model.Example, error) {
	args := m.Called(ctx, tr, name)
	if e, ok := args.Get(0).(*model.Example); ok {
		return e, args.Error(1)
	}
	return nil, args.Error(1)
}

// Create Mock cache repository
type MockExampleCacheRepo struct {
	mock.Mock
}

func (m *MockExampleCacheRepo) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockExampleCacheRepo) GetByID(ctx context.Context, id int) (*model.Example, error) {
	args := m.Called(ctx, id)
	if e, ok := args.Get(0).(*model.Example); ok {
		return e, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExampleCacheRepo) GetByName(ctx context.Context, name string) (*model.Example, error) {
	args := m.Called(ctx, name)
	if e, ok := args.Get(0).(*model.Example); ok {
		return e, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExampleCacheRepo) Set(ctx context.Context, example *model.Example) error {
	args := m.Called(ctx, example)
	return args.Error(0)
}

func (m *MockExampleCacheRepo) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExampleCacheRepo) Invalidate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Create Mock event bus
type MockEventBus struct {
	mock.Mock
}

// Publish implements EventBus interface's Publish method
func (m *MockEventBus) Publish(ctx context.Context, event event.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Subscribe implements EventBus interface's Subscribe method
func (m *MockEventBus) Subscribe(handler event.EventHandler) {
	m.Called(handler)
}

// Unsubscribe implements EventBus interface's Unsubscribe method
func (m *MockEventBus) Unsubscribe(handler event.EventHandler) {
	m.Called(handler)
}

// Create Mock transaction object
type MockTransaction struct {
	mock.Mock
}

// Define a helper function for testing to set EventBus
func withEventBus(service *ExampleService, bus event.EventBus) *ExampleService {
	service.EventBus = bus
	return service
}

// Test ExampleService creation
func TestNewExampleService(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)
	mockEventBus := new(MockEventBus)

	// Test creation without optional parameters
	service := NewExampleService(mockRepo, nil)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Repository)
	assert.Nil(t, service.CacheRepo)
	assert.Nil(t, service.EventBus)

	// Test creation with cache
	service = NewExampleService(mockRepo, mockCacheRepo)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Repository)
	assert.Equal(t, mockCacheRepo, service.CacheRepo)
	assert.Nil(t, service.EventBus)

	// Test setting event bus
	service = NewExampleService(mockRepo, mockCacheRepo)
	service.EventBus = mockEventBus
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Repository)
	assert.Equal(t, mockCacheRepo, service.CacheRepo)
	assert.Equal(t, mockEventBus, service.EventBus)

	// Test using helper function to set event bus
	service = withEventBus(NewExampleService(mockRepo, mockCacheRepo), mockEventBus)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Repository)
	assert.Equal(t, mockCacheRepo, service.CacheRepo)
	assert.Equal(t, mockEventBus, service.EventBus)
}

// TestExampleService_Create tests Create method
func TestExampleService_Create(t *testing.T) {
	// Prepare
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)
	mockEventBus := new(MockEventBus)

	service := NewExampleService(mockRepo, mockCacheRepo)
	service.EventBus = mockEventBus

	// Create a correct Example instance to get events
	input, err := model.NewExample("Test", "test-alias")
	assert.NoError(t, err)
	input.Id = 1 // Ensure ID is set

	// Mock dependency behavior
	mockRepo.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Run(func(args mock.Arguments) {
		// When Create is called, add events to created object and ensure ID is set
		example := args.Get(2).(*model.Example)
		example.Id = 1
	}).Return(input, nil) // Return object with events

	mockCacheRepo.On("Set", mock.Anything, mock.Anything).Return(nil)
	mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("event.ExampleCreatedEvent")).Return(nil)

	// Execute
	result, err := service.Create(context.Background(), "Test", "test-alias")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test", result.Name)
	assert.Equal(t, "test-alias", result.Alias)

	// Verify all mock calls
	mockRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// Test Delete method
func TestExampleService_Delete(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)
	mockEventBus := new(MockEventBus)

	// Create service instance
	service := NewExampleService(mockRepo, mockCacheRepo)
	service.EventBus = mockEventBus

	testCases := []struct {
		name        string
		setupMocks  func()
		exampleId   int
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Successfully delete example",
			setupMocks: func() {
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 1).Return(&model.Example{
					Id:    1,
					Name:  "Test Example",
					Alias: "test-alias",
				}, nil)
				mockRepo.On("Delete", mock.Anything, mock.Anything, 1).Return(nil)
				mockCacheRepo.On("Delete", mock.Anything, 1).Return(nil)
				mockEventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)
			},
			exampleId: 1,
			wantErr:   false,
		},
		{
			name: "Example does not exist",
			setupMocks: func() {
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 2).Return(nil, repo.ErrNotFound)
			},
			exampleId:   2,
			wantErr:     true,
			expectedErr: repo.ErrNotFound,
		},
		{
			name: "Delete error",
			setupMocks: func() {
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 3).Return(&model.Example{
					Id:    3,
					Name:  "Test Example",
					Alias: "test-alias",
				}, nil)
				mockRepo.On("Delete", mock.Anything, mock.Anything, 3).Return(errors.New("delete error"))
			},
			exampleId: 3,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set mock behavior
			mockRepo.ExpectedCalls = nil
			mockCacheRepo.ExpectedCalls = nil
			mockEventBus.ExpectedCalls = nil
			tc.setupMocks()

			// Execute test
			ctx := context.Background()
			err := service.Delete(ctx, tc.exampleId)

			// Verify results
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify Mock calls
			mockRepo.AssertExpectations(t)
			mockCacheRepo.AssertExpectations(t)
			mockEventBus.AssertExpectations(t)
		})
	}
}

// Test Update method
func TestExampleService_Update(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)
	mockEventBus := new(MockEventBus)

	// Create service instance
	service := NewExampleService(mockRepo, mockCacheRepo)
	service.EventBus = mockEventBus

	testCases := []struct {
		name        string
		setupMocks  func()
		exampleId   int
		newName     string
		newAlias    string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Successfully update example",
			setupMocks: func() {
				example := &model.Example{
					Id:    1,
					Name:  "Original Name",
					Alias: "original-alias",
				}
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 1).Return(example, nil)
				mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)
				mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)
				mockEventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)
			},
			exampleId: 1,
			newName:   "Updated Name",
			newAlias:  "updated-alias",
			wantErr:   false,
		},
		{
			name: "Example does not exist",
			setupMocks: func() {
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 2).Return(nil, repo.ErrNotFound)
			},
			exampleId:   2,
			newName:     "Updated Name",
			newAlias:    "updated-alias",
			wantErr:     true,
			expectedErr: repo.ErrNotFound,
		},
		{
			name: "Update validation failed",
			setupMocks: func() {
				example := &model.Example{
					Id:    3,
					Name:  "Original Name",
					Alias: "original-alias",
				}
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 3).Return(example, nil)
				// Update won't be called due to validation failure
			},
			exampleId:   3,
			newName:     "", // Empty name will cause validation error
			newAlias:    "updated-alias",
			wantErr:     true,
			expectedErr: model.ErrEmptyExampleName,
		},
		{
			name: "Update storage error",
			setupMocks: func() {
				example := &model.Example{
					Id:    4,
					Name:  "Original Name",
					Alias: "original-alias",
				}
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 4).Return(example, nil)
				mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Return(errors.New("update error"))
			},
			exampleId: 4,
			newName:   "Updated Name",
			newAlias:  "updated-alias",
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set mock behavior
			mockRepo.ExpectedCalls = nil
			mockCacheRepo.ExpectedCalls = nil
			mockEventBus.ExpectedCalls = nil
			tc.setupMocks()

			// Execute test
			ctx := context.Background()
			err := service.Update(ctx, tc.exampleId, tc.newName, tc.newAlias)

			// Verify results
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify Mock calls
			mockRepo.AssertExpectations(t)
			mockCacheRepo.AssertExpectations(t)
			mockEventBus.AssertExpectations(t)
		})
	}
}

// Test Get method
func TestExampleService_Get(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)

	// Create service instance
	service := NewExampleService(mockRepo, mockCacheRepo)

	testCases := []struct {
		name        string
		setupMocks  func()
		exampleId   int
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Get example from cache",
			setupMocks: func() {
				// Mock cache hit
				mockCacheRepo.On("GetByID", mock.Anything, 1).Return(&model.Example{
					Id:    1,
					Name:  "Cached Example",
					Alias: "cached-alias",
				}, nil)
				// Repository should not be called
			},
			exampleId: 1,
			wantErr:   false,
		},
		{
			name: "Get example from repository (cache miss)",
			setupMocks: func() {
				// Mock cache miss
				mockCacheRepo.On("GetByID", mock.Anything, 2).Return(nil, errors.New("cache miss"))
				// Get from repository
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 2).Return(&model.Example{
					Id:    2,
					Name:  "Database Example",
					Alias: "db-alias",
				}, nil)
				// Update cache
				mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)
			},
			exampleId: 2,
			wantErr:   false,
		},
		{
			name: "Example not found in both cache and repository",
			setupMocks: func() {
				// Mock cache miss
				mockCacheRepo.On("GetByID", mock.Anything, 3).Return(nil, errors.New("cache miss"))
				// Repository also not found
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 3).Return(nil, repo.ErrNotFound)
			},
			exampleId:   3,
			wantErr:     true,
			expectedErr: repo.ErrNotFound,
		},
		{
			name: "Get from repository without cache",
			setupMocks: func() {
				// Create service without cache
				service = NewExampleService(mockRepo, nil)
				// Only use repository
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 4).Return(&model.Example{
					Id:    4,
					Name:  "Database Example",
					Alias: "db-alias",
				}, nil)
			},
			exampleId: 4,
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set mock behavior
			mockRepo.ExpectedCalls = nil
			mockCacheRepo.ExpectedCalls = nil
			tc.setupMocks()

			// Execute test
			ctx := context.Background()
			result, err := service.Get(ctx, tc.exampleId)

			// Verify results
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.exampleId, result.Id)
			}

			// Verify Mock calls
			mockRepo.AssertExpectations(t)
			mockCacheRepo.AssertExpectations(t)
		})
	}
}

// Test FindByName method
func TestExampleService_FindByName(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)

	// Create service instance
	service := NewExampleService(mockRepo, mockCacheRepo)

	testCases := []struct {
		name        string
		setupMocks  func()
		searchName  string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Get example from cache",
			setupMocks: func() {
				// Mock cache hit
				mockCacheRepo.On("GetByName", mock.Anything, "test-name").Return(&model.Example{
					Id:    1,
					Name:  "test-name",
					Alias: "cached-alias",
				}, nil)
				// Repository should not be called
			},
			searchName: "test-name",
			wantErr:    false,
		},
		{
			name: "Get example from repository (cache miss)",
			setupMocks: func() {
				// Mock cache miss
				mockCacheRepo.On("GetByName", mock.Anything, "db-name").Return(nil, errors.New("cache miss"))
				// Get from repository
				mockRepo.On("FindByName", mock.Anything, mock.Anything, "db-name").Return(&model.Example{
					Id:    2,
					Name:  "db-name",
					Alias: "db-alias",
				}, nil)
				// Update cache
				mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)
			},
			searchName: "db-name",
			wantErr:    false,
		},
		{
			name: "Example not found in both cache and repository",
			setupMocks: func() {
				// Mock cache miss
				mockCacheRepo.On("GetByName", mock.Anything, "missing-name").Return(nil, errors.New("cache miss"))
				// Repository also not found
				mockRepo.On("FindByName", mock.Anything, mock.Anything, "missing-name").Return(nil, repo.ErrNotFound)
			},
			searchName:  "missing-name",
			wantErr:     true,
			expectedErr: repo.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set mock behavior
			mockRepo.ExpectedCalls = nil
			mockCacheRepo.ExpectedCalls = nil
			tc.setupMocks()

			// Execute test
			ctx := context.Background()
			result, err := service.FindByName(ctx, tc.searchName)

			// Verify results
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.searchName, result.Name)
			}

			// Verify Mock calls
			mockRepo.AssertExpectations(t)
			mockCacheRepo.AssertExpectations(t)
		})
	}
}
