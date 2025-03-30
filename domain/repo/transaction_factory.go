package repo

import (
	"context"
)

// TransactionFactory defines an interface for creating transactions
type TransactionFactory interface {
	// NewTransaction creates a new transaction with the specified store type and options
	NewTransaction(ctx context.Context, store StoreType, opts any) (Transaction, error)
}

// NoopTransactionFactory is a no-operation transaction factory implementation
type NoopTransactionFactory struct{}

// NewNoOpTransactionFactory creates a new NoopTransactionFactory
func NewNoOpTransactionFactory() TransactionFactory {
	return &NoopTransactionFactory{}
}

// NewTransaction creates a new no-operation transaction
func (f *NoopTransactionFactory) NewTransaction(ctx context.Context, store StoreType, opts any) (Transaction, error) {
	return NewNoopTransaction(nil), nil
}
