package runtime

import (
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

// grantRegistry is the shared mock, loaded with a small tool taxonomy.
func grantRegistry() *mockMCP {
	return &mockMCP{inGroup: map[string][]string{
		"usage/read":        {"compose_record_lookup", "compose_record_report"},
		"usage/write":       {"compose_record_create", "compose_record_lookup", "compose_record_report"},
		"configuring/read":  {"compose_module_lookup"},
		"configuring/write": {"compose_module_create", "compose_module_lookup"},
	}}
}

// grantingAgent builds an agent whose access is exactly the given entries.
func grantingAgent(tools ...types.AgentAccessTool) *types.Agent {
	return &types.Agent{Access: types.AgentAccess{Tools: tools}}
}

func namesOf(a *types.Agent) []string {
	out := make([]string, 0, len(a.Access.Tools))
	for _, t := range a.Access.Tools {
		out = append(out, t.Name)
	}
	return out
}

// TL;DR: a group grant becomes the tools it stands for, before anything else
// looks at the allow-list.
// Example: granting an agent read access meant naming five tools and repeating
// the scope on each; adding a tool later meant editing every agent.
func TestExpandToolGrants(t *testing.T) {
	scope := []types.AgentAccessAllow{{NamespaceID: 100}}

	t.Run("a group becomes its members", func(t *testing.T) {
		reg := grantRegistry()
		out := expandToolGrants(grantingAgent(
			types.AgentAccessTool{Group: "usage", MaxRisk: "read", Allow: scope},
		), reg)

		require.Equal(t, []string{"compose_record_lookup", "compose_record_report"}, namesOf(out))
	})

	t.Run("the group's scope reaches every member", func(t *testing.T) {
		out := expandToolGrants(grantingAgent(
			types.AgentAccessTool{Group: "usage", MaxRisk: "read", Allow: scope},
		), grantRegistry())

		for _, e := range out.Access.Tools {
			require.Equal(t, scope, e.Allow, "%s lost the grant's scope", e.Name)
			require.Empty(t, e.Group, "the expanded entry must read as an ordinary one")
			require.Empty(t, e.MaxRisk)
		}
	})

	// The ceiling is the point of the feature; it must actually be passed on.
	t.Run("maxRisk selects the set", func(t *testing.T) {
		reg := grantRegistry()
		out := expandToolGrants(grantingAgent(
			types.AgentAccessTool{Group: "usage", MaxRisk: "write", Allow: scope},
		), reg)
		require.Contains(t, namesOf(out), "compose_record_create")
		require.Equal(t, "usage/write", reg.calls[0], "the grant is expanded at the stated ceiling")
	})

	t.Run("an unstated risk permits only reading", func(t *testing.T) {
		reg := grantRegistry()
		expandToolGrants(grantingAgent(
			types.AgentAccessTool{Group: "usage", Allow: scope},
		), reg)
		require.Equal(t, "usage/read", reg.calls[0], "the grant itself is expanded at the read ceiling")
	})

	t.Run("named tools are kept alongside", func(t *testing.T) {
		out := expandToolGrants(grantingAgent(
			types.AgentAccessTool{Name: "system_user_lookup", Allow: scope},
			types.AgentAccessTool{Group: "configuring", MaxRisk: "read", Allow: scope},
		), grantRegistry())
		require.Equal(t, []string{"system_user_lookup", "compose_module_lookup"}, namesOf(out))
	})

	// Overlapping grants are ordinary: two groups, or a group and a name.
	t.Run("a tool granted twice appears once", func(t *testing.T) {
		out := expandToolGrants(grantingAgent(
			types.AgentAccessTool{Name: "compose_record_lookup", Allow: scope},
			types.AgentAccessTool{Group: "usage", MaxRisk: "read", Allow: scope},
		), grantRegistry())
		require.Equal(t, []string{"compose_record_lookup", "compose_record_report"}, namesOf(out))
	})

	// The agent comes from a shared registry; widening the caller's copy would
	// leak one request's grant into the next.
	t.Run("the caller's agent is not modified", func(t *testing.T) {
		in := grantingAgent(types.AgentAccessTool{Group: "usage", MaxRisk: "read", Allow: scope})
		out := expandToolGrants(in, grantRegistry())
		require.Len(t, in.Access.Tools, 1)
		require.Equal(t, "usage", in.Access.Tools[0].Group)
		require.NotSame(t, in, out)
	})

	t.Run("nothing to expand still gets its modes resolved", func(t *testing.T) {
		in := grantingAgent(types.AgentAccessTool{Name: "compose_record_lookup", Allow: scope})
		out := expandToolGrants(in, grantRegistry())
		require.NotSame(t, in, out, "the mode is filled in on a copy")
		require.Equal(t, "always", out.Access.Tools[0].Permission)
		require.Empty(t, in.Access.Tools[0].Permission, "the caller's agent is left alone")

		require.Nil(t, expandToolGrants(nil, grantRegistry()))
	})

	t.Run("a grant that states its mode keeps it", func(t *testing.T) {
		in := grantingAgent(types.AgentAccessTool{Name: "compose_record_lookup", Permission: "ask", Allow: scope})
		require.Same(t, in, expandToolGrants(in, grantRegistry()), "nothing to resolve, nothing to copy")
	})

	t.Run("a group nothing matches yields nothing, not everything", func(t *testing.T) {
		out := expandToolGrants(grantingAgent(
			types.AgentAccessTool{Group: "nonsense", MaxRisk: "read", Allow: scope},
		), grantRegistry())
		require.Empty(t, out.Access.Tools)
	})
}

// An agent that names no tools inherits what the invoking user can already do.
// Listing tools by hand was the single thing that made an agent tedious to set
// up; the agent runs as that user, so RBAC bounds it either way.
func TestExpandToolGrantsInherits(t *testing.T) {
	reg := grantRegistry()
	out := expandToolGrants(&types.Agent{}, reg)

	require.NotEmpty(t, out.Access.Tools)
	require.Contains(t, namesOf(out), "compose_record_create", "writing is inherited")
	require.Contains(t, namesOf(out), "compose_module_lookup", "both groups are inherited")

	// Both groups are asked for at the write ceiling — never destructive.
	require.Contains(t, reg.calls, "usage/write")
	require.Contains(t, reg.calls, "configuring/write")
	for _, c := range reg.calls {
		require.NotContains(t, c, "destructive", "a deletion is never inherited")
	}

	// Nothing is narrowed: the agent's own scope is what bounds an inherited
	// tool, and it carries no allow list of its own.
	for _, e := range out.Access.Tools {
		require.Empty(t, e.Allow, "%s must not invent a scope", e.Name)
	}
}

// Reading changes nothing and runs unannounced; anything that writes is put to
// the user first.
func TestExpandToolGrantsResolvesModesByRisk(t *testing.T) {
	out := expandToolGrants(&types.Agent{}, grantRegistry())

	modes := map[string]string{}
	for _, e := range out.Access.Tools {
		modes[e.Name] = e.Permission
	}

	require.Equal(t, "always", modes["compose_record_lookup"])
	require.Equal(t, "always", modes["compose_record_report"])
	require.Equal(t, "always", modes["compose_module_lookup"])
	require.Equal(t, "ask", modes["compose_record_create"])
	require.Equal(t, "ask", modes["compose_module_create"])
}

// A configured agent is narrowed exactly as it was written — inheritance is
// what happens when nothing was written, not a floor under everything.
func TestExpandToolGrantsDoesNotInheritOverAConfiguredAgent(t *testing.T) {
	out := expandToolGrants(grantingAgent(
		types.AgentAccessTool{Name: "compose_record_lookup"},
	), grantRegistry())

	require.Equal(t, []string{"compose_record_lookup"}, namesOf(out))
}

// "All data tools, but ask before deleting" is written as a group grant plus a
// named entry. Expanding the group first would claim the name and take its mode
// with it, leaving the override with no effect.
func TestExpandToolGrantsPrefersTheNamedEntry(t *testing.T) {
	out := expandToolGrants(grantingAgent(
		types.AgentAccessTool{Group: "usage", MaxRisk: "write"},
		types.AgentAccessTool{Name: "compose_record_create", Permission: "deny"},
	), grantRegistry())

	modes := map[string]string{}
	for _, e := range out.Access.Tools {
		modes[e.Name] = e.Permission
	}

	require.Equal(t, "deny", modes["compose_record_create"], "the named entry must win")
	require.Equal(t, "always", modes["compose_record_lookup"], "the rest still comes from the group")
}
