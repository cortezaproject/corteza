package service

import (
	"sync"
)

var (
	o                   sync.Once
	DefaultRecord       RecordService
	DefaultModule       ModuleService
	DefaultTrigger      TriggerService
	DefaultChart        ChartService
	DefaultPage         PageService
	DefaultNotification NotificationService
)

func Init() {
	o.Do(func() {
		DefaultRecord = Record()
		DefaultModule = Module()
		DefaultTrigger = Trigger()
		DefaultPage = Page()
		DefaultChart = Chart()
		DefaultNotification = Notification()
	})
}
