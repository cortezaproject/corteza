package automation

import (
	"context"

	agenticRuntime "github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/types"
)

type (
	agentHandler struct {
		reg     agentHandlerRegistry
		runtime agentHandlerRuntime
		agents  agentHandlerAgentService
	}

	agentHandlerRuntime interface {
		Run(ctx context.Context, req *agenticRuntime.AgentRequest) (*agenticRuntime.AgentResponse, error)
	}

	agentHandlerAgentService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Agent, error)
	}
)

func AgentHandler(reg agentHandlerRegistry, rt agentHandlerRuntime, agents agentHandlerAgentService) *agentHandler {
	h := &agentHandler{reg: reg, runtime: rt, agents: agents}
	h.register()
	return h
}

// run invokes the agent as whoever the workflow runs as; the workflow's own
// runAs is the identity, the same way it is for every other step.
func (h agentHandler) run(ctx context.Context, args *agentRunArgs) (*agentRunResults, error) {
	if _, err := h.agents.FindByID(ctx, args.AgentID); err != nil {
		return nil, err
	}

	resp, err := h.runtime.Run(ctx, &agenticRuntime.AgentRequest{
		AgentID:        args.AgentID,
		Unattended:     true,
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
