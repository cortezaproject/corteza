package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/project_fria_scenario_actions.yaml

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/system/types"
	"strings"
	"time"
)

type (
	projectFriaScenarioActionProps struct {
		projectFriaScenario *types.ProjectFriaScenario
		new                 *types.ProjectFriaScenario
		update              *types.ProjectFriaScenario
		search              *types.ProjectFriaScenarioFilter
		diff                []*revisions.Change
		old                 json.RawMessage
	}

	projectFriaScenarioAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *projectFriaScenarioActionProps
	}

	projectFriaScenarioLogMetaKey   struct{}
	projectFriaScenarioPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setProjectFriaScenario updates projectFriaScenarioActionProps's projectFriaScenario
//
// This function is auto-generated.
func (p *projectFriaScenarioActionProps) setProjectFriaScenario(projectFriaScenario *types.ProjectFriaScenario) *projectFriaScenarioActionProps {
	p.projectFriaScenario = projectFriaScenario
	return p
}

// setNew updates projectFriaScenarioActionProps's new
//
// This function is auto-generated.
func (p *projectFriaScenarioActionProps) setNew(new *types.ProjectFriaScenario) *projectFriaScenarioActionProps {
	p.new = new
	return p
}

// setUpdate updates projectFriaScenarioActionProps's update
//
// This function is auto-generated.
func (p *projectFriaScenarioActionProps) setUpdate(update *types.ProjectFriaScenario) *projectFriaScenarioActionProps {
	p.update = update
	return p
}

