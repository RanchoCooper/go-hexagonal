package repo

import (
	"context"
)

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
