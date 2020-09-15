package app

import (
	"context"
	"errors"
	"fmt"
	cmpService "github.com/cortezaproject/corteza-server/compose/service"
	cmpEvent "github.com/cortezaproject/corteza-server/compose/service/event"
	msgService "github.com/cortezaproject/corteza-server/messaging/service"
	msgEvent "github.com/cortezaproject/corteza-server/messaging/service/event"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/mail"
	"github.com/cortezaproject/corteza-server/pkg/monitor"
	"github.com/cortezaproject/corteza-server/pkg/scheduler"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/provision/compose"
	"github.com/cortezaproject/corteza-server/provision/messaging"
	"github.com/cortezaproject/corteza-server/provision/system"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/auth/external"
	sysService "github.com/cortezaproject/corteza-server/system/service"
	sysEvent "github.com/cortezaproject/corteza-server/system/service/event"
	"go.uber.org/zap"
	"time"
)

type (
	storeUpgrader interface {
		Upgrade(context.Context, *zap.Logger) error
	}
)

const (
	bootLevelWaiting = iota
	bootLevelSetup
	bootLevelStoreInitialized
	bootLevelServicesInitialized
	bootLevelUpgraded
	bootLevelProvisioned
	bootLevelActivated
)

// Setup configures all required services
func (app *CortezaApp) Setup() (err error) {
	app.Log = logger.Default()

	if app.lvl >= bootLevelSetup {
		// Are basics already set-up?
		return nil
	}

	hcd := healthcheck.Defaults()
	hcd.Add(scheduler.Healthcheck, "Scheduler")
	hcd.Add(mail.Healthcheck, "Mail")
	hcd.Add(corredor.Healthcheck, "Corredor")

	if err = sentry.Init(app.Opt.Sentry); err != nil {
		return fmt.Errorf("could not initialize Sentry: %w", err)
	}

	// Use Sentry right away to handle any panics
	// that might occur inside auth, mail setup...
	defer sentry.Recover()

	auth.SetupDefault(app.Opt.Auth.Secret, int(app.Opt.Auth.Expiry/time.Minute))
	mail.SetupDialer(app.Opt.SMTP.Host, app.Opt.SMTP.Port, app.Opt.SMTP.User, app.Opt.SMTP.Pass, app.Opt.SMTP.From)

	http.SetupDefaults(
		app.Opt.HTTPClient.HttpClientTimeout,
		app.Opt.HTTPClient.ClientTSLInsecure,
	)

	monitor.Setup(app.Log, app.Opt.Monitor)

	scheduler.Setup(app.Log, eventbus.Service(), 0)
	scheduler.Service().OnTick(
		sysEvent.SystemOnInterval(),
		sysEvent.SystemOnTimestamp(),
		cmpEvent.ComposeOnInterval(),
		cmpEvent.ComposeOnTimestamp(),
		msgEvent.MessagingOnInterval(),
		msgEvent.MessagingOnTimestamp(),
	)

	if err = corredor.Setup(app.Log, app.Opt.Corredor); err != nil {
		return err
	}

	app.lvl = bootLevelSetup
	return
}

// InitStore initializes store backend(s) and runs upgrade procedures
func (app *CortezaApp) InitStore(ctx context.Context) (err error) {
	if app.lvl >= bootLevelStoreInitialized {
		// Is store already initialised?
		return nil
	} else if err = app.Setup(); err != nil {
		// Initialize previous level
		return err
	}

	// Do not re-initialize store
	// This will make integration test setup a bit more painless
	if app.Store == nil {
		defer sentry.Recover()

		app.Store, err = store.Connect(ctx, app.Opt.DB.DSN)
		if err != nil {
			return err
		}
	}

	if upgradableStore, ok := app.Store.(storeUpgrader); !ok {
		app.Log.Debug("store does not support upgrades")
	} else if !app.Opt.Upgrade.Always {
		app.Log.Debug("store upgrade skipped (UPGRADE_ALWAYS=false)")
	} else {
		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Upgrade)

		// If not explicitly set (UPGRADE_DEBUG=true) suppress logging in upgrader
		log := zap.NewNop()
		if app.Opt.Upgrade.Debug {
			log = app.Log.Named("store.upgrade")
			log.Debug("store upgrade running in debug mode (UPGRADE_DEBUG=true)")
		} else {
			app.Log.Debug("store upgrade running (to enable upgrade debug logging set UPGRADE_DEBUG=true)")

		}

		if err = upgradableStore.Upgrade(ctx, log); err != nil {
			return err
		}
	}

	app.lvl = bootLevelStoreInitialized
	return nil
}

