package repo

import (
	"context"

	"gorm.io/gorm"

	"go-hexagonal/internal/domain/model"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

type IExampleRepo interface {
	Create(ctx context.Context, tx *gorm.DB, entity *model.Example) (*model.Example, error)
	Delete(ctx context.Context, tx *gorm.DB, Id int) error
	Update(ctx context.Context, tx *gorm.DB, entity *model.Example) error
	GetByID(ctx context.Context, Id int) (*model.Example, error)
	FindByName(ctx context.Context, name string) (*model.Example, error)
}

type IExampleCacheRepo interface {
	HealthCheck(ctx context.Context) error
}
