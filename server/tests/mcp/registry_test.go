// Package mcp_test holds the structural assertions that keep the MCP tool
// surface honest as it grows from ~36 tools toward ~200.
//
// It cannot live in server/system/agentic/mcp: every handler package imports
// that package, so a test there that imported the handlers back would close a
// cycle. It lives here instead, as a leaf that imports everything and is
// imported by nothing.
//
// The rules asserted here are specified in
// server/system/agentic/mcp/CONVENTIONS.md; the resource-name table it checks
// against is RESOURCES.md alongside it. When a rule changes, change the spec
// first.
package mcp_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	autoAgentic "github.com/crusttech/human/server/automation/agentic"
	cmpAgentic "github.com/crusttech/human/server/compose/agentic"
	a "github.com/crusttech/human/server/pkg/auth"
	sysAgentic "github.com/crusttech/human/server/system/agentic"
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/crusttech/human/server/system/agentic/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubSigner satisfies the discovery handler's unexported signer interface.
// Discovery registers only when it has a base URL, so the test supplies one
// that is never dialled — registration is all this file cares about.
type stubSigner struct{}

func (stubSigner) Sign(...a.IssueOptFn) ([]byte, error) { return []byte("stub"), nil }

// buildRegistry wires every handler into a fresh registry, mirroring
// app/boot_levels.go. Handlers only touch the registrar at construction, so no
// services or database are needed.
//
// A handler added to boot_levels.go and not added here is invisible to every
// assertion below — TestRegistryMatchesBootWiring guards that.
func buildRegistry(t *testing.T) *hmcp.Registry {
	t.Helper()

	reg := hmcp.NewRegistry()

	cmpAgentic.RecordHandler(reg)
	cmpAgentic.NamespaceHandler(reg, nil)
	cmpAgentic.ModuleHandler(reg, nil)
	cmpAgentic.PageHandler(reg)
	cmpAgentic.ChartHandler(reg)
	autoAgentic.TAQHandler(reg)
	autoAgentic.WorkflowHandler(reg)
	autoAgentic.TriggerHandler(reg)
	autoAgentic.EventTypeHandler(reg)
	sysAgentic.ReminderHandler(reg)
	sysAgentic.UserHandler(reg)
	sysAgentic.UserGroupHandler(reg)
	sysAgentic.RoleHandler(reg)
	sysAgentic.AuthClientHandler(reg)
	sysAgentic.ApplicationHandler(reg)
	sysAgentic.DiscoveryHandler(reg, "http://discovery.invalid", stubSigner{})

	return reg
}

// legacyNames are the three shipped tool names that predate the grammar and are
// grandfathered rather than renamed. discovery_search is additionally hardcoded
// in eight places in runtime/executor.go. See CONVENTIONS.md §6.
//
// This list must not grow. A new tool that needs an entry here is a tool named
// wrongly.
var legacyNames = map[string]bool{
	"discovery_search":               true,
	"compose_page_block_schema":      true,
	"automation_taq_execution_trace": true,
}

// validOps is the op vocabulary from CONVENTIONS.md §6. Domain ops are allowed
// beyond this set, so a name failing here is reported rather than failed —
// except that its resource segment must still appear in RESOURCES.md.
var validOps = map[string]bool{
	"lookup": true, "create": true, "update": true,
	"delete": true, "undelete": true, "exec": true,
}

// splitName breaks a tool name into its app, resource and op segments.
//
// The resource segment is snake_case and so contains underscores, which is why
// nothing in the codebase parses tool names structurally — RESOURCES.md is the
// authority. Splitting on the first and last underscore is good enough for the
// assertions here.
func splitName(name string) (app, resource, op string) {
	first := strings.Index(name, "_")
	last := strings.LastIndex(name, "_")
	if first < 0 || first == last {
		return name, "", ""
	}
	return name[:first], name[first+1 : last], name[last+1:]
}

// camelCase converts a snake_case resource segment to the camelCase form a
// param uses: auth_client -> authClient, user_group -> userGroup.
func camelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func TestToolNamesFollowTheGrammar(t *testing.T) {
	for _, tool := range buildRegistry(t).Tools() {
		if legacyNames[tool.Name] {
			continue
		}

		app, resource, op := splitName(tool.Name)

		assert.Containsf(t, []string{"compose", "automation", "system"}, app,
			"tool %q has an unknown app segment", tool.Name)
		assert.NotEmptyf(t, resource, "tool %q has no resource segment", tool.Name)
		assert.NotEmptyf(t, op, "tool %q has no op segment", tool.Name)

		// Domain ops beyond the standard vocabulary are allowed (§6), so an
		// unrecognised op is reported rather than failed — it is a prompt to
		// check the name was deliberate, not a defect on its own.
		if op != "" && !validOps[op] {
			t.Logf("note: %q uses domain op %q, outside the standard vocabulary", tool.Name, op)
		}

		assert.Equalf(t, strings.ToLower(tool.Name), tool.Name,
			"tool %q is not all lower snake_case; see RESOURCES.md", tool.Name)
	}
}

