package compose

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/compose/commands"
	migrate "github.com/cortezaproject/corteza-server/compose/db"
	"github.com/cortezaproject/corteza-server/compose/rest"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

const (
	compose = "compose"
)

func Configure() *cli.Config {
	var servicesInitialized bool

	return &cli.Config{
		ServiceName: compose,

		RootCommandPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				return
			},
		},

		InitServices: func(ctx context.Context, c *cli.Config) {
			if servicesInitialized {
				return
			}
			servicesInitialized = true

			cli.HandleError(service.Init(ctx, c.Log, service.Config{
				Storage:          *c.StorageOpt,
				Corredor:         *c.ScriptRunner,
				GRPCClientSystem: *c.GRPCServerSystem,
			}))
		},

		ApiServerPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				if c.ProvisionOpt.MigrateDatabase {
					cli.HandleError(c.ProvisionMigrateDatabase.Run(ctx, cmd, c))
				}

				c.InitServices(ctx, c)

				if c.ProvisionOpt.Configuration {
					cli.HandleError(provisionConfig(ctx, cmd, c))
				}

				go service.Watchers(ctx)
				return nil
			},
		},

		ApiServerRoutes: cli.Mounters{
			rest.MountRoutes,
		},

		AdtSubCommands: cli.CommandMakers{
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Importer(ctx, c)
			},
			func(ctx context.Context, c *cli.Config) *cobra.Command {
				return commands.Exporter(ctx, c)
			},
		},

		ProvisionMigrateDatabase: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				var db, err = factory.Database.Get(compose)
				if err != nil {
					return err
				}

				db = db.With(ctx).Quiet()

				return migrate.Migrate(db, c.Log)
			},
		},

		ProvisionConfig: cli.Runners{
			provisionConfig,
		},
	}
}
