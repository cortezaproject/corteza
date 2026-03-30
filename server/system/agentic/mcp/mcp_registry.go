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

func NewRegistry() *Registry {
	return &Registry{
		tools:     make(map[string]registeredTool),
		resources: make(map[string]registeredResource),
	}
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

func (r *Registry) HasTool(name string) bool {
	_, ok := r.tools[name]
	return ok
}

func (r *Registry) RegisterResource(resource mcp.Resource, handler server.ResourceHandlerFunc){
	r.resources[resource.URI] = registeredResource{Resource: resource, Handler: handler}
}

func (r *Registry) GetTools(ctx context.Context, allowedTools []string) ([]rt.Tool, error) {
	// nil means no filter — return all registered tools
	if allowedTools == nil {
		out := make([]rt.Tool, 0, len(r.tools))
		for _, t := range r.tools {
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
		t, ok := r.tools[name]
		if !ok {
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