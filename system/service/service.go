package service

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/objstore"
	"github.com/cortezaproject/corteza-server/pkg/objstore/minio"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	ngStore "github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	"time"
)

type (
	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}

	Config struct {
		ActionLog options.ActionLogOpt
		Storage   options.ObjectStoreOpt
	}

	permitChecker interface {
		Validate(string, bool) error
		CanCreateUser(uint) error
		CanRegister(uint) error
	}

	eventDispatcher interface {
		WaitFor(ctx context.Context, ev eventbus.Event) (err error)
		Dispatch(ctx context.Context, ev eventbus.Event)
	}
)

var (
	DefaultObjectStore objstore.Store

	// DefaultNgStore is an interface to storage backend(s)
	// ng (next-gen) is a temporary prefix
	// so that we can differentiate between it and the file-only store
	DefaultNgStore ngStore.Storer

	DefaultLogger *zap.Logger

	// CurrentSubscription holds current subscription info,
	// and functions for domain validation, user limit checks and
	// warning texts
	//
	// By default, Corteza (community edition) has this set to nil
	// and with that all checks & validations are skipped
	//
	// Other flavours or distributions can set this to
	// something that suits their needs.
	CurrentSubscription permitChecker

	// DefaultPermissions Retrieves & stores permissions
	DefaultPermissions permissionServicer

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
	DefaultUser        UserService
	DefaultRole        RoleService
	DefaultApplication *application
	DefaultReminder    ReminderService
	DefaultAttachment  AttachmentService

	DefaultStatistics *statistics

	// wrapper around time.Now() that will aid service testing
	now = func() time.Time {
		return time.Now()
	}

	// returns pointer to time.Time struct that is set to current time
	nowPtr = func() *time.Time {
		n := now()
		return &n
	}

	// wrapper around id.Next() that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

func Initialize(ctx context.Context, log *zap.Logger, s ngStore.Storer, c Config) (err error) {
	var (
		hcd = healthcheck.Defaults()
	)

	// we're doing conversion to avoid having
	// store interface exposed or generated inside app package
	DefaultNgStore = s

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
		DefaultPermissions = permissions.Service(ctx, DefaultLogger, DefaultNgStore)
	}

	DefaultAccessControl = AccessControl(DefaultPermissions)

	DefaultSettings = Settings(DefaultNgStore, DefaultLogger, DefaultAccessControl, CurrentSettings)

	if DefaultObjectStore == nil {
		const svcPath = "compose"
		if c.Storage.MinioEndpoint != "" {
			var bucket = svcPath
			if c.Storage.MinioBucket != "" {
				bucket = c.Storage.MinioBucket + "/" + svcPath
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

	DefaultAuthNotification = AuthNotification(CurrentSettings)
	DefaultAuth = Auth()
	DefaultUser = User(ctx)
	DefaultRole = Role(ctx)
	DefaultApplication = Application(DefaultNgStore, DefaultAccessControl, DefaultActionlog, eventbus.Service())
	DefaultReminder = Reminder(ctx)
	DefaultSink = Sink()
	DefaultStatistics = Statistics()
	DefaultAttachment = Attachment(DefaultObjectStore)

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
	// Reloading permissions on change
	DefaultPermissions.Watch(ctx)
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
