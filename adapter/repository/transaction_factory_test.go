package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-hexagonal/domain/repo"
)

const (
	mysqlStartTimeout = 2 * time.Minute
)

// setupMySQLContainer creates a MySQL container for testing
func setupMySQLContainer(t *testing.T) (string, int, string, string, string) {
	t.Log("Setting up an instance of MySQL with testcontainers-go")
	ctx := context.Background()

	user, password, dbName := "user", "123456", "test"

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_USER":          user,
			"MYSQL_ROOT_PASSWORD": password,
			"MYSQL_PASSWORD":      password,
			"MYSQL_DATABASE":      dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("3306/tcp").WithStartupTimeout(mysqlStartTimeout),
			wait.ForLog("ready for connections"),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("could not start Docker container, err: %s", err)
	}

	t.Cleanup(func() {
		t.Log("Removing MySQL container from Docker")
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate MySQL container, err: %s", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host where the container is exposed, err: %s", err)
	}

	port, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		t.Fatalf("failed to get externally mapped port to MySQL database, err: %s", err)
	}

	portInt := port.Int()
	t.Logf("MySQL container running at %s:%d", host, portInt)

	return host, portInt, user, password, dbName
}

// getTestDB creates a MySQL test database connection
func getTestDB(t *testing.T, host string, port int, user, password, dbName string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
		user, password, host, port, dbName)

	// Create GORM connection
	dialector := mysql.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

func TestNewTransactionFactory(t *testing.T) {
	// 创建一个空的clients映射
	clients := map[StoreType]any{}

	// Test that NewTransactionFactory returns a non-nil implementation of repo.TransactionFactory
	factory := NewTransactionFactory(clients)
	assert.NotNil(t, factory, "Factory should not be nil")
	assert.IsType(t, &TransactionFactoryImpl{}, factory, "Factory should be of type *TransactionFactoryImpl")
}

func TestTransactionFactoryImpl_NewTransaction(t *testing.T) {
	// Use testcontainers to create MySQL test container
	host, port, user, password, dbName := setupMySQLContainer(t)

	// Get test database connection
	db := getTestDB(t, host, port, user, password, dbName)

	// Save original clients
	originalClients := Clients
	defer func() {
		// Restore original clients after test
		Clients = originalClients
	}()

	// Set MySQL client
	mysqlClient := &MySQL{}
	mysqlClient.SetDB(db)

	// Set test clients
	Clients = &ClientContainer{
		MySQL: mysqlClient,
	}

	// Create test table
	err := db.Exec("CREATE TABLE IF NOT EXISTS test_transaction (id INT PRIMARY KEY, name VARCHAR(255))").Error
	assert.NoError(t, err, "Should successfully create test table")

	// Create clients map
	clients := map[StoreType]any{
		MySQLStore: mysqlClient,
	}

	// Create transaction factory
	factory := NewTransactionFactory(clients)

	t.Run("MySQL Transaction", func(t *testing.T) {
		// Create context
		ctx := context.Background()

		// Create transaction options
		txOpts := &sql.TxOptions{
			Isolation: sql.LevelDefault,
			ReadOnly:  false,
		}

		// Create transaction
		tr, err := factory.NewTransaction(ctx, repo.MySQLStore, txOpts)

		// Should be successful
		assert.NoError(t, err)
		assert.NotNil(t, tr)

		// Test transaction commit
		err = tr.Begin()
		assert.NoError(t, err)

		// Execute a simple query in transaction
		sqlDB, err := db.DB()
		assert.NoError(t, err)

		// Insert data in transaction
		_, err = sqlDB.Exec("INSERT INTO test_transaction (id, name) VALUES (1, 'test')")
		assert.NoError(t, err)

		// Commit transaction
		err = tr.Commit()
		assert.NoError(t, err)

		// Check if data was inserted
		var count int
		err = db.Raw("SELECT COUNT(*) FROM test_transaction WHERE id = 1").Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, 1, count, "Should successfully insert one record")
	})

	t.Run("Unsupported Store Type", func(t *testing.T) {
		// Create context
		ctx := context.Background()

		// Create transaction with unsupported store type
		tr, err := factory.NewTransaction(ctx, "UnsupportedStore", nil)

		// Should return client not found error
		assert.Error(t, err)
		assert.Nil(t, tr)
		assert.Contains(t, err.Error(), "no client found for store type")
	})

	t.Run("Redis Transaction Not Implemented", func(t *testing.T) {
		// Create context
		ctx := context.Background()

		// Add Redis client
		redisClient := NewRedisClient()
		clients[RedisStore] = redisClient

		// Try to create Redis transaction
		tr, err := factory.NewTransaction(ctx, repo.RedisStore, nil)

		// Should return method not supported error
		assert.Error(t, err)
		assert.Nil(t, tr)
		assert.Contains(t, err.Error(), "redis transaction not implemented")
	})

	t.Run("Non-SQL Transaction Options", func(t *testing.T) {
		// Create context
		ctx := context.Background()

		// Create transaction with non-SQL options
		nonSqlOpts := "not sql options"
		tr, err := factory.NewTransaction(ctx, repo.MySQLStore, nonSqlOpts)

		// Should create transaction with default options
		assert.NoError(t, err)
		assert.NotNil(t, tr)
	})
}
