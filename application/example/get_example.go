package example

import (
	"context"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/service"
)

// GetExampleInput represents input for retrieving an example
type GetExampleInput struct {
	ID int `json:"id"`
}

// GetExampleOutput represents output after retrieving an example
type GetExampleOutput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// GetExampleHandler handles example retrieval
type GetExampleHandler struct {
	ExampleService *service.ExampleService
}

// NewGetExampleHandler creates a new handler instance
func NewGetExampleHandler(exampleService *service.ExampleService) *GetExampleHandler {
	return &GetExampleHandler{
		ExampleService: exampleService,
	}
}

// Handle processes the example retrieval request
func (h *GetExampleHandler) Handle(ctx context.Context, input interface{}) (interface{}, error) {
	getInput, ok := input.(GetExampleInput)
	if !ok {
		return nil, core.ErrInvalidInput
	}

	example, err := h.ExampleService.Get(ctx, getInput.ID)
	if err != nil {
		return nil, err
	}

	if example == nil {
		return nil, core.ErrNotFound
	}

	return GetExampleOutput{
		ID:    example.Id,
		Name:  example.Name,
		Alias: example.Alias,
	}, nil
}

// GetExampleUseCase represents the use case for retrieving examples
type GetExampleUseCase struct {
	handler core.UseCaseHandler
}

// NewGetExampleUseCase creates a new use case instance
func NewGetExampleUseCase(exampleService *service.ExampleService) *GetExampleUseCase {
	return &GetExampleUseCase{
		handler: NewGetExampleHandler(exampleService),
	}
}

// Execute executes the use case
func (u *GetExampleUseCase) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	return u.handler.Handle(ctx, input)
}
