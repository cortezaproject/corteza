package corteza

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/app"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/db"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/mail"
	"github.com/cortezaproject/corteza-server/pkg/monitor"
	"github.com/cortezaproject/corteza-server/pkg/scheduler"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	App struct {
		opt *app.Options
		log *zap.Logger
	}
)

var _ app.Runnable = &App{}

func (app *App) Setup(log *zap.Logger, opts *app.Options) (err error) {
	app.log = log
	app.opt = opts

	logger.SetDefault(log)

	if err = sentry.Init(opts.Sentry); err != nil {
		return errors.Wrap(err, "could not initialize Sentry")
	}

	// Use Sentry right away to handle any panics
	// that might occur inside auth, mail setup...
	defer sentry.Recover()

	auth.SetupDefault(opts.Auth.Secret, int(opts.Auth.Expiry/time.Minute))
	mail.SetupDialer(opts.SMTP.Host, opts.SMTP.Port, opts.SMTP.User, opts.SMTP.Pass, opts.SMTP.From)

	http.SetupDefaults(
		opts.HTTPClient.HttpClientTimeout,
		opts.HTTPClient.ClientTSLInsecure,
	)

	monitor.Setup(app.log, opts.Monitor)

	scheduler.Setup(log, eventbus.Service(), 0)

	if err = corredor.Setup(log, opts.Corredor); err != nil {
		return err
	}

	return
}

func (app *App) Initialize(ctx context.Context) (err error) {
	defer sentry.Recover()

	_, err = db.TryToConnect(ctx, app.log, app.opt.DB)
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}

	if err = corredor.Service().Connect(ctx); err != nil {
		return
	}

	return
}

func (app *App) Upgrade(ctx context.Context) error {
	return nil
}

func (app *App) Activate(ctx context.Context) (err error) {

	// Start scheduler
	scheduler.Service().Start(ctx)

	// Load corredor scripts & init watcher (script reloader)
	corredor.Service().Load(ctx)
	corredor.Service().Watch(ctx)

	return nil
}

func (app *App) Provision(ctx context.Context) error {
	return nil
}
