package app

import (
	"context"
	"crypto/tls"
	"fmt"
	authService "github.com/cortezaproject/corteza-server/auth"
	"github.com/cortezaproject/corteza-server/auth/settings"
	cmpService "github.com/cortezaproject/corteza-server/compose/service"
	cmpEvent "github.com/cortezaproject/corteza-server/compose/service/event"
	fdrService "github.com/cortezaproject/corteza-server/federation/service"
	fedService "github.com/cortezaproject/corteza-server/federation/service"
	msgService "github.com/cortezaproject/corteza-server/messaging/service"
	msgEvent "github.com/cortezaproject/corteza-server/messaging/service/event"
	"github.com/cortezaproject/corteza-server/messaging/websocket"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/mail"
	"github.com/cortezaproject/corteza-server/pkg/monitor"
	"github.com/cortezaproject/corteza-server/pkg/provision"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/scheduler"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/store"
	sysService "github.com/cortezaproject/corteza-server/system/service"
	sysEvent "github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
	"strings"
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

	{
		// Raise warnings about experimental parts that are enabled
		log := app.Log.WithOptions(zap.AddStacktrace(zap.PanicLevel), zap.WithCaller(false))

		if app.Opt.Federation.Enabled {
			log.Warn("Record Federation is still in EXPERIMENTAL phase")
		}

		if app.Opt.SCIM.Enabled {
			log.Warn("Support for SCIM protocol is still in EXPERIMENTAL phase")
		}

		if app.Opt.DB.IsSQLite() {
			log.Warn("You're using SQLite as a storage backend")
			log.Warn("Should be used only for testing")
			log.Warn("You may experience unstability and data loss")
		}
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

	auth.SetupDefault(app.Opt.Auth.Secret, app.Opt.Auth.Expiry)

	mail.SetupDialer(
		app.Opt.SMTP.Host,
		app.Opt.SMTP.Port,
		app.Opt.SMTP.User,
		app.Opt.SMTP.Pass,
		app.Opt.SMTP.From,

		// Apply TLS configuration
		func(d *gomail.Dialer) {
			if d.TLSConfig == nil {
				d.TLSConfig = &tls.Config{ServerName: d.Host}
			}

			if app.Opt.SMTP.TlsInsecure {
				d.TLSConfig.InsecureSkipVerify = true
			}

			if app.Opt.SMTP.TlsServerName != "" {
				d.TLSConfig.ServerName = app.Opt.SMTP.TlsServerName
			}
		},
	)

	http.SetupDefaults(
		app.Opt.HTTPClient.HttpClientTimeout,
		app.Opt.HTTPClient.ClientTSLInsecure,
	)

	monitor.Setup(app.Log, app.Opt.Monitor)

	if app.Opt.Eventbus.SchedulerEnabled {
		scheduler.Setup(app.Log, eventbus.Service(), app.Opt.Eventbus.SchedulerInterval)
		scheduler.Service().OnTick(
			sysEvent.SystemOnInterval(),
			sysEvent.SystemOnTimestamp(),
			cmpEvent.ComposeOnInterval(),
			cmpEvent.ComposeOnTimestamp(),
			msgEvent.MessagingOnInterval(),
			msgEvent.MessagingOnTimestamp(),
		)
	} else {
		app.Log.Debug("eventbus scheduler disabled (EVENTBUS_SCHEDULER_ENABLED=false)")
	}

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

	app.Log.Info("running store update")

	if !app.Opt.Upgrade.Always {
		app.Log.Info("store upgrade skipped (UPGRADE_ALWAYS=false)")
	} else {
		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Upgrade)

		// If not explicitly set (UPGRADE_DEBUG=true) suppress logging in upgrader
		log := zap.NewNop()
		if app.Opt.Upgrade.Debug {
			log = app.Log.Named("store.upgrade")
			log.Info("store upgrade running in debug mode (UPGRADE_DEBUG=true)")
		} else {
			app.Log.Info("store upgrade running (to enable upgrade debug logging set UPGRADE_DEBUG=true)")
		}

		if err = store.Upgrade(ctx, log, app.Store); err != nil {
			return err
		}

		// @todo refactor this to make more sense and put it where it belongs
		{
			var set types.SettingValueSet
			set, _, err = store.SearchSettings(ctx, app.Store, types.SettingsFilter{Prefix: "auth.external"})
			if err != nil {
				return err
			}

			err = set.Walk(func(old *types.SettingValue) error {
				if strings.HasSuffix(old.Name, ".redirect-url") {
					// remove obsolete redirect-url
					if err = store.DeleteSetting(ctx, app.Store, old); err != nil {
						return err
					}

					return nil
				}

				if strings.Contains(old.Name, ".provider.gplus.") {
					var new = *old
					new.Name = strings.Replace(new.Name, "provider.gplus.", "provider.google.", 1)

					log.Info("renaming settings", zap.String("old", old.Name), zap.String("new", new.Name))

					if err = store.CreateSetting(ctx, app.Store, &new); err != nil {
						if store.ErrNotUnique != err {
							return err
						}
					}

					if err = store.DeleteSetting(ctx, app.Store, old); err != nil {
						return err
					}
				}

				return nil
			})

			if err != nil {
				return err
			}
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

	{
		// Initialize RBAC subsystem
		// and (re)load rules from the storage backend
		err = rbac.Initialize(app.Log, app.Store)
		if err != nil {
			return
		}

		rbac.Global().Reload(ctx)
	}

	// Initializes system services
	//
	// Note: this is a legacy approach, all services from all 3 apps
	// will most likely be merged in the future
	err = sysService.Initialize(ctx, app.Log, app.Store, sysService.Config{
		ActionLog: app.Opt.ActionLog,
		Storage:   app.Opt.ObjStore,
		Template:  app.Opt.Template,
		Auth:      app.Opt.Auth,
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

	corredor.Service().SetUserFinder(sysService.DefaultUser)
	corredor.Service().SetRoleFinder(sysService.DefaultRole)

	app.WsServer = websocket.New(&websocket.Config{
		Timeout:     app.Opt.Websocket.Timeout,
		PingTimeout: app.Opt.Websocket.PingTimeout,
		PingPeriod:  app.Opt.Websocket.PingPeriod,
	})

	if app.Opt.Federation.Enabled {
		// Initializes federation services
		//
		// Note: this is a legacy approach, all services from all 3 apps
		// will most likely be merged in the future
		err = fdrService.Initialize(ctx, app.Log, app.Store, fdrService.Config{
			ActionLog:  app.Opt.ActionLog,
			Federation: app.Opt.Federation,
		})

		if err != nil {
			return
		}
	}

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
		defer sentry.Recover()

		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Provision)
		ctx = auth.SetSuperUserContext(ctx)

		if err = provision.Run(ctx, app.Log, app.Store, app.Opt.Provision, app.Opt.Auth); err != nil {
			return err
		}

		// Provisioning doesn't automatically reload rbac rules, so this is required
		rbac.Global().Reload(ctx)
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
	if app.Opt.Eventbus.SchedulerEnabled {
		scheduler.Service().Start(ctx)
	}

	// Load corredor scripts & init watcher (script reloader)
	corredor.Service().Load(ctx)
	corredor.Service().Watch(ctx)

	sysService.Watchers(ctx)
	cmpService.Watchers(ctx)
	msgService.Watchers(ctx)

	if app.Opt.Federation.Enabled {
		fedService.Watchers(ctx)
	}

	monitor.Watcher(ctx)

	rbac.Global().Watch(ctx)

	if err = sysService.Activate(ctx); err != nil {
		return err
	}

	if err = cmpService.Activate(ctx); err != nil {
		return err
	}

	if err = msgService.Activate(ctx); err != nil {
		return err
	}

	if app.WsServer != nil {
		websocket.Watch(ctx)
	}

	if app.AuthService, err = authService.New(ctx, app.Log, app.Store, app.Opt.Auth); err != nil {
		return fmt.Errorf("failed to init auth service: %w", err)
	}

	updateAuthSettings(app.AuthService, sysService.CurrentSettings)
	sysService.DefaultSettings.Register("auth.", func(ctx context.Context, current interface{}, set types.SettingValueSet) {
		appSettings, is := current.(*types.AppSettings)
		if !is {
			return
		}

		updateAuthSettings(app.AuthService, appSettings)
	})

	app.AuthService.Watch(ctx)

	app.lvl = bootLevelActivated
	return nil
}

func updateAuthSettings(svc authServicer, current *types.AppSettings) {
	var (
		// current auth settings
		cas       = current.Auth
		providers []settings.Provider
	)

	for _, p := range cas.External.Providers {
		if !p.Enabled {
			continue
		}

		providers = append(providers, settings.Provider{
			Handle:      p.Handle,
			Label:       p.Label,
			IssuerUrl:   p.IssuerUrl,
			Key:         p.Key,
			RedirectUrl: p.RedirectUrl,
			Secret:      p.Secret,
		})
	}

	svc.UpdateSettings(&settings.Settings{
		LocalEnabled:              cas.Internal.Enabled,
		SignupEnabled:             cas.Internal.Signup.Enabled,
		EmailConfirmationRequired: cas.Internal.Signup.EmailConfirmationRequired,
		PasswordResetEnabled:      cas.Internal.PasswordReset.Enabled,
		ExternalEnabled:           cas.External.Enabled,
		Providers:                 providers,
	})
}
