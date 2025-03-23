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
)

const (
	// PostgreSQLStartTimeout defines the timeout for starting the PostgreSQL container
	PostgreSQLStartTimeout = 2 * time.Minute
)

// PostgreConfig stores PostgreSQL connection configuration
type PostgreConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
	SSLMode  string
	TimeZone string
}

// PostgreSQLContainerConfig holds configuration for PostgreSQL test container
type PostgreSQLContainerConfig struct {
	DatabaseName string
	User         string
	Password     string
	Port         string
	Host         string
	DSN          string
	Database     string
	SSLMode      string
	TimeZone     string
}

// SetupPostgreSQLContainer creates and starts a PostgreSQL test container
func SetupPostgreSQLContainer(t *testing.T) *PostgreSQLContainerConfig {
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
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
			"POSTGRES_DB":       "test_db",
		},
		Mounts: []testcontainers.ContainerMount{
			testcontainers.BindMount(initScriptPath, "/docker-entrypoint-initdb.d/init.sql"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithStartupTimeout(time.Minute),
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

	// Create config
	config := &PostgreSQLContainerConfig{
		DatabaseName: "test_db",
		User:         "test_user",
		Password:     "test_password",
		Host:         host,
		Port:         port.Port(),
		DSN:          fmt.Sprintf("host=%s port=%s user=test_user password=test_password dbname=test_db sslmode=disable", host, port.Port()),
		Database:     "test_db",
		SSLMode:      "disable",
		TimeZone:     "UTC",
	}

	// Wait a bit for initialization to complete
	time.Sleep(2 * time.Second)

	return config
}

// GetTestDB creates a GORM connection based on PostgreSQL configuration
func GetTestDB(t *testing.T, config *PostgreSQLContainerConfig) *PostgreSQLClient {
	t.Helper()

	client, err := NewPostgreSQLClient(config.DSN)
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
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
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
	// Execute all SQL statements
	for _, sql := range sqls {
		tx := db.Exec(sql)
		if tx.Error != nil {
			t.Fatalf("Unable to insert data: %v", tx.Error)
		}
	}
}
