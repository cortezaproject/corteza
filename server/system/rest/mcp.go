package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/system/agentic/runtime"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
)

type (
	Mcp struct {
		registry mcpRegistry
	}

	mcpRegistry interface {
		GetTools(ctx context.Context, agentID uint64) ([]runtime.Tool, error)
	}
)

func (Mcp) New() Mcp {
	return Mcp{registry: service.DefaultMCPRegistry}
}

func (ctrl Mcp) ListTools(ctx context.Context, _ *request.McpListTools) (interface{}, error) {
	return ctrl.registry.GetTools(ctx, 0)
}
