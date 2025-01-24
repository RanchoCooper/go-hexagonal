package service

import (
	"context"
	"testing"

	"go-hexagonal/adapter/repository"
	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init("../../config", "config")
	log.Init()

	repository.Init(
		repository.WithMySQL(),
		repository.WithRedis(),
	)
	m.Run()
}
