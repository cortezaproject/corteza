// +build migrations

package db

import (
	"os"
	"testing"

	"github.com/titpetric/factory"
)

func TestMigrations(t *testing.T) {
	factory.Database.Add("messaging", os.Getenv("MESSAGING_DB_DSN"))
	db := factory.Database.MustGet("messaging")
	db.Profiler = &factory.Database.ProfilerStdout
	if err := Migrate(db); err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}
}
