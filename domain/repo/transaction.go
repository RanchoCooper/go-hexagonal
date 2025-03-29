package repo

import (
	"context"
	"time"

	"go-hexagonal/util/metrics"
)

// StoreType represents the type of data store
type StoreType string

const (
	// MySQLStore represents MySQL data store
	MySQLStore StoreType = "mysql"
	// PostgresStore represents PostgreSQL data store
	PostgresStore StoreType = "postgres"
	// MongoStore represents MongoDB data store
	MongoStore StoreType = "mongo"
	// RedisStore represents Redis data store
	RedisStore StoreType = "redis"
)

// TransactionOptions represents options for a transaction
type TransactionOptions struct {
	// ReadOnly indicates if the transaction is read-only
	ReadOnly bool
	// Timeout specifies the transaction timeout duration
	Timeout time.Duration
	// Isolation specifies the isolation level
	Isolation string
}

// DefaultTransactionOptions returns default transaction options
func DefaultTransactionOptions() *TransactionOptions {
	return &TransactionOptions{
		ReadOnly:  false,
		Timeout:   time.Second * 10,
		Isolation: "",
	}
}

// Transaction defines the interface for database transactions
type Transaction interface {
	// Begin starts the transaction
	Begin() error
	// Commit commits the transaction
	Commit() error
	// Rollback rolls back the transaction
	Rollback() error
	// WithContext returns a new transaction with the given context
	WithContext(ctx context.Context) Transaction
	// Context returns the transaction's context
	Context() context.Context
	// StoreType returns the store type
	StoreType() StoreType
	// Options returns the transaction options
	Options() *TransactionOptions
}

// TransactionHandler defines a higher-level interface for transaction handling with metrics
type TransactionHandler interface {
	// ExecuteInTransaction executes the given function within a transaction
	ExecuteInTransaction(ctx context.Context, opts *TransactionOptions, fn func(ctx context.Context, tx Transaction) error) error
}

// BaseTransaction provides a base implementation for transactions
type BaseTransaction struct {
	ctx       context.Context
	storeType StoreType
	options   *TransactionOptions
}

// NewBaseTransaction creates a new base transaction
func NewBaseTransaction(ctx context.Context, storeType StoreType, options *TransactionOptions) *BaseTransaction {
	if options == nil {
		options = DefaultTransactionOptions()
	}
	return &BaseTransaction{
		ctx:       ctx,
		storeType: storeType,
		options:   options,
	}
}

// Context returns the transaction's context
func (tx *BaseTransaction) Context() context.Context {
	return tx.ctx
}

// WithContext returns a new transaction with the given context
func (tx *BaseTransaction) WithContext(ctx context.Context) Transaction {
	newTx := *tx
	newTx.ctx = ctx
	return &newTx
}

// StoreType returns the store type
func (tx *BaseTransaction) StoreType() StoreType {
	return tx.storeType
}

// Options returns the transaction options
func (tx *BaseTransaction) Options() *TransactionOptions {
	return tx.options
}

// Begin starts the transaction (to be overridden by implementors)
func (tx *BaseTransaction) Begin() error {
	if metrics.Initialized() {
		metrics.RecordTransactionOperation("begin", string(tx.storeType))
	}
	return nil
}

// Commit commits the transaction (to be overridden by implementors)
func (tx *BaseTransaction) Commit() error {
	if metrics.Initialized() {
		metrics.RecordTransactionOperation("commit", string(tx.storeType))
	}
	return nil
}

// Rollback rolls back the transaction (to be overridden by implementors)
func (tx *BaseTransaction) Rollback() error {
	if metrics.Initialized() {
		metrics.RecordTransactionOperation("rollback", string(tx.storeType))
	}
	return nil
}

// ExecuteWithMetrics executes a transaction with metrics
func ExecuteWithMetrics(ctx context.Context, tx Transaction, fn func(ctx context.Context, tx Transaction) error) error {
	if metrics.Initialized() {
		return metrics.MeasureTransaction(string(tx.StoreType()), func() error {
			return fn(ctx, tx)
		})
	}
	return fn(ctx, tx)
}

// NoopTransaction provides a no-operation transaction for cases where a transaction is required by an interface but not needed
type NoopTransaction struct {
	*BaseTransaction
	repo interface{}
}

// NewNoopTransaction creates a new no-operation transaction
func NewNoopTransaction(repo interface{}) *NoopTransaction {
	return &NoopTransaction{
		BaseTransaction: NewBaseTransaction(context.Background(), "noop", nil),
		repo:            repo,
	}
}

// Begin implements the Transaction interface for NoopTransaction
func (tx *NoopTransaction) Begin() error {
	if metrics.Initialized() {
		metrics.RecordTransactionOperation("begin", "noop")
	}
	return nil
}

// Commit implements the Transaction interface for NoopTransaction
func (tx *NoopTransaction) Commit() error {
	if metrics.Initialized() {
		metrics.RecordTransactionOperation("commit", "noop")
	}
	return nil
}

// Rollback implements the Transaction interface for NoopTransaction
func (tx *NoopTransaction) Rollback() error {
	if metrics.Initialized() {
		metrics.RecordTransactionOperation("rollback", "noop")
	}
	return nil
}
