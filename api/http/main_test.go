package http

import (
	"context"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/internal/domain/entity"
	"go-hexagonal/internal/domain/service"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2022/1/7
 */

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init()
	log.Init()
	repository.Init(repository.WithMySQL(), repository.WithRedis())
	db := repository.Clients.MySQL.GetDB(ctx)
	_ = db.AutoMigrate(&entity.Example{})

	service.Init(ctx)

	m.Run()
}
