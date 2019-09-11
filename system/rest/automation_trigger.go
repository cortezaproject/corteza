package rest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
)

var _ = errors.Wrap

type (
	automationTriggerPayload struct {
		*automation.Trigger

		CanRun bool `json:"canRun"`
	}

	automationTriggerSetPayload struct {
		Filter automation.TriggerFilter    `json:"filter"`
		Set    []*automationTriggerPayload `json:"set"`
	}

	AutomationTrigger struct {
		triggers automationTriggerService
		scripts  automationScriptFinderService

		ac automationTriggerAccessController
	}

	automationTriggerService interface {
		FindByID(context.Context, uint64) (*automation.Trigger, error)
		Find(context.Context, automation.TriggerFilter) (automation.TriggerSet, automation.TriggerFilter, error)
		Create(context.Context, *automation.Script, *automation.Trigger) error
		Update(context.Context, *automation.Script, *automation.Trigger) error
		Delete(context.Context, *automation.Script, *automation.Trigger) error
	}

	automationScriptFinderService interface {
		FindByID(context.Context, uint64) (*automation.Script, error)
	}

	automationTriggerAccessController interface {
		CanGrant(context.Context) bool

		CanRunAutomationTrigger(context.Context, *automation.Trigger) bool
	}
)

func (AutomationTrigger) New() *AutomationTrigger {
	return &AutomationTrigger{
		scripts:  service.DefaultAutomationScriptManager,
		triggers: service.DefaultAutomationTriggerManager,
		ac:       service.DefaultAccessControl,
	}
}

func (ctrl AutomationTrigger) List(ctx context.Context, r *request.AutomationTriggerList) (interface{}, error) {
	set, filter, err := ctrl.triggers.Find(ctx, automation.TriggerFilter{
		Resource: r.Resource,
		Event:    r.Event,
		ScriptID: r.ScriptID,

		IncDeleted: false,
		PageFilter: rh.Paging(r.Page, r.PerPage),
	})

	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl AutomationTrigger) Create(ctx context.Context, r *request.AutomationTriggerCreate) (interface{}, error) {
	s, _, err := ctrl.loadCombo(ctx, r.ScriptID, 0)
	if err != nil {
		return nil, err
	}

	var (
		t = &automation.Trigger{
			Event:     r.Event,
			Resource:  r.Resource,
			Condition: r.Condition,
			ScriptID:  s.ID,
			Enabled:   r.Enabled,
		}
	)

	return ctrl.makePayload(ctx, t, ctrl.triggers.Create(ctx, s, t))
}

func (ctrl AutomationTrigger) Read(ctx context.Context, r *request.AutomationTriggerRead) (interface{}, error) {
	_, t, err := ctrl.loadCombo(ctx, r.ScriptID, r.TriggerID)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, t, err)
}

func (ctrl AutomationTrigger) Update(ctx context.Context, r *request.AutomationTriggerUpdate) (interface{}, error) {
	s, t, err := ctrl.loadCombo(ctx, r.ScriptID, r.TriggerID)
	if err != nil {
		return nil, err
	}

	t.Event = r.Event
	t.Resource = r.Resource
	t.Condition = r.Condition
	t.ScriptID = r.ScriptID
	t.Enabled = r.Enabled

	return ctrl.makePayload(ctx, t, ctrl.triggers.Update(ctx, s, t))
}

func (ctrl AutomationTrigger) Delete(ctx context.Context, r *request.AutomationTriggerDelete) (interface{}, error) {
	s, t, err := ctrl.loadCombo(ctx, r.ScriptID, r.TriggerID)
	if err != nil {
		return nil, err
	}

	return resputil.OK(), ctrl.triggers.Delete(ctx, s, t)
}

func (ctrl AutomationTrigger) loadCombo(ctx context.Context, scriptID, triggerID uint64) (s *automation.Script, t *automation.Trigger, err error) {
	if triggerID > 0 {
		if t, err = ctrl.triggers.FindByID(ctx, triggerID); err != nil {
			return
		}
	}

	if scriptID > 0 {
		s, err = ctrl.scripts.FindByID(ctx, scriptID)
		return
	}

	return
}

func (ctrl AutomationTrigger) makePayload(ctx context.Context, t *automation.Trigger, err error) (*automationTriggerPayload, error) {
	if err != nil || t == nil {
		return nil, err
	}

	return &automationTriggerPayload{
		Trigger: t,

		CanRun: ctrl.ac.CanRunAutomationTrigger(ctx, t),
	}, nil
}

func (ctrl AutomationTrigger) makeFilterPayload(ctx context.Context, nn automation.TriggerSet, f automation.TriggerFilter, err error) (*automationTriggerSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &automationTriggerSetPayload{Filter: f, Set: make([]*automationTriggerPayload, len(nn))}

	for i := range nn {
		modp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return modp, nil
}
