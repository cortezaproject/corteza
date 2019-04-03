// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/system/types"
)

func TestApplication(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	db := factory.Database.MustGet()

	// Create application repository.
	crepo := Application(context.Background(), db)

	// Run tests in transaction to maintain DB state.
	test.Error(t, db.Transaction(func() error {
		db.Exec("DELETE FROM sys_application WHERE 1=1")

		app := &types.Application{
			Name:    "created",
			Enabled: true,
			OwnerID: 1,
			Unify: &types.ApplicationUnify{
				Name:   "created",
				Listed: true,
				Order:  1,
				Icon:   "...ico",
			},
		}

		app, err := crepo.Create(app)
		test.NoError(t, err, "Application.Create error: %+v", err)
		test.Assert(t, app.Valid(), "Expecting application to be valid after creation")
		test.Assert(t, app.Name == "created", "Expecting application name to be set, got %q", app.Name)
		test.Assert(t, app.Enabled, "Expecting application to be enabled")
		test.Assert(t, app.Unify.Name == "created", "Expecting application name to be set in unify, got %q", app.Name)
		test.Assert(t, app.Unify.Listed, "Expecting application to be listed in unify")
		test.Assert(t, app.Unify.Order == 1, "Expecting application name to have order val 1")

		app.Name = "updated"
		app.Enabled = false
		app.Unify.Name = "updated"
		app.Unify.Listed = false
		app, err = crepo.Update(app)

		test.NoError(t, err, "Application.Create error: %+v", err)
		test.Assert(t, err == nil, "Application.Create error: %+v", err)
		test.Assert(t, app.Name == "updated", "Expecting application name to be updated")
		test.Assert(t, !app.Enabled, "Expecting application to be disabled")
		test.Assert(t, app.Unify.Name == "updated", "Expecting application name to be updated in unify")
		test.Assert(t, !app.Unify.Listed, "Expecting application to be unlisted in unify")

		return errors.New("Rollback")
	}), "expected rollback error")

}
