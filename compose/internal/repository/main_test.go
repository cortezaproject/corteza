// +build integration

package repository

import (
	"os"
	"testing"
	"time"

	"github.com/namsral/flag"
	"github.com/titpetric/factory"
	"go.uber.org/zap/zapcore"

	"github.com/crusttech/crust/internal/logger"
)

func TestMain(m *testing.M) {
	logger.Init(zapcore.DebugLevel)

	dsn := ""
	flag.StringVar(&dsn, "compose-db-dsn", "", "")
	flag.Parse()
	factory.Database.Add("compose", dsn)

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
