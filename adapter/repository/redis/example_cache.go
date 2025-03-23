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
)

const (
	// Example key prefixes for Redis
	exampleKeyPrefix  = "example:id:"
	exampleNamePrefix = "example:name:"

	// Default cache durations
	defaultCacheDuration = 30 * time.Minute
	shortCacheDuration   = 5 * time.Minute
)

// ErrCacheMiss is returned when a requested item is not found in cache
var ErrCacheMiss = errors.New("cache miss")

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
	key := fmt.Sprintf("%s%d", exampleKeyPrefix, id)

	// Try to get from cache
	data, err := c.client.Client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to get example from cache: %w", err)
	}

	// Unmarshal data
	var example model.Example
	if err := json.Unmarshal(data, &example); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached example: %w", err)
	}

	return &example, nil
}

// GetByName gets an example by name from the cache
func (c *ExampleCacheRepo) GetByName(ctx context.Context, name string) (*model.Example, error) {
	key := fmt.Sprintf("%s%s", exampleNamePrefix, name)

	// Try to get example ID from cache
	idStr, err := c.client.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to get example ID by name from cache: %w", err)
	}

	// Get the example data using the ID
	data, err := c.client.Client.Get(ctx, fmt.Sprintf("%s%s", exampleKeyPrefix, idStr)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to get example from cache: %w", err)
	}

	// Unmarshal data
	var example model.Example
	if err := json.Unmarshal(data, &example); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached example: %w", err)
	}

	return &example, nil
}

// Set adds or updates an example in the cache
func (c *ExampleCacheRepo) Set(ctx context.Context, example *model.Example) error {
	// Marshal the example to JSON
	data, err := json.Marshal(example)
	if err != nil {
		return fmt.Errorf("failed to marshal example: %w", err)
	}

	// Set the example data
	exampleKey := fmt.Sprintf("%s%d", exampleKeyPrefix, example.Id)
	if err := c.client.Client.Set(ctx, exampleKey, data, defaultCacheDuration).Err(); err != nil {
		return fmt.Errorf("failed to cache example: %w", err)
	}

	// Set the name to ID mapping
	nameKey := fmt.Sprintf("%s%s", exampleNamePrefix, example.Name)
	if err := c.client.Client.Set(ctx, nameKey, example.Id, defaultCacheDuration).Err(); err != nil {
		return fmt.Errorf("failed to cache example name mapping: %w", err)
	}

	return nil
}

// Delete removes an example from the cache
func (c *ExampleCacheRepo) Delete(ctx context.Context, id int) error {
	// Get the example to find its name
	example, err := c.GetByID(ctx, id)
	if err != nil && !errors.Is(err, ErrCacheMiss) {
		return fmt.Errorf("failed to get example for deletion: %w", err)
	}

	// Delete the name to ID mapping if the example exists
	if example != nil {
		nameKey := fmt.Sprintf("%s%s", exampleNamePrefix, example.Name)
		if err := c.client.Client.Del(ctx, nameKey).Err(); err != nil {
			return fmt.Errorf("failed to delete example name mapping: %w", err)
		}
	}

	// Delete the example data
	exampleKey := fmt.Sprintf("%s%d", exampleKeyPrefix, id)
	if err := c.client.Client.Del(ctx, exampleKey).Err(); err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}

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
