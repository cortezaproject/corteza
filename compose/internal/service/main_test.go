// +build integration

package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/titpetric/factory"
	"go.uber.org/zap"

	composeMigrate "github.com/cortezaproject/corteza-server/compose/db"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/test"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	dbLogger "github.com/titpetric/factory/logger"
)

type (
	mockDB struct{}
)

func (mockDB) Transaction(callback func() error) error { return callback() }

func TestMain(m *testing.M) {
	logger.SetDefault(logger.MakeDebugLogger())

	factory.Database.Add("compose", os.Getenv("COMPOSE_DB_DSN"))
	db := factory.Database.MustGet("compose")
	db.SetLogger(dbLogger.Default{})

	// migrate database schema
	if err := composeMigrate.Migrate(db, logger.Default()); err != nil {
		fmt.Printf("Error running migrations: %+v\n", err)
		return
	}

	// clean up tables
	{
		// @todo remove this asap, service should not access db at all.
		for _, name := range []string{
			"compose_chart",
			"compose_trigger",
			"compose_module_field",
			"compose_module",
			"compose_record_value",
			"compose_record",
			"compose_page",
			"compose_attachment",
			"compose_namespace",
		} {
			_, err := db.Exec("DELETE FROM " + name)
			if err != nil {
				panic("Error when clearing " + name + ": " + err.Error())
			}
		}
	}

	ctx := context.Background()

	Init(ctx, zap.NewNop(), "/tmp/corteza-compose-store")

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
