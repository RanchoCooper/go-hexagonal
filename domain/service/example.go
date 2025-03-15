package service

import (
	"context"

	"go-hexagonal/adapter/repository/mysql/entity"
	"go-hexagonal/domain/event"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// ExampleService handles business logic for Example entity
type ExampleService struct {
	Repository repo.IExampleRepo
	EventBus   event.EventBus
}

// NewExampleService creates a new example service instance
func NewExampleService(ctx context.Context) *ExampleService {
	srv := &ExampleService{Repository: entity.NewExample()}
	return srv
}

// Create creates a new example
func (e *ExampleService) Create(ctx context.Context, model *model.Example) (*model.Example, error) {
	example, err := e.Repository.Create(ctx, nil, model)
	if err != nil {
		return nil, err
	}

	// Publish event if event bus is available
	if e.EventBus != nil {
		evt := event.NewExampleCreatedEvent(example.Id, example.Name, example.Alias)
		e.EventBus.Publish(ctx, evt)
	}

	return example, nil
}

// Delete deletes an example by ID
func (e *ExampleService) Delete(ctx context.Context, id int) error {
	err := e.Repository.Delete(ctx, nil, id)
	if err != nil {
		return err
	}

	// Publish event if event bus is available
	if e.EventBus != nil {
		evt := event.NewExampleDeletedEvent(id)
		e.EventBus.Publish(ctx, evt)
	}

	return nil
}

// Update updates an existing example
func (e *ExampleService) Update(ctx context.Context, model *model.Example) error {
	err := e.Repository.Update(ctx, nil, model)
	if err != nil {
		return err
	}

	// Publish event if event bus is available
	if e.EventBus != nil {
		evt := event.NewExampleUpdatedEvent(model.Id, model.Name, model.Alias)
		e.EventBus.Publish(ctx, evt)
	}

	return nil
}

// Get retrieves an example by ID
func (e *ExampleService) Get(ctx context.Context, id int) (*model.Example, error) {
	example, err := e.Repository.GetByID(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	return example, nil
}
