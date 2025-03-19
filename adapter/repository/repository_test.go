package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ctx = context.TODO()

func TestNewRepository(t *testing.T) {
	// Clean up any existing clients first
	Clients = &clients{}

	// Initialize new connections
	Init(WithMySQL(), WithRedis())
	defer Close(ctx)
}

func TestTransaction_Conn(t *testing.T) {
	// Clean up any existing clients first
	Clients = &clients{
		MySQL: &MySQL{},
		Redis: &Redis{},
	}

	// Create mock database connection
	sqlDB, sqlMock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Set up sqlmock expectations
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	// Manually set up the mock MySQL client
	mockDB, err := createMockDB(sqlDB)
	assert.NoError(t, err)
	Clients.MySQL.SetDB(mockDB)

	// Create transaction
	tr, err := NewTransaction(ctx, MySQLStore, nil)
	assert.NoError(t, err)
	assert.NotNil(t, tr)

	// Test transaction connection
	db := tr.Conn(ctx)
	assert.NotNil(t, db)

	// Rollback transaction
	err = tr.Rollback()
	assert.NoError(t, err)

	// Verify all expectations are satisfied
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

// createMockDB creates a mock GORM database connection
func createMockDB(db *sql.DB) (*gorm.DB, error) {
	dialect := driver.New(
		driver.Config{
			Conn:                      db,
			DriverName:                "mysql-mock",
			SkipInitializeWithVersion: true,
		},
	)

	return gorm.Open(dialect, buildGormConfig())
}
