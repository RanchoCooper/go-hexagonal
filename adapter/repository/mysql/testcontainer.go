package mysql

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
)

// MySQLContainerConfig holds configuration for MySQL test container
type MySQLContainerConfig struct {
	DatabaseName string
	User         string
	Password     string
	Port         string
	Host         string
	DSN          string
	Database     string
	CharSet      string
	ParseTime    bool
	TimeZone     string
}

// SetupMySQLContainer creates a MySQL container for testing
func SetupMySQLContainer(t *testing.T) *MySQLContainerConfig {
	t.Helper()

	ctx := context.Background()

	// Create a temporary SQL file with init script
	tempFile, err := os.CreateTemp("", "mysql-init-*.sql")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write SQL schema directly
	initSQL := "CREATE TABLE IF NOT EXISTS `example` (\n" +
		"    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',\n" +
		"    `name` VARCHAR(255) NOT NULL COMMENT 'Name',\n" +
		"    `alias` VARCHAR(255) DEFAULT NULL COMMENT 'Alias',\n" +
		"    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',\n" +
		"    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',\n" +
		"    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Deletion time',\n" +
		"    PRIMARY KEY (`id`),\n" +
		"    KEY `idx_name` (`name`),\n" +
		"    KEY `idx_deleted_at` (`deleted_at`)\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Example table for Hexagonal Architecture';"

	if _, err := tempFile.WriteString(initSQL); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Define MySQL port
	mysqlPort := "3306/tcp"

	// Get the absolute path to the init SQL script
	initScriptPath, err := filepath.Abs(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to get absolute path to init script: %v", err)
	}

	// MySQL container configuration
	containerReq := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{mysqlPort},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "test_db",
			"MYSQL_USER":          "test_user",
			"MYSQL_PASSWORD":      "test_password",
		},
		Mounts: []testcontainers.ContainerMount{
			testcontainers.BindMount(initScriptPath, "/docker-entrypoint-initdb.d/init.sql"),
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL").
			WithStartupTimeout(time.Minute),
	}

	// Start MySQL container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start MySQL container: %v", err)
	}

	// Add cleanup function to terminate container after test
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate MySQL container: %v", err)
		}
	})

	// Get host and port
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get MySQL container host: %v", err)
	}

	port, err := container.MappedPort(ctx, nat.Port(mysqlPort))
	if err != nil {
		t.Fatalf("Failed to get MySQL container port: %v", err)
	}

	// Create config
	config := &MySQLContainerConfig{
		DatabaseName: "test_db",
		User:         "test_user",
		Password:     "test_password",
		Host:         host,
		Port:         port.Port(),
		DSN:          fmt.Sprintf("test_user:test_password@tcp(%s:%s)/test_db?charset=utf8mb4&parseTime=true&loc=Local", host, port.Port()),
		Database:     "test_db",
		CharSet:      "utf8mb4",
		ParseTime:    true,
		TimeZone:     "UTC",
	}

	// Wait a bit for initialization to complete
	time.Sleep(2 * time.Second)

	return config
}

// GetTestDB returns a test MySQL client
func GetTestDB(t *testing.T, config *MySQLContainerConfig) *MySQLClient {
	t.Helper()

	client, err := NewMySQLClient(config.DSN)
	if err != nil {
		t.Fatalf("Failed to create MySQL client: %v", err)
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

// MockMySQLData executes SQL statements in the test database
func MockMySQLData(t *testing.T, db *MySQLClient, sqls []string) {
	t.Helper()

	for _, sql := range sqls {
		tx := db.DB.Exec(sql)
		if tx.Error != nil {
			t.Fatalf("Unable to insert data: %v", tx.Error)
		}
	}
}
