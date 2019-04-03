// +build integration external

package service

import (
	"log"
	"os"
	"testing"

	"github.com/namsral/flag"
	"github.com/titpetric/factory"

	messagingMigrate "github.com/crusttech/crust/messaging/db"
	"github.com/crusttech/crust/messaging/internal/repository"
	systemMigrate "github.com/crusttech/crust/system/db"
	systemService "github.com/crusttech/crust/system/service"
)

type mockDB struct{}

func (mockDB) Transaction(callback func() error) error { return callback() }

func TestMain(m *testing.M) {
	dsn := ""
	new(repository.Flags).Init("messaging")
	flag.StringVar(&dsn, "db-dsn", "crust:crust@tcp(crust-db:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.Parse()

	factory.Database.Add("default", dsn)
	factory.Database.Add("messaging", dsn)
	factory.Database.Add("system", dsn)

	db := factory.Database.MustGet()
	db.Profiler = &factory.Database.ProfilerStdout

	// migrate database schema
	if err := systemMigrate.Migrate(db); err != nil {
		log.Printf("Error running migrations: %+v\n", err)
		return
	}
	if err := messagingMigrate.Migrate(db); err != nil {
		log.Printf("Error running migrations: %+v\n", err)
		return
	}

	// clean up tables
	{
		for _, name := range []string{"sys_user", "sys_role", "sys_role_member", "sys_organisation", "sys_rules"} {
			_, err := db.Exec("truncate " + name)
			if err != nil {
				panic("Error when clearing " + name + ": " + err.Error())
			}
		}
	}

	systemService.Init()
	Init()

	os.Exit(m.Run())
}
