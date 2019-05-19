// +build integration

package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/namsral/flag"
	"github.com/titpetric/factory"

	migrate "github.com/crusttech/crust/messaging/db"
)

func TestMain(m *testing.M) {
	dsn := ""
	new(Flags).Init("messaging")
	flag.StringVar(&dsn, "db-dsn", "crust:crust@tcp(crust-db:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.Parse()

	factory.Database.Add("default", dsn)
	factory.Database.Add("messaging", dsn)

	db := factory.Database.MustGet()
	db.Profiler = &factory.Database.ProfilerStdout

	// migrate database schema
	if err := migrate.Migrate(db); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		return
	}

	os.Exit(m.Run())
}
