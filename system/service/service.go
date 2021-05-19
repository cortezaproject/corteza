package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"strings"
	"time"

	automationService "github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/objstore"
	"github.com/cortezaproject/corteza-server/pkg/objstore/minio"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/automation"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

type (
	websocketSender interface {
		Send(kind string, payload interface{}, userIDs ...uint64) error
	}

	Config struct {
		ActionLog options.ActionLogOpt
		Storage   options.ObjectStoreOpt
		Template  options.TemplateOpt
		Auth      options.AuthOpt
		RBAC      options.RBACOpt
	}

	eventDispatcher interface {
		WaitFor(ctx context.Context, ev eventbus.Event) (err error)
		Dispatch(ctx context.Context, ev eventbus.Event)
	}
)

var (
	DefaultObjectStore objstore.Store

	// DefaultStore is an interface to storage backend(s)
	// ng (next-gen) is a temporary prefix
	// so that we can differentiate between it and the file-only store
	DefaultStore store.Storer

	DefaultLogger *zap.Logger

	// DefaultSettings controls system's settings
	DefaultSettings *settings

	// DefaultAccessControl Access control checking
	DefaultAccessControl *accessControl

	DefaultAuthNotification AuthNotificationService

	// CurrentSettings represents current system settings
	CurrentSettings = &types.AppSettings{}

	DefaultActionlog actionlog.Recorder

	DefaultSink *sink

	DefaultAuth        *auth
	DefaultAuthClient  *authClient
	DefaultUser        UserService
	DefaultRole        *role
	DefaultApplication *application
	DefaultReminder    ReminderService
	DefaultAttachment  AttachmentService
	DefaultRenderer    TemplateService
	DefaultQueue       *queue

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

	DefaultAccessControl = AccessControl()

	DefaultSettings = Settings(ctx, DefaultStore, DefaultLogger, DefaultAccessControl, CurrentSettings)

	if DefaultObjectStore == nil {
		const svcPath = "system"
		if c.Storage.MinioEndpoint != "" {
			var bucket = svcPath
			if c.Storage.MinioBucket != "" {
				bucket = c.Storage.MinioBucket + c.Storage.MinioBucketSep + svcPath
			}

			DefaultObjectStore, err = minio.New(bucket, minio.Options{
				Endpoint:        c.Storage.MinioEndpoint,
				Secure:          c.Storage.MinioSecure,
				Strict:          c.Storage.MinioStrict,
				AccessKeyID:     c.Storage.MinioAccessKey,
				SecretAccessKey: c.Storage.MinioSecretKey,

				ServerSideEncryptKey: []byte(c.Storage.MinioSSECKey),
			})

			log.Info("initializing minio",
				zap.String("bucket", bucket),
				zap.String("endpoint", c.Storage.MinioEndpoint),
				zap.Error(err))
		} else {
			path := c.Storage.Path + "/" + svcPath
			DefaultObjectStore, err = plain.New(path)
			log.Info("initializing store",
				zap.String("path", path),
				zap.Error(err))
		}

		if err != nil {
			return err
		}
	}

	hcd.Add(objstore.Healthcheck(DefaultObjectStore), "ObjectStore/System")

	DefaultRenderer = Renderer(c.Template)
	DefaultAuthNotification = AuthNotification(CurrentSettings, DefaultRenderer, c.Auth)
	DefaultAuth = Auth()
	DefaultAuthClient = AuthClient(DefaultStore, DefaultAccessControl, DefaultActionlog, eventbus.Service())
	DefaultUser = User(ctx)
	DefaultRole = Role()
	DefaultApplication = Application(DefaultStore, DefaultAccessControl, DefaultActionlog, eventbus.Service())
	DefaultReminder = Reminder(ctx, DefaultLogger.Named("reminder"), ws)
	DefaultSink = Sink()
	DefaultStatistics = Statistics()
	DefaultAttachment = Attachment(DefaultObjectStore)
	DefaultQueue = Queue()

	automationService.DefaultUser = DefaultUser

	automationService.Registry().AddTypes(
		automation.User{},
		automation.Role{},
		automation.Template{},
		automation.RenderOptions{},
		automation.RenderedDocument{},
	)

	automation.UsersHandler(
		automationService.Registry(),
		DefaultUser,
	)

	automation.TemplatesHandler(
		automationService.Registry(),
		DefaultRenderer,
	)

	automation.RolesHandler(
		automationService.Registry(),
		DefaultRole,
	)

	return
}

