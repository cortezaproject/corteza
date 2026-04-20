package rest

import (
	"context"

	"github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
)

type (
	Mcp struct {
		registry mcpRegistry
	}

	mcpRegistry interface {
		GetTools(ctx context.Context, allowedTools []string) ([]runtime.Tool, error)
	}
)

func (Mcp) New() Mcp {
	return Mcp{registry: service.DefaultMCPRegistry}
}

func (ctrl Mcp) ListTools(ctx context.Context, _ *request.McpListTools) (interface{}, error) {
	return ctrl.registry.GetTools(ctx, nil)
}
