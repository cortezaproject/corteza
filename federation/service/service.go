package service

import (
	"context"
	"time"

	cs "github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/pkg/objstore"
	"github.com/cortezaproject/corteza-server/pkg/objstore/minio"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

type (
	Config struct {
		ActionLog  options.ActionLogOpt
		Storage    options.ObjectStoreOpt
		Federation options.FederationOpt
	}
)

var (
	DefaultObjectStore objstore.Store

	// DefaultNgStore is an interface to storage backend(s)
	// ng (next-gen) is a temporary prefix
	// so that we can differentiate between it and the file-only store
	DefaultStore store.Storer

	DefaultLogger *zap.Logger

	// CurrentSettings represents current system settings
	CurrentSettings = &types.AppSettings{}

	DefaultOptions options.FederationOpt

	DefaultActionlog actionlog.Recorder

	DefaultNode          *node
	DefaultNodeSync      NodeSyncService
	DefaultExposedModule ExposedModuleService
	DefaultSharedModule  SharedModuleService
	DefaultModuleMapping ModuleMappingService

	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now()
		return &c
	}

	// wrapper around id.Next() that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

func Initialize(ctx context.Context, log *zap.Logger, s store.Storer, c Config) (err error) {
	var (
		hcd = healthcheck.Defaults()
	)

	DefaultOptions = c.Federation

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
			tee = log
		}

		DefaultActionlog = actionlog.NewService(DefaultStore, log, tee, policy)
	}

	if DefaultStore == nil {
		const svcPath = "federation"
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

	hcd.Add(objstore.Healthcheck(DefaultObjectStore), "Store/Federation")

	DefaultNode = Node(DefaultStore, service.DefaultUser, DefaultActionlog, auth.DefaultJwtHandler, c.Federation)
	DefaultNodeSync = NodeSync()
	DefaultExposedModule = ExposedModule()
	DefaultSharedModule = SharedModule()
	DefaultModuleMapping = ModuleMapping()

	return
}

func Watchers(ctx context.Context) {
	syncService := NewSync(
		&Syncer{},
		&Mapper{},
		DefaultSharedModule,
		cs.DefaultRecord)

	syncStructure := WorkerStructure(syncService, DefaultLogger)
	syncData := WorkerData(syncService, service.DefaultLogger)

	go syncStructure.Watch(ctx, time.Second*DefaultOptions.StructureMonitorInterval, DefaultOptions.StructurePageSize)
	go syncData.Watch(ctx, time.Second*DefaultOptions.DataMonitorInterval, DefaultOptions.DataPageSize)
}

func AddFederationLabel(entity label.LabeledResource, value string) {
	entity.SetLabel("federation", value)
}
