package automation

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"go.uber.org/zap"
)

type (
	logHandler struct {
		reg logHandlerRegistry
	}

	logMessenger interface {
		GetMessage() (bool, string, io.Reader)
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
	msg, err := castLogArgs(args)
	if err != nil {
		return err
	}

	logger.ContextValue(ctx, zap.NewNop()).Debug(msg, zapFields(args.Fields)...)
	return nil
}

func (h logHandler) info(ctx context.Context, args *logInfoArgs) (err error) {
	msg, err := castLogArgs(args)
	if err != nil {
		return err
	}

	logger.ContextValue(ctx, zap.NewNop()).Info(msg, zapFields(args.Fields)...)
	return nil
}

func (h logHandler) warn(ctx context.Context, args *logWarnArgs) (err error) {
	msg, err := castLogArgs(args)
	if err != nil {
		return err
	}

	logger.ContextValue(ctx, zap.NewNop()).Warn(msg, zapFields(args.Fields)...)
	return nil
}

func (h logHandler) error(ctx context.Context, args *logErrorArgs) (err error) {
	msg, err := castLogArgs(args)
	if err != nil {
		return err
	}

	logger.ContextValue(ctx, zap.NewNop()).Error(msg, zapFields(args.Fields)...)
	return nil
}

func castLogArgs(msg logMessenger) (rr string, err error) {
	var msgr []byte

	_, msgS, msgR := msg.GetMessage()
	if msgR != nil {
		msgr, err = ioutil.ReadAll(msgR)
		if err != nil {
			return "", err
		}

		rr = string(msgr)
	} else {
		rr = msgS
	}

	return
}
