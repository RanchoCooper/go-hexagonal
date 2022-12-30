package redis

import (
	"context"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2022/12/29
 */

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init()
	log.Init()

	repository.Clients.Redis = repository.NewRedisClient()
	m.Run()
}
