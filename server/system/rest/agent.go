package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/system/runtime"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	Agent struct {
		svc agentService
		ac  agentAccessController
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
		svc: service.DefaultAgent,
		ac:  service.DefaultAccessControl,
	}
}

func (ctrl *Agent) List(ctx context.Context, r *request.AgentList) (interface{}, error) {
	var (
		err error
		f   = types.AgentFilter{
			Query:   r.Query,
			Handle:  r.Handle,
			Status:  r.Status,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Agent) Create(ctx context.Context, r *request.AgentCreate) (interface{}, error) {
	var (
		err error
		a   = &types.Agent{
			Handle:     r.Handle,
			Status:     r.Status,
			Meta:       r.Meta,
			Behavior:   r.Behavior,
			Execution:  r.Execution,
			Access:     r.Access,
			Invocation: r.Invocation,
		}
	)

	a, err = ctrl.svc.Create(ctx, a)
	return ctrl.makePayload(ctx, a, err)
}

func (ctrl *Agent) Read(ctx context.Context, r *request.AgentRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.AgentID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Agent) Update(ctx context.Context, r *request.AgentUpdate) (interface{}, error) {
	var (
		err error
		a   = &types.Agent{
			ID:         r.AgentID,
			Handle:     r.Handle,
			Status:     r.Status,
			Meta:       r.Meta,
			Behavior:   r.Behavior,
			Execution:  r.Execution,
			Access:     r.Access,
			Invocation: r.Invocation,
			UpdatedAt:  r.UpdatedAt,
		}
	)

	a, err = ctrl.svc.Update(ctx, a)
	return ctrl.makePayload(ctx, a, err)
}

func (ctrl *Agent) Delete(ctx context.Context, r *request.AgentDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.AgentID)
}

func (ctrl *Agent) Undelete(ctx context.Context, r *request.AgentUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.AgentID)
}

func (ctrl *Agent) Exec(ctx context.Context, r *request.AgentExec) (interface{}, error) {
	return service.DefaultAgenticRuntime.Run(ctx, &runtime.AgentRequest{
		AgentID:        r.AgentID,
		Input:          r.Input,
		ConversationID: r.ConversationID,
	})
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
