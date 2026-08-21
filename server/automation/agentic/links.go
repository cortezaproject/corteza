package agentic

import (
	autTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/weburl"
)

// Links a single-item result carries, so a caller can say where to go and look
// at what it changed. A builder that cannot resolve an ID returns "", and
// toolkit.JSONResultWith drops the key rather than emitting a blank one.

func taqLinks(taq *autTypes.NgAutomation) map[string]string {
	if taq == nil {
		return nil
	}
	return map[string]string{"url": weburl.TAQ(taq.ID)}
}

func workflowLinks(wf *autTypes.Workflow) map[string]string {
	if wf == nil {
		return nil
	}
	return map[string]string{"url": weburl.Workflow(wf.ID)}
}

// triggerLinks point at the workflow the trigger fires: a trigger has no screen
// of its own, it is edited inside that editor.
func triggerLinks(t *autTypes.Trigger) map[string]string {
	if t == nil {
		return nil
	}
	return map[string]string{"url": weburl.Workflow(t.WorkflowID)}
}
