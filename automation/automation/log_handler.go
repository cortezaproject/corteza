package automation

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"go.uber.org/zap"
)

type (
	logHandler struct {
		reg logHandlerRegistry
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
	}

	h.register()
	return h
}

func (h logHandler) debug(ctx context.Context, args *logDebugArgs) (err error) {
	logger.ContextValue(ctx, zap.NewNop()).Debug(args.Message, zapFields(args.Fields)...)
	return nil
}

func (h logHandler) info(ctx context.Context, args *logInfoArgs) (err error) {
	logger.ContextValue(ctx, zap.NewNop()).Info(args.Message, zapFields(args.Fields)...)
	return nil
}

func (h logHandler) warn(ctx context.Context, args *logWarnArgs) (err error) {
	logger.ContextValue(ctx, zap.NewNop()).Warn(args.Message, zapFields(args.Fields)...)
	return nil
}

func (h logHandler) error(ctx context.Context, args *logErrorArgs) (err error) {
	logger.ContextValue(ctx, zap.NewNop()).Error(args.Message, zapFields(args.Fields)...)
	return nil
}
