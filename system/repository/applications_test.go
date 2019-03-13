package repository

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/types"

	. "github.com/crusttech/crust/internal/test"
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
	Error(t, db.Transaction(func() error {
		db.Delete("sys_application", "1=1")

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
		NoError(t, err, "Application.Create error: %+v", err)
		Assert(t, app.Valid(), "Expecting application to be valid after creation")
		Assert(t, app.Name == "created", "Expecting application name to be set, got %q", app.Name)
		Assert(t, app.Enabled, "Expecting application to be enabled")
		Assert(t, app.Unify.Name == "created", "Expecting application name to be set in unify, got %q", app.Name)
		Assert(t, app.Unify.Listed, "Expecting application to be listed in unify")
		Assert(t, app.Unify.Order == 1, "Expecting application name to have order val 1")

		app.Name = "updated"
		app.Enabled = false
		app.Unify.Name = "updated"
		app.Unify.Listed = false
		app, err = crepo.Update(app)

		NoError(t, err, "Application.Create error: %+v", err)
		Assert(t, err == nil, "Application.Create error: %+v", err)
		Assert(t, app.Name == "updated", "Expecting application name to be updated")
		Assert(t, !app.Enabled, "Expecting application to be disabled")
		Assert(t, app.Unify.Name == "updated", "Expecting application name to be updated in unify")
		Assert(t, !app.Unify.Listed, "Expecting application to be unlisted in unify")

		return errors.New("Rollback")
	}), "expected rollback error")

}
