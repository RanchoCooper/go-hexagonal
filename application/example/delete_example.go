package example

import (
	"context"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/service"
)

// DeleteExampleInput represents input for deleting an example
type DeleteExampleInput struct {
	ID int `json:"id" validate:"required"`
}

// DeleteExampleHandler handles example deletion
type DeleteExampleHandler struct {
	ExampleService *service.ExampleService
}

// NewDeleteExampleHandler creates a new handler instance
func NewDeleteExampleHandler(exampleService *service.ExampleService) *DeleteExampleHandler {
	return &DeleteExampleHandler{
		ExampleService: exampleService,
	}
}

// Handle processes the example deletion request
func (h *DeleteExampleHandler) Handle(ctx context.Context, input interface{}) (interface{}, error) {
	deleteInput, ok := input.(DeleteExampleInput)
	if !ok {
		return nil, core.NewValidationError(400, "invalid input type", core.ErrInvalidInput)
	}

	// Delete example
	if err := h.ExampleService.Delete(ctx, deleteInput.ID); err != nil {
		if err == core.ErrNotFound {
			return nil, core.NewNotFoundError(404, "example not found", err)
		}
		return nil, core.NewInternalError(500, "failed to delete example", err)
	}

	return nil, nil
}

// DeleteExampleUseCase represents the use case for deleting examples
type DeleteExampleUseCase struct {
	handler core.UseCaseHandler
}

// NewDeleteExampleUseCase creates a new use case instance
func NewDeleteExampleUseCase(exampleService *service.ExampleService) *DeleteExampleUseCase {
	return &DeleteExampleUseCase{
		handler: NewDeleteExampleHandler(exampleService),
	}
}

// Execute executes the use case
func (u *DeleteExampleUseCase) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	return u.handler.Handle(ctx, input)
}
