package rest

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/rh"
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
		RunInUA  bool                `json:"runInUA"`
	}

	automationScriptRun struct {
		Record *types.Record `json:"record,omitempty"`
	}

	AutomationScript struct {
		scripts automationScriptService
		runner  automationScriptRunner
		ac      automationScriptAccessController

		namespace service.NamespaceService
		module    service.ModuleService
		record    service.RecordService
	}

	automationScriptService interface {
		FindByID(context.Context, uint64, uint64) (*automation.Script, error)
		Find(context.Context, uint64, automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error)
		Create(context.Context, uint64, *automation.Script) error
		Update(context.Context, uint64, *automation.Script) error
		Delete(context.Context, uint64, uint64) error
	}

	automationScriptRunner interface {
		UserScripts(context.Context) automation.ScriptSet
		RecordManual(context.Context, uint64, *types.Namespace, *types.Module, *types.Record) error
		RecordScriptTester(context.Context, string, *types.Namespace, *types.Module, *types.Record) error
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

		namespace: service.DefaultNamespace,
		module:    service.DefaultModule,
		record:    service.DefaultRecord,
	}
}

func (ctrl AutomationScript) List(ctx context.Context, r *request.AutomationScriptList) (interface{}, error) {
	set, filter, err := ctrl.scripts.Find(ctx, r.NamespaceID, automation.ScriptFilter{
		NamespaceID: r.NamespaceID,

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
			NamespaceID: r.NamespaceID,
			Name:        r.Name,
			SourceRef:   r.SourceRef,
			Source:      r.Source,
			Async:       r.Async,
			RunAs:       r.RunAs,
			RunInUA:     r.RunInUA,
			Timeout:     r.Timeout,
			Critical:    r.Critical,
			Enabled:     r.Enabled,
		}
	)

	script.AddTrigger(automation.STMS_FRESH, r.Triggers...)

	return ctrl.makePayload(ctx, script, ctrl.scripts.Create(ctx, r.NamespaceID, script))
}

func (ctrl AutomationScript) Read(ctx context.Context, r *request.AutomationScriptRead) (interface{}, error) {
	script, err := ctrl.scripts.FindByID(ctx, r.NamespaceID, r.ScriptID)
	return ctrl.makePayload(ctx, script, err)
}

func (ctrl AutomationScript) Update(ctx context.Context, r *request.AutomationScriptUpdate) (interface{}, error) {
	mod := &automation.Script{
		ID:          r.ScriptID,
		NamespaceID: r.NamespaceID,
		Name:        r.Name,
		SourceRef:   r.SourceRef,
		Source:      r.Source,
		Async:       r.Async,
		RunAs:       r.RunAs,
		RunInUA:     r.RunInUA,
		Timeout:     r.Timeout,
		Critical:    r.Critical,
		Enabled:     r.Enabled,
	}

	mod.AddTrigger(automation.STMS_UPDATE, r.Triggers...)

	return ctrl.makePayload(ctx, mod, ctrl.scripts.Update(ctx, r.NamespaceID, mod))
}

func (ctrl AutomationScript) Delete(ctx context.Context, r *request.AutomationScriptDelete) (interface{}, error) {
	return resputil.OK(), ctrl.scripts.Delete(ctx, r.NamespaceID, r.ScriptID)
}

func (ctrl AutomationScript) Runnable(ctx context.Context, r *request.AutomationScriptRunnable) (interface{}, error) {
	var (
		rval = &automationScriptRunnablePayload{
			Set: make([]*automationScriptRunnable, 0),
		}
	)

	return rval, ctrl.runner.UserScripts(ctx).Walk(func(script *automation.Script) error {
		if script.NamespaceID != r.NamespaceID {
			return nil
		}

		// @todo filter out all modules (by t.Condition) we do not have access to
		out := &automationScriptRunnable{
			ScriptID: script.ID,
			Name:     script.Name,
			Events:   map[string][]string{},
			Async:    script.Async,
			RunInUA:  script.RunInUA,
		}

		if script.RunInUA {
			out.Source = script.Source
		}

		_ = script.Triggers().Walk(func(t *automation.Trigger) error {
			if r.Condition != "" && r.Condition != t.Condition {
				// When not requesting explicit module and condition does not match (module id or 0)
				// ignore
				return nil
			}

			if _, ok := out.Events[t.Event]; ok {
				out.Events[t.Event] = append(out.Events[t.Event], t.Condition)
			} else {
				out.Events[t.Event] = []string{t.Condition}
			}

			return nil
		})

		if len(out.Events) == 0 {
			return nil
		}
		rval.Set = append(rval.Set, out)
		return nil
	})
}

func (ctrl AutomationScript) Run(ctx context.Context, r *request.AutomationScriptRun) (interface{}, error) {
	var (
		rval               automationScriptRun
		ns, m, record, err = ctrl.loadRecordScriptRunningCombo(ctx, r.NamespaceID, r.ModuleID, r.RecordID, r.Record)
	)

	if err != nil {
		return nil, err
	}

	if err = ctrl.runner.RecordManual(ctx, r.ScriptID, ns, m, record); err != nil {
		return nil, err
	}

	// When record was passed return it.
	if record != nil {
		rval.Record = record
	}

	return rval, err
}

func (ctrl AutomationScript) Test(ctx context.Context, r *request.AutomationScriptTest) (interface{}, error) {
	var (
		rval               automationScriptRun
		ns, m, record, err = ctrl.loadRecordScriptRunningCombo(ctx, r.NamespaceID, r.ModuleID, 0, r.Record)
	)

	if err != nil {
		return nil, err
	}

	if err = ctrl.runner.RecordScriptTester(ctx, r.Source, ns, m, record); err != nil {
		return nil, err
	}

	// When record was passed return it.
	if record != nil {
		rval.Record = record
	}

	return rval, err
}

func (ctrl AutomationScript) loadRecordScriptRunningCombo(ctx context.Context, namespaceID, moduleID, recordID uint64, record json.RawMessage) (ns *types.Namespace, m *types.Module, r *types.Record, err error) {
	r = &types.Record{}

	// Load requested namespace
	if ns, err = ctrl.namespace.With(ctx).FindByID(namespaceID); err != nil {
		return
	}

	if moduleID > 0 {
		// Unmarshal given module or find existing one from ID
		if m, err = ctrl.module.With(ctx).FindByID(ns.ID, moduleID); err != nil {
			return
		}
	}

	// Unmarshal given record or find existing one from ID
	if record != nil {
		if err = json.Unmarshal(record, &r); err != nil {
			err = errors.New("Could not parse record payload")
			return
		}
	} else if r, err = ctrl.record.With(ctx).FindByID(ns.ID, recordID); err != nil {
		return
	}

	return
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
