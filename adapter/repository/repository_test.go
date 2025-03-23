package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-hexagonal/adapter/repository"
	"go-hexagonal/util/log"
)

// MockDB is a mock implementation of gorm.DB
type MockDB struct {
	mock.Mock
}

func TestNewMySQLClient(t *testing.T) {
	// Create a nil db
	var db *gorm.DB = nil

	// Create MySQL client
	mysqlClient := repository.NewMySQLClient(db)

	// Verify client is not nil
	assert.NotNil(t, mysqlClient)

	// Verify GetDB returns nil (because db itself is nil)
	ctx := context.Background()
	assert.Nil(t, mysqlClient.GetDB(ctx))

	// Verify Close method doesn't return an error
	err := mysqlClient.Close(ctx)
	assert.NoError(t, err)
}

func TestNewPostgreSQLClient(t *testing.T) {
	// Create a nil db
	var db *gorm.DB = nil

	// Create PostgreSQL client
	pgClient := repository.NewPostgreSQLClient(db)

	// Verify client is not nil
	assert.NotNil(t, pgClient)

	// Verify GetDB returns nil (because db itself is nil)
	ctx := context.Background()
	assert.Nil(t, pgClient.GetDB(ctx))

	// Verify Close method doesn't return an error
	err := pgClient.Close(ctx)
	assert.NoError(t, err)
}

func TestNewRedisClient(t *testing.T) {
	// Create Redis client
	redisClient := repository.NewRedisClient()

	// Verify client is not nil
	assert.NotNil(t, redisClient)

	// Verify Close method doesn't return an error
	ctx := context.Background()
	err := redisClient.Close(ctx)
	assert.NoError(t, err)
}

func TestClientContainer_Close(t *testing.T) {
	// Create ClientContainer and set test clients
	container := &repository.ClientContainer{
		MySQL:      repository.NewMySQLClient(nil),
		PostgreSQL: repository.NewPostgreSQLClient(nil),
		Redis:      repository.NewRedisClient(),
	}

	// Verify Close method doesn't panic
	ctx := context.Background()
	assert.NotPanics(t, func() {
		container.Close(ctx)
	})
}

func TestClose(t *testing.T) {
	// Save original Logger
	originalLogger := log.Logger
	defer func() {
		// Restore original Logger after test
		log.Logger = originalLogger
	}()

	// Set a temporary Logger
	log.Logger = zap.NewNop()

	// Ensure Close function doesn't panic
	ctx := context.Background()
	assert.NotPanics(t, func() {
		repository.Close(ctx)
	})
}
