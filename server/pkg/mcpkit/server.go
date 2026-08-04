package mcpkit

import (
	"context"
	"encoding/json"
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
	reg        *Registry
	disclosed  *disclosure
}

// NewMCPServer exposes the registry over the HTTP (streamable) transport.
//
// Three things decide what a remote client sees, and they are enforced
// differently on purpose:
//
//   - Hidden marks a tool as in-process only. It is never added to the mcp-go
//     server, so it is neither listed nor callable over HTTP. Filtering alone
//     would not do: mcp-go tool filters apply to list_tools, not to dispatch.
//   - Available is a runtime predicate (discovery_search health-checks the
//     discovery service on every call), so it cannot be evaluated once here at
//     boot. It goes through the tool filter, which mcp-go re-runs per request.
//   - Scope (group, risk ceiling) comes from the request URL; see scope.go.
//
// On top of that, a session sees only the always-on tools until it searches for
// more — see disclosure.go for why.
func NewMCPServer(reg *Registry, name, version string) *MCPServer {
	m := &MCPServer{reg: reg, disclosed: newDisclosure()}

	hooks := &server.Hooks{}
	hooks.AddOnUnregisterSession(func(_ context.Context, s server.ClientSession) {
		m.disclosed.forget(s.SessionID())
	})

	m.server = server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
		server.WithToolFilter(m.listFilter()),
		server.WithToolHandlerMiddleware(m.riskCeiling()),
		server.WithHooks(hooks),
	)

	for _, t := range reg.tools {
		if t.Hidden {
			continue
		}
		m.server.AddTool(t.Tool, t.Handler)
	}

	for _, r := range reg.resources {
		m.server.AddResource(r.Resource, r.Handler)
	}

	m.registerMetaTools()

	m.httpServer = server.NewStreamableHTTPServer(m.server,
		server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			return WithScope(ctx, scopeFromRequest(mcpBasePath, r))
		}),
	)

	return m
}

