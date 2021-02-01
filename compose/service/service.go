package service

import (
	"context"
	"fmt"
	automationService "github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/compose/automation"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/objstore"
	"github.com/cortezaproject/corteza-server/pkg/objstore/minio"
	"github.com/cortezaproject/corteza-server/pkg/objstore/plain"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type (
	RBACServicer interface {
		accessControlRBACServicer
		Watch(ctx context.Context)
	}

	Config struct {
		ActionLog options.ActionLogOpt
		Storage   options.ObjectStoreOpt
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

	DefaultActionlog actionlog.Recorder

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

// Initializes compose-only services
func Initialize(ctx context.Context, log *zap.Logger, s store.Storer, c Config) (err error) {
	var (
		hcd = healthcheck.Defaults()
	)

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

		hcd.Add(objstore.Healthcheck(DefaultObjectStore), "ObjectStore/Compose")

		if err != nil {
			return err
		}
	}

	DefaultNamespace = Namespace()
	DefaultModule = Module()

	DefaultImportSession = ImportSession()
	DefaultRecord = Record()
	DefaultPage = Page()
	DefaultChart = Chart()
	DefaultNotification = Notification()
	DefaultAttachment = Attachment(DefaultObjectStore)

	RegisterIteratorProviders()

	automationService.Registry().AddTypes(
		automation.ComposeNamespace{},
		automation.ComposeModule{},
		automation.ComposeRecord{},
		automation.ComposeRecordValues{},
	)

	automation.RecordsHandler(
		automationService.Registry(),
		DefaultNamespace,
		DefaultModule,
		DefaultRecord,
	)

	automation.ModulesHandler(
		automationService.Registry(),
		DefaultNamespace,
		DefaultModule,
	)

	automation.NamespacesHandler(
		automationService.Registry(),
		DefaultNamespace,
	)

	return nil
}

func Activate(ctx context.Context) (err error) {
	return
}

func Watchers(ctx context.Context) {
	//
}

func RegisterIteratorProviders() {
	// Register resource finders on iterator
	corredor.Service().RegisterIteratorProvider(
		"compose:record",
		func(ctx context.Context, f map[string]string, h eventbus.HandlerFn, action string) (err error) {
			rf := types.RecordFilter{Query: f["query"]}
			rf.Sort.Set(f["sort"])

			var limit uint
			if lString, has := f["limit"]; has {
				if limit64, err := strconv.ParseUint(lString, 10, 32); err != nil {
					return fmt.Errorf("can not parse iterator limit param: %w", err)
				} else {
					// We specify the bit size during parsing, so this is fine
					limit = uint(limit64)
				}
			}

			page := f["page"]
			rf.Paging, err = filter.NewPaging(uint(limit), page)
			if err != nil {
				return err
			}

			if nsLookup, has := f["namespace"]; !has {
				return fmt.Errorf("namespace for record iteration filter not defined")
			} else if ns, err := DefaultNamespace.FindByAny(ctx, nsLookup); err != nil {
				return err
			} else {
				rf.NamespaceID = ns.ID
			}

			if mLookup, has := f["module"]; !has {
				return fmt.Errorf("module for record iteration filter not defined")
			} else if m, err := DefaultModule.FindByAny(ctx, rf.NamespaceID, mLookup); err != nil {
				return err
			} else {
				rf.ModuleID = m.ID
			}

			return DefaultRecord.Iterator(ctx, rf, h, action)
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

// trim1st removes 1st param and returns only error
func trim1st(_ interface{}, err error) error {
	return err
}
