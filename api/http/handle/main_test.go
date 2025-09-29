package handle

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-hexagonal/util/log"
)

func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize logging configuration
	initTestLogger()

	// Run tests
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

// Initialize logging configuration for test environment
func initTestLogger() {
	// Use simplest console logging configuration
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	// Initialize global logger variable
	log.Logger = logger
	log.SugaredLogger = logger.Sugar()
}
