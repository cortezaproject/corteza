package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/internal/service"
	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/types"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Trigger struct {
	trigger service.TriggerService
}

func (Trigger) New() *Trigger {
	return &Trigger{
		trigger: service.DefaultTrigger,
	}
}

func (ctrl *Trigger) List(ctx context.Context, r *request.TriggerList) (interface{}, error) {
	filter := types.TriggerFilter{}

	if r.ModuleID > 0 {
		filter.ModuleID = r.ModuleID
	}

	return ctrl.trigger.With(ctx).Find(filter)
}

func (ctrl *Trigger) Create(ctx context.Context, r *request.TriggerCreate) (interface{}, error) {
	trigger := &types.Trigger{
		Name:    r.Name,
		Actions: r.Actions,
		Enabled: r.Enabled,
		Source:  r.Source,
	}

	if r.ModuleID > 0 {
		trigger.ModuleID = r.ModuleID
	}

	return ctrl.trigger.With(ctx).Create(trigger)
}

func (ctrl *Trigger) Read(ctx context.Context, r *request.TriggerRead) (interface{}, error) {
	return ctrl.trigger.With(ctx).FindByID(r.TriggerID)
}

func (ctrl *Trigger) Update(ctx context.Context, r *request.TriggerUpdate) (interface{}, error) {
	if trigger, err := ctrl.trigger.FindByID(r.TriggerID); err != nil {
		return nil, err
	} else {
		trigger.Name = r.Name
		trigger.Actions = r.Actions
		trigger.Enabled = r.Enabled
		trigger.Source = r.Source
		trigger.ModuleID = r.ModuleID

		return ctrl.trigger.With(ctx).Update(trigger)
	}
}

func (ctrl *Trigger) Delete(ctx context.Context, r *request.TriggerDelete) (interface{}, error) {
	return resputil.OK(), ctrl.trigger.With(ctx).DeleteByID(r.TriggerID)
}
