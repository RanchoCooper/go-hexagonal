package transaction

import (
	"context"

	"go-hexagonal/domain/repo"
	"go-hexagonal/util/errors"
	"go-hexagonal/util/log"
)

// ExecuteWithRollback executes a function within a transaction with automatic rollback handling
func ExecuteWithRollback(
	ctx context.Context,
	txFactory repo.TransactionFactory,
	storeType repo.StoreType,
	fn func(context.Context, repo.Transaction) (any, error),
) (any, error) {
	// Create transaction
	tx, err := txFactory.NewTransaction(ctx, storeType, nil)
	if err != nil {
		if log.SugaredLogger != nil {
			log.SugaredLogger.Errorf("Failed to create transaction: %v", err)
		}
		return nil, errors.Wrapf(err, errors.ErrorTypeSystem, "failed to create transaction")
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	// Execute function within transaction
	if fn == nil {
		return nil, errors.New(errors.ErrorTypeSystem, "function parameter cannot be nil")
	}
	result, err := fn(ctx, tx)
	if err != nil {
		if log.SugaredLogger != nil {
			log.SugaredLogger.Errorf("Transaction execution failed: %v", err)
		}
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		if log.SugaredLogger != nil {
			log.SugaredLogger.Errorf("Failed to commit transaction: %v", err)
		}
		return nil, errors.Wrapf(err, errors.ErrorTypeSystem, "failed to commit transaction")
	}

	return result, nil
}

// ExecuteWithMetrics executes a function within a transaction with metrics monitoring
func ExecuteWithMetrics(
	ctx context.Context,
	txFactory repo.TransactionFactory,
	storeType repo.StoreType,
	useCaseName string,
	fn func(context.Context, repo.Transaction) (any, error),
) (any, error) {
	// Create transaction
	tx, err := txFactory.NewTransaction(ctx, storeType, nil)
	if err != nil {
		if log.SugaredLogger != nil {
			log.SugaredLogger.Errorf("Failed to create transaction: %v", err)
		}
		return nil, errors.Wrapf(err, errors.ErrorTypeSystem, "failed to create transaction")
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	// Execute function within transaction
	if fn == nil {
		return nil, errors.New(errors.ErrorTypeSystem, "function parameter cannot be nil")
	}
	result, err := fn(ctx, tx)
	if err != nil {
		if log.SugaredLogger != nil {
			log.SugaredLogger.Errorf("Transaction execution failed: %v", err)
		}
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		if log.SugaredLogger != nil {
			log.SugaredLogger.Errorf("Failed to commit transaction: %v", err)
		}
		return nil, errors.Wrapf(err, errors.ErrorTypeSystem, "failed to commit transaction")
	}

	return result, nil
}

// WithBackgroundContext creates a proper background context for long-running operations
func WithBackgroundContext() context.Context {
	return context.Background()
}

// WithTestContext creates a context suitable for testing
func WithTestContext() context.Context {
	return context.Background()
}
