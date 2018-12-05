package service

import (
	"sync"
)

var (
	o               sync.Once
	DefaultContent  ContentService
	DefaultField    FieldService
	DefaultModule   ModuleService
	DefaultPage     PageService
	DefaultWorkflow WorkflowService
)

func Init() {
	o.Do(func() {
		DefaultContent = Content()
		DefaultField = Field()
		DefaultModule = Module()
		DefaultPage = Page()
		DefaultWorkflow = Workflow()
	})
}
