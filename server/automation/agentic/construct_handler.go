package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	autoService "github.com/crusttech/human/server/automation/service"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

type constructHandler struct {
	reg toolRegistrar
}

// The declarations for this handler are in construct_tools.go.
func ConstructHandler(reg toolRegistrar) *constructHandler {
	h := &constructHandler{reg: reg}
	h.register()
	return h
}

// There is no authorization question to answer for either catalogue, and the
// reason is worth stating rather than the conclusion.
//
// Both are read straight off the process-wide registries the TAQ and workflow
// editors already read over REST, with no per-caller dimension: the same 19
// construct functions and 93 workflow functions for everybody. The construct
// library is not quite a compile-time constant — a configured connection adds
// one function per operation and one trigger per webhook — but what it adds is a
// ref, a label and a parameter list, the identical payload
// /automation/construct-library/functions already serves to any authenticated
// caller. There is no per-caller narrowing to preserve and no secret in it.

// taqConstructLookup returns the construct library: TAQ step refs and trigger
// pairs.
func (h *constructHandler) taqConstructLookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	var (
		lib       = autoService.ConstructLibrary()
		catalogue = strings.ToLower(strings.TrimSpace(toolkit.Str(args, "catalogue")))
		out       = map[string]any{}
	)

	switch catalogue {
	case "", "functions", "triggers":
	default:
		return nil, fmt.Errorf("construct library lookup failed: unknown catalogue %q; it is \"functions\", \"triggers\", or omitted for both", catalogue)
	}

	if catalogue != "triggers" {
		fns, err := slimSet(lib.Functions())
		if err != nil {
			return nil, err
		}

		refs := refsOf(fns)
		if q := strings.TrimSpace(toolkit.Str(args, "ref")); q != "" {
			matched := filterByRef(fns, q)
			// A ref that matches nothing is the likely failure — the names are
			// not guessable, and guessing is what puts an unknown ref in a step
			// — so answer with the ones that do exist.
			if len(matched) == 0 {
				out["knownRefs"] = refs
			}
			fns = matched
		}

		out["functions"] = fns
		out["functionCount"] = len(fns)
	}

	if catalogue != "functions" {
		triggers, err := slimSet(lib.Triggers())
		if err != nil {
			return nil, err
		}

		if q := strings.TrimSpace(toolkit.Str(args, "resourceType")); q != "" {
			matched := make([]map[string]any, 0, len(triggers))
			for _, t := range triggers {
				if matchesResourceType(t, q) {
					matched = append(matched, t)
				}
			}
			if len(matched) == 0 {
				out["knownResourceTypes"] = knownResourceTypes(triggers)
			}
			triggers = matched
		}

		out["triggers"] = triggers
		out["triggerCount"] = len(triggers)
	}

	return toolkit.JSONResult(out)
}

// workflowFunctionLookup returns the workflow function registry.
func (h *constructHandler) workflowFunctionLookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	fns, err := slimSet(autoService.Registry().Functions())
	if err != nil {
		return nil, err
	}

	out := map[string]any{}
	if q := strings.TrimSpace(toolkit.Str(args, "ref")); q != "" {
		refs := refsOf(fns)
		matched := filterByRef(fns, q)
		if len(matched) == 0 {
			out["knownRefs"] = refs
		}
		fns = matched
	}

	out["functions"] = fns
	out["count"] = len(fns)

	return toolkit.JSONResult(out)
}

// slimSet re-reads a catalogue as generic JSON and drops the webapp's form
// spec.
//
// Generic rather than a typed mirror for the same reason the event-type handler
// gives: a parameter added to a definition passes through untouched instead of
// being silently dropped by a local struct that has not been updated. What is
// removed is named explicitly — "segments" is the input layout the TAQ editor
// draws, four fifths of the payload, and says nothing an author binding
// arguments by argumentName needs.
func slimSet(set any) ([]map[string]any, error) {
	raw, err := json.Marshal(set)
	if err != nil {
		return nil, fmt.Errorf("construct catalogue decode failed: %w", err)
	}

	var decoded []map[string]any
	if err = json.Unmarshal(raw, &decoded); err != nil {
		return nil, fmt.Errorf("construct catalogue decode failed: %w", err)
	}

	for _, def := range decoded {
		delete(def, "segments")
	}
	return decoded, nil
}

// filterByRef matches an exact ref first, then any ref containing the query.
// Exact first matters: "loopDo" is a substring of nothing but "composeRecords"
// matches six, and a caller naming one exactly wants that one.
func filterByRef(set []map[string]any, query string) []map[string]any {
	for _, def := range set {
		if ref, _ := def["ref"].(string); strings.EqualFold(ref, query) {
			return []map[string]any{def}
		}
	}

	out := make([]map[string]any, 0, len(set))
	for _, def := range set {
		if ref, _ := def["ref"].(string); strings.Contains(strings.ToLower(ref), strings.ToLower(query)) {
			out = append(out, def)
		}
	}
	return out
}

func refsOf(set []map[string]any) []string {
	out := make([]string, 0, len(set))
	for _, def := range set {
		if ref, _ := def["ref"].(string); ref != "" {
			out = append(out, ref)
		}
	}
	sort.Strings(out)
	return out
}
