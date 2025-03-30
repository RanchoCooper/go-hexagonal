package example

import (
	"context"
	"fmt"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
	"go-hexagonal/util/log"
)

// FindByNameUseCase handles the find example by name use case
type FindByNameUseCase struct {
	*core.UseCaseHandler
	exampleService service.IExampleService
}

// NewFindByNameUseCase creates a new FindByNameUseCase instance
func NewFindByNameUseCase(
	exampleService service.IExampleService,
	txFactory repo.TransactionFactory,
) *FindByNameUseCase {
	return &FindByNameUseCase{
		UseCaseHandler: core.NewUseCaseHandler(txFactory),
		exampleService: exampleService,
	}
}

// Execute processes the find example by name request
func (uc *FindByNameUseCase) Execute(ctx context.Context, input any) (any, error) {
	// Convert and validate input
	findInput, ok := input.(*FindByNameInput)
	if !ok {
		return nil, core.ValidationError("invalid input type", nil)
	}

	if err := findInput.Validate(); err != nil {
		return nil, err
	}

	// Find example by name (no transaction needed for read-only operation)
	example, err := uc.exampleService.FindByName(ctx, findInput.Name)
	if err != nil {
		log.SugaredLogger.Errorf("Failed to find example by name: %v", err)
		return nil, fmt.Errorf("failed to find example by name: %w", err)
	}

	if example == nil {
		return nil, core.NotFoundError(fmt.Sprintf("example with name '%s' not found", findInput.Name))
	}

	// Create output DTO
	return NewExampleOutput(example), nil
}
