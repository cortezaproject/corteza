package provision

import (
	"context"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"time"
)

func migrateApplications(ctx context.Context, s store.Storer) error {
	rd := releaseDate(2021, time.March)

	set, _, err := store.SearchApplications(ctx, s, types.ApplicationFilter{})
	if err != nil {
		return err
	}

	return set.Walk(func(app *types.Application) (err error) {
		if app.Unify == nil {
			return
		}

		if app.Unify.Url == "/messaging" && (app.UpdatedAt == nil || app.UpdatedAt.Before(rd)) {
			// Disable and un-list messaging app
			// but only if it was not updated after 21.3 release date
			app.Unify.Listed = false
			app.Enabled = false

			return store.UpdateApplication(ctx, s, app)
		}

		return nil
	})
}
