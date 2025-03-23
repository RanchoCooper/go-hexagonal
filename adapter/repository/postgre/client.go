package postgre

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-hexagonal/adapter/repository"
	"go-hexagonal/config"
)

// PostgreSQLClient represents a PostgreSQL database client using GORM
type PostgreSQLClient struct {
	DB *gorm.DB
}

// NewPostgreSQLClient creates a new PostgreSQL client
func NewPostgreSQLClient(dsn string) (*PostgreSQLClient, error) {
	if dsn == "" {
		return nil, repository.ErrMissingPostgreSQLConfig
	}

	db, err := openPostgreSQLDB(dsn)
	if err != nil {
		return nil, err
	}

	return &PostgreSQLClient{DB: db}, nil
}

// GetDB returns the GORM database instance with context
func (c *PostgreSQLClient) GetDB(ctx context.Context) *gorm.DB {
	return c.DB.WithContext(ctx)
}

// SetDB sets the GORM database instance
func (c *PostgreSQLClient) SetDB(db *gorm.DB) {
	c.DB = db
}

// Close closes the PostgreSQL database connection
func (c *PostgreSQLClient) Close(ctx context.Context) error {
	sqlDB, err := c.GetDB(ctx).DB()
	if err != nil {
		return fmt.Errorf("failed to get PostgreSQL DB: %w", err)
	}

	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close PostgreSQL connection: %w", err)
		}
	}

	return nil
}

// openPostgreSQLDB creates and opens a new GORM database connection for PostgreSQL
func openPostgreSQLDB(dsn string) (*gorm.DB, error) {
	// Create PostgreSQL dialect
	dialect := postgres.Open(dsn)

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
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
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

// ConfigureConnectionPool configures the PostgreSQL connection pool
func ConfigureConnectionPool(db *gorm.DB, maxIdleConns, maxOpenConns int, maxLifetime, maxIdleTime time.Duration) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	return nil
}

// NewConnPool creates a new PostgreSQL connection pool using pgx
func NewConnPool(pgConfig *config.PostgreSQLConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		pgConfig.User,
		pgConfig.Password,
		pgConfig.Host,
		pgConfig.Port,
		pgConfig.Database,
		pgConfig.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Configure connection pool
	if pgConfig.MaxConnections > 0 {
		poolConfig.MaxConns = int32(pgConfig.MaxConnections)
	}
	if pgConfig.MinConnections > 0 {
		poolConfig.MinConns = int32(pgConfig.MinConnections)
	}
	if pgConfig.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = time.Duration(pgConfig.MaxConnLifetime) * time.Second
	}
	if pgConfig.IdleTimeout > 0 {
		poolConfig.MaxConnIdleTime = time.Duration(pgConfig.IdleTimeout) * time.Second
	}
	if pgConfig.ConnectTimeout > 0 {
		poolConfig.ConnConfig.ConnectTimeout = time.Duration(pgConfig.ConnectTimeout) * time.Second
	}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}
