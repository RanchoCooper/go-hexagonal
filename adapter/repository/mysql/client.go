package mysql

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-hexagonal/adapter/repository"
)

// MySQLClient represents a MySQL database client using GORM
type MySQLClient struct {
	DB *gorm.DB
}

// NewMySQLClient creates a new MySQL client
func NewMySQLClient(dsn string) (*MySQLClient, error) {
	if dsn == "" {
		return nil, repository.ErrMissingMySQLConfig
	}

	db, err := openMySQLDB(dsn)
	if err != nil {
		return nil, err
	}

	return &MySQLClient{DB: db}, nil
}

// GetDB returns the GORM database instance with context
func (c *MySQLClient) GetDB(ctx context.Context) *gorm.DB {
	return c.DB.WithContext(ctx)
}

// SetDB sets the GORM database instance
func (c *MySQLClient) SetDB(db *gorm.DB) {
	c.DB = db
}

// Close closes the MySQL database connection
func (c *MySQLClient) Close(ctx context.Context) error {
	sqlDB, err := c.GetDB(ctx).DB()
	if err != nil {
		return fmt.Errorf("failed to get MySQL DB: %w", err)
	}

	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close MySQL connection: %w", err)
		}
	}

	return nil
}

// openMySQLDB creates and opens a new GORM database connection
func openMySQLDB(dsn string) (*gorm.DB, error) {
	// Create MySQL dialect
	dialect := mysql.Open(dsn)

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Do not ignore ErrRecordNotFound
			Colorful:                  true,        // Enable colorful output
		},
	)

	// GORM configuration
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         gormLogger,
	}

	// Open database connection
	db, err := gorm.Open(dialect, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Default connection pool settings
	// In a real application, these should be configured from config files
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return db, nil
}

// ConfigureConnectionPool configures the MySQL connection pool
func ConfigureConnectionPool(db *gorm.DB, maxIdleConns, maxOpenConns int, maxLifetime, maxIdleTime time.Duration) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(cast.ToDuration(maxLifetime))
	sqlDB.SetConnMaxIdleTime(cast.ToDuration(maxIdleTime))

	return nil
}
