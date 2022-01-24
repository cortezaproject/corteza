package automation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func observableLog() (*zap.Logger, *observer.ObservedLogs) {
	core, o := observer.New(zap.NewAtomicLevelAt(zap.DebugLevel))
	return zap.New(core), o
}

func TestLogHandler_Debug(t *testing.T) {
	var (
		req      = require.New(t)
		log, obs = observableLog()
		handler  = &logHandler{logger: log}
		err      = handler.debug(context.Background(), &logDebugArgs{Message: "123abc", Fields: map[string]string{"foo": "bar"}})
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
		handler  = &logHandler{logger: log}
		err      = handler.info(context.Background(), &logInfoArgs{Message: "123abc", Fields: map[string]string{"foo": "bar"}})
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
		handler  = &logHandler{logger: log}
		err      = handler.warn(context.Background(), &logWarnArgs{Message: "123abc", Fields: map[string]string{"foo": "bar"}})
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
		handler  = &logHandler{logger: log}
		err      = handler.error(context.Background(), &logErrorArgs{Message: "123abc", Fields: map[string]string{"foo": "bar"}})
	)

	req.NoError(err)
	req.Equal(1, obs.Len())

	entry := obs.TakeAll()[0]
	req.Contains(entry.Message, "123abc")
	req.Equal(entry.Level, zap.ErrorLevel)
	req.Contains(entry.ContextMap(), "foo")
}
