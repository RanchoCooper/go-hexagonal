package service

import (
	"context"

	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/internal/domain/entity"
	"go-hexagonal/internal/domain/repo"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

type ExampleService struct {
	Repository repo.IExampleRepo
}

func NewExampleService(ctx context.Context) *ExampleService {
	srv := &ExampleService{Repository: repository.Example}
	log.Logger.Info("example service init successfully")
	return srv
}

func (e *ExampleService) Create(ctx context.Context, example *entity.Example) (*entity.Example, error) {
	example, err := e.Repository.Create(ctx, nil, example)
	if err != nil {
		return nil, err
	}
	return example, nil
}

func (e *ExampleService) Delete(ctx context.Context, id int) error {
	err := e.Repository.Delete(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}

func (e *ExampleService) Update(ctx context.Context, example *entity.Example) error {
	err := e.Repository.Update(ctx, nil, example)
	if err != nil {
		return err
	}
	return nil
}

func (e *ExampleService) Get(ctx context.Context, id int) (*entity.Example, error) {
	example, err := e.Repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return example, nil
}
