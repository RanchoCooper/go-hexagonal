package repository

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"go-hexagonal/util/log"
)

func TestMain(m *testing.M) {
	// Initialize logging configuration
	initTestLogger()

	// Run tests
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

// initTestLogger Initialize logging configuration for test environment
func initTestLogger() {
	// Use simplest console logging configuration
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	// Initialize global logger variable
	log.Logger = logger
	log.SugaredLogger = logger.Sugar()
}
