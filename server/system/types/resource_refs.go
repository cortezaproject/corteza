package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/resourceref"
)

// ResourceRefs returns configuration-level references to other resources
func (a Agent) ResourceRefs() (out []resourceref.Ref) {
	for i, id := range a.Behavior.KnowledgeBases {
		out = resourceref.Append(out, resourceref.Make(
			resourceref.KindKnowledgeBase,
			id,
			resourceref.ReasonAgentKnowledgeBase,
			fmt.Sprintf("Behavior.KnowledgeBases.%d", i),
		))
	}

	out = resourceref.Append(out, resourceref.Make(
		resourceref.KindLlmProvider,
		a.Execution.Model.LLMProviderID,
		resourceref.ReasonAgentLlmProvider,
		"Execution.Model.LLMProviderID",
	))

	for i, t := range a.Access.TAQs {
		out = resourceref.Append(out, resourceref.Make(
			resourceref.KindNgAutomation,
			t.ID,
			resourceref.ReasonAgentAutomation,
			fmt.Sprintf("Access.TAQs.%d.ID", i),
		))
	}

	for i, w := range a.Access.Workflows {
		out = resourceref.Append(out, resourceref.Make(
			resourceref.KindAutomationWorkflow,
			w.ID,
			resourceref.ReasonAgentAutomation,
			fmt.Sprintf("Access.Workflows.%d.ID", i),
		))
	}

	for i, tool := range a.Access.Tools {
		for j, allow := range tool.Allow {
			for k, id := range allow.ModuleIDs {
				out = resourceref.Append(out, resourceref.Make(
					resourceref.KindComposeModule,
					id,
					resourceref.ReasonAgentModule,
					fmt.Sprintf("Access.Tools.%d.Allow.%d.ModuleIDs.%d", i, j, k),
				))
			}
		}
	}

	return
}

// ResourceRefs returns configuration-level references to other resources
func (c Chatbot) ResourceRefs() (out []resourceref.Ref) {
	out = resourceref.Append(out,
		automationHookRef(c.Handoff.Automation.OnRequested, "Handoff.Automation.OnRequested.Automation"),
		automationHookRef(c.Handoff.Automation.OnAccepted, "Handoff.Automation.OnAccepted.Automation"),
	)

	for i, s := range c.Scenarios {
		out = resourceref.Append(out,
			resourceref.Make(
				resourceref.KindAgent,
				s.AgentID,
				resourceref.ReasonChatbotAgent,
				fmt.Sprintf("Scenarios.%d.AgentID", i),
			),
			automationHookRef(s.Automation.Before, fmt.Sprintf("Scenarios.%d.Automation.Before.Automation", i)),
			automationHookRef(s.Automation.After, fmt.Sprintf("Scenarios.%d.Automation.After.Automation", i)),
		)
	}

	return
}

// automationHookRef parses hook automation identifiers in the
// "corteza::automation:ng-automation/{ID}" format (see chatbotSession.invokeAutomation)
func automationHookRef(h ChatbotAutomationHook, path string) resourceref.Ref {
	if h.Automation == "" {
		return resourceref.Ref{}
	}

	id, err := strconv.ParseUint(strings.TrimPrefix(h.Automation, resourceref.KindNgAutomation+"/"), 10, 64)
	if err != nil {
		return resourceref.Ref{}
	}

	return resourceref.Make(resourceref.KindNgAutomation, id, resourceref.ReasonChatbotAutomation, path)
}

// ResourceRefs returns configuration-level references to other resources
func (kb KnowledgeBase) ResourceRefs() (out []resourceref.Ref) {
	if kb.Context == nil {
		return
	}

	for i, ns := range kb.Context.Namespaces {
		for j, id := range ns.ModuleIDs {
			out = resourceref.Append(out, resourceref.Make(
				resourceref.KindComposeModule,
				id,
				resourceref.ReasonKnowledgeBaseModule,
				fmt.Sprintf("Context.Namespaces.%d.ModuleIDs.%d", i, j),
			))
		}
	}

	return
}
