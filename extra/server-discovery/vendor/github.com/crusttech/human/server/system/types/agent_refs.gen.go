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
func (r Agent) ResourceRefs() (out []resourceref.Ref) {
	for _, id := range r.Behavior.KnowledgeBases {
		out = resourceref.Append(out, resourceref.Make(
			resourceref.KindKnowledgeBase, id, resourceref.ReasonAgentKnowledgeBase,
		))
	}
	out = resourceref.Append(out, resourceref.Make(
		resourceref.KindLlmProvider, r.Execution.Model.LLMProviderID, resourceref.ReasonAgentLlmProvider,
	))
	for _, item := range r.Access.TAQs {
		out = resourceref.Append(out, resourceref.Make(
			resourceref.KindNgAutomation, item.ID, resourceref.ReasonAgentAutomation,
		))
	}
	for _, item := range r.Access.Workflows {
		out = resourceref.Append(out, resourceref.Make(
			resourceref.KindAutomationWorkflow, item.ID, resourceref.ReasonAgentAutomation,
		))
	}
	return r.resourceRefsExt(out)
}
