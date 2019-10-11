package messaging

import (
	"context"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/commands"
	migrate "github.com/cortezaproject/corteza-server/messaging/db"
	"github.com/cortezaproject/corteza-server/messaging/rest"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/websocket"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
)

const (
	messaging = "messaging"
)

func Configure() *cli.Config {
	var (
		servicesInitialized bool

		// Websocket handler
		ws *websocket.Websocket
	)

	return &cli.Config{
		ServiceName: messaging,

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
				Storage: *c.StorageOpt,
			}))
		},

		ApiServerPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				if c.ProvisionOpt.MigrateDatabase {
					cli.HandleError(c.ProvisionMigrateDatabase.Run(ctx, cmd, c))
				}

				if c.ProvisionOpt.Configuration {
					cli.HandleError(provisionConfig(ctx, cmd, c))
				}

				c.InitServices(ctx, c)

				var websocketOpt = options.Websocket(messaging)

				ws = websocket.Init(ctx, &websocket.Config{
					Timeout:     websocketOpt.Timeout,
					PingTimeout: websocketOpt.PingTimeout,
					PingPeriod:  websocketOpt.PingPeriod,
				})

				go service.Watchers(ctx)
				return nil
			},
		},

		ApiServerRoutes: cli.Mounters{
			rest.MountRoutes,
			// Wrap in func() to assure ws is set when mounted
			func(r chi.Router) { ws.ApiServerRoutes(r) },
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
				var db, err = factory.Database.Get(messaging)
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
