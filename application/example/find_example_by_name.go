package example

import (
	"context"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/service"
)

// FindExampleByNameInput represents input for finding an example by name
type FindExampleByNameInput struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// FindExampleByNameOutput represents output after finding an example by name
type FindExampleByNameOutput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// FindExampleByNameHandler handles finding example by name
type FindExampleByNameHandler struct {
	ExampleService *service.ExampleService
}

// NewFindExampleByNameHandler creates a new handler instance
func NewFindExampleByNameHandler(exampleService *service.ExampleService) *FindExampleByNameHandler {
	return &FindExampleByNameHandler{
		ExampleService: exampleService,
	}
}

// Handle processes the find example by name request
func (h *FindExampleByNameHandler) Handle(ctx context.Context, input interface{}) (interface{}, error) {
	findInput, ok := input.(FindExampleByNameInput)
	if !ok {
		return nil, core.NewValidationError(400, "invalid input type", core.ErrInvalidInput)
	}

	example, err := h.ExampleService.FindByName(ctx, findInput.Name)
	if err != nil {
		return nil, core.NewInternalError(500, "failed to find example by name", err)
	}
	if example == nil {
		return nil, core.NewNotFoundError(404, "example not found", core.ErrNotFound)
	}

	return FindExampleByNameOutput{
		ID:    example.Id,
		Name:  example.Name,
		Alias: example.Alias,
	}, nil
}

// FindExampleByNameUseCase represents the use case for finding examples by name
type FindExampleByNameUseCase struct {
	handler core.UseCaseHandler
}

// NewFindExampleByNameUseCase creates a new use case instance
func NewFindExampleByNameUseCase(exampleService *service.ExampleService) *FindExampleByNameUseCase {
	return &FindExampleByNameUseCase{
		handler: NewFindExampleByNameHandler(exampleService),
	}
}

// Execute executes the use case
func (u *FindExampleByNameUseCase) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	return u.handler.Handle(ctx, input)
}
