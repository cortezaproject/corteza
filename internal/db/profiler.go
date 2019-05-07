package db

import (
	"time"

	"github.com/titpetric/factory"
	"go.uber.org/zap"
)

// zapProfiler logs query statistics to zap.logger
type (
	zapProfiler struct {
		logger *zap.Logger
	}
)

func ZapProfiler(logger *zap.Logger) *zapProfiler {
	return &zapProfiler{
		logger: logger,
	}
}

// Post prints the query statistics to stdout
func (p zapProfiler) Post(c *factory.DatabaseProfilerContext) {
	// @todo when factory.DatabaseProfilerContext gets access to context from
	//       db functions, try to extract RequestID with middleware.GetReqID()

	p.logger.Debug(
		c.Query,
		zap.Any("args", c.Args),
		zap.Float64("duration", time.Since(c.Time).Seconds()))
}

// Flush stdout (no-op for this profiler)
func (zapProfiler) Flush() {
}
