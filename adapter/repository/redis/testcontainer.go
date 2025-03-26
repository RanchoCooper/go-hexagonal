package redis

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"

	"go-hexagonal/config"
)

// SetupRedisContainer creates a Redis mock for testing using miniredis
func SetupRedisContainer(t *testing.T) *config.RedisConfig {
	t.Helper()

	// Start miniredis server (in-memory implementation)
	s := miniredis.RunT(t)

	// Convert port to int
	portInt, err := strconv.Atoi(s.Port())
	if err != nil {
		t.Fatalf("Failed to convert port to integer: %v", err)
	}

	// Create config using config.RedisConfig
	redisConfig := &config.RedisConfig{
		Host:         s.Host(),
		Port:         portInt,
		Password:     "", // No password for test instance
		DB:           0,
		PoolSize:     10,
		IdleTimeout:  300,
		MinIdleConns: 2,
	}

	// Return server instance via test cleanup to ensure proper shutdown
	t.Cleanup(func() {
		s.Close()
	})

	return redisConfig
}

// GetRedisClient returns a Redis client for testing
func GetRedisClient(t *testing.T, config *config.RedisConfig) *RedisClient {
	t.Helper()

	opts := ClientOptionsFromConfig(config)
	client, err := NewClient(opts)
	if err != nil {
		t.Fatalf("Failed to create Redis client: %v", err)
	}

	return client
}

// MockRedisData adds test data to Redis
func MockRedisData(t *testing.T, client *RedisClient, data map[string]interface{}) {
	t.Helper()

	ctx := context.Background()

	// Clear existing data
	keys, err := client.Client.Keys(ctx, "*").Result()
	if err != nil {
		t.Fatalf("Failed to get Redis keys: %v", err)
	}

	if len(keys) > 0 {
		if _, err := client.Client.Del(ctx, keys...).Result(); err != nil {
			t.Fatalf("Failed to clear Redis data: %v", err)
		}
	}

	// Add the test data
	for k, v := range data {
		if err := client.Client.Set(ctx, k, v, 0).Err(); err != nil {
			t.Fatalf("Failed to set Redis data for key %s: %v", k, err)
		}
	}
}

// AssertRedisData checks if Redis data matches expected values
func AssertRedisData(t *testing.T, client *RedisClient, key string, expected interface{}) {
	t.Helper()

	ctx := context.Background()

	// Get value from Redis
	val, err := client.Client.Get(ctx, key).Result()
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
