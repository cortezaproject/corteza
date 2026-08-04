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

	"github.com/crusttech/human/server/pkg/mcpkit"
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
	return &Registry{Registry: mcpkit.NewRegistry()}
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
		out = append(out, rt.Tool{
			Name:        d.Name,
			Title:       d.Title,
			Description: d.Description,
			InputSchema: d.InputSchema,
		})
	}

	return out, nil
}
