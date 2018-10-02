package db

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/config"
)

func TestMigrations(t *testing.T) {
	// @todo this is a very optimistic initialization, make it more robust
	godotenv.Load("../../.env")

	dbConfig := new(config.Database).Init("sam")
	flag.Parse()

	if err := dbConfig.Validate(); err != nil {
		t.Fatalf("Error in database config: %+v", err)
	}

	factory.Database.Add("default", dbConfig.DSN)

	db := factory.Database.MustGet()
	if err := Migrate(db); err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}
}
