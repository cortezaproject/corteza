// Package mcp is the Human side of the MCP surface: the tool families, the
// policy that classifies them, and the registry the agentic runtime talks to.
//
// The transport-level machinery — the registry itself, tagging, scope,
// progressive disclosure and the HTTP server — lives in pkg/mcpkit, which is
// shared with the developer MCP under dev/ and may not import anything from
// Human's domain. What stays here is what does know about Human.
package mcp

import (
	"context"
	"strings"

	herrors "github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	rt "github.com/crusttech/human/server/system/agentic/runtime"
)

// Registry is mcpkit's registry plus the one thing that cannot live there: a
// projection into the agentic runtime's own Tool type.
//
// mcpkit is shared with a server that knows nothing about Human, so it cannot
// name rt.Tool. Embedding rather than wrapping keeps every registration and
// execution method reachable, so this type is a drop-in wherever the registry
// was passed before.
type Registry struct {
	*mcpkit.Registry
}

// NewRegistry builds the Human registry.
func NewRegistry() *Registry {
	r := &Registry{Registry: mcpkit.NewRegistry()}
	r.SetErrorClassifier(ClassifyError)
	return r
}

// ClassifyError names the code a Human error carries on the wire.
//
// Generated action errors (NamespaceErrNotFound, ModuleErrNotAllowedToUpdate…)
// are all KindInternal and tell their kind through the "type" meta, so that is
// read first; the error kinds cover what the store and auth layers raise
// directly. "" means mcpkit falls back to "failed".
func ClassifyError(err error) (code, next string) {
	var e *herrors.Error
	if herrors.As(err, &e) {
		if t, ok := e.MetaValue("type"); ok {
			switch s, _ := t.(string); {
			case s == "notFound":
				return toolkit.CodeNotFound, ""
			case strings.HasPrefix(s, "notAllowedTo"):
				return toolkit.CodeForbidden, ""
			case s == "invalidID" || s == "invalidHandle":
				return toolkit.CodeInvalid, "the identifier is malformed; IDs are numeric strings and handles are letters, digits and underscores"
			case s == "staleData":
				return toolkit.CodeConflict, ""
			}
		}
	}

	switch {
	case herrors.IsNotFound(err):
		return toolkit.CodeNotFound, ""
	case herrors.IsUnauthorized(err):
		return toolkit.CodeForbidden, ""
	case herrors.IsUnauthenticated(err):
		return toolkit.CodeUnauthenticated, ""
	case herrors.IsDuplicateData(err), herrors.IsStaleData(err):
		return toolkit.CodeConflict, ""
	case herrors.IsInvalidData(err):
		return toolkit.CodeInvalid, ""
	}
	return "", ""
}

// GetTools projects the registry down to the agentic runtime's own Tool type,
// which drops Meta and annotations.
func (r *Registry) GetTools(ctx context.Context, allowedTools []string) ([]rt.Tool, error) {
	defs, err := r.Select(allowedTools)
	if err != nil {
		return nil, err
	}

	out := make([]rt.Tool, 0, len(defs))
	for _, d := range defs {
		groups := make([]string, 0, len(d.Groups))
		for _, g := range d.Groups {
			groups = append(groups, string(g))
		}

		out = append(out, rt.Tool{
			Name:        d.Name,
			Title:       d.Title,
			Description: d.Description,
			InputSchema: d.InputSchema,
			Groups:      groups,
			Risk:        string(d.Risk),
		})
	}

	return out, nil
}

// ToolNamesIn projects mcpkit's group/risk lookup into the plain strings the
// agentic runtime's MCPClient interface speaks, so the runtime does not have to
// import mcpkit's vocabulary.
func (r *Registry) ToolNamesIn(group, maxRisk string) []string {
	return r.Registry.ToolNamesIn(mcpkit.Group(group), mcpkit.Risk(maxRisk))
}
