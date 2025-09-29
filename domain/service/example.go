package service

import (
	"context"

	"go-hexagonal/domain/event"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
	"go-hexagonal/util/error_handler"
	"go-hexagonal/util/log"
)

// Ensure ExampleService implements IExampleService
var _ IExampleService = (*ExampleService)(nil)

// ExampleService handles business logic for Example entity
type ExampleService struct {
	Repository repo.IExampleRepo
	CacheRepo  repo.IExampleCacheRepo
	EventBus   event.EventBus
}

// NewExampleService creates a new example service instance
func NewExampleService(repository repo.IExampleRepo, cacheRepo repo.IExampleCacheRepo) *ExampleService {
	return &ExampleService{
		Repository: repository,
		CacheRepo:  cacheRepo,
	}
}

// Create creates a new example
func (s *ExampleService) Create(ctx context.Context, name string, alias string) (*model.Example, error) {
	// Create a new example entity
	example, err := model.NewExample(name, alias)
	if err != nil {
		return nil, error_handler.HandleAndConvertError(ctx, err, "create example entity", "validation")
	}

	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	// Persist the entity
	createdExample, err := s.Repository.Create(ctx, tr, example)
	if err != nil {
		return nil, error_handler.HandleAndWrapError(ctx, err, "persist example", "failed to create example")
	}

	// Update cache if available
	if s.CacheRepo != nil {
		if err := s.CacheRepo.Set(ctx, createdExample); err != nil {
			log.SugaredLogger.Warnf("Failed to update cache: %v", err)
		}
	}

	// Publish domain events if event bus is available
	if s.EventBus != nil {
		domainEvents := createdExample.Events()
		for _, evt := range domainEvents {
			if domainEvt, ok := evt.(model.ExampleCreatedEvent); ok {
				// Map domain event to integration event
				integrationEvent := event.NewExampleCreatedEvent(
					domainEvt.ExampleID,
					domainEvt.Name,
					domainEvt.Alias,
				)

				if err := s.EventBus.Publish(ctx, integrationEvent); err != nil {
					log.SugaredLogger.Warnf("Failed to publish event: %v", err)
				}
			}
		}
	}

	return createdExample, nil
}

// Delete deletes an example by ID
func (s *ExampleService) Delete(ctx context.Context, id int) error {
	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	// Get the example to be deleted
	example, err := s.Repository.GetByID(ctx, tr, id)
	if err != nil {
		return error_handler.HandleAndWrapError(ctx, err, "get example for deletion", "example not found")
	}

	// Mark example as deleted (generates domain event)
	example.MarkDeleted()

	// Delete from repository
	if err := s.Repository.Delete(ctx, tr, id); err != nil {
		return error_handler.HandleAndWrapError(ctx, err, "delete example", "failed to delete example")
	}

	// Invalidate cache if available
	if s.CacheRepo != nil {
		if err := s.CacheRepo.Delete(ctx, id); err != nil {
			log.SugaredLogger.Warnf("Failed to invalidate cache: %v", err)
		}
	}

	// Publish domain events if event bus is available
	if s.EventBus != nil {
		domainEvents := example.Events()
		for _, evt := range domainEvents {
			if domainEvt, ok := evt.(model.ExampleDeletedEvent); ok {
				// Map domain event to integration event
				integrationEvent := event.NewExampleDeletedEvent(domainEvt.ExampleID)

				if err := s.EventBus.Publish(ctx, integrationEvent); err != nil {
					log.SugaredLogger.Warnf("Failed to publish event: %v", err)
				}
			}
		}
	}

	return nil
}

// Update updates an existing example
func (s *ExampleService) Update(ctx context.Context, id int, name string, alias string) error {
	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	// Get the example to be updated
	example, err := s.Repository.GetByID(ctx, tr, id)
	if err != nil {
		return error_handler.HandleAndWrapError(ctx, err, "get example for update", "example not found")
	}

	// Update the entity (generates domain event)
	if err := example.Update(name, alias); err != nil {
		return error_handler.HandleAndConvertError(ctx, err, "update example entity", "invalid update data")
	}

	// Persist the changes
	if err := s.Repository.Update(ctx, tr, example); err != nil {
		return error_handler.HandleAndWrapError(ctx, err, "persist example update", "failed to update example")
	}

	// Update cache if available
	if s.CacheRepo != nil {
		if err := s.CacheRepo.Set(ctx, example); err != nil {
			log.SugaredLogger.Warnf("Failed to update cache: %v", err)
		}
	}

	// Publish domain events if event bus is available
	if s.EventBus != nil {
		domainEvents := example.Events()
		for _, evt := range domainEvents {
			if domainEvt, ok := evt.(model.ExampleUpdatedEvent); ok {
				// Map domain event to integration event
				integrationEvent := event.NewExampleUpdatedEvent(
					domainEvt.ExampleID,
					domainEvt.Name,
					domainEvt.Alias,
				)

				if err := s.EventBus.Publish(ctx, integrationEvent); err != nil {
					log.SugaredLogger.Warnf("Failed to publish event: %v", err)
				}
			}
		}
	}

	return nil
}

// Get retrieves an example by ID
func (s *ExampleService) Get(ctx context.Context, id int) (*model.Example, error) {
	// Try to get from cache first
	if s.CacheRepo != nil {
		example, err := s.CacheRepo.GetByID(ctx, id)
		if err == nil {
			return example, nil
		}
		log.SugaredLogger.Debugf("Cache miss for example ID %d: %v", id, err)
	}

	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	// Get from repository
	example, err := s.Repository.GetByID(ctx, tr, id)
	if err != nil {
		return nil, error_handler.HandleError(ctx, err, "get example by ID")
	}

	// Update cache if available
	if s.CacheRepo != nil {
		if err := s.CacheRepo.Set(ctx, example); err != nil {
			log.SugaredLogger.Warnf("Failed to update cache: %v", err)
		}
	}

	return example, nil
}

// FindByName retrieves an example by name
func (s *ExampleService) FindByName(ctx context.Context, name string) (*model.Example, error) {
	// Try to get from cache first
	if s.CacheRepo != nil {
		example, err := s.CacheRepo.GetByName(ctx, name)
		if err == nil {
			return example, nil
		}
		log.SugaredLogger.Debugf("Cache miss for example name %s: %v", name, err)
	}

	// Create a no-operation transaction
	tr := repo.NewNoopTransaction(s.Repository)

	// Get from repository
	example, err := s.Repository.FindByName(ctx, tr, name)
	if err != nil {
		return nil, error_handler.HandleAndWrapError(ctx, err, "find example by name", "failed to find example")
	}

	// Update cache if available
	if s.CacheRepo != nil {
		if err := s.CacheRepo.Set(ctx, example); err != nil {
			log.SugaredLogger.Warnf("Failed to update cache: %v", err)
		}
	}

	return example, nil
}