func TestEveryToolIsTagged(t *testing.T) {
	validGroups := map[hmcp.Group]bool{
		hmcp.GroupDevelopment: true,
		hmcp.GroupConfiguring: true,
		hmcp.GroupUsage:       true,
	}
	validRisks := map[hmcp.Risk]bool{
		hmcp.RiskRead: true, hmcp.RiskWrite: true, hmcp.RiskDestructive: true,
	}

	for _, tool := range buildRegistry(t).Tools() {
		groups := hmcp.GroupsOf(tool)
		require.NotEmpty(t, groups, "tool %q has no group; see CONVENTIONS.md §3", tool.Name)
		for _, g := range groups {
			assert.True(t, validGroups[g], "tool %q has unknown group %q", tool.Name, g)
		}

		risk := hmcp.RiskOf(tool)
		require.NotEmpty(t, risk, "tool %q has no risk; see CONVENTIONS.md §4", tool.Name)
		assert.True(t, validRisks[risk], "tool %q has unknown risk %q", tool.Name, risk)
	}
}

// TestRiskDrivesAnnotations guards the defect that motivated making WithRisk the
// sole writer of the hints: mcp.NewTool defaults DestructiveHint and
// OpenWorldHint to true, so a tool that does not go through WithRisk advertises
// itself to every client as destructive and open-world.
func TestRiskDrivesAnnotations(t *testing.T) {
	for _, tool := range buildRegistry(t).Tools() {
		risk := hmcp.RiskOf(tool)
		require.NotNil(t, tool.Annotations.ReadOnlyHint, tool.Name)
		require.NotNil(t, tool.Annotations.DestructiveHint, tool.Name)
		require.NotNil(t, tool.Annotations.OpenWorldHint, tool.Name)

		assert.False(t, *tool.Annotations.OpenWorldHint,
			"tool %q claims openWorld; every Human tool acts on this instance", tool.Name)

		switch risk {
		case hmcp.RiskRead:
			assert.True(t, *tool.Annotations.ReadOnlyHint, "read tool %q is not readOnly", tool.Name)
			assert.False(t, *tool.Annotations.DestructiveHint, "read tool %q claims destructive", tool.Name)
		case hmcp.RiskDestructive:
			assert.True(t, *tool.Annotations.DestructiveHint, "destructive tool %q does not say so", tool.Name)
		case hmcp.RiskWrite:
			assert.False(t, *tool.Annotations.ReadOnlyHint, "write tool %q claims readOnly", tool.Name)
		}
	}
}

// minDescription is a floor, not a standard. It catches "Execute a TAQ" — the
// kind of description that is useless to a caller without the source. Quality
// is a review criterion; see CONVENTIONS.md §8.5.
const minDescription = 60

func TestDescriptionsAreSubstantial(t *testing.T) {
	for _, tool := range buildRegistry(t).Tools() {
		assert.GreaterOrEqualf(t, len(tool.Description), minDescription,
			"tool %q has a %d-char description; a caller without the repo has only this",
			tool.Name, len(tool.Description))
		assert.NotContains(t, tool.Description, "Corteza",
			"tool %q says Corteza; the product is Human", tool.Name)
	}
}

// TestIDParamsAreDeclaredAsStrings enforces CONVENTIONS.md §8.3. Human IDs are
// uint64 and exceed JavaScript's safe integer range, so a numeric param would
// silently truncate.
func TestIDParamsAreDeclaredAsStrings(t *testing.T) {
	for _, tool := range buildRegistry(t).Tools() {
		for name, raw := range tool.InputSchema.Properties {
			if !strings.HasSuffix(name, "ID") && !strings.HasSuffix(name, "IDs") {
				continue
			}
			prop, ok := raw.(map[string]any)
			require.Truef(t, ok, "tool %q param %q has an unreadable schema", tool.Name, name)
			assert.Equalf(t, "string", prop["type"],
				"tool %q param %q must be a string to avoid precision loss", tool.Name, name)
		}
	}
}

