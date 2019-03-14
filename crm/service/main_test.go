package service

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
	"github.com/titpetric/factory"

	crmMigrate "github.com/crusttech/crust/crm/db"
	systemMigrate "github.com/crusttech/crust/system/db"
	systemService "github.com/crusttech/crust/system/service"
)

func TestMain(m *testing.M) {
	// @todo this is a very optimistic initialization, make it more robust
	godotenv.Load("../../.env")

	prefix := "crm"
	dsn := ""

	p := func(s string) string {
		return prefix + "-" + s
	}

	flag.StringVar(&dsn, p("db-dsn"), "crust:crust@tcp(db1:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.Parse()

	factory.Database.Add("default", dsn)
	factory.Database.Add("crm", dsn)
	factory.Database.Add("system", dsn)

	db := factory.Database.MustGet()
	db.Profiler = &factory.Database.ProfilerStdout

	// migrate database schema
	if err := systemMigrate.Migrate(db); err != nil {
		log.Printf("Error running migrations: %+v\n", err)
		return
	}
	if err := crmMigrate.Migrate(db); err != nil {
		log.Printf("Error running migrations: %+v\n", err)
		return
	}

	// clean up tables
	{
		for _, name := range []string{"crm_chart", "crm_trigger", "crm_module", "crm_module_form", "crm_record", "crm_record_value", "crm_page", "sys_user"} {
			_, err := db.Exec("truncate " + name)
			if err != nil {
				panic("Error when clearing " + name + ": " + err.Error())
			}
		}
	}

	systemService.Init()

	os.Exit(m.Run())
}

func assert(t *testing.T, ok bool, format string, args ...interface{}) bool {
	if !ok {
		t.Fatalf(format, args...)
	}
	return ok
}
