package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	authService "github.com/crusttech/human/server/auth"
	"github.com/crusttech/human/server/auth/saml"
	authSettings "github.com/crusttech/human/server/auth/settings"
	autoAgentic "github.com/crusttech/human/server/automation/agentic"
	autService "github.com/crusttech/human/server/automation/service"
	cmpAgentic "github.com/crusttech/human/server/compose/agentic"
	cmpService "github.com/crusttech/human/server/compose/service"
	cmpEvent "github.com/crusttech/human/server/compose/service/event"
	discoveryService "github.com/crusttech/human/server/discovery/service"
	fedService "github.com/crusttech/human/server/federation/service"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/apigw"
	apigwTypes "github.com/crusttech/human/server/pkg/apigw/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/corredor"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/healthcheck"
	"github.com/crusttech/human/server/pkg/http"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/pkg/mail"
	"github.com/crusttech/human/server/pkg/messagebus"
	"github.com/crusttech/human/server/pkg/monitor"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/provision"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/pkg/scheduler"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/pkg/sentry"
	"github.com/crusttech/human/server/pkg/valuestore"
	"github.com/crusttech/human/server/pkg/version"
	"github.com/crusttech/human/server/pkg/websocket"
	"github.com/crusttech/human/server/store"
	sysAgentic "github.com/crusttech/human/server/system/agentic"
	mcpkg "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/service"
	sysService "github.com/crusttech/human/server/system/service"
	sysEvent "github.com/crusttech/human/server/system/service/event"
	"github.com/crusttech/human/server/system/types"
	"github.com/lestrrat-go/jwx/jwt"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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
func (app *HumanApp) Setup() (err error) {
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
			log.Warn("You may experience instability and data loss")
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

		if _, is := os.LookupEnv("AUTH_JWT_EXPIRY"); is {
			log.Warn("AUTH_JWT_EXPIRY is removed. " +
				"JWT expiration value is set from AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME")
		}

		if app.Opt.Auth.SessionLifetime < time.Hour {
			log.Warn("AUTH_SESSION_LIFETIME is set to less then an hour, this might not be what you want." +
				"When user logs-in without 'remember-me',  AUTH_SESSION_LIFETIME is used to set a maximum time before session is expired if user does not interacts with Human. " +
				"Recommended session lifetime value is between one hour (default) and a day")
		}

		if app.Opt.Auth.SessionPermLifetime < time.Hour {
			log.Warn("AUTH_SESSION_PERM_LIFETIME is set to less then an hour, this might not be what you want. " +
				"Recommended permanent session lifetime values are between a day and a year (default)")
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
			return fmt.Errorf(
				"locale service setup: %w; "+
					"if this is development environment set ENVIRONMENT=dev or LOCALE_DEVELOPMENT_MODE=true, "+
					"or run make -C pkg/locale if you want to embed languages before bulding server binary", err)
		} else {
			locale.SetGlobal(languages)
		}
	}

	http.SetupDefaults(
		app.Opt.HTTPClient.Timeout,
		app.Opt.HTTPClient.TlsInsecure,
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
		return fmt.Errorf("corredor setup failed: %w", err)
	}

	{
		// load only setup even if disabled, so we can fail gracefuly
		// on queue push
		messagebus.Setup(options.Messagebus(), app.Log)

		if !app.Opt.Messagebus.Enabled {
			app.Log.Debug("messagebus disabled (MESSAGEBUS_ENABLED=false)")
		}
	}

	app.lvl = bootLevelSetup
	return
}

