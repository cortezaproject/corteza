package agentic

import (
	"fmt"
	"sort"
	"strings"

	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// knownTools is the slice of the registry needed to tell a real tool name from
// a typo. Narrow on purpose: toolRegistrar is shared by every agentic handler
// and only registers.
type knownTools interface {
	HasTool(name string) bool
	Tools() []mcp.Tool
}

// validateAgentTools refuses an agent naming a tool the registry does not know.
//
// Nothing downstream tolerates one. The runtime resolves an agent's whole
// allow-list through Registry.Select, which fails on the first unknown name, so
// a single typo stops the agent running at all rather than costing it one
// capability. Meanwhile the misspelt entry's description is still appended to
// the system prompt, telling the model about a tool nothing can serve. Both
// only surface when someone runs the agent, which is why the name is checked
// where it is written.
func (h *agentHandler) validateAgentTools(a *sysTypes.Agent) error {
	reg, ok := h.reg.(knownTools)
	if !ok {
		// A registrar that cannot be asked; nothing to check against.
		return nil
	}

	for i, t := range a.Access.Tools {
		if strings.TrimSpace(t.Name) == "" {
			return fmt.Errorf(`access.tools[%d]: "name" is required — the MCP tool this entry grants`, i)
		}
		if reg.HasTool(t.Name) {
			continue
		}

		msg := fmt.Sprintf(
			"access.tools[%d]: no tool named %q. An agent's allow-list is resolved as a whole, "+
				"so one unknown name stops the agent running at all",
			i, t.Name,
		)
		if near := nearestToolNames(t.Name, reg.Tools()); len(near) > 0 {
			return fmt.Errorf("%s. Did you mean %s?", msg, strings.Join(quoteAll(near), " or "))
		}
		return fmt.Errorf("%s", msg)
	}

	return nil
}

// nearestToolNames returns the registered names closest to a misspelt one.
// Tool names are long and structured (compose_record_lookup), so a typo is
// nearly always within an edit or two, and a wrong family still shares a
// prefix.
func nearestToolNames(name string, tools []mcp.Tool) []string {
	want := strings.ToLower(strings.TrimSpace(name))

	type scored struct {
		name string
		dist int
	}
	var out []scored

	for _, t := range tools {
		k := strings.ToLower(t.Name)
		d := editDistance(want, k)
		// Either a near-miss spelling, or the same thing said about another
		// resource — both are what someone reaching for this name meant.
		if d <= 2 || strings.HasPrefix(k, want) || strings.HasPrefix(want, k) {
			out = append(out, scored{t.Name, d})
		}
	}

	if len(out) == 0 {
		return nil
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].dist != out[j].dist {
			return out[i].dist < out[j].dist
		}
		return out[i].name < out[j].name
	})

	names := make([]string, 0, 3)
	for _, s := range out {
		if len(names) == 3 {
			break
		}
		names = append(names, s.name)
	}
	return names
}

// editDistance is Levenshtein, over two rows rather than a full matrix.
func editDistance(a, b string) int {
	if a == b {
		return 0
	}
	ra, rb := []rune(a), []rune(b)
	if len(ra) == 0 {
		return len(rb)
	}
	if len(rb) == 0 {
		return len(ra)
	}

	prev := make([]int, len(rb)+1)
	curr := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(ra); i++ {
		curr[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			curr[j] = min(prev[j]+1, min(curr[j-1]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}

	return prev[len(rb)]
}

func quoteAll(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = fmt.Sprintf("%q", s)
	}
	return out
}
