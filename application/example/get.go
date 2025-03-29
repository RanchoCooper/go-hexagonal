package example

import (
	"context"
	"fmt"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
	"go-hexagonal/util/log"
)

// GetUseCase handles the get example use case
type GetUseCase struct {
	*core.UseCaseHandler
	exampleService service.IExampleService
}

// NewGetUseCase creates a new GetUseCase instance
func NewGetUseCase(
	exampleService service.IExampleService,
	txFactory repo.TransactionFactory,
) *GetUseCase {
	return &GetUseCase{
		UseCaseHandler: core.NewUseCaseHandler(txFactory),
		exampleService: exampleService,
	}
}

// Execute processes the get example request
func (uc *GetUseCase) Execute(ctx context.Context, input any) (any, error) {
	// Convert and validate input
	getInput, ok := input.(*GetInput)
	if !ok {
		return nil, core.ValidationError("invalid input type", nil)
	}

	if err := getInput.Validate(); err != nil {
		return nil, err
	}

	// Retrieve example directly (no transaction needed)
	example, err := uc.exampleService.Get(ctx, getInput.ID)
	if err != nil {
		log.SugaredLogger.Errorf("Failed to get example: %v", err)
		return nil, fmt.Errorf("failed to get example: %w", err)
	}

	if example == nil {
		return nil, core.NotFoundError(fmt.Sprintf("example with ID %d not found", getInput.ID))
	}

	// Create output DTO
	return NewExampleOutput(example), nil
}
