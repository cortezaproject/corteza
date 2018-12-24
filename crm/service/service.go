package service

import (
	"sync"
)

var (
	o               sync.Once
	DefaultRecord   RecordService
	DefaultModule   ModuleService
	DefaultChart    ChartService
	DefaultPage     PageService
	DefaultWorkflow WorkflowService
	DefaultJob      JobService
)

func Init() {
	o.Do(func() {
		DefaultRecord = Record()
		DefaultModule = Module()
		DefaultPage = Page()
		DefaultChart = Chart()
		DefaultWorkflow = Workflow()
		DefaultJob = Job()
	})
}
