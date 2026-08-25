package runtime

import (
	"strings"
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func agentWith(names ...string) *types.Agent {
	tools := make([]types.AgentAccessTool, 0, len(names))
	for _, n := range names {
		tools = append(tools, types.AgentAccessTool{Name: n})
	}
	return &types.Agent{Access: types.AgentAccess{Tools: tools}}
}

// TL;DR: an agent that cannot create anything is not told to create things.
// Example: the build guidance says not to tell the user something cannot be
// done when creating a namespace, module or record would do it. Handed to a
// read-only agent that is an instruction to attempt what the policy refuses —
// the denial handler already argues against it — and it buries the author's own
// prompt under a page the agent cannot act on.
func TestSystemContextFor(t *testing.T) {
	const buildHeading = "## Building Data Structures"

	t.Run("a read-only agent does not get the build section", func(t *testing.T) {
		out := systemContextFor(agentWith("compose_record_lookup", "compose_record_report"))
		require.NotContains(t, out, buildHeading)
		require.NotContains(t, out, "build:start")
		require.NotContains(t, out, "build:end")
	})

	t.Run("everything else survives the cut", func(t *testing.T) {
		out := systemContextFor(agentWith("compose_record_lookup"))
		for _, keep := range []string{
			"# System Context", "## Platform Concepts", "## General Rules", "## Response Style",
		} {
			require.Contains(t, out, keep)
		}
		require.NotEmpty(t, strings.TrimSpace(out))
	})

	t.Run("an agent that can create keeps it", func(t *testing.T) {
		out := systemContextFor(agentWith("compose_record_lookup", "compose_record_create"))
		require.Contains(t, out, buildHeading)
		// The guidance itself, not just its heading.
		require.Contains(t, out, "figure out what namespaces, modules, and fields would represent that data")
	})

	t.Run("an agent granted nothing gets no build section", func(t *testing.T) {
		require.NotContains(t, systemContextFor(&types.Agent{}), buildHeading)
		require.NotContains(t, systemContextFor(nil), buildHeading)
	})

	// The markers must never reach a model, whichever branch is taken.
	t.Run("markers never leak", func(t *testing.T) {
		require.NotContains(t, humanSystemContext, "<!-- build:start -->\n<!-- build:end -->")
		for _, a := range []*types.Agent{agentWith("compose_record_create"), agentWith("compose_record_lookup")} {
			out := systemContextFor(a)
			require.NotContains(t, out, "build:end")
		}
	})
}

func TestAgentCanBuild(t *testing.T) {
	require.True(t, agentCanBuild(agentWith("compose_namespace_create")))
	require.True(t, agentCanBuild(agentWith("compose_record_lookup", "system_user_create")))
	require.False(t, agentCanBuild(agentWith("compose_record_lookup", "compose_record_update")))
	require.False(t, agentCanBuild(agentWith()))
	require.False(t, agentCanBuild(nil))
}
