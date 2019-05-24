package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

var (
	DefaultLevel  = zap.NewAtomicLevel()
	defaultLogger = zap.NewNop()
)

func MakeDebugLogger() *zap.Logger {
	conf := zap.NewDevelopmentConfig()
	conf.Level = DefaultLevel

	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func Init() {
	var (
		err          error
		isProduction bool
	)

	if env := os.Getenv("ENVIRONMENT"); strings.Index(env, "prod") == 0 {
		//  Try to guess if production logging from environment
		isProduction = true
	} else if vh := os.Getenv("VIRTUAL_HOST"); len(vh) > 0 {
		// Try to guess if in production logging from the fact that VIRTUAL_HOST env is set
		isProduction = true
	}

	if !isProduction {
		defaultLogger = MakeDebugLogger()
		return
	}

	conf := zap.NewProductionConfig()
	conf.Level = DefaultLevel

	// We do not want sampling
	conf.Sampling = nil
	conf.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	defaultLogger, err = conf.Build()
	if err != nil {
		panic(err)
	}
}

func Default() *zap.Logger {
	return defaultLogger
}

func SetDefault(logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}

	defaultLogger = logger
}
