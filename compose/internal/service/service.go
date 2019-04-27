package service

import (
	"github.com/crusttech/crust/internal/store"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	DefaultRecord       RecordService
	DefaultModule       ModuleService
	DefaultTrigger      TriggerService
	DefaultChart        ChartService
	DefaultPage         PageService
	DefaultNotification NotificationService
	DefaultPermissions  PermissionsService
	DefaultAttachment   AttachmentService
)

func Init() error {
	fs, err := store.New("var/store")
	if err != nil {
		return err
	}

	DefaultPermissions = Permissions()
	DefaultRecord = Record()
	DefaultModule = Module()
	DefaultTrigger = Trigger()
	DefaultPage = Page()
	DefaultChart = Chart()
	DefaultNotification = Notification()
	DefaultAttachment = Attachment(fs)

	return nil
}
