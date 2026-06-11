package types

import (
	"testing"

	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/stretchr/testify/require"
)

func TestAgentResourceRefs(t *testing.T) {
	a := Agent{
		Behavior: AgentBehavior{
			KnowledgeBases: KnowledgeBaseIDList{901, 902},
		},
		Execution: AgentExecution{
			Model: AgentExecutionModel{LLMProviderID: 700},
		},
		Access: AgentAccess{
			TAQs:      []AgentAccessTAQ{{ID: 501}},
			Workflows: []AgentAccessWorkflow{{ID: 601}},
			Tools: []AgentAccessTool{
				{Allow: []AgentAccessAllow{{NamespaceID: 1, ModuleIDs: AgentAccessIDList{1001, 1002}}}},
			},
		},
	}

	require.Equal(t, []resourceref.Ref{
		{Kind: resourceref.KindKnowledgeBase, ID: 901, Reason: resourceref.ReasonAgentKnowledgeBase, Path: "Behavior.KnowledgeBases.0"},
		{Kind: resourceref.KindKnowledgeBase, ID: 902, Reason: resourceref.ReasonAgentKnowledgeBase, Path: "Behavior.KnowledgeBases.1"},
		{Kind: resourceref.KindLlmProvider, ID: 700, Reason: resourceref.ReasonAgentLlmProvider, Path: "Execution.Model.LLMProviderID"},
		{Kind: resourceref.KindNgAutomation, ID: 501, Reason: resourceref.ReasonAgentAutomation, Path: "Access.TAQs.0.ID"},
		{Kind: resourceref.KindAutomationWorkflow, ID: 601, Reason: resourceref.ReasonAgentAutomation, Path: "Access.Workflows.0.ID"},
		{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonAgentModule, Path: "Access.Tools.0.Allow.0.ModuleIDs.0"},
		{Kind: resourceref.KindComposeModule, ID: 1002, Reason: resourceref.ReasonAgentModule, Path: "Access.Tools.0.Allow.0.ModuleIDs.1"},
	}, a.ResourceRefs())
}

func TestAgentResourceRefsEmpty(t *testing.T) {
	require.Empty(t, Agent{}.ResourceRefs())
}

func TestChatbotResourceRefs(t *testing.T) {
	c := Chatbot{
		Handoff: ChatbotHandoff{
			Automation: ChatbotHandoffAutomation{
				OnRequested: ChatbotAutomationHook{Automation: "corteza::automation:ng-automation/77"},
			},
		},
		Scenarios: ChatbotScenarios{
			{
				AgentID: 333,
				Automation: ChatbotScenarioAutomation{
					Before: ChatbotAutomationHook{Automation: "corteza::automation:ng-automation/88"},
					After:  ChatbotAutomationHook{Automation: "not-a-valid-ref"},
				},
			},
		},
	}

	require.Equal(t, []resourceref.Ref{
		{Kind: resourceref.KindNgAutomation, ID: 77, Reason: resourceref.ReasonChatbotAutomation, Path: "Handoff.Automation.OnRequested.Automation"},
		{Kind: resourceref.KindAgent, ID: 333, Reason: resourceref.ReasonChatbotAgent, Path: "Scenarios.0.AgentID"},
		{Kind: resourceref.KindNgAutomation, ID: 88, Reason: resourceref.ReasonChatbotAutomation, Path: "Scenarios.0.Automation.Before.Automation"},
	}, c.ResourceRefs())
}

func TestKnowledgeBaseResourceRefs(t *testing.T) {
	require.Empty(t, KnowledgeBase{}.ResourceRefs())

	kb := KnowledgeBase{
		Context: &KnowledgeBaseContext{
			Namespaces: []KnowledgeBaseNamespaceContext{
				{NamespaceID: 1, ModuleIDs: KnowledgeBaseIDList{1001, 1002}},
			},
		},
	}

	require.Equal(t, []resourceref.Ref{
		{Kind: resourceref.KindComposeModule, ID: 1001, Reason: resourceref.ReasonKnowledgeBaseModule, Path: "Context.Namespaces.0.ModuleIDs.0"},
		{Kind: resourceref.KindComposeModule, ID: 1002, Reason: resourceref.ReasonKnowledgeBaseModule, Path: "Context.Namespaces.0.ModuleIDs.1"},
	}, kb.ResourceRefs())
}
