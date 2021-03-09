package provision

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

func apps(ctx context.Context, s store.Storer) error {
	var (
		v202103 = time.Date(2021, 4, 31, 0, 0, 0, 0, time.UTC)
	)

	set, _, err := store.SearchApplications(ctx, s, types.ApplicationFilter{})
	if err != nil {
		return err
	}

	return set.Walk(func(app *types.Application) (err error) {
		if app.Unify.Url == "/messaging" && (app.UpdatedAt == nil || app.UpdatedAt.Before(v202103)) {
			// Disable messaging app but only if it was not updated after 21.3 release date
			app.Enabled = false
			err = store.UpdateApplication(ctx, s, app)
		}

		return err
	})
}
