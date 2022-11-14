package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
)

type (
	Trigger struct {
		svc interface {
			Search(ctx context.Context, filter types.TriggerFilter) (types.TriggerSet, types.TriggerFilter, error)
			LookupByID(ctx context.Context, triggerID uint64) (*types.Trigger, error)
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

func (Trigger) New() *Trigger {
	ctrl := &Trigger{}
	ctrl.svc = service.DefaultTrigger
	return ctrl
}

func (ctrl Trigger) List(ctx context.Context, r *request.TriggerList) (interface{}, error) {
	var (
		err error
		f   = types.TriggerFilter{
			WorkflowID:   payload.ParseUint64s(r.WorkflowID),
			TriggerID:    payload.ParseUint64s(r.TriggerID),
			EventType:    r.EventType,
			ResourceType: r.ResourceType,
			Labels:       r.Labels,
			Deleted:      filter.State(r.Deleted),
			Disabled:     filter.State(r.Disabled),
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

func (ctrl Trigger) Create(ctx context.Context, r *request.TriggerCreate) (interface{}, error) {
	trigger := &types.Trigger{
		Enabled:      r.Enabled,
		WorkflowID:   r.WorkflowID,
		StepID:       r.WorkflowStepID,
		ResourceType: r.ResourceType,
		EventType:    r.EventType,
		Constraints:  r.Constraints,
		Input:        r.Input,
		Labels:       r.Labels,
		OwnedBy:      r.OwnedBy,
		Meta:         r.Meta,
	}

	return ctrl.svc.Create(ctx, trigger)
}

func (ctrl Trigger) Update(ctx context.Context, r *request.TriggerUpdate) (interface{}, error) {
	trigger := &types.Trigger{
		ID:           r.TriggerID,
		Enabled:      r.Enabled,
		WorkflowID:   r.WorkflowID,
		StepID:       r.WorkflowStepID,
		ResourceType: r.ResourceType,
		EventType:    r.EventType,
		Constraints:  r.Constraints,
		Input:        r.Input,
		Labels:       r.Labels,
		OwnedBy:      r.OwnedBy,
		Meta:         r.Meta,
	}

	return ctrl.svc.Update(ctx, trigger)
}

func (ctrl Trigger) Read(ctx context.Context, r *request.TriggerRead) (interface{}, error) {
	return ctrl.svc.LookupByID(ctx, r.TriggerID)
}

func (ctrl Trigger) Delete(ctx context.Context, r *request.TriggerDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.TriggerID)
}

func (ctrl Trigger) Undelete(ctx context.Context, r *request.TriggerUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.TriggerID)
}

func (ctrl Trigger) makeFilterPayload(ctx context.Context, uu types.TriggerSet, f types.TriggerFilter, err error) (*triggerSetPayload, error) {
	if err != nil {
		return nil, err
	}

	if len(uu) == 0 {
		uu = make([]*types.Trigger, 0)
	}

	return &triggerSetPayload{Filter: f, Set: uu}, nil
}
