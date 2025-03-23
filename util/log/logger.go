// Package log provides logging functionality for the application
package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"go-hexagonal/config"
	"go-hexagonal/util"
)

// Global logger instances, maintained for compatibility with existing code
var (
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
)

// Options defines the configuration options for the logger
type Options struct {
	// Log level
	Level zapcore.Level
	// Whether to output to console
	EnableConsole bool
	// Whether to output to file
	EnableFile bool
	// Whether to enable colored output
	EnableColor bool
	// Whether to add caller information
	EnableCaller bool
	// Whether to add stack traces
	EnableStacktrace bool
	// File configuration, required if EnableFile is true
	FileConfig *FileConfig
}

// FileConfig defines the configuration for log files
type FileConfig struct {
	// Directory path to save log files
	SavePath string
	// Log file name
	FileName string
	// Maximum size in MB
	MaxSize int
	// Maximum age in days
	MaxAge int
	// Whether to use local time
	LocalTime bool
	// Whether to compress old log files
	Compress bool
	// Maximum number of backup files
	MaxBackups int
}

// AppLogger wraps the zap logger
type AppLogger struct {
	zap            *zap.Logger
	sugar          *zap.SugaredLogger
	options        *Options
	lumberjackHook func(zapcore.Entry) error
}

// ParseLogLevel converts a string log level to zapcore.Level
func ParseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		// Default to info level if invalid or empty
		return zapcore.InfoLevel
	}
}

// DefaultOptions returns the default configuration options
func DefaultOptions() *Options {
	return &Options{
		Level:            zapcore.InfoLevel,
		EnableConsole:    true,
		EnableFile:       false,
		EnableColor:      true,
		EnableCaller:     true,
		EnableStacktrace: true,
		FileConfig:       nil,
	}
}

// Option defines a function type for configuring options
type Option func(*Options)

// WithLevel sets the log level
func WithLevel(level zapcore.Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithConsole enables or disables console output
func WithConsole(enable bool) Option {
	return func(o *Options) {
		o.EnableConsole = enable
	}
}

// WithFile enables file output and sets the file configuration
func WithFile(fileConfig *FileConfig) Option {
	return func(o *Options) {
		o.EnableFile = true
		o.FileConfig = fileConfig
	}
}

// WithColor enables or disables colored output
func WithColor(enable bool) Option {
	return func(o *Options) {
		o.EnableColor = enable
	}
}

// WithCaller enables or disables caller information
func WithCaller(enable bool) Option {
	return func(o *Options) {
		o.EnableCaller = enable
	}
}

// WithStacktrace enables or disables stack traces
func WithStacktrace(enable bool) Option {
	return func(o *Options) {
		o.EnableStacktrace = enable
	}
}

// FileConfigFromGlobal creates a FileConfig from global configuration
func FileConfigFromGlobal() *FileConfig {
	return &FileConfig{
		SavePath:   config.GlobalConfig.Log.SavePath,
		FileName:   config.GlobalConfig.Log.FileName,
		MaxSize:    config.GlobalConfig.Log.MaxSize,
		MaxAge:     config.GlobalConfig.Log.MaxAge,
		LocalTime:  config.GlobalConfig.Log.LocalTime,
		Compress:   config.GlobalConfig.Log.Compress,
		MaxBackups: 1,
	}
}

// New creates a new logger instance
func New(opts ...Option) (*AppLogger, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Return error if file output is enabled but no file config is provided
	if options.EnableFile && options.FileConfig == nil {
		return nil, fmt.Errorf("file config is required when file output is enabled")
	}

	logger := &AppLogger{
		options: options,
	}

	// Create core logging components
	cores := logger.buildCores()
	zapOptions := logger.buildZapOptions()

	// Create zap logger
	zapLogger := zap.New(zapcore.NewTee(cores...), zapOptions...)
	logger.zap = zapLogger
	logger.sugar = zapLogger.Sugar()

	return logger, nil
}

// buildCores creates the logging output cores
func (l *AppLogger) buildCores() []zapcore.Core {
	encoderConfig := l.getEncoderConfig()
	cores := make([]zapcore.Core, 0)

	// Add console output
	if l.options.EnableConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(l.options.Level),
		)
		cores = append(cores, consoleCore)
	}

	// Add file output
	if l.options.EnableFile && l.options.FileConfig != nil {
		fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)

		l.lumberjackHook = l.createLumberjackHook()

		// Create a core that doesn't write directly to file, we'll use the hook for actual file output
		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(os.Stdout), // Temporary output to stdout, actual output handled by hook
			zap.NewAtomicLevelAt(l.options.Level),
		)
		cores = append(cores, fileCore)
	}

	return cores
}

