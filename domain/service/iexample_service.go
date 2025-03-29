package service

import (
	"context"

	"go-hexagonal/domain/model"
)

// IExampleService defines the interface for example service
// This allows the application layer to depend on interfaces rather than concrete implementations,
// facilitating testing and adhering to the dependency inversion principle
type IExampleService interface {
	// Create creates a new example with the given name and alias
	// Returns the created example or an error if validation or persistence fails
	Create(ctx context.Context, name string, alias string) (*model.Example, error)

	// Delete deletes an example by ID
	// Returns an error if the example doesn't exist or deletion fails
	Delete(ctx context.Context, id int) error

	// Update updates an example with the given ID, name and alias
	// Returns an error if the example doesn't exist, validation fails, or update fails
	Update(ctx context.Context, id int, name string, alias string) error

	// Get retrieves an example by ID
	// Returns the example or an error if not found
	Get(ctx context.Context, id int) (*model.Example, error)

	// FindByName finds examples by name
	// Returns the example or an error if not found
	FindByName(ctx context.Context, name string) (*model.Example, error)
}
