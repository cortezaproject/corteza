package rest

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/agentic/tcl"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	Agent struct {
		agent agentService
		ac    agentAccessController
	}

	agentPayload struct {
		*types.Agent

		CanGrant       bool `json:"canGrant"`
		CanUpdateAgent bool `json:"canUpdateAgent"`
		CanDeleteAgent bool `json:"canDeleteAgent"`
	}

	agentSetPayload struct {
		Filter types.AgentFilter `json:"filter"`
		Set    []*agentPayload   `json:"set"`
	}

	agentService interface {
		FindByID(ctx context.Context, ID uint64) (a *types.Agent, err error)
		Create(ctx context.Context, new *types.Agent) (a *types.Agent, err error)
		Update(ctx context.Context, upd *types.Agent) (a *types.Agent, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		UndeleteByID(ctx context.Context, ID uint64) (err error)
		Search(ctx context.Context, filter types.AgentFilter) (set types.AgentSet, f types.AgentFilter, err error)
	}

	agentAccessController interface {
		CanGrant(context.Context) bool

		CanCreateAgent(context.Context) bool
		CanUpdateAgent(context.Context, *types.Agent) bool
		CanDeleteAgent(context.Context, *types.Agent) bool
	}
)

func (Agent) New() *Agent {
	return &Agent{
		agent: service.DefaultAgent,
		ac:    service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl *Agent) makeFilter(ctx context.Context, r *request.AgentList) (types.AgentFilter, error) {
	var (
		err error
		f   = types.AgentFilter{
			Query:     r.Query,
			ProjectID: r.ProjectID,
			Handle:    r.Handle,
			Status:    r.Status,
			Labels:    r.Labels,
			Deleted:   filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl *Agent) beforeCreate(ctx context.Context, res *types.Agent, r *request.AgentCreate) error {
	// Project-scoped: stamp the owning project so the agent is created inside it
	// (the generated controller maps only the simple value params).
	res.ProjectID = r.ProjectID
	res.Meta = r.Meta
	res.Behavior = r.Behavior
	res.Execution = r.Execution
	res.Access = r.Access
	res.Invocation = r.Invocation
	res.Labels = r.Labels
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl *Agent) beforeUpdate(ctx context.Context, res *types.Agent, r *request.AgentUpdate) error {
	res.Meta = r.Meta
	res.Behavior = r.Behavior
	res.Execution = r.Execution
	res.Access = r.Access
	res.Invocation = r.Invocation
	res.Labels = r.Labels
	return nil
}

func (ctrl *Agent) Exec(ctx context.Context, r *request.AgentExec) (interface{}, error) {
	a, err := ctrl.agent.FindByID(ctx, r.AgentID)
	if err != nil {
		return nil, err
	}

	if !a.Invocation.User.Enabled {
		return nil, fmt.Errorf("agent is not available for user invocation")
	}

	return service.DefaultAgenticRuntime.Run(ctx, &runtime.AgentRequest{
		AgentID:        r.AgentID,
		Input:          r.Input,
		ConversationID: r.ConversationID,
		ExecContext:    r.Context,
	})
}

func (ctrl *Agent) TclMasterList(_ context.Context, _ *request.AgentTclMasterList) (interface{}, error) {
	return tcl.Master(), nil
}

func (ctrl *Agent) makePayload(ctx context.Context, a *types.Agent, err error) (*agentPayload, error) {
	if err != nil || a == nil {
		return nil, err
	}

	p := &agentPayload{
		Agent: a,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateAgent: ctrl.ac.CanUpdateAgent(ctx, a),
		CanDeleteAgent: ctrl.ac.CanDeleteAgent(ctx, a),
	}

	return p, nil
}

func (ctrl *Agent) makeFilterPayload(ctx context.Context, nn types.AgentSet, f types.AgentFilter, err error) (*agentSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &agentSetPayload{Filter: f, Set: make([]*agentPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
