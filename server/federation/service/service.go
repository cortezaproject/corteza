package service

import (
	"context"
	cs "github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/store"
	ss "github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
	"time"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	Config struct {
		ActionLog  options.ActionLogOpt
		Storage    options.ObjectStoreOpt
		Federation options.FederationOpt
		Server     options.HttpServerOpt
	}
)

var (
	DefaultStore store.Storer

	DefaultLogger *zap.Logger

	DefaultAccessControl *accessControl

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
		c := time.Now().Round(time.Second)
		return &c
	}

	// wrapper around id.Next() that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

func isStale(new *time.Time, updatedAt *time.Time, createdAt time.Time) bool {
	if new == nil {
		return false
	}
	if updatedAt != nil {
		return !new.Equal(*updatedAt)
	}
	return !new.Equal(createdAt)
}

func Initialize(_ context.Context, log *zap.Logger, s store.Storer, c Config) (err error) {
	DefaultOptions = c.Federation

	// we're doing conversion to avoid having
	// store interface exposed or generated inside app package
	DefaultStore = s

	DefaultLogger = log.Named("service")

	// @todo add healthcheck(s) for federation sync services

	{
		tee := zap.NewNop()
		policy := actionlog.MakeProductionPolicy()

		if !c.ActionLog.Enabled {
			policy = actionlog.MakeDisabledPolicy()
		} else if c.ActionLog.Debug {
			policy = actionlog.MakeDebugPolicy()
			tee = logger.MakeDebugLogger()
		}

		// @todo allow configuring a separate store DSN for the actionlog
		DefaultActionlog = actionlog.NewService(DefaultStore, log, tee, policy)
	}

	DefaultAccessControl = AccessControl(s)

	DefaultNode = Node(
		DefaultStore,
		ss.DefaultUser,
		DefaultActionlog,
		func(ctx context.Context, i auth.Identifiable) (token []byte, err error) {
			return auth.TokenIssuer.Issue(
				ctx,
				auth.WithIdentity(i),
				auth.WithScope("api"),
				auth.WithAudience("federation"),
			)
		},
		c.Federation,
		c.Server,
		DefaultAccessControl,
	)
	DefaultNodeSync = NodeSync()
	DefaultExposedModule = ExposedModule()
	DefaultSharedModule = SharedModule()
	DefaultModuleMapping = ModuleMapping()

	return
}

func Watchers(ctx context.Context) {
	DefaultLogger.Info("Starting federation - warning, this is still an experimental feature")

	syncService := NewSync(
		&Syncer{},
		&Mapper{},
		DefaultSharedModule,
		cs.DefaultRecord,
		ss.DefaultUser,
		ss.DefaultRole)

	syncStructure := WorkerStructure(syncService, DefaultLogger)
	syncData := WorkerData(syncService, DefaultLogger)

	go syncStructure.Watch(
		ctx,
		DefaultOptions.StructureMonitorInterval,
		DefaultOptions.StructurePageSize)

	go syncData.Watch(
		ctx,
		DefaultOptions.DataMonitorInterval,
		DefaultOptions.DataPageSize)
}

func AddFederationLabel(entity label.LabeledResource, key string, value string) {
	entity.SetLabel(key, labelTypes.LabelValue{Val: value})
}
