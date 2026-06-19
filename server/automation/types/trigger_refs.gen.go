package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/crusttech/human/server/pkg/resourceref"
)

// ResourceRefs returns configuration-level references to other resources
func (r Trigger) ResourceRefs() (out []resourceref.Ref) {
	out = resourceref.Append(out, resourceref.Make(
		resourceref.KindAutomationWorkflow, r.WorkflowID, resourceref.ReasonTriggerWorkflow,
	))
	return r.resourceRefsExt(out)
}
