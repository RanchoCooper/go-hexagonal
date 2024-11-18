package postgre

import (
	"testing"

	"go-hexagonal/config"
)

func TestMain(m *testing.M) {
	config.Init("../../config", "config")

	m.Run()
}
