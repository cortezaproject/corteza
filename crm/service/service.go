package service

import (
	"sync"
)

var (
	o              sync.Once
	DefaultRecord  RecordService
	DefaultModule  ModuleService
	DefaultTrigger TriggerService
	DefaultChart   ChartService
	DefaultPage    PageService
)

func Init() {
	o.Do(func() {
		DefaultRecord = Record()
		DefaultModule = Module()
		DefaultTrigger = Trigger()
		DefaultPage = Page()
		DefaultChart = Chart()
	})
}
