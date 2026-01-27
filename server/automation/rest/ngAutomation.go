package rest

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/automation/rest/request"
	"github.com/cortezaproject/corteza/server/automation/service"
	"github.com/cortezaproject/corteza/server/automation/types"
	cmpService "github.com/cortezaproject/corteza/server/compose/service"
	cmpTypes "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	NgAutomation struct {
		svc interface {
			Search(ctx context.Context, filter types.NgAutomationFilter) (types.NgAutomationSet, types.NgAutomationFilter, error)
			LookupByID(ctx context.Context, automationID uint64) (*types.NgAutomation, error)
			Create(ctx context.Context, new *types.NgAutomation) (*types.NgAutomation, error)
			Update(ctx context.Context, upd *types.NgAutomation) (*types.NgAutomation, error)
			DeleteByID(ctx context.Context, automationID uint64) error
			UndeleteByID(ctx context.Context, automationID uint64) error
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

	ngAutomationExecPayload struct {
		Results   *expr.Vars       `json:"results"`
		Trace     types.Stacktrace `json:"trace,omitempty"`
		SessionID uint64           `json:"sessionID,string,omitempty"`
		Error     string           `json:"error,omitempty"`
	}
)

func (NgAutomation) New() *NgAutomation {
	ctrl := &NgAutomation{}
	ctrl.svc = service.DefaultNgAutomation
	ctrl.svcModule = cmpService.DefaultModule
	ctrl.ac = service.DefaultAccessControl
	return ctrl
}

func (ctrl NgAutomation) List(ctx context.Context, r *request.NgAutomationList) (interface{}, error) {
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
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl NgAutomation) Create(ctx context.Context, r *request.NgAutomationCreate) (interface{}, error) {
	ngAutomation := &types.NgAutomation{
		Handle:   r.Handle,
		Labels:   r.Labels,
		Meta:     r.Meta,
		Enabled:  r.Enabled,
		Scope:    r.Scope,
		Triggers: r.Triggers,
		Steps:    r.Steps,
		Paths:    r.Paths,
		RunAs:    r.RunAs,
		OwnedBy:  r.OwnedBy,
	}

	wf, err := ctrl.svc.Create(ctx, ngAutomation)
	return ctrl.makePayload(ctx, wf, err)
}

func (ctrl NgAutomation) Update(ctx context.Context, r *request.NgAutomationUpdate) (interface{}, error) {
	ngAutomation := &types.NgAutomation{
		ID:        r.AutomationID,
		Handle:    r.Handle,
		Labels:    r.Labels,
		Meta:      r.Meta,
		Enabled:   r.Enabled,
		Scope:     r.Scope,
		Triggers:  r.Triggers,
		Steps:     r.Steps,
		Paths:     r.Paths,
		RunAs:     r.RunAs,
		OwnedBy:   r.OwnedBy,
		UpdatedAt: r.UpdatedAt,
	}

	wf, err := ctrl.svc.Update(ctx, ngAutomation)
	return ctrl.makePayload(ctx, wf, err)
}

func (ctrl NgAutomation) Read(ctx context.Context, r *request.NgAutomationRead) (interface{}, error) {
	wf, err := ctrl.svc.LookupByID(ctx, r.AutomationID)
	return ctrl.makePayload(ctx, wf, err)
}

func (ctrl NgAutomation) Test(ctx context.Context, r *request.NgAutomationTest) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (ctrl NgAutomation) Delete(ctx context.Context, r *request.NgAutomationDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.AutomationID)
}

func (ctrl NgAutomation) Undelete(ctx context.Context, r *request.NgAutomationUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.AutomationID)
}

func (ctrl NgAutomation) Exec(ctx context.Context, r *request.NgAutomationExec) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
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
