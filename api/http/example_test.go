package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/application"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// MockExampleRepo mocks the IExampleRepo interface
type MockExampleRepo struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockExampleRepo) Create(ctx context.Context, tr repo.Transaction, example *model.Example) (*model.Example, error) {
	args := m.Called(ctx, tr, example)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockExampleRepo) Delete(ctx context.Context, tr repo.Transaction, id int) error {
	args := m.Called(ctx, tr, id)
	return args.Error(0)
}

// Update mocks the Update method
func (m *MockExampleRepo) Update(ctx context.Context, tr repo.Transaction, example *model.Example) error {
	args := m.Called(ctx, tr, example)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockExampleRepo) GetByID(ctx context.Context, tr repo.Transaction, id int) (*model.Example, error) {
	args := m.Called(ctx, tr, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// FindByName mocks the FindByName method
func (m *MockExampleRepo) FindByName(ctx context.Context, tr repo.Transaction, name string) (*model.Example, error) {
	args := m.Called(ctx, tr, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// MockConverter mocks the Converter interface
type MockConverter struct {
	mock.Mock
}

func (m *MockConverter) ToExampleResponse(example *model.Example) (any, error) {
	args := m.Called(example)
	return args.Get(0), args.Error(1)
}

func (m *MockConverter) FromCreateRequest(req any) (*model.Example, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *MockConverter) FromUpdateRequest(req any) (*model.Example, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// MockTransaction mocks the Transaction interface
type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Begin() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Conn(ctx context.Context) any {
	args := m.Called(ctx)
	return args.Get(0)
}

// MockTransactionFactory mocks the TransactionFactory interface
type MockTransactionFactory struct {
	mock.Mock
}

func (m *MockTransactionFactory) NewTransaction(ctx context.Context, store repo.StoreType, opts any) (repo.Transaction, error) {
	args := m.Called(ctx, store, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(repo.Transaction), args.Error(1)
}

// setupTest sets up the test environment
func setupTest(t *testing.T) (*gin.Engine, *MockExampleRepo, *service.ExampleService, *MockConverter, func()) {
	// Save original services and factory
	originalServices := services
	originalAppFactory := appFactory

	// Create mock repository
	mockRepo := new(MockExampleRepo)

	// Create mock transaction
	mockTx := new(MockTransaction)
	mockTx.On("Begin").Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Create mock transaction factory
	mockTxFactory := new(MockTransactionFactory)
	mockTxFactory.On("NewTransaction", mock.Anything, mock.Anything, mock.Anything).Return(mockTx, nil)

	// Create mock converter
	mockConverter := new(MockConverter)

	// Create test service
	testService := service.NewExampleService(mockRepo)

	// Create test services
	testServices := &service.Services{
		ExampleService: testService,
	}

	// Register test services
	RegisterServices(testServices)

	// Initialize application factory
	testAppFactory := application.NewFactory(testService, mockConverter, mockTxFactory)
	SetAppFactory(testAppFactory)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Return cleanup function
	cleanup := func() {
		services = originalServices
		appFactory = originalAppFactory
	}

	return router, mockRepo, testService, mockConverter, cleanup
}

func TestCreateExample(t *testing.T) {
	router, mockRepo, _, mockConverter, cleanup := setupTest(t)
	defer cleanup()

	router.POST("/api/examples", CreateExample)

	// Prepare test data
	expectedExample := &model.Example{
		Id:    1,
		Name:  "Test Example",
		Alias: "test",
	}

	// Set up mock behavior
	mockConverter.On("FromCreateRequest", mock.AnythingOfType("dto.CreateExampleReq")).Return(expectedExample, nil)
	mockConverter.On("ToExampleResponse", expectedExample).Return(expectedExample, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Return(expectedExample, nil)

	// Create request
	requestBody := map[string]any{
		"name":  "Test Example",
		"alias": "test",
	}
	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/examples", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusCreated, recorder.Code)
	mockRepo.AssertExpectations(t)
	mockConverter.AssertExpectations(t)
}

func TestGetExample(t *testing.T) {
	router, mockRepo, _, mockConverter, cleanup := setupTest(t)
	defer cleanup()

	router.GET("/api/examples/:id", GetExample)

	// Prepare test data
	expectedExample := &model.Example{
		Id:    1,
		Name:  "Test Example",
		Alias: "test",
	}

	// Set up mock behavior
	mockRepo.On("GetByID", mock.Anything, mock.Anything, 1).Return(expectedExample, nil)
	mockConverter.On("ToExampleResponse", expectedExample).Return(expectedExample, nil)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/api/examples/1", nil)

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
	mockConverter.AssertExpectations(t)
}

func TestUpdateExample(t *testing.T) {
	router, mockRepo, _, mockConverter, cleanup := setupTest(t)
	defer cleanup()

	router.PUT("/api/examples/:id", UpdateExample)

	// Prepare test data
	exampleModel := &model.Example{
		Id:    1,
		Name:  "Updated Example",
		Alias: "updated",
	}

	// Set up mock behavior
	mockConverter.On("FromUpdateRequest", mock.AnythingOfType("dto.UpdateExampleReq")).Return(exampleModel, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)

	// Create request
	requestBody := map[string]any{
		"name":  "Updated Example",
		"alias": "updated",
	}
	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPut, "/api/examples/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
	mockConverter.AssertExpectations(t)
}

func TestDeleteExample(t *testing.T) {
	router, mockRepo, _, _, cleanup := setupTest(t)
	defer cleanup()

	router.DELETE("/api/examples/:id", DeleteExample)

	// Set up mock behavior
	mockRepo.On("Delete", mock.Anything, mock.Anything, 1).Return(nil)

	// Create request
	req, _ := http.NewRequest(http.MethodDelete, "/api/examples/1", nil)

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestFindExampleByName(t *testing.T) {
	router, mockRepo, _, mockConverter, cleanup := setupTest(t)
	defer cleanup()

	router.GET("/api/examples/name/:name", FindExampleByName)

	// Prepare test data
	expectedExample := &model.Example{
		Id:    1,
		Name:  "Test Example",
		Alias: "test",
	}

	// Set up mock behavior - success case
	mockRepo.On("FindByName", mock.Anything, mock.Anything, "test").Return(expectedExample, nil)
	mockConverter.On("ToExampleResponse", expectedExample).Return(expectedExample, nil)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/api/examples/name/test", nil)

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Test not found case
	mockRepo.On("FindByName", mock.Anything, mock.Anything, "nonexistent").
		Return(nil, fmt.Errorf("record not found"))

	// Create request
	req, _ = http.NewRequest(http.MethodGet, "/api/examples/name/nonexistent", nil)

	// Execute request
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	mockRepo.AssertExpectations(t)
	mockConverter.AssertExpectations(t)
}
