package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/internal/repository"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/internal/store"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
)

type (
	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}

	Config struct {
		Storage options.StorageOpt
	}
)

var (
	DefaultPermissions permissionServicer

	DefaultLogger *zap.Logger

	DefaultAccessControl *accessControl

	DefaultRecord       RecordService
	DefaultModule       ModuleService
	DefaultTrigger      TriggerService
	DefaultChart        ChartService
	DefaultPage         PageService
	DefaultNotification NotificationService
	DefaultAttachment   AttachmentService
	DefaultNamespace    NamespaceService
)

func Init(ctx context.Context, log *zap.Logger, c Config) (err error) {
	DefaultLogger = log.Named("service")

	fs, err := store.New(c.Storage.Path)
	log.Info("initializing store", zap.String("path", c.Storage.Path), zap.Error(err))
	if err != nil {
		return err
	}

	DefaultPermissions = permissions.Service(
		ctx,
		DefaultLogger,
		permissions.Repository(repository.DB(ctx), "compose_permission_rules"))

	DefaultAccessControl = AccessControl(DefaultPermissions)

	DefaultRecord = Record()
	DefaultModule = Module()
	DefaultTrigger = Trigger()
	DefaultPage = Page()
	DefaultChart = Chart()
	DefaultNotification = Notification()
	DefaultAttachment = Attachment(fs)
	DefaultNamespace = Namespace()

	return nil
}

func Watchers(ctx context.Context) {
	DefaultPermissions.Watch(ctx)
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
