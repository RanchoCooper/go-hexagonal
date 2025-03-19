package repo

import "context"

// Transaction defines the transaction interface
type Transaction interface {
	Begin() error
	Commit() error
	Rollback() error
	Conn(ctx context.Context) interface{}
}

// NoopTransaction is a no-operation transaction implementation
type NoopTransaction struct {
	conn interface{}
}

func NewNoopTransaction(conn interface{}) *NoopTransaction {
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

func (t *NoopTransaction) Conn(ctx context.Context) interface{} {
	return t.conn
}
