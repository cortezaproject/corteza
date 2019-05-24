// +build integration

package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/titpetric/factory"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	systemMigrate "github.com/cortezaproject/corteza-server/system/db"
)

func TestMain(m *testing.M) {
	logger.Init(zapcore.DebugLevel)

	factory.Database.Add("system", os.Getenv("SYSTEM_DB_DSN"))
	db := factory.Database.MustGet("system")
	db.Profiler = &factory.Database.ProfilerStdout

	// migrate database schema
	if err := systemMigrate.Migrate(db); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		return
	}

	os.Exit(m.Run())
}
