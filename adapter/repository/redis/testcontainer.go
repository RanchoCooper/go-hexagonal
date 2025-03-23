package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RedisContainerConfig holds configuration for Redis test container
type RedisContainerConfig struct {
	Password string
	Host     string
	Port     string
	DB       int
}

// SetupRedisContainer creates a Redis container for testing
func SetupRedisContainer(t *testing.T) *RedisContainerConfig {
	t.Helper()

	ctx := context.Background()

	// Define Redis port
	redisPort := "6379/tcp"

	// Redis container configuration
	containerReq := testcontainers.ContainerRequest{
		Image:        "redis:6-alpine",
		ExposedPorts: []string{redisPort},
		Env: map[string]string{
			"REDIS_PASSWORD": "", // No password for test container
		},
		Cmd: []string{"redis-server", "--requirepass", ""},
		WaitingFor: wait.ForLog("Ready to accept connections").
			WithStartupTimeout(time.Minute),
	}

	// Start Redis container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start Redis container: %v", err)
	}

	// Add cleanup function to terminate container after test
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate Redis container: %v", err)
		}
	})

	// Get host and port
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get Redis container host: %v", err)
	}

	port, err := container.MappedPort(ctx, nat.Port(redisPort))
	if err != nil {
		t.Fatalf("Failed to get Redis container port: %v", err)
	}

	// Create config
	config := &RedisContainerConfig{
		Password: "",
		Host:     host,
		Port:     port.Port(),
		DB:       0,
	}

	// Wait a bit for initialization to complete
	time.Sleep(2 * time.Second)

	return config
}

// GetRedisClient returns a test Redis client
func GetRedisClient(t *testing.T, config *RedisContainerConfig) *RedisClient {
	t.Helper()

	client, err := NewRedisClient(
		fmt.Sprintf("%s:%s", config.Host, config.Port),
		config.Password,
		config.DB,
	)
	if err != nil {
		t.Fatalf("Failed to create Redis client: %v", err)
	}

	return client
}

// MockRedisData adds test data to Redis
func MockRedisData(t *testing.T, client *RedisClient, data map[string]interface{}) {
	t.Helper()

	ctx := context.Background()

	// Clear any existing data
	redisClient := client.GetClient()
	keys, err := redisClient.Keys(ctx, "*").Result()
	if err != nil {
		t.Fatalf("Failed to get Redis keys: %v", err)
	}

	if len(keys) > 0 {
		if _, err := redisClient.Del(ctx, keys...).Result(); err != nil {
			t.Fatalf("Failed to clear Redis data: %v", err)
		}
	}

	// Add test data with default expiration
	for key, value := range data {
		if err := redisClient.Set(ctx, key, value, 10*time.Minute).Err(); err != nil {
			t.Fatalf("Failed to set Redis data for key %s: %v", key, err)
		}
	}
}

// AssertRedisData checks if the Redis data matches expected values
func AssertRedisData(t *testing.T, client *RedisClient, key string, expected interface{}) {
	t.Helper()

	ctx := context.Background()

	// Get the value
	val, err := client.GetClient().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			t.Fatalf("Key %s does not exist in Redis", key)
		}
		t.Fatalf("Failed to get Redis data for key %s: %v", key, err)
	}

	// Convert expected to string for comparison
	expectedStr := fmt.Sprintf("%v", expected)

	// Compare values
	if val != expectedStr {
		t.Fatalf("Redis value mismatch for key %s. Expected: %v, Got: %s", key, expected, val)
	}
}
