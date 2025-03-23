package repo

import "context"

// StoreType defines the type of storage
type StoreType string

const (
	MySQLStore      StoreType = "MySQL"
	RedisStore      StoreType = "Redis"
	MongoStore      StoreType = "Mongo"
	PostgreSQLStore StoreType = "PostgreSQL"
)

// Transaction defines the transaction interface
type Transaction interface {
	Begin() error
	Commit() error
	Rollback() error
	Conn(ctx context.Context) any
}

// TransactionFactory defines an interface for creating transactions
type TransactionFactory interface {
	// NewTransaction creates a new transaction with the specified store type and options
	NewTransaction(ctx context.Context, store StoreType, opts any) (Transaction, error)
}

// NoopTransaction is a no-operation transaction implementation
type NoopTransaction struct {
	conn any
}

func NewNoopTransaction(conn any) *NoopTransaction {
	return &NoopTransaction{
		conn: conn,
	}
}

func (t *NoopTransaction) Begin() error {
	return nil
}

func (t *NoopTransaction) Commit() error {
	return nil
}

func (t *NoopTransaction) Rollback() error {
	return nil
}

func (t *NoopTransaction) Conn(ctx context.Context) any {
	return t.conn
}
