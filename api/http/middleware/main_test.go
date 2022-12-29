package middleware

import (
	"os"
	"testing"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2022/12/29
 */

func TestMain(m *testing.M) {
	config.Init()
	log.Init()

	exitCode := m.Run()
	os.Exit(exitCode)
}
