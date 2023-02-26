package envoy

import (
	"time"

	"github.com/cortezaproject/corteza/server/automation/types"
)

func (e StoreEncoder) setWorkflowDefaults(res *types.Workflow) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateWorkflow(res *types.Workflow) (err error) {
	return
}

func (e StoreEncoder) setTriggerDefaults(res *types.Trigger) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateTrigger(res *types.Trigger) (err error) {
	return
}
