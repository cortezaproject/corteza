package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bep/godartsass/v2"

	automationService "github.com/crusttech/human/server/automation/service"
	discoveryService "github.com/crusttech/human/server/discovery/service"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/pkg/healthcheck"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/pkg/objstore"
	"github.com/crusttech/human/server/pkg/objstore/minio"
	"github.com/crusttech/human/server/pkg/objstore/plain"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/pkg/valuestore"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/api/cred_registry"
	"github.com/crusttech/human/server/system/service/appstore"
	agenticGuard "github.com/crusttech/human/server/system/agentic/guard"
	agenticKnowledge "github.com/crusttech/human/server/system/agentic/knowledge"
	agenticMcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/crusttech/human/server/system/agentic/observability"
	agenticRuntime "github.com/crusttech/human/server/system/agentic/runtime"
	agenticSkills "github.com/crusttech/human/server/system/agentic/skills"
	"github.com/crusttech/human/server/system/automation"
	"github.com/crusttech/human/server/system/llm"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
)

type (
	websocketSender interface {
		Send(kind string, payload interface{}, userIDs ...uint64) error
	}

	Config struct {
		ActionLog       options.ActionLogOpt
		Discovery       options.DiscoveryOpt
		Storage         options.ObjectStoreOpt
		DB              options.DBOpt
		Template        options.TemplateOpt
		Auth            options.AuthOpt
		RBAC            options.RbacOpt
		Limit           options.LimitOpt
		Attachment      options.AttachmentOpt
		Webapps         options.WebappOpt
		Agentic         options.AgenticOpt
		Appstore        options.AppstoreOpt
		ObsBus          *observability.Bus
		NamespaceLookup agenticKnowledge.NamespaceLookup
		ModuleLookup    agenticKnowledge.ModuleLookup
	}

	eventDispatcher interface {
		WaitFor(ctx context.Context, ev eventbus.Event) (err error)
		Dispatch(ctx context.Context, ev eventbus.Event)
	}

	// AgenticRunner abstracts the agentic runtime execution.
	AgenticRunner interface {
		Run(ctx context.Context, req *agenticRuntime.AgentRequest) (*agenticRuntime.AgentResponse, error)
		SetLookups(ns agenticKnowledge.NamespaceLookup, mod agenticKnowledge.ModuleLookup)
		SetTAQService(s agenticRuntime.TAQService)
		SetWorkflowService(s agenticRuntime.WorkflowService)
		SetNsModResolver(s agenticRuntime.NsModResolver)
		SetBuiltinGuard(g agenticGuard.GuardService)
		SetProviderGuard(g agenticGuard.GuardService)
		SetSkillRegistry(s agenticSkills.Registry)
	}
)

var (
	// DefaultObsBus is the agentic observability bus, exposed so that
	// surface-level integrations (chatbot widget SSE) can attach dispatchers.
	DefaultObsBus *observability.Bus

	DefaultObjectStore objstore.Store

	// DefaultStore is an interface to storage backend(s)
	// ng (next-gen) is a temporary prefix
	// so that we can differentiate between it and the file-only store
	DefaultStore store.Storer

	DefaultLogger *zap.Logger

	// DefaultSettings controls system's settings
	DefaultSettings *settings

	DefaultStylesheet *stylesheet

	// DefaultAccessControl Access control checking
	DefaultAccessControl *accessControl

	DefaultAuthNotification AuthNotificationService

	// CurrentSettings represents current system settings
	CurrentSettings = &types.AppSettings{}

	DefaultActionlog actionlog.Recorder

	DefaultSink *sink

	DefaultAuth                 *auth
	DefaultAuthClient           *authClient
	DefaultUser                 *user
	DefaultCredentials          *credentials
	DefaultDalConnection        *dalConnection
	DefaultDalSensitivityLevel  *dalSensitivityLevel
	DefaultDalSchemaAlteration  *dalSchemaAlteration
	DefaultRole                 *role
	DefaultUserGroup            *userGroup
	DefaultApplication          *application
	DefaultReminder             ReminderService
	DefaultNotification         NotificationService
	DefaultAttachment           AttachmentService
	DefaultRenderer             TemplateService
	DefaultResourceTranslation  ResourceTranslationService
	DefaultQueue                *queue
	DefaultAgent                *agent
	DefaultAiConversation       *aiConversation
	DefaultKnowledgeBase        *knowledgeBase
	DefaultChatbot              *chatbot
	DefaultChatbotSession       *chatbotSession
	DefaultChatbotPreview       *chatbotPreview
	DefaultAgenticRuntime       AgenticRunner
	DefaultMCPRegistry          *agenticMcp.Registry
	DefaultLlmService           *llm.Service
	DefaultApigwRoute           *apigwRoute
	DefaultApigwFilter          *apigwFilter
	DefaultApigwProfiler        *apigwProfiler
	DefaultReport               *report
	DefaultDataPrivacy          *dataPrivacy
	DefaultSMTPChecker          *smtpConfigurationChecker
	DefaultExpression           *expression
	DefaultConnection           *connection
	DefaultConfiguredConnection *configuredConnection

	DefaultStatistics *statistics

	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}

	// wrapper around nextID that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

