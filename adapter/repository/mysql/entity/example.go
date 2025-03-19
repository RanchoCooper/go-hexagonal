package entity

import (
	"context"
	"time"

	"gorm.io/gorm"

	"go-hexagonal/adapter/repository"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
	"go-hexagonal/util/log"
)

// NewExample creates a new EntityExample instance
func NewExample() *EntityExample {
	return &EntityExample{}
}

// EntityExample represents the database entity for Example
type EntityExample struct {
	ID        int                    `json:"id" gorm:"column:id;primaryKey;autoIncrement" structs:",omitempty,underline"`
	Name      string                 `json:"name" gorm:"column:name;type:varchar(255);not null" structs:",omitempty,underline"`
	Alias     string                 `json:"alias" gorm:"column:alias;type:varchar(255)" structs:",omitempty,underline"`
	CreatedAt time.Time              `json:"created_at" gorm:"column:created_at;autoCreateTime" structs:",omitempty,underline"`
	UpdatedAt time.Time              `json:"updated_at" gorm:"column:updated_at;autoUpdateTime" structs:",omitempty,underline"`
	DeletedAt gorm.DeletedAt         `json:"deleted_at" gorm:"column:deleted_at" structs:",omitempty,underline"`
	ChangeMap map[string]interface{} `json:"-" gorm:"-" structs:"-"`
}

// TableName returns the table name for the EntityExample
func (e EntityExample) TableName() string {
	return "example"
}

// Ensure EntityExample implements IExampleRepo
var _ repo.IExampleRepo = &EntityExample{}

// ToModel converts the EntityExample to a model.Example
func (e EntityExample) ToModel() *model.Example {
	return &model.Example{
		Id:        e.ID,
		Name:      e.Name,
		Alias:     e.Alias,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// FromModel converts a model.Example to an EntityExample
func (e *EntityExample) FromModel(m *model.Example) {
	e.ID = m.Id
	e.Name = m.Name
	e.Alias = m.Alias
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt
}

// Create creates a new example record
func (e *EntityExample) Create(ctx context.Context, tr repo.Transaction, example *model.Example) (*model.Example, error) {
	entity := &EntityExample{}
	entity.FromModel(example)

	db := repository.Clients.MySQL.GetDB(ctx)

	if err := db.Create(entity).Error; err != nil {
		log.SugaredLogger.Errorf("Failed to create example: %v", err)
		return nil, err
	}

	return entity.ToModel(), nil
}

// Delete deletes an example record by ID
func (e *EntityExample) Delete(ctx context.Context, tr repo.Transaction, id int) error {
	db := repository.Clients.MySQL.GetDB(ctx)
	result := db.Delete(&EntityExample{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Update updates an example record
func (e *EntityExample) Update(ctx context.Context, tr repo.Transaction, example *model.Example) error {
	entity := &EntityExample{}
	entity.FromModel(example)

	db := repository.Clients.MySQL.GetDB(ctx)
	result := db.Model(&EntityExample{}).Where("id = ?", entity.ID).Updates(entity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetByID retrieves an example record by ID
func (e *EntityExample) GetByID(ctx context.Context, tr repo.Transaction, id int) (*model.Example, error) {
	entity := &EntityExample{}

	db := repository.Clients.MySQL.GetDB(ctx)
	if err := db.Model(&EntityExample{}).Where("id = ?", id).First(entity).Error; err != nil {
		return nil, err
	}

	return entity.ToModel(), nil
}

// FindByName retrieves an example record by name
func (e *EntityExample) FindByName(ctx context.Context, tr repo.Transaction, name string) (*model.Example, error) {
	entity := &EntityExample{}

	db := repository.Clients.MySQL.GetDB(ctx)
	if err := db.Where("name = ?", name).First(entity).Error; err != nil {
		return nil, err
	}

	return entity.ToModel(), nil
}
