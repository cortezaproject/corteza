package agentic

import (
	"testing"

	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestApplyChatbotSectionsReplacesScenariosWholesale(t *testing.T) {
	c := &sysTypes.Chatbot{
		Scenarios: sysTypes.ChatbotScenarios{
			{ID: "one", Type: "conversation", AgentID: 7},
			{ID: "two", Type: "conversation", AgentID: 8},
		},
	}

	require.NoError(t, applyChatbotSections(c, map[string]any{
		"scenarios": `[{"id":"three","type":"conversation","agentID":"9"}]`,
	}))

	// A collection always replaces (§8.2), which is exactly why the tool
	// description tells the caller to send every scenario it wants to keep.
	require.Len(t, c.Scenarios, 1)
	require.Equal(t, "three", c.Scenarios[0].ID)
	require.Equal(t, uint64(9), c.Scenarios[0].AgentID)
}

func TestApplyChatbotSectionsLeavesUnsentSectionsAlone(t *testing.T) {
	c := &sysTypes.Chatbot{
		Handoff:        sysTypes.ChatbotHandoff{Enabled: true},
		Styling:        sysTypes.ChatbotStyling{FontFamily: "Inter"},
		AllowedOrigins: sysTypes.ChatbotAllowedOrigins{"https://example.com"},
	}

	require.NoError(t, applyChatbotSections(c, map[string]any{
		"styling": `{"fontFamily":"Georgia"}`,
	}))

	require.Equal(t, "Georgia", c.Styling.FontFamily)
	require.True(t, c.Handoff.Enabled)
	require.Len(t, c.AllowedOrigins, 1)
}

// allowedOrigins is what actually restricts where the widget may run, so
// clearing it has to be expressible — and has to be distinguishable from not
// mentioning it at all.
func TestApplyChatbotSectionsCanClearAllowedOrigins(t *testing.T) {
	c := &sysTypes.Chatbot{AllowedOrigins: sysTypes.ChatbotAllowedOrigins{"https://example.com"}}

	require.NoError(t, applyChatbotSections(c, map[string]any{"allowedOrigins": `[]`}))
	require.Empty(t, c.AllowedOrigins)

	c.AllowedOrigins = sysTypes.ChatbotAllowedOrigins{"https://example.com"}
	require.NoError(t, applyChatbotSections(c, map[string]any{}))
	require.Len(t, c.AllowedOrigins, 1, "an unmentioned list must survive")
}
