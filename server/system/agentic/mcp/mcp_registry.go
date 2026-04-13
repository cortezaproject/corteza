package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	rt "github.com/cortezaproject/corteza/server/system/agentic/runtime"
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
)

// ToolAliases maps old tool names to their current equivalents.
// Add entries here when a tool is renamed so existing agents keep working.
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

func (r *Registry) RegisterHiddenTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc) {
	schema := map[string]any{
		"type":       tool.InputSchema.Type,
		"properties": tool.InputSchema.Properties,
	}
	if len(tool.InputSchema.Required) > 0 {
		schema["required"] = tool.InputSchema.Required
	}
	r.tools[tool.Name] = registeredTool{Tool: tool, Handler: handler, InputSchema: schema, Title: title, Hidden: true}
}

func (r *Registry) RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc) {
	schema := map[string]any{
		"type":       tool.InputSchema.Type,
		"properties": tool.InputSchema.Properties,
	}
	if len(tool.InputSchema.Required) > 0 {
		schema["required"] = tool.InputSchema.Required
	}

	r.tools[tool.Name] = registeredTool{Tool: tool, Handler: handler, InputSchema: schema, Title: title}
}

func (r *Registry) RegisterToolWithAvailability(tool mcp.Tool, title string, handler server.ToolHandlerFunc, available func() bool) {
	schema := map[string]any{
		"type":       tool.InputSchema.Type,
		"properties": tool.InputSchema.Properties,
	}
	if len(tool.InputSchema.Required) > 0 {
		schema["required"] = tool.InputSchema.Required
	}

	r.tools[tool.Name] = registeredTool{Tool: tool, Handler: handler, InputSchema: schema, Title: title, Available: available}
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
