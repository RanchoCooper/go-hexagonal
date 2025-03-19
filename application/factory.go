// Package application provides application layer functionality
package application

import (
	"go-hexagonal/application/example"
	"go-hexagonal/domain/event"
	"go-hexagonal/domain/service"
)

// UseCaseFactory creates and provides use cases
type UseCaseFactory struct {
	exampleService *service.ExampleService
	eventBus       event.EventBus
}

// NewUseCaseFactory creates a new use case factory
func NewUseCaseFactory(exampleService *service.ExampleService, eventBus event.EventBus) *UseCaseFactory {
	return &UseCaseFactory{
		exampleService: exampleService,
		eventBus:       eventBus,
	}
}

// CreateExampleUseCase creates an example creation use case
func (f *UseCaseFactory) CreateExampleUseCase() *example.CreateExampleUseCase {
	return example.NewCreateExampleUseCase(f.exampleService)
}

// GetExampleUseCase creates an example retrieval use case
func (f *UseCaseFactory) GetExampleUseCase() *example.GetExampleUseCase {
	return example.NewGetExampleUseCase(f.exampleService)
}

// UpdateExampleUseCase creates an example update use case
func (f *UseCaseFactory) UpdateExampleUseCase() *example.UpdateExampleUseCase {
	return example.NewUpdateExampleUseCase(f.exampleService)
}

// DeleteExampleUseCase creates an example deletion use case
func (f *UseCaseFactory) DeleteExampleUseCase() *example.DeleteExampleUseCase {
	return example.NewDeleteExampleUseCase(f.exampleService)
}

// FindExampleByNameUseCase creates a use case for finding example by name
func (f *UseCaseFactory) FindExampleByNameUseCase() *example.FindExampleByNameUseCase {
	return example.NewFindExampleByNameUseCase(f.exampleService)
}
