package service

import (
	"context"
	"fmt"

	"go-hexagonal/domain/event"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
	"go-hexagonal/util/log"
)

// Ensure ExampleService implements IExampleService
var _ IExampleService = (*ExampleService)(nil)

// ExampleService handles business logic for Example entity
type ExampleService struct {
	Repository repo.IExampleRepo
	EventBus   event.EventBus
}

// NewExampleService creates a new example service instance
func NewExampleService(repository repo.IExampleRepo) *ExampleService {
	return &ExampleService{
		Repository: repository,
	}
}

// Create creates a new example
func (s *ExampleService) Create(ctx context.Context, example *model.Example) (*model.Example, error) {
	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	createdExample, err := s.Repository.Create(ctx, tr, example)
	if err != nil {
		log.SugaredLogger.Errorf("Failed to create example: %v", err)
		return nil, fmt.Errorf("failed to create example: %w", err)
	}

	// Publish event if event bus is available
	if s.EventBus != nil {
		evt := event.NewExampleCreatedEvent(createdExample.Id, createdExample.Name, createdExample.Alias)
		if err := s.EventBus.Publish(ctx, evt); err != nil {
			log.SugaredLogger.Warnf("Failed to publish event: %v", err)
			return createdExample, fmt.Errorf("failed to publish example created event: %w", err)
		}
	}

	return createdExample, nil
}

// Delete deletes an example by ID
func (s *ExampleService) Delete(ctx context.Context, id int) error {
	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	if err := s.Repository.Delete(ctx, tr, id); err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}

	// Publish event if event bus is available
	if s.EventBus != nil {
		evt := event.NewExampleDeletedEvent(id)
		if err := s.EventBus.Publish(ctx, evt); err != nil {
			return fmt.Errorf("failed to publish example deleted event: %w", err)
		}
	}

	return nil
}

// Update updates an existing example
func (s *ExampleService) Update(ctx context.Context, example *model.Example) error {
	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	if err := s.Repository.Update(ctx, tr, example); err != nil {
		return fmt.Errorf("failed to update example: %w", err)
	}

	// Publish event if event bus is available
	if s.EventBus != nil {
		evt := event.NewExampleUpdatedEvent(example.Id, example.Name, example.Alias)
		if err := s.EventBus.Publish(ctx, evt); err != nil {
			return fmt.Errorf("failed to publish example updated event: %w", err)
		}
	}

	return nil
}

// Get retrieves an example by ID
func (s *ExampleService) Get(ctx context.Context, id int) (*model.Example, error) {
	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	example, err := s.Repository.GetByID(ctx, tr, id)
	if err != nil {
		return nil, err
	}

	return example, nil
}

// FindByName retrieves an example by name
func (s *ExampleService) FindByName(ctx context.Context, name string) (*model.Example, error) {
	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	example, err := s.Repository.FindByName(ctx, tr, name)
	if err != nil {
		return nil, fmt.Errorf("failed to find example: %w", err)
	}

	return example, nil
}
