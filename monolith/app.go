package monolith

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/webapp"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"strings"

	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/corteza"
	"github.com/cortezaproject/corteza-server/messaging"
	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/system"
	systemService "github.com/cortezaproject/corteza-server/system/service"
)

type (
	App struct {
		Opts      *app.Options
		Core      *corteza.App
		System    *system.App
		Compose   *compose.App
		Messaging *messaging.App
	}
)

func (monolith *App) Setup(log *zap.Logger, opts *app.Options) (err error) {
	monolith.Opts = opts
	// Make sure system behaves properly
	//
	// This will alter the auth settings provision procedure
	system.IsMonolith = true

	if err = monolith.CheckOptions(); err != nil {
		return
	}

	err = app.RunSetup(
		log,
		opts,
		monolith.System,
		monolith.Compose,
		monolith.Messaging,
	)

	if err != nil {
		return
	}

	return
}

func (monolith *App) CheckOptions() error {
	o := monolith.Opts.HTTPServer
	if o.ApiEnabled && o.WebappEnabled && o.ApiBaseUrl == o.WebappBaseUrl {
		return errors.Errorf("cannot serve api and web apps form the same base url (%v)", o.ApiBaseUrl)
	}

	return nil
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
	err = app.RunInitialize(
		ctx,
		monolith.System,
		monolith.Compose,
		monolith.Messaging,
	)

	if err != nil {
		return
	}

	corredor.Service().SetUserFinder(systemService.DefaultUser)
	corredor.Service().SetRoleFinder(systemService.DefaultRole)

	return
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
	var (
		apiEnabled = monolith.Opts.HTTPServer.ApiEnabled
		apiBaseUrl = strings.Trim(monolith.Opts.HTTPServer.ApiBaseUrl, "/")

		webappEnabled = monolith.Opts.HTTPServer.WebappEnabled
		webappBaseUrl = strings.Trim(monolith.Opts.HTTPServer.WebappBaseUrl, "/")
	)

	if apiEnabled {
		r.Route("/"+apiBaseUrl, func(r chi.Router) {
			r.Route("/system", func(r chi.Router) {
				monolith.System.MountApiRoutes(r)
			})

			r.Route("/compose", func(r chi.Router) {
				monolith.Compose.MountApiRoutes(r)
			})

			r.Route("/messaging", func(r chi.Router) {
				monolith.Messaging.MountApiRoutes(r)
			})
		})
	}

	if webappEnabled {
		r.Route("/"+webappBaseUrl, webapp.MakeWebappServer(monolith.Opts.HTTPServer))
	}
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