// TestLookupContract enforces CONVENTIONS.md §8.1.
//
// discovery_search is exempt from paging only: its backing service offsets with
// `from` and returns no cursor, so a pageCursor param would advertise paging
// that cannot work.
func TestLookupContract(t *testing.T) {
	// Tools whose backing data cannot be paged, and why. §8.1's rule is that a
	// lookup must be bounded — not that every lookup must page. Advertising a
	// cursor that cannot be honoured is the failure the rule exists to prevent,
	// so a genuine exemption is safer than a broken cursor.
	//
	// An entry here is a claim about the data source, not a convenience. Both
	// were verified against the source, not assumed.
	noCursor := map[string]bool{
		// Elasticsearch-style API: offsets with `from`, returns no cursor.
		"discovery_search": true,
		// A compile-time constant slice, not a store. Measured at 118 entries
		// and 37 KB — 14% of the JSONResult ceiling — so it returns whole.
		"automation_event_type_lookup": true,
	}

	// Of those, the ones that cannot meaningfully bound either. discovery_search
	// still caps result size, so it keeps its `limit`.
	noLimit := map[string]bool{"automation_event_type_lookup": true}

	for _, tool := range buildRegistry(t).Tools() {
		if !strings.HasSuffix(tool.Name, "_lookup") && tool.Name != "discovery_search" {
			continue
		}

		// Only the tool's own reference param must stay optional. A parent-scope
		// param is a different thing and may well be required: listing modules
		// without saying which namespace is meaningless, so
		// compose_module_lookup requires `namespace` and leaves `module` free.
		//
		// The resource segment is snake_case while a param is camelCase, so
		// `auth_client` has to be matched against `authClient`. Comparing the two
		// verbatim silently matched nothing and quietly exempted every
		// multi-word resource from this assertion.
		_, resource, _ := splitName(tool.Name)
		refNames := map[string]bool{
			resource: true, resource + "ID": true,
			camelCase(resource): true, camelCase(resource) + "ID": true,
		}

		for _, req := range tool.InputSchema.Required {
			assert.Falsef(t, refNames[req],
				"lookup %q marks its own reference param %q required, which makes list mode unreachable",
				tool.Name, req)
		}

		if !noLimit[tool.Name] {
			assert.Containsf(t, tool.InputSchema.Properties, "limit",
				"lookup %q declares no limit; an unbounded list can drain a module into the caller's context", tool.Name)
		}

		if !noCursor[tool.Name] {
			assert.Containsf(t, tool.InputSchema.Properties, "pageCursor",
				"lookup %q declares no pageCursor", tool.Name)
		}
	}
}

// TestEveryToolIsClassifiedInPolicy keeps the in-process authorization gate in
// step with the registry. policy.go cannot import the registry (it would close
// the cycle mcp → runtime → policy), so this is the seam that catches a new
// tool arriving with no resource scoping. See CONVENTIONS.md §2.4.
func TestEveryToolIsClassifiedInPolicy(t *testing.T) {
	for _, tool := range buildRegistry(t).Tools() {
		assert.Truef(t, policy.IsClassified(tool.Name),
			"tool %q has no buildResource case and is not in resourceScopeExempt "+
				"(server/system/agentic/policy/policy.go)", tool.Name)
	}
}

func TestAliasMapsAgree(t *testing.T) {
	assert.Equal(t, hmcp.ToolAliases, policy.ToolAliases(),
		"the alias maps in mcp and policy have drifted; both must list every rename")
}

// TestDeclarationsLiveInToolsFiles enforces the file split in CONVENTIONS.md §5:
// reading every *_tools.go must give a complete picture of the tool surface,
// which only holds if no declaration hides in a handler file.
func TestDeclarationsLiveInToolsFiles(t *testing.T) {
	for _, dir := range agenticDirs(t) {
		entries, err := os.ReadDir(dir)
		require.NoError(t, err)

		for _, e := range entries {
			name := e.Name()
			if !strings.HasSuffix(name, "_handler.go") {
				continue
			}
			body, err := os.ReadFile(filepath.Join(dir, name))
			require.NoError(t, err)

			assert.NotContainsf(t, string(body), "mcp.NewTool(",
				"%s declares a tool; declarations belong in the matching _tools.go", filepath.Join(dir, name))
		}
	}
}