func Activate(ctx context.Context) (err error) {
	// Run initial update of current settings with super-user credentials
	err = DefaultSettings.UpdateCurrent(intAuth.SetSuperUserContext(ctx))
	if err != nil {
		return
	}

	return
}

func Watchers(ctx context.Context) {
	DefaultReminder.Watch(ctx)
	return
}

// Configures RBAC with roles
//
// Sets all closed & im
func initRoles(ctx context.Context, log *zap.Logger, opt options.RBACOpt) (err error) {
	var (
		// splits space separated string into map
		s = func(l string) (map[string]bool, error) {
			m := make(map[string]bool)
			for _, r := range strings.Split(l, " ") {
				if r = strings.TrimSpace(r); len(r) == 0 {
					continue
				}

				if !handle.IsValid(r) {
					return nil, fmt.Errorf("invalid handle '%s'", r)
				}

				m[r] = true
			}

			return m, nil
		}

		// joins map keys into string slice
		j = func(mm ...map[string]bool) []string {
			o := make([]string, 0)

			for _, m := range mm {
				for r := range m {
					o = append(o, r)
				}
			}

			return o
		}

		system, service, bypass, authenticated, anonymous map[string]bool
	)

	if bypass, err = s(opt.BypassRoles); err != nil {
		return fmt.Errorf("failed to process list of bypass roles (RBAC_BYPASS_ROLES): %w", err)
	}
	if authenticated, err = s(opt.AuthenticatedRoles); err != nil {
		return fmt.Errorf("failed to process list of authenticated roles (RBAC_AUTHENTICATED_ROLES): %w", err)
	}
	if anonymous, err = s(opt.AnonymousRoles); err != nil {
		return fmt.Errorf("failed to process list of anonymous roles (RBAC_ANONYMOUS_ROLES): %w", err)
	}

	if len(service) != 1 {
		return fmt.Errorf("role %s used for authenticated users can not be used as bypass role")
	}

	for r := range authenticated {
		if bypass[r] {
			return fmt.Errorf("role %s used for authenticated users must not be used as bypass role")
		}
	}

	for r := range anonymous {
		if bypass[r] {
			return fmt.Errorf("role %s used for anonymous users must not be used as bypass role")
		}

		if authenticated[r] {
			return fmt.Errorf("role %s used for anonymous users must not be used as bypass role")
		}
	}

	DefaultRole.SetSystem(j(system)...)
	DefaultRole.SetClosed(j(authenticated, anonymous)...)

	// Hook to role create, update & delete events and
	// re-apply all roles to RBAC
	eventbus.Service().Register(
		func(_ context.Context, ev eventbus.Event) (err error) {
			var (
				p  = expr.NewParser()
				f  = types.RoleFilter{}
				rr []*rbac.Role
			)

			f.Paging.Total = 0
			roles, _, err := DefaultStore.SearchRoles(ctx, f)
			for _, r := range roles {
				log := log.With(
					zap.Uint64("id", r.ID),
					zap.String("handle", r.Handle),
					zap.String("expr", r.Meta.ContextExpr),
				)

				switch {
				case bypass[r.Handle]:
					rr = append(rr, rbac.BypassRole.Make(r.ID, r.Handle))

				case anonymous[r.Handle]:
					rr = append(rr, rbac.AnonymousRole.Make(r.ID, r.Handle))

				case authenticated[r.Handle]:
					rr = append(rr, rbac.AuthenticatedRole.Make(r.ID, r.Handle))

				case r.Meta != nil && len(r.Meta.ContextExpr) > 0:
					log := log.With(zap.String("expr", r.Meta.ContextExpr))
					eval, err := p.Parse(r.Meta.ContextExpr)
					if err != nil {
						log.Error("failed to parse role context expression", zap.Error(err))
						continue
					}

					check := func(s map[string]interface{}) bool {
						vars, err := expr.NewVars(s)
						if err != nil {
							log.Error("failed to convert check scope to expr.Vars", zap.Error(err))
							return false
						}

						test, err := eval.Test(ctx, vars)
						if err != nil {
							log.Error("failed to evaluate role context expression", zap.Error(err))
							return false
						}

						return test
					}

					rr = append(rr, rbac.MakeContextRole(r.ID, r.Handle, check))

				default:
					rr = append(rr, rbac.CommonRole.Make(r.ID, r.Handle))

				}
			}

			log.Info("role changed " + ev.EventType())
			rbac.Global().UpdateRoles(rr...)

			return nil
		},
		eventbus.For("system:role"),
		eventbus.On("afterUpdate", "afterCreate", "afterDelete"),
	)

	return nil
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
