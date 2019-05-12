package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/crusttech/crust/compose/internal/repository"
	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/permissions"
	"github.com/crusttech/crust/internal/store"
)

type (
	db interface {
		Transaction(callback func() error) error
	}

	permissionServicer interface {
		accessControlPermissionServicer
		Watch(ctx context.Context)
	}
)

var (
	permSvc permissionServicer

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

func Init(ctx context.Context) error {
	DefaultLogger = logger.Default().Named("compose.service")

	fs, err := store.New("var/store")
	if err != nil {
		return err
	}

	permSvc = permissions.Service(
		ctx,
		DefaultLogger,
		permissions.Repository(repository.DB(ctx), "compose_permission_rules"))

	DefaultAccessControl = AccessControl(permSvc)

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
	permSvc.Watch(ctx)
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
