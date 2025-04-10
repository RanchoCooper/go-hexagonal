package clean_arch

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"go-hexagonal/util/log"
)

func TestMain(m *testing.M) {
	// 初始化日志配置
	initTestLogger()

	// 运行测试
	exitCode := m.Run()

	// 退出
	os.Exit(exitCode)
}

// initTestLogger 初始化测试环境的日志配置
func initTestLogger() {
	// 使用最简单的控制台日志配置
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	// 初始化全局日志变量
	log.Logger = logger
	log.SugaredLogger = logger.Sugar()
}
