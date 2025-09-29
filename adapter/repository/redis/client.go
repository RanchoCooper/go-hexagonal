package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

// ClientOptions holds Redis client configuration options
type ClientOptions struct {
	// Address is the Redis server address
	Address string
	// Password is the Redis server password
	Password string
	// DB is the Redis database index
	DB int
	// PoolSize is the maximum number of socket connections
	PoolSize int
	// MinIdleConns is the minimum number of idle connections
	MinIdleConns int
	// DialTimeout is the timeout for establishing new connections
	DialTimeout time.Duration
	// ReadTimeout is the timeout for socket reads
	ReadTimeout time.Duration
	// WriteTimeout is the timeout for socket writes
	WriteTimeout time.Duration
	// PoolTimeout is the timeout for getting a connection from the pool
	PoolTimeout time.Duration
	// IdleTimeout is the timeout for idle connections
	IdleTimeout time.Duration
	// MaxRetries is the maximum number of retries before giving up
	MaxRetries int
	// MinRetryBackoff is the minimum backoff between retries
	MinRetryBackoff time.Duration
	// MaxRetryBackoff is the maximum backoff between retries
	MaxRetryBackoff time.Duration
}

// DefaultClientOptions returns the default Redis client options
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Address:         "localhost:6379",
		Password:        "",
		DB:              0,
		PoolSize:        10,
		MinIdleConns:    5,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     4 * time.Second,
		IdleTimeout:     5 * time.Minute,
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	}
}

// ClientOptionsFromConfig creates client options from application config
func ClientOptionsFromConfig(cfg *config.RedisConfig) *ClientOptions {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	opts := DefaultClientOptions()
	opts.Address = addr
	opts.Password = cfg.Password
	opts.DB = cfg.DB

	if cfg.PoolSize > 0 {
		opts.PoolSize = cfg.PoolSize
	}

	if cfg.MinIdleConns > 0 {
		opts.MinIdleConns = cfg.MinIdleConns
	}

	if cfg.IdleTimeout > 0 {
		opts.IdleTimeout = time.Duration(cfg.IdleTimeout) * time.Second
	}

	return opts
}

// RedisClient wraps a Redis client with additional functionality
type RedisClient struct {
	Client *redis.Client
	opts   *ClientOptions
}

// NewClient creates a new Redis client with the given options
func NewClient(opts *ClientOptions) (*RedisClient, error) {
	if opts == nil {
		opts = DefaultClientOptions()
	}

	redisOpts := &redis.Options{
		Addr:            opts.Address,
		Password:        opts.Password,
		DB:              opts.DB,
		PoolSize:        opts.PoolSize,
		MinIdleConns:    opts.MinIdleConns,
		DialTimeout:     opts.DialTimeout,
		ReadTimeout:     opts.ReadTimeout,
		WriteTimeout:    opts.WriteTimeout,
		PoolTimeout:     opts.PoolTimeout,
		IdleTimeout:     opts.IdleTimeout,
		MaxRetries:      opts.MaxRetries,
		MinRetryBackoff: opts.MinRetryBackoff,
		MaxRetryBackoff: opts.MaxRetryBackoff,
	}

	client := redis.NewClient(redisOpts)
	redisClient := &RedisClient{
		Client: client,
		opts:   opts,
	}

	// Verify connection on creation
	if err := redisClient.HealthCheck(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return redisClient, nil
}

// HealthCheck performs a ping to verify the Redis connection is working
func (c *RedisClient) HealthCheck(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.opts.DialTimeout)
	defer cancel()

	status := c.Client.Ping(timeoutCtx)
	if status.Err() != nil {
		return fmt.Errorf("redis health check failed: %w", status.Err())
	}
	return nil
}

// Close closes the Redis connection
func (c *RedisClient) Close() error {
	log.Logger.Info("Closing Redis connection")
	return c.Client.Close()
}

// Stats returns the Redis client connection stats
func (c *RedisClient) Stats() *redis.PoolStats {
	return c.Client.PoolStats()
}

// WithTimeout returns a new context with timeout
func (c *RedisClient) WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

// NewClientFromConfig creates a new Redis client from application config
func NewClientFromConfig(cfg *config.RedisConfig) (*RedisClient, error) {
	opts := ClientOptionsFromConfig(cfg)
	return NewClient(opts)
}
