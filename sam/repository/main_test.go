package repository

import (
	"flag"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/titpetric/factory"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		return
	}

	// @todo this is a very optimistic initialization, make it more robust
	godotenv.Load("../../.env")
	factory.Database.Add("default", os.Getenv("SAM_DB_DSN"))
	factory.Database.MustGet().Profiler = &factory.Database.ProfilerStdout

	os.Exit(m.Run())
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}
