package app

import (
	"context"

	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type (
	httpApiServer interface {
		MountRoutes(mm ...func(chi.Router))
		Serve(ctx context.Context)
	}

	grpcServer interface {
		RegisterServices(func(server *grpc.Server))
		Serve(ctx context.Context)
	}

	wsServer interface {
		MountRoutes(chi.Router)
		Send(kind string, payload interface{}, userIDs ...uint64) error
	}

	authServicer interface {
		MountHttpRoutes(string, chi.Router)
		UpdateSettings(*settings.Settings)
		Watch(ctx context.Context)
	}

	CortezaApp struct {
		Opt *options.Options
		lvl int
		Log *zap.Logger

		// Store interface
		//
		// Just a blank interface{} because we want to avoid generating
		// whole store interface (as we do for other packages).
		//
		// Value will be type-casted when assigned to sys/msg/cmp services
		// with warnings when incompatible
		Store store.Storer

		// CLI Commands
		Command *cobra.Command

		// Servers
		HttpServer httpApiServer
		GrpcServer grpcServer
		WsServer   wsServer

		AuthService authServicer
	}
)

func New() *CortezaApp {
	app := &CortezaApp{lvl: bootLevelWaiting}
	app.Opt = options.Init()
	app.InitCLI()
	return app
}

func (app *CortezaApp) Options() *options.Options {
	return app.Opt
}
