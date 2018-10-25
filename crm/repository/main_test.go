package repository

import (
	"github.com/joho/godotenv"
	"github.com/namsral/flag"
	"github.com/titpetric/factory"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// @todo this is a very optimistic initialization, make it more robust
	godotenv.Load("../../.env")

	prefix := "sam"
	dsn := ""

	p := func(s string) string {
		return prefix + "-" + s
	}

	flag.StringVar(&dsn, p("db-dsn"), "crust:crust@tcp(db1:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.Parse()

	if testing.Short() {
		return
	}

	factory.Database.Add("default", dsn)

	db := factory.Database.MustGet()
	db.Profiler = &factory.Database.ProfilerStdout

	// clean up tables
	{
		for _, name := range []string{"crm_module", "crm_module_form", "crm_content", "crm_content_column", "crm_page"} {
			_, err := db.Exec("truncate " + name)
			if err != nil {
				panic("Error when clearing " + name + ": " + err.Error())
			}
		}
	}

	os.Exit(m.Run())
}

func db() *factory.DB {
	return factory.Database.MustGet()
}

func must(t *testing.T, err error, message ...string) {
	prefix := "Error"
	if len(message) > 0 {
		prefix = message[0]
	}
	if err != nil {
		t.Fatalf(prefix+": %+v", err)
	}
}

func assert(t *testing.T, ok bool, format string, args ...interface{}) bool {
	if !ok {
		t.Fatalf(format, args...)
	}
	return ok
}
