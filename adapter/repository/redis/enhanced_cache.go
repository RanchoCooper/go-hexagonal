package redis

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	apperrors "go-hexagonal/util/errors"
)

// CacheOptions defines configuration options for the enhanced cache
type CacheOptions struct {
	// Default TTL for cache entries
	DefaultTTL time.Duration
	// TTL for negative cache entries (for cache miss protection)
	NegativeTTL time.Duration
	// Enable/disable negative caching
	EnableNegativeCache bool
	// Enable/disable key tracking for cache protection
	EnableKeyTracking bool
	// Maximum number of tracked keys
	MaxTrackedKeys int
	// Lock expiration for distributed locks
	LockExpiration time.Duration
	// Lock retry attempts
	LockRetryAttempts int
	// Lock retry delay
	LockRetryDelay time.Duration
	// Lock timeout
	LockTimeout time.Duration
}

// DefaultCacheOptions returns the default cache options
func DefaultCacheOptions() CacheOptions {
	return CacheOptions{
		DefaultTTL:          30 * time.Minute,
		NegativeTTL:         5 * time.Minute,
		EnableNegativeCache: true,
		EnableKeyTracking:   true,
		MaxTrackedKeys:      10000,
		LockExpiration:      5 * time.Second,
		LockRetryAttempts:   5,
		LockRetryDelay:      100 * time.Millisecond,
		LockTimeout:         3 * time.Second,
	}
}

// CacheValue represents a value stored in the cache
type CacheValue struct {
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	// IsNegative indicates this is a negative cache entry (miss)
	IsNegative bool `json:"is_negative"`
}

// EnhancedCache provides advanced caching capabilities
type EnhancedCache struct {
	client  *RedisClient
	options CacheOptions
	// Simple key tracking map
	trackedKeys map[string]struct{}
	keysMutex   sync.RWMutex
}

// NewEnhancedCache creates a new enhanced cache instance
func NewEnhancedCache(client *RedisClient, options CacheOptions) *EnhancedCache {
	cache := &EnhancedCache{
		client:      client,
		options:     options,
		trackedKeys: make(map[string]struct{}, options.MaxTrackedKeys),
	}

	return cache
}

// Get retrieves a value from the cache
func (c *EnhancedCache) Get(ctx context.Context, key string, dest interface{}) error {
	// Check tracked keys first if enabled
	if c.options.EnableKeyTracking {
		c.keysMutex.RLock()
		_, exists := c.trackedKeys[key]
		c.keysMutex.RUnlock()

		if !exists {
			// Key is definitely not in the cache
			return apperrors.New(apperrors.ErrorTypeNotFound, "key not tracked in local cache")
		}
	}

	// Try to get from Redis
	data, err := c.client.Client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return apperrors.New(apperrors.ErrorTypeNotFound, "cache miss")
		}
		return apperrors.Wrapf(err, apperrors.ErrorTypePersistence, "failed to get value from cache: %s", key)
	}

	// Unmarshal the cache value wrapper
	var cacheValue CacheValue
	if err := json.Unmarshal(data, &cacheValue); err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypeSystem, "failed to unmarshal cache value: %s", key)
	}

	// Check if this is a negative cache entry
	if cacheValue.IsNegative {
		return apperrors.New(apperrors.ErrorTypeNotFound, "negative cache hit")
	}

	// Unmarshal the actual data
	if err := json.Unmarshal(cacheValue.Data, dest); err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypeSystem, "failed to unmarshal cached data: %s", key)
	}

	return nil
}

// Set stores a value in the cache
func (c *EnhancedCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Marshal the value
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypeSystem, "failed to marshal value for cache: %s", key)
	}

	// Create cache value wrapper
	cacheValue := CacheValue{
		Data:       valueBytes,
		CreatedAt:  time.Now(),
		IsNegative: false,
	}

	// Marshal the cache value wrapper
	cacheValueBytes, err := json.Marshal(cacheValue)
	if err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypeSystem, "failed to marshal cache value wrapper: %s", key)
	}

	// If TTL is zero, use default
	if ttl == 0 {
		ttl = c.options.DefaultTTL
	}

	// Store in Redis
	if err := c.client.Client.Set(ctx, key, cacheValueBytes, ttl).Err(); err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypePersistence, "failed to set cache value: %s", key)
	}

	// Add to tracked keys if enabled
	if c.options.EnableKeyTracking {
		c.keysMutex.Lock()
		// If map is full, clear half of it (simple LRU approximation)
		if len(c.trackedKeys) >= c.options.MaxTrackedKeys {
			newMap := make(map[string]struct{}, c.options.MaxTrackedKeys)
			c.trackedKeys = newMap
		}
		c.trackedKeys[key] = struct{}{}
		c.keysMutex.Unlock()
	}

	return nil
}

// SetNegative stores a negative cache entry (for cache miss protection)
func (c *EnhancedCache) SetNegative(ctx context.Context, key string) error {
	// Only proceed if negative caching is enabled
	if !c.options.EnableNegativeCache {
		return nil
	}

	// Create negative cache value
	cacheValue := CacheValue{
		Data:       nil,
		CreatedAt:  time.Now(),
		IsNegative: true,
	}

	// Marshal the cache value
	cacheValueBytes, err := json.Marshal(cacheValue)
	if err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypeSystem, "failed to marshal negative cache value: %s", key)
	}

	// Store in Redis with negative TTL
	if err := c.client.Client.Set(ctx, key, cacheValueBytes, c.options.NegativeTTL).Err(); err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypePersistence, "failed to set negative cache: %s", key)
	}

	// Add to tracked keys if enabled
	if c.options.EnableKeyTracking {
		c.keysMutex.Lock()
		c.trackedKeys[key] = struct{}{}
		c.keysMutex.Unlock()
	}

	return nil
}

