package example

import (
	"context"
	"fmt"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
	"go-hexagonal/util/log"
)

// UpdateUseCase handles the update example use case
type UpdateUseCase struct {
	*core.UseCaseHandler
	exampleService service.IExampleService
}

// NewUpdateUseCase creates a new UpdateUseCase instance
func NewUpdateUseCase(
	exampleService service.IExampleService,
	txFactory repo.TransactionFactory,
) *UpdateUseCase {
	return &UpdateUseCase{
		UseCaseHandler: core.NewUseCaseHandler(txFactory),
		exampleService: exampleService,
	}
}

// Execute processes the update example request
func (uc *UpdateUseCase) Execute(ctx context.Context, input any) (any, error) {
	// Convert and validate input
	updateInput, ok := input.(*UpdateInput)
	if !ok {
		return nil, core.ValidationError("invalid input type", nil)
	}

	if err := updateInput.Validate(); err != nil {
		return nil, err
	}

	// Execute in transaction
	result, err := uc.ExecuteInTransaction(ctx, repo.MySQLStore, func(ctx context.Context, tx repo.Transaction) (any, error) {
		// Call domain service to update the example
		err := uc.exampleService.Update(ctx, updateInput.ID, updateInput.Name, updateInput.Alias)
		if err != nil {
			log.SugaredLogger.Errorf("Failed to update example: %v", err)
			return nil, fmt.Errorf("failed to update example: %w", err)
		}

		// Get the updated example
		updatedExample, err := uc.exampleService.Get(ctx, updateInput.ID)
		if err != nil {
			log.SugaredLogger.Errorf("Failed to get updated example: %v", err)
			return nil, fmt.Errorf("failed to get updated example: %w", err)
		}

		// Create output DTO
		return NewExampleOutput(updatedExample), nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
