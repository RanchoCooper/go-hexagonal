package repository

import (
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

func TestMain(m *testing.M) {
	config.Init("../../config", "config.yaml")
	log.Init()

	m.Run()
}
