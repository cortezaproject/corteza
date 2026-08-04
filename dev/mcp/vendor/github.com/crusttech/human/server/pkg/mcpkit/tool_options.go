package mcpkit

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// Group and Risk are tool-definition concerns: they describe the tool itself,
// not how a particular registration exposes it, so they travel on the mcp.Tool
// via ToolOptions rather than through the registrar. Two consequences that were
// the point of choosing this shape:
//
//   - Registrar interfaces and every call site keep their arity, so adding tags
//     to ~200 tools is a per-declaration edit and never an interface migration.
//   - server.ToolFilterFunc receives []mcp.Tool, so per-request filtering by
//     group or risk can read the tag straight off the tool.
type (
	Group string
	Risk  string
)

const (
	// GroupDevelopment covers repo-level tooling: it operates on Human's source,
	// not on a running instance. Currently empty; defined so L1 tooling has a
	// home without a later retrofit.
	GroupDevelopment Group = "development"

	// GroupConfiguring covers the schema and definition level — what the system
	// *is*. Modules, pages, charts, TAQ and workflow definitions, roles.
	GroupConfiguring Group = "configuring"

	// GroupUsage covers the data and execution level — what the system *holds*,
	// or running it. Records, executions, chatbots, reports.
	GroupUsage Group = "usage"
)

const (
	// RiskRead makes no state change. Note this is about state, not disclosure:
	// a read-shaped tool can still expose sensitive data, which is governed by
	// RBAC and by the authoring rule on services that lack it.
	RiskRead Risk = "read"

	// RiskWrite creates or modifies. Execution is always write, even when the
	// executed thing happens to be read-only — the registry cannot know what a
	// TAQ or workflow does.
	RiskWrite Risk = "write"

	// RiskDestructive removes. Deletes in Human are mostly soft and reversible
	// via undelete, but a soft-deleted resource is invisible to every consumer,
	// so the blast radius is real.
	RiskDestructive Risk = "destructive"
)

// Meta keys under which the tags travel. They are namespaced because _meta is
// shared with the protocol and with anything else that writes to it.
const (
	MetaGroups = "human.dev/groups"
	MetaRisk   = "human.dev/risk"
)

func ensureMeta(t *mcp.Tool) {
	if t.Meta == nil {
		t.Meta = &mcp.Meta{}
	}
	if t.Meta.AdditionalFields == nil {
		t.Meta.AdditionalFields = map[string]any{}
	}
}

// InGroup tags a tool with one or more groups. Prefer one: a tool that belongs
// in two groups is usually two tools.
func InGroup(groups ...Group) mcp.ToolOption {
	return func(t *mcp.Tool) {
		ensureMeta(t)
		out := make([]string, 0, len(groups))
		for _, g := range groups {
			out = append(out, string(g))
		}
		t.Meta.AdditionalFields[MetaGroups] = out
	}
}

// WithRisk tags a tool's risk level and is the sole writer of the protocol's
// four annotation hints.
//
// Authors must not set the hints by hand. mcp.NewTool defaults DestructiveHint
// and OpenWorldHint to true and always serializes annotations, so a tool that
// does not override them advertises itself to every client as destructive and
// open-world — which is what all 29 tools did before this existed, read-only
// lookups included. Deriving the hints from one declared risk level makes that
// state unreachable.
func WithRisk(r Risk) mcp.ToolOption {
	return func(t *mcp.Tool) {
		ensureMeta(t)
		t.Meta.AdditionalFields[MetaRisk] = string(r)

		readOnly, destructive, idempotent := false, false, false
		switch r {
		case RiskRead:
			readOnly, destructive, idempotent = true, false, true
		case RiskWrite:
			// Create is not idempotent; update is, but the registry cannot tell
			// them apart, so the conservative answer is the honest one.
			readOnly, destructive, idempotent = false, false, false
		case RiskDestructive:
			// Deleting twice lands in the same state, hence idempotent.
			readOnly, destructive, idempotent = false, true, true
		}

		t.Annotations.ReadOnlyHint = mcp.ToBoolPtr(readOnly)
		t.Annotations.DestructiveHint = mcp.ToBoolPtr(destructive)
		t.Annotations.IdempotentHint = mcp.ToBoolPtr(idempotent)
		// Every Human tool acts on this instance's own data, never on an open
		// world of external entities.
		t.Annotations.OpenWorldHint = mcp.ToBoolPtr(false)
	}
}

// GroupsOf reports the groups a tool is tagged with, or nil if untagged.
func GroupsOf(t mcp.Tool) []Group {
	if t.Meta == nil || t.Meta.AdditionalFields == nil {
		return nil
	}
	raw, ok := t.Meta.AdditionalFields[MetaGroups]
	if !ok {
		return nil
	}
	// Written as []string by InGroup; tolerate []any in case a tool arrives
	// having been through a JSON round-trip.
	switch v := raw.(type) {
	case []string:
		out := make([]Group, 0, len(v))
		for _, g := range v {
			out = append(out, Group(g))
		}
		return out
	case []any:
		out := make([]Group, 0, len(v))
		for _, g := range v {
			if s, ok := g.(string); ok {
				out = append(out, Group(s))
			}
		}
		return out
	}
	return nil
}

// RiskOf reports a tool's declared risk level, or "" if untagged.
func RiskOf(t mcp.Tool) Risk {
	if t.Meta == nil || t.Meta.AdditionalFields == nil {
		return ""
	}
	if s, ok := t.Meta.AdditionalFields[MetaRisk].(string); ok {
		return Risk(s)
	}
	return ""
}

// AtOrBelow reports whether r is within the ceiling. Unknown levels are treated
// as exceeding every ceiling, so an untagged tool cannot slip through.
func (r Risk) AtOrBelow(ceiling Risk) bool {
	rank := map[Risk]int{RiskRead: 1, RiskWrite: 2, RiskDestructive: 3}
	got, ok := rank[r]
	if !ok {
		return false
	}
	max, ok := rank[ceiling]
	if !ok {
		return false
	}
	return got <= max
}
