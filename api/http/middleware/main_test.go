package middleware

import (
	"os"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

func TestMain(m *testing.M) {
	config.Init("../../../config", "config")
	log.Init()

	exitCode := m.Run()
	os.Exit(exitCode)
}
