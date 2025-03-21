package example

import (
	"context"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/service"
)

// UpdateExampleInput represents input for updating an example
type UpdateExampleInput struct {
	ID    int    `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Alias string `json:"alias" validate:"omitempty,max=255"`
}

// UpdateExampleOutput represents output after updating an example
type UpdateExampleOutput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// UpdateExampleHandler handles example updates
type UpdateExampleHandler struct {
	ExampleService *service.ExampleService
}

// NewUpdateExampleHandler creates a new handler instance
func NewUpdateExampleHandler(exampleService *service.ExampleService) *UpdateExampleHandler {
	return &UpdateExampleHandler{
		ExampleService: exampleService,
	}
}

// Handle processes the example update request
func (h *UpdateExampleHandler) Handle(ctx context.Context, input interface{}) (interface{}, error) {
	updateInput, ok := input.(UpdateExampleInput)
	if !ok {
		return nil, core.NewValidationError(core.StatusBadRequest, "invalid input type", core.ErrInvalidInput)
	}

	// Check if example exists
	_, err := h.ExampleService.Get(ctx, updateInput.ID)
	if err != nil {
		return nil, core.NewInternalError(core.StatusInternalServerError, "failed to check example existence", err)
	}
	if err == core.ErrNotFound {
		return nil, core.NewNotFoundError(core.StatusNotFound, "example not found", core.ErrNotFound)
	}

	// Create domain model from input
	example := &model.Example{
		Id:    updateInput.ID,
		Name:  updateInput.Name,
		Alias: updateInput.Alias,
	}

	// Update example - ExampleService.Update only returns error
	if err := h.ExampleService.Update(ctx, example); err != nil {
		return nil, core.NewInternalError(core.StatusInternalServerError, "failed to update example", err)
	}

	return UpdateExampleOutput{
		ID:    example.Id,
		Name:  example.Name,
		Alias: example.Alias,
	}, nil
}

// UpdateExampleUseCase represents the use case for updating examples
type UpdateExampleUseCase struct {
	handler core.UseCaseHandler
}

// NewUpdateExampleUseCase creates a new use case instance
func NewUpdateExampleUseCase(exampleService *service.ExampleService) *UpdateExampleUseCase {
	return &UpdateExampleUseCase{
		handler: NewUpdateExampleHandler(exampleService),
	}
}

// Execute executes the use case
func (u *UpdateExampleUseCase) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	return u.handler.Handle(ctx, input)
}
