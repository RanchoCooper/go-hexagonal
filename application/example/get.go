package example

import (
	"context"
	"fmt"

	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// GetUseCase handles the get example by ID use case
type GetUseCase struct {
	exampleService service.IExampleService
	converter      service.Converter
	txFactory      repo.TransactionFactory
}

// NewGetUseCase creates a new GetUseCase instance
func NewGetUseCase(
	exampleService service.IExampleService,
	converter service.Converter,
	txFactory repo.TransactionFactory,
) *GetUseCase {
	return &GetUseCase{
		exampleService: exampleService,
		converter:      converter,
		txFactory:      txFactory,
	}
}

// Execute processes the get example by ID request
func (uc *GetUseCase) Execute(ctx context.Context, id int) (any, error) {
	// Create a transaction for consistent read
	tx, err := uc.txFactory.NewTransaction(ctx, repo.MySQLStore, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Call domain service
	example, err := uc.exampleService.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get example: %w", err)
	}

	// Commit read-only transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert domain model to response using the converter
	result, err := uc.converter.ToExampleResponse(example)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response: %w", err)
	}

	return result, nil
}
