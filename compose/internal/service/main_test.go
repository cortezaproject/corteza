// +build integration

package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/titpetric/factory"
	"go.uber.org/zap/zapcore"

	composeMigrate "github.com/cortezaproject/corteza-server/compose/db"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/logger"
	"github.com/cortezaproject/corteza-server/internal/test"
)

type (
	mockDB struct{}
)

func (mockDB) Transaction(callback func() error) error { return callback() }

func TestMain(m *testing.M) {
	logger.Init(zapcore.DebugLevel)

	factory.Database.Add("compose", os.Getenv("COMPOSE_DB_DSN"))
	db := factory.Database.MustGet("compose")
	db.Profiler = &factory.DatabaseProfilerStdout{}

	// migrate database schema
	if err := composeMigrate.Migrate(db); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		return
	}

	// clean up tables
	{
		// @todo remove this asap, service should not access db at all.
		for _, name := range []string{
			"compose_chart",
			"compose_trigger",
			"compose_module_field",
			"compose_module",
			"compose_record_value",
			"compose_record",
			"compose_page",
			"compose_attachment",
			"compose_namespace",
		} {
			_, err := db.Exec("DELETE FROM " + name)
			if err != nil {
				panic("Error when clearing " + name + ": " + err.Error())
			}
		}
	}

	ctx := context.Background()

	Init(ctx)

	os.Exit(m.Run())
}

func createTestNamespaces(ctx context.Context, t *testing.T) (ns1 *types.Namespace, ns2 *types.Namespace) {
	var err error

	ns1, err = Namespace().With(ctx).Create(&types.Namespace{Enabled: true, Name: "TestNamespace"})
	test.Assert(t, err == nil, "Error when creating namespace: %+v", err)

	ns2, err = Namespace().With(ctx).Create(&types.Namespace{Enabled: true, Name: "TestNamespace"})
	test.Assert(t, err == nil, "Error when creating namespace: %+v", err)

	return ns1, ns2
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
