package repo

import (
	"context"

	"go-hexagonal/domain/model"
)

// IExampleRepo defines the interface for example repository
type IExampleRepo interface {
	Create(ctx context.Context, tr Transaction, example *model.Example) (*model.Example, error)
	Delete(ctx context.Context, tr Transaction, id int) error
	Update(ctx context.Context, tr Transaction, entity *model.Example) error
	GetByID(ctx context.Context, tr Transaction, Id int) (*model.Example, error)
	FindByName(ctx context.Context, tr Transaction, name string) (*model.Example, error)
}

// IExampleCacheRepo defines the interface for example cache repository
type IExampleCacheRepo interface {
	HealthCheck(ctx context.Context) error
}
