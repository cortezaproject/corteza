package service

import (
	"time"

	"go.uber.org/zap"

	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/store"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	DefaultLogger *zap.Logger

	DefaultRecord       RecordService
	DefaultModule       ModuleService
	DefaultTrigger      TriggerService
	DefaultChart        ChartService
	DefaultPage         PageService
	DefaultNotification NotificationService
	DefaultPermissions  PermissionsService
	DefaultAttachment   AttachmentService
	DefaultNamespace    NamespaceService
)

func Init() error {
	fs, err := store.New("var/store")
	if err != nil {
		return err
	}

	DefaultLogger = logger.Default().Named("compose.service")

	DefaultPermissions = Permissions()
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
