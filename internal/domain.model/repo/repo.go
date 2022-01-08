package repo

import (
    "context"

    "go-hexagonal/api/http/dto"
    "go-hexagonal/internal/domain.model/entity"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

type IExampleRepo interface {
    Create(ctx context.Context, dto dto.CreateExampleReq) (*entity.Example, error)
    Delete(ctx context.Context, Id int) error
    Save(ctx context.Context, entity *entity.Example) error
    Get(ctx context.Context, Id int) (entity *entity.Example, e error)
    FindByName(ctx context.Context, name string) (*entity.Example, error)
}

type IHealthCheckRepository interface {
    HealthCheck(ctx context.Context) error
}
