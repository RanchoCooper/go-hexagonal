package entity

import (
	"context"
	"time"

	"github.com/RanchoCooper/structs"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"go-hexagonal/adapter/repository"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

func NewExample() *EntityExample {
	return &EntityExample{}
}

type EntityExample struct {
	Id        int                    `json:"id" gorm:"primarykey" structs:",omitempty,underline"`
	Name      string                 `json:"name" structs:",omitempty,underline"`
	Alias     string                 `json:"alias" structs:",omitempty,underline"`
	CreatedAt time.Time              `json:"created_at" structs:",omitempty,underline"`
	UpdatedAt time.Time              `json:"updated_at" structs:",omitempty,underline"`
	DeletedAt gorm.DeletedAt         `json:"deleted_at" structs:",omitempty,underline"`
	ChangeMap map[string]interface{} `json:"-" gorm:"-" structs:"-"`
}

func (e EntityExample) TableName() string {
	return "example"
}

var _ repo.IExampleRepo = &EntityExample{}

func autoCommit(db *gorm.DB) error {
	return db.Commit().Error
}

func (e *EntityExample) Create(ctx context.Context, tr *repository.Transaction, model *model.Example) (result *model.Example, err error) {
	entity := &EntityExample{}
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

func (e *EntityExample) Delete(ctx context.Context, tr *repository.Transaction, id int) (err error) {
	entity := &EntityExample{}

	db := tr.Conn(ctx)
	err = db.Delete(entity, id).Error
	// hard delete
	// err := tx.Unscoped().Delete(entity, Id).Error
	return err
}

func (e *EntityExample) Update(ctx context.Context, tr *repository.Transaction, model *model.Example) (err error) {
	entity := &EntityExample{}
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

func (e *EntityExample) GetByID(ctx context.Context, tr *repository.Transaction, id int) (domain *model.Example, err error) {
	entity := &EntityExample{}

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

func (e *EntityExample) FindByName(ctx context.Context, tr *repository.Transaction, name string) (model *model.Example, err error) {
	entity := &EntityExample{}

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
