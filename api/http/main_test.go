package http

import (
	"context"
	"testing"

	"go-hexagonal/adapter/dependency"
	"go-hexagonal/adapter/repository"
	"go-hexagonal/adapter/repository/mysql/entity"
	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init("../../config", "config")
	log.Init()

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

	m.Run()
}
