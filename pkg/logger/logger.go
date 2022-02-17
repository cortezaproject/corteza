package logger

import (
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/zapfilter"
)

var (
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

// Init (re)initializes global logger according to the settings
//
// It also peaks into http-server options to determinate if log events
// should be buffered for use from web console
func Init() {
	var (
		// @todo this should probably be refactored by adding a new option to LogOpt
		//       that controls if we create a buffered output as well; and when not explicitly
		//       set, we take state of web-console as a base
		hSrvOpt = options.HttpServer()
		logger  = Must(Make(options.Log()))
	)

	if hSrvOpt.WebConsoleEnabled {
		// web console is the only thing right now
		// that needs logger to buffer events for later access
		logger = withDebugBuffer(logger)
	}

	defaultLogger = logger
}

// Make creates a logger (debug or production) according to options
func Make(opt *options.LogOpt) (logger *zap.Logger, err error) {
	if opt.Debug {
		// Do we want to enable debug logger
		// with a bit more dev-friendly output
		logger, err = Debug(opt)
	} else {
		logger, err = Production(opt)
	}

	if err != nil {
		return nil, err
	}

	logger = withFilter(logger, opt.Filter)
	logger = withStacktraceLevel(logger, opt.StacktraceLevel)

	return logger, nil
}

func MakeDebugLogger() *zap.Logger {
	return Must(Debug(options.Log()))
}

// Must is a utility function that panics if given log maker returns an error
func Must(logger *zap.Logger, err error) *zap.Logger {
	if err != nil {
		panic(fmt.Errorf("failed to configure logger: %w", err))
	}

	return logger
}

// Debug prepares debug logger using options
func Debug(opt *options.LogOpt) (*zap.Logger, error) {
	var (
		// make a copy of debug options so that we do not
		dbgOpt = &options.LogOpt{
			Debug:           true,
			Level:           "debug",
			Filter:          opt.Filter,
			IncludeCaller:   opt.IncludeCaller,
			StacktraceLevel: opt.StacktraceLevel,
		}
	)

	var (
		conf = applyOptionsToConfig(zap.NewDevelopmentConfig(), dbgOpt)
	)

	// Print log level in colors
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Shorten timestamp, we do not care about the date
	conf.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05.000"))
	}

	return conf.Build()
}

func Production(opt *options.LogOpt) (*zap.Logger, error) {
	var (
		conf = applyOptionsToConfig(zap.NewProductionConfig(), opt)
	)

	return conf.Build()
}

// Applies options from environment variables
func applyOptionsToConfig(conf zap.Config, opt *options.LogOpt) zap.Config {
	// LOG_LEVEL
	conf.Level = zap.NewAtomicLevelAt(mustParseLevel(opt.Level))

	// LOG_INCLUDE_CALLER
	conf.DisableCaller = !opt.IncludeCaller

	conf.Sampling = nil

	return conf
}

// Applies filtering options
//
// This is controlled with LOG_FILTER environmental var
func withFilter(l *zap.Logger, filter string) *zap.Logger {
	if len(filter) > 0 {
		l = zap.New(zapfilter.NewFilteringCore(l.Core(), zapfilter.MustParseRules(filter)))
	}

	return l
}

// Applies stacktrace level options
//
// This is controlled with LOG_STACKTRACE_LEVEL environmental var
func withStacktraceLevel(l *zap.Logger, level string) *zap.Logger {
	return l.WithOptions(zap.AddStacktrace(mustParseLevel(level)))
}

// Adds Tee logger that copies all log messages to debug buffer
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
