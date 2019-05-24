package compose

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	migrate "github.com/cortezaproject/corteza-server/compose/db"
	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

const (
	compose = "compose"
)

func Configure() *cli.Config {
	var (
		accessControlSetup = func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
			// Calling grant directly on internal permissions service to avoid AC check for "grant"
			var p = service.DefaultPermissions
			var ac = service.DefaultAccessControl
			return p.Grant(ctx, ac.Whitelist(), ac.DefaultRules()...)
		}
	)

	return &cli.Config{
		ServiceName: compose,

		RootCommandPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				if c.ProvisionOpt.MigrateDatabase {
					cli.HandleError(c.ProvisionMigrateDatabase.Run(ctx, cmd, c))
				}

				cli.HandleError(service.Init(ctx, c.Log))

				if c.ProvisionOpt.AutoSetup {
					cli.HandleError(accessControlSetup(ctx, cmd, c))
				}

				return
			},
		},

		ApiServerPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				go service.Watchers(ctx)
				return nil
			},
		},

		ApiServerRoutes: cli.Mounters{
			rest.MountRoutes,
		},

		ProvisionMigrateDatabase: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				var db, err = factory.Database.Get(compose)
				if err != nil {
					return err
				}

				db = db.With(ctx)
				// Disable profiler for migrations
				db.Profiler = nil

				return migrate.Migrate(db, c.Log)
			},
		},

		ProvisionAccessControl: cli.Runners{
			accessControlSetup,
		},
	}
}
