package agentic

import (
	"testing"

	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/require"
)

// fakeRegistry stands in for the MCP registry with a fixed set of tool names.
type fakeRegistry struct{ names []string }

func (f fakeRegistry) RegisterTool(mcp.Tool, string, server.ToolHandlerFunc, ...hmcp.RegisterOption) {
}

func (f fakeRegistry) HasTool(name string) bool {
	for _, n := range f.names {
		if n == name {
			return true
		}
	}
	return false
}

func (f fakeRegistry) Tools() []mcp.Tool {
	out := make([]mcp.Tool, 0, len(f.names))
	for _, n := range f.names {
		out = append(out, mcp.Tool{Name: n})
	}
	return out
}

func handlerWith(names ...string) *agentHandler {
	return &agentHandler{reg: fakeRegistry{names: names}}
}

func agentWithTools(names ...string) *sysTypes.Agent {
	tools := make([]sysTypes.AgentAccessTool, 0, len(names))
	for _, n := range names {
		tools = append(tools, sysTypes.AgentAccessTool{Name: n})
	}
	return &sysTypes.Agent{Access: sysTypes.AgentAccess{Tools: tools}}
}

var registered = []string{
	"compose_record_lookup", "compose_record_create", "compose_record_update",
	"compose_module_lookup", "compose_namespace_lookup", "compose_chart_lookup",
	"system_user_lookup",
}

// TL;DR: a tool name the registry does not know is refused where it is written.
// Example: an agent granted "compose_record_lookkup" stored happily and then
// failed to run at all — the runtime resolves the whole allow-list at once, so
// one typo costs every tool, and the typo's description still told the model
// the capability was there.
func TestValidateAgentTools_UnknownName(t *testing.T) {
	h := handlerWith(registered...)

	t.Run("accepts registered names", func(t *testing.T) {
		require.NoError(t, h.validateAgentTools(agentWithTools("compose_record_lookup", "system_user_lookup")))
	})

	t.Run("refuses a typo and suggests the real name", func(t *testing.T) {
		err := h.validateAgentTools(agentWithTools("compose_record_lookkup"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "compose_record_lookkup")
		require.Contains(t, err.Error(), `"compose_record_lookup"`)
	})

	t.Run("refuses a name with no near match", func(t *testing.T) {
		err := h.validateAgentTools(agentWithTools("utterly_made_up_tool"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "utterly_made_up_tool")
		require.NotContains(t, err.Error(), "Did you mean")
	})

	t.Run("names the offending entry", func(t *testing.T) {
		err := h.validateAgentTools(agentWithTools("compose_record_lookup", "nope_not_here"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "access.tools[1]")
	})

	t.Run("refuses an empty name", func(t *testing.T) {
		err := h.validateAgentTools(agentWithTools(""))
		require.Error(t, err)
		require.Contains(t, err.Error(), "name")
	})

	t.Run("an agent granting nothing is not an error", func(t *testing.T) {
		require.NoError(t, h.validateAgentTools(&sysTypes.Agent{}))
	})

	// The registrar every other agentic handler is built with only registers.
	// Validation has to stand down rather than refuse everything.
	t.Run("a registrar that cannot be asked skips the check", func(t *testing.T) {
		bare := &agentHandler{reg: registrarOnly{}}
		require.NoError(t, bare.validateAgentTools(agentWithTools("anything_at_all")))
	})
}

type registrarOnly struct{}

func (registrarOnly) RegisterTool(mcp.Tool, string, server.ToolHandlerFunc, ...hmcp.RegisterOption) {
}

// TL;DR: the suggestion is the name the author was reaching for.
func TestNearestToolNames(t *testing.T) {
	tools := make([]mcp.Tool, 0, len(registered))
	for _, n := range registered {
		tools = append(tools, mcp.Tool{Name: n})
	}

	cases := map[string]string{
		"compose_record_lookkup": "compose_record_lookup",
		"compose_record_lookup ": "compose_record_lookup",
		"COMPOSE_RECORD_LOOKUP":  "compose_record_lookup",
		"compose_module_lookups": "compose_module_lookup",
		"system_user_lookop":     "system_user_lookup",
	}
	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got := nearestToolNames(in, tools)
			require.NotEmpty(t, got)
			require.Equal(t, want, got[0])
		})
	}

	t.Run("nothing close returns nothing", func(t *testing.T) {
		require.Empty(t, nearestToolNames("zzzzzzzzzzzzzzzz", tools))
	})

	t.Run("at most three", func(t *testing.T) {
		require.LessOrEqual(t, len(nearestToolNames("compose_record_looku", tools)), 3)
	})
}
