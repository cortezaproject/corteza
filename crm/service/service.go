package service

import (
	"sync"
)

type (
	db interface {
		Transaction(callback func() error) error
	}
)

var (
	o                   sync.Once
	DefaultRecord       RecordService
	DefaultModule       ModuleService
	DefaultTrigger      TriggerService
	DefaultChart        ChartService
	DefaultPage         PageService
	DefaultNotification NotificationService
	DefaultPermissions  PermissionsService
)

func Init() {
	o.Do(func() {
		DefaultRecord = Record()
		DefaultModule = Module()
		DefaultTrigger = Trigger()
		DefaultPage = Page()
		DefaultChart = Chart()
		DefaultNotification = Notification()
		DefaultPermissions = Permissions()
	})
}
