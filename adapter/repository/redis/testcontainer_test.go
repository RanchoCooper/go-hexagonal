package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRedisContainer(t *testing.T) {
	// Create Redis mock
	config := SetupRedisContainer(t)

	// Validate configuration
	assert.NotEmpty(t, config.Host, "Host should not be empty")
	assert.NotZero(t, config.Port, "Port should be greater than 0")
	assert.Empty(t, config.Password, "Password should be empty for test instance")
	assert.Equal(t, 0, config.DB, "DB should be 0 for test instance")

	// Validate additional config fields
	assert.Equal(t, 10, config.PoolSize)
	assert.Equal(t, 300, config.IdleTimeout)
	assert.Equal(t, 2, config.MinIdleConns)

	// Get Redis client
	client := GetRedisClient(t, config)

	// Test redis connection by executing simple commands
	ctx := context.Background()
	err := client.Client.Set(ctx, "test_key", "test_value", 0).Err()
	assert.NoError(t, err, "Should be able to set a value")

	val, err := client.Client.Get(ctx, "test_key").Result()
	assert.NoError(t, err, "Should be able to get a value")
	assert.Equal(t, "test_value", val, "Value should match what was set")

	// Test MockRedisData function
	testData := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": "1",
	}

	// Use the MockRedisData function
	MockRedisData(t, client, testData)

	// Verify data was inserted
	val, err = client.Client.Get(ctx, "key1").Result()
	assert.NoError(t, err, "Should be able to get key1")
	assert.Equal(t, "value1", val, "Value should match what was set")

	// Test AssertRedisData function
	AssertRedisData(t, client, "key1", "value1")
	AssertRedisData(t, client, "key2", 42)
	AssertRedisData(t, client, "key3", "1")
}
