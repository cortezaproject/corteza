package runtime

import (
	"context"
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
			"# System Context", "## Platform Concepts", "## General Rules",
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

// TL;DR: the platform context follows the agent's tools, and the confirm rule
// follows whether anybody is there to confirm.
// Example: an automation agent granted compose_record_delete is told what a
// module is, but not to ask a person who does not exist before deleting.
func TestBuildSystemPrompt(t *testing.T) {
	const (
		contextHeading = "# System Context"
		denialRule     = "If a tool call is denied, stop and tell the user exactly what was denied."
		confirmRule    = "Never perform a destructive action (delete, bulk delete) without confirming with the user first."
	)

	prompt := func(a *types.Agent, unattended bool) string {
		out, _ := (&runtime{}).buildSystemPrompt(context.Background(), a, nil, unattended)
		return out
	}

	t.Run("an agent granted a tool gets the platform context", func(t *testing.T) {
		require.Contains(t, prompt(agentWith("compose_record_lookup"), false), contextHeading)
	})

	t.Run("an agent granted nothing does not", func(t *testing.T) {
		require.NotContains(t, prompt(&types.Agent{}, false), contextHeading)
	})

	t.Run("every run is told to stop on a denial", func(t *testing.T) {
		for _, a := range []*types.Agent{agentWith("compose_record_lookup"), {}} {
			for _, unattended := range []bool{false, true} {
				require.Contains(t, prompt(a, unattended), denialRule)
			}
		}
	})

	t.Run("an attended run is told to confirm destructive actions", func(t *testing.T) {
		require.Contains(t, prompt(agentWith("compose_record_delete"), false), confirmRule)
		require.Contains(t, prompt(&types.Agent{}, false), confirmRule)
	})

	t.Run("an unattended run is not", func(t *testing.T) {
		require.NotContains(t, prompt(agentWith("compose_record_delete"), true), confirmRule)
	})
}
