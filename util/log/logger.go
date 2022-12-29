package log

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"

    "go-hexagonal/config"
    "go-hexagonal/util"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

var (
    Logger        *zap.Logger
    SugaredLogger *zap.SugaredLogger
)

func initCore() zapcore.Core {
    opts := make([]zapcore.WriteSyncer, 0)
    opts = append(opts, zapcore.AddSync(os.Stdout))
    // opts = append(opts,
    //     zapcore.AddSync(&lumberjack.Logger{
    //         Filename:  filepath.Join(util.GetProjectRootPath(), config.Config.Log.SavePath, config.Config.Log.FileName),
    //         MaxSize:   config.Config.Log.MaxSize,
    //         MaxAge:    config.Config.Log.MaxAge,
    //         LocalTime: config.Config.Log.LocalTime,
    //         Compress:  config.Config.Log.Compress,
    //     }))
    syncWriter := zapcore.NewMultiWriteSyncer(opts...)

    customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
        enc.AppendString("[" + t.Format("2006-01-02T15:04:05.000Z0700") + "]")
    }
    customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
        enc.AppendString("[" + level.CapitalString() + "]")
    }
    customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
        enc.AppendString("[" + caller.TrimmedPath() + "]")
    }

    encoderConf := zapcore.EncoderConfig{
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

    encoderConf.EncodeLevel = zapcore.CapitalColorLevelEncoder

    var level zap.AtomicLevel
    if config.Config.Env.IsProd() {
        level = zap.NewAtomicLevelAt(zap.InfoLevel)
    } else {
        level = zap.NewAtomicLevelAt(zap.DebugLevel)
    }
    core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConf), syncWriter, level)
    return core
}

func lumberjackZapHook(e zapcore.Entry) error {
    lum := &lumberjack.Logger{
        Filename:   filepath.Join(util.GetProjectRootPath(), config.Config.Log.SavePath, config.Config.Log.FileName),
        MaxSize:    config.Config.Log.MaxSize,
        MaxAge:     config.Config.Log.MaxAge,
        MaxBackups: 1,
        LocalTime:  config.Config.Log.LocalTime,
        Compress:   config.Config.Log.Compress,
    }

    format := "[%-32s]\t %s\t [%s]\t %s\n"
    _, err := lum.Write([]byte(fmt.Sprintf(format,
        e.Time.Format(time.RFC3339Nano),
        e.Level.CapitalString(),
        e.Caller.TrimmedPath(),
        e.Message)),
    )
    if err != nil {
        log.Fatalf("write log fail: %s", err.Error())
    }
    return nil
}

func Init() {
    zapCore := initCore()
    Logger = zap.New(zapCore,
        zap.Hooks(lumberjackZapHook),
        zap.AddCaller(),
        zap.AddStacktrace(zap.ErrorLevel),
    )
    SugaredLogger = Logger.Sugar()
    defer Logger.Sync()
    // defer func(Logger *zap.Logger) {
    //     err := Logger.Sync()
    //     if err != nil {
    //         Logger.Error("Zap Logger fail to sync", zap.String("err", err.Error()))
    //     }
    // }(Logger)
}
