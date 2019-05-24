// +build migrations

package db

import (
	"os"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/titpetric/factory"
)

func TestMigrations(t *testing.T) {
	factory.Database.Add("compose", os.Getenv("COMPOSE_DB_DSN"))
	db := factory.Database.MustGet("compose")
	if err := Migrate(db, logger.Default()); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
}
