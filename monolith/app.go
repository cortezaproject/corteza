package monolith

import (
	"context"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/corteza"
	"github.com/cortezaproject/corteza-server/messaging"
	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/system"
)

type (
	App struct {
		Core      *corteza.App
		System    *system.App
		Compose   *compose.App
		Messaging *messaging.App
	}
)

func (monolith *App) Setup(log *zap.Logger, opts *app.Options) (err error) {

	// Make sure system behaves properly
	//
	// This will alter the auth settings provision procedure
	system.IsMonolith = true

	return app.RunSetup(
		log,
		opts,
		monolith.System,
		monolith.Compose,
		monolith.Messaging,
	)
}

func (monolith *App) Upgrade(ctx context.Context) (err error) {
	return app.RunUpgrade(
		ctx,
		monolith.System,
		monolith.Compose,
		monolith.Messaging,
	)
}

func (monolith App) Initialize(ctx context.Context) (err error) {
	return app.RunInitialize(
		ctx,
		monolith.System,
		monolith.Compose,
		monolith.Messaging,
	)
}

func (monolith *App) Activate(ctx context.Context) (err error) {
	err = app.RunActivate(
		ctx,
		monolith.System,
		monolith.Compose,
		monolith.Messaging,
	)

	if err != nil {
		return
	}

	// Override JWT maker for Corredor with internal
	corredor.Service().SetJwtMaker(corredor.InternalAuthTokenMaker())

	return
}

func (monolith *App) Provision(ctx context.Context) (err error) {
	return app.RunProvision(
		ctx,
		monolith.System,
		monolith.Compose,
		monolith.Messaging,
	)
}

func (monolith *App) MountApiRoutes(r chi.Router) {
	r.Route("/system", func(r chi.Router) {
		monolith.System.MountApiRoutes(r)
	})

	r.Route("/compose", func(r chi.Router) {
		monolith.Compose.MountApiRoutes(r)
	})

	r.Route("/messaging", func(r chi.Router) {
		monolith.Messaging.MountApiRoutes(r)
	})
}

func (monolith *App) RegisterGrpcServices(srv *grpc.Server) {
	monolith.System.RegisterGrpcServices(srv)
}

// RegisterCliCommands on monolith will wrapp all commands
func (monolith *App) RegisterCliCommands(rootCmd *cobra.Command) {
	systemCmd := &cobra.Command{
		Use:     "system",
		Aliases: []string{"sys"},
		Short:   "Commands from system service",
	}

	composeCmd := &cobra.Command{
		Use:     "compose",
		Aliases: []string{"cmp"},
		Short:   "Commands from messaging service",
	}

	messagingCmd := &cobra.Command{
		Use:     "messaging",
		Aliases: []string{"msg"},
		Short:   "Commands from compose service",
	}

	monolith.System.RegisterCliCommands(systemCmd)
	monolith.Compose.RegisterCliCommands(composeCmd)
	monolith.Messaging.RegisterCliCommands(messagingCmd)

	rootCmd.AddCommand(
		systemCmd,
		composeCmd,
		messagingCmd,
	)
}
