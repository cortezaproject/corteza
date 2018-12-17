package service

import (
	"sync"
)

var (
	o               sync.Once
	DefaultContent  ContentService
	DefaultField    FieldService
	DefaultModule   ModuleService
	DefaultChart    ChartService
	DefaultPage     PageService
	DefaultWorkflow WorkflowService
	DefaultJob      JobService
)

func Init() {
	o.Do(func() {
		DefaultContent = Content()
		DefaultField = Field()
		DefaultModule = Module()
		DefaultPage = Page()
		DefaultChart = Chart()
		DefaultWorkflow = Workflow()
		DefaultJob = Job()
	})
}
