package rest

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

var _ = errors.Wrap

type (
	automationTriggerPayload struct {
		*automation.Trigger
	}

	automationTriggerSetPayload struct {
		Filter automation.TriggerFilter    `json:"filter"`
		Set    []*automationTriggerPayload `json:"set"`
	}

	AutomationTrigger struct {
		triggers automationTriggerService
		scripts  automationScriptFinderService
	}

	automationTriggerService interface {
		FindByID(context.Context, uint64) (*automation.Trigger, error)
		Find(context.Context, automation.TriggerFilter) (automation.TriggerSet, automation.TriggerFilter, error)
		Create(context.Context, *automation.Script, *automation.Trigger) error
		Update(context.Context, *automation.Script, *automation.Trigger) error
		Delete(context.Context, *automation.Trigger) error
	}

	automationScriptFinderService interface {
		FindByID(context.Context, uint64) (*automation.Script, error)
	}
)

func (AutomationTrigger) New() *AutomationTrigger {
	return &AutomationTrigger{
		scripts:  service.DefaultAutomationScriptManager,
		triggers: service.DefaultAutomationTriggerManager,
	}
}

func (ctrl AutomationTrigger) List(ctx context.Context, r *request.AutomationTriggerList) (interface{}, error) {
	set, filter, err := ctrl.triggers.Find(ctx, automation.TriggerFilter{
		// @todo namespace filtering
		//   Might be a bit tricky as triggers themselves not know about namespaces
		//   Namespace: r.NamespaceID

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
		return nil, errors.Wrap(err, "can not create trigger")
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
	_, t, err := ctrl.loadCombo(ctx, r.ScriptID, 0)
	if err != nil {
		return nil, errors.Wrap(err, "can not read trigger")
	}

	return ctrl.makePayload(ctx, t, err)
}

func (ctrl AutomationTrigger) Update(ctx context.Context, r *request.AutomationTriggerUpdate) (interface{}, error) {
	s, t, err := ctrl.loadCombo(ctx, r.ScriptID, r.TriggerID)
	if err != nil {
		return nil, errors.Wrap(err, "can not update trigger")
	}

	t.Event = r.Event
	t.Resource = r.Resource
	t.Condition = r.Condition
	t.ScriptID = r.ScriptID
	t.Enabled = r.Enabled

	return ctrl.makePayload(ctx, t, ctrl.triggers.Update(ctx, s, t))
}

func (ctrl AutomationTrigger) Delete(ctx context.Context, r *request.AutomationTriggerDelete) (interface{}, error) {
	trigger, err := ctrl.triggers.FindByID(ctx, r.TriggerID)
	if err != nil {
		return nil, errors.Wrap(err, "can not delete trigger")
	}

	return resputil.OK(), ctrl.triggers.Delete(ctx, trigger)
}

func (ctrl AutomationTrigger) loadCombo(ctx context.Context, scriptID, triggerID uint64) (s *automation.Script, t *automation.Trigger, err error) {
	if triggerID > 0 {
		t, err = ctrl.triggers.FindByID(ctx, triggerID)
		return
	}

	if scriptID > 0 {
		s, err = ctrl.scripts.FindByID(ctx, scriptID)
		return
	}

	return
}

func (ctrl AutomationTrigger) makePayload(ctx context.Context, s *automation.Trigger, err error) (*automationTriggerPayload, error) {
	if err != nil || s == nil {
		return nil, err
	}

	return &automationTriggerPayload{
		Trigger: s,

		// CanUpdateModule: ctrl.ac.CanUpdateModule(ctx, s),
		// CanDeleteModule: ctrl.ac.CanDeleteModule(ctx, s),
		// CanCreateRecord: ctrl.ac.CanCreateRecord(ctx, s),
		// CanReadRecord:   ctrl.ac.CanReadRecord(ctx, s),
		// CanUpdateRecord: ctrl.ac.CanUpdateRecord(ctx, s),
		// CanDeleteRecord: ctrl.ac.CanDeleteRecord(ctx, s),
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
