package rdbms

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/logger"
	"go.uber.org/zap"
)

// log() returns named logger with caller skip and stacktrace set to Fatal
//
// It check the given context logger
func (s Store) log(ctx context.Context) *zap.Logger {
	return logger.ContextValue(ctx, s.logger).
		Named("store.rdbms").
		WithOptions(zap.AddCallerSkip(2))
}

func (s Store) SetLogger(logger *zap.Logger) {
	s.logger = logger
}
