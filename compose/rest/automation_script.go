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
	automationScriptPayload struct {
		*automation.Script
	}

	automationScriptSetPayload struct {
		Filter automation.ScriptFilter    `json:"filter"`
		Set    []*automationScriptPayload `json:"set"`
	}

	AutomationScript struct {
		scripts automationScriptService
	}

	automationScriptService interface {
		FindByID(context.Context, uint64) (*automation.Script, error)
		Find(context.Context, automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error)
		Create(context.Context, *automation.Script) error
		Update(context.Context, *automation.Script) error
		Delete(context.Context, *automation.Script) error
	}
)

func (AutomationScript) New() *AutomationScript {
	return &AutomationScript{
		scripts: service.DefaultAutomationScriptManager,
	}
}

func (ctrl AutomationScript) List(ctx context.Context, r *request.AutomationScriptList) (interface{}, error) {
	set, filter, err := ctrl.scripts.Find(ctx, automation.ScriptFilter{
		// @todo namespace filtering
		//   Might be a bit tricky as scripts themselves not know about namespaces
		//   Namespace: r.NamespaceID

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
			RunInUA:   r.RunInUA,
			Timeout:   r.Timeout,
			Critical:  r.Critical,
			Enabled:   r.Enabled,
		}
	)

	return ctrl.makePayload(ctx, script, ctrl.scripts.Create(ctx, script))
}

func (ctrl AutomationScript) Read(ctx context.Context, r *request.AutomationScriptRead) (interface{}, error) {
	script, err := ctrl.scripts.FindByID(ctx, r.ScriptID)
	return ctrl.makePayload(ctx, script, err)
}

func (ctrl AutomationScript) Update(ctx context.Context, r *request.AutomationScriptUpdate) (interface{}, error) {
	script, err := ctrl.scripts.FindByID(ctx, r.ScriptID)
	if err != nil {
		return nil, errors.Wrap(err, "can not update script")
	}

	script.Name = r.Name
	script.SourceRef = r.SourceRef
	script.Source = r.Source
	script.Async = r.Async
	script.RunAs = r.RunAs
	script.RunInUA = r.RunInUA
	script.Timeout = r.Timeout
	script.Critical = r.Critical
	script.Enabled = r.Enabled

	return ctrl.makePayload(ctx, script, ctrl.scripts.Update(ctx, script))
}

func (ctrl AutomationScript) Delete(ctx context.Context, r *request.AutomationScriptDelete) (interface{}, error) {
	script, err := ctrl.scripts.FindByID(ctx, r.ScriptID)
	if err != nil {
		return nil, errors.Wrap(err, "can not delete script")
	}

	return resputil.OK(), ctrl.scripts.Delete(ctx, script)
}

func (ctrl AutomationScript) makePayload(ctx context.Context, s *automation.Script, err error) (*automationScriptPayload, error) {
	if err != nil || s == nil {
		return nil, err
	}

	return &automationScriptPayload{
		Script: s,

		// CanUpdateModule: ctrl.ac.CanUpdateModule(ctx, s),
		// CanDeleteModule: ctrl.ac.CanDeleteModule(ctx, s),
		// CanCreateRecord: ctrl.ac.CanCreateRecord(ctx, s),
		// CanReadRecord:   ctrl.ac.CanReadRecord(ctx, s),
		// CanUpdateRecord: ctrl.ac.CanUpdateRecord(ctx, s),
		// CanDeleteRecord: ctrl.ac.CanDeleteRecord(ctx, s),
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
