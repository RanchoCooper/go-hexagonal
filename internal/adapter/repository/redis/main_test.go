package redis

import (
	"context"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/util/log"
)

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init("../../../../config", "config")
	log.Init()

	repository.Clients.Redis = repository.NewRedisClient()
	m.Run()
}
