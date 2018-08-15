package repository

import (
	"github.com/namsral/flag"
	"github.com/joho/godotenv"
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
	factory.Database.MustGet().Profiler = &factory.Database.ProfilerStdout

	os.Exit(m.Run())
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}
