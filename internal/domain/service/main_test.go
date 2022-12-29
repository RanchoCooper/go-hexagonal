package service

import (
	"context"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2021/12/25
 */

var ctx = context.TODO()

func TestMain(m *testing.M) {
	config.Init()
	log.Init()

	repository.Init(
		repository.WithMySQL(),
		repository.WithRedis(),
	)
	m.Run()
}
