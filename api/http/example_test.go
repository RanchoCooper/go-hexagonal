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

	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// Create mock service interface
type MockExampleRepo struct {
	mock.Mock
}

func (m *MockExampleRepo) Create(ctx context.Context, tr repo.Transaction, example *model.Example) (*model.Example, error) {
	args := m.Called(ctx, tr, example)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *MockExampleRepo) Delete(ctx context.Context, tr repo.Transaction, id int) error {
	args := m.Called(ctx, tr, id)
	return args.Error(0)
}

func (m *MockExampleRepo) Update(ctx context.Context, tr repo.Transaction, example *model.Example) error {
	args := m.Called(ctx, tr, example)
	return args.Error(0)
}

func (m *MockExampleRepo) GetByID(ctx context.Context, tr repo.Transaction, id int) (*model.Example, error) {
	args := m.Called(ctx, tr, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *MockExampleRepo) FindByName(ctx context.Context, tr repo.Transaction, name string) (*model.Example, error) {
	args := m.Called(ctx, tr, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// Setup test
func setupTest(t *testing.T) (*gin.Engine, *MockExampleRepo, *service.ExampleService, func()) {
	// Save original service
	originalService := service.ExampleSvc

	// Create mock and new service
	mockRepo := new(MockExampleRepo)
	testService := service.NewExampleService(mockRepo)

	// Replace global service
	service.ExampleSvc = testService

	// Set up Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Return cleanup function
	cleanup := func() {
		service.ExampleSvc = originalService
	}

	return router, mockRepo, testService, cleanup
}

func TestCreateExample(t *testing.T) {
	router, mockRepo, _, cleanup := setupTest(t)
	defer cleanup()

	router.POST("/api/examples", CreateExample)

	// Prepare test data
	expectedExample := &model.Example{
		Id:    1,
		Name:  "Test Example",
		Alias: "test",
	}

	// Set up mock behavior
	mockRepo.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Return(expectedExample, nil)

	// Create request
	requestBody := map[string]interface{}{
		"name":  "Test Example",
		"alias": "test",
	}
	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/examples", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusCreated, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetExample(t *testing.T) {
	router, mockRepo, _, cleanup := setupTest(t)
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

	// Create request
	req, _ := http.NewRequest("GET", "/api/examples/1", nil)

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdateExample(t *testing.T) {
	router, mockRepo, _, cleanup := setupTest(t)
	defer cleanup()

	router.PUT("/api/examples/:id", UpdateExample)

	// Set up mock behavior
	mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)

	// Create request
	requestBody := map[string]interface{}{
		"name":  "Updated Example",
		"alias": "updated",
	}
	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("PUT", "/api/examples/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeleteExample(t *testing.T) {
	router, mockRepo, _, cleanup := setupTest(t)
	defer cleanup()

	router.DELETE("/api/examples/:id", DeleteExample)

	// Set up mock behavior
	mockRepo.On("Delete", mock.Anything, mock.Anything, 1).Return(nil)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/examples/1", nil)

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusOK, recorder.Code)
	mockRepo.AssertExpectations(t)
}

func TestFindExampleByName(t *testing.T) {
	router, mockRepo, _, cleanup := setupTest(t)
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

	// Create request
	req, _ := http.NewRequest("GET", "/api/examples/name/test", nil)

	// Execute request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Test not found case
	mockRepo.On("FindByName", mock.Anything, mock.Anything, "nonexistent").
		Return(nil, fmt.Errorf("record not found"))

	// Create request
	req, _ = http.NewRequest("GET", "/api/examples/name/nonexistent", nil)

	// Execute request
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check results
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	mockRepo.AssertExpectations(t)
}
