package rest

import (
	"context"

	"github.com/crusttech/human/server/automation/rest/request"
	"github.com/crusttech/human/server/automation/service"
	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	Trigger struct {
		trigger interface {
			Search(ctx context.Context, filter types.TriggerFilter) (types.TriggerSet, types.TriggerFilter, error)
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

func (Trigger) New() *Trigger {
	ctrl := &Trigger{}
	ctrl.trigger = service.DefaultTrigger
	return ctrl
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl Trigger) makeFilter(ctx context.Context, r *request.TriggerList) (types.TriggerFilter, error) {
	var (
		err error
		f   = types.TriggerFilter{
			WorkflowID:   r.WorkflowID,
			TriggerID:    r.TriggerID,
			EventType:    r.EventType,
			ResourceType: r.ResourceType,
			Labels:       r.Labels,
			Deleted:      filter.State(r.Deleted),
			Disabled:     filter.State(r.Disabled),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl Trigger) beforeCreate(ctx context.Context, res *types.Trigger, r *request.TriggerCreate) error {
	res.StepID = r.WorkflowStepID
	res.Constraints = r.Constraints
	res.Input = r.Input
	res.Labels = r.Labels
	res.Meta = r.Meta
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl Trigger) beforeUpdate(ctx context.Context, res *types.Trigger, r *request.TriggerUpdate) error {
	res.StepID = r.WorkflowStepID
	res.Constraints = r.Constraints
	res.Input = r.Input
	res.Labels = r.Labels
	res.Meta = r.Meta
	return nil
}

func (ctrl Trigger) makePayload(ctx context.Context, m *types.Trigger, err error) (*types.Trigger, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return m, nil
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
