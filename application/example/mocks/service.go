// Package mocks contains mock implementations for testing
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"go-hexagonal/domain/model"
)

// ExampleService is a mock implementation of service.ExampleService
type ExampleService struct {
	mock.Mock
}

// Create mocks the Create method
func (m *ExampleService) Create(ctx context.Context, example *model.Example) (*model.Example, error) {
	args := m.Called(ctx, example)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// Delete mocks the Delete method
func (m *ExampleService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Update mocks the Update method
func (m *ExampleService) Update(ctx context.Context, example *model.Example) error {
	args := m.Called(ctx, example)
	return args.Error(0)
}

// Get mocks the Get method
func (m *ExampleService) Get(ctx context.Context, id int) (*model.Example, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// FindByName mocks the FindByName method
func (m *ExampleService) FindByName(ctx context.Context, name string) (*model.Example, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}
