package example

import (
	"context"
	"fmt"

	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// CreateUseCase handles the create example use case
type CreateUseCase struct {
	exampleService service.IExampleService
	converter      service.Converter
	txFactory      repo.TransactionFactory
}

// NewCreateUseCase creates a new CreateUseCase instance
func NewCreateUseCase(
	exampleService service.IExampleService,
	converter service.Converter,
	txFactory repo.TransactionFactory,
) *CreateUseCase {
	return &CreateUseCase{
		exampleService: exampleService,
		converter:      converter,
		txFactory:      txFactory,
	}
}

// Execute processes the create example request
func (uc *CreateUseCase) Execute(ctx context.Context, input any) (any, error) {
	// Convert input to domain model using the converter
	example, err := uc.converter.FromCreateRequest(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	// Create a real transaction for atomic operations
	tx, err := uc.txFactory.NewTransaction(ctx, repo.MySQLStore, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Call domain service
	createdExample, err := uc.exampleService.Create(ctx, example)
	if err != nil {
		return nil, fmt.Errorf("failed to create example: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert domain model to response using the converter
	result, err := uc.converter.ToExampleResponse(createdExample)
	if err != nil {
		return nil, fmt.Errorf("failed to convert response: %w", err)
	}

	return result, nil
}