// InitStore initializes store backend(s) and runs upgrade procedures
func (app *HumanApp) InitStore(ctx context.Context) (err error) {
	if app.lvl >= bootLevelStoreInitialized {
		// Is store already initialised?
		return nil
	} else if err = app.Setup(); err != nil {
		// Initialize previous level
		return fmt.Errorf("app setup failed: %w", err)
	}

	// Do not re-initialize store
	// This will make integration test setup a bit more painless
	if app.Store == nil {
		defer sentry.Recover()

		app.Store, err = store.Connect(ctx, app.Log, app.Opt.DB.DSN, app.Opt.Environment.IsDevelopment())
		if err != nil {
			return fmt.Errorf("could not connect to primary store: %w", err)
		}
	}

	if app.ActionlogStore == nil {
		// dedicated pool/DB when ACTIONLOG_DB_DSN is set; share main store otherwise
		if dsn := app.Opt.ActionLog.DbDsn; dsn != "" {
			app.Log.Info("connecting dedicated actionlog store", zap.String("dsn", "***"))
			als, alsErr := store.Connect(ctx, app.Log, dsn, app.Opt.Environment.IsDevelopment())
			if alsErr != nil {
				return fmt.Errorf("could not connect to actionlog store: %w", alsErr)
			}
			if alsErr = store.UpgradeActionlog(ctx, app.Log, als); alsErr != nil {
				return fmt.Errorf("could not upgrade actionlog store: %w", alsErr)
			}
			app.ActionlogStore = als
		} else {
			app.ActionlogStore = app.Store
		}
	}

	if !app.Opt.Upgrade.Always {
		app.Log.Info("store upgrade skipped (UPGRADE_ALWAYS=false)")
	} else {
		log := app.Log.Named("store")
		log.Info("running schema upgrade")
		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Upgrade)

		// If not explicitly set (UPGRADE_DEBUG=true) suppress logging in upgrader
		if app.Opt.Upgrade.Debug {
			log.Info("store upgrade running in debug mode (UPGRADE_DEBUG=true)")
		} else {
			log.Info("store upgrade running (to enable upgrade debug logging set UPGRADE_DEBUG=true)")
			log = zap.NewNop()
		}

		if err = store.Upgrade(ctx, log, app.Store); err != nil {
			return fmt.Errorf("could not upgrade primary store: %w", err)
		}

		healthcheck.Defaults().Add(app.Store.Healthcheck, "Primary store")
	}

	{
		// Initialize Data Access Layer (DAL)
		if err = app.initDAL(ctx, app.Log); err != nil {
			return fmt.Errorf("can not initialize DAL: %w", err)
		}
	}

	app.lvl = bootLevelStoreInitialized
	return nil
}

// Provision instance with configuration and settings
// by importing preset configurations and running autodiscovery procedures
func (app *HumanApp) Provision(ctx context.Context) (err error) {
	if app.lvl >= bootLevelProvisioned {
		return
	}

	if err = app.InitStore(ctx); err != nil {
		return
	}

	if err = app.initSystemEntities(ctx); err != nil {
		return fmt.Errorf("could not initialize system entities: %w", err)
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
			return fmt.Errorf("could not run provision: %w", err)
		}
	}

	app.lvl = bootLevelProvisioned
	return
}

