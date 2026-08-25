package runtime

import (
	"slices"

	"github.com/crusttech/human/server/system/types"
)

// defaultGrantRisk is what a group grant permits when it does not say.
//
// Reading is the level that cannot damage anything, and a grant written without
// thinking about risk should not be the one that lets an agent delete records.
const defaultGrantRisk = "read"

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

	if !slices.ContainsFunc(agent.Access.Tools, func(t types.AgentAccessTool) bool { return t.Group != "" }) {
		return agent
	}

	out := agent.Clone()
	expanded := make([]types.AgentAccessTool, 0, len(out.Access.Tools))
	seen := make(map[string]bool, len(out.Access.Tools))

	for _, t := range out.Access.Tools {
		if t.Group == "" {
			if !seen[t.Name] {
				seen[t.Name] = true
				expanded = append(expanded, t)
			}
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
	return out
}
