package messaging

import (
	"context"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	migrate "github.com/cortezaproject/corteza-server/messaging/db"
	"github.com/cortezaproject/corteza-server/messaging/internal/service"
	"github.com/cortezaproject/corteza-server/messaging/rest"
	"github.com/cortezaproject/corteza-server/messaging/websocket"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
)

const (
	messaging = "messaging"
)

func Configure() *cli.Config {
	var (
		// Websocket handler
		ws *websocket.Websocket
	)

	return &cli.Config{
		ServiceName: messaging,

		RootCommandPreRun: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) (err error) {
				if c.ProvisionOpt.MigrateDatabase {
					cli.HandleError(c.ProvisionMigrateDatabase.Run(ctx, cmd, c))
				}

				storagePath := options.EnvString("", "MESSAGING_STORAGE_PATH", "var/store")

				cli.HandleError(service.Init(ctx, c.Log, storagePath))

				var websocketOpt = options.Websocket(messaging)

				ws = websocket.Init(ctx, &websocket.Config{
					Timeout:     websocketOpt.Timeout,
					PingTimeout: websocketOpt.PingTimeout,
					PingPeriod:  websocketOpt.PingPeriod,
				})

				if c.ProvisionOpt.AutoSetup {
					cli.HandleError(accessControlSetup(ctx, cmd, c))
					cli.HandleError(makeDefaultChannels(ctx, cmd, c))
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
			// Wrap in func() to assure ws is set when mounted
			func(r chi.Router) { ws.ApiServerRoutes(r) },
		},

		ProvisionMigrateDatabase: cli.Runners{
			func(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
				var db, err = factory.Database.Get(messaging)
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