// listFilter decides what a given request may see: available, in scope, and
// either always-on or already pulled in by this session.
func (m *MCPServer) listFilter() server.ToolFilterFunc {
	return func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
		var (
			scope = ScopeFromContext(ctx)
			sid   = sessionID(ctx)
		)

		out := make([]mcp.Tool, 0, len(tools))
		for _, tool := range tools {
			if t, ok := m.reg.tools[tool.Name]; ok && t.Available != nil && !t.Available() {
				continue
			}
			if !scope.Permits(GroupsOf(tool), RiskOf(tool)) {
				continue
			}
			if !scope.AllTools && !alwaysOn[tool.Name] && !m.disclosed.isLoaded(sid, tool.Name) {
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
// Neither group nor disclosure is enforced here. Both are presentation: a
// caller who names a tool it was never shown has done nothing RBAC would not
// already permit, and refusing would break the very clients that ignore
// list_changed and call straight from a search result.
func (m *MCPServer) riskCeiling() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scope := ScopeFromContext(ctx)
			if scope.MaxRisk == "" {
				return next(ctx, req)
			}

			t, ok := m.reg.tools[ResolveToolAlias(req.Params.Name)]
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

// registerMetaTools adds the search and load tools.
//
// They go straight onto the mcp-go server rather than into the Registry
// because they are a property of this transport, not of Human: the in-process
// runtime scopes an agent with allowedTools and has no listing to shrink.
// Keeping them out of the Registry also keeps them out of the coverage matrix,
// where they would read as resources they are not.
func (m *MCPServer) registerMetaTools() {
	m.server.AddTool(
		mcp.NewTool(toolSearchName,
			mcp.WithDescription(
				"Find tools by what you are trying to do, and load them for this session. "+
					"Only a few tools are listed up front — this server has far more than it advertises, "+
					"because listing them all would cost tens of thousands of tokens per request. "+
					"Search returns each matching tool's full definition, so you can call anything it "+
					"returns immediately; the tools also become visible in the tool list from now on. "+
					"Query by resource or action: \"page\", \"delete role\", \"create user\", \"workflow\". "+
					"All terms must match. If you get nothing back, try one broader word.",
			),
			mcp.WithString("query", mcp.Required(), mcp.Description("What you want to do, e.g. \"page\", \"delete role\", \"reminder snooze\".")),
			InGroup(GroupConfiguring, GroupUsage),
			WithRisk(RiskRead),
		),
		m.handleToolSearch,
	)

	m.server.AddTool(
		mcp.NewTool(toolLoadName,
			mcp.WithDescription(
				"Load named tools into this session when you already know their names, without searching. "+
					"Returns their full definitions. Use "+toolSearchName+" instead when you know what you "+
					"want to do but not what it is called.",
			),
			mcp.WithString("names", mcp.Required(), mcp.Description("JSON array of exact tool names, e.g. [\"system_role_create\",\"system_role_member_add\"].")),
			InGroup(GroupConfiguring, GroupUsage),
			WithRisk(RiskRead),
		),
		m.handleToolLoad,
	)
}

func (m *MCPServer) handleToolSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid request: arguments must be an object")
	}

	query, _ := args["query"].(string)
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	hits := m.searchTools(query, ScopeFromContext(ctx))
	return m.discloseResult(ctx, hits, fmt.Sprintf("no tool matches %q — try one broader word", query))
}

func (m *MCPServer) handleToolLoad(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid request: arguments must be an object")
	}

	var names []string
	switch v := args["names"].(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &names); err != nil {
			return nil, fmt.Errorf("names must be a JSON array of tool names: %w", err)
		}
	case []any:
		for _, n := range v {
			if s, ok := n.(string); ok {
				names = append(names, s)
			}
		}
	default:
		return nil, fmt.Errorf("names must be a JSON array of tool names")
	}

	var (
		scope   = ScopeFromContext(ctx)
		hits    []mcp.Tool
		unknown []string
	)

	for _, name := range names {
		t, ok := m.reg.tools[ResolveToolAlias(name)]
		if !ok || t.Hidden || !scope.Permits(GroupsOf(t.Tool), RiskOf(t.Tool)) {
			unknown = append(unknown, name)
			continue
		}
		hits = append(hits, t.Tool)
	}

	if len(unknown) > 0 && len(hits) == 0 {
		return nil, fmt.Errorf("no such tool: %v — use %s to find one", unknown, toolSearchName)
	}

	return m.discloseResult(ctx, hits, "")
}

// discloseResult marks the tools loaded for this session, tells the client its
// list changed, and returns the definitions inline.
//
// Returning them inline is deliberate belt and braces. A client that honours
// notifications/tools/list_changed will re-list and see them properly; one that
// ignores it still has everything it needs to make the call, straight out of
// this result. Without that, disclosure would be silently broken on any client
// that does not re-list, and the failure would look like the tool not existing.
func (m *MCPServer) discloseResult(ctx context.Context, hits []mcp.Tool, emptyMsg string) (*mcp.CallToolResult, error) {
	if len(hits) == 0 {
		if emptyMsg == "" {
			emptyMsg = "nothing matched"
		}
		return mcp.NewToolResultText(emptyMsg), nil
	}

	names := make([]string, 0, len(hits))
	for _, t := range hits {
		names = append(names, t.Name)
	}

	if sid := sessionID(ctx); sid != "" {
		m.disclosed.load(sid, names...)
		// Best effort: a client that cannot receive it still has the payload.
		_ = m.server.SendNotificationToSpecificClient(sid, "notifications/tools/list_changed", nil)
	}

	out, err := json.Marshal(map[string]any{
		"tools": hits,
		"note": "These tools are now available for the rest of this session. " +
			"You can call them straight away using the schemas above.",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tools: %w", err)
	}

	return mcp.NewToolResultText(string(out)), nil
}

func (m *MCPServer) MountRoutes(r chi.Router) {
	r.Handle("/*", m.httpServer)
}
