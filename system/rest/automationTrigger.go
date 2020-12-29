package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	AutomationTrigger struct {
		svc interface {
			Find(ctx context.Context, filter types.TriggerFilter) (types.TriggerSet, types.TriggerFilter, error)
			FindByID(ctx context.Context, triggerID uint64) (*types.Trigger, error)
			Create(ctx context.Context, new *types.Trigger) (*types.Trigger, error)
			Update(ctx context.Context, upd *types.Trigger) (*types.Trigger, error)
			DeleteByID(ctx context.Context, triggerID uint64) error
			UndeleteByID(ctx context.Context, triggerID uint64) error
		}
	}

	triggerSetPayload struct {
		Filter types.TriggerFilter `json:"filter"`
		Set    types.TriggerSet    `json:"set"`
	}
)

func (AutomationTrigger) New() *AutomationTrigger {
	ctrl := &AutomationTrigger{}
	ctrl.svc = service.DefaultTrigger
	return ctrl
}

func (ctrl AutomationTrigger) List(ctx context.Context, r *request.AutomationTriggerList) (interface{}, error) {
	var (
		err error
		f   = types.TriggerFilter{
			WorkflowID:   payload.ParseUint64s(r.WorkflowID),
			TriggerID:    payload.ParseUint64s(r.TriggerID),
			EventType:    r.EventType,
			ResourceType: r.ResourceType,
			Labels:       r.Labels,
			Deleted:      filter.State(r.Deleted),
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

func (ctrl AutomationTrigger) Create(ctx context.Context, r *request.AutomationTriggerCreate) (interface{}, error) {
	trigger := &types.Trigger{
		Labels:  r.Labels,
		Enabled: r.Enabled,
		OwnedBy: r.OwnedBy,
	}

	return ctrl.svc.Create(ctx, trigger)
}

func (ctrl AutomationTrigger) Update(ctx context.Context, r *request.AutomationTriggerUpdate) (interface{}, error) {
	trigger := &types.Trigger{
		ID:      r.TriggerID,
		Labels:  r.Labels,
		Enabled: r.Enabled,
		OwnedBy: r.OwnedBy,
	}

	return ctrl.svc.Update(ctx, trigger)
}

func (ctrl AutomationTrigger) Read(ctx context.Context, r *request.AutomationTriggerRead) (interface{}, error) {
	return ctrl.svc.FindByID(ctx, r.TriggerID)
}

func (ctrl AutomationTrigger) Delete(ctx context.Context, r *request.AutomationTriggerDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.TriggerID)
}

func (ctrl AutomationTrigger) Undelete(ctx context.Context, r *request.AutomationTriggerUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.TriggerID)
}

func (ctrl AutomationTrigger) makeFilterPayload(ctx context.Context, uu types.TriggerSet, f types.TriggerFilter, err error) (*triggerSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(uu) == 0 {
		uu = make([]*types.Trigger, 0)
	}

	return &triggerSetPayload{Filter: f, Set: uu}, nil
}
