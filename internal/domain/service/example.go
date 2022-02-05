package service

import (
    "context"

    "go-hexagonal/internal/adapter/repository"
    "go-hexagonal/internal/domain/entity"
    "go-hexagonal/internal/domain/repo"
    "go-hexagonal/util/logger"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

type ExampleSvc struct {
    Repository repo.IExampleRepo
}

func NewExampleService(ctx context.Context) *ExampleSvc {
    srv := &ExampleSvc{Repository: repository.Example}
    logger.Log.Info(ctx, "example service init successfully")
    return srv
}

func (e *ExampleSvc) Create(ctx context.Context, example *entity.Example) (*entity.Example, error) {
    example, err := e.Repository.Create(ctx, nil, example)
    if err != nil {
        return nil, err
    }
    return example, nil
}

func (e *ExampleSvc) Delete(ctx context.Context, id int) error {
    err := e.Repository.Delete(ctx, nil, id)
    if err != nil {
        return err
    }
    return nil
}

func (e *ExampleSvc) Update(ctx context.Context, example *entity.Example) error {
    err := e.Repository.Update(ctx, nil, example)
    if err != nil {
        return err
    }
    return nil
}

func (e *ExampleSvc) Get(ctx context.Context, id int) (*entity.Example, error) {
    example, err := e.Repository.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    return example, nil
}
