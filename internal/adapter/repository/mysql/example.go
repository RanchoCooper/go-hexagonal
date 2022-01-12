package mysql

import (
    "context"

    "github.com/jinzhu/copier"
    "github.com/pkg/errors"
    "gorm.io/gorm"

    "go-hexagonal/api/http/dto"
    "go-hexagonal/internal/domain/entity"
    "go-hexagonal/internal/domain/repo"
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

func (e *Example) Create(ctx context.Context, tx *gorm.DB, dto dto.CreateExampleReq) (result *entity.Example, err error) {
    if tx == nil {
        tx = e.GetDB(ctx).Begin()
        defer func() {
            if r := recover(); r != nil {
                tx.Rollback()
                return
            }
            if err != nil {
                tx.Rollback()
                return
            }
            err = errors.WithStack(tx.Commit().Error)
        }()
    }
    record := &entity.Example{}
    err = copier.Copy(record, dto)
    if err != nil {
        return nil, err
    }
    err = tx.Table(record.TableName()).Create(record).Error
    if err != nil {
        return nil, err
    }

    return record, nil
}

func (e *Example) Delete(ctx context.Context, tx *gorm.DB, id int) (err error) {
    if tx == nil {
        tx = e.GetDB(ctx).Begin()
        defer func() {
            if r := recover(); r != nil {
                tx.Rollback()
                return
            }
            if err != nil {
                tx.Rollback()
                return
            }
            err = errors.WithStack(tx.Commit().Error)
        }()
    }
    if id == 0 {
        return errors.New("delete fail. need Id")
    }
    example := &entity.Example{}
    err = tx.Table(example.TableName()).Delete(example, id).Error
    // hard delete with .Unscoped()
    // err := e.GetDB(ctx).Table(example.TableName()).Unscoped().Delete(example, Id).Error
    return err
}

func (e *Example) Save(ctx context.Context, tx *gorm.DB, example *entity.Example) (err error) {
    if tx == nil {
        tx = e.GetDB(ctx).Begin()
        defer func() {
            if r := recover(); r != nil {
                tx.Rollback()
                return
            }
            if err != nil {
                tx.Rollback()
                return
            }
            err = errors.WithStack(tx.Commit().Error)
        }()
    }
    return tx.Table(example.TableName()).Where("id = ? AND deleted_at IS NULL", example.Id).Updates(example.ChangeMap).Error
}

func (e *Example) Get(ctx context.Context, id int) (*entity.Example, error) {
    record := &entity.Example{}
    if id == 0 {
        return nil, errors.New("get fail. need Id")
    }
    err := e.GetDB(ctx).Table(record.TableName()).Find(record, id).Error
    return record, err
}

func (e *Example) FindByName(ctx context.Context, name string) (*entity.Example, error) {
    record := &entity.Example{}
    if name == "" {
        return nil, errors.New("FindByName fail. need name")
    }
    err := e.GetDB(ctx).Table(record.TableName()).Where("name = ?", name).Last(record).Error
    return record, err
}
