package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DefaultLevel  = zap.NewAtomicLevel()
	defaultLogger *zap.Logger
)

func Init(level zapcore.Level) {
	DefaultLevel.SetLevel(level)

	var (
		err  error
		conf = zap.Config{
			Level:            DefaultLevel,
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	)

	// Make logstash's life a bit easier
	conf.EncoderConfig.LevelKey = "appLogLevel"
	conf.EncoderConfig.MessageKey = "message"
	conf.EncoderConfig.TimeKey = "@timestamp"
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	defaultLogger, err = conf.Build()
	if err != nil {
		panic(err)
	}
}

func Default() *zap.Logger {
	if defaultLogger == nil {
		panic("default logger not initialised")
	}

	return defaultLogger
}
