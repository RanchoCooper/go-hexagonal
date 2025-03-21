package example

import (
	"context"
	"errors"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/service"
)

// DeleteExampleInput represents input for deleting an example
type DeleteExampleInput struct {
	ID int `json:"id" validate:"required"`
}

// DeleteExampleOutput represents output for deleting an example
type DeleteExampleOutput struct {
	Success bool `json:"success"`
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
		return nil, core.NewValidationError(core.StatusBadRequest, "invalid input type", core.ErrInvalidInput)
	}

	err := h.ExampleService.Delete(ctx, deleteInput.ID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil, core.NewNotFoundError(core.StatusNotFound, "example not found", err)
		}
		return nil, core.NewInternalError(core.StatusInternalServerError, "failed to delete example", err)
	}

	return DeleteExampleOutput{Success: true}, nil
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
