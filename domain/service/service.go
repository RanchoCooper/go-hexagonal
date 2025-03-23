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
	Converter      Converter
}

// NewServices creates a services collection
func NewServices(exampleService *ExampleService, eventBus event.EventBus) *Services {
	return &Services{
		ExampleService: exampleService,
		EventBus:       eventBus,
	}
}

// WithConverter adds a converter to the services
func (s *Services) WithConverter(converter Converter) *Services {
	s.Converter = converter
	return s
}
