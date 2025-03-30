package example

import (
	"context"
	"fmt"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
	"go-hexagonal/util/log"
)

// DeleteUseCase handles the delete example use case
type DeleteUseCase struct {
	*core.UseCaseHandler
	exampleService service.IExampleService
}

// NewDeleteUseCase creates a new DeleteUseCase instance
func NewDeleteUseCase(
	exampleService service.IExampleService,
	txFactory repo.TransactionFactory,
) *DeleteUseCase {
	return &DeleteUseCase{
		UseCaseHandler: core.NewUseCaseHandler(txFactory),
		exampleService: exampleService,
	}
}

// Execute processes the delete example request
func (uc *DeleteUseCase) Execute(ctx context.Context, input any) (any, error) {
	// Convert and validate input
	deleteInput, ok := input.(*DeleteInput)
	if !ok {
		return nil, core.ValidationError("invalid input type", nil)
	}

	if err := deleteInput.Validate(); err != nil {
		return nil, err
	}

	// Execute in transaction
	_, err := uc.ExecuteInTransaction(ctx, repo.MySQLStore, func(ctx context.Context, tx repo.Transaction) (any, error) {
		// Call domain service to delete the example
		err := uc.exampleService.Delete(ctx, deleteInput.ID)
		if err != nil {
			log.SugaredLogger.Errorf("Failed to delete example: %v", err)
			return nil, fmt.Errorf("failed to delete example: %w", err)
		}

		return core.NewSuccessOutput(), nil
	})

	if err != nil {
		return nil, err
	}

	return core.NewSuccessOutput(), nil
}
