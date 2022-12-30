package repo

import (
	"context"

	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/internal/domain/model"
)

/**
 * @author Rancho
 * @date 2022/12/30
 */

type IExampleRepo interface {
	Create(ctx context.Context, tr *repository.Transaction, entity *model.Example) (*model.Example, error)
	Delete(ctx context.Context, tr *repository.Transaction, Id int) error
	Update(ctx context.Context, tr *repository.Transaction, entity *model.Example) error
	GetByID(ctx context.Context, tr *repository.Transaction, Id int) (*model.Example, error)
	FindByName(ctx context.Context, tr *repository.Transaction, name string) (*model.Example, error)
}

type IExampleCacheRepo interface {
	HealthCheck(ctx context.Context) error
}
