package mcp

import (
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