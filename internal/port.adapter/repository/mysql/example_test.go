package mysql

import (
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/RanchoCooper/structs"
    "github.com/stretchr/testify/assert"

    "go-hexagonal/api/http/dto"
    "go-hexagonal/internal/domain.model/entity"
)

/**
 * @author Rancho
 * @date 2022/1/8
 */

func TestExample_Create(t *testing.T) {
    exampleRepo := NewExample(NewMySQLClient())
    DB, mock := exampleRepo.MockClient()
    exampleRepo.SetDB(DB)
    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO `example`").WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()
    d := dto.CreateExampleReq{
        Name:  "rancho",
        Alias: "cooper",
    }
    example, err := exampleRepo.Create(ctx, d)
    assert.NoError(t, err)
    assert.NotEmpty(t, example.Id)
    assert.Equal(t, 1, example.Id)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestExample_Delete(t *testing.T) {
    exampleRepo := NewExample(NewMySQLClient())
    DB, mock := exampleRepo.MockClient()
    exampleRepo.SetDB(DB)
    mock.ExpectBegin()
    mock.ExpectExec("UPDATE `example`").WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()
    d := dto.DeleteExampleReq{
        Id: 1,
    }
    err := exampleRepo.Delete(ctx, d.Id)
    assert.NoError(t, err)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestExample_Save(t *testing.T) {
    exampleRepo := NewExample(NewMySQLClient())
    DB, mock := exampleRepo.MockClient()
    exampleRepo.SetDB(DB)
    mock.ExpectBegin()
    mock.ExpectExec("UPDATE `example`").WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()
    d := &entity.Example{
        Id:   1,
        Name: "random",
    }
    d.ChangeMap = structs.Map(d)
    err := exampleRepo.Save(ctx, d)
    assert.NoError(t, err)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestExample_Get(t *testing.T) {
    exampleRepo := NewExample(NewMySQLClient())
    DB, mock := exampleRepo.MockClient()
    exampleRepo.SetDB(DB)
    // FIXME
    mock.ExpectExec("SELECT * FROM `example` WHERE `example`.`id` = ? AND `example`.`deleted_at` IS NULL").WithArgs(uint(1)).WillReturnResult(sqlmock.NewResult(1, 1))
    example, err := exampleRepo.Get(ctx, 1)
    assert.NoError(t, err)
    assert.Equal(t, 1, example.Id)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}
