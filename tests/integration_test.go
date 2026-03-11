package tests

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	rediscache "go-hexagonal/adapter/repository/redis"
	"go-hexagonal/domain/model"
)

// TestExampleCacheRepoIntegration tests the ExampleCacheRepo with a real Redis connection
// This test requires a running Redis instance
func TestExampleCacheRepoIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Connect to Redis (use testcontainers or local Redis in real CI)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	// Check Redis connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Skipping integration test: Redis not available: %v", err)
	}

	// Create Redis client wrapper
	redisClient := rediscache.NewClientFromRedis(client)

	// Create cache repository
	cacheRepo := rediscache.NewExampleCacheRepo(redisClient)

	// Test data
	testExample := &model.Example{
		Id:        999,
		Name:      "test-integration",
		Alias:     "test-alias",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test Set
	err := cacheRepo.Set(ctx, testExample)
	if err != nil {
		t.Errorf("Failed to set cache: %v", err)
	}

	// Test GetByID
	result, err := cacheRepo.GetByID(ctx, 999)
	if err != nil {
		t.Errorf("Failed to get by ID: %v", err)
	}
	if result == nil || result.Name != "test-integration" {
		t.Errorf("Expected name 'test-integration', got '%v'", result)
	}

	// Test GetByName
	result, err = cacheRepo.GetByName(ctx, "test-integration")
	if err != nil {
		t.Errorf("Failed to get by name: %v", err)
	}
	if result == nil || result.Id != 999 {
		t.Errorf("Expected id 999, got '%v'", result)
	}

	// Test Delete
	err = cacheRepo.Delete(ctx, 999)
	if err != nil {
		t.Errorf("Failed to delete: %v", err)
	}

	// Verify deletion
	_, err = cacheRepo.GetByID(ctx, 999)
	if err == nil {
		t.Errorf("Expected cache miss after deletion")
	}
}

// TestEnhancedCacheIntegration tests the EnhancedCache with a real Redis connection
func TestEnhancedCacheIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Skipping integration test: Redis not available: %v", err)
	}

	// Create Redis client wrapper
	redisClient := rediscache.NewClientFromRedis(client)

	// Create enhanced cache
	cache := rediscache.NewEnhancedCache(redisClient, rediscache.CacheOptions{
		DefaultTTL:        10 * time.Minute,
		EnableKeyTracking: true,
		MaxTrackedKeys:    100,
	})

	// Test data
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	testData := TestData{Name: "integration-test", Value: 42}

	// Test Set
	err := cache.Set(ctx, "test:key:1", testData, 5*time.Minute)
	if err != nil {
		t.Errorf("Failed to set cache: %v", err)
	}

	// Test Get
	var result TestData
	err = cache.Get(ctx, "test:key:1", &result)
	if err != nil {
		t.Errorf("Failed to get cache: %v", err)
	}
	if result.Name != "integration-test" || result.Value != 42 {
		t.Errorf("Expected {Name: integration-test, Value: 42}, got %v", result)
	}

	// Test Delete
	err = cache.Delete(ctx, "test:key:1")
	if err != nil {
		t.Errorf("Failed to delete cache: %v", err)
	}

	// Verify deletion
	err = cache.Get(ctx, "test:key:1", &result)
	if err == nil {
		t.Errorf("Expected cache miss after deletion")
	}
}
