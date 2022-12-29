package service

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/internal/adapter/repository/mysql"
	"go-hexagonal/internal/domain/model"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2021/12/25
 */

func TestExampleService_Create(t *testing.T) {
	config.Init()
	log.Init()
	repository.Init(repository.WithMySQL())

	_, mock := mysql.Client.MockClient()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `example`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	srv := NewExampleService(ctx)
	assert.NotNil(t, srv)
	assert.NotNil(t, srv.Repository)
	resp, err := srv.Create(ctx, &model.Example{
		Name:  "RanchoCooper",
		Alias: "Rancho",
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	fmt.Println(resp)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
