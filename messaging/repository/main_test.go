// +build integration

package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/titpetric/factory"

	migrate "github.com/cortezaproject/corteza-server/messaging/db"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	dbLogger "github.com/titpetric/factory/logger"
)

func TestMain(m *testing.M) {
	factory.Database.Add("messaging", os.Getenv("MESSAGING_DB_DSN"))
	db := factory.Database.MustGet("messaging")
	db.SetLogger(dbLogger.Default{})

	// migrate database schema
	if err := migrate.Migrate(db, logger.Default()); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		return
	}

	os.Exit(m.Run())
}
