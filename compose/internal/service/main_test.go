// +build integration

package service

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/namsral/flag"
	"github.com/titpetric/factory"

	composeMigrate "github.com/crusttech/crust/compose/db"
	"github.com/crusttech/crust/compose/types"
	"github.com/crusttech/crust/internal/test"
	systemMigrate "github.com/crusttech/crust/system/db"
	systemService "github.com/crusttech/crust/system/service"
)

type mockDB struct{}

func (mockDB) Transaction(callback func() error) error { return callback() }

func TestMain(m *testing.M) {
	dsn := ""
	flag.StringVar(&dsn, "db-dsn", "crust:crust@tcp(crust-db:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.Parse()

	factory.Database.Add("default", dsn)
	factory.Database.Add("compose", dsn)
	factory.Database.Add("system", dsn)

	db := factory.Database.MustGet()
	db.Profiler = &factory.Database.ProfilerStdout

	// migrate database schema
	if err := systemMigrate.Migrate(db); err != nil {
		log.Printf("Error running migrations: %+v\n", err)
		return
	}
	if err := composeMigrate.Migrate(db); err != nil {
		log.Printf("Error running migrations: %+v\n", err)
		return
	}

	// clean up tables
	{
		for _, name := range []string{"compose_chart", "compose_trigger", "compose_module", "compose_module_form", "compose_record", "compose_record_value", "compose_page", "sys_user"} {
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

func createTestNamespaces(ctx context.Context, t *testing.T) (ns1 *types.Namespace, ns2 *types.Namespace) {
	var err error

	ns1, err = Namespace().With(ctx).Create(&types.Namespace{Enabled: true, Name: "TestNamespace"})
	test.Assert(t, err == nil, "Error when creating namespace: %+v", err)

	ns2, err = Namespace().With(ctx).Create(&types.Namespace{Enabled: true, Name: "TestNamespace"})
	test.Assert(t, err == nil, "Error when creating namespace: %+v", err)

	return ns1, ns2
}
