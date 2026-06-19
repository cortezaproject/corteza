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
		resourceref.Make(resourceref.KindKnowledgeBase, 901, resourceref.ReasonAgentKnowledgeBase),
		resourceref.Make(resourceref.KindKnowledgeBase, 902, resourceref.ReasonAgentKnowledgeBase),
		resourceref.Make(resourceref.KindLlmProvider, 700, resourceref.ReasonAgentLlmProvider),
		resourceref.Make(resourceref.KindNgAutomation, 501, resourceref.ReasonAgentAutomation),
		resourceref.Make(resourceref.KindAutomationWorkflow, 601, resourceref.ReasonAgentAutomation),
		resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonAgentModule),
		resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonAgentModule),
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
		resourceref.Make(resourceref.KindNgAutomation, 77, resourceref.ReasonChatbotAutomation),
		resourceref.Make(resourceref.KindAgent, 333, resourceref.ReasonChatbotAgent),
		resourceref.Make(resourceref.KindNgAutomation, 88, resourceref.ReasonChatbotAutomation),
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
		resourceref.Make(resourceref.KindComposeModule, 1001, resourceref.ReasonKnowledgeBaseModule),
		resourceref.Make(resourceref.KindComposeModule, 1002, resourceref.ReasonKnowledgeBaseModule),
	}, kb.ResourceRefs())
}
