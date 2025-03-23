package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/domain/event"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

type mockExampleRepo struct {
	mock.Mock
}

func (m *mockExampleRepo) Create(ctx context.Context, tr repo.Transaction, example *model.Example) (*model.Example, error) {
	args := m.Called(ctx, tr, example)
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *mockExampleRepo) Delete(ctx context.Context, tr repo.Transaction, id int) error {
	args := m.Called(ctx, tr, id)
	return args.Error(0)
}

func (m *mockExampleRepo) Update(ctx context.Context, tr repo.Transaction, example *model.Example) error {
	args := m.Called(ctx, tr, example)
	return args.Error(0)
}

func (m *mockExampleRepo) GetByID(ctx context.Context, tr repo.Transaction, id int) (*model.Example, error) {
	args := m.Called(ctx, tr, id)
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *mockExampleRepo) FindByName(ctx context.Context, tr repo.Transaction, name string) (*model.Example, error) {
	args := m.Called(ctx, tr, name)
	return args.Get(0).(*model.Example), args.Error(1)
}

type mockTransaction struct {
	mock.Mock
}

func (m *mockTransaction) Begin() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockTransaction) Conn(ctx context.Context) any {
	args := m.Called(ctx)
	return args.Get(0)
}

type mockEventBus struct {
	mock.Mock
}

func (m *mockEventBus) Publish(ctx context.Context, event event.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockEventBus) Subscribe(handler event.EventHandler) {
	m.Called(handler)
}

func (m *mockEventBus) Unsubscribe(handler event.EventHandler) {
	m.Called(handler)
}

func TestExampleService_Create(t *testing.T) {
	mockRepo := new(mockExampleRepo)
	mockBus := new(mockEventBus)
	service := NewExampleService(mockRepo)
	service.EventBus = mockBus

	ctx := context.Background()
	example := &model.Example{
		Name:  "test",
		Alias: "test_alias",
	}

	expectedExample := &model.Example{
		Id:        1,
		Name:      "test",
		Alias:     "test_alias",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*repo.NoopTransaction"), example).Return(expectedExample, nil)
	mockBus.On("Publish", ctx, mock.AnythingOfType("event.ExampleCreatedEvent")).Return(nil)

	result, err := service.Create(ctx, example)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedExample.Id, result.Id)
	mockRepo.AssertExpectations(t)
	mockBus.AssertExpectations(t)
}

func TestExampleService_Get(t *testing.T) {
	mockRepo := new(mockExampleRepo)
	mockBus := new(mockEventBus)
	service := NewExampleService(mockRepo)
	service.EventBus = mockBus

	ctx := context.Background()
	expectedExample := &model.Example{
		Id:        1,
		Name:      "test",
		Alias:     "test_alias",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", ctx, mock.AnythingOfType("*repo.NoopTransaction"), 1).Return(expectedExample, nil)

	result, err := service.Get(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedExample.Id, result.Id)
	mockRepo.AssertExpectations(t)
}

func TestExampleService_Update(t *testing.T) {
	mockRepo := new(mockExampleRepo)
	mockBus := new(mockEventBus)
	service := NewExampleService(mockRepo)
	service.EventBus = mockBus

	ctx := context.Background()
	example := &model.Example{
		Id:        1,
		Name:      "test",
		Alias:     "test_alias",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("Update", ctx, mock.AnythingOfType("*repo.NoopTransaction"), example).Return(nil)
	mockBus.On("Publish", ctx, mock.AnythingOfType("event.ExampleUpdatedEvent")).Return(nil)

	err := service.Update(ctx, example)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockBus.AssertExpectations(t)
}

func TestExampleService_Delete(t *testing.T) {
	mockRepo := new(mockExampleRepo)
	mockBus := new(mockEventBus)
	service := NewExampleService(mockRepo)
	service.EventBus = mockBus

	ctx := context.Background()
	mockRepo.On("Delete", ctx, mock.AnythingOfType("*repo.NoopTransaction"), 1).Return(nil)
	mockBus.On("Publish", ctx, mock.AnythingOfType("event.ExampleDeletedEvent")).Return(nil)

	err := service.Delete(ctx, 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockBus.AssertExpectations(t)
}

func TestExampleService_FindByName(t *testing.T) {
	mockRepo := new(mockExampleRepo)
	service := NewExampleService(mockRepo)

	ctx := context.Background()
	expectedExample := &model.Example{
		Id:        1,
		Name:      "test",
		Alias:     "test_alias",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByName", ctx, mock.AnythingOfType("*repo.NoopTransaction"), "test").Return(expectedExample, nil)

	result, err := service.FindByName(ctx, "test")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedExample.Id, result.Id)
	mockRepo.AssertExpectations(t)
}
