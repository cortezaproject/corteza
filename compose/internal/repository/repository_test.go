// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestRepository(t *testing.T) {
	repo := &repository{}
	repo.With(context.Background(), nil)
}

func tx(t *testing.T, f func(context.Context, *factory.DB, *types.Namespace) error) {
	var (
		err error
		ctx = context.Background()
		db  = DB(ctx)
		ns  *types.Namespace
	)

	err = db.Begin()
	test.Assert(t, err == nil, "Could not begin transaction: %+v", err)

	ns, err = Namespace(ctx, db).Create(&types.Namespace{})
	test.Assert(t, err == nil, "Test transaction setup (namespace creation) resulted in an error: %+v", err)

	// Setup test log profiler and route all db logs to test log facility
	// We do this right after transaction is started and setup is done, no need for extra queries there...
	db.Profiler = newTestLogProfiler(t)

	err = f(ctx, db, ns)
	test.Assert(t, err == nil, "Test transaction resulted in an error: %+v", err)

	// Remove profiler to omit final rollback statement in the logs
	db.Profiler = nil

	err = db.Rollback()
	test.Assert(t, err == nil, "Could not rollback transaction: %+v", err)
}
