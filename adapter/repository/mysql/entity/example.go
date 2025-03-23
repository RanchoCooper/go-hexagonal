package entity

import (
	"context"
	"time"

	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// Example represents the MySQL implementation of IExampleRepo
type Example struct {
	// Database connection or any dependencies could be added here
}

// NewExample creates a new Example repository
func NewExample() *Example {
	return &Example{}
}

// Create implements IExampleRepo.Create
func (e *Example) Create(ctx context.Context, tr repo.Transaction, example *model.Example) (*model.Example, error) {
	// Implement actual database logic for creation
	return example, nil
}

// GetByID implements IExampleRepo.GetByID
func (e *Example) GetByID(ctx context.Context, tr repo.Transaction, id int) (*model.Example, error) {
	// Implement actual database logic for fetching
	return &model.Example{
		Id:        id,
		Name:      "Example from MySQL",
		Alias:     "MySQL Demo",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Update implements IExampleRepo.Update
func (e *Example) Update(ctx context.Context, tr repo.Transaction, example *model.Example) error {
	// Implement actual database logic for updating
	return nil
}

// Delete implements IExampleRepo.Delete
func (e *Example) Delete(ctx context.Context, tr repo.Transaction, id int) error {
	// Implement actual database logic for deletion
	return nil
}

// FindByName implements IExampleRepo.FindByName
func (e *Example) FindByName(ctx context.Context, tr repo.Transaction, name string) (*model.Example, error) {
	// Implement actual database logic for finding by name
	return &model.Example{
		Id:        1,
		Name:      name,
		Alias:     "MySQL Mock",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// WithTransaction implements IExampleRepo.WithTransaction
func (e *Example) WithTransaction(ctx context.Context, tx repo.Transaction) repo.IExampleRepo {
	// Return the same repository for now, as it's a mock
	// In a real implementation, this would create a new repository that uses the transaction
	return e
}
