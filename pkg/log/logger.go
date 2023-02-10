package log

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/3d0c/storage/pkg/config"
)

var (
	instance *zap.Logger
	once     sync.Once

	logLevelMap = map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
	}
)

// InitLogger setups logger instance based on provided config
// not a thread safe, should be called from main goroutine during program startup.
// It's used in unit tests to initialize logger
func InitLogger(cfg config.Logger) {
	var (
		level       zapcore.Level
		ok          bool
		outputPaths = cfg.OutputPaths
	)

	if len(cfg.OutputPaths) == 0 {
		outputPaths = []string{"stdout"}
	}

	level, ok = logLevelMap[cfg.Level]
	if !ok {
		level = zapcore.InfoLevel
	}

	zapCfg := zap.Config{
		Level:         zap.NewAtomicLevelAt(level),
		DisableCaller: !cfg.AddCaller,
		Encoding:      "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			TimeKey:        "time",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
		},
		OutputPaths: outputPaths,
	}

	logger, err := zapCfg.Build()
	if err != nil {
		TheLogger().Error("error initializing logger, going to use default", zap.Error(err))
	}

	instance = logger
}

// TheLogger logger singleton
func TheLogger() *zap.Logger {
	once.Do(func() {
		if instance == nil {
			instance, _ = zap.NewProduction()
			InitLogger(config.Logger{})
		}
	})

	return instance
}
