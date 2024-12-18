package http

import (
	"context"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/internal/adapter/repository/mysql/entity"
	"go-hexagonal/internal/domain/service"
	"go-hexagonal/util/log"
)

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init("../../config", "config")
	log.Init()

	repository.Init(repository.WithMySQL(), repository.WithRedis())
	_ = repository.Clients.MySQL.GetDB(ctx).AutoMigrate(&entity.Example{})

	service.Init(ctx)

	m.Run()
}