// InitServices initializes all services used
func (app *HumanApp) InitServices(ctx context.Context) (err error) {
	if app.lvl >= bootLevelServicesInitialized {
		return nil
	}

	// @todo place this somewhere better
	id.Init(ctx)

	err = app.initEnvoy(ctx, app.Log)
	if err != nil {
		return
	}

	if err := app.Provision(ctx); err != nil {
		return err
	}

	if err = app.initSystemEntities(ctx); err != nil {
		return
	}

	if err = app.initAuth(ctx); err != nil {
		return fmt.Errorf("can not initialize auth: %w", err)
	}

	initValuestore(app.Opt)

	app.WsServer = websocket.Server(
		app.Log,
		app.Opt.Websocket,
		func(ctx context.Context, s string) (_ auth.Identifiable, err error) {
			var token jwt.Token

			if token, err = jwt.Parse([]byte(s)); err != nil {
				return
			}

			if err = auth.TokenIssuer.Validate(ctx, token); err != nil {
				return
			}

			return auth.IdentityFromToken(token), nil
		},
	)

	corredor.Service().SetAuthTokenMaker(func(i auth.Identifiable) (signed []byte, err error) {
		return auth.TokenIssuer.Issue(ctx,
			auth.WithIdentity(i),
			auth.WithScope("api", "profile"),
			auth.WithAudience("corredor"),
		)
	})

	ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Init)
	defer sentry.Recover()

	if err = corredor.Service().Connect(ctx); err != nil {
		return fmt.Errorf("could not connecto to corredor service: %w", err)
	}

	// Register global (system-level) runtime in the scope registry.
	// Per-tenant and per-project runtimes are added later, keyed by non-zero Scope.
	scope.Default.Set(scope.NewRuntime(scope.Scope{}))

	if rbac.Global() == nil {
		log := zap.NewNop()
		if app.Opt.RBAC.Log {
			log = app.Log
		}

		// Initialize RBAC subsystem
		ac := rbac.NewService(log, app.Store)

		// and (re)load rules from the storage backend
		ac.Reload(ctx)

		rbac.SetGlobal(ac)
	}

	// Initialize resource translation stuff
	locale.Global().BindStore(app.Store)
	if err = locale.Global().ReloadResourceTranslations(ctx); err != nil {
		return fmt.Errorf("could not reload resource translations: %w", err)
	}

	// Initializes system services
	//
	// Note: this is a legacy approach, all services from all 3 apps
	// will most likely be merged in the future
	obs := observability.NewBus(observability.NewLogDispatcher())

	if o := app.Opt.Observability; o.LangfuseHost != "" && o.LangfusePublicKey != "" && o.LangfuseSecretKey != "" {
		obs.Register(observability.NewLangfuseDispatcher(o.LangfuseHost, o.LangfusePublicKey, o.LangfuseSecretKey))
		app.Log.Info("Langfuse dispatcher registered", zap.String("host", o.LangfuseHost))
	}

	if o := app.Opt.Observability; o.OtelExporterOtlpEndpoint != "" {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			app.Log.Warn("failed to initialize OTel exporter", zap.Error(err))
		} else {
			tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
			obs.Register(observability.NewOtelDispatcher(tp))
			app.Log.Info("OTel dispatcher registered", zap.String("endpoint", o.OtelExporterOtlpEndpoint))
		}
	}

	err = sysService.Initialize(ctx, app.Log, app.Store, app.WsServer, sysService.Config{
		ActionLog:      app.Opt.ActionLog,
		ActionlogStore: app.ActionlogStore,
		Discovery:      app.Opt.Discovery,
		Storage:        app.Opt.ObjStore,
		Template:       app.Opt.Template,
		DB:             app.Opt.DB,
		Auth:           app.Opt.Auth,
		RBAC:           app.Opt.RBAC,
		Limit:          app.Opt.Limit,
		Attachment:     app.Opt.Attachment,
		Webapps:        app.Opt.Webapp,
		Agentic:        app.Opt.Agentic,
		Appstore:       app.Opt.Appstore,
		ObsBus:         obs,
	})
	if err != nil {
		return
	}

	cmpAgentic.RecordHandler(sysService.DefaultMCPRegistry)
	cmpAgentic.NamespaceHandler(sysService.DefaultMCPRegistry, sysService.DefaultAgent)
	cmpAgentic.ModuleHandler(sysService.DefaultMCPRegistry, sysService.DefaultAgent)
	cmpAgentic.PageHandler(sysService.DefaultMCPRegistry)
	cmpAgentic.ChartHandler(sysService.DefaultMCPRegistry)
	autoAgentic.TAQHandler(sysService.DefaultMCPRegistry)
	autoAgentic.WorkflowHandler(sysService.DefaultMCPRegistry)
	if app.Opt.Discovery.Enabled && app.Opt.Discovery.BaseUrl != "" && app.Opt.Discovery.JwtSecret != "" {
		discoverySigner, err := auth.NewTokenIssuer(auth.WithSecretSigner(app.Opt.Discovery.JwtSecret))
		if err != nil {
			return fmt.Errorf("could not initialize discovery token signer: %w", err)
		}
		sysAgentic.DiscoveryHandler(sysService.DefaultMCPRegistry, app.Opt.Discovery.BaseUrl, discoverySigner)
	}
	app.McpServer = mcpkg.NewMCPServer(sysService.DefaultMCPRegistry, app.Opt.Agentic.McpServerName, app.Opt.Agentic.McpServerVersion)
	app.LlmService = service.DefaultLlmService

	if app.Opt.Messagebus.Enabled {
		// initialize all the queue handlers
		messagebus.Service().Init(ctx, service.DefaultQueue)
	}

	// Initializes automation services
	//
	// Note: this is a legacy approach, all services from all 3 apps
	// will most likely be merged in the future
	err = autService.Initialize(ctx, app.Log, app.Store, app.WsServer, autService.Config{
		ActionLog:      app.Opt.ActionLog,
		ActionlogStore: app.ActionlogStore,
		Workflow:       app.Opt.Workflow,
		Corredor:       app.Opt.Corredor,
	})

	if err != nil {
		return fmt.Errorf("could not initialize automation services: %w", err)
	}

	// Initializes compose services
	//
	// Note: this is a legacy approach, all services from all 3 apps
	// will most likely be merged in the future
	err = cmpService.Initialize(ctx, app.Log, app.Store, cmpService.Config{
		ActionLog:        app.Opt.ActionLog,
		ActionlogStore:   app.ActionlogStore,
		Discovery:        app.Opt.Discovery,
		Storage:          app.Opt.ObjStore,
		Limit:            app.Opt.Limit,
		UserFinder:       sysService.DefaultUser,
		SchemaAltManager: sysService.DefaultDalSchemaAlteration,
	})

	if err != nil {
		return fmt.Errorf("could not initialize compose services: %w", err)
	}

	sysService.SetDmlComposeDeps(sysService.DmlComposeDeps{
		ModuleSvc:    cmpService.DefaultModule,
		NamespaceSvc: cmpService.DefaultNamespace,
		RecordSvc:    cmpService.DefaultRecord,
	}, dal.Service())

	sysService.SetProjectRevisionDeps(
		cmpService.DefaultNamespace,
		sysService.DefaultDmlImporter,
		dal.Service(),
		cmpService.DefaultRecord,
	)

	sysService.DefaultAgenticRuntime.SetLookups(cmpService.DefaultNamespace, cmpService.DefaultModule)
	sysService.DefaultAgenticRuntime.SetTAQService(autService.DefaultNgAutomation)
	sysService.DefaultAgenticRuntime.SetWorkflowService(autService.DefaultWorkflow)
	sysService.DefaultAgenticRuntime.SetNsModResolver(cmpAgentic.NsModResolver())

	corredor.Service().SetUserFinder(sysService.DefaultUser)
	corredor.Service().SetRoleFinder(sysService.DefaultRole)

	{
		var c = apigwTypes.Config{Enabled: true}

		c.Profiler.Global = sysService.CurrentSettings.Apigw.Profiler.Global
		c.Profiler.Enabled = sysService.CurrentSettings.Apigw.Profiler.Enabled
		c.Proxy.FollowRedirects = sysService.CurrentSettings.Apigw.Proxy.FollowRedirects
		c.Proxy.OutboundTimeout = sysService.CurrentSettings.Apigw.Proxy.OutboundTimeout

		// Initialize API GW bits
		apigw.Setup(c, app.Log, app.Store)
	}

	if app.Opt.Federation.Enabled {
		// Initializes federation services
		//
		// Note: this is a legacy approach, all services from all 3 apps
		// will most likely be merged in the future
		err = fedService.Initialize(ctx, app.Log, app.Store, fedService.Config{
			ActionLog:      app.Opt.ActionLog,
			ActionlogStore: app.ActionlogStore,
			Federation:     app.Opt.Federation,
			Server:         app.Opt.HTTPServer,
		})

		if err != nil {
			return fmt.Errorf("could not initialize federation services: %w", err)
		}
	}

	// Initializing discovery
	if app.Opt.Discovery.Enabled {
		err = discoveryService.Initialize(ctx, app.Log, app.Opt.Discovery, app.Store)
		if err != nil {
			return fmt.Errorf("could not initialize discovery services: %w", err)
		}
	}

	app.lvl = bootLevelServicesInitialized
	return
}

