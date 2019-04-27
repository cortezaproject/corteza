// +build migrations

package db

import (
	"testing"

	"github.com/namsral/flag"
	"github.com/titpetric/factory"
)

func TestMigrations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	var dsn string

	flag.StringVar(&dsn, "db-dsn", "crust:crust@tcp(crust-db:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.Parse()

	factory.Database.Add("default", dsn)
	factory.Database.Add("system", dsn)

	db := factory.Database.MustGet()
	if err := Migrate(db); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
}
