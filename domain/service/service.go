package service

import (
	"context"
	"sync"

	"go-hexagonal/domain/event"
)

var (
	once       sync.Once
	ExampleSvc *ExampleService
	EventBus   event.EventBus
)

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
func Init(ctx context.Context) {
	once.Do(func() {
		// Initialize event bus
		EventBus = event.NewInMemoryEventBus()

		// Register event handlers
		loggingHandler := event.NewLoggingEventHandler()
		exampleHandler := event.NewExampleEventHandler()
		EventBus.Subscribe(loggingHandler)
		EventBus.Subscribe(exampleHandler)

		// Initialize services
		ExampleSvc = NewExampleService(ctx)
		ExampleSvc.EventBus = EventBus
	})
}
