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
	automationScriptPayload struct {
		*automation.Script

		CanGrant         bool `json:"canGrant"`
		CanUpdate        bool `json:"canUpdate"`
		CanDelete        bool `json:"canDelete"`
		CanSetRunner     bool `json:"canSetRunner"`
		CanSetAsAsync    bool `json:"canSetAsAsync"`
		CanSetAsCritical bool `json:"canAsCritical"`
	}

	automationScriptSetPayload struct {
		Filter automation.ScriptFilter    `json:"filter"`
		Set    []*automationScriptPayload `json:"set"`
	}

	automationScriptRunnablePayload struct {
		Set []*automationScriptRunnable `json:"set"`
	}

	automationScriptRunnable struct {
		ScriptID uint64              `json:"scriptID,string"`
		Name     string              `json:"name"`
		Events   map[string][]string `json:"events"`
		Source   string              `json:"source,omitempty"`
		Async    bool                `json:"async"`
	}

	AutomationScript struct {
		scripts automationScriptService
		runner  automationScriptRunner
		ac      automationScriptAccessController
	}

	automationScriptService interface {
		FindByID(context.Context, uint64) (*automation.Script, error)
		Find(context.Context, automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error)
		Create(context.Context, *automation.Script) error
		Update(context.Context, *automation.Script) error
		Delete(context.Context, uint64) error
	}

	automationScriptRunner interface {
		RecordScriptTester(context.Context, string, interface{}) error
	}

	automationScriptAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateAutomationScript(context.Context, *automation.Script) bool
		CanDeleteAutomationScript(context.Context, *automation.Script) bool
	}
)

func (AutomationScript) New() *AutomationScript {
	return &AutomationScript{
		scripts: service.DefaultAutomationScriptManager,
		runner:  service.DefaultAutomationRunner,
		ac:      service.DefaultAccessControl,
	}
}

func (ctrl AutomationScript) List(ctx context.Context, r *request.AutomationScriptList) (interface{}, error) {
	set, filter, err := ctrl.scripts.Find(ctx, automation.ScriptFilter{
		Query:    r.Query,
		Resource: r.Resource,

		IncDeleted: false,
		PageFilter: rh.Paging(r.Page, r.PerPage),
	})

	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl AutomationScript) Create(ctx context.Context, r *request.AutomationScriptCreate) (interface{}, error) {
	var (
		script = &automation.Script{
			Name:      r.Name,
			SourceRef: r.SourceRef,
			Source:    r.Source,
			Async:     r.Async,
			RunAs:     r.RunAs,
			Timeout:   r.Timeout,
			Critical:  r.Critical,
			Enabled:   r.Enabled,
		}
	)

	script.AddTrigger(automation.STMS_FRESH, r.Triggers...)

	return ctrl.makePayload(ctx, script, ctrl.scripts.Create(ctx, script))
}

func (ctrl AutomationScript) Read(ctx context.Context, r *request.AutomationScriptRead) (interface{}, error) {
	script, err := ctrl.scripts.FindByID(ctx, r.ScriptID)
	return ctrl.makePayload(ctx, script, err)
}

func (ctrl AutomationScript) Update(ctx context.Context, r *request.AutomationScriptUpdate) (interface{}, error) {
	mod := &automation.Script{
		ID:        r.ScriptID,
		Name:      r.Name,
		SourceRef: r.SourceRef,
		Source:    r.Source,
		Async:     r.Async,
		RunAs:     r.RunAs,
		Timeout:   r.Timeout,
		Critical:  r.Critical,
		Enabled:   r.Enabled,
	}

	mod.AddTrigger(automation.STMS_UPDATE, r.Triggers...)

	return ctrl.makePayload(ctx, mod, ctrl.scripts.Update(ctx, mod))
}

func (ctrl AutomationScript) Delete(ctx context.Context, r *request.AutomationScriptDelete) (interface{}, error) {
	return resputil.OK(), ctrl.scripts.Delete(ctx, r.ScriptID)
}

func (ctrl AutomationScript) Test(ctx context.Context, r *request.AutomationScriptTest) (interface{}, error) {
	var (
		err error
	)

	if err = ctrl.runner.RecordScriptTester(ctx, r.Source, r.Payload); err != nil {
		return nil, err
	}

	return r.Payload, err
}

func (ctrl AutomationScript) makePayload(ctx context.Context, s *automation.Script, err error) (*automationScriptPayload, error) {
	if err != nil || s == nil {
		return nil, err
	}

	return &automationScriptPayload{
		Script: s,

		CanGrant:  ctrl.ac.CanGrant(ctx),
		CanUpdate: ctrl.ac.CanUpdateAutomationScript(ctx, s),
		CanDelete: ctrl.ac.CanDeleteAutomationScript(ctx, s),

		CanSetRunner:     ctrl.ac.CanGrant(ctx),
		CanSetAsCritical: true,
		CanSetAsAsync:    true,
	}, nil
}

func (ctrl AutomationScript) makeFilterPayload(ctx context.Context, nn automation.ScriptSet, f automation.ScriptFilter, err error) (*automationScriptSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &automationScriptSetPayload{Filter: f, Set: make([]*automationScriptPayload, len(nn))}

	for i := range nn {
		modp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return modp, nil
}
