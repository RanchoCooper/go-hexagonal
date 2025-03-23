package service

import (
	"testing"

	"go-hexagonal/adapter/repository"
	"go-hexagonal/config"
	"go-hexagonal/util/log"

	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// Initialize configuration and logging
	config.Init("../../config", "config")
	log.Init()

	repository.Clients = &repository.ClientContainer{
		MySQL: repository.NewMySQLClient(&gorm.DB{}),
		Redis: repository.NewRedisClient(),
	}

	m.Run()
}
