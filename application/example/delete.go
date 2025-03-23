package example

import (
	"context"
	"fmt"

	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// DeleteUseCase handles the delete example use case
type DeleteUseCase struct {
	exampleService service.IExampleService
	txFactory      repo.TransactionFactory
}

// NewDeleteUseCase creates a new DeleteUseCase instance
func NewDeleteUseCase(
	exampleService service.IExampleService,
	txFactory repo.TransactionFactory,
) *DeleteUseCase {
	return &DeleteUseCase{
		exampleService: exampleService,
		txFactory:      txFactory,
	}
}

// Execute processes the delete example request
func (uc *DeleteUseCase) Execute(ctx context.Context, id int) error {
	// Create a real transaction for atomic operations
	tx, err := uc.txFactory.NewTransaction(ctx, repo.MySQLStore, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Call domain service
	if err := uc.exampleService.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
