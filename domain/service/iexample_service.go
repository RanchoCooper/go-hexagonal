package service

import (
	"context"

	"go-hexagonal/domain/model"
)

// IExampleService defines the interface for example service
// This allows the application layer to depend on interfaces rather than concrete implementations,
// facilitating testing and adhering to the dependency inversion principle
type IExampleService interface {
	// Create creates a new example
	Create(ctx context.Context, example *model.Example) (*model.Example, error)

	// Delete deletes an example by ID
	Delete(ctx context.Context, id int) error

	// Update updates an example
	Update(ctx context.Context, example *model.Example) error

	// Get retrieves an example by ID
	Get(ctx context.Context, id int) (*model.Example, error)

	// FindByName finds examples by name
	FindByName(ctx context.Context, name string) (*model.Example, error)
}
