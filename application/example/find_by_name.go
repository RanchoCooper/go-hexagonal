package example

import (
	"context"
	"fmt"
	"strings"

	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// FindByNameUseCase handles the find example by name use case
type FindByNameUseCase struct {
	exampleService service.IExampleService
	converter      service.Converter
	txFactory      repo.TransactionFactory
}

// NewFindByNameUseCase creates a new FindByNameUseCase instance
func NewFindByNameUseCase(
	exampleService service.IExampleService,
	converter service.Converter,
	txFactory repo.TransactionFactory,
) *FindByNameUseCase {
	return &FindByNameUseCase{
		exampleService: exampleService,
		converter:      converter,
		txFactory:      txFactory,
	}
}

// Execute processes the find example by name request
func (uc *FindByNameUseCase) Execute(ctx context.Context, name string) (any, error) {
	// Create a transaction for consistent read
	tx, err := uc.txFactory.NewTransaction(ctx, repo.MySQLStore, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// Call domain service
	example, err := uc.exampleService.FindByName(ctx, name)
	if err != nil {
		// Check for "record not found" error
		if strings.Contains(err.Error(), "record not found") {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("failed to find example by name: %w", err)
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
