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
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
)

const (
	compose = "compose"
)

func Configure() *cli.Config {
	return &cli.Config{
		ServiceName: compose,

		RootCommandPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				if c.ProvisionOpt.MigrateDatabase {
					cli.HandleError(c.ProvisionMigrateDatabase.Run(ctx, cmd, c))
				}

				storagePath := options.EnvString("", "COMPOSE_STORAGE_PATH", "var/store")

				cli.HandleError(service.Init(ctx, c.Log, storagePath))

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
