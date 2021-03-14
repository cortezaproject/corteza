package provision

import (
	"context"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"time"
)

func migrateApplications(ctx context.Context, s store.Storer) error {
	rd := releaseDate(2021, time.March)

	set, _, err := store.SearchApplications(ctx, s, types.ApplicationFilter{})
	if err != nil {
		return err
	}

	return set.Walk(func(app *types.Application) (err error) {
		// Disable messaging app but only if it was not updated after 21.3 release date
		if app.Unify.Url == "/messaging" && (app.UpdatedAt == nil || app.UpdatedAt.Before(rd)) {
			app.Enabled = false
			err = store.UpdateApplication(ctx, s, app)
		}

		return err
	})
}
