package http

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"go-hexagonal/adapter/dependency"
	"go-hexagonal/adapter/repository"
	"go-hexagonal/adapter/repository/mysql"
	"go-hexagonal/adapter/repository/mysql/entity"
	"go-hexagonal/adapter/repository/redis"
	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

var ctx = context.Background()

func TestMain(m *testing.M) {
	// Parse command line arguments, support -short flag
	flag.Parse()

	// Initialize configuration and logging
	config.Init("../../config", "config")
	log.Init()

	// Skip integration tests in short mode
	if testing.Short() {
		fmt.Println("Skipping integration tests in short mode")
		os.Exit(0)
		return
	}

	// Use test containers
	t := &testing.T{}
	mysqlConfig := mysql.SetupMySQLContainer(t)
	redisConfig := redis.SetupRedisContainer(t)

	// Set global config to use test containers
	config.GlobalConfig.MySQL = mysqlConfig
	config.GlobalConfig.Redis = redisConfig

	// Initialize repositories using dependency injection
	clients, err := dependency.InitializeRepositories(
		dependency.WithMySQL(),
		dependency.WithRedis(),
	)
	if err != nil {
		log.SugaredLogger.Fatalf("Failed to initialize repositories: %v", err)
	}
	repository.Clients = clients
	_ = repository.Clients.MySQL.GetDB(ctx).AutoMigrate(&entity.Example{})

	// Initialize services using dependency injection
	svcs, err := dependency.InitializeServices(ctx, dependency.WithExampleService())
	if err != nil {
		log.SugaredLogger.Fatalf("Failed to initialize services: %v", err)
	}

	// Register services for API handlers
	RegisterServices(svcs)

	// Run tests
	os.Exit(m.Run())
}
