package example

import (
	"context"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
	"go-hexagonal/util/error_handler"
)

// CreateUseCase handles the create example use case
type CreateUseCase struct {
	*core.UseCaseHandler
	exampleService service.IExampleService
}

// NewCreateUseCase creates a new CreateUseCase instance
func NewCreateUseCase(
	exampleService service.IExampleService,
	txFactory repo.TransactionFactory,
) *CreateUseCase {
	return &CreateUseCase{
		UseCaseHandler: core.NewUseCaseHandler(txFactory),
		exampleService: exampleService,
	}
}

// Execute processes the create example request
func (uc *CreateUseCase) Execute(ctx context.Context, input any) (any, error) {
	// Convert and validate input
	createInput, ok := input.(*CreateInput)
	if !ok {
		return nil, core.ValidationError("invalid input type", nil)
	}

	if err := createInput.Validate(); err != nil {
		return nil, error_handler.HandleAndConvertError(ctx, err, "input validation", "validation")
	}

	// Execute in transaction
	result, err := uc.ExecuteInTransaction(ctx, repo.MySQLStore, func(ctx context.Context, tx repo.Transaction) (any, error) {
		// Call domain service
		example, err := uc.exampleService.Create(ctx, createInput.Name, createInput.Alias)
		if err != nil {
			return nil, error_handler.HandleAndWrapError(ctx, err, "create example", "failed to create example")
		}

		// Create output DTO
		return NewExampleOutput(example), nil
	})

	if err != nil {
		return nil, error_handler.HandleError(ctx, err, "execute create use case")
	}

	return result, nil
}