// InitServices initializes all services used
func (app *CortezaApp) InitServices(ctx context.Context) (err error) {
	if app.lvl >= bootLevelServicesInitialized {
		return nil
	} else if err := app.InitStore(ctx); err != nil {
		return err

	}

	ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Init)
	defer sentry.Recover()

	if err = corredor.Service().Connect(ctx); err != nil {
		return
	}

	corredor.Service().SetUserFinder(sysService.DefaultUser)
	corredor.Service().SetRoleFinder(sysService.DefaultRole)

	// Initializes system services
	//
	// Note: this is a legacy approach, all services from all 3 apps
	// will most likely be merged in the future
	err = sysService.Initialize(ctx, app.Log, app.Store, sysService.Config{
		ActionLog: app.Opt.ActionLog,
		Storage:   app.Opt.ObjStore,
	})

	if err != nil {
		return
	}

	// Initializes compose services
	//
	// Note: this is a legacy approach, all services from all 3 apps
	// will most likely be merged in the future
	err = cmpService.Initialize(ctx, app.Log, app.Store, cmpService.Config{
		ActionLog: app.Opt.ActionLog,
		Storage:   app.Opt.ObjStore,
	})

	if err != nil {
		return
	}

	// Initializes messaging services
	//
	// Note: this is a legacy approach, all services from all 3 apps
	// will most likely be merged in the future
	err = msgService.Initialize(ctx, app.Log, app.Store, msgService.Config{
		ActionLog: app.Opt.ActionLog,
		Storage:   app.Opt.ObjStore,
	})

	if err != nil {
		return
	}

	// Initialize external authentication (from default settings)
	external.Init()

	app.lvl = bootLevelServicesInitialized
	return
}

// Provision instance with configuration and settings
// by importing preset configurations and running autodiscovery procedures
func (app *CortezaApp) Provision(ctx context.Context) (err error) {
	if app.lvl >= bootLevelProvisioned {
		return
	}

	if err = app.InitServices(ctx); err != nil {
		return err
	}

	if !app.Opt.Provision.Always {
		app.Log.Debug("provisioning skipped (PROVISION_ALWAYS=false)")
	} else {
		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Provision)
		defer sentry.Recover()

		ctx = auth.SetSuperUserContext(ctx)

		if err = system.Provision(ctx, app.Log, app.Store); err != nil {
			return fmt.Errorf("could not provision system: %w", err)
		}

		if err = compose.Provision(ctx, app.Log, app.Store); err != nil {
			return fmt.Errorf("could not provision compose: %w", err)
		}

		if err = messaging.Provision(ctx, app.Log, app.Store); err != nil {
			return fmt.Errorf("could not provision messaging: %w", err)
		}

		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}

		if err != nil {
			return err
		}
	}

	app.lvl = bootLevelProvisioned
	return
}

// Activate start all internal services and watchers
func (app *CortezaApp) Activate(ctx context.Context) (err error) {
	if app.lvl >= bootLevelActivated {
		return
	} else if err := app.Provision(ctx); err != nil {
		return err
	}

	ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Activate)
	defer sentry.Recover()

	// Start scheduler
	scheduler.Service().Start(ctx)

	// Load corredor scripts & init watcher (script reloader)
	corredor.Service().Load(ctx)
	corredor.Service().Watch(ctx)

	sysService.Watchers(ctx)
	cmpService.Watchers(ctx)
	msgService.Watchers(ctx)

	if err = sysService.Activate(ctx); err != nil {
		return err
	}

	if err = cmpService.Activate(ctx); err != nil {
		return err
	}

	if err = msgService.Activate(ctx); err != nil {
		return err
	}

	app.lvl = bootLevelActivated
	return nil
}
