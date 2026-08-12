package mcpkit

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummarize(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{
			"cuts at the first sentence",
			"Create a new page in a namespace. A page is a screen in the navigation.",
			"Create a new page in a namespace.",
		},
		{
			"cuts at a line break that comes first",
			"Create a module.\n\nModules define a data structure with typed fields.",
			"Create a module.",
		},
		{
			// Splitting naively on ". " truncates this to "Query by resource, e.g."
			// which says nothing about what the tool does.
			"does not stop at an abbreviation",
			"Query by resource, e.g. a page or a role. Ranked by term count.",
			"Query by resource, e.g. a page or a role.",
		},
		{
			"a single sentence is returned whole",
			"Delete a role by name, handle, or ID",
			"Delete a role by name, handle, or ID",
		},
		{"empty stays empty", "", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, summarize(c.in))
		})
	}
}

// The registry hands out tools whose schema maps are shared, so slimming a copy
// must not reach back into the original. Getting this wrong would delete the
// documentation on first use, for human_tool_load and for the in-process
// runtime as well as for the listing — and it would look like it worked.
func TestSlimToolDoesNotMutateTheOriginal(t *testing.T) {
	full := mcp.NewTool("compose_page_create",
		mcp.WithDescription("Create a page. The grid is 48 columns wide."),
		mcp.WithString("namespace", mcp.Description("Namespace name, handle, slug, or ID.")),
		InGroup(GroupConfiguring), WithRisk(RiskWrite),
	)

	slim := slimTool(full)

	assert.Equal(t, "Create a page.", slim.Description)
	require.Contains(t, slim.InputSchema.Properties, "namespace")
	assert.NotContains(t, slim.InputSchema.Properties["namespace"], "description",
		"parameter prose is what the listing drops")
	assert.Contains(t, slim.InputSchema.Properties["namespace"], "type",
		"the shape a caller needs to build the call must survive")

	assert.Equal(t, "Create a page. The grid is 48 columns wide.", full.Description)
	assert.Equal(t, "Namespace name, handle, slug, or ID.",
		full.InputSchema.Properties["namespace"].(map[string]any)["description"],
		"the registry's own copy must be untouched")
}

func TestSlimToolFlagsToolsThatNeedTheirDocs(t *testing.T) {
	plain := slimTool(mcp.NewTool("system_role_create",
		mcp.WithDescription("Create a role. It grants nothing until rules are added."),
		InGroup(GroupConfiguring), WithRisk(RiskWrite)))
	assert.Equal(t, "Create a role.", plain.Description)

	flagged := slimTool(mcp.NewTool("compose_page_create",
		mcp.WithDescription("Create a page. The grid is 48 columns wide."),
		InGroup(GroupConfiguring), WithRisk(RiskWrite), NeedsFullDocs()))
	assert.Contains(t, flagged.Description, toolLoadName,
		"a tool whose rules live in its prose must say so in the summary")
}

func searchFixture(t *testing.T) *MCPServer {
	t.Helper()

	reg := NewRegistry()
	add := func(name, desc string, g Group, r Risk, opts ...mcp.ToolOption) {
		reg.RegisterTool(
			mcp.NewTool(name, append([]mcp.ToolOption{
				mcp.WithDescription(desc), InGroup(g), WithRisk(r),
			}, opts...)...),
			name, nil,
		)
	}

	add("system_role_create", "Create a role that grants permissions.", GroupConfiguring, RiskWrite)
	add("system_role_delete", "Delete a role. Members lose what it granted.", GroupConfiguring, RiskDestructive)
	add("compose_record_create", "Create a record in a module.", GroupUsage, RiskWrite)
	add("compose_page_update", "Update a page. Blocks are merged by block id.", GroupConfiguring, RiskWrite)
	add("system_user_suspend", "Suspend a user so they cannot sign in.", GroupConfiguring, RiskWrite)

	// Sorts after compose_chart_create alphabetically and is longer, so it
	// loses the tie-break rather than winning it by name.
	add("automation_taq_create", "Create a TAQ. Useful for scheduled reporting.", GroupConfiguring, RiskWrite)
	add("compose_chart_create", "Create a chart in a namespace.", GroupConfiguring, RiskWrite,
		WithKeywords("report", "dashboard"))

	return &MCPServer{reg: reg}
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

// TestSearchRanking pins the three defects measured against the live surface:
// a one-letter word matched every tool, ties fell through to alphabetical order
// so automation_* always won, and a word Human does not use could not reach the
// tool that answers it however the description was written.
func TestSearchRanking(t *testing.T) {
	m := searchFixture(t)

	t.Run("a keyword reaches a tool whose name does not contain the word", func(t *testing.T) {
		got := names(m.searchTools("dashboard", Scope{}))
		require.NotEmpty(t, got, "a word Human does not use must still find the tool that answers it")
		assert.Equal(t, "compose_chart_create", got[0])
	})

	t.Run("a keyword outranks an incidental description mention", func(t *testing.T) {
		// automation_taq_create says "reporting" in its description; the chart
		// tool declares "report" as a keyword. The keyword has to win, or no
		// amount of description editing can fix the ordering.
		got := names(m.searchTools("report", Scope{}))
		require.NotEmpty(t, got)
		assert.Equal(t, "compose_chart_create", got[0])
		assert.Contains(t, got, "automation_taq_create", "the weaker match is still a candidate")
	})

	t.Run("a filler word contributes nothing", func(t *testing.T) {
		// "a" is a substring of almost every name and description, so before
		// the minimum term length it pulled in the entire surface: the phrased
		// question scored worse than the bare word it contained.
		bare := names(m.searchTools("chart", Scope{}))
		phrased := names(m.searchTools("a chart", Scope{}))
		require.NotEmpty(t, bare)
		assert.Equal(t, bare, phrased,
			"the filler word must leave the result identical")
	})

	t.Run("a query of only short words still matches", func(t *testing.T) {
		// Dropping every term would turn a real question into silence.
		assert.NotEmpty(t, names(m.searchTools("id", Scope{})))
	})

	t.Run("an equal match does not default to alphabetical order", func(t *testing.T) {
		// Both name-match "create" on their op segment and score identically.
		// Alphabetically automation_taq_create wins; by name length the more
		// general compose_chart_create does.
		got := names(m.searchTools("create", Scope{}))
		require.NotEmpty(t, got)
		assert.NotEqual(t, "automation_taq_create", got[0],
			"automation must not win a tie merely by sorting first")
	})
}

// Search used to skip the five always-on tools, because a tool already in the
// listing could not be "disclosed" and returning it wasted a slot. Now search
// returns documentation rather than access, so every listed tool is a legitimate
// result — a caller asking about a lookup tool wants its parameter docs, and
// excluding it would answer "no such tool" about one they can plainly see.
func TestSearchReturnsToolsThatAreAlreadyListed(t *testing.T) {
	m := searchFixture(t)
	m.reg.RegisterTool(
		mcp.NewTool("compose_module_lookup", mcp.WithDescription("Look up modules."),
			InGroup(GroupConfiguring), WithRisk(RiskRead)),
		"compose_module_lookup", nil,
	)

	assert.Contains(t, names(m.searchTools("module lookup", Scope{})), "compose_module_lookup")
}
