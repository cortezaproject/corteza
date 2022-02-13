package logger

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/zapfilter"
)

var (
	opt           = options.Log()
	defaultLogger = zap.NewNop()
)

func Default() *zap.Logger {
	if defaultLogger == nil {
		return zap.NewNop()
	}

	return defaultLogger
}

func SetDefault(logger *zap.Logger) {
	defaultLogger = logger
}

// Init (re)initializes logger according to the settings
func Init() {
	if opt.Debug {
		// Do we want to enable debug logger
		// with a bit more dev-friendly output
		defaultLogger = MakeDebugLogger()
		defaultLogger.Debug("full debug mode enabled")
		return
	}

	var (
		conf        = applyOptions(zap.NewProductionConfig(), opt)
		logger, err = conf.Build()
	)

	if err != nil {
		panic(err)
	}

	logger = applySpecials(defaultLogger, opt)
	logger = withDebugBuffer(logger)

	defaultLogger = logger
}

func MakeDebugLogger() *zap.Logger {
	dbgOpt := *opt
	dbgOpt.Debug = true
	dbgOpt.Level = "debug"

	var (
		conf = applyOptions(zap.NewDevelopmentConfig(), &dbgOpt)
	)

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

	logger = withDebugBuffer(logger)

	return applySpecials(logger, &dbgOpt)
}

// Applies options from environment variables
func applyOptions(conf zap.Config, opt *options.LogOpt) zap.Config {
	// LOG_LEVEL
	conf.Level = zap.NewAtomicLevelAt(mustParseLevel(opt.Level))

	// LOG_INCLUDE_CALLER
	conf.DisableCaller = !opt.IncludeCaller

	conf.Sampling = nil

	return conf
}

// Applies "special" options - filtering and conditional stack-level
func applySpecials(l *zap.Logger, opt *options.LogOpt) *zap.Logger {
	if len(opt.Filter) > 0 {
		// LOG_FILTER
		l = zap.New(zapfilter.NewFilteringCore(l.Core(), zapfilter.MustParseRules(opt.Filter)))
	}

	// LOG_STACKTRACE_LEVEL
	return l.WithOptions(zap.AddStacktrace(mustParseLevel(opt.StacktraceLevel)))
}

// adds Tee logger that copies all log messages to debug buffer
func withDebugBuffer(in *zap.Logger) *zap.Logger {
	return zap.New(zapcore.NewTee(
		in.Core(),

		DebugBufferedLogger(debugLogRR),
	))
}

func mustParseLevel(l string) (o zapcore.Level) {
	if err := o.Set(l); err != nil {
		panic(err)
	}

	return
}
