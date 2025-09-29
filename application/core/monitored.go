package core

import (
	"context"
	"fmt"

	"go-hexagonal/domain/repo"
	"go-hexagonal/util/errors"
	"go-hexagonal/util/log"
	"go-hexagonal/util/metrics"
)

// MonitoredUseCaseHandler extends UseCaseHandler with monitoring capabilities
type MonitoredUseCaseHandler struct {
	*UseCaseHandler
	useCaseName string
	storeType   repo.StoreType
}

// NewMonitoredUseCaseHandler creates a new monitored use case handler
func NewMonitoredUseCaseHandler(base *UseCaseHandler, storeType repo.StoreType, useCaseName string) *MonitoredUseCaseHandler {
	return &MonitoredUseCaseHandler{
		UseCaseHandler: base,
		storeType:      storeType,
		useCaseName:    useCaseName,
	}
}

// ExecuteInTransaction executes the given function within a transaction with monitoring
func (h *MonitoredUseCaseHandler) ExecuteInTransaction(
	ctx context.Context,
	storeType repo.StoreType,
	fn func(context.Context, repo.Transaction) (any, error),
) (any, error) {
	// If metrics is not initialized, fallback to base implementation
	if !metrics.Initialized() {
		return h.UseCaseHandler.ExecuteInTransaction(ctx, storeType, fn)
	}

	// Measure transaction execution time
	var result any
	var txErr error

	err := metrics.MeasureTransaction(fmt.Sprintf("usecase_%s", h.useCaseName), func() error {
		// Create transaction
		tx, err := h.TxFactory.NewTransaction(ctx, storeType, nil)
		if err != nil {
			metrics.RecordError("transaction_factory", string(storeType))
			log.SugaredLogger.Errorf("Failed to create transaction: %v", err)
			return errors.Wrapf(err, errors.ErrorTypeSystem, "failed to create transaction")
		}
		defer func() { _ = tx.Rollback() }()

		// Execute function within transaction
		result, txErr = fn(ctx, tx)
		if txErr != nil {
			metrics.RecordError("transaction_execution", string(storeType))
			log.SugaredLogger.Errorf("Transaction execution failed: %v", txErr)
			return txErr
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			metrics.RecordError("transaction_commit", string(storeType))
			log.SugaredLogger.Errorf("Failed to commit transaction: %v", err)
			return errors.Wrapf(err, errors.ErrorTypeSystem, "failed to commit transaction")
		}

		return nil
	})

	if err != nil {
		metrics.RecordError("usecase", h.useCaseName)
		return nil, err
	}

	return result, nil
}

// ValidateInput validates the input using the provided validator with monitoring
func (h *MonitoredUseCaseHandler) ValidateInput(input Input) error {
	if !metrics.Initialized() {
		return input.Validate()
	}

	err := input.Validate()
	if err != nil {
		metrics.RecordError("validation", h.useCaseName)
	}
	return err
}

// HandleResultWithMetrics wraps a result with metrics monitoring
func HandleResultWithMetrics(useCaseName string, result any, err error) (any, error) {
	if !metrics.Initialized() {
		return result, err
	}

	if err != nil {
		metrics.RecordError("usecase_result", useCaseName)
	}
	return result, err
}

// MonitoredUseCase wraps a use case with monitoring
type MonitoredUseCase struct {
	useCase     UseCase
	useCaseName string
}

// NewMonitoredUseCase creates a new monitored use case
func NewMonitoredUseCase(useCase UseCase, useCaseName string) *MonitoredUseCase {
	return &MonitoredUseCase{
		useCase:     useCase,
		useCaseName: useCaseName,
	}
}

// Execute executes the use case with monitoring
func (u *MonitoredUseCase) Execute(ctx context.Context, input any) (any, error) {
	if !metrics.Initialized() {
		return u.useCase.Execute(ctx, input)
	}

	var result any
	var execErr error

	err := metrics.MeasureTransaction(fmt.Sprintf("usecase_%s", u.useCaseName), func() error {
		result, execErr = u.useCase.Execute(ctx, input)
		return execErr
	})

	if err != nil {
		metrics.RecordError("usecase", u.useCaseName)
	}

	return result, execErr
}
