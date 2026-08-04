package app

import (
	"context"
	"net/http"

	"github.com/crusttech/human/server/auth/settings"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/llm"
	mcpkg "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/system/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type (
	httpApiServer interface {
		Serve(ctx context.Context)
		Activate(mm ...func(chi.Router))
		Shutdown()
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
		WellKnownOpenIDConfiguration() http.HandlerFunc
		UpdateSettings(*settings.Settings)
		Watch(ctx context.Context)
	}

	apigwServicer interface {
		http.Handler
	}

	HumanApp struct {
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

		// separate pool/DB when ACTIONLOG_DB_DSN is set; otherwise equals Store
		ActionlogStore store.Storer

		// CLI Commands
		Command *cobra.Command

		oa2m oauth2.Manager

		DefaultAuthClient *types.AuthClient

		// Servers
		HttpServer httpApiServer
		GrpcServer grpcServer
		WsServer   wsServer
		McpServer  *mcpkg.MCPServer

		LlmService   *llm.Service
		AuthService  authServicer
		ApigwService apigwServicer

		systemEntitiesInitialized bool
	}
)

func New() *HumanApp {
	app := &HumanApp{
		lvl: bootLevelWaiting,
		Log: logger.Default(),
	}

	app.InitCLI()
	return app
}

func (app *HumanApp) Options() *options.Options {
	return app.Opt
}
