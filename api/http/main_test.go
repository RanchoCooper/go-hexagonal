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

	// Initialize repositories
	repository.Init(repository.WithMySQL(), repository.WithRedis())
	_ = repository.Clients.MySQL.GetDB(ctx).AutoMigrate(&entity.EntityExample{})

	// Initialize services using dependency injection
	svcs, err := dependency.InitializeServices(ctx)
	if err != nil {
		log.SugaredLogger.Fatalf("Failed to initialize services: %v", err)
	}

	// Register services for API handlers
	RegisterServices(svcs)

	m.Run()
}
