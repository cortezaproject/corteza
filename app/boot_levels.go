package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	authService "github.com/cortezaproject/corteza-server/auth"
	authHandlers "github.com/cortezaproject/corteza-server/auth/handlers"
	"github.com/cortezaproject/corteza-server/auth/saml"
	authSettings "github.com/cortezaproject/corteza-server/auth/settings"
	autService "github.com/cortezaproject/corteza-server/automation/service"
	cmpService "github.com/cortezaproject/corteza-server/compose/service"
	cmpEvent "github.com/cortezaproject/corteza-server/compose/service/event"
	fdrService "github.com/cortezaproject/corteza-server/federation/service"
	fedService "github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/apigw"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/mail"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/pkg/monitor"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/provision"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/scheduler"
	"github.com/cortezaproject/corteza-server/pkg/seeder"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/cortezaproject/corteza-server/pkg/websocket"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	sysService "github.com/cortezaproject/corteza-server/system/service"
	sysEvent "github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

const (
	bootLevelWaiting = iota
	bootLevelSetup
	bootLevelStoreInitialized
	bootLevelProvisioned
	bootLevelServicesInitialized
	bootLevelActivated
)

// Setup configures all required services
func (app *CortezaApp) Setup() (err error) {
	if app.lvl >= bootLevelSetup {
		// Are basics already set-up?
		return nil
	}

	{
		// Raise warnings about experimental parts that are enabled
		log := app.Log.WithOptions(zap.WithCaller(false))

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

		if _, is := os.LookupEnv("MINIO_BUCKET_SEP"); is {
			log.Warn("Found MINIO_BUCKET_SEP in environment variables, it has been removed")

			return fmt.Errorf(
				"invalid minio configurtion: " +
					"found MINIO_BUCKET_SEP in environment variables, " +
					"which is removed due to latest versions of min.io " +
					"bucket names can only consist of lowercase letters, numbers, dots (.), and hyphens (-). " +
					"so instead use environment variable MINIO_BUCKET, " +
					"we have extended it to have more flexibility over minio bucket name")
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

	{
		var (
			localeLog = zap.NewNop()
		)

		if app.Opt.Locale.Log {
			localeLog = app.Log
		}

		if languages, err := locale.Service(localeLog, app.Opt.Locale); err != nil {
			return err
		} else {
			locale.SetGlobal(languages)
		}
	}

	// set base path for links&routes in auth server
	authHandlers.BasePath = app.Opt.HTTPServer.BaseUrl

	auth.SetupDefault(app.Opt.Auth.Secret, app.Opt.Auth.Expiry)

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
		)
	} else {
		app.Log.Debug("eventbus scheduler disabled (EVENTBUS_SCHEDULER_ENABLED=false)")
	}

	if err = corredor.Setup(app.Log, app.Opt.Corredor); err != nil {
		return err
	}

	{
		// load only setup even if disabled, so we can fail gracefuly
		// on queue push
		messagebus.Setup(options.Messagebus(), app.Log)

		if !app.Opt.Messagebus.Enabled {
			app.Log.Debug("messagebus disabled (MESSAGEBUS_ENABLED=false)")
		}
	}

	if err = app.plugins.Setup(app.Log.Named("plugin")); err != nil {
		return
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

		app.Store, err = store.Connect(ctx, app.Log, app.Opt.DB.DSN, app.Opt.Environment.IsDevelopment())
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

// Provision instance with configuration and settings
// by importing preset configurations and running autodiscovery procedures
func (app *CortezaApp) Provision(ctx context.Context) (err error) {
	if app.lvl >= bootLevelProvisioned {
		return
	}

	if err = app.InitStore(ctx); err != nil {
		return err
	}

	if err = app.initSystemEntities(ctx); err != nil {
		return
	}

	{
		// register temporary RBAC with bypass roles
		// this is needed because envoy relies on availability of access-control
		//
		// @todo envoy should be decoupled from RBAC and import directly into store,
		//       w/o using any access control

		var (
			ac  = rbac.NewService(zap.NewNop(), app.Store)
			acr = make([]*rbac.Role, 0)
		)
		for _, r := range auth.ProvisionUser().Roles() {
			acr = append(acr, rbac.BypassRole.Make(r, auth.BypassRoleHandle))
		}
		ac.UpdateRoles(acr...)
		rbac.SetGlobal(ac)
		defer rbac.SetGlobal(nil)
	}

	if !app.Opt.Provision.Always {
		app.Log.Debug("provisioning skipped (PROVISION_ALWAYS=false)")
	} else {
		defer sentry.Recover()

		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Provision)
		ctx = auth.SetIdentityToContext(ctx, auth.ProvisionUser())

		if err = provision.Run(ctx, app.Log, app.Store, app.Opt.Provision, app.Opt.Auth); err != nil {
			return err
		}
	}

	app.lvl = bootLevelProvisioned
	return
}

// InitServices initializes all services used
func (app *CortezaApp) InitServices(ctx context.Context) (err error) {
	if app.lvl >= bootLevelServicesInitialized {
		return nil
	}

	if err := app.Provision(ctx); err != nil {
		return err
	}

	if err = app.initSystemEntities(ctx); err != nil {
		return
	}

	app.WsServer = websocket.Server(app.Log, app.Opt.Websocket)

	ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Init)
	defer sentry.Recover()

	if err = corredor.Service().Connect(ctx); err != nil {
		return
	}

	if rbac.Global() == nil {
		log := zap.NewNop()
		if app.Opt.RBAC.Log {
			log = app.Log
		}

		//Initialize RBAC subsystem
		ac := rbac.NewService(log, app.Store)

		// and (re)load rules from the storage backend
		ac.Reload(ctx)

		rbac.SetGlobal(ac)
	}

	// Initialize resource translation stuff
	locale.Global().BindStore(app.Store)
	if err = locale.Global().ReloadResourceTranslations(ctx); err != nil {
		return err
	}

	// Initializes system services
	//
	// Note: this is a legacy approach, all services from all 3 apps
	// will most likely be merged in the future
	err = sysService.Initialize(ctx, app.Log, app.Store, app.WsServer, sysService.Config{
		ActionLog: app.Opt.ActionLog,
		Storage:   app.Opt.ObjStore,
		Template:  app.Opt.Template,
		Auth:      app.Opt.Auth,
		RBAC:      app.Opt.RBAC,
		Limit:     app.Opt.Limit,
	})

	if err != nil {
		return
	}

	if app.Opt.Messagebus.Enabled {
		// initialize all the queue handlers
		messagebus.Service().Init(ctx, service.DefaultQueue)
	}

	// Initializes automation services
	//
	// Note: this is a legacy approach, all services from all 3 apps
	// will most likely be merged in the future
	err = autService.Initialize(ctx, app.Log, app.Store, app.WsServer, autService.Config{
		ActionLog: app.Opt.ActionLog,
		Workflow:  app.Opt.Workflow,
		Corredor:  app.Opt.Corredor,
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

	auth.SetJWTStore(app.Store)

	corredor.Service().SetUserFinder(sysService.DefaultUser)
	corredor.Service().SetRoleFinder(sysService.DefaultRole)

	// Initialize API GW bits
	apigw.Setup(options.Apigw(), app.Log, app.Store)
	if err = apigw.Service().Reload(ctx); err != nil {
		return err
	}

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

	// Initializing seeder
	_ = seeder.Seeder(ctx, app.Store, seeder.Faker())

	if err = app.plugins.Initialize(ctx, app.Log); err != nil {
		return
	}

	if err = app.plugins.RegisterAutomation(autService.Registry()); err != nil {
		return
	}

	app.lvl = bootLevelServicesInitialized
	return
}

// Activate start all internal services and watchers
func (app *CortezaApp) Activate(ctx context.Context) (err error) {
	if app.lvl >= bootLevelActivated {
		return
	}

	if err := app.InitServices(ctx); err != nil {
		return err
	}

	ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Activate)
	defer sentry.Recover()

	// Start scheduler
	if app.Opt.Eventbus.SchedulerEnabled {
		scheduler.Service().Start(ctx)
	}

	// Load corredor scripts & init watcher (script reloader)
	{
		ctx := auth.SetIdentityToContext(ctx, auth.ServiceUser())
		corredor.Service().Load(ctx)
		corredor.Service().Watch(ctx)
	}

	sysService.Watchers(ctx)
	autService.Watchers(ctx)
	cmpService.Watchers(ctx)

	if app.Opt.Federation.Enabled {
		fedService.Watchers(ctx)
	}

	monitor.Watcher(ctx)

	rbac.Global().Watch(ctx)

	if err = sysService.Activate(ctx); err != nil {
		return err
	}

	if err = autService.Activate(ctx); err != nil {
		return err
	}

	if err = cmpService.Activate(ctx); err != nil {
		return err
	}

	if err = applySmtpOptionsToSettings(ctx, app.Log, app.Opt.SMTP, sysService.CurrentSettings); err != nil {
		return err
	}

	updateSmtpSettings(app.Log, sysService.CurrentSettings)

	if app.AuthService, err = authService.New(ctx, app.Log, app.Store, app.Opt.Auth); err != nil {
		return fmt.Errorf("failed to init auth service: %w", err)
	}

	app.ApigwService = apigw.Service()

	updateFederationSettings(app.Opt.Federation, sysService.CurrentSettings)
	updateAuthSettings(app.AuthService, sysService.CurrentSettings)
	updatePasswdSettings(app.Opt.Auth, sysService.CurrentSettings)
	sysService.DefaultSettings.Register("auth.", func(ctx context.Context, current interface{}, set types.SettingValueSet) {
		appSettings, is := current.(*types.AppSettings)
		if !is {
			return
		}

		updateAuthSettings(app.AuthService, appSettings)
		updatePasswdSettings(app.Opt.Auth, sysService.CurrentSettings)
	})

	updateLocaleSettings(app.Opt.Locale)

	app.AuthService.Watch(ctx)

	// messagebus reloader and consumer listeners
	if app.Opt.Messagebus.Enabled {

		// set messagebus listener on input channel
		messagebus.Service().Listen(ctx)

		// watch for queue changes and restart on update
		messagebus.Service().Watch(ctx, service.DefaultQueue)
	}

	app.lvl = bootLevelActivated
	return nil
}

// Provisions and initializes system roles and users
func (app *CortezaApp) initSystemEntities(ctx context.Context) (err error) {
	if app.systemEntitiesInitialized {
		// make sure we do this once.
		return nil
	}

	app.systemEntitiesInitialized = true

	var (
		uu types.UserSet
		rr types.RoleSet
	)

	// Basic provision for system resources that we need before anything else
	if rr, err = provision.SystemRoles(ctx, app.Log, app.Store); err != nil {
		return
	}

	// Basic provision for system users that we need before anything else
	if uu, err = provision.SystemUsers(ctx, app.Log, app.Store); err != nil {
		return
	}

	// set system users & roles with so that the whole app knows what to use
	auth.SetSystemUsers(uu, rr)
	auth.SetSystemRoles(rr)

	app.Log.Debug(
		"system entities set",
		zap.Uint64s("users", uu.IDs()),
		zap.Uint64s("roles", rr.IDs()),
	)

	return nil
}

func updateAuthSettings(svc authServicer, current *types.AppSettings) {
	as := &authSettings.Settings{
		LocalEnabled:              current.Auth.Internal.Enabled,
		SignupEnabled:             current.Auth.Internal.Signup.Enabled,
		EmailConfirmationRequired: current.Auth.Internal.Signup.EmailConfirmationRequired,
		PasswordResetEnabled:      current.Auth.Internal.PasswordReset.Enabled,
		PasswordCreateEnabled:     current.Auth.Internal.PasswordCreate.Enabled,
		SplitCredentialsCheck:     current.Auth.Internal.SplitCredentialsCheck,
		ExternalEnabled:           current.Auth.External.Enabled,
		MultiFactor: authSettings.MultiFactor{
			TOTP: authSettings.TOTP{
				Enabled:  current.Auth.MultiFactor.TOTP.Enabled,
				Enforced: current.Auth.MultiFactor.TOTP.Enforced,
				Issuer:   current.Auth.MultiFactor.TOTP.Issuer,
			},
			EmailOTP: authSettings.EmailOTP{
				Enabled:  current.Auth.MultiFactor.EmailOTP.Enabled,
				Enforced: current.Auth.MultiFactor.EmailOTP.Enforced,
			},
		},
	}

	for _, p := range current.Auth.External.Providers {
		if p.ValidConfiguration() {
			as.Providers = append(as.Providers, authSettings.Provider{
				Handle:      p.Handle,
				Label:       p.Label,
				IssuerUrl:   p.IssuerUrl,
				Key:         p.Key,
				RedirectUrl: p.RedirectUrl,
				Secret:      p.Secret,
				Scope:       p.Scope,
			})
		}
	}

	// SAML
	saml.UpdateSettings(current, as)

	svc.UpdateSettings(as)
}

// Checks if federation is enabled in the options
func updateFederationSettings(opt options.FederationOpt, current *types.AppSettings) {
	current.Federation.Enabled = opt.Enabled
}

// Checks if password security is enabled in the options
func updatePasswdSettings(opt options.AuthOpt, current *types.AppSettings) {
	current.Auth.Internal.PasswordConstraints.PasswordSecurity = opt.PasswordSecurity
}

// Sanitizes application (current) settings with languages from options
//
// It updates resource-translations.languages slice
// These do not need to be subset of LOCALE_LANGUAGES but need to be valid language tags!
func updateLocaleSettings(opt options.LocaleOpt) {
	updateResourceLanguages := func(appSettings *types.AppSettings) {
		out := make([]string, 0, 8)

		if opt.ResourceTranslationsEnabled {
			for _, t := range locale.Global().Tags() {
				out = append(out, t.String())
			}
		} else {
			// when resource translation is disabled,
			// add only default (first) language to the list
			out = append(out, locale.Global().Default().Tag.String())
		}

		appSettings.ResourceTranslations.Languages = out
	}

	updateResourceLanguages(sysService.CurrentSettings)
	sysService.DefaultSettings.Register("resource-translations.languages", func(ctx context.Context, current interface{}, _ types.SettingValueSet) {
		appSettings, is := current.(*types.AppSettings)
		if !is {
			return
		}

		updateResourceLanguages(appSettings)
	})
}

// takes current options (SMTP_* env variables) and copies their values to settings
func applySmtpOptionsToSettings(ctx context.Context, log *zap.Logger, opt options.SMTPOpt, current *types.AppSettings) (err error) {
	if len(opt.Host) == 0 {
		// nothing to do here, SMTP_HOST not set
		return
	}

	// Create SMTP server settings struct
	// from the environmental variables (SMTP_*)
	// we'll use it for provisioning empty SMTP settings
	// and for comparison to issue a warning
	optServer := &types.SmtpServers{
		Host:          opt.Host,
		Port:          opt.Port,
		User:          opt.User,
		Pass:          opt.Pass,
		From:          opt.From,
		TlsInsecure:   opt.TlsInsecure,
		TlsServerName: opt.TlsServerName,
	}

	if len(current.SMTP.Servers) > 0 {
		if current.SMTP.Servers[0] != *optServer {
			// ENV variables changed OR settings changed.
			// One way or the other, this can lead to unexpected situations
			//
			// Let's log a warning
			log.Warn(
				"Environmental variables (SMTP_*) and SMTP settings " +
					"(most likely changed via admin console) are not the same. " +
					"When server was restarted, values from environmental" +
					"variables were copied to settings for easier management. " +
					"To avoid confusion and potential issues, we suggest you to " +
					"remove all SMTP_* variables")
		}

		return
	}

	// SMTP server settings do not exist but
	// there is something in the options (SMTP_HOST)
	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	// When settings for the SMTP servers are missing,
	// we'll try to use one from the options (environmental vars)
	s := &types.SettingValue{Name: "smtp.servers"}
	err = s.SetValue([]*types.SmtpServers{optServer})

	if err != nil {
		return
	}

	if err = sysService.DefaultSettings.Set(ctx, s); err != nil {
		return
	}

	if err = sysService.DefaultSettings.UpdateCurrent(ctx); err != nil {
		return
	}

	return
}

func updateSmtpSettings(log *zap.Logger, current *types.AppSettings) {
	sysService.DefaultSettings.Register("smtp", func(ctx context.Context, current interface{}, _ types.SettingValueSet) {
		appSettings, is := current.(*types.AppSettings)
		if !is {
			return
		}

		setupSmtpDialer(log, appSettings.SMTP.Servers...)
	})
	setupSmtpDialer(log, current.SMTP.Servers...)
}

func setupSmtpDialer(log *zap.Logger, servers ...types.SmtpServers) {
	if len(servers) == 0 {
		log.Warn("no SMTP servers found, email sending will be disabled")
		return
	}

	// Supporting only one server for now
	s := servers[0]

	if s.Host == "" {
		log.Warn("SMTP server configured without host/server, email sending will be disabled")
		return
	}

	log.Info("reloading SMTP configuration",
		zap.String("host", s.Host),
		zap.Int("port", s.Port),
		zap.String("user", s.User),
		logger.Mask("pass", s.Pass),
		zap.Bool("tsl-insecure", s.TlsInsecure),
		zap.String("tls-server-name", s.TlsServerName),
	)

	mail.SetupDialer(
		s.Host,
		s.Port,
		s.User,
		s.Pass,
		s.From,

		// Apply TLS configuration
		func(d *gomail.Dialer) {
			if d.TLSConfig == nil {
				d.TLSConfig = &tls.Config{ServerName: d.Host}
			}

			if s.TlsInsecure {
				d.TLSConfig.InsecureSkipVerify = true
			}

			if s.TlsServerName != "" {
				d.TLSConfig.ServerName = s.TlsServerName
			}
		},
	)

}
