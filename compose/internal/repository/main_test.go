// +build integration

package repository

import (
	"os"
	"testing"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.SetDefault(logger.MakeDebugLogger())

	factory.Database.Add("compose", os.Getenv("COMPOSE_DB_DSN"))
	db := factory.Database.MustGet("compose")
	db.Profiler = &factory.DatabaseProfilerStdout{}

	os.Exit(m.Run())
}

// zapProfiler logs query statistics to zap.logger
type (
	testLogProfiler struct {
		logger testProfilerLogger
	}

	testProfilerLogger interface {
		Logf(format string, args ...interface{})
	}
)

func newTestLogProfiler(logger testProfilerLogger) *testLogProfiler {
	return &testLogProfiler{
		logger: logger,
	}
}

// Post prints the query statistics to stdout
func (p testLogProfiler) Post(c *factory.DatabaseProfilerContext) {
	p.logger.Logf(
		"%s\nArgs: %v\nDuration: %fs",
		c.Query,
		c.Args,
		time.Since(c.Time).Seconds(),
	)
}

// Flush stdout (no-op for this profiler)
func (testLogProfiler) Flush() {}
