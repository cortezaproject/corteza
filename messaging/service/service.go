package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/store"
	"github.com/cortezaproject/corteza-server/pkg/store/minio"
	"github.com/cortezaproject/corteza-server/pkg/store/plain"
)

type (
	db interface {
		Transaction(callback func() error) error
	}

	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}

	Config struct {
		ActionLog options.ActionLogOpt
		Storage   options.StorageOpt
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

	DefaultPermissions permissionServicer

	DefaultLogger *zap.Logger

	DefaultActionlog actionlog.Recorder

	DefaultAccessControl *accessControl

	DefaultAttachment AttachmentService
	DefaultChannel    ChannelService
	DefaultMessage    MessageService
	DefaultEvent      EventService
	DefaultCommand    CommandService
)

func Initialize(ctx context.Context, log *zap.Logger, s interface{}, c Config) (err error) {
	var (
		hcd = healthcheck.Defaults()

		msgStore storeInterface
		ok       bool
	)

	// we're doing conversion to avoid having
	// store interface exposed or generated inside app package
	if msgStore, ok = s.(storeInterface); !ok {
		return fmt.Errorf("store %T is incompatible with compose storeInterface", s)
	}

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
		DefaultPermissions = permissions.Service(ctx, DefaultLogger, msgStore)
	}

	DefaultAccessControl = AccessControl(DefaultPermissions)

	if DefaultStore == nil {
		const svcPath = "messaging"
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

	hcd.Add(store.Healthcheck(DefaultStore), "Store/Messaging")

	DefaultEvent = Event(ctx)
	DefaultChannel = Channel(ctx)
	DefaultAttachment = Attachment(ctx, DefaultStore)
	DefaultMessage = Message(ctx)
	DefaultCommand = Command(ctx)

	return nil
}

func Activate(ctx context.Context) (err error) {
	return
}

func Watchers(ctx context.Context) {
	DefaultPermissions.Watch(ctx)
}

func timeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
