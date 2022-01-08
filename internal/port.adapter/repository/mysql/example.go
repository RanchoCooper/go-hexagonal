package mysql

import (
    "context"
    "errors"

    "github.com/jinzhu/copier"

    "go-hexagonal/api/http/dto"
    "go-hexagonal/internal/domain.model/entity"
    "go-hexagonal/internal/domain.model/repo"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

func NewExample(mysql IMySQL) *Example {
    return &Example{IMySQL: mysql}
}

type Example struct {
    IMySQL
}

var _ repo.IExampleRepo = &Example{}

func (e *Example) Create(ctx context.Context, dto dto.CreateExampleReq) (*entity.Example, error) {
    record := &entity.Example{}
    record.Name = dto.Name
    record.Alias = dto.Alias
    _ = copier.Copy(record, dto)
    err := e.GetDB(ctx).Table(record.TableName()).Create(record).Error
    if err != nil {
        return nil, err
    }

    return record, nil
}

func (e *Example) Delete(ctx context.Context, ID int) error {
    if ID == 0 {
        return errors.New("delete fail. need Id")
    }
    example := &entity.Example{}
    err := e.GetDB(ctx).Table(example.TableName()).Delete(example, ID).Error
    // hard delete with .Unscoped()
    // err := e.GetDB(ctx).Table(example.TableName()).Unscoped().Delete(example, Id).Error
    return err
}

func (e *Example) Get(ctx context.Context, ID int) (*entity.Example, error) {
    var record *entity.Example
    if ID == 0 {
        return nil, errors.New("get fail. need Id")
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
    return e.GetDB(ctx).Table(example.TableName()).Where("id = ? AND deleted_at IS NULL", example.Id).Updates(example.ChangeMap).Error
}
