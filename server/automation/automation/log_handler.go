package automation

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/logger"
	"go.uber.org/zap"
)

type (
	logHandler struct {
		reg    logHandlerRegistry
		logger *zap.Logger
	}
)

func zapFields(fields map[string]string) []zap.Field {
	ff := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		ff = append(ff, zap.String(k, v))
	}

	return ff
}

func LogHandler(reg logHandlerRegistry) *logHandler {
	h := &logHandler{
		reg: reg,

		// this is a temporary solution how to get to the logging facility
		// without being dependent on workflow logging settings
		//
		// there are still general logging settings that have a final say what is logged and what not
		//
		// With LOG_DEBUG=true you'll see all log levels from these functions.
		// With debug logging disabled and LOG_LEVEL set on error, warn, info or debug, you'll filter out
		// logging by verbosity:
		//   error will show
		// and with LOG_LEVEL=error on you'll just see logged errors
		// while LOG_LEVEL=info will show info, warn and error but ignore debug.
		logger: logger.Default().
			Named("workflow").
			WithOptions(zap.WithCaller(false)),
	}

	h.register()
	return h
}

func (h logHandler) debug(_ context.Context, args *logDebugArgs) (err error) {
	h.logger.Debug(args.Message, zapFields(args.Fields)...)
	return nil
}

func (h logHandler) info(_ context.Context, args *logInfoArgs) (err error) {
	h.logger.Info(args.Message, zapFields(args.Fields)...)
	return nil
}

func (h logHandler) warn(_ context.Context, args *logWarnArgs) (err error) {
	h.logger.Warn(args.Message, zapFields(args.Fields)...)
	return nil
}

func (h logHandler) error(_ context.Context, args *logErrorArgs) (err error) {
	h.logger.Error(args.Message, zapFields(args.Fields)...)
	return nil
}
