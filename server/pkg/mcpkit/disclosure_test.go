package mcpkit

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisclosureIsPerSession(t *testing.T) {
	d := newDisclosure()

	d.load("session-a", "system_role_create", "system_role_delete")

	assert.True(t, d.isLoaded("session-a", "system_role_create"))
	assert.False(t, d.isLoaded("session-a", "system_user_create"), "only what was loaded")

	// One session pulling a tool in must not reveal it to another. This is the
	// property that makes disclosure per-session state rather than a global
	// cache.
	assert.False(t, d.isLoaded("session-b", "system_role_create"))

	d.forget("session-a")
	assert.False(t, d.isLoaded("session-a", "system_role_create"), "teardown clears state")
}

// A request with no session cannot accumulate anything. It still works — the
// search result carries the schemas inline — but nothing is remembered.
func TestDisclosureWithoutSession(t *testing.T) {
	d := newDisclosure()
	d.load("", "system_role_create")
	assert.False(t, d.isLoaded("", "system_role_create"))
}

func searchFixture(t *testing.T) *MCPServer {
	t.Helper()

	reg := NewRegistry()
	add := func(name, desc string, g Group, r Risk) {
		reg.RegisterTool(
			mcp.NewTool(name, mcp.WithDescription(desc), InGroup(g), WithRisk(r)),
			name, nil,
		)
	}

	add("system_role_create", "Create a role that grants permissions.", GroupConfiguring, RiskWrite)
	add("system_role_delete", "Delete a role. Members lose what it granted.", GroupConfiguring, RiskDestructive)
	add("compose_record_create", "Create a record in a module.", GroupUsage, RiskWrite)
	add("compose_page_update", "Update a page. Blocks are merged by block id.", GroupConfiguring, RiskWrite)
	add("system_user_suspend", "Suspend a user so they cannot sign in.", GroupConfiguring, RiskWrite)

	return &MCPServer{reg: reg, disclosed: newDisclosure()}
}

func names(tools []mcp.Tool) []string {
	out := make([]string, 0, len(tools))
	for _, t := range tools {
		out = append(out, t.Name)
	}
	return out
}

func TestSearchMatching(t *testing.T) {
	m := searchFixture(t)

	t.Run("name hits outrank description hits", func(t *testing.T) {
		got := names(m.searchTools("role", Scope{}))
		require.NotEmpty(t, got)
		// Both role tools name-match; system_user_suspend does not match at all.
		assert.Subset(t, []string{"system_role_create", "system_role_delete"}, got[:2])
	})

	t.Run("matching more terms ranks higher", func(t *testing.T) {
		got := names(m.searchTools("delete role", Scope{}))
		require.NotEmpty(t, got)
		assert.Equal(t, "system_role_delete", got[0],
			"the tool matching both terms comes first")
		assert.Contains(t, got, "system_role_create",
			"a tool matching one term is still a candidate")
	})

	t.Run("a term that matches nothing does not eliminate the tool", func(t *testing.T) {
		// Requiring every term made a natural multi-word question the worst
		// possible input, and "no tool matches" reads as "no such capability".
		assert.Contains(t, names(m.searchTools("role nonexistentword", Scope{})), "system_role_create")
	})

	t.Run("a query matching nothing at all is still empty", func(t *testing.T) {
		assert.Empty(t, names(m.searchTools("nonexistentword", Scope{})))
	})

	t.Run("an empty query matches nothing rather than everything", func(t *testing.T) {
		assert.Empty(t, names(m.searchTools("", Scope{})))
	})

	// Search must not become a way around the scope the session asked for.
	t.Run("scope still applies", func(t *testing.T) {
		got := names(m.searchTools("create", Scope{Group: GroupUsage}))
		assert.Equal(t, []string{"compose_record_create"}, got)

		got = names(m.searchTools("role", Scope{MaxRisk: RiskWrite}))
		assert.NotContains(t, got, "system_role_delete", "destructive is above the ceiling")
		assert.Contains(t, got, "system_role_create")
	})
}

func TestAlwaysOnToolsAreNotSearchResults(t *testing.T) {
	m := searchFixture(t)
	m.reg.RegisterTool(
		mcp.NewTool("compose_module_lookup", mcp.WithDescription("Look up modules."),
			InGroup(GroupConfiguring), WithRisk(RiskRead)),
		"compose_module_lookup", nil,
	)

	// It is already listed, so returning it would waste a search slot.
	assert.NotContains(t, names(m.searchTools("module lookup", Scope{})), "compose_module_lookup")
	assert.True(t, alwaysOn["compose_module_lookup"])
}
