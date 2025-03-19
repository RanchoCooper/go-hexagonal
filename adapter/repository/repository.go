package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

// DefaultTimeout defines the default context timeout for database operations
const DefaultTimeout = 30 * time.Second

var Clients = &clients{}

// Transaction represents a database transaction
type Transaction struct {
	Session *gorm.DB
	TxOpt   *sql.TxOptions
}

type StoreType string

const (
	MySQLStore      StoreType = "MySQL"
	RedisStore      StoreType = "Redis"
	MongoStore      StoreType = "Mongo"
	PostgreSQLStore StoreType = "PostgreSQL"
)

type clients struct {
	MySQL *MySQL
	Redis *Redis
}

type Option func(*clients)

// Conn returns a database connection with transaction support
func (tr *Transaction) Conn(ctx context.Context) interface{} {
	if tr == nil {
		return Clients.MySQL.GetDB(ctx)
	}
	if tr.Session == nil {
		tr.Session = Clients.MySQL.GetDB(ctx).Begin(tr.TxOpt)
	}
	return tr.Session.WithContext(ctx)
}

// Begin starts a new transaction
func (tr *Transaction) Begin() error {
	if tr == nil {
		return fmt.Errorf("invalid transaction")
	}
	if tr.Session == nil {
		return fmt.Errorf("invalid session")
	}
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
func NewTransaction(ctx context.Context, store StoreType, opt *sql.TxOptions) (*Transaction, error) {
	tr := &Transaction{TxOpt: opt}

	switch store {
	case MySQLStore:
		session := Clients.MySQL.GetDB(ctx)
		if session == nil {
			return nil, fmt.Errorf("failed to get database session")
		}
		tr.Session = session.Begin(opt)
		if tr.Session.Error != nil {
			return nil, fmt.Errorf("failed to begin transaction: %w", tr.Session.Error)
		}
	case RedisStore:
		// TODO: Implement Redis transaction support
		return nil, fmt.Errorf("redis transaction not implemented")
	case MongoStore:
		// TODO: Implement MongoDB transaction support
		return nil, fmt.Errorf("mongo transaction not implemented")
	case PostgreSQLStore:
		// TODO: Implement PostgreSQL transaction support
		return nil, fmt.Errorf("postgresql transaction not implemented")
	default:
		return nil, fmt.Errorf("unsupported store type: %s", store)
	}

	return tr, nil
}

func (c *clients) close(ctx context.Context) {
	if c.MySQL != nil {
		if err := c.MySQL.Close(ctx); err != nil {
			log.Logger.Error("failed to close MySQL connection", zap.Error(err))
		}
	}
	if c.Redis != nil {
		if err := c.Redis.Close(ctx); err != nil {
			log.Logger.Error("failed to close Redis connection", zap.Error(err))
		}
	}
}

// WithMySQL initializes MySQL client with configuration
func WithMySQL() Option {
	return func(c *clients) {
		if c.MySQL == nil {
			if config.GlobalConfig.MySQL == nil {
				panic("repository init fail, MySQL config is empty")
			}
			c.MySQL = NewMySQLClient()
			if c.MySQL == nil {
				panic("failed to create MySQL client")
			}
		}
	}
}

// WithRedis initializes Redis client with configuration
func WithRedis() Option {
	return func(c *clients) {
		if c.Redis == nil {
			if config.GlobalConfig.Redis == nil {
				panic("repository init fail, Redis config is empty")
			}
			c.Redis = NewRedisClient()
			if c.Redis == nil {
				panic("failed to create Redis client")
			}
		}
	}
}

// Init initializes the repository with the provided options
func Init(opts ...Option) {
	for _, opt := range opts {
		opt(Clients)
	}
	log.Logger.Info("repository init successfully")
}

// Close closes all repository connections
func Close(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	Clients.close(ctx)
	log.Logger.Info("repository closed")
}
