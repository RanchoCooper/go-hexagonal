package example

import (
	"context"
	"fmt"

	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// UpdateUseCase handles the update example use case
type UpdateUseCase struct {
	exampleService service.IExampleService
	converter      service.Converter
	txFactory      repo.TransactionFactory
}

// NewUpdateUseCase creates a new UpdateUseCase instance
func NewUpdateUseCase(
	exampleService service.IExampleService,
	converter service.Converter,
	txFactory repo.TransactionFactory,
) *UpdateUseCase {
	return &UpdateUseCase{
		exampleService: exampleService,
		converter:      converter,
		txFactory:      txFactory,
	}
}

// Execute processes the update example request
func (uc *UpdateUseCase) Execute(ctx context.Context, input any) error {
	// Convert input to domain model using the converter
	example, err := uc.converter.FromUpdateRequest(input)
	if err != nil {
		return fmt.Errorf("failed to convert request: %w", err)
	}

	// Create a real transaction for atomic operations
	tx, err := uc.txFactory.NewTransaction(ctx, repo.MySQLStore, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Call domain service
	if err := uc.exampleService.Update(ctx, example); err != nil {
		return fmt.Errorf("failed to update example: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
