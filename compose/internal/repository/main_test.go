// +build integration

package repository

import (
	"os"
	"testing"

	"github.com/titpetric/factory"
	dbLogger "github.com/titpetric/factory/logger"

	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.SetDefault(logger.MakeDebugLogger())

	factory.Database.Add("compose", os.Getenv("COMPOSE_DB_DSN"))
	db := factory.Database.MustGet("compose")
	db.SetLogger(dbLogger.Default{})

	os.Exit(m.Run())
}
