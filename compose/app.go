package compose

import (
	"context"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/commands"
	migrate "github.com/cortezaproject/corteza-server/compose/db"
	"github.com/cortezaproject/corteza-server/compose/rest"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/scheduler"
	"github.com/cortezaproject/corteza-server/system/auth/external"
)

type (
	App struct {
		Opts *app.Options
		Log  *zap.Logger
	}
)

const SERVICE = "compose"

func (app *App) Setup(log *zap.Logger, opts *app.Options) (err error) {
	app.Log = log.Named(SERVICE)
	app.Opts = opts

	scheduler.OnTick(
		event.ComposeOnInterval(),
		event.ComposeOnTimestamp(),
	)

	return
}

func (app *App) Upgrade(ctx context.Context) (err error) {
	db := factory.Database.MustGet(SERVICE, "default").With(ctx).Quiet()
	err = migrate.Migrate(db, app.Log)
	if err != nil {
		return
	}

	return
}

// Initialized
func (app *App) Initialize(ctx context.Context) (err error) {
	// Connects to all services it needs to
	err = service.Initialize(ctx, app.Log, service.Config{
		Storage:  app.Opts.Storage,
		Corredor: app.Opts.Corredor,
	})

	if err != nil {
		return
	}

	// Initialize external authentication (from default settings)
	external.Init()
	return
}

func (app *App) Activate(ctx context.Context) (err error) {
	if err = service.Activate(ctx); err != nil {
		return
	}

	service.Watchers(ctx)

	// Wire in cross service JWT maker for Corredor
	corredor.Service().SetJwtMaker(corredor.CrossServiceAuthTokenMaker())

	return
}

func (app *App) Provision(ctx context.Context) (err error) {
	ctx = auth.SetSuperUserContext(ctx)

	if err = provisionConfig(ctx, app.Log); err != nil {
		return
	}

	return
}

func (app *App) MountApiRoutes(r chi.Router) {
	rest.MountRoutes(r)
}

func (app *App) RegisterCliCommands(p *cobra.Command) {
	p.AddCommand(
		commands.Importer(),
		commands.Exporter(),
	)
}
