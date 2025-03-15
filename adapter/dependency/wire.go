// Package dependency provides dependency injection configuration
//go:build wireinject
// +build wireinject

package dependency

import (
	"context"

	"github.com/google/wire"

	"go-hexagonal/adapter/repository/mysql/entity"
	"go-hexagonal/domain/event"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// InitializeServices initializes all services
func InitializeServices(ctx context.Context) (*service.Services, error) {
	wire.Build(
		// Repository dependencies
		entity.NewExample,
		wire.Bind(new(repo.IExampleRepo), new(*entity.EntityExample)),

		// Event bus dependencies
		provideEventBus,
		wire.Bind(new(event.EventBus), new(*event.InMemoryEventBus)),

		// Service dependencies
		provideExampleService,
		provideServices,
	)
	return nil, nil
}

// provideEventBus creates and configures the event bus
func provideEventBus() *event.InMemoryEventBus {
	eventBus := event.NewInMemoryEventBus()

	// Register event handlers
	loggingHandler := event.NewLoggingEventHandler()
	exampleHandler := event.NewExampleEventHandler()
	eventBus.Subscribe(loggingHandler)
	eventBus.Subscribe(exampleHandler)

	return eventBus
}

// provideExampleService creates and configures the example service
func provideExampleService(ctx context.Context, repo repo.IExampleRepo, eventBus event.EventBus) *service.ExampleService {
	exampleService := service.NewExampleService(ctx)
	exampleService.Repository = repo
	exampleService.EventBus = eventBus
	return exampleService
}

// provideServices creates the services container
func provideServices(exampleService *service.ExampleService, eventBus event.EventBus) *service.Services {
	return service.NewServices(exampleService, eventBus)
}
