package postgre

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"

	"go-hexagonal/config"
)

const (
	// PostgreSQLStartTimeout defines the timeout for starting the PostgreSQL container
	PostgreSQLStartTimeout = 2 * time.Minute
)

// SetupPostgreSQLContainer creates and starts a PostgreSQL test container
func SetupPostgreSQLContainer(t *testing.T) *config.PostgreSQLConfig {
	t.Helper()

	ctx := context.Background()

	// Create a temporary SQL file with init script
	tempFile, err := os.CreateTemp("", "postgres-init-*.sql")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write SQL schema directly - note PostgreSQL syntax differences from MySQL
	initSQL := "CREATE TABLE IF NOT EXISTS example (\n" +
		"    id SERIAL PRIMARY KEY,\n" +
		"    name VARCHAR(255) NOT NULL,\n" +
		"    alias VARCHAR(255),\n" +
		"    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
		"    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
		"    deleted_at TIMESTAMP\n" +
		");\n\n" +
		"CREATE INDEX idx_example_name ON example(name);\n" +
		"CREATE INDEX idx_example_deleted_at ON example(deleted_at);\n" +
		"COMMENT ON TABLE example IS 'Example table for Hexagonal Architecture';\n" +
		"COMMENT ON COLUMN example.id IS 'Primary key ID';\n" +
		"COMMENT ON COLUMN example.name IS 'Name';\n" +
		"COMMENT ON COLUMN example.alias IS 'Alias';\n" +
		"COMMENT ON COLUMN example.created_at IS 'Creation time';\n" +
		"COMMENT ON COLUMN example.updated_at IS 'Update time';\n" +
		"COMMENT ON COLUMN example.deleted_at IS 'Deletion time';"

	if _, err := tempFile.WriteString(initSQL); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Define PostgreSQL port
	postgresPort := "5432/tcp"

	// Get the absolute path to the init SQL script
	initScriptPath, err := filepath.Abs(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to get absolute path to init script: %v", err)
	}

	// PostgreSQL container configuration
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{postgresPort},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "123456",
			"POSTGRES_DB":       "postgres",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      initScriptPath,
				ContainerFilePath: "/docker-entrypoint-initdb.d/init.sql",
				FileMode:          0644,
			},
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForExec([]string{"pg_isready"}).
				WithPollInterval(1*time.Second).
				WithExitCodeMatcher(func(exitCode int) bool {
					return exitCode == 0
				}),
			wait.ForLog("database system is ready to accept connections"),
		).WithStartupTimeout(PostgreSQLStartTimeout),
	}

	// Start PostgreSQL container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// Add cleanup function to terminate container after test
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate PostgreSQL container: %v", err)
		}
	})

	// Get host and port
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get PostgreSQL container host: %v", err)
	}

	port, err := container.MappedPort(ctx, nat.Port(postgresPort))
	if err != nil {
		t.Fatalf("Failed to get PostgreSQL container port: %v", err)
	}

	// Get port as integer
	portInt := port.Int()

	// Create config using config.PostgreSQLConfig
	postgresConfig := &config.PostgreSQLConfig{
		User:            "postgres",
		Password:        "123456",
		Host:            host,
		Port:            portInt,
		Database:        "postgres",
		SSLMode:         "disable",
		Options:         "",
		MaxConnections:  100,
		MinConnections:  10,
		MaxConnLifetime: 3600,
		IdleTimeout:     300,
		ConnectTimeout:  10,
		TimeZone:        "UTC",
	}

	// Create the test database and user
	dsn := fmt.Sprintf("host=%s port=%d user=postgres password=123456 dbname=postgres sslmode=disable",
		host, portInt)

	// Create a temporary client to setup test database and user
	tmpClient, err := NewPostgreSQLClient(dsn)
	if err != nil {
		t.Fatalf("Failed to create temporary PostgreSQL client: %v", err)
	}

	// Create test database and user - use simple statements that won't fail if objects already exist
	queries := []string{
		"SELECT 1", // Replace with appropriate statements for your database
		// Skip database creation as we use the default postgres database
		// Skip user creation as we use the postgres user
	}

	for _, query := range queries {
		if err := tmpClient.DB.Exec(query).Error; err != nil {
			t.Fatalf("Failed to execute query %s: %v", query, err)
		}
	}

	// Wait a bit for initialization to complete
	time.Sleep(2 * time.Second)

	return postgresConfig
}

// GetTestDB creates a GORM connection based on PostgreSQL configuration
func GetTestDB(t *testing.T, config *config.PostgreSQLConfig) *PostgreSQLClient {
	t.Helper()

	// Create DSN
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode,
		config.TimeZone,
	)

	client, err := NewPostgreSQLClient(dsn)
	if err != nil {
		t.Fatalf("Failed to create PostgreSQL client: %v", err)
	}

	// Auto migrate to ensure schema is up to date
	if err := client.DB.AutoMigrate(&ExampleTable{}); err != nil {
		t.Fatalf("Failed to run auto migration: %v", err)
	}

	return client
}

// ExampleTable represents the example table for testing
type ExampleTable struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;type:varchar(255);not null"`
	Alias     string    `gorm:"column:alias;type:varchar(255)"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt time.Time `gorm:"column:deleted_at;index"`
}

// TableName returns the table name
func (ExampleTable) TableName() string {
	return "example"
}

// MockPostgreSQLData executes SQL statements in the PostgreSQL database
func MockPostgreSQLData(t *testing.T, db *gorm.DB, sqls []string) {
	t.Helper()

	// Execute all SQL statements
	for _, sql := range sqls {
		tx := db.Exec(sql)
		if tx.Error != nil {
			t.Fatalf("Unable to insert data: %v", tx.Error)
		}
	}
}