// TestRegistryMatchesBootWiring fails when a handler is wired into the app but
// not into buildRegistry, which would silently exempt its tools from every
// assertion in this file.
func TestRegistryMatchesBootWiring(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(repoRoot(t), "server", "app", "boot_levels.go"))
	require.NoError(t, err)

	src := string(body)
	for _, ctor := range []string{
		"RecordHandler(", "NamespaceHandler(", "ModuleHandler(", "PageHandler(",
		"ChartHandler(", "TAQHandler(", "WorkflowHandler(", "TriggerHandler(", "EventTypeHandler(", "ReminderHandler(",
		"DiscoveryHandler(", "UserHandler(", "UserGroupHandler(", "RoleHandler(",
		"AuthClientHandler(", "ApplicationHandler(",
	} {
		assert.Containsf(t, src, ctor,
			"buildRegistry wires %s but boot_levels.go does not; one of them is wrong", ctor)
	}

	// The reverse direction: every *Handler( call in the wiring block must be
	// one this test knows about.
	known := map[string]bool{
		"RecordHandler": true, "NamespaceHandler": true, "ModuleHandler": true,
		"PageHandler": true, "ChartHandler": true, "TAQHandler": true,
		"WorkflowHandler": true, "TriggerHandler": true, "EventTypeHandler": true,
		"ReminderHandler": true, "DiscoveryHandler": true,
		"UserHandler": true, "UserGroupHandler": true, "RoleHandler": true,
		"AuthClientHandler": true, "ApplicationHandler": true,
	}
	for _, line := range strings.Split(src, "\n") {
		for _, prefix := range []string{"cmpAgentic.", "autoAgentic.", "sysAgentic."} {
			idx := strings.Index(line, prefix)
			if idx < 0 {
				continue
			}
			rest := line[idx+len(prefix):]
			end := strings.Index(rest, "(")
			if end < 0 {
				continue
			}
			ctor := rest[:end]

			// Only tool handlers matter here. These packages also export
			// non-registry helpers wired elsewhere in boot — NsModResolver
			// serves the agentic runtime, not the MCP registry.
			if !strings.HasSuffix(ctor, "Handler") {
				continue
			}

			assert.Truef(t, known[ctor],
				"boot_levels.go wires %s%s, which buildRegistry does not know about", prefix, ctor)
		}
	}
}

func agenticDirs(t *testing.T) []string {
	t.Helper()
	root := repoRoot(t)
	return []string{
		filepath.Join(root, "server", "compose", "agentic"),
		filepath.Join(root, "server", "automation", "agentic"),
		filepath.Join(root, "server", "system", "agentic"),
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd() // server/tests/mcp
	require.NoError(t, err)
	return filepath.Join(wd, "..", "..", "..")
}

// toolRow is one line of the coverage matrix.
type toolRow struct {
	Name   string
	Groups string
	Risk   string
	Hidden bool
}

func rows(reg *hmcp.Registry) []toolRow {
	tools := reg.Tools()
	out := make([]toolRow, 0, len(tools))
	for _, tool := range tools {
		groups := make([]string, 0, 2)
		for _, g := range hmcp.GroupsOf(tool) {
			groups = append(groups, string(g))
		}
		sort.Strings(groups)
		out = append(out, toolRow{
			Name:   tool.Name,
			Groups: strings.Join(groups, ", "),
			Risk:   string(hmcp.RiskOf(tool)),
			Hidden: reg.IsHidden(tool.Name),
		})
	}
	return out
}

func renderMatrix(reg *hmcp.Registry) string {
	var b strings.Builder

	b.WriteString("# MCP tool coverage\n\n")
	b.WriteString("Generated. Do not edit by hand — run:\n\n")
	b.WriteString("```sh\ncd server && go test ./tests/mcp/ -run TestToolsMatrix -update\n```\n\n")
	b.WriteString("Resource-to-tool naming is fixed by `RESOURCES.md`; the rules these\n")
	b.WriteString("tools are held to are in `CONVENTIONS.md`.\n\n")

	rr := rows(reg)
	fmt.Fprintf(&b, "## Registered tools (%d)\n\n", len(rr))
	b.WriteString("| Tool | Group | Risk | Surface |\n|---|---|---|---|\n")
	for _, r := range rr {
		surface := "both"
		if r.Hidden {
			surface = "in-process"
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s |\n", r.Name, r.Groups, r.Risk, surface)
	}

	byRisk := map[string]int{}
	byGroup := map[string]int{}
	for _, r := range rr {
		byRisk[r.Risk]++
		byGroup[r.Groups]++
	}

	b.WriteString("\n## Totals\n\n")
	b.WriteString("| Dimension | Value | Tools |\n|---|---|---|\n")
	for _, k := range sortedKeys(byGroup) {
		fmt.Fprintf(&b, "| group | %s | %d |\n", k, byGroup[k])
	}
	for _, k := range sortedKeys(byRisk) {
		fmt.Fprintf(&b, "| risk | %s | %d |\n", k, byRisk[k])
	}

	return b.String()
}

func sortedKeys(m map[string]int) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
