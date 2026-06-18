package rest

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/automation/rest/request"
	"github.com/crusttech/human/server/automation/service"
	"github.com/crusttech/human/server/automation/types"
	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	execTypes "github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/modern-go/reflect2"
)

type (
	NgAutomation struct {
		ngAutomation interface {
			Search(ctx context.Context, filter types.NgAutomationFilter) (types.NgAutomationSet, types.NgAutomationFilter, error)
			FindByID(ctx context.Context, automationID uint64) (*types.NgAutomation, error)
			Create(ctx context.Context, new *types.NgAutomation) (*types.NgAutomation, error)
			Update(ctx context.Context, upd *types.NgAutomation) (*types.NgAutomation, error)
			DeleteByID(ctx context.Context, automationID uint64) error
			UndeleteByID(ctx context.Context, automationID uint64) error

			ExecAndWait(ctx context.Context, automationID uint64, p types.NgAutomationExecParams) (out *execTypes.ExecutionResult, err error)
			GetExecutions(ctx context.Context, automationID uint64) ([]*execTypes.ExecutionResult, error)
			GetExecutionTrace(ctx context.Context, exeID, executionID uint64, rev int) ([]execTypes.StackFrame, error)
			GetAllExecutions(ctx context.Context, f execTypes.ExecutionFilter) ([]*execTypes.ExecutionResult, error)
		}

		// cross-link with compose service to load module on resolved records
		svcModule interface {
			FindByID(ctx context.Context, namespaceID, moduleID uint64) (*cmpTypes.Module, error)
		}

		ac ngAutomationAccessControl
	}

	ngAutomationAccessControl interface {
		CanGrant(context.Context) bool

		CanUpdateNgAutomation(context.Context, *types.NgAutomation) bool
		CanDeleteNgAutomation(context.Context, *types.NgAutomation) bool
		CanUndeleteNgAutomation(context.Context, *types.NgAutomation) bool
		CanExecuteNgAutomation(context.Context, *types.NgAutomation) bool
	}

	ngAutomationPayload struct {
		*types.NgAutomation

		CanGrant                      bool `json:"canGrant"`
		CanUpdateNgAutomation         bool `json:"canUpdateNgAutomation"`
		CanDeleteNgAutomation         bool `json:"canDeleteNgAutomation"`
		CanUndeleteNgAutomation       bool `json:"canUndeleteNgAutomation"`
		CanExecuteNgAutomation        bool `json:"canExecuteNgAutomation"`
		CanManageNgAutomationTriggers bool `json:"canManageNgAutomationTriggers"`
		CanManageNgAutomationSessions bool `json:"canManageNgAutomationSessions"`
	}

	ngAutomationSetPayload struct {
		Filter types.NgAutomationFilter `json:"filter"`
		Set    []*ngAutomationPayload   `json:"set"`
	}
)

func (NgAutomation) New() *NgAutomation {
	ctrl := &NgAutomation{}
	ctrl.ngAutomation = service.DefaultNgAutomation
	ctrl.svcModule = cmpService.DefaultModule
	ctrl.ac = service.DefaultAccessControl
	return ctrl
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl NgAutomation) makeFilter(ctx context.Context, r *request.NgAutomationList) (types.NgAutomationFilter, error) {
	var (
		err error
		f   = types.NgAutomationFilter{
			// AutomationID:    r.AutomationID,
			Query:    r.Query,
			Labels:   r.Labels,
			Deleted:  filter.State(r.Deleted),
			Disabled: filter.State(r.Disabled),
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
func (ctrl NgAutomation) beforeCreate(ctx context.Context, res *types.NgAutomation, r *request.NgAutomationCreate) error {
	res.Meta = r.Meta
	res.Scope = r.Scope
	res.Triggers = r.Triggers
	res.Steps = r.Steps
	res.Paths = r.Paths
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl NgAutomation) beforeUpdate(ctx context.Context, res *types.NgAutomation, r *request.NgAutomationUpdate) error {
	res.Meta = r.Meta
	res.Scope = r.Scope
	res.Triggers = r.Triggers
	res.Steps = r.Steps
	res.Paths = r.Paths
	return nil
}

func (ctrl NgAutomation) Test(ctx context.Context, r *request.NgAutomationTest) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ctrl NgAutomation) Exec(ctx context.Context, r *request.NgAutomationExec) (interface{}, error) {
	input := r.Input

	if !reflect2.IsNil(input) {
		err := input.ResolveTypes(service.Registry().Type)
		if err != nil {
			return nil, err
		}
	}

	return ctrl.ngAutomation.ExecAndWait(ctx, r.AutomationID, types.NgAutomationExecParams{
		Input: input,
	})
}

func (ctrl NgAutomation) Executions(ctx context.Context, r *request.NgAutomationExecutions) (interface{}, error) {
	return ctrl.ngAutomation.GetExecutions(ctx, r.AutomationID)
}

func (ctrl NgAutomation) AllExecutions(ctx context.Context, r *request.NgAutomationAllExecutions) (interface{}, error) {
	return ctrl.ngAutomation.GetAllExecutions(ctx, execTypes.ExecutionFilter{
		AutomationID: r.AutomationID,
		EventType:    r.EventType,
		ResourceType: r.ResourceType,
		Status:       r.Status,
	})
}

func (ctrl NgAutomation) ExecutionTrace(ctx context.Context, r *request.NgAutomationExecutionTrace) (interface{}, error) {
	return ctrl.ngAutomation.GetExecutionTrace(
		ctx,
		r.AutomationID,
		r.ExecutionID,
		// @todo revisions
		0,
	)
}

func (ctrl NgAutomation) makeFilterPayload(ctx context.Context, set types.NgAutomationSet, f types.NgAutomationFilter, err error) (*ngAutomationSetPayload, error) {
	if err != nil {
		return nil, err
	}

	wfsp := &ngAutomationSetPayload{Filter: f, Set: make([]*ngAutomationPayload, len(set))}

	for i, wf := range set {
		wfsp.Set[i], _ = ctrl.makePayload(ctx, wf, nil)
	}

	return wfsp, nil
}

func (ctrl NgAutomation) makePayload(ctx context.Context, wf *types.NgAutomation, err error) (*ngAutomationPayload, error) {
	if err != nil {
		return nil, err
	}

	return &ngAutomationPayload{
		NgAutomation: wf,

		CanGrant:                ctrl.ac.CanGrant(ctx),
		CanUpdateNgAutomation:   ctrl.ac.CanUpdateNgAutomation(ctx, wf),
		CanDeleteNgAutomation:   ctrl.ac.CanDeleteNgAutomation(ctx, wf),
		CanUndeleteNgAutomation: ctrl.ac.CanUndeleteNgAutomation(ctx, wf),
		CanExecuteNgAutomation:  ctrl.ac.CanExecuteNgAutomation(ctx, wf),
	}, nil
}
