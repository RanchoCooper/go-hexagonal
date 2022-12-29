package entity

import (
	"context"
	"time"

	"github.com/RanchoCooper/structs"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"go-hexagonal/internal/adapter/repository/mysql"
	"go-hexagonal/internal/domain/model"
	"go-hexagonal/internal/domain/repo"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

func NewExample() *Example {
	return &Example{}
}

type Example struct {
	Id        int                    `json:"id" structs:",omitempty,underline" gorm:"primarykey"`
	Name      string                 `json:"name" structs:",omitempty,underline"`
	Alias     string                 `json:"alias" structs:",omitempty,underline"`
	CreatedAt time.Time              `json:"created_at" structs:",omitempty,underline"`
	UpdatedAt time.Time              `json:"updated_at" structs:",omitempty,underline"`
	DeletedAt gorm.DeletedAt         `json:"deleted_at" structs:",omitempty,underline"`
	ChangeMap map[string]interface{} `json:"-" structs:"-" gorm:"-"`
}

func (e Example) TableName() string {
	return "example"
}

var _ repo.IExampleRepo = &Example{}

func (e *Example) Create(ctx context.Context, tx *gorm.DB, model *model.Example) (result *model.Example, err error) {
	if tx == nil {
		tx = mysql.Client.GetDB(ctx).Begin()
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
	row := &Example{}
	err = copier.Copy(row, model)
	if err != nil {
		return nil, err
	}
	err = tx.Create(row).Error
	if err != nil {
		return nil, err
	}
	err = copier.Copy(model, row)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (e *Example) Delete(ctx context.Context, tx *gorm.DB, id int) (err error) {
	if tx == nil {
		tx = mysql.Client.GetDB(ctx).Begin()
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
	row := &Example{}
	err = tx.Delete(row, id).Error
	// hard delete
	// err := tx.Unscoped().Delete(row, Id).Error
	return err
}

func (e *Example) Update(ctx context.Context, tx *gorm.DB, model *model.Example) (err error) {
	if tx == nil {
		tx = mysql.Client.GetDB(ctx).Begin()
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
	row := &Example{}
	err = copier.Copy(row, model)
	if err != nil {
		return err
	}
	row.ChangeMap = structs.Map(row)
	row.ChangeMap["updated_at"] = time.Now()
	return tx.Table(row.TableName()).Where("id = ? AND deleted_at IS NULL", row.Id).Updates(row.ChangeMap).Error
}

func (e *Example) GetByID(ctx context.Context, id int) (*model.Example, error) {
	query := mysql.Client.GetDB(ctx)
	row := &Example{}
	err := query.Table(row.TableName()).Find(row, id).Error
	if err != nil {
		return nil, err
	}

	result := &model.Example{}
	err = copier.Copy(result, row)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (e *Example) FindByName(ctx context.Context, name string) (*model.Example, error) {
	query := mysql.Client.GetDB(ctx)
	row := &Example{}
	err := query.Table(row.TableName()).Where("name = ?", name).Last(row).Error
	if err != nil {
		return nil, err
	}

	result := &model.Example{}
	err = copier.Copy(result, row)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (e *Example) BeforeCreate(tx *gorm.DB) (err error) {
	return
}