// Activate start all internal services and watchers
func (app *HumanApp) Activate(ctx context.Context) (err error) {
	if app.lvl >= bootLevelActivated {
		return
	}

	err = app.InitExpr(ctx)
	if err != nil {
		return
	}

	if err = app.InitServices(ctx); err != nil {
		return err
	}

	ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Activate)
	defer sentry.Recover()

	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

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
		return fmt.Errorf("could not activate system services: %w", err)

	}

	if err = autService.Activate(ctx); err != nil {
		return fmt.Errorf("could not activate automation services: %w", err)

	}

	if err = cmpService.Activate(ctx); err != nil {
		return fmt.Errorf("could not activate compose services: %w", err)

	}

	if err = applySmtpOptionsToSettings(ctx, app.Log, app.Opt.SMTP, sysService.CurrentSettings); err != nil {
		return fmt.Errorf("could not apply SMTP options to settings: %w", err)
	}

	updateSmtpSettings(app.Log, sysService.CurrentSettings)

	if app.AuthService, err = authService.New(ctx, app.Log, app.oa2m, app.Store, app.Opt.Auth, app.DefaultAuthClient); err != nil {
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

	cmpService.DefaultPage.UpdateConfig(sysService.CurrentSettings)
	sysService.DefaultSettings.Register("compose.ui.record-toolbar", func(ctx context.Context, current interface{}, set types.SettingValueSet) {
		appSettings, is := current.(*types.AppSettings)
		if !is {
			return
		}

		cmpService.DefaultPage.UpdateConfig(appSettings)
	})

	updateDiscoverySettings(app.Opt.Discovery, service.CurrentSettings)
	updateLocaleSettings(app.Opt.Locale)

	app.AuthService.Watch(ctx)

	//Generate CSS for webapps
	if err = sysService.DefaultStylesheet.GenerateCSS(sysService.CurrentSettings, app.Log); err != nil {
		return fmt.Errorf("could not generate css for webapps: %w", err)
	}

	// messagebus reloader and consumer listeners
	if app.Opt.Messagebus.Enabled {

		// set messagebus listener on input channel
		messagebus.Service().Listen(ctx)

		// watch for queue changes and restart on update
		messagebus.Service().Watch(ctx, service.DefaultQueue)
	}

	{
		if err = applyApigwOptionsToSettings(ctx, app.Log, app.Opt.Apigw, sysService.CurrentSettings); err != nil {
			return fmt.Errorf("could not apply integration gateway options to settings: %w", err)
		}

		updateApigwSettings(ctx, sysService.CurrentSettings)

		// // Reload routes
		if err = apigw.Service().Reload(ctx); err != nil {
			return fmt.Errorf("could not initialize api gateway services: %w", err)
		}
	}

	app.lvl = bootLevelActivated
	return nil
}

