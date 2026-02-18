package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	rt "github.com/cortezaproject/corteza/server/pkg/agentic/runtime"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type (
	registeredTool struct {
		Tool    mcp.Tool
		Handler server.ToolHandlerFunc
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

func (r *Registry) RegisterTool(tool mcp.Tool, handler server.ToolHandlerFunc){
	r.tools[tool.Name]=registeredTool{Tool: tool, Handler: handler}
}

func (r *Registry) RegisterResource(resource mcp.Resource, handler server.ResourceHandlerFunc){
	r.resources[resource.URI] = registeredResource{Resource: resource, Handler: handler}
}

func (r *Registry) GetTools(ctx context.Context, agentID uint64) ([]rt.Tool, error) {
	out := make([]rt.Tool, 0, len(r.tools))
	for _, t := range r.tools {
		schema := map[string]any{
			"type":       t.Tool.InputSchema.Type,
			"properties": t.Tool.InputSchema.Properties,
		}
		if len(t.Tool.InputSchema.Required) > 0 {
			schema["required"] = t.Tool.InputSchema.Required
		}

		out = append(out, rt.Tool{
			Name:        t.Tool.Name,
			Description: t.Tool.Description,
			InputSchema: schema,
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
		if text, ok := c.(mcp.TextContent); ok {
			var parsed any
			if err := json.Unmarshal([]byte(text.Text), &parsed); err == nil {
				return parsed, nil
			}
			return text.Text, nil
		}
	}

	return nil, nil
}