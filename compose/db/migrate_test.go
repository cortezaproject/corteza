// +build migrations

package db

import (
	"os"
	"testing"

	"github.com/titpetric/factory"
)

func TestMigrations(t *testing.T) {
	factory.Database.Add("compose", os.Getenv("COMPOSE_DB_DSN"))
	db := factory.Database.MustGet("compose")
	if err := Migrate(db); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
}
