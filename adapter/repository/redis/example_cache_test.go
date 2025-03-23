package redis

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"

	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

var testCtx = context.Background()

// ExampleCache is a simplified implementation for testing
type ExampleCache struct {
	client redis.Cmdable
}

// NewExampleCache creates a new ExampleCache for testing
func NewExampleCache(client redis.Cmdable) repo.IExampleCacheRepo {
	return &ExampleCache{
		client: client,
	}
}

// HealthCheck implements the health check functionality
func (c *ExampleCache) HealthCheck(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Required interface methods with stub implementations for testing
func (c *ExampleCache) GetByID(ctx context.Context, id int) (*model.Example, error) {
	return &model.Example{}, nil
}

func (c *ExampleCache) GetByName(ctx context.Context, name string) (*model.Example, error) {
	return &model.Example{}, nil
}

func (c *ExampleCache) Set(ctx context.Context, example *model.Example) error {
	return nil
}

func (c *ExampleCache) Delete(ctx context.Context, id int) error {
	return nil
}

func (c *ExampleCache) Invalidate(ctx context.Context) error {
	return nil
}

func TestExampleCache_HealthCheck(t *testing.T) {
	// Create Redis client and mock
	db, mock := redismock.NewClientMock()
	mock.ExpectPing().SetVal("PONG")

	// Create Redis cache instance
	cache := NewExampleCache(db)

	// Execute health check
	err := cache.HealthCheck(testCtx)
	assert.Nil(t, err)

	// Verify mock expectations
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
