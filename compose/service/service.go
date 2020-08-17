package service

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"go.uber.org/zap"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/store"
	"github.com/cortezaproject/corteza-server/pkg/store/minio"
	"github.com/cortezaproject/corteza-server/pkg/store/plain"
	systemService "github.com/cortezaproject/corteza-server/system/service"
)

type (
	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}

	Config struct {
		ActionLog options.ActionLogOpt
		Storage   options.StorageOpt
	}

	eventDispatcher interface {
		WaitFor(ctx context.Context, ev eventbus.Event) (err error)
		Dispatch(ctx context.Context, ev eventbus.Event)
	}

	// storeInterface wraps generated interfaces to enable extensions
	storeInterface interface {
		// Include generated interfaces
		storeGeneratedInterfaces

		// And all additional required functions
		// ...
	}
)

var (
	DefaultStore store.Store

	// DefaultNgStore is an interface to storage backend(s)
	// ng (next-gen) is a temporary prefix
	// so that we can differentiate between it and the file-only store
	DefaultNgStore storeInterface

	DefaultLogger *zap.Logger

	DefaultActionlog actionlog.Recorder

	// DefaultPermissions Retrieves & stores permissions
	DefaultPermissions permissionServicer

	// DefaultAccessControl Access control checking
	DefaultAccessControl *accessControl

	DefaultNamespace     NamespaceService
	DefaultImportSession ImportSessionService
	DefaultRecord        RecordService
	DefaultModule        ModuleService
	DefaultChart         ChartService
	DefaultPage          PageService
	DefaultAttachment    AttachmentService
	DefaultNotification  *notification

	// DefaultSystemUser is a bridge to users in a system service
	// @todo this is ad-hoc solution that connects compose to system it breaks microservice
	//       architecture and service separation and should be refactored properly
	//       (that is, if we want to continue with microservice architecture)
	DefaultSystemUser systemService.UserService
	//DefaultSystemUser *systemUser
	//DefaultSystemRole *systemRole
)

// Initializes compose-only services
func Initialize(ctx context.Context, log *zap.Logger, s interface{}, c Config) (err error) {
	var (
		hcd = healthcheck.Defaults()

		cmpStore storeInterface
	)

	// we're doing conversion to avoid having
	// store interface exposed or generated inside app package
	cmpStore = s.(storeInterface)

	DefaultLogger = log.Named("service")

	{
		tee := zap.NewNop()
		policy := actionlog.MakeProductionPolicy()

		if !c.ActionLog.Enabled {
			policy = actionlog.MakeDisabledPolicy()
		} else if c.ActionLog.Debug {
			policy = actionlog.MakeDebugPolicy()
			tee = log
		}

		DefaultActionlog = actionlog.NewService(DefaultNgStore, log, tee, policy)
	}

	if DefaultPermissions == nil {
		// Do not override permissions service stored under DefaultPermissions
		// to allow integration tests to inject own permission service
		DefaultPermissions = permissions.Service(ctx, DefaultLogger, cmpStore)
	}

	DefaultAccessControl = AccessControl(DefaultPermissions)

	if DefaultStore == nil {
		const svcPath = "compose"
		if c.Storage.MinioEndpoint != "" {
			var bucket = svcPath
			if c.Storage.MinioBucket != "" {
				bucket = c.Storage.MinioBucket + "/" + svcPath
			}

			DefaultStore, err = minio.New(bucket, minio.Options{
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
			DefaultStore, err = plain.New(path)
			log.Info("initializing store",
				zap.String("path", path),
				zap.Error(err))

		}

		hcd.Add(store.Healthcheck(DefaultStore), "Store/Compose")

		if err != nil {
			return err
		}
	}

	DefaultNamespace = Namespace()
	DefaultModule = Module()

	DefaultSystemUser = systemService.DefaultUser

	DefaultImportSession = ImportSession()
	DefaultRecord = Record()
	DefaultPage = Page()
	DefaultChart = Chart()
	DefaultNotification = Notification()
	DefaultAttachment = Attachment(DefaultStore)

	RegisterIteratorProviders()

	return nil
}

func Activate(ctx context.Context) (err error) {
	return
}

func Watchers(ctx context.Context) {
	// Reloading permissions on change
	DefaultPermissions.Watch(ctx)
}

func RegisterIteratorProviders() {
	// Register resource finders on iterator
	corredor.Service().RegisterIteratorProvider(
		"compose:record",
		func(ctx context.Context, f map[string]string, h eventbus.HandlerFn, action string) error {
			rf := types.RecordFilter{
				Query: f["query"],
				Sort:  f["sort"],
			}

			rf.ParsePagination(f)

			if nsLookup, has := f["namespace"]; !has {
				return errors.New("namespace for record iteration filter not defined")
			} else if ns, err := DefaultNamespace.With(ctx).FindByAny(nsLookup); err != nil {
				return err
			} else {
				rf.NamespaceID = ns.ID
			}

			if mLookup, has := f["module"]; !has {
				return errors.New("module for record iteration filter not defined")
			} else if m, err := DefaultModule.With(ctx).FindByAny(rf.NamespaceID, mLookup); err != nil {
				return err
			} else {
				rf.ModuleID = m.ID
			}

			return DefaultRecord.With(ctx).Iterator(rf, h, action)
		},
	)
}

// Data is stale when new date does not match updatedAt or createdAt (before first update)
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

func nowPtr() *time.Time {
	now := time.Now()
	return &now
}
