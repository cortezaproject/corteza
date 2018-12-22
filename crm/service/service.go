package service

import (
	"sync"
)

var (
	o               sync.Once
	DefaultRecord   RecordService
	DefaultField    FieldService
	DefaultModule   ModuleService
	DefaultChart    ChartService
	DefaultPage     PageService
	DefaultWorkflow WorkflowService
	DefaultJob      JobService
)

func Init() {
	o.Do(func() {
		DefaultRecord = Record()
		DefaultField = Field()
		DefaultModule = Module()
		DefaultPage = Page()
		DefaultChart = Chart()
		DefaultWorkflow = Workflow()
		DefaultJob = Job()
	})
}
