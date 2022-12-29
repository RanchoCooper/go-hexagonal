package mysql

import (
	"context"
	"time"

	"github.com/RanchoCooper/structs"
	"github.com/pkg/errors"
	"gorm.io/gorm"

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

func (e *Example) Create(ctx context.Context, tx *gorm.DB, example *entity.Example) (result *entity.Example, err error) {
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
	err = tx.Create(example).Error
	if err != nil {
		return nil, err
	}

	return example, nil
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
	err = tx.Delete(example, id).Error
	// hard delete with .Unscoped()
	// err := e.GetDB(ctx).Table(example.TableName()).Unscoped().Delete(example, Id).Error
	return err
}

func (e *Example) Update(ctx context.Context, tx *gorm.DB, example *entity.Example) (err error) {
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
	example.ChangeMap = structs.Map(example)
	example.ChangeMap["updated_at"] = time.Now()
	return tx.Table(example.TableName()).Where("id = ? AND deleted_at IS NULL", example.Id).Updates(example.ChangeMap).Error
}

func (e *Example) Get(ctx context.Context, id int) (*entity.Example, error) {
	record := &entity.Example{}
	if id == 0 {
		return nil, errors.New("get fail. need Id")
	}
	err := e.GetDB(ctx).Find(record, id).Error
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

func (e *Example) BeforeCreate(tx *gorm.DB) (err error) {
	return
}
