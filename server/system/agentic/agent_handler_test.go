package agentic

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/auth"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

// The section contract (§8.2): a section the caller sent REPLACES that section,
// and a section absent from the arguments is left exactly as it was. Both
// halves matter — the first is what the descriptions promise, and the second is
// what makes a one-section update safe on a service that writes whole records.
func TestApplyAgentSectionsReplacesWhatIsSent(t *testing.T) {
	a := &sysTypes.Agent{
		Meta:     sysTypes.AgentMeta{Short: "Old", Description: "Old description"},
		Behavior: sysTypes.AgentBehavior{SystemPrompt: "old prompt"},
	}

	require.NoError(t, applyAgentSections(a, map[string]any{
		"meta": `{"short":"New"}`,
	}))

	require.Equal(t, "New", a.Meta.Short)

	// The section replaces, so a field the new object omits is gone. Decoding
	// into the populated struct would have merged and kept the old description,
	// which is a behaviour nothing documents.
	require.Empty(t, a.Meta.Description, "meta merged instead of replacing")

	// An untouched section is untouched.
	require.Equal(t, "old prompt", a.Behavior.SystemPrompt)
}

func TestApplyAgentSectionsAcceptsADecodedObject(t *testing.T) {
	a := &sysTypes.Agent{}

	// Clients differ on whether a JSON param arrives as a string or already
	// decoded; both have to work.
	require.NoError(t, applyAgentSections(a, map[string]any{
		"execution": map[string]any{
			"model": map[string]any{"model": "claude-sonnet-5"},
		},
	}))

	require.Equal(t, "claude-sonnet-5", a.Execution.Model.Model)
}

func TestApplyAgentSectionsLeavesEverythingAloneWhenNothingIsSent(t *testing.T) {
	a := &sysTypes.Agent{
		Meta:       sysTypes.AgentMeta{Short: "Kept"},
		Access:     sysTypes.AgentAccess{Tools: []sysTypes.AgentAccessTool{{Name: "compose_record_lookup"}}},
		Invocation: sysTypes.AgentInvocation{User: sysTypes.AgentInvocationUser{Enabled: true}},
	}

	require.NoError(t, applyAgentSections(a, map[string]any{"handle": "unrelated"}))

	require.Equal(t, "Kept", a.Meta.Short)
	require.Len(t, a.Access.Tools, 1)
	require.True(t, a.Invocation.User.Enabled)
}

func TestApplyAgentSectionsRejectsMalformedJSON(t *testing.T) {
	err := applyAgentSections(&sysTypes.Agent{}, map[string]any{"behavior": `{"systemPrompt":`})
	require.Error(t, err)
	require.Contains(t, err.Error(), "agent behavior")
}

// An empty string is not an empty object: it means the caller sent nothing, so
// the section keeps what it had rather than being wiped.
func TestApplyAgentSectionsTreatsAnEmptyStringAsAbsent(t *testing.T) {
	a := &sysTypes.Agent{Meta: sysTypes.AgentMeta{Short: "Kept"}}
	require.NoError(t, applyAgentSections(a, map[string]any{"meta": "  "}))
	require.Equal(t, "Kept", a.Meta.Short)
}

// An agent must not be able to start another agent through system_agent_exec.
// Nothing in the runtime bounds recursion, and the tool sits in the usage group
// at write risk, so a group grant would otherwise hand every agent the ability
// to run agents — itself included.
func TestAgentExecRefusesAnAgentInvocation(t *testing.T) {
	ctx := auth.SetIdentityToContext(context.Background(),
		auth.IdentityWithAgent(auth.Authenticated(42), 511))

	h := &agentHandler{}
	_, err := h.exec(ctx, mcp.CallToolRequest{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "agent cannot start another agent")
	require.Contains(t, err.Error(), "511")
}

// The guard must not fire for an ordinary user, which is every call over the
// MCP transport. It falls through to argument validation instead.
func TestAgentExecAllowsAPlainUser(t *testing.T) {
	ctx := auth.SetIdentityToContext(context.Background(), auth.Authenticated(42))

	h := &agentHandler{}
	_, err := h.exec(ctx, mcp.CallToolRequest{})

	require.Error(t, err)
	require.NotContains(t, err.Error(), "agent cannot start another agent")
}
