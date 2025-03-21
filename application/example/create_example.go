// Package example provides example-related use cases
package example

import (
	"context"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/service"
	"go-hexagonal/util/errors"
)

// CreateExampleInput represents input for creating an example
type CreateExampleInput struct {
	Name  string `json:"name"  validate:"required"`
	Alias string `json:"alias"`
}

// CreateExampleOutput represents output after creating an example
type CreateExampleOutput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// CreateExampleHandler handles example creation
type CreateExampleHandler struct {
	ExampleService *service.ExampleService
}

// NewCreateExampleHandler creates a new handler instance
func NewCreateExampleHandler(exampleService *service.ExampleService) *CreateExampleHandler {
	return &CreateExampleHandler{
		ExampleService: exampleService,
	}
}

// Handle processes the example creation request
func (h *CreateExampleHandler) Handle(ctx context.Context, input interface{}) (interface{}, error) {
	createInput, ok := input.(CreateExampleInput)
	if !ok {
		return nil, errors.NewValidationError("Invalid input data", core.ErrInvalidInput)
	}

	// Validate input
	if createInput.Name == "" {
		return nil, errors.NewValidationError("Name cannot be empty", nil)
	}

	example := &model.Example{
		Name:  createInput.Name,
		Alias: createInput.Alias,
	}

	createdExample, err := h.ExampleService.Create(ctx, example)
	if err != nil {
		// Wrap domain error as application error
		return nil, errors.NewPersistenceError("Failed to create example", err)
	}

	return CreateExampleOutput{
		ID:    createdExample.Id,
		Name:  createdExample.Name,
		Alias: createdExample.Alias,
	}, nil
}

// CreateExampleUseCase represents the use case for creating examples
type CreateExampleUseCase struct {
	handler core.UseCaseHandler
}

// NewCreateExampleUseCase creates a new use case instance
func NewCreateExampleUseCase(exampleService *service.ExampleService) *CreateExampleUseCase {
	return &CreateExampleUseCase{
		handler: NewCreateExampleHandler(exampleService),
	}
}

// Execute executes the use case
func (u *CreateExampleUseCase) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	return u.handler.Handle(ctx, input)
}
