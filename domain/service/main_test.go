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
	// Initialize configuration and logging
	config.Init("../../config", "config")
	log.Init()

	// Initialize repositories for testing
	repository.Init(
		repository.WithMySQL(),
		repository.WithRedis(),
	)

	m.Run()
}