// getEncoderConfig gets the encoder configuration
func (l *AppLogger) getEncoderConfig() zapcore.EncoderConfig {
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02T15:04:05.000Z0700") + "]")
	}
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,
		EncodeLevel:    customLevelEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeCaller:   customCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	if l.options.EnableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return encoderConfig
}

// buildZapOptions builds zap options
func (l *AppLogger) buildZapOptions() []zap.Option {
	options := make([]zap.Option, 0)

	if l.options.EnableFile && l.lumberjackHook != nil {
		options = append(options, zap.Hooks(l.lumberjackHook))
	}

	if l.options.EnableCaller {
		options = append(options, zap.AddCaller())
	}

	if l.options.EnableStacktrace {
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}

	return options
}

// createLumberjackHook creates a lumberjack hook for file output
func (l *AppLogger) createLumberjackHook() func(zapcore.Entry) error {
	return func(e zapcore.Entry) error {
		if !l.options.EnableFile || l.options.FileConfig == nil {
			return nil
		}

		fc := l.options.FileConfig
		lum := &lumberjack.Logger{
			Filename:   filepath.Join(util.GetProjectRootPath(), fc.SavePath, fc.FileName),
			MaxSize:    fc.MaxSize,
			MaxAge:     fc.MaxAge,
			MaxBackups: fc.MaxBackups,
			LocalTime:  fc.LocalTime,
			Compress:   fc.Compress,
		}

		format := "[%-32s]\t %s\t [%s]\t %s\n"
		_, err := lum.Write([]byte(fmt.Sprintf(format,
			e.Time.Format(time.RFC3339Nano),
			e.Level.CapitalString(),
			e.Caller.TrimmedPath(),
			e.Message)),
		)
		return err
	}
}

// Zap returns the wrapped zap logger
func (l *AppLogger) Zap() *zap.Logger {
	return l.zap
}

// Sugar returns the wrapped sugared logger
func (l *AppLogger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// Sync synchronizes the log buffer
func (l *AppLogger) Sync() error {
	return l.zap.Sync()
}

// Close closes the logger and synchronizes the buffer
func (l *AppLogger) Close() error {
	return l.Sync()
}

// Init initializes the global logger instances (for compatibility with existing code)
func Init() {
	var opts []Option

	// Determine log level from configuration
	logLevel := zapcore.InfoLevel
	if config.GlobalConfig.Log != nil && config.GlobalConfig.Log.Level != "" {
		// Use configured log level if available
		logLevel = ParseLogLevel(config.GlobalConfig.Log.Level)
	} else if !config.GlobalConfig.Env.IsProd() {
		// Fall back to debug level in non-production environments
		logLevel = zapcore.DebugLevel
	}
	opts = append(opts, WithLevel(logLevel))

	// Configure console output (defaults to true if not specified)
	enableConsole := true
	if config.GlobalConfig.Log != nil {
		enableConsole = config.GlobalConfig.Log.EnableConsole
	}
	opts = append(opts, WithConsole(enableConsole))

	// Configure colorized output (defaults to enabled in non-production)
	enableColor := !config.GlobalConfig.Env.IsProd()
	if config.GlobalConfig.Log != nil {
		enableColor = config.GlobalConfig.Log.EnableColor
	}
	opts = append(opts, WithColor(enableColor))

	// Configure caller information (defaults to true if not specified)
	enableCaller := true
	if config.GlobalConfig.Log != nil {
		enableCaller = config.GlobalConfig.Log.EnableCaller
	}
	opts = append(opts, WithCaller(enableCaller))

	// Configure stack traces (defaults to true if not specified)
	enableStacktrace := true
	if config.GlobalConfig.Log != nil {
		enableStacktrace = config.GlobalConfig.Log.EnableStacktrace
	}
	opts = append(opts, WithStacktrace(enableStacktrace))

	// Add file output if global config has log file settings
	if config.GlobalConfig.Log != nil && config.GlobalConfig.Log.SavePath != "" {
		opts = append(opts, WithFile(FileConfigFromGlobal()))
	}

	// Create new logger instance
	logger, err := New(opts...)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Set global variables
	Logger = logger.Zap()
	SugaredLogger = logger.Sugar()

	// Register deferred sync
	// Direct Sync() call is commented out because it might cause errors on program shutdown
	// defer Logger.Sync()
}
