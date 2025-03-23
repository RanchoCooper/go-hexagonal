package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// DefaultTimeout defines the default context timeout for database operations
const DefaultTimeout = 30 * time.Second

// Transaction represents a database transaction
type Transaction struct {
	Session *gorm.DB
	TxOpt   *sql.TxOptions
}

// Conn returns a database connection with transaction support
func (tr *Transaction) Conn(ctx context.Context) any {
	if tr == nil {
		return nil
	}
	if tr.Session == nil {
		return nil
	}
	return tr.Session.WithContext(ctx)
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

// NewTransaction creates a new transaction with the specified store type and options
func NewTransaction(ctx context.Context, store StoreType, client any, opt *sql.TxOptions) (*Transaction, error) {
	tr := &Transaction{TxOpt: opt}

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
