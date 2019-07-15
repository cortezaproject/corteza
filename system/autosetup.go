package system

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/internal/repository"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

func accessControlSetup(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	c.InitServices(ctx, c)

	// Calling grant directly on internal permissions service to avoid AC check for "grant"
	var p = service.DefaultPermissions
	var ac = service.DefaultAccessControl
	return p.Grant(ctx, ac.Whitelist(), ac.DefaultRules()...)
}

func makeDefaultApplications(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	db, err := factory.Database.Get(system)
	if err != nil {
		return err
	}

	repo := repository.Application(ctx, db)

	aa, err := repo.Find()
	if err != nil {
		return err
	}

	// List of apps to create.
	//
	// We use Unify.Url field for matching,
	// so make sure it's always present!
	defApps := types.ApplicationSet{
		&types.Application{
			Name:    "CRM",
			Enabled: true,
			Unify: &types.ApplicationUnify{
				Name:   "CRM",
				Listed: true,
				Icon:   "/applications/crust_favicon.png",
				Logo:   "/applications/crust.jpg",
				Url:    "/compose/ns/crm/pages",
			},
		},
	}

	return defApps.Walk(func(defApp *types.Application) error {
		for _, a := range aa {
			if a.Unify != nil && a.Unify.Url == defApp.Unify.Url {
				// App already added.
				return nil
			}
		}

		defApp, err = repo.Create(defApp)
		c.Log.Info(
			"creating default application",
			zap.String("name", defApp.Name),
			zap.Uint64("name", defApp.ID),
			zap.Error(err),
		)
		return err

		return nil
	})
}

func discoverSettings(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	return service.DefaultSettings.With(ctx).AutoDiscovery()
}
