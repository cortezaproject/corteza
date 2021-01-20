package rest

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
)

type (
	Workflow struct {
		svc interface {
			Search(ctx context.Context, filter types.WorkflowFilter) (types.WorkflowSet, types.WorkflowFilter, error)
			LookupByID(ctx context.Context, workflowID uint64) (*types.Workflow, error)
			Create(ctx context.Context, new *types.Workflow) (*types.Workflow, error)
			Update(ctx context.Context, upd *types.Workflow) (*types.Workflow, error)
			DeleteByID(ctx context.Context, workflowID uint64) error
			UndeleteByID(ctx context.Context, workflowID uint64) error
		}
	}

	workflowSetPayload struct {
		Filter types.WorkflowFilter `json:"filter"`
		Set    types.WorkflowSet    `json:"set"`
	}
)

func (Workflow) New() *Workflow {
	ctrl := &Workflow{}
	ctrl.svc = service.DefaultWorkflow
	return ctrl
}

func (ctrl Workflow) List(ctx context.Context, r *request.WorkflowList) (interface{}, error) {
	var (
		err error
		f   = types.WorkflowFilter{
			WorkflowID: payload.ParseUint64s(r.WorkflowID),
			Query:      r.Query,
			Labels:     r.Labels,
			Deleted:    filter.State(r.Deleted),
			Disabled:   filter.State(r.Disabled),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Workflow) Create(ctx context.Context, r *request.WorkflowCreate) (interface{}, error) {
	workflow := &types.Workflow{
		Handle:       r.Handle,
		Labels:       r.Labels,
		Meta:         r.Meta,
		Enabled:      r.Enabled,
		Trace:        r.Trace,
		KeepSessions: r.KeepSessions,
		Scope:        r.Scope,
		Steps:        r.Steps,
		Paths:        r.Paths,
		RunAs:        r.RunAs,
		OwnedBy:      r.OwnedBy,
	}

	return ctrl.svc.Create(ctx, workflow)
}

func (ctrl Workflow) Update(ctx context.Context, r *request.WorkflowUpdate) (interface{}, error) {
	workflow := &types.Workflow{
		ID:           r.WorkflowID,
		Handle:       r.Handle,
		Labels:       r.Labels,
		Meta:         r.Meta,
		Enabled:      r.Enabled,
		Trace:        r.Trace,
		KeepSessions: r.KeepSessions,
		Scope:        r.Scope,
		Steps:        r.Steps,
		Paths:        r.Paths,
		RunAs:        r.RunAs,
		OwnedBy:      r.OwnedBy,
	}

	return ctrl.svc.Update(ctx, workflow)
}

func (ctrl Workflow) Read(ctx context.Context, r *request.WorkflowRead) (interface{}, error) {
	return ctrl.svc.LookupByID(ctx, r.WorkflowID)
}

func (ctrl Workflow) Test(ctx context.Context, r *request.WorkflowTest) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ctrl Workflow) Delete(ctx context.Context, r *request.WorkflowDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.WorkflowID)
}

func (ctrl Workflow) Undelete(ctx context.Context, r *request.WorkflowUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.WorkflowID)
}

func (ctrl Workflow) makeFilterPayload(ctx context.Context, uu types.WorkflowSet, f types.WorkflowFilter, err error) (*workflowSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(uu) == 0 {
		uu = make([]*types.Workflow, 0)
	}

	return &workflowSetPayload{Filter: f, Set: uu}, nil
}
