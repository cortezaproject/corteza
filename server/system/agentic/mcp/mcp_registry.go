package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	rt "github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type (
	registeredTool struct {
		Tool        mcp.Tool
		Handler     server.ToolHandlerFunc
		InputSchema map[string]any
		Title       string
		Hidden      bool
		Available   func() bool // optional; if set, tool is skipped when it returns false
	}

	registeredResource struct {
		Resource mcp.Resource
		Handler  server.ResourceHandlerFunc
	}

	Registry struct {
		tools     map[string]registeredTool
		resources map[string]registeredResource
	}

	// RegisterOption adjusts how a tool is registered. These are
	// registration-time concerns — on which surface a tool appears — as
	// opposed to tool-definition concerns (group, risk, annotations), which
	// travel on the mcp.Tool itself via the options in tool_options.go.
	RegisterOption func(*registeredTool)
)

// Hidden marks a tool as in-process only: the agentic runtime can see and call
// it, remote MCP clients cannot. See NewMCPServer for how this is enforced.
func Hidden() RegisterOption {
	return func(t *registeredTool) { t.Hidden = true }
}

// Available gates a tool on a runtime predicate. It is re-evaluated on every
// listing, so it may depend on live state (a health check, an option).
func Available(fn func() bool) RegisterOption {
	return func(t *registeredTool) { t.Available = fn }
}

// ToolAliases maps old tool names to their current equivalents.
// Add entries here when a tool is renamed so existing agents keep working.
//
// Aliases matter for the in-process surface only: agent definitions store tool
// names by value, so a rename would otherwise strand them. Remote MCP clients
// list tools live and never hold a stale name, which is why aliases are
// deliberately not registered with the mcp-go server — doing so would
// advertise every historical name to every client.
//
// Keep in sync with server/system/agentic/policy/policy.go; the structural
// test asserts the two maps agree.
var ToolAliases = map[string]string{
	"compose_namespace_list": "compose_namespace_lookup",
	"compose_module_list":    "compose_module_lookup",
}

// ResolveToolAlias returns the current tool name for a given name, resolving any alias.
func ResolveToolAlias(name string) string {
	if alias, ok := ToolAliases[name]; ok {
		return alias
	}
	return name
}

func NewRegistry() *Registry {
	return &Registry{
		tools:     make(map[string]registeredTool),
		resources: make(map[string]registeredResource),
	}
}

// inputSchema flattens a tool's declared schema into the shape the agentic
// runtime hands to an LLM.
func inputSchema(tool mcp.Tool) map[string]any {
	schema := map[string]any{
		"type":       tool.InputSchema.Type,
		"properties": tool.InputSchema.Properties,
	}
	if len(tool.InputSchema.Required) > 0 {
		schema["required"] = tool.InputSchema.Required
	}
	return schema
}

// RegisterTool adds a tool to the registry.
//
// Registration happens once, at boot. A duplicate name is a programming error
// — previously it silently overwrote the earlier tool, which meant two
// handlers claiming one name produced a green build with one tool missing —
// so it panics rather than returning an error no caller would check.
func (r *Registry) RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc, opts ...RegisterOption) {
	if _, exists := r.tools[tool.Name]; exists {
		panic(fmt.Sprintf("mcp: duplicate tool registration for %q", tool.Name))
	}

	t := registeredTool{
		Tool:        tool,
		Handler:     handler,
		InputSchema: inputSchema(tool),
		Title:       title,
	}
	for _, opt := range opts {
		opt(&t)
	}

	r.tools[tool.Name] = t
}

// Tools returns every registered tool definition, sorted by name.
//
// GetTools deliberately projects down to rt.Tool for the agentic runtime, which
// drops Meta and annotations. The structural test and the coverage matrix need
// exactly those, so they get the definitions themselves.
func (r *Registry) Tools() []mcp.Tool {
	out := make([]mcp.Tool, 0, len(r.tools))
	for _, t := range r.tools {
		out = append(out, t.Tool)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// IsHidden reports whether a tool is registered as in-process only.
func (r *Registry) IsHidden(name string) bool {
	t, ok := r.tools[ResolveToolAlias(name)]
	return ok && t.Hidden
}

func (r *Registry) HasTool(name string) bool {
	_, ok := r.tools[ResolveToolAlias(name)]
	return ok
}

func (r *Registry) RegisterResource(resource mcp.Resource, handler server.ResourceHandlerFunc) {
	r.resources[resource.URI] = registeredResource{Resource: resource, Handler: handler}
}

func (r *Registry) GetTools(ctx context.Context, allowedTools []string) ([]rt.Tool, error) {
	isAvailable := func(t registeredTool) bool {
		return t.Available == nil || t.Available()
	}

	// nil means no filter — return all non-hidden registered tools
	// @note should this be len(allowedTools) == 0 instead? Depends on the caller logic
	if allowedTools == nil {
		out := make([]rt.Tool, 0, len(r.tools))
		for _, t := range r.tools {
			if t.Hidden || !isAvailable(t) {
				continue
			}
			out = append(out, rt.Tool{
				Name:        t.Tool.Name,
				Title:       t.Title,
				Description: t.Tool.Description,
				InputSchema: t.InputSchema,
			})
		}
		return out, nil
	}

	out := make([]rt.Tool, 0, len(allowedTools))
	for _, name := range allowedTools {
		t, ok := r.tools[ResolveToolAlias(name)]
		if !ok {
			return nil, fmt.Errorf("tool not found: %s", name)
		}
		if !isAvailable(t) {
			continue
		}
		out = append(out, rt.Tool{
			Name:        t.Tool.Name,
			Title:       t.Title,
			Description: t.Tool.Description,
			InputSchema: t.InputSchema,
		})
	}
	return out, nil
}

func (r *Registry) ExecuteTool(ctx context.Context, toolName string, args map[string]any) (any, error) {
	// Resolve first: agent definitions store tool names by value, so a stored
	// pre-rename name reaches dispatch verbatim. GetTools and HasTool already
	// resolve, so without this an aliased tool lists fine and fails on call.
	toolName = ResolveToolAlias(toolName)

	t, ok := r.tools[toolName]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", toolName)
	}

	req := mcp.CallToolRequest{}
	req.Params.Name = toolName
	req.Params.Arguments = args

	result, err := t.Handler(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, c := range result.Content {
		switch v := c.(type) {
		case mcp.TextContent:
			var parsed any
			if err := json.Unmarshal([]byte(v.Text), &parsed); err == nil {
				return parsed, nil
			}
			return v.Text, nil
		default:
			return nil, fmt.Errorf("tool %q returned unsupported content type %T", toolName, c)
		}
	}

	return nil, nil
}
