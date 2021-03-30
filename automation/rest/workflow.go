package rest

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/automation"
	cmpService "github.com/cortezaproject/corteza-server/compose/service"
	cmpTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/expr"
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
			Exec(ctx context.Context, workflowID uint64, p types.WorkflowExecParams) (*expr.Vars, types.Stacktrace, error)
		}

		// cross-link with compose service to load module on resolved records
		svcModule interface {
			FindByID(ctx context.Context, namespaceID, moduleID uint64) (*cmpTypes.Module, error)
		}
	}

	workflowSetPayload struct {
		Filter types.WorkflowFilter `json:"filter"`
		Set    types.WorkflowSet    `json:"set"`
	}

	workflowExecPayload struct {
		Results *expr.Vars       `json:"results"`
		Trace   types.Stacktrace `json:"trace,omitempty"`
		Error   string           `json:"error,omitempty"`
	}
)

func (Workflow) New() *Workflow {
	ctrl := &Workflow{}
	ctrl.svc = service.DefaultWorkflow
	ctrl.svcModule = cmpService.DefaultModule
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

func (ctrl Workflow) Exec(ctx context.Context, r *request.WorkflowExec) (interface{}, error) {
	var (
		wep = &workflowExecPayload{}
		err error

		execParams = types.WorkflowExecParams{
			StepID: r.StepID,
			Trace:  r.Trace,
			Input:  r.Input,
			Async:  r.Async,
			Wait:   r.Wait,
		}
	)

	if execParams.Input != nil {
		if err = execParams.Input.ResolveTypes(service.Registry().Type); err != nil {
			return nil, err
		}
	}

	// Now that all types are resolved we have to load modules and link them to records
	//
	// Very naive approach for now.
	execParams.Input.Each(func(k string, v expr.TypedValue) error {
		switch c := v.(type) {
		case *automation.ComposeRecord:
			rec := c.GetValue()
			if rec == nil {
				return nil
			}

			mod, err := ctrl.svcModule.FindByID(ctx, rec.NamespaceID, rec.ModuleID)
			if err != nil {
				return fmt.Errorf("failed to resolve ComposeRecord type: %w", err)
			}
			c.GetValue().SetModule(mod)
		}

		return nil
	})

	wep.Results, wep.Trace, err = ctrl.svc.Exec(ctx, r.WorkflowID, execParams)

	if err != nil && wep.Trace != nil && r.Trace {
		// in case of an error & trace enabled (and stacktrace present)
		// we'll suppress the error
		wep.Error = err.Error()
		return wep, nil
	}

	return wep, err
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
