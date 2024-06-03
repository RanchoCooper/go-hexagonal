package service

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/internal/domain/model"
)

func TestExampleService_Create(t *testing.T) {
	_, mock := repository.Clients.MySQL.MockClient()
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
