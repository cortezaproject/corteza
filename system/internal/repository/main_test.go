// +build integration

package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	systemMigrate "github.com/cortezaproject/corteza-server/system/db"
)

func TestMain(m *testing.M) {
	factory.Database.Add("system", os.Getenv("SYSTEM_DB_DSN"))
	db := factory.Database.MustGet("system")
	db.Profiler = &factory.Database.ProfilerStdout

	// migrate database schema
	if err := systemMigrate.Migrate(db, logger.Default()); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
