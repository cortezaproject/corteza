package automation

import (
	"context"
	"fmt"

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
		Search(ctx context.Context, f types.AgentFilter) (types.AgentSet, types.AgentFilter, error)
	}

	agentLookup interface {
		GetLookup() (bool, uint64, string, *types.Agent)
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

func (h agentHandler) lookup(ctx context.Context, args *agentLookupArgs) (results *agentLookupResults, err error) {
	results = &agentLookupResults{}
	results.Agent, err = lookupAgent(ctx, h.agents, args)
	return
}

// lookupAgent resolves an agent given as a resource, an ID or a handle.
//
// Agent handles are not unique, so a handle resolves to the first agent the
// search returns; a workflow that needs one agent in particular names its ID.
func lookupAgent(ctx context.Context, svc agentHandlerAgentService, args agentLookup) (*types.Agent, error) {
	_, ID, handle, agent := args.GetLookup()

	switch {
	case agent != nil:
		return agent, nil
	case ID > 0:
		return svc.FindByID(ctx, ID)
	case len(handle) > 0:
		set, _, err := svc.Search(ctx, types.AgentFilter{Handle: handle})
		if err != nil {
			return nil, err
		}
		if len(set) == 0 {
			return nil, fmt.Errorf("agent %q not found", handle)
		}
		return set[0], nil
	}

	return nil, fmt.Errorf("empty lookup params")
}
