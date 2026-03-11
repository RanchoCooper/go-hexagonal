package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
	"go-hexagonal/util/metrics"
)

const (
	// Example key prefixes for Redis
	exampleKeyPrefix  = "example:id:"
	exampleNamePrefix = "example:name:"

	// Default cache duration
	defaultCacheDuration = 30 * time.Minute
)

// ErrCacheMiss is returned when a requested item is not found in cache
var ErrCacheMiss = errors.New("cache miss")

// exampleIDKey generates a Redis key for example by ID
func exampleIDKey(id int) string {
	return fmt.Sprintf("%s%d", exampleKeyPrefix, id)
}

// exampleNameKey generates a Redis key for example by name
func exampleNameKey(name string) string {
	return fmt.Sprintf("%s%s", exampleNamePrefix, name)
}

// exampleIDStrKey generates a Redis key for example by string ID
func exampleIDStrKey(idStr string) string {
	return fmt.Sprintf("%s%s", exampleKeyPrefix, idStr)
}

// ExampleCacheRepo implements the example cache repository
type ExampleCacheRepo struct {
	client *RedisClient
}

// NewExampleCacheRepo creates a new Redis example cache repository
func NewExampleCacheRepo(client *RedisClient) repo.IExampleCacheRepo {
	return &ExampleCacheRepo{
		client: client,
	}
}

// HealthCheck checks if Redis is available
func (c *ExampleCacheRepo) HealthCheck(ctx context.Context) error {
	return c.client.HealthCheck(ctx)
}

// GetByID gets an example by ID from the cache
func (c *ExampleCacheRepo) GetByID(ctx context.Context, id int) (*model.Example, error) {
	key := exampleIDKey(id)

	// Try to get from cache
	data, err := c.client.Client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			metrics.RecordCacheHit("example_cache", "GetByID_miss")
			return nil, ErrCacheMiss
		}
		metrics.RecordCacheHit("example_cache", "GetByID_error")
		return nil, fmt.Errorf("failed to get example from cache: %w", err)
	}

	// Unmarshal data
	var example model.Example
	if err := json.Unmarshal(data, &example); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached example: %w", err)
	}

	metrics.RecordCacheHit("example_cache", "GetByID_hit")
	return &example, nil
}

// GetByName gets an example by name from the cache
func (c *ExampleCacheRepo) GetByName(ctx context.Context, name string) (*model.Example, error) {
	key := exampleNameKey(name)

	// Try to get example ID from cache
	idStr, err := c.client.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			metrics.RecordCacheHit("example_cache", "GetByName_miss")
			return nil, ErrCacheMiss
		}
		metrics.RecordCacheHit("example_cache", "GetByName_error")
		return nil, fmt.Errorf("failed to get example ID by name from cache: %w", err)
	}

	// Get the example data using the ID
	data, err := c.client.Client.Get(ctx, exampleIDStrKey(idStr)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			metrics.RecordCacheHit("example_cache", "GetByName_miss")
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to get example from cache: %w", err)
	}

	// Unmarshal data
	var example model.Example
	if err := json.Unmarshal(data, &example); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached example: %w", err)
	}

	metrics.RecordCacheHit("example_cache", "GetByName_hit")

	return &example, nil
}

// Set adds or updates an example in the cache
func (c *ExampleCacheRepo) Set(ctx context.Context, example *model.Example) error {
	// Marshal the example to JSON
	data, err := json.Marshal(example)
	if err != nil {
		return fmt.Errorf("failed to marshal example: %w", err)
	}

	// Use pipeline for batch operations
	pipe := c.client.Client.Pipeline()

	// Set the example data
	exampleKey := exampleIDKey(example.Id)
	pipe.Set(ctx, exampleKey, data, defaultCacheDuration)

	// Set the name to ID mapping
	nameKey := exampleNameKey(example.Name)
	pipe.Set(ctx, nameKey, example.Id, defaultCacheDuration)

	// Execute pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		// Check if it's a Redis nil response (not an actual error)
		if err.Error() != "redis: nil" {
			metrics.RecordCacheHit("example_cache", "Set_error")
			return fmt.Errorf("failed to cache example: %w", err)
		}
	}

	metrics.RecordCacheHit("example_cache", "Set_success")
	return nil
}

// Delete removes an example from the cache
func (c *ExampleCacheRepo) Delete(ctx context.Context, id int) error {
	// Get the example to find its name
	example, err := c.GetByID(ctx, id)
	if err != nil && !errors.Is(err, ErrCacheMiss) {
		return fmt.Errorf("failed to get example for deletion: %w", err)
	}

	// Use pipeline for batch operations
	pipe := c.client.Client.Pipeline()

	// Delete the name to ID mapping if the example exists
	if example != nil {
		nameKey := exampleNameKey(example.Name)
		pipe.Del(ctx, nameKey)
	}

	// Delete the example data
	exampleKey := exampleIDKey(id)
	pipe.Del(ctx, exampleKey)

	// Execute pipeline
	_, err = pipe.Exec(ctx)
	if err != nil {
		// Check if it's a Redis nil response (not an actual error)
		if err.Error() != "redis: nil" {
			metrics.RecordCacheHit("example_cache", "Delete_error")
			return fmt.Errorf("failed to delete example: %w", err)
		}
	}

	metrics.RecordCacheHit("example_cache", "Delete_success")
	return nil
}

// Invalidate removes all example related data from the cache
func (c *ExampleCacheRepo) Invalidate(ctx context.Context) error {
	// Find all keys matching the example prefixes
	exampleKeys, err := c.client.Client.Keys(ctx, exampleKeyPrefix+"*").Result()
	if err != nil {
		return fmt.Errorf("failed to get example keys: %w", err)
	}

	nameKeys, err := c.client.Client.Keys(ctx, exampleNamePrefix+"*").Result()
	if err != nil {
		return fmt.Errorf("failed to get example name keys: %w", err)
	}

	// Combine all keys
	var allKeys []string
	allKeys = append(allKeys, exampleKeys...)
	allKeys = append(allKeys, nameKeys...)

	// Delete all keys if there are any
	if len(allKeys) > 0 {
		if err := c.client.Client.Del(ctx, allKeys...).Err(); err != nil {
			return fmt.Errorf("failed to invalidate example cache: %w", err)
		}
	}

	return nil
}
