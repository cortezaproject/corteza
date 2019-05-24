// +build integration

package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/titpetric/factory"

	messagingMigrate "github.com/cortezaproject/corteza-server/messaging/db"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type mockDB struct{}

func (mockDB) Transaction(callback func() error) error { return callback() }

func TestMain(m *testing.M) {
	logger.SetDefault(logger.MakeDebugLogger())

	factory.Database.Add("messaging", os.Getenv("MESSAGING_DB_DSN"))
	db := factory.Database.MustGet("messaging")
	db.Profiler = &factory.Database.ProfilerStdout

	// migrate database schema
	if err := messagingMigrate.Migrate(db, logger.Default()); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		return
	}

	Init(context.Background())

	os.Exit(m.Run())
}
