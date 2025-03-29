package repo

import (
	"context"
	"fmt"
	"sync"

	"go-hexagonal/domain/repo"
	"go-hexagonal/util/metrics"
)

// TransactionFactory defines the interface for creating transactions
type TransactionFactory interface {
	// NewTransaction creates a new transaction with the specified store type and options
	NewTransaction(ctx context.Context, storeType repo.StoreType, opts *repo.TransactionOptions) (repo.Transaction, error)
	// ExecuteInTransaction executes the given function within a transaction
	ExecuteInTransaction(ctx context.Context, storeType repo.StoreType, opts *repo.TransactionOptions, fn func(ctx context.Context, tx repo.Transaction) error) error
}

// BaseTransactionFactory provides a basic implementation of TransactionFactory
type BaseTransactionFactory struct {
	mu         sync.RWMutex
	factories  map[repo.StoreType]func(ctx context.Context, opts *repo.TransactionOptions) (repo.Transaction, error)
	isTestMode bool
}

// NewBaseTransactionFactory creates a new BaseTransactionFactory
func NewBaseTransactionFactory(isTestMode bool) *BaseTransactionFactory {
	return &BaseTransactionFactory{
		factories:  make(map[repo.StoreType]func(ctx context.Context, opts *repo.TransactionOptions) (repo.Transaction, error)),
		isTestMode: isTestMode,
	}
}

// RegisterFactory registers a transaction factory for a specific store type
func (f *BaseTransactionFactory) RegisterFactory(storeType repo.StoreType, factory func(ctx context.Context, opts *repo.TransactionOptions) (repo.Transaction, error)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.factories[storeType] = factory
}

// NewTransaction creates a new transaction with the specified store type and options
func (f *BaseTransactionFactory) NewTransaction(ctx context.Context, storeType repo.StoreType, opts *repo.TransactionOptions) (repo.Transaction, error) {
	if f.isTestMode {
		return repo.NewNoopTransaction(nil), nil
	}

	f.mu.RLock()
	factory, ok := f.factories[storeType]
	f.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no transaction factory registered for store type: %s", storeType)
	}

	return factory(ctx, opts)
}

// ExecuteInTransaction executes the given function within a transaction
func (f *BaseTransactionFactory) ExecuteInTransaction(ctx context.Context, storeType repo.StoreType, opts *repo.TransactionOptions, fn func(ctx context.Context, tx repo.Transaction) error) error {
	if opts == nil {
		opts = repo.DefaultTransactionOptions()
	}

	tx, err := f.NewTransaction(ctx, storeType, opts)
	if err != nil {
		if metrics.Initialized() {
			metrics.RecordError("transaction_factory", string(storeType))
		}
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	err = tx.Begin()
	if err != nil {
		if metrics.Initialized() {
			metrics.RecordError("transaction_begin", string(storeType))
		}
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Measure transaction execution time if metrics are initialized
	var execErr error
	if metrics.Initialized() {
		execErr = metrics.MeasureTransaction(string(storeType), func() error {
			return fn(ctx, tx)
		})
	} else {
		execErr = fn(ctx, tx)
	}

	if execErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			if metrics.Initialized() {
				metrics.RecordError("transaction_rollback", string(storeType))
			}
			return fmt.Errorf("execution failed: %v, and rollback failed: %v", execErr, rollbackErr)
		}
		return execErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		if metrics.Initialized() {
			metrics.RecordError("transaction_commit", string(storeType))
		}
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return nil
}

// MockTransactionFactory is a mock implementation of TransactionFactory for testing
type MockTransactionFactory struct {
	Transactions map[repo.StoreType]repo.Transaction
	Errors       map[string]error
}

// NewMockTransactionFactory creates a new MockTransactionFactory
func NewMockTransactionFactory() *MockTransactionFactory {
	return &MockTransactionFactory{
		Transactions: make(map[repo.StoreType]repo.Transaction),
		Errors:       make(map[string]error),
	}
}

// SetTransactionForStoreType sets a transaction for a specific store type
func (f *MockTransactionFactory) SetTransactionForStoreType(storeType repo.StoreType, tx repo.Transaction) {
	f.Transactions[storeType] = tx
}

// SetError sets an error for a specific operation
func (f *MockTransactionFactory) SetError(operation string, err error) {
	f.Errors[operation] = err
}

// NewTransaction returns a mock transaction for the specified store type
func (f *MockTransactionFactory) NewTransaction(ctx context.Context, storeType repo.StoreType, opts *repo.TransactionOptions) (repo.Transaction, error) {
	if err, ok := f.Errors["NewTransaction"]; ok {
		return nil, err
	}

	tx, ok := f.Transactions[storeType]
	if !ok {
		return repo.NewNoopTransaction(nil), nil
	}

	return tx, nil
}

// ExecuteInTransaction executes the given function within a mock transaction
func (f *MockTransactionFactory) ExecuteInTransaction(ctx context.Context, storeType repo.StoreType, opts *repo.TransactionOptions, fn func(ctx context.Context, tx repo.Transaction) error) error {
	if err, ok := f.Errors["ExecuteInTransaction"]; ok {
		return err
	}

	tx, err := f.NewTransaction(ctx, storeType, opts)
	if err != nil {
		return err
	}

	if err := tx.Begin(); err != nil {
		return err
	}

	execErr := fn(ctx, tx)
	if execErr != nil {
		_ = tx.Rollback()
		return execErr
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
