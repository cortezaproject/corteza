package logger

import (
	"os"
	"time"

	// Make sure we read the ENV from .env
	_ "github.com/joho/godotenv/autoload"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/app/options"
)

var (
	DefaultLevel  = zap.NewAtomicLevel()
	defaultLogger = zap.NewNop()
)

func MakeDebugLogger() *zap.Logger {
	conf := zap.NewDevelopmentConfig()
	conf.Level = DefaultLevel

	// Print log level in colors
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Shorten timestamp, we do not care about the date
	conf.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05.000"))
	}

	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func Init() {
	var (
		err error

		// Set INFO as defaut log level
		logLevel = zapcore.InfoLevel

		// Do we want to enable debug logger
		// with a bit more dev-friendly output
		debuggingLogger = options.EnvBool("", "LOG_DEBUG", false)
	)

	if debuggingLogger {
		logLevel = zapcore.DebugLevel
	}

	if ll, has := os.LookupEnv("LOG_LEVEL"); has {
		_ = logLevel.Set(ll)
	}

	DefaultLevel.SetLevel(logLevel)

	if debuggingLogger {
		defaultLogger = MakeDebugLogger()
		return
	}

	conf := zap.NewProductionConfig()
	conf.Level = DefaultLevel

	// We do not want sampling
	conf.Sampling = nil

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
