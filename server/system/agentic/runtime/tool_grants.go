package runtime

import (
	"slices"

	"github.com/crusttech/human/server/system/agentic/policy"
	"github.com/crusttech/human/server/system/types"
)

// defaultGrantRisk is what a group grant permits when it does not say.
//
// Reading is the level that cannot damage anything, and a grant written without
// thinking about risk should not be the one that lets an agent delete records.
const defaultGrantRisk = "read"

// inheritedGroups is what an agent that configures no tools at all is given.
//
// It runs as the user who invoked it, so RBAC already bounds everything it can
// touch; listing tools by hand on top of that was the single thing that made an
// agent tedious to set up and stale afterwards. Destructive tools are the
// exception and still have to be named: a deletion is the one action no
// approval prompt reliably takes back.
var inheritedGroups = []types.AgentAccessTool{
	{Group: "usage", MaxRisk: "write", Description: "Inherited: data and execution, within what you can reach"},
	{Group: "configuring", MaxRisk: "write", Description: "Inherited: schema and definitions, within what you can reach"},
}

// expandToolGrants replaces every group grant with the tools it stands for.
//
// A grant may name one tool, or a group and a risk ceiling — "everything in
// 'usage' up to 'read'" — which is how an agent is given a working set without
// listing it, and how it keeps working when a tool is added later. Expanding
// here, once, means the policy check and the tool list both see an ordinary
// allow-list and neither has to know groups exist.
//
// The agent is copied rather than edited: it comes from a shared registry, and
// widening the caller's copy would leak one request's grant into the next.
func expandToolGrants(agent *types.Agent, reg MCPClient) *types.Agent {
	if agent == nil || reg == nil {
		return agent
	}

	configured := agent.Access.Tools
	if len(configured) == 0 {
		configured = inheritedGroups
	} else if !slices.ContainsFunc(configured, func(t types.AgentAccessTool) bool { return t.Group != "" }) {
		return withResolvedPermissions(agent, reg)
	}

	out := agent.Clone()
	out.Access.Tools = configured
	expanded := make([]types.AgentAccessTool, 0, len(out.Access.Tools))
	seen := make(map[string]bool, len(out.Access.Tools))

	// Named entries first, whatever order they were written in: an entry that
	// names a tool is the more specific statement about it, and a group grant
	// expanding over it would silently take its permission mode with it —
	// which is how "all data tools, but ask before delete" is written.
	for _, t := range out.Access.Tools {
		if t.Group != "" || seen[t.Name] {
			continue
		}
		seen[t.Name] = true
		expanded = append(expanded, t)
	}

	for _, t := range out.Access.Tools {
		if t.Group == "" {
			continue
		}

		risk := t.MaxRisk
		if risk == "" {
			risk = defaultGrantRisk
		}

		for _, name := range reg.ToolNamesIn(t.Group, risk) {
			if seen[name] {
				continue
			}
			seen[name] = true

			// Each expanded entry carries the group's own scope and context, so
			// the grant means the same whichever tool it turned into.
			member := t
			member.Name = name
			member.Group = ""
			member.MaxRisk = ""
			expanded = append(expanded, member)
		}
	}

	out.Access.Tools = expanded
	return withResolvedPermissions(out, reg)
}

// withResolvedPermissions fills in the mode of every grant that names none.
//
// Reading changes nothing and runs unannounced; anything that writes is put to
// the user first. A grant that states its own mode is left exactly as written,
// including one set to deny.
func withResolvedPermissions(agent *types.Agent, reg MCPClient) *types.Agent {
	if !slices.ContainsFunc(agent.Access.Tools, func(t types.AgentAccessTool) bool { return t.Permission == "" }) {
		return agent
	}

	readOnly := readOnlyTools(reg)

	out := agent.Clone()
	for i := range out.Access.Tools {
		if out.Access.Tools[i].Permission != "" {
			continue
		}
		if readOnly[out.Access.Tools[i].Name] {
			out.Access.Tools[i].Permission = policy.PermissionAlways
		} else {
			out.Access.Tools[i].Permission = policy.PermissionAsk
		}
	}
	return out
}

// readOnlyTools is every tool that only reads, across both groups.
//
// The registry answers "in this group, at or below this risk", so the read-only
// set is the answer at the read ceiling — no separate risk lookup needed.
func readOnlyTools(reg MCPClient) map[string]bool {
	out := map[string]bool{}
	for _, group := range []string{"usage", "configuring"} {
		for _, name := range reg.ToolNamesIn(group, "read") {
			out[name] = true
		}
	}
	return out
}
