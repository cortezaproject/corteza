// +build integration

package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/titpetric/factory"
	"go.uber.org/zap"

	messagingMigrate "github.com/cortezaproject/corteza-server/messaging/db"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	dbLogger "github.com/titpetric/factory/logger"
)

type mockDB struct{}

func (mockDB) Transaction(callback func() error) error { return callback() }

func TestMain(m *testing.M) {
	logger.SetDefault(logger.MakeDebugLogger())

	factory.Database.Add("messaging", os.Getenv("MESSAGING_DB_DSN"))
	db := factory.Database.MustGet("messaging")
	db.SetLogger(dbLogger.Default{})

	// migrate database schema
	if err := messagingMigrate.Migrate(db, logger.Default()); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		return
	}

	Init(context.Background(), zap.NewNop(), Config{
		Storage: options.StorageOpt{
			Path: "/tmp/corteza-messaging-store",
		},
	})

	os.Exit(m.Run())
}
