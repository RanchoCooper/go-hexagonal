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

// RepositorySet defines repository dependencies
var RepositorySet = wire.NewSet(
	entity.NewExample,
	wire.Bind(new(repo.IExampleRepo), new(*entity.EntityExample)),
)

// EventSet defines event-related dependencies
var EventSet = wire.NewSet(
	event.NewInMemoryEventBus,
	wire.Bind(new(event.EventBus), new(*event.InMemoryEventBus)),
	event.NewLoggingEventHandler,
	event.NewExampleEventHandler,
)

// ServiceSet defines service dependencies
var ServiceSet = wire.NewSet(
	service.NewExampleService,
)

// ApplicationSet defines application layer dependencies
var ApplicationSet = wire.NewSet(
// Application layer dependencies will be added here
)

// AllSets combines all dependency sets
var AllSets = wire.NewSet(
	RepositorySet,
	EventSet,
	ServiceSet,
	ApplicationSet,
)

// InitializeServices initializes all services
func InitializeServices(ctx context.Context) (*service.Services, error) {
	wire.Build(
		AllSets,
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
