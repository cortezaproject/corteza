package rest

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/automation/rest/request"
	"github.com/cortezaproject/corteza/server/automation/service"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/automation"
	cmpService "github.com/cortezaproject/corteza/server/compose/service"
	cmpTypes "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
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
			Exec(ctx context.Context, workflowID uint64, p types.WorkflowExecParams) (*expr.Vars, uint64, types.Stacktrace, error)
		}

		// cross-link with compose service to load module on resolved records
		svcModule interface {
			FindByID(ctx context.Context, namespaceID, moduleID uint64) (*cmpTypes.Module, error)
		}

		ac workflowAccessControl
	}

	workflowAccessControl interface {
		CanGrant(context.Context) bool

		CanUpdateWorkflow(context.Context, *types.Workflow) bool
		CanDeleteWorkflow(context.Context, *types.Workflow) bool
		CanUndeleteWorkflow(context.Context, *types.Workflow) bool
		CanExecuteWorkflow(context.Context, *types.Workflow) bool
		CanManageTriggersOnWorkflow(context.Context, *types.Workflow) bool
		CanManageSessionsOnWorkflow(context.Context, *types.Workflow) bool
	}

	workflowPayload struct {
		*types.Workflow

		CanGrant                  bool `json:"canGrant"`
		CanUpdateWorkflow         bool `json:"canUpdateWorkflow"`
		CanDeleteWorkflow         bool `json:"canDeleteWorkflow"`
		CanUndeleteWorkflow       bool `json:"canUndeleteWorkflow"`
		CanExecuteWorkflow        bool `json:"canExecuteWorkflow"`
		CanManageWorkflowTriggers bool `json:"canManageWorkflowTriggers"`
		CanManageWorkflowSessions bool `json:"canManageWorkflowSessions"`
	}

	workflowSetPayload struct {
		Filter types.WorkflowFilter `json:"filter"`
		Set    []*workflowPayload   `json:"set"`
	}

	workflowExecPayload struct {
		Results   *expr.Vars       `json:"results"`
		Trace     types.Stacktrace `json:"trace,omitempty"`
		SessionID uint64           `json:"sessionID,string,omitempty"`
		Error     string           `json:"error,omitempty"`
	}
)

func (Workflow) New() *Workflow {
	ctrl := &Workflow{}
	ctrl.svc = service.DefaultWorkflow
	ctrl.svcModule = cmpService.DefaultModule
	ctrl.ac = service.DefaultAccessControl
	return ctrl
}

func (ctrl Workflow) List(ctx context.Context, r *request.WorkflowList) (interface{}, error) {
	var (
		err error
		f   = types.WorkflowFilter{
			WorkflowID:  r.WorkflowID,
			Query:       r.Query,
			Labels:      r.Labels,
			Deleted:     filter.State(r.Deleted),
			Disabled:    filter.State(r.Disabled),
			SubWorkflow: filter.State(r.SubWorkflow),
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

	wf, err := ctrl.svc.Create(ctx, workflow)
	return ctrl.makePayload(ctx, wf, err)
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

	wf, err := ctrl.svc.Update(ctx, workflow)
	return ctrl.makePayload(ctx, wf, err)
}

func (ctrl Workflow) Read(ctx context.Context, r *request.WorkflowRead) (interface{}, error) {
	wf, err := ctrl.svc.LookupByID(ctx, r.WorkflowID)
	return ctrl.makePayload(ctx, wf, err)
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

	// Now when all types are resolved we have to load modules and link them to records
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

	wep.Results, wep.SessionID, wep.Trace, err = ctrl.svc.Exec(ctx, r.WorkflowID, execParams)

	if err != nil && wep.Trace != nil && r.Trace {
		// in case of an error & trace enabled (and stacktrace present)
		// we'll suppress the error
		wep.Error = err.Error()
		return wep, nil
	}

	return wep, err
}

func (ctrl Workflow) makeFilterPayload(ctx context.Context, set types.WorkflowSet, f types.WorkflowFilter, err error) (*workflowSetPayload, error) {
	if err != nil {
		return nil, err
	}

	wfsp := &workflowSetPayload{Filter: f, Set: make([]*workflowPayload, len(set))}

	for i, wf := range set {
		wfsp.Set[i], _ = ctrl.makePayload(ctx, wf, nil)
	}

	return wfsp, nil
}

func (ctrl Workflow) makePayload(ctx context.Context, wf *types.Workflow, err error) (*workflowPayload, error) {
	if err != nil {
		return nil, err
	}

	return &workflowPayload{
		Workflow: wf,

		CanGrant:                  ctrl.ac.CanGrant(ctx),
		CanUpdateWorkflow:         ctrl.ac.CanUpdateWorkflow(ctx, wf),
		CanDeleteWorkflow:         ctrl.ac.CanDeleteWorkflow(ctx, wf),
		CanUndeleteWorkflow:       ctrl.ac.CanUndeleteWorkflow(ctx, wf),
		CanExecuteWorkflow:        ctrl.ac.CanExecuteWorkflow(ctx, wf),
		CanManageWorkflowTriggers: ctrl.ac.CanManageTriggersOnWorkflow(ctx, wf),
		CanManageWorkflowSessions: ctrl.ac.CanManageSessionsOnWorkflow(ctx, wf),
	}, nil
}