// Delete removes a value from the cache
func (c *EnhancedCache) Delete(ctx context.Context, key string) error {
	// Delete from Redis
	if err := c.client.Client.Del(ctx, key).Err(); err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypePersistence, "failed to delete from cache: %s", key)
	}

	// Remove from tracked keys if enabled
	if c.options.EnableKeyTracking {
		c.keysMutex.Lock()
		delete(c.trackedKeys, key)
		c.keysMutex.Unlock()
	}

	return nil
}

// DeleteWithPattern removes values matching a pattern from the cache
func (c *EnhancedCache) DeleteWithPattern(ctx context.Context, pattern string) error {
	// Find all keys matching the pattern
	keys, err := c.client.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypePersistence, "failed to get keys with pattern: %s", pattern)
	}

	// Delete all matching keys
	if len(keys) > 0 {
		if err := c.client.Client.Del(ctx, keys...).Err(); err != nil {
			return apperrors.Wrapf(err, apperrors.ErrorTypePersistence, "failed to delete keys with pattern: %s", pattern)
		}

		// Remove from tracked keys if enabled
		if c.options.EnableKeyTracking {
			c.keysMutex.Lock()
			for _, key := range keys {
				delete(c.trackedKeys, key)
			}
			c.keysMutex.Unlock()
		}
	}

	return nil
}

// WithLock executes a function with a distributed lock
func (c *EnhancedCache) WithLock(ctx context.Context, lockKey string, fn func() error) error {
	lockKey = "lock:" + lockKey

	// Try to acquire lock
	acquired, err := c.client.Client.SetNX(ctx, lockKey, "1", c.options.LockExpiration).Result()
	if err != nil {
		return apperrors.Wrapf(err, apperrors.ErrorTypeSystem, "failed to acquire lock: %s", lockKey)
	}

	if !acquired {
		// Lock is already held, wait and retry
		for i := 0; i < c.options.LockRetryAttempts; i++ {
			// Sleep before retry
			time.Sleep(c.options.LockRetryDelay)

			// Try to acquire lock again
			acquired, err = c.client.Client.SetNX(ctx, lockKey, "1", c.options.LockExpiration).Result()
			if err != nil {
				return apperrors.Wrapf(err, apperrors.ErrorTypeSystem, "failed to acquire lock on retry: %s", lockKey)
			}

			if acquired {
				break
			}
		}

		if !acquired {
			return apperrors.New(apperrors.ErrorTypeSystem, "failed to acquire lock after retries")
		}
	}

	// Ensure lock is released
	defer func() {
		_ = c.client.Client.Del(ctx, lockKey).Err()
	}()

	// Execute the function
	return fn()
}

// TryGetSet tries to get a value from cache, if not found executes the loader and sets the result
func (c *EnhancedCache) TryGetSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, loader func() (interface{}, error)) error {
	// Try to get from cache first
	err := c.Get(ctx, key, dest)
	if err == nil {
		// Cache hit, return success
		return nil
	}

	// If it's not a NotFound error, return the error
	if !apperrors.IsNotFoundError(err) {
		return err
	}

	// Execute within a lock to prevent cache stampede
	lockKey := "lock:" + key
	return c.WithLock(ctx, lockKey, func() error {
		// Try to get again (might have been set by another process while waiting for lock)
		err := c.Get(ctx, key, dest)
		if err == nil {
			// Cache hit, return success
			return nil
		}

		// If it's not a NotFound error, return the error
		if !apperrors.IsNotFoundError(err) {
			return err
		}

		// Execute the loader
		result, err := loader()
		if err != nil {
			// If loader failed, set negative cache if enabled
			if c.options.EnableNegativeCache {
				_ = c.SetNegative(ctx, key)
			}
			return err
		}

		// If nil result, set negative cache
		if result == nil {
			if c.options.EnableNegativeCache {
				_ = c.SetNegative(ctx, key)
			}
			return apperrors.New(apperrors.ErrorTypeNotFound, "loader returned nil result")
		}

		// Set the value in cache
		if err := c.Set(ctx, key, result, ttl); err != nil {
			return err
		}

		// Convert the result back to the destination
		resultBytes, err := json.Marshal(result)
		if err != nil {
			return apperrors.Wrapf(err, apperrors.ErrorTypeSystem, "failed to marshal loader result")
		}

		return json.Unmarshal(resultBytes, dest)
	})
}

// RefreshTrackedKeys refreshes the tracked keys from existing Redis keys
func (c *EnhancedCache) RefreshTrackedKeys(ctx context.Context, patterns []string) error {
	if !c.options.EnableKeyTracking {
		return nil
	}

	// Create a new tracked keys map
	newTrackedKeys := make(map[string]struct{}, c.options.MaxTrackedKeys)

	// Scan all keys matching the patterns
	for _, pattern := range patterns {
		var cursor uint64
		for {
			var keys []string
			var err error
			keys, cursor, err = c.client.Client.Scan(ctx, cursor, pattern, 1000).Result()
			if err != nil {
				return apperrors.Wrapf(err, apperrors.ErrorTypePersistence, "failed to scan keys for tracking")
			}

			// Add all keys to the tracking map
			for _, key := range keys {
				// If map is full, stop adding
				if len(newTrackedKeys) >= c.options.MaxTrackedKeys {
					break
				}
				newTrackedKeys[key] = struct{}{}
			}

			if cursor == 0 || len(newTrackedKeys) >= c.options.MaxTrackedKeys {
				break
			}
		}
	}

	// Replace the old tracked keys with the new ones
	c.keysMutex.Lock()
	c.trackedKeys = newTrackedKeys
	c.keysMutex.Unlock()

	return nil
}
