package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"go-hexagonal/adapter/repository"
)

// RedisClient represents a Redis database client
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	if addr == "" {
		return nil, repository.ErrMissingRedisConfig
	}

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     10,
		MinIdleConns: 5,
		IdleTimeout:  5 * time.Minute,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		Client: client,
	}, nil
}

// GetClient returns the Redis client
func (c *RedisClient) GetClient() *redis.Client {
	return c.Client
}

// Close closes the Redis client connection
func (c *RedisClient) Close(ctx context.Context) error {
	if err := c.Client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}
	return nil
}

// ConfigureConnectionPool configures the Redis connection pool
func ConfigureConnectionPool(client *redis.Client, poolSize, minIdleConns int, idleTimeout time.Duration) {
	client.Options().PoolSize = poolSize
	client.Options().MinIdleConns = minIdleConns
	client.Options().IdleTimeout = idleTimeout
}

// HealthCheck performs a Redis health check
func (c *RedisClient) HealthCheck(ctx context.Context) error {
	if err := c.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}
	return nil
}