// setSearch updates projectFriaScenarioActionProps's search
//
// This function is auto-generated.
func (p *projectFriaScenarioActionProps) setSearch(search *types.ProjectFriaScenarioFilter) *projectFriaScenarioActionProps {
	p.search = search
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *projectFriaScenarioActionProps) setDiff(diff []*revisions.Change) *projectFriaScenarioActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *projectFriaScenarioActionProps) setOld(v any) *projectFriaScenarioActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts projectFriaScenarioActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p projectFriaScenarioActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.projectFriaScenario != nil {
		m.Set("projectFriaScenario.title", p.projectFriaScenario.Title, true)
		m.Set("projectFriaScenario.ID", p.projectFriaScenario.ID, true)
	}
	if p.new != nil {
		m.Set("new.title", p.new.Title, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.title", p.update.Title, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.search != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p projectFriaScenarioActionProps) Format(in string, err error) string {
	var (
		pairs = []string{"{{err}}"}
		// first non-empty string
		fns = func(ii ...interface{}) string {
			for _, i := range ii {
				if s := fmt.Sprintf("%v", i); len(s) > 0 {
					return s
				}
			}

			return ""
		}
	)

	if err != nil {
		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}

	if p.projectFriaScenario != nil {
		// replacement for "{{projectFriaScenario}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{projectFriaScenario}}",
			fns(
				p.projectFriaScenario.Title,
				p.projectFriaScenario.ID,
			),
		)
		pairs = append(pairs, "{{projectFriaScenario.title}}", fns(p.projectFriaScenario.Title))
		pairs = append(pairs, "{{projectFriaScenario.ID}}", fns(p.projectFriaScenario.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Title,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.title}}", fns(p.new.Title))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Title,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.title}}", fns(p.update.Title))
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.search != nil {
		// replacement for "{{search}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{search}}",
			fns(),
		)
	}
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
func (a *projectFriaScenarioAction) String() string {
	var props = &projectFriaScenarioActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *projectFriaScenarioAction) ToAction() *actionlog.Action {
	resource := e.resource
	var resourceProjectID uint64
	var resourceRevisionID uint64
	if e.props != nil && e.props.projectFriaScenario != nil {
		if r, ok := any(e.props.projectFriaScenario).(actionlog.RbacResourcer); ok {
			resource = r.RbacResource()
		}

		// Attribute the event to the project owning the affected resource. This is
		// independent of the request scope, so the log stays filterable per project
		// on routes that carry no project.
		if r, ok := any(e.props.projectFriaScenario).(actionlog.ProjectResourcer); ok {
			resourceProjectID = r.ProjectRef()
		}

		// Attribute the event to the revision owning (or assigned to) the affected
		// resource. Work items (RevisionResourcer) file against the chain root and
		// carry their own, independent revision assignment. Everything else already
		// denormalises the revision onto ProjectRef() -- every projects row is its
		// own revision, so whichever revision project row a resource lives under
		// already is one.
		if r, ok := any(e.props.projectFriaScenario).(actionlog.RevisionResourcer); ok {
			resourceRevisionID = r.RevisionRef()
		} else {
			resourceRevisionID = resourceProjectID
		}
	}
	return &actionlog.Action{
		Resource:           resource,
		ResourceProjectID:  resourceProjectID,
		ResourceRevisionID: resourceRevisionID,
		Action:             e.action,
		Severity:           e.severity,
		Description:        e.String(),
		Meta:               e.props.Serialize(),
		Delta:              actionlog.Delta(e.props.diff),
		OldState:           actionlog.OldState(e.props.old),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// ProjectFriaScenarioActionSearch returns "system:project-fria-scenario.search" action
//
// This function is auto-generated.
func ProjectFriaScenarioActionSearch(props ...*projectFriaScenarioActionProps) *projectFriaScenarioAction {
	a := &projectFriaScenarioAction{
		timestamp: time.Now(),
		resource:  "system:project-fria-scenario",
		action:    "search",
		log:       "searched for FRIA risk scenarios",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectFriaScenarioActionLookup returns "system:project-fria-scenario.lookup" action
//
// This function is auto-generated.
func ProjectFriaScenarioActionLookup(props ...*projectFriaScenarioActionProps) *projectFriaScenarioAction {
	a := &projectFriaScenarioAction{
		timestamp: time.Now(),
		resource:  "system:project-fria-scenario",
		action:    "lookup",
		log:       "looked-up for a {{projectFriaScenario}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectFriaScenarioActionCreate returns "system:project-fria-scenario.create" action
//
// This function is auto-generated.
func ProjectFriaScenarioActionCreate(props ...*projectFriaScenarioActionProps) *projectFriaScenarioAction {
	a := &projectFriaScenarioAction{
		timestamp: time.Now(),
		resource:  "system:project-fria-scenario",
		action:    "create",
		log:       "created {{projectFriaScenario}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectFriaScenarioActionUpdate returns "system:project-fria-scenario.update" action
//
// This function is auto-generated.
func ProjectFriaScenarioActionUpdate(props ...*projectFriaScenarioActionProps) *projectFriaScenarioAction {
	a := &projectFriaScenarioAction{
		timestamp: time.Now(),
		resource:  "system:project-fria-scenario",
		action:    "update",
		log:       "updated {{projectFriaScenario}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectFriaScenarioActionDelete returns "system:project-fria-scenario.delete" action
//
// This function is auto-generated.
func ProjectFriaScenarioActionDelete(props ...*projectFriaScenarioActionProps) *projectFriaScenarioAction {
	a := &projectFriaScenarioAction{
		timestamp: time.Now(),
		resource:  "system:project-fria-scenario",
		action:    "delete",
		log:       "deleted {{projectFriaScenario}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// ProjectFriaScenarioErrGeneric returns "system:project-fria-scenario.generic" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrGeneric(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:project-fria-scenario"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFriaScenarioLogMetaKey{}, "{err}"),
		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFriaScenarioErrNotFound returns "system:project-fria-scenario.notFound" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrNotFound(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("FRIA risk scenario not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:project-fria-scenario"),

		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFriaScenarioErrInvalidID returns "system:project-fria-scenario.invalidID" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrInvalidID(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:project-fria-scenario"),

		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFriaScenarioErrMissingProject returns "system:project-fria-scenario.missingProject" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrMissingProject(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project is required", nil),

		errors.Meta("type", "missingProject"),
		errors.Meta("resource", "system:project-fria-scenario"),

		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.missingProject"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFriaScenarioErrStaleData returns "system:project-fria-scenario.staleData" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrStaleData(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("FRIA risk scenario was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:project-fria-scenario"),

		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFriaScenarioErrNotAllowedToCreate returns "system:project-fria-scenario.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrNotAllowedToCreate(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a FRIA risk scenario", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:project-fria-scenario"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFriaScenarioLogMetaKey{}, "failed to create a FRIA risk scenario; insufficient permissions"),
		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFriaScenarioErrNotAllowedToRead returns "system:project-fria-scenario.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrNotAllowedToRead(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this FRIA risk scenario", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:project-fria-scenario"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFriaScenarioLogMetaKey{}, "failed to read {{projectFriaScenario.title}}; insufficient permissions"),
		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFriaScenarioErrNotAllowedToSearch returns "system:project-fria-scenario.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrNotAllowedToSearch(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list FRIA risk scenarios", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:project-fria-scenario"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFriaScenarioLogMetaKey{}, "failed to search or list FRIA risk scenarios; insufficient permissions"),
		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFriaScenarioErrNotAllowedToUpdate returns "system:project-fria-scenario.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrNotAllowedToUpdate(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this FRIA risk scenario", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:project-fria-scenario"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFriaScenarioLogMetaKey{}, "failed to update {{projectFriaScenario.title}}; insufficient permissions"),
		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFriaScenarioErrNotAllowedToDelete returns "system:project-fria-scenario.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ProjectFriaScenarioErrNotAllowedToDelete(mm ...*projectFriaScenarioActionProps) *errors.Error {
	var p = &projectFriaScenarioActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this FRIA risk scenario", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:project-fria-scenario"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFriaScenarioLogMetaKey{}, "failed to delete {{projectFriaScenario.title}}; insufficient permissions"),
		errors.Meta(projectFriaScenarioPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-fria-scenario.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// It will wrap unrecognized/internal errors with generic errors.
//
// This function is auto-generated.
func (svc projectFriaScenario) recordAction(ctx context.Context, props *projectFriaScenarioActionProps, actionFn func(...*projectFriaScenarioActionProps) *projectFriaScenarioAction, err error, diff ...any) error {
	if len(diff) == 2 && diff[0] != nil {
		props.setDiff(actionlog.DiffResourceState(diff[0], diff[1]))
		props.setOld(diff[0])
	}
	if svc.actionlog == nil || actionFn == nil {
		// action log disabled or no action fn passed, return error as-is
		return err
	} else if err == nil {
		// action completed w/o error, record it
		svc.actionlog.Record(ctx, actionFn(props).ToAction())
		return nil
	}

	a := actionFn(props).ToAction()

	// Extracting error information and recording it as action
	a.Error = err.Error()

	switch c := err.(type) {
	case *errors.Error:
		m := c.Meta()

		a.Error = err.Error()
		a.Severity = actionlog.Severity(m.AsInt("severity"))
		a.Description = props.Format(m.AsString(projectFriaScenarioLogMetaKey{}), err)

		if p, has := m[projectFriaScenarioPropsMetaKey{}]; has {
			a.Meta = p.(*projectFriaScenarioActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
