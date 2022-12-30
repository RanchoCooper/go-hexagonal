package entity

import (
	"context"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2022/1/8
 */

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init()
	log.Init()

	repository.Clients.MySQL = repository.NewMySQLClient()
	_ = repository.Clients.MySQL.GetDB(ctx).AutoMigrate(&Example{})
	m.Run()
}
