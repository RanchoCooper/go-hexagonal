package service

import (
	"context"

	"go-hexagonal/internal/adapter/repository/mysql/entity"
	"go-hexagonal/internal/domain/model"
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
	srv := &ExampleService{Repository: entity.NewExample()}
	log.Logger.Info("example service init successfully")
	return srv
}

func (e *ExampleService) Create(ctx context.Context, model *model.Example) (*model.Example, error) {
	example, err := e.Repository.Create(ctx, nil, model)
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

func (e *ExampleService) Update(ctx context.Context, model *model.Example) error {
	err := e.Repository.Update(ctx, nil, model)
	if err != nil {
		return err
	}
	return nil
}

func (e *ExampleService) Get(ctx context.Context, id int) (*model.Example, error) {
	example, err := e.Repository.GetByID(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	return example, nil
}
