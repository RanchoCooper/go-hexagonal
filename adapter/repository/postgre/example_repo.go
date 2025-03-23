package postgre

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// ExampleRepo implements the example repository for PostgreSQL
type ExampleRepo struct {
	client *PostgreSQLClient
}

// NewExampleRepo creates a new PostgreSQL example repository
func NewExampleRepo(client *PostgreSQLClient) repo.IExampleRepo {
	return &ExampleRepo{
		client: client,
	}
}

// Create creates a new example in the database
func (r *ExampleRepo) Create(ctx context.Context, tr repo.Transaction, example *model.Example) (*model.Example, error) {
	// Set timestamps
	now := time.Now()
	example.CreatedAt = now
	example.UpdatedAt = now

	// Get DB connection (from transaction or direct client)
	db := r.getDB(ctx, tr)

	// Create record
	if err := db.Create(example).Error; err != nil {
		return nil, err
	}

	return example, nil
}

// Update updates an existing example
func (r *ExampleRepo) Update(ctx context.Context, tr repo.Transaction, example *model.Example) error {
	// Set update timestamp
	example.UpdatedAt = time.Now()

	// Get DB connection (from transaction or direct client)
	db := r.getDB(ctx, tr)

	// Update record - note the use of updates map to handle zero values properly in PostgreSQL
	updates := map[string]interface{}{
		"name":       example.Name,
		"alias":      example.Alias,
		"updated_at": example.UpdatedAt,
	}

	result := db.Model(&model.Example{}).Where("id = ?", example.Id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	// Check if record exists
	if result.RowsAffected == 0 {
		return repo.ErrNotFound
	}

	return nil
}

// Delete deletes an example by ID
func (r *ExampleRepo) Delete(ctx context.Context, tr repo.Transaction, id int) error {
	// Get DB connection (from transaction or direct client)
	db := r.getDB(ctx, tr)

	// Delete record
	result := db.Delete(&model.Example{}, id)
	if result.Error != nil {
		return result.Error
	}

	// Check if record exists
	if result.RowsAffected == 0 {
		return repo.ErrNotFound
	}

	return nil
}

// GetByID retrieves an example by ID
func (r *ExampleRepo) GetByID(ctx context.Context, tr repo.Transaction, id int) (*model.Example, error) {
	// Get DB connection (from transaction or direct client)
	db := r.getDB(ctx, tr)

	// Find record
	var example model.Example
	if err := db.Where("id = ?", id).First(&example).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}

	return &example, nil
}

// FindByName retrieves an example by name
func (r *ExampleRepo) FindByName(ctx context.Context, tr repo.Transaction, name string) (*model.Example, error) {
	// Get DB connection (from transaction or direct client)
	db := r.getDB(ctx, tr)

	// Find record
	var example model.Example
	if err := db.Where("name = ?", name).First(&example).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repo.ErrNotFound
		}
		return nil, err
	}

	return &example, nil
}

// getDB returns the appropriate database connection based on transaction
func (r *ExampleRepo) getDB(ctx context.Context, tr repo.Transaction) *gorm.DB {
	if tr != nil {
		if db, ok := tr.Conn(ctx).(*gorm.DB); ok {
			return db
		}
	}
	return r.client.GetDB(ctx)
}
