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
// On top of that, every tool in the listing is summarised unless the request
// asks for full documentation — see listing.go for why.
func NewMCPServer(reg *Registry, name, version string) *MCPServer {
	m := &MCPServer{reg: reg}

	m.server = server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
		server.WithToolFilter(m.listFilter()),
		server.WithToolHandlerMiddleware(m.riskCeiling()),
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

// listFilter decides what a given request sees: every available, in-scope tool,
// summarised unless full documentation was asked for.
//
// mcp-go assigns this function's return value straight into the result, so it
// rewrites as well as filters. That is the whole implementation of slimming —
// no second transport and no protocol extension, because a listing is the only
// place a tool description is ever sent.
func (m *MCPServer) listFilter() server.ToolFilterFunc {
	return func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
		scope := ScopeFromContext(ctx)

		out := make([]mcp.Tool, 0, len(tools))
		for _, tool := range tools {
			if t, ok := m.reg.tools[tool.Name]; ok && t.Available != nil && !t.Available() {
				continue
			}
			if !scope.Permits(GroupsOf(tool), RiskOf(tool)) {
				continue
			}
			if !scope.FullDocs {
				tool = slimTool(tool)
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
// Group is not enforced here. It is presentation: a caller who names a tool
// outside its group has done nothing RBAC would not already permit.
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
				"Find the right tool for what you are trying to do, and get its full documentation. "+
					"Every tool on this server is already listed and callable, but the listing shows only "+
					"a one-line summary of each and omits the per-parameter detail; this returns the "+
					"complete description and parameter documentation for the tools that match. Use it "+
					"when you are unsure which tool does the job, and whenever a tool's summary tells you "+
					"to load it first. Query by resource or action, in as many words as you like: "+
					"\"page\", \"delete role\", \"create a user\", \"workflow\". Tools matching more of "+
					"your words rank first, so an extra word narrows the ranking rather than emptying the "+
					"result. If you get nothing back, try a different word.",
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
				"Get the full documentation for tools you can already name, without searching. The tool "+
					"list summarises every tool to one line and drops the per-parameter detail; this "+
					"returns it in full. Call this before using any tool whose summary says to, and any "+
					"time a call was rejected for arguments you are unsure about. Use "+toolSearchName+
					" instead when you know what you want to do but not what it is called.",
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
	return m.toolDefsResult(hits, fmt.Sprintf("no tool matches %q — try one broader word", query))
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

	return m.toolDefsResult(hits, "")
}

// toolDefsResult returns full documentation for the named tools.
//
// It changes nothing about what the session can call. Every tool is already
// listed and callable; the listing merely summarised it. That is what makes
// this work on clients where the old disclosure did not — nothing has to be
// registered, so no client has to cooperate, and the result is only text.
func (m *MCPServer) toolDefsResult(hits []mcp.Tool, emptyMsg string) (*mcp.CallToolResult, error) {
	if len(hits) == 0 {
		if emptyMsg == "" {
			emptyMsg = "nothing matched"
		}
		return mcp.NewToolResultText(emptyMsg), nil
	}

	out, err := json.Marshal(map[string]any{
		"tools": hits,
		"note": "Full documentation for these tools, including the parameter detail the tool " +
			"list leaves out. They were already callable — this adds the guidance, not the tools.",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tools: %w", err)
	}

	return mcp.NewToolResultText(string(out)), nil
}

func (m *MCPServer) MountRoutes(r chi.Router) {
	r.Handle("/*", m.httpServer)
}

// ServeStdio runs the same server over stdio, for a client that launches it as
// a child process — which is how the developer MCP under dev/ is started.
//
// Everything a session sees is unchanged: the tool filter, the risk-ceiling
// middleware, the meta-tools and progressive disclosure all hang off the mcp-go
// server, not off the transport. The one thing that does not carry over is
// Scope, which is parsed from the request URL and has nowhere to come from on a
// pipe. The zero Scope permits everything, which is the right default here: a
// server the developer launched themselves, in their own checkout, has no
// caller to narrow.
func (m *MCPServer) ServeStdio() error {
	return server.ServeStdio(m.server)
}
