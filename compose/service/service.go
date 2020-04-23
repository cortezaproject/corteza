package service

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/app/options"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/pkg/store"
	"github.com/cortezaproject/corteza-server/pkg/store/minio"
	"github.com/cortezaproject/corteza-server/pkg/store/plain"
	systemProto "github.com/cortezaproject/corteza-server/system/proto"
)

type (
	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}

	Config struct {
		Storage          options.StorageOpt
		GRPCClientSystem options.GRPCServerOpt
	}

	eventDispatcher interface {
		WaitFor(ctx context.Context, ev eventbus.Event) (err error)
		Dispatch(ctx context.Context, ev eventbus.Event)
	}
)

var (
	DefaultStore store.Store

	DefaultLogger *zap.Logger

	DefaultSettings settings.Service

	// DefaultPermissions Retrieves & stores permissions
	DefaultPermissions permissionServicer

	// DefaultAccessControl Access control checking
	DefaultAccessControl *accessControl

	// CurrentSettings represents current compose settings
	CurrentSettings = &types.Settings{}

	DefaultNamespace     NamespaceService
	DefaultImportSession ImportSessionService
	DefaultRecord        RecordService
	DefaultModule        ModuleService
	DefaultChart         ChartService
	DefaultPage          PageService
	DefaultAttachment    AttachmentService
	DefaultNotification  *notification

	DefaultSystemUser *systemUser
	DefaultSystemRole *systemRole
)

// Initializes compose-only services
func Initialize(ctx context.Context, log *zap.Logger, c Config) (err error) {
	var db = repository.DB(ctx)

	DefaultLogger = log.Named("service")

	if DefaultPermissions == nil {
		// Do not override permissions service stored under DefaultPermissions
		// to allow integration tests to inject own permission service
		DefaultPermissions = permissions.Service(ctx, DefaultLogger, db, "compose_permission_rules")
	}

	DefaultAccessControl = AccessControl(DefaultPermissions)

	DefaultSettings = settings.NewService(
		settings.NewRepository(repository.DB(ctx), "compose_settings"),
		DefaultLogger,
		DefaultAccessControl,
		CurrentSettings,
	)

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

		if err != nil {
			return err
		}
	}

	DefaultNamespace = Namespace()
	DefaultModule = Module()

	{
		systemClientConn, err := NewSystemGRPCClient(ctx, c.GRPCClientSystem, DefaultLogger)
		if err != nil {
			return err
		}

		DefaultSystemUser = SystemUser(systemProto.NewUsersClient(systemClientConn))
		DefaultSystemRole = SystemRole(systemProto.NewRolesClient(systemClientConn))
	}

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
	// Run initial update of current settings with super-user credentials
	err = DefaultSettings.UpdateCurrent(auth.SetSuperUserContext(ctx))
	if err != nil {
		return
	}

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
				Filter: f["query"],
				Sort:   f["sort"],
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
