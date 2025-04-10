package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go-hexagonal/config"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected zapcore.Level
	}{
		{"debug", "debug", zapcore.DebugLevel},
		{"info", "info", zapcore.InfoLevel},
		{"warn", "warn", zapcore.WarnLevel},
		{"warning", "warning", zapcore.WarnLevel},
		{"error", "error", zapcore.ErrorLevel},
		{"dpanic", "dpanic", zapcore.DPanicLevel},
		{"panic", "panic", zapcore.PanicLevel},
		{"fatal", "fatal", zapcore.FatalLevel},
		{"empty", "", zapcore.InfoLevel},
		{"invalid", "invalid", zapcore.InfoLevel},
		{"case insensitive", "INFO", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := ParseLogLevel(tt.level)
			assert.Equal(t, tt.expected, level)
		})
	}
}

func TestDefaultOptions(t *testing.T) {
	options := DefaultOptions()

	assert.Equal(t, zapcore.InfoLevel, options.Level)
	assert.True(t, options.EnableConsole)
	assert.False(t, options.EnableFile)
	assert.True(t, options.EnableColor)
	assert.True(t, options.EnableCaller)
	assert.True(t, options.EnableStacktrace)
	assert.Nil(t, options.FileConfig)
}

func TestWithOptions(t *testing.T) {
	// Test various option functions
	options := DefaultOptions()

	WithLevel(zapcore.DebugLevel)(options)
	assert.Equal(t, zapcore.DebugLevel, options.Level)

	WithConsole(false)(options)
	assert.False(t, options.EnableConsole)

	WithColor(false)(options)
	assert.False(t, options.EnableColor)

	WithCaller(false)(options)
	assert.False(t, options.EnableCaller)

	WithStacktrace(false)(options)
	assert.False(t, options.EnableStacktrace)

	fileConfig := &FileConfig{
		SavePath:   "logs",
		FileName:   "test.log",
		MaxSize:    10,
		MaxAge:     7,
		LocalTime:  true,
		Compress:   true,
		MaxBackups: 3,
	}

	WithFile(fileConfig)(options)
	assert.True(t, options.EnableFile)
	assert.Equal(t, fileConfig, options.FileConfig)
}

func TestNew(t *testing.T) {
	// Test creating a new logger
	logger, err := New()
	require.NoError(t, err)
	require.NotNil(t, logger)
	require.NotNil(t, logger.zap)
	require.NotNil(t, logger.sugar)

	// Test creation with file configuration
	tempDir, err := ioutil.TempDir("", "logger-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fileConfig := &FileConfig{
		SavePath:   tempDir,
		FileName:   "test.log",
		MaxSize:    1,
		MaxAge:     1,
		LocalTime:  true,
		Compress:   false,
		MaxBackups: 1,
	}

	logger, err = New(WithFile(fileConfig))
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test file output enabled but no file config
	logger, err = New(WithConsole(false))
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Test error case
	_, err = New(func(o *Options) {
		o.EnableFile = true
		o.FileConfig = nil
	})
	require.Error(t, err)
}

func TestLogContext(t *testing.T) {
	ctx := NewLogContext()
	assert.Equal(t, "system", ctx.Component)

	// Test chained calls
	ctx = ctx.WithRequestID("req-123")
	assert.Equal(t, "req-123", ctx.RequestID)

	ctx = ctx.WithUserID("user-456")
	assert.Equal(t, "user-456", ctx.UserID)

	ctx = ctx.WithTraceID("trace-789")
	assert.Equal(t, "trace-789", ctx.TraceID)

	ctx = ctx.WithSpanID("span-012")
	assert.Equal(t, "span-012", ctx.SpanID)

	ctx = ctx.WithOperation("test-op")
	assert.Equal(t, "test-op", ctx.Operation)

	ctx = ctx.WithComponent("test-comp")
	assert.Equal(t, "test-comp", ctx.Component)

	// Test conversion to fields
	fields := ctx.ToFields()
	assert.Len(t, fields, 6)

	fieldMap := make(map[string]string)
	for _, field := range fields {
		switch field.Key {
		case "request_id", "user_id", "trace_id", "span_id", "operation", "component":
			fieldMap[field.Key] = field.String
		}
	}

	assert.Equal(t, "req-123", fieldMap["request_id"])
	assert.Equal(t, "user-456", fieldMap["user_id"])
	assert.Equal(t, "trace-789", fieldMap["trace_id"])
	assert.Equal(t, "span-012", fieldMap["span_id"])
	assert.Equal(t, "test-op", fieldMap["operation"])
	assert.Equal(t, "test-comp", fieldMap["component"])
}

