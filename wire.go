// Package dependency provides dependency injection configuration
//go:build wireinject
// +build wireinject

package dependency

import (
	"context"

	"github.com/google/wire"

	"go-hexagonal/adapter/repository/mysql/entity"
	"go-hexagonal/domain/event"
	"go-hexagonal/domain/service"
)

// InitializeServices initializes all services
func InitializeServices(ctx context.Context) (*service.Services, error) {
	wire.Build(
		setupRepositories,
		setupEventBus,
		setupServices,
	)
	return nil, nil
}

// setupRepositories creates and configures repositories
func setupRepositories() *entity.EntityExample {
	return entity.NewExample()
}

// setupEventBus creates and configures the event bus with handlers
func setupEventBus() *event.InMemoryEventBus {
	eventBus := event.NewInMemoryEventBus()
	loggingHandler := event.NewLoggingEventHandler()
	exampleHandler := event.NewExampleEventHandler()
	eventBus.Subscribe(loggingHandler)
	eventBus.Subscribe(exampleHandler)
	return eventBus
}

// setupServices creates and configures services
func setupServices(ctx context.Context, repo *entity.EntityExample, eventBus *event.InMemoryEventBus) *service.Services {
	exampleService := service.NewExampleService(ctx)
	exampleService.Repository = repo
	exampleService.EventBus = eventBus
	return service.NewServices(exampleService, eventBus)
}
