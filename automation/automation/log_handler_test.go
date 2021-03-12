package automation

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func observableLog() (*zap.Logger, *observer.ObservedLogs) {
	core, o := observer.New(zap.NewAtomicLevelAt(zap.DebugLevel))
	return zap.New(core), o
}

func TestLogHandler_Debug(t *testing.T) {
	var (
		req      = require.New(t)
		log, obs = observableLog()
		handler  = &logHandler{}
		ctx      = logger.ContextWithValue(context.Background(), log)
		err      = handler.debug(ctx, &logDebugArgs{Message: "123abc", Fields: map[string]string{"foo": "bar"}})
	)

	req.NoError(err)
	req.Equal(1, obs.Len())

	entry := obs.TakeAll()[0]
	req.Contains(entry.Message, "123abc")
	req.Equal(entry.Level, zap.DebugLevel)
	req.Contains(entry.ContextMap(), "foo")
}

func TestLogHandler_Info(t *testing.T) {
	var (
		req      = require.New(t)
		log, obs = observableLog()
		handler  = &logHandler{}
		ctx      = logger.ContextWithValue(context.Background(), log)
		err      = handler.info(ctx, &logInfoArgs{Message: "123abc", Fields: map[string]string{"foo": "bar"}})
	)

	req.NoError(err)
	req.Equal(1, obs.Len())

	entry := obs.TakeAll()[0]
	req.Contains(entry.Message, "123abc")
	req.Equal(entry.Level, zap.InfoLevel)
	req.Contains(entry.ContextMap(), "foo")
}

func TestLogHandler_Warn(t *testing.T) {
	var (
		req      = require.New(t)
		log, obs = observableLog()
		handler  = &logHandler{}
		ctx      = logger.ContextWithValue(context.Background(), log)
		err      = handler.warn(ctx, &logWarnArgs{Message: "123abc", Fields: map[string]string{"foo": "bar"}})
	)

	req.NoError(err)
	req.Equal(1, obs.Len())

	entry := obs.TakeAll()[0]
	req.Contains(entry.Message, "123abc")
	req.Equal(entry.Level, zap.WarnLevel)
	req.Contains(entry.ContextMap(), "foo")
}

func TestLogHandler_Error(t *testing.T) {
	var (
		req      = require.New(t)
		log, obs = observableLog()
		handler  = &logHandler{}
		ctx      = logger.ContextWithValue(context.Background(), log)
		err      = handler.error(ctx, &logErrorArgs{Message: "123abc", Fields: map[string]string{"foo": "bar"}})
	)

	req.NoError(err)
	req.Equal(1, obs.Len())

	entry := obs.TakeAll()[0]
	req.Contains(entry.Message, "123abc")
	req.Equal(entry.Level, zap.ErrorLevel)
	req.Contains(entry.ContextMap(), "foo")
}
