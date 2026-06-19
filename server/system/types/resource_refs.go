package types

import (
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/resourceref"
)

// resourceRefsExt extends the generated Agent.ResourceRefs() with the
// Access.Tools triple-nested loop that cannot be expressed as a simple ref.
func (a Agent) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	for _, tool := range a.Access.Tools {
		for _, allow := range tool.Allow {
			for _, id := range allow.ModuleIDs {
				out = resourceref.Append(out, resourceref.Make(
					resourceref.KindComposeModule,
					id,
					resourceref.ReasonAgentModule,
				))
			}
		}
	}

	return out
}

// resourceRefsExt extends the generated Chatbot.ResourceRefs() with the
// handoff + scenario refs; preserve the exact emission order (handoff hooks
// first, then per scenario: agent, before, after) — a test asserts it.
func (c Chatbot) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	out = resourceref.Append(out,
		automationHookRef(c.Handoff.Automation.OnRequested),
		automationHookRef(c.Handoff.Automation.OnAccepted),
	)

	for _, s := range c.Scenarios {
		out = resourceref.Append(out,
			resourceref.Make(
				resourceref.KindAgent,
				s.AgentID,
				resourceref.ReasonChatbotAgent,
			),
			automationHookRef(s.Automation.Before),
			automationHookRef(s.Automation.After),
		)
	}

	return out
}

// automationHookRef parses hook automation identifiers in the
// "corteza::automation:ng-automation/{ID}" format (see chatbotSession.invokeAutomation)
func automationHookRef(h ChatbotAutomationHook) resourceref.Ref {
	if h.Automation == "" {
		return resourceref.Ref{}
	}

	id, err := strconv.ParseUint(strings.TrimPrefix(h.Automation, resourceref.KindNgAutomation+"/"), 10, 64)
	if err != nil {
		return resourceref.Ref{}
	}

	return resourceref.Make(resourceref.KindNgAutomation, id, resourceref.ReasonChatbotAutomation)
}

// resourceRefsExt extends the generated KnowledgeBase.ResourceRefs() with the
// context module refs.
func (kb KnowledgeBase) resourceRefsExt(out []resourceref.Ref) []resourceref.Ref {
	if kb.Context == nil {
		return out
	}

	for _, ns := range kb.Context.Namespaces {
		for _, id := range ns.ModuleIDs {
			out = resourceref.Append(out, resourceref.Make(
				resourceref.KindComposeModule,
				id,
				resourceref.ReasonKnowledgeBaseModule,
			))
		}
	}

	return out
}
