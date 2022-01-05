package mysql

import (
    "context"
    "errors"

    "go-hexagonal/internal/domain.model/dto"
    "go-hexagonal/internal/domain.model/entity"
    "go-hexagonal/internal/domain.model/repo"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

func NewExampleInstance(mysql IMySQL) *Example {
    return &Example{
        IMySQL: mysql,
    }
}

type Example struct {
    IMySQL
}

var _ repo.IExampleRepo = &Example{}

func (e *Example) Create(ctx context.Context, dto dto.CreateExampleReq) (*entity.Example, error) {
    record := &entity.Example{}
    record.Name = dto.Name
    record.Alias = dto.Alias
    err := e.GetDB(ctx).Create(record).Error
    if err != nil {
        return nil, err
    }

    return record, nil
}

func (e *Example) Delete(ctx context.Context, ID int) error {
    if ID == 0 {
        return errors.New("delete fail. need ID")
    }
    err := e.GetDB(ctx).Delete(&entity.Example{}, ID).Error
    return err
}

func (e *Example) Get(ctx context.Context, ID int) (*entity.Example, error) {
    var record *entity.Example
    if ID == 0 {
        return nil, errors.New("get fail. need ID")
    }
    err := e.GetDB(ctx).Find(record, ID).Error
    return record, err
}

func (e *Example) FindByName(ctx context.Context, name string) (*entity.Example, error) {
    var record *entity.Example
    if name == "" {
        return nil, errors.New("FindByName fail. need name")
    }
    err := e.GetDB(ctx).Where("name = ?", name).Last(record).Error
    return record, err
}

func (e *Example) Save(ctx context.Context, example *entity.Example) error {
    return e.GetDB(ctx).Table(example.TableName()).Updates(example.GetChangeMap()).Error
}
