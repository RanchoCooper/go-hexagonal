package repository

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-hexagonal/util/log"

	"github.com/go-redis/redis/v8"
)

// DefaultRepositoryTimeout defines the default context timeout for database operations
const DefaultRepositoryTimeout = 30 * time.Second

// ClientContainer holds all repository client instances
type ClientContainer struct {
	MySQL      *MySQL
	Redis      *Redis
	PostgreSQL *PostgreSQL
}

// Global instance for backward compatibility
// Note: It's recommended to use dependency injection with Wire instead of this global instance
var Clients = &ClientContainer{}

// ISQLClient defines the interface for SQL database clients (MySQL, PostgreSQL)
type ISQLClient interface {
	// GetDB returns the database instance with context
	GetDB(ctx context.Context) interface{}

	// SetDB sets the database instance
	SetDB(db interface{})

	// Close closes the database connection
	Close(ctx context.Context) error
}

// IRedisClient defines the interface for Redis clients
type IRedisClient interface {
	// Close closes the Redis connection
	Close(ctx context.Context) error
}

// Initialize creates a new client container if not already initialized
func Initialize() {
	if Clients == nil {
		Clients = &ClientContainer{}
	}
}

// Close closes all repository connections
func (c *ClientContainer) Close(ctx context.Context) {
	if c.MySQL != nil {
		if err := c.MySQL.Close(ctx); err != nil {
			log.Logger.Error("failed to close MySQL connection", zap.Error(err))
		}
	}
	if c.PostgreSQL != nil {
		if err := c.PostgreSQL.Close(ctx); err != nil {
			log.Logger.Error("failed to close PostgreSQL connection", zap.Error(err))
		}
	}
	if c.Redis != nil {
		if err := c.Redis.Close(ctx); err != nil {
			log.Logger.Error("failed to close Redis connection", zap.Error(err))
		}
	}
}

// Close closes all repository connections
func Close(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, DefaultRepositoryTimeout)
	defer cancel()

	Clients.Close(ctx)

	log.Logger.Info("repository closed")
}

// MySQL represents a MySQL database client
type MySQL struct {
	DB *gorm.DB
}

// SetDB sets the GORM database connection
func (m *MySQL) SetDB(db *gorm.DB) {
	m.DB = db
}

// GetDB returns the GORM database connection
func (m *MySQL) GetDB(ctx context.Context) *gorm.DB {
	if m.DB == nil {
		return nil
	}
	return m.DB.WithContext(ctx)
}

// Close closes the MySQL connection
func (m *MySQL) Close(ctx context.Context) error {
	// No-op for now, as GORM manages connection pooling
	return nil
}

// NewMySQLClient creates a new MySQL client
func NewMySQLClient(db *gorm.DB) *MySQL {
	return &MySQL{DB: db}
}

// PostgreSQL represents a PostgreSQL database client
type PostgreSQL struct {
	DB *gorm.DB
}

// SetDB sets the GORM database connection
func (p *PostgreSQL) SetDB(db *gorm.DB) {
	p.DB = db
}

// GetDB returns the GORM database connection
func (p *PostgreSQL) GetDB(ctx context.Context) *gorm.DB {
	if p.DB == nil {
		return nil
	}
	return p.DB.WithContext(ctx)
}

// Close closes the PostgreSQL connection
func (p *PostgreSQL) Close(ctx context.Context) error {
	// No-op for now, as GORM manages connection pooling
	return nil
}

// NewPostgreSQLClient creates a new PostgreSQL client
func NewPostgreSQLClient(db *gorm.DB) *PostgreSQL {
	return &PostgreSQL{DB: db}
}

// Redis represents a Redis client
type Redis struct {
	DB *redis.Client
}

// Close closes the Redis connection
func (r *Redis) Close(ctx context.Context) error {
	if r.DB != nil {
		if err := r.DB.Close(); err != nil {
			return fmt.Errorf("failed to close Redis connection: %w", err)
		}
	}
	return nil
}

// NewRedisClient creates a new Redis client
func NewRedisClient() *Redis {
	return &Redis{}
}
