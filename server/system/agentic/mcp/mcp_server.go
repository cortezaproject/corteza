package mcp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// mcpBasePath is where MountRoutes is mounted under the API router. The scope
// parser needs it to find the group segment that follows.
const mcpBasePath = "/api/mcp"

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
//
// On top of that, each request carries a Scope resolved from its URL, which
// narrows by group and caps by risk. See scope.go.
func NewMCPServer(reg *Registry, name, version string) *MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
		server.WithToolFilter(listFilter(reg)),
		server.WithToolHandlerMiddleware(riskCeiling(reg)),
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

	httpServer := server.NewStreamableHTTPServer(s,
		server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			return WithScope(ctx, scopeFromRequest(mcpBasePath, r))
		}),
	)

	return &MCPServer{server: s, httpServer: httpServer}
}

// listFilter drops tools the request's scope does not cover, and tools whose
// Available predicate currently reports false.
func listFilter(reg *Registry) server.ToolFilterFunc {
	return func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
		scope := ScopeFromContext(ctx)

		out := make([]mcp.Tool, 0, len(tools))
		for _, tool := range tools {
			if t, ok := reg.tools[tool.Name]; ok && t.Available != nil && !t.Available() {
				continue
			}
			if !scope.Permits(GroupsOf(tool), RiskOf(tool)) {
				continue
			}
			out = append(out, tool)
		}
		return out
	}
}

// riskCeiling rejects a call above the session's declared ceiling.
//
// The filter above only hides tools from list_tools; a client that already knew
// a name could still call it. The ceiling has to refuse at dispatch to mean
// anything, which is why this middleware exists as well as the filter.
//
// Group is deliberately NOT enforced here. It is a presentation concern, and a
// caller who names a tool outside its listed group has done nothing RBAC would
// not already permit.
func riskCeiling(reg *Registry) server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope := ScopeFromContext(ctx)
			if scope.MaxRisk == "" {
				return next(ctx, req)
			}

			t, ok := reg.tools[ResolveToolAlias(req.Params.Name)]
			if !ok {
				return next(ctx, req)
			}

			if risk := RiskOf(t.Tool); !risk.AtOrBelow(scope.MaxRisk) {
				return nil, fmt.Errorf(
					"tool %q is %s but this session is capped at %s; reconnect without maxRisk to use it",
					req.Params.Name, risk, scope.MaxRisk,
				)
			}

			return next(ctx, req)
		}
	}
}

func (m *MCPServer) MountRoutes(r chi.Router) {
	r.Handle("/*", m.httpServer)
}
