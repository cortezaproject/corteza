package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/objstore"
	"github.com/cortezaproject/corteza-server/pkg/objstore/minio"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

type (
	db interface {
		Transaction(callback func() error) error
	}

	RBACServicer interface {
		accessControlRBACServicer
		Watch(ctx context.Context)
	}

	Config struct {
		ActionLog options.ActionLogOpt
		Storage   options.ObjectStoreOpt
	}
)

var (
	DefaultObjectStore objstore.Store

	// DefaultStore is an interface to storage backend(s)
	// ng (next-gen) is a temporary prefix
	// so that we can differentiate between it and the file-only store
	DefaultStore store.Storer

	DefaultLogger *zap.Logger

	DefaultActionlog actionlog.Recorder

	DefaultAccessControl *accessControl

	DefaultAttachment AttachmentService
	DefaultChannel    ChannelService
	DefaultMessage    MessageService
	DefaultEvent      EventService
	DefaultCommand    CommandService

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

func Initialize(ctx context.Context, log *zap.Logger, s store.Storer, c Config) (err error) {
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

	DefaultAccessControl = AccessControl(rbac.Global())

	if DefaultObjectStore == nil {
		const svcPath = "messaging"
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

	hcd.Add(objstore.Healthcheck(DefaultObjectStore), "ObjectStore/Messaging")

	DefaultEvent = Event()
	DefaultChannel = Channel()
	DefaultAttachment = Attachment(DefaultObjectStore)
	DefaultMessage = Message()
	DefaultCommand = Command()

	return nil
}

func Activate(ctx context.Context) (err error) {
	return
}

func Watchers(ctx context.Context) {
	//
}
