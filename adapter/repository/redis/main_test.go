package redis

import (
	"context"
	"testing"

	repository2 "go-hexagonal/adapter/repository"
	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init("../../../config", "config")
	log.Init()

	repository2.Clients.Redis = repository2.NewRedisClient()
	m.Run()
}