func Initialize(ctx context.Context, log *zap.Logger, s store.Storer, ws websocketSender, c Config) (err error) {
	var (
		hcd = healthcheck.Defaults()
	)

	// we're doing conversion to avoid having
	// store interface exposed or generated inside app package
	DefaultStore = s

	DefaultLogger = log.Named("service")

	{
		tee := zap.NewNop()
		policy := actionlog.MakeProductionPolicy()

		if !c.ActionLog.Enabled {
			policy = actionlog.MakeDisabledPolicy()
		} else if c.ActionLog.Debug {
			policy = actionlog.MakeDebugPolicy()
			tee = logger.MakeDebugLogger()
		}

		DefaultActionlog = actionlog.NewService(DefaultStore, log, tee, policy)
	}

	// Activity log for system resources
	{
		l := log
		if !c.Discovery.Debug {
			l = zap.NewNop()
		}

		DefaultResourceActivity := discoveryService.ResourceActivity(l, c.Discovery, DefaultStore, eventbus.Service())
		err = DefaultResourceActivity.InitResourceActivityLog(ctx, []string{
			// (types.User{}).RbacResource(), // @todo user?? suppose to be system:user
			"system:user",
		})
		if err != nil {
			return err
		}
	}

	sassTranspiler := dartSassTranspiler(log)

	DefaultAccessControl = AccessControl(s)

	DefaultSettings = Settings(ctx, DefaultStore, DefaultLogger, DefaultAccessControl, DefaultActionlog, CurrentSettings, c.Webapps)
	DefaultStylesheet = Stylesheet(sassTranspiler, log)

	// Initialize credential registry for DAL API connections
	if err = initializeCredentialRegistry(DefaultStore, DefaultLogger); err != nil {
		return fmt.Errorf("failed to initialize credential registry: %w", err)
	}

	DefaultDalConnection = DalConnection(ctx, dal.Service(), c.DB)

	DefaultDalSensitivityLevel = SensitivityLevel(ctx, dal.Service())

	DefaultDalSchemaAlteration = DalSchemaAlteration(dal.Service())

	if DefaultObjectStore == nil {
		var (
			opt    = c.Storage
			bucket string
		)
		const svcPath = "system"
		if opt.MinioEndpoint != "" {
			bucket = minio.GetBucket(opt.MinioBucket, svcPath)

			DefaultObjectStore, err = minio.New(bucket, opt.MinioPathPrefix, svcPath, minio.Options{
				Endpoint:        opt.MinioEndpoint,
				Secure:          opt.MinioSecure,
				Strict:          opt.MinioStrict,
				AccessKeyID:     opt.MinioAccessKey,
				SecretAccessKey: opt.MinioSecretKey,

				ServerSideEncryptKey: []byte(opt.MinioSSECKey),
			})

			log.Info("initializing minio",
				zap.String("bucket", bucket),
				zap.String("endpoint", opt.MinioEndpoint),
				zap.Error(err))
		} else {
			path := opt.Path + "/" + svcPath
			DefaultObjectStore, err = plain.New(path)
			log.Info("initializing store",
				zap.String("path", path),
				zap.Error(err))
		}

		hcd.Add(objstore.Healthcheck(DefaultObjectStore), "ObjectStore/System")

		if err != nil {
			return err
		}

	}

	DefaultRenderer = Renderer(c.Template)
	DefaultResourceTranslation = ResourceTranslation()
	DefaultAuthNotification = AuthNotification(CurrentSettings, DefaultRenderer, c.Auth)
	DefaultAuth = Auth(AuthOptions{LimitUsers: c.Limit.SystemUsers})
	DefaultAuthClient = AuthClient(DefaultStore, DefaultAccessControl, DefaultActionlog, eventbus.Service(), c.Auth)
	DefaultAttachment = Attachment(DefaultObjectStore, c.Attachment, DefaultLogger)
	DefaultUser = User(UserOptions{LimitUsers: c.Limit.SystemUsers})
	DefaultCredentials = Credentials()
	DefaultReport = Report(DefaultStore, DefaultAccessControl, DefaultActionlog, eventbus.Service())
	DefaultRole = Role(rbac.Global())
	DefaultUserGroup = UserGroup(rbac.Global())
	DefaultApplication = Application(DefaultStore, DefaultAccessControl, DefaultActionlog, eventbus.Service())
	DefaultReminder = Reminder(ctx, DefaultLogger.Named("reminder"), ws)
	DefaultNotification = Notification(ctx, DefaultLogger.Named("notification"), ws)
	DefaultSink = Sink()
	DefaultStatistics = Statistics()
	DefaultQueue = Queue()
	DefaultAgent = Agent()
	DefaultAiConversation = AiConversation()
	DefaultKnowledgeBase = KnowledgeBase()
	DefaultChatbot = Chatbot()
	DefaultChatbotSession = ChatbotSession()

	DefaultLlmService, err = llm.New(s, DefaultAccessControl, c.Agentic.AnthropicApiVersion)
	if err != nil {
		return fmt.Errorf("could not initialize LLM service: %w", err)
	}
	DefaultAgent.WithLLMValidator(DefaultLlmService)

	DefaultMCPRegistry = agenticMcp.NewRegistry()

	DefaultObsBus = c.ObsBus

	// Depends on DefaultObsBus being set so SSE emissions reach subscribers.
	DefaultChatbotPreview = ChatbotPreview()

	DefaultAgenticRuntime = agenticRuntime.Runtime(
		DefaultAgent,
		DefaultLlmService,
		DefaultMCPRegistry,
		DefaultAiConversation,
		DefaultKnowledgeBase,
		c.NamespaceLookup,
		c.ModuleLookup,
		c.ObsBus,
	)

	if skillReg, err := agenticSkills.LoadLibrary(); err != nil {
		log.Error("failed to load skill library", zap.Error(err))
	} else {
		DefaultAgenticRuntime.SetSkillRegistry(skillReg)
	}

	// Always-on built-in guard
	DefaultAgenticRuntime.SetBuiltinGuard(agenticGuard.NewBuiltinGuard())

	// Wire guard provider if an LlmProvider with guard config exists
	if guardProvider, guardKey, err := findGuardProvider(ctx, s); err != nil {
		log.Error("failed to find guard provider", zap.Error(err))
	} else if guardProvider != nil {
		DefaultAgenticRuntime.SetProviderGuard(agenticGuard.NewLlamaGuard(guardProvider, guardKey))
		log.Info("guard provider configured", zap.String("provider", guardProvider.Handle))
	} else {
		log.Info("no guard provider configured")
	}

	// Wire preview service deps now that runtime + agent + conv are constructed.
	DefaultChatbotPreview.WithDeps(DefaultAgenticRuntime, DefaultAgent, DefaultAiConversation, DefaultStore)
	DefaultChatbotSession.WithDeps(DefaultObsBus, DefaultAgenticRuntime, DefaultAgent, DefaultAiConversation)

	DefaultApigwRoute = Route()
	DefaultApigwProfiler = Profiler()
	DefaultApigwFilter = Filter()
	DefaultDataPrivacy = DataPrivacy(DefaultStore, DefaultAccessControl, DefaultActionlog, eventbus.Service())
	DefaultSMTPChecker = SmtpConfigurationChecker(CurrentSettings, DefaultRenderer, DefaultAccessControl, c.Auth)
	DefaultExpression = Expression()
	catalogClient := appstore.New(c.Appstore.URL, c.Appstore.APIKey)
	if catalogClient != nil {
		DefaultLogger.Info("appstore catalog enabled", zap.String("url", c.Appstore.URL))
	} else {
		DefaultLogger.Info("appstore catalog disabled (APPSTORE_URL not set)")
	}
	DefaultConnection = Connection().WithCatalog(catalogClient)
	DefaultConfiguredConnection = ConfiguredConnectionSvc().WithDalConnection(DefaultDalConnection)

	DefaultConnection.WithConfiguredConnection(DefaultConfiguredConnection)

	// Register automation functions from all active configured connections
	DefaultConfiguredConnection.RegisterAllOperations(ctx)

	// Start background Google resource discovery refresh (every 1 hour)
	DefaultConfiguredConnection.StartDiscoveryRefreshLoop(ctx, time.Hour)

	if err = initRoles(ctx, log.Named("rbac.roles"), c.RBAC, eventbus.Service(), rbac.Global()); err != nil {
		return err
	}

	automationService.DefaultUser = DefaultUser

	automationService.Registry().AddTypes(
		automation.User{},
		automation.Role{},
		automation.Template{},
		automation.RenderOptions{},
		automation.RenderedDocument{},
		automation.RbacResource{},
		automation.Reminder{},
	)

	automation.UsersHandler(
		automationService.Registry(),
		DefaultUser,
		DefaultRole,
	)

	automation.NgUsersHandler(
		automationService.ConstructLibrary(),
		automationService.Registry(),
		DefaultUser,
		DefaultRole,
	)

	automation.TemplatesHandler(
		automationService.Registry(),
		DefaultRenderer,
	)

	automation.RolesHandler(
		automationService.Registry(),
		DefaultRole,
		DefaultUser,
	)

	automation.RemindersHandler(
		automationService.Registry(),
		DefaultReminder,
	)

	automation.NgRemindersHandler(
		automationService.ConstructLibrary(),
		automationService.Registry(),
		DefaultReminder,
	)

	automation.RbacHandler(
		automationService.Registry(),
		rbac.Global(),
		DefaultUser,
		DefaultRole,
	)

	// Register notification handler
	automation.NotificationHandler(
		automationService.Registry(),
		DefaultNotification,
		DefaultUser,
		log,
	)

	automation.AgentHandler(
		automationService.Registry(),
		DefaultAgenticRuntime,
		DefaultAgent,
		DefaultUser,
	)

	automation.NgAgentHandler(
		automationService.ConstructLibrary(),
		automationService.Registry(),
		DefaultAgenticRuntime,
		DefaultAiConversation,
	)

	automation.NgNotificationHandler(
		automationService.ConstructLibrary(),
		automationService.Registry(),
		DefaultNotification,
		DefaultUser,
		log,
	)

	// ValuestoreHandler isn't (yet) a system thing but this initialization resides
	// here just so we can easily register it

	automation.ValuestoreHandler(
		automationService.Registry(),
		valuestore.Global(),
	)

	if c.ActionLog.WorkflowFunctionsEnabled {
		// register action-log functions & types only when enabled
		automation.ActionlogHandler(
			automationService.Registry(),
			DefaultActionlog,
		)

		automationService.Registry().AddTypes(
			automation.Action{},
		)
	}

	// Reload DAL sensitivity levels
	err = DefaultDalSensitivityLevel.ReloadSensitivityLevels(ctx, DefaultStore)
	if err != nil {
		return
	}

	// Reload DAL connections
	err = DefaultDalConnection.ReloadConnections(ctx)
	if err != nil {
		return
	}

	return
}

