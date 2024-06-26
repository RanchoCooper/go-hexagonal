package entity

import (
	"context"
	"time"

	"github.com/RanchoCooper/structs"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/internal/domain/model"
	"go-hexagonal/internal/domain/repo"
)

func NewExample() *Example {
	return &Example{}
}

type Example struct {
	Id        int                    `json:"id" gorm:"primarykey" structs:",omitempty,underline"`
	Name      string                 `json:"name" structs:",omitempty,underline"`
	Alias     string                 `json:"alias" structs:",omitempty,underline"`
	CreatedAt time.Time              `json:"created_at" structs:",omitempty,underline"`
	UpdatedAt time.Time              `json:"updated_at" structs:",omitempty,underline"`
	DeletedAt gorm.DeletedAt         `json:"deleted_at" structs:",omitempty,underline"`
	ChangeMap map[string]interface{} `json:"-" gorm:"-" structs:"-"`
}

func (e Example) TableName() string {
	return "example"
}

var _ repo.IExampleRepo = &Example{}

func autoCommit(db *gorm.DB) error {
	return db.Commit().Error
}

func (e *Example) Create(ctx context.Context, tr *repository.Transaction, model *model.Example) (result *model.Example, err error) {
	entity := &Example{}
	err = copier.Copy(entity, model)
	if err != nil {
		return nil, errors.Wrap(err, "copier fail")
	}

	db := tr.Conn(ctx)
	err = db.Create(entity).Error
	if err != nil {
		return nil, err
	}

	err = copier.Copy(model, entity)
	if err != nil {
		return nil, errors.Wrap(err, "copier fail")
	}

	return model, nil
}

func (e *Example) Delete(ctx context.Context, tr *repository.Transaction, id int) (err error) {
	entity := &Example{}

	db := tr.Conn(ctx)
	err = db.Delete(entity, id).Error
	// hard delete
	// err := tx.Unscoped().Delete(entity, Id).Error
	return err
}

func (e *Example) Update(ctx context.Context, tr *repository.Transaction, model *model.Example) (err error) {
	entity := &Example{}
	err = copier.Copy(entity, model)
	if err != nil {
		return errors.Wrap(err, "copier fail")
	}
	entity.ChangeMap = structs.Map(entity)
	entity.ChangeMap["updated_at"] = time.Now()

	db := tr.Conn(ctx)
	db = db.Table(entity.TableName()).Where("id = ? AND deleted_at IS NULL", entity.Id).Updates(entity.ChangeMap)

	return db.Error
}

func (e *Example) GetByID(ctx context.Context, tr *repository.Transaction, id int) (domain *model.Example, err error) {
	entity := &Example{}

	db := tr.Conn(ctx)
	db = db.Table(entity.TableName()).Find(entity, id)

	if db.Error != nil {
		return nil, err
	}

	domain = &model.Example{}
	err = copier.Copy(domain, entity)
	if err != nil {
		return nil, errors.Wrap(err, "copier fail")
	}
	return domain, nil
}

func (e *Example) FindByName(ctx context.Context, tr *repository.Transaction, name string) (model *model.Example, err error) {
	entity := &Example{}

	db := tr.Conn(ctx)
	db.Table(entity.TableName()).Where("name = ?", name).Last(entity)
	if db.Error != nil {
		return nil, err
	}

	err = copier.Copy(model, entity)
	if err != nil {
		return nil, errors.Wrap(err, "copier fail")
	}
	return model, err
}
