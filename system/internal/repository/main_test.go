// +build integration

package repository

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
	factory.Database.Add("system", os.Getenv("SYSTEM_DB_DSN"))
	db := factory.Database.MustGet("system")
	db.SetLogger(dbLogger.Default{})

	// migrate database schema
	if err := systemMigrate.Migrate(db, logger.Default()); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