// Provisions and initializes system roles and users
func (app *HumanApp) initSystemEntities(ctx context.Context) (err error) {
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
		return fmt.Errorf("could not provision system roles: %w", err)
	}

	// Basic provision for system users that we need before anything else
	if uu, err = provision.SystemUsers(ctx, app.Log, app.Store); err != nil {
		return fmt.Errorf("could not provision system users: %w", err)
	}

	// set system users & roles with so that the whole app knows what to use
	auth.SetSystemUsers(uu, rr)
	auth.SetSystemRoles(rr)

	app.Log.Debug(
		"system entities set",
		logger.Uint64s("users", uu.IDs()),
		logger.Uint64s("roles", rr.IDs()),
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
		ProfileAvatarEnabled:      current.Auth.Internal.ProfileAvatar.Enabled,
		SendUserInviteEmail:       current.Auth.Internal.SendUserInviteEmail.Enabled,
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
		BackgroundUI: authSettings.BackgroundUI{
			BackgroundImageSrcUrl: setAuthBgImageSrcUrl(current.Auth.UI.BackgroundImageSrc),
			Styles:                setAuthBgStyles(current.Auth.UI.Styles),
		},

		AutoLogoutEnabled: current.Auth.AutoLogOut.Enabled,
		AutoLogoutTimeout: current.Auth.AutoLogOut.Timeout,
	}

	for _, p := range current.Auth.External.Providers {
		if p.ValidConfiguration() {
			usage := p.Usage

			// By default, use as an identity provider
			if len(p.Usage) == 0 {
				p.Usage = []string{types.ExternalProviderUsageIdentity}
			}

			as.Providers = append(as.Providers, authSettings.Provider{
				Handle:      p.Handle,
				Label:       p.Label,
				IssuerUrl:   p.IssuerUrl,
				Key:         p.Key,
				RedirectUrl: p.RedirectUrl,
				Secret:      p.Secret,
				Scope:       p.Scope,
				Usage:       usage,
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

// Loads current settings into integration gateway and handles the updates / reloads
func updateApigwSettings(ctx context.Context, current *types.AppSettings) {

	updateCurrentSettings := func(ctx context.Context, s *types.AppSettings) {
		var c = apigwTypes.Config{Enabled: true}

		c.Profiler.Enabled = s.Apigw.Profiler.Enabled
		c.Profiler.Global = s.Apigw.Profiler.Global
		c.Proxy.FollowRedirects = s.Apigw.Proxy.FollowRedirects
		c.Proxy.OutboundTimeout = s.Apigw.Proxy.OutboundTimeout

		apigw.Service().UpdateSettings(ctx, c)
	}

	sysService.DefaultSettings.Register("apigw", func(ctx context.Context, current interface{}, _ types.SettingValueSet) {
		appSettings, is := current.(*types.AppSettings)
		if !is {
			return
		}

		updateCurrentSettings(ctx, appSettings)
	})

	// on first load, the options (env) can be different than loaded settings
	// we need to update them here
	updateCurrentSettings(ctx, current)
}

// Checks if discovery is enabled in the options
func updateDiscoverySettings(opt options.DiscoveryOpt, current *types.AppSettings) {
	current.Discovery.Enabled = opt.Enabled
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
			def := locale.Global().Default()
			if def != nil {
				out = append(out, def.Tag.String())
			}
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

// initValuestore initializes and sets the global valuestore with environment variables
func initValuestore(opt *options.Options) {
	s := valuestore.New()

	apiHostname := options.GuessApiHostname()

	// Base variables
	vars := map[string]any{
		// General environment variables such as environment name and version info
		"name":           opt.Environment.Environment,
		"is-development": opt.Environment.IsDevelopment(),
		"is-test":        opt.Environment.IsTest(),
		"is-production":  opt.Environment.IsProduction(),

		"version":    version.Version,
		"build-time": version.BuildTime,
	}

	// Auth variables
	vars["auth.base-url"] = opt.Auth.BaseURL
	vars["auth.domain"] = apiHostname

	// In case there is a missmatch in the auth base URL and server domain,
	// guess the domain from the auth baseURL.
	if !strings.Contains(opt.Auth.BaseURL, apiHostname) {
		u, err := url.Parse(opt.Auth.BaseURL)
		if err != nil {
			panic(err.Error())
		}
		vars["auth.domain"] = u.Host
	}

	// API variables

	// API related values -- domain, base url, base sink route, ...
	vars["api.domain"] = apiHostname
	vars["api.base-url"] = options.FullURL(opt.HTTPServer.BaseUrl, opt.HTTPServer.ApiBaseUrl)

	// Web application (single unified app served at the root)
	webappDomain := ""
	webappBaseURL := ""

	if opt.HTTPServer.WebappEnabled {
		// When served from the server container, use server variables
		webappDomain = apiHostname
		webappBaseURL = options.FullURL(opt.HTTPServer.BaseUrl, opt.HTTPServer.WebappBaseUrl)
	} else {
		// When not served from the server, use client variables
		webappDomain = options.GuessWebappHostname()
		webappBaseURL = options.FullWebappURL(opt.HTTPServer.BaseUrl, opt.HTTPServer.WebappBaseUrl)
	}

	// Webapp related values -- domain, base url (for webapps), ...
	// Splitting the two since the webapp can be served somewhere else on
	// a completely different domain
	vars["webapp.domain"] = webappDomain
	vars["webapp.base-url"] = webappBaseURL

	for k, v := range humanEnvVars() {
		vars[k] = v
	}

	s.SetEnv(vars)
	valuestore.SetGlobal(s)
}

// humanEnvVars returns a set of environment variables that should be passed
// through to Human for use in workflows and elsewhere.
func humanEnvVars() (vars map[string]string) {
	vars = make(map[string]string, 2)

	for _, e := range os.Environ() {
		if !strings.HasPrefix(strings.ToLower(e), "human_env") {
			continue
		}

		pp := strings.SplitN(e, "=", 2)
		if len(pp) != 2 {
			continue
		}

		vars[strings.ToLower(pp[0])] = pp[1]
	}

	return
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
					"When server was restarted, values from environmental " +
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
	err = s.SetSetting([]*types.SmtpServers{optServer})

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

func applyApigwOptionsToSettings(ctx context.Context, log *zap.Logger, opt options.ApigwOpt, current *types.AppSettings) (err error) {
	optApigw := &types.ApigwSettings{Enabled: opt.Enabled}
	optApigw.Profiler.Enabled = opt.ProfilerEnabled
	optApigw.Profiler.Global = opt.ProfilerGlobal
	optApigw.Proxy.FollowRedirects = opt.ProxyFollowRedirects
	optApigw.Proxy.OutboundTimeout = opt.ProxyOutboundTimeout

	if current.Apigw != *optApigw {
		log.Warn(
			"Environmental variables (APIGW_*) and integration gateway settings " +
				"(most likely changed via admin console) are not the same. " +
				"When server was restarted, values from environmental " +
				"variables were copied to settings for easier management. " +
				"To avoid confusion and potential issues, we suggest you to " +
				"remove all APIGW_* variables")
	}

	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	if updateSetting(ctx, "apigw.enabled", optApigw.Enabled) != nil {
		return
	}

	if updateSetting(ctx, "apigw.profiler.enabled", optApigw.Profiler.Enabled) != nil {
		return
	}

	if updateSetting(ctx, "apigw.profiler.global", optApigw.Profiler.Global) != nil {
		return
	}

	if updateSetting(ctx, "apigw.proxy.follow-redirects", optApigw.Proxy.FollowRedirects) != nil {
		return
	}

	if updateSetting(ctx, "apigw.proxy.outbound-timeout", optApigw.Proxy.OutboundTimeout) != nil {
		return
	}

	if err = sysService.DefaultSettings.UpdateCurrent(ctx); err != nil {
		return
	}

	return
}

func updateSetting(ctx context.Context, path string, val interface{}) (err error) {
	s := &types.SettingValue{Name: path}

	err = s.SetSetting(val)

	if err != nil {
		return
	}

	if err = sysService.DefaultSettings.Set(ctx, s); err != nil {
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
		zap.String("from", s.From),
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

func setAuthBgImageSrcUrl(imgAttachment string) string {
	if imgAttachment == "" {
		return ""
	}
	imgAttachmentValues := strings.Split(imgAttachment, ":")
	imgSrcUrl := fmt.Sprintf("/api/system/%s/settings/%s/original/auth.ui.background-image-src", imgAttachmentValues[0], imgAttachmentValues[1])
	return imgSrcUrl
}

// initScopeRBAC initializes the RBAC service for a tenant+project runtime.
func (app *HumanApp) initScopeRBAC(_ context.Context, _ *scope.Runtime) error { return nil }

// initScopeScheduler initializes the scheduler for a tenant+project runtime.
func (app *HumanApp) initScopeScheduler(_ context.Context, _ *scope.Runtime) error { return nil }

// initScopeReminders initializes the reminder watcher for a tenant+project runtime.
func (app *HumanApp) initScopeReminders(_ context.Context, _ *scope.Runtime) error { return nil }

// initScopeCredentialRefresher initializes the credential refresher for a tenant+project runtime.
func (app *HumanApp) initScopeCredentialRefresher(_ context.Context, _ *scope.Runtime) error {
	return nil
}

// initScopeFederationStructureSync initializes federation structure sync for a tenant+project runtime.
func (app *HumanApp) initScopeFederationStructureSync(_ context.Context, _ *scope.Runtime) error {
	return nil
}

// initScopeFederationDataSync initializes federation data sync for a tenant+project runtime.
func (app *HumanApp) initScopeFederationDataSync(_ context.Context, _ *scope.Runtime) error {
	return nil
}

// initScopeConfiguredConnections initializes configured connection refresh for a tenant+project runtime.
func (app *HumanApp) initScopeConfiguredConnections(_ context.Context, _ *scope.Runtime) error {
	return nil
}

// initScopeAutomationSessions initializes the automation session manager for a tenant+project runtime.
func (app *HumanApp) initScopeAutomationSessions(_ context.Context, _ *scope.Runtime) error {
	return nil
}

// initScopeAutomationWorkflows initializes the automation workflow executor for a tenant+project runtime.
func (app *HumanApp) initScopeAutomationWorkflows(_ context.Context, _ *scope.Runtime) error {
	return nil
}

// initScopeMessagebus initializes the messagebus for a tenant+project runtime.
func (app *HumanApp) initScopeMessagebus(_ context.Context, _ *scope.Runtime) error { return nil }

// loadAllScopes fetches every tenant and every project from the store.
// Used during boot to seed the scope registry.
func (app *HumanApp) loadAllScopes(ctx context.Context) (tenants types.TenantSet, projects types.ProjectSet, err error) {
	tf := types.TenantFilter{Paging: filter.Paging{Limit: 200}}
	for {
		var page types.TenantSet
		page, tf, err = app.Store.SearchTenants(ctx, tf)
		if err != nil {
			return
		}
		tenants = append(tenants, page...)
		if tf.NextPage == nil {
			break
		}
		tf.PageCursor = tf.NextPage
	}

	// Archived projects are still projects at runtime -- shelving one hides it
	// from listings, it does not stop it resolving a scope -- and the filter
	// excludes them by default.
	pf := types.ProjectFilter{
		Archived: filter.StateInclusive,
		Paging:   filter.Paging{Limit: 200},
	}
	for {
		var page types.ProjectSet
		page, pf, err = app.Store.SearchProjects(ctx, pf)
		if err != nil {
			return
		}
		projects = append(projects, page...)
		if pf.NextPage == nil {
			break
		}
		pf.PageCursor = pf.NextPage
	}

	return
}

func setAuthBgStyles(styles string) string {
	re := regexp.MustCompile(`\{(.*?)\}`)

	styles = strings.Replace(styles, "\n", "", -1)
	matches := re.FindAllStringSubmatch(styles, -1)

	if matches != nil {
		return matches[0][1]
	}

	return ""
}
