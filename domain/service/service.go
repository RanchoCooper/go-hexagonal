package service

import (
	"go-hexagonal/domain/event"
	"go-hexagonal/domain/repo"
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
