// Package core provides core interfaces and abstractions for the application layer
package core

import (
	"context"

	"go-hexagonal/domain/repo"
	"go-hexagonal/util/errors"
	"go-hexagonal/util/log"
)

// UseCase defines the interface for all use cases in the application
type UseCase interface {
	// Execute processes the use case with the given input and returns the result or an error
	Execute(ctx context.Context, input any) (any, error)
}

// UseCaseHandler provides a base implementation for use cases
type UseCaseHandler struct {
	TxFactory repo.TransactionFactory
}

// NewUseCaseHandler creates a new use case handler
func NewUseCaseHandler(txFactory repo.TransactionFactory) *UseCaseHandler {
	return &UseCaseHandler{
		TxFactory: txFactory,
	}
}

// ExecuteInTransaction executes the given function within a transaction
func (h *UseCaseHandler) ExecuteInTransaction(
	ctx context.Context,
	storeType repo.StoreType,
	fn func(context.Context, repo.Transaction) (any, error),
) (any, error) {
	// Create transaction
	tx, err := h.TxFactory.NewTransaction(ctx, storeType, nil)
	if err != nil {
		log.SugaredLogger.Errorf("Failed to create transaction: %v", err)
		return nil, errors.Wrapf(err, errors.ErrorTypeSystem, "failed to create transaction")
	}
	defer tx.Rollback()

	// Execute function within transaction
	result, err := fn(ctx, tx)
	if err != nil {
		log.SugaredLogger.Errorf("Transaction execution failed: %v", err)
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.SugaredLogger.Errorf("Failed to commit transaction: %v", err)
		return nil, errors.Wrapf(err, errors.ErrorTypeSystem, "failed to commit transaction")
	}

	return result, nil
}

// Input defines the base interface for all use case inputs
type Input interface {
	Validate() error
}

// BaseInput provides common input validation functionality
type BaseInput struct{}

// Validate performs basic validation on the input
func (b *BaseInput) Validate() error {
	return nil
}

// Output defines the base interface for all use case outputs
type Output interface {
	GetStatus() string
}

// BaseOutput provides common output functionality
type BaseOutput struct {
	Status string `json:"status,omitempty"`
}

// GetStatus returns the status of the output
func (o *BaseOutput) GetStatus() string {
	return o.Status
}

// NewSuccessOutput creates a new success output
func NewSuccessOutput() *BaseOutput {
	return &BaseOutput{Status: "success"}
}

// ValidationError returns a validation error with the given message and details
func ValidationError(message string, details map[string]any) error {
	return errors.NewValidationError(message, nil).WithDetails(details)
}

// NotFoundError returns a not found error with the given message
func NotFoundError(message string) error {
	return errors.New(errors.ErrorTypeNotFound, message)
}

// SystemError returns a system error with the given message and cause
func SystemError(message string, cause error) error {
	return errors.NewSystemError(message, cause)
}

// BusinessError returns a business error with the given message
func BusinessError(message string) error {
	return errors.New(errors.ErrorTypeBusiness, message)
}
