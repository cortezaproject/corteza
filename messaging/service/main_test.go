package service

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
	"github.com/titpetric/factory"

	systemMigrate "github.com/crusttech/crust/system/db"
	systemService "github.com/crusttech/crust/system/service"
)

type mockDB struct{}

func (mockDB) Transaction(callback func() error) error { return callback() }

func TestMain(m *testing.M) {
	// @todo this is a very optimistic initialization, make it more robust
	godotenv.Load("../../.env")

	prefix := "messaging"
	dsn := ""

	p := func(s string) string {
		return prefix + "-" + s
	}

	flag.StringVar(&dsn, p("db-dsn"), "crust:crust@tcp(db1:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
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

	os.Exit(m.Run())
}

func assert(t *testing.T, ok bool, format string, args ...interface{}) bool {
	if !ok {
		_, file, line, _ := runtime.Caller(1)
		caller := fmt.Sprintf("\nAsserted at:%s:%d", file, line)

		t.Fatalf(format+caller, args...)
	}
	return ok
}
