package agentic

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	autoService "github.com/crusttech/human/server/automation/service"
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
func (h *agentHandler) validateAgentAccess(ctx context.Context, a *sysTypes.Agent) error {
	if err := h.validateAgentTools(a); err != nil {
		return err
	}
	return validateAgentAutomations(ctx, a)
}

// validateAgentAutomations refuses a TAQ or workflow the agent cannot be given.
//
// Unlike a tool name, a dangling automation ID costs the agent only that one
// capability: loadTAQInfos drops what it cannot find and getAvailableTools
// skips it, both without a word. So the agent answers as though it were never
// granted the thing its configuration plainly lists.
func validateAgentAutomations(ctx context.Context, a *sysTypes.Agent) error {
	for i, t := range a.Access.TAQs {
		if t.ID == 0 {
			return fmt.Errorf(`access.taqs[%d]: "id" is required — the TAQ this entry lets the agent run`, i)
		}
		if svc := autoService.DefaultNgAutomation; svc != nil {
			if _, err := svc.FindByID(ctx, t.ID); err != nil {
				return fmt.Errorf(
					"access.taqs[%d]: no TAQ with ID %d. A TAQ the runtime cannot find is skipped in silence, so the agent would simply not have it: %w",
					i, t.ID, err,
				)
			}
		}
	}

	for i, w := range a.Access.Workflows {
		if w.ID == 0 {
			return fmt.Errorf(`access.workflows[%d]: "id" is required — the workflow this entry lets the agent run`, i)
		}
		if svc := autoService.DefaultWorkflow; svc != nil {
			if _, err := svc.FindByID(ctx, w.ID); err != nil {
				return fmt.Errorf(
					"access.workflows[%d]: no workflow with ID %d. The agent is still handed automation_workflow_exec, so this fails when it calls it: %w",
					i, w.ID, err,
				)
			}
		}
	}

	return nil
}

func (h *agentHandler) validateAgentTools(a *sysTypes.Agent) error {
	reg, ok := h.reg.(knownTools)
	if !ok {
		// A registrar that cannot be asked; nothing to check against.
		return nil
	}

	for i, t := range a.Access.Tools {
		// A grant names one tool, or a group and a risk ceiling standing for
		// every tool in it.
		if t.Group != "" {
			if err := validateGrantGroup(i, t, reg); err != nil {
				return err
			}
			continue
		}

		if strings.TrimSpace(t.Name) == "" {
			return fmt.Errorf(`access.tools[%d]: needs a "name" (one tool) or a "group" (every tool in it, capped by "maxRisk")`, i)
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

// validateGrantGroup checks a group grant names a real group and risk, and that
// the pair actually stands for something.
//
// An empty expansion is refused rather than stored: it reads as a grant and
// behaves as none, which is the failure that is hardest to see from the config.
func validateGrantGroup(i int, t sysTypes.AgentAccessTool, reg knownTools) error {
	if t.Name != "" {
		return fmt.Errorf(
			`access.tools[%d]: set "name" or "group", not both — %q would be granted twice over and it is not clear which scope wins`,
			i, t.Name,
		)
	}

	if !slices.Contains(grantGroups, t.Group) {
		return fmt.Errorf(
			"access.tools[%d]: group %q is not one of %s",
			i, t.Group, strings.Join(quoteAll(grantGroups), ", "),
		)
	}

	risk := t.MaxRisk
	if risk == "" {
		risk = defaultGrantRisk
	}

	if !slices.Contains(grantRisks, risk) {
		return fmt.Errorf(
			"access.tools[%d]: maxRisk %q is not one of %s",
			i, t.MaxRisk, strings.Join(quoteAll(grantRisks), ", "),
		)
	}

	lister, ok := reg.(groupLister)
	if !ok {
		return nil
	}

	if len(lister.ToolNamesIn(t.Group, risk)) == 0 {
		return fmt.Errorf(
			"access.tools[%d]: group %q at maxRisk %q covers no tools, so the entry grants nothing",
			i, t.Group, risk,
		)
	}

	return nil
}

// groupLister is the part of the registry that can expand a group grant.
type groupLister interface {
	ToolNamesIn(group, maxRisk string) []string
}

// What a grant may name. Mirrors mcpkit's Group and Risk without importing the
// vocabulary, and "development" is deliberately absent: it is the developer
// MCP's surface, not something to hand an agent.
var (
	grantGroups = []string{"configuring", "usage"}
	grantRisks  = []string{"read", "write", "destructive"}
)

// defaultGrantRisk mirrors the runtime's: a grant that does not say what risk it
// permits permits only reading.
const defaultGrantRisk = "read"

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