func TestLogContextMethods(t *testing.T) {
	// Create a temporary file to capture log output
	tempFile, err := ioutil.TempFile("", "logger-test-*.log")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Create a logger for testing
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(tempFile),
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
	)

	zapLogger := zap.New(core)
	logger := &AppLogger{
		zap:     zapLogger,
		sugar:   zapLogger.Sugar(),
		options: DefaultOptions(),
	}

	// Create test context
	ctx := NewLogContext().
		WithRequestID("req-test").
		WithComponent("test-component")

	// Test different log levels
	logger.DebugContext(ctx, "debug message", zap.String("key", "value"))
	logger.InfoContext(ctx, "info message", zap.Int("code", 200))
	logger.WarnContext(ctx, "warn message")
	logger.ErrorContext(ctx, "error message", zap.Error(assert.AnError))

	// Test nil context
	logger.InfoContext(nil, "info with nil context")

	// Test Sugar versions
	logger.SugaredDebugContext(ctx, "sugar debug %s", "message")
	logger.SugaredInfoContext(ctx, "sugar info %d", 123)
	logger.SugaredWarnContext(ctx, "sugar warn %v", true)
	logger.SugaredErrorContext(ctx, "sugar error %v", assert.AnError)

	// Ensure all logs are written to file
	logger.Sync()

	// Read log content for verification
	tempFile.Seek(0, 0)
	content, err := ioutil.ReadAll(tempFile)
	require.NoError(t, err)

	logContent := string(content)

	// Verify log content contains expected messages and fields
	assert.Contains(t, logContent, "debug message")
	assert.Contains(t, logContent, "info message")
	assert.Contains(t, logContent, "warn message")
	assert.Contains(t, logContent, "error message")
	assert.Contains(t, logContent, "req-test")
	assert.Contains(t, logContent, "test-component")
	assert.Contains(t, logContent, "sugar debug message")
	assert.Contains(t, logContent, "sugar info 123")
}

func TestInit(t *testing.T) {
	// Save original GlobalConfig
	originalConfig := config.GlobalConfig
	defer func() { config.GlobalConfig = originalConfig }()

	// Initialize GlobalConfig if it's nil to prevent panic
	if config.GlobalConfig == nil {
		config.GlobalConfig = &config.Config{
			Env: config.Env("test"),
			Log: &config.LogConfig{
				Level:            "debug",
				EnableConsole:    true,
				EnableColor:      true,
				EnableCaller:     true,
				EnableStacktrace: true,
				SavePath:         os.TempDir(),
				FileName:         "test.log",
				MaxSize:          1,
				MaxAge:           1,
				LocalTime:        true,
				Compress:         false,
			},
		}
	} else {
		// Ensure Log is not nil
		if config.GlobalConfig.Log == nil {
			config.GlobalConfig.Log = &config.LogConfig{
				Level:            "debug",
				EnableConsole:    true,
				EnableColor:      true,
				EnableCaller:     true,
				EnableStacktrace: true,
				SavePath:         os.TempDir(),
				FileName:         "test.log",
				MaxSize:          1,
				MaxAge:           1,
				LocalTime:        true,
				Compress:         false,
			}
		}
	}

	// Initialize logger
	Init()

	// Verify global variables are set
	assert.NotNil(t, Logger)
	assert.NotNil(t, SugaredLogger)
}
