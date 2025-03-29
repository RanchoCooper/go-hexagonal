package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"go-hexagonal/domain/repo"

	"gorm.io/gorm"
)

// DefaultTimeout defines the default context timeout for database operations
const DefaultTimeout = 30 * time.Second

// Transaction represents a database transaction
type Transaction struct {
	ctx     context.Context
	Session *gorm.DB
	TxOpt   *sql.TxOptions
	store   repo.StoreType
	options *repo.TransactionOptions
}

// Conn returns a database connection with transaction support
func (tr *Transaction) Conn() any {
	if tr == nil {
		return nil
	}
	if tr.Session == nil {
		return nil
	}
	return tr.Session.WithContext(tr.ctx)
}

// Begin starts a new transaction
func (tr *Transaction) Begin() error {
	if tr == nil {
		return ErrInvalidTransaction
	}
	if tr.Session == nil {
		return ErrInvalidSession
	}
	// Begin a new transaction
	tr.Session = tr.Session.Begin(tr.TxOpt)
	if tr.Session.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tr.Session.Error)
	}
	return nil
}

// Commit commits the transaction
func (tr *Transaction) Commit() error {
	if tr != nil && tr.Session != nil {
		return tr.Session.Commit().Error
	}
	return nil
}

// Rollback rolls back the transaction
func (tr *Transaction) Rollback() error {
	if tr != nil && tr.Session != nil {
		return tr.Session.Rollback().Error
	}
	return nil
}

// WithContext returns a new transaction with the given context
func (tr *Transaction) WithContext(ctx context.Context) repo.Transaction {
	if tr == nil {
		return nil
	}
	newTr := *tr
	newTr.ctx = ctx
	if newTr.Session != nil {
		newTr.Session = newTr.Session.WithContext(ctx)
	}
	return &newTr
}

// Context returns the transaction's context
func (tr *Transaction) Context() context.Context {
	return tr.ctx
}

// StoreType returns the store type
func (tr *Transaction) StoreType() repo.StoreType {
	return tr.store
}

// Options returns the transaction options
func (tr *Transaction) Options() *repo.TransactionOptions {
	return tr.options
}

// NewTransaction creates a new transaction with the specified store type and options
func NewTransaction(ctx context.Context, store StoreType, client any, sqlOpt *sql.TxOptions) (*Transaction, error) {
	// Convert store type to domain store type
	domainStore := repo.StoreType(store)

	// Convert SQL options to transaction options
	var options *repo.TransactionOptions
	if sqlOpt != nil {
		options = &repo.TransactionOptions{
			ReadOnly: sqlOpt.ReadOnly,
		}
	} else {
		options = repo.DefaultTransactionOptions()
	}

	tr := &Transaction{
		ctx:     ctx,
		TxOpt:   sqlOpt,
		store:   domainStore,
		options: options,
	}

	switch store {
	case MySQLStore, PostgreSQLStore:
		// Handle SQL-based databases with GORM
		var db *gorm.DB

		// Use reflection to check type and call appropriate method
		clientValue := reflect.ValueOf(client)
		if clientValue.Kind() == reflect.Ptr && !clientValue.IsNil() {
			// Try to call GetDB method
			method := clientValue.MethodByName("GetDB")
			if method.IsValid() {
				result := method.Call([]reflect.Value{reflect.ValueOf(ctx)})
				if len(result) > 0 && !result[0].IsNil() {
					if gormDB, ok := result[0].Interface().(*gorm.DB); ok {
						db = gormDB
					}
				}
			}
		}

		if db == nil {
			return nil, fmt.Errorf("failed to get database session")
		}
		// Initialize the session
		tr.Session = db
	case RedisStore:
		// Redis transactions would be implemented here
		return nil, fmt.Errorf("redis transaction not implemented")
	default:
		return nil, ErrUnsupportedStoreType
	}

	return tr, nil
}

// StoreType defines the type of storage
type StoreType string

// Available store types
const (
	MySQLStore      StoreType = "MySQL"
	RedisStore      StoreType = "Redis"
	PostgreSQLStore StoreType = "PostgreSQL"
)
