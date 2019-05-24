package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	triggerPayload struct {
		*types.Trigger

		CanGrant         bool `json:"canGrant"`
		CanUpdateTrigger bool `json:"canUpdateTrigger"`
		CanDeleteTrigger bool `json:"canDeleteTrigger"`
	}

	triggerSetPayload struct {
		Filter types.TriggerFilter `json:"filter"`
		Set    []*triggerPayload   `json:"set"`
	}

	Trigger struct {
		trigger service.TriggerService
		ac      triggerAccessController
	}

	triggerAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateTrigger(context.Context, *types.Trigger) bool
		CanDeleteTrigger(context.Context, *types.Trigger) bool
	}
)

func (Trigger) New() *Trigger {
	return &Trigger{
		trigger: service.DefaultTrigger,
		ac:      service.DefaultAccessControl,
	}
}

func (ctrl Trigger) List(ctx context.Context, r *request.TriggerList) (interface{}, error) {
	f := types.TriggerFilter{
		NamespaceID: r.NamespaceID,
		Query:       r.Query,
		PerPage:     r.PerPage,
		Page:        r.Page,
	}

	set, filter, err := ctrl.trigger.With(ctx).Find(f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Trigger) Create(ctx context.Context, r *request.TriggerCreate) (interface{}, error) {
	var (
		err error
		ns  = &types.Trigger{
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			Name:        r.Name,
			Actions:     r.Actions,
			Enabled:     r.Enabled,
			Source:      r.Source,
		}
	)

	ns, err = ctrl.trigger.With(ctx).Create(ns)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Trigger) Read(ctx context.Context, r *request.TriggerRead) (interface{}, error) {
	mod, err := ctrl.trigger.With(ctx).FindByID(r.NamespaceID, r.TriggerID)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl Trigger) Update(ctx context.Context, r *request.TriggerUpdate) (interface{}, error) {
	var (
		mod = &types.Trigger{}
		err error
	)

	mod.ID = r.TriggerID
	mod.NamespaceID = r.NamespaceID
	mod.ModuleID = r.ModuleID
	mod.Name = r.Name
	mod.Actions = r.Actions
	mod.Enabled = r.Enabled
	mod.Source = r.Source

	mod, err = ctrl.trigger.With(ctx).Update(mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl Trigger) Delete(ctx context.Context, r *request.TriggerDelete) (interface{}, error) {
	_, err := ctrl.trigger.With(ctx).FindByID(r.NamespaceID, r.TriggerID)
	if err != nil {
		return nil, err
	}

	return resputil.OK(), ctrl.trigger.With(ctx).DeleteByID(r.NamespaceID, r.TriggerID)
}

func (ctrl Trigger) makePayload(ctx context.Context, t *types.Trigger, err error) (*triggerPayload, error) {
	if err != nil || t == nil {
		return nil, err
	}

	return &triggerPayload{
		Trigger: t,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateTrigger: ctrl.ac.CanUpdateTrigger(ctx, t),
		CanDeleteTrigger: ctrl.ac.CanDeleteTrigger(ctx, t),
	}, nil
}

func (ctrl Trigger) makeFilterPayload(ctx context.Context, nn types.TriggerSet, f types.TriggerFilter, err error) (*triggerSetPayload, error) {
	if err != nil {
		return nil, err
	}

	nsp := &triggerSetPayload{Filter: f, Set: make([]*triggerPayload, len(nn))}

	for i := range nn {
		nsp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return nsp, nil
}
