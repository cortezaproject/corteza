package store

import (
	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/envoy"
	"github.com/crusttech/human/server/pkg/envoy/resource"
)

func newAutomationWorkflow(wf *types.Workflow, tt types.TriggerSet, ux *userIndex) *automationWorkflow {
	return &automationWorkflow{
		wf: wf,
		tt: tt,

		ux: ux,
	}
}

func (awf *automationWorkflow) MarshalEnvoy() ([]resource.Interface, error) {
	rs := resource.NewAutomationWorkflow(awf.wf)
	syncUserStamps(rs.Userstamps(), awf.ux)

	for _, t := range awf.tt {
		rt := rs.AddAutomationTrigger(t)
		syncUserStamps(rt.Userstamps(), awf.ux)
	}

	for _, s := range awf.wf.Steps {
		rs.AddAutomationWorkflowStep(s)
	}

	for _, p := range awf.wf.Paths {
		rs.AddAutomationWorkflowPath(p)
	}

	return envoy.CollectNodes(
		rs,
	)
}
