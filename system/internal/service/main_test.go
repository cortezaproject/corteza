// +build integration

package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	systemMigrate "github.com/cortezaproject/corteza-server/system/db"
	dbLogger "github.com/titpetric/factory/logger"
)

func TestMain(m *testing.M) {
	logger.SetDefault(logger.MakeDebugLogger())

	factory.Database.Add("system", os.Getenv("SYSTEM_DB_DSN"))
	db := factory.Database.MustGet("system")
	db.SetLogger(dbLogger.Default{})

	// migrate database schema
	if err := systemMigrate.Migrate(db, logger.Default()); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		return
	}

	os.Exit(m.Run())
}
