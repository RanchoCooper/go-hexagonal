package service

import (
	"context"
	"sync"

	"go-hexagonal/domain/event"
	"go-hexagonal/domain/repo"
)

var (
	once       sync.Once
	ExampleSvc *ExampleService
	EventBus   event.EventBus
)

// ExampleRepoFactory defines the interface for example repository factory
type ExampleRepoFactory interface {
	CreateExampleRepo() repo.IExampleRepo
}

// Services contains all service instances
type Services struct {
	ExampleService *ExampleService
	EventBus       event.EventBus
}

// NewServices creates a services collection
func NewServices(exampleService *ExampleService, eventBus event.EventBus) *Services {
	return &Services{
		ExampleService: exampleService,
		EventBus:       eventBus,
	}
}

// Init initializes services (legacy method for backward compatibility)
// Note: This method is deprecated, services should be initialized through dependency injection
func Init(ctx context.Context) {
	// This method is deprecated, new code should use dependency injection
	once.Do(func() {
		// Initialize event bus
		EventBus = event.NewInMemoryEventBus()

		// Register event handlers
		loggingHandler := event.NewLoggingEventHandler()
		exampleHandler := event.NewExampleEventHandler()
		EventBus.Subscribe(loggingHandler)
		EventBus.Subscribe(exampleHandler)

		// Service instances will be injected by the infrastructure layer
		ExampleSvc = NewExampleService(nil)
		ExampleSvc.EventBus = EventBus
	})
}
