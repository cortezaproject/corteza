package automation

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/pkg/auth"
	agenticRuntime "github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/types"
)

type (
	agentHandler struct {
		reg     agentHandlerRegistry
		runtime agentHandlerRuntime
		agents  agentHandlerAgentService
		users   agentHandlerUserService
	}

	agentHandlerRuntime interface {
		Run(ctx context.Context, req *agenticRuntime.AgentRequest) (*agenticRuntime.AgentResponse, error)
	}

	agentHandlerAgentService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Agent, error)
	}

	agentHandlerUserService interface {
		FindByID(ctx context.Context, ID uint64) (*types.User, error)
	}
)

func AgentHandler(reg agentHandlerRegistry, rt agentHandlerRuntime, agents agentHandlerAgentService, users agentHandlerUserService) *agentHandler {
	h := &agentHandler{reg: reg, runtime: rt, agents: agents, users: users}
	h.register()
	return h
}

func (h agentHandler) run(ctx context.Context, args *agentRunArgs) (*agentRunResults, error) {
	a, err := h.agents.FindByID(ctx, args.AgentID)
	if err != nil {
		return nil, err
	}

	if !a.Invocation.System.Enabled {
		return nil, fmt.Errorf("agent is not available for system invocation")
	}

	if a.Invocation.System.ServiceAccount != 0 {
		u, err := h.users.FindByID(ctx, a.Invocation.System.ServiceAccount)
		if err != nil {
			return nil, fmt.Errorf("could not resolve service account %d: %w", a.Invocation.System.ServiceAccount, err)
		}
		ctx = auth.SetIdentityToContext(ctx, auth.Authenticated(u.ID, u.Roles()...))
	}

	resp, err := h.runtime.Run(ctx, &agenticRuntime.AgentRequest{
		AgentID:        args.AgentID,
		Input:          args.Input,
		ConversationID: args.ConversationID,
	})
	if err != nil {
		return nil, err
	}

	return &agentRunResults{
		Output:         resp.Output,
		ConversationID: resp.ConversationID,
	}, nil
}