func Watchers(ctx context.Context) {
	DefaultReminder.Watch(ctx)
}

func Activate(ctx context.Context) (err error) {
	// Run initial update of current settings
	err = DefaultSettings.UpdateCurrent(ctx)
	if err != nil {
		return
	}

	err = DefaultUserGroup.Activate(ctx)
	if err != nil {
		return
	}

	return
}

// isGeneric returns true if given error is generic
func isGeneric(err error) bool {
	g, ok := err.(interface{ IsGeneric() bool })
	return ok && g != nil && g.IsGeneric()
}

// unwrapGeneric unwraps error if error is generic (and wrapped)
func unwrapGeneric(err error) error {
	for {
		if isGeneric(err) {
			err = errors.Unwrap(err)
			continue
		}

		return err
	}
}

// Data is stale when new date does not match updatedAt or createdAt (before first update)
//
// @todo This is the same as in compose.service; do we want to make an util thing?
func isStale(new *time.Time, updatedAt *time.Time, createdAt time.Time) bool {
	if new == nil {
		// Change to true for stale-data-check
		return false
	}

	if updatedAt != nil {
		return !new.Equal(*updatedAt)
	}

	return new.Equal(createdAt)
}

func dartSassTranspiler(log *zap.Logger) *godartsass.Transpiler {
	transpiler, err := godartsass.Start(godartsass.Options{
		DartSassEmbeddedFilename: "sass",
	})

	if err != nil {
		log.Warn("dart sass is not installed in your system", zap.Error(err))
		return nil
	}

	return transpiler
}

func initializeCredentialRegistry(s store.Storer, log *zap.Logger) error {
	reg, err := cred_registry.New(s, log)
	if err != nil {
		return err
	}

	cred_registry.SetDefault(reg)
	log.Info("credential registry initialized")
	return nil
}

// findGuardProvider scans all LlmProvider records for one with Config.Guard.Enabled == true
// and Config.Guard.Provider == "llama-guard". Returns the provider and its API key, or nil
// if no guard provider is configured.
func findGuardProvider(ctx context.Context, s store.Storer) (*types.LlmProvider, string, error) {
	set, _, err := store.SearchLlmProviders(ctx, s, types.LlmProviderFilter{})
	if err != nil {
		return nil, "", err
	}

	for _, p := range set {
		if p.Config.Guard == nil || !p.Config.Guard.Enabled || p.Config.Guard.Provider != "llama-guard" {
			continue
		}
		apiKey := ""
		if p.CredentialID != 0 {
			cred, err := store.LookupCredentialByID(ctx, s, p.CredentialID)
			if err == nil && cred != nil {
				apiKey = cred.Credentials
			}
		}
		return p, apiKey, nil
	}

	return nil, "", nil
}
