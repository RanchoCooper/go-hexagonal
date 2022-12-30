package entity

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"go-hexagonal/api/dto"
	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/internal/adapter/repository/mysql"
	"go-hexagonal/internal/domain/model"
)

/**
 * @author Rancho
 * @date 2022/1/8
 */

func TestExample_Create(t *testing.T) {
	exampleRepo := NewExample()
	t.Run("run with nil transaction", func(t *testing.T) {
		_, mock := repository.Clients.MySQL.MockClient()
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `example`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		e := &model.Example{
			Name:  "RanchoCooper",
			Alias: "Rancho",
		}
		example, err := exampleRepo.Create(ctx, nil, e)
		assert.NoError(t, err)
		assert.NotEmpty(t, example.Id)
		assert.Equal(t, 1, example.Id)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("run with new transaction", func(t *testing.T) {
		_, mock := repository.Clients.MySQL.MockClient()
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `example`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		e := &model.Example{
			Name:  "rancho",
			Alias: "cooper",
		}
		tr := mysql.NewTransaction(ctx, &sql.TxOptions{
			Isolation: sql.LevelReadUncommitted,
			ReadOnly:  false,
		})
		example, err := exampleRepo.Create(ctx, tr, e)
		assert.NoError(t, err)
		assert.NotEmpty(t, example.Id)
		assert.Equal(t, 1, example.Id)
		tr.Session.Commit()

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestExample_Delete(t *testing.T) {
	exampleRepo := NewExample()
	_, mock := repository.Clients.MySQL.MockClient()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `example`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	d := dto.DeleteExampleReq{
		Id: 1,
	}
	err := exampleRepo.Delete(ctx, nil, d.Id)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestExample_Update(t *testing.T) {
	exampleRepo := NewExample()
	_, mock := repository.Clients.MySQL.MockClient()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `example`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	e := &model.Example{
		Id:   1,
		Name: "random",
	}
	err := exampleRepo.Update(ctx, nil, e)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestExample_GetByID(t *testing.T) {
	exampleRepo := NewExample()
	_, mock := repository.Clients.MySQL.MockClient()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `example` WHERE `example`.`id` = ? AND `example`.`deleted_at` IS NULL")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test1"))
	example, err := exampleRepo.GetByID(ctx, nil, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, example.Id)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
