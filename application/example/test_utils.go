package example

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"go-hexagonal/domain/repo"
)

// TestTransaction 测试用事务实现
type TestTransaction struct {
	mock.Mock
}

// Begin 开始事务
func (tx *TestTransaction) Begin() error {
	args := tx.Called()
	return args.Error(0)
}

// Commit 提交事务
func (tx *TestTransaction) Commit() error {
	args := tx.Called()
	return args.Error(0)
}

// Rollback 回滚事务
func (tx *TestTransaction) Rollback() error {
	args := tx.Called()
	return args.Error(0)
}

// Conn 获取底层连接
func (tx *TestTransaction) Conn(ctx context.Context) any {
	args := tx.Called(ctx)
	return args.Get(0)
}

// Context 获取事务上下文
func (tx *TestTransaction) Context() context.Context {
	args := tx.Called()
	return args.Get(0).(context.Context)
}

// WithContext 设置事务上下文
func (tx *TestTransaction) WithContext(ctx context.Context) repo.Transaction {
	args := tx.Called(ctx)
	return args.Get(0).(repo.Transaction)
}

// StoreType 获取存储类型
func (tx *TestTransaction) StoreType() repo.StoreType {
	args := tx.Called()
	return args.Get(0).(repo.StoreType)
}

// Options 获取事务选项
func (tx *TestTransaction) Options() *repo.TransactionOptions {
	args := tx.Called()
	return args.Get(0).(*repo.TransactionOptions)
}

// CreateTestTransaction 创建测试用事务
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

// ErrorTestTransaction 创建一个会失败的事务
func ErrorTestTransaction(ctx context.Context) (repo.Transaction, error) {
	return nil, errors.New("创建事务失败")
}

// CommitErrorTestTransaction 创建一个提交会失败的事务
func CommitErrorTestTransaction(ctx context.Context) (repo.Transaction, error) {
	tx := new(TestTransaction)
	tx.On("Begin").Return(nil)
	tx.On("Commit").Return(errors.New("提交事务失败"))
	tx.On("Rollback").Return(nil)
	tx.On("Context").Return(ctx)
	tx.On("StoreType").Return(repo.MySQLStore)
	tx.On("Options").Return(repo.DefaultTransactionOptions())
	tx.On("WithContext", mock.Anything).Return(tx)
	return tx, nil
}
