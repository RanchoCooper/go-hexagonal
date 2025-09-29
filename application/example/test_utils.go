package example

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"go-hexagonal/domain/repo"
)

// TestTransaction test transaction implementation
type TestTransaction struct {
	mock.Mock
}

// Begin starts transaction
func (tx *TestTransaction) Begin() error {
	args := tx.Called()
	return args.Error(0)
}

// Commit commits transaction
func (tx *TestTransaction) Commit() error {
	args := tx.Called()
	return args.Error(0)
}

// Rollback rolls back transaction
func (tx *TestTransaction) Rollback() error {
	args := tx.Called()
	return args.Error(0)
}

// Conn gets underlying connection
func (tx *TestTransaction) Conn(ctx context.Context) any {
	args := tx.Called(ctx)
	return args.Get(0)
}

// Context gets transaction context
func (tx *TestTransaction) Context() context.Context {
	args := tx.Called()
	return args.Get(0).(context.Context)
}

// WithContext sets transaction context
func (tx *TestTransaction) WithContext(ctx context.Context) repo.Transaction {
	args := tx.Called(ctx)
	return args.Get(0).(repo.Transaction)
}

// StoreType gets storage type
func (tx *TestTransaction) StoreType() repo.StoreType {
	args := tx.Called()
	return args.Get(0).(repo.StoreType)
}

// Options gets transaction options
func (tx *TestTransaction) Options() *repo.TransactionOptions {
	args := tx.Called()
	return args.Get(0).(*repo.TransactionOptions)
}

// CreateTestTransaction creates test transaction
func CreateTestTransaction(ctx context.Context) (repo.Transaction, error) {
	tx := new(TestTransaction)
	tx.On("Begin").Return(nil)
	tx.On("Commit").Return(nil)
	tx.On("Rollback").Return(nil)
	tx.On("Context").Return(ctx)
	tx.On("StoreType").Return(repo.MySQLStore)
	tx.On("Options").Return(repo.DefaultTransactionOptions())
	tx.On("WithContext", mock.Anything).Return(tx)
	return tx, nil
}

// ErrorTestTransaction creates a transaction that will fail
func ErrorTestTransaction(ctx context.Context) (repo.Transaction, error) {
	return nil, errors.New("failed to create transaction")
}

// CommitErrorTestTransaction creates a transaction that will fail on commit
func CommitErrorTestTransaction(ctx context.Context) (repo.Transaction, error) {
	tx := new(TestTransaction)
	tx.On("Begin").Return(nil)
	tx.On("Commit").Return(errors.New("failed to commit transaction"))
	tx.On("Rollback").Return(nil)
	tx.On("Context").Return(ctx)
	tx.On("StoreType").Return(repo.MySQLStore)
	tx.On("Options").Return(repo.DefaultTransactionOptions())
	tx.On("WithContext", mock.Anything).Return(tx)
	return tx, nil
}
