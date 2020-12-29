package rest

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	AutomationWorkflow struct {
		svc interface {
			Find(ctx context.Context, filter types.WorkflowFilter) (types.WorkflowSet, types.WorkflowFilter, error)
			FindByID(ctx context.Context, workflowID uint64) (*types.Workflow, error)
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

func (AutomationWorkflow) New() *AutomationWorkflow {
	ctrl := &AutomationWorkflow{}
	ctrl.svc = service.DefaultWorkflow
	return ctrl
}

func (ctrl AutomationWorkflow) List(ctx context.Context, r *request.AutomationWorkflowList) (interface{}, error) {
	var (
		err error
		f   = types.WorkflowFilter{
			WorkflowID: payload.ParseUint64s(r.WorkflowID),
			Query:      r.Query,
			Labels:     r.Labels,
			Deleted:    filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl AutomationWorkflow) Create(ctx context.Context, r *request.AutomationWorkflowCreate) (interface{}, error) {
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

func (ctrl AutomationWorkflow) Update(ctx context.Context, r *request.AutomationWorkflowUpdate) (interface{}, error) {
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

func (ctrl AutomationWorkflow) Read(ctx context.Context, r *request.AutomationWorkflowRead) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.WorkflowID)
}

func (ctrl AutomationWorkflow) Test(ctx context.Context, r *request.AutomationWorkflowTest) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ctrl AutomationWorkflow) Delete(ctx context.Context, r *request.AutomationWorkflowDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.WorkflowID)
}

func (ctrl AutomationWorkflow) Undelete(ctx context.Context, r *request.AutomationWorkflowUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.WorkflowID)
}

func (ctrl AutomationWorkflow) makeFilterPayload(ctx context.Context, uu types.WorkflowSet, f types.WorkflowFilter, err error) (*workflowSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(uu) == 0 {
		uu = make([]*types.Workflow, 0)
	}

	return &workflowSetPayload{Filter: f, Set: uu}, nil
}
