package system

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	provision "github.com/cortezaproject/corteza-server/provision/system"
	"github.com/cortezaproject/corteza-server/system/importer"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

func provisionConfig(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	c.Log.Debug("running configuration provision")
	c.InitServices(ctx, c)

	// Make sure we have all full access for provisioning
	ctx = auth.SetSuperUserContext(ctx)

	if provisioned, err := isProvisioned(ctx); err != nil {
		return err
	} else if provisioned {
		c.Log.Debug("configuration already provisioned")
		return nil
	}

	readers, err := impAux.ReadStatic(provision.Asset)
	if err != nil {
		return err
	}

	return errors.Wrap(
		importer.Import(ctx, readers...),
		"could not provision configuration for system service",
	)
}

// Provision ONLY when there are no rules for role admins / everyone
func isProvisioned(ctx context.Context) (bool, error) {
	return len(service.DefaultPermissions.FindRulesByRoleID(permissions.EveryoneRoleID)) > 0 &&
		len(service.DefaultPermissions.FindRulesByRoleID(permissions.AdminsRoleID)) > 0, nil
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
