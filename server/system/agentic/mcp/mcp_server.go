package mcp

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server     *server.MCPServer
	httpServer *server.StreamableHTTPServer
}

// NewMCPServer exposes the registry over the HTTP (streamable) transport.
//
// Two registry flags decide what a remote client may reach, and they are
// enforced differently on purpose:
//
//   - Hidden marks a tool as in-process only. It is never added to the mcp-go
//     server, so it is neither listed nor callable over HTTP. Filtering alone
//     would not do: mcp-go tool filters apply to list_tools, not to dispatch.
//   - Available is a runtime predicate (discovery_search health-checks the
//     discovery service on every call), so it cannot be evaluated once here at
//     boot. It goes through a tool filter, which mcp-go re-runs per list_tools
//     request.
func NewMCPServer(reg *Registry, name, version string) *MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
		server.WithToolFilter(availabilityFilter(reg)),
	)

	for _, t := range reg.tools {
		if t.Hidden {
			continue
		}
		s.AddTool(t.Tool, t.Handler)
	}

	for _, r := range reg.resources {
		s.AddResource(r.Resource, r.Handler)
	}

	httpServer := server.NewStreamableHTTPServer(s)

	return &MCPServer{
		server:     s,
		httpServer: httpServer,
	}
}

// availabilityFilter drops tools whose Available predicate currently reports
// false. Tools without a predicate always pass.
func availabilityFilter(reg *Registry) server.ToolFilterFunc {
	return func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
		out := make([]mcp.Tool, 0, len(tools))
		for _, tool := range tools {
			t, ok := reg.tools[tool.Name]
			if ok && t.Available != nil && !t.Available() {
				continue
			}
			out = append(out, tool)
		}
		return out
	}
}

func (m *MCPServer) MountRoutes(r chi.Router) {
	r.Handle("/*", m.httpServer)
}
