package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/project_incident_actions.yaml

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
	projectIncidentActionProps struct {
		projectIncident *types.ProjectIncident
		new             *types.ProjectIncident
		update          *types.ProjectIncident
		filter          *types.ProjectIncidentFilter
		diff            []*revisions.Change
		old             json.RawMessage
	}

	projectIncidentAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *projectIncidentActionProps
	}

	projectIncidentLogMetaKey   struct{}
	projectIncidentPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setProjectIncident updates projectIncidentActionProps's projectIncident
//
// This function is auto-generated.
func (p *projectIncidentActionProps) setProjectIncident(projectIncident *types.ProjectIncident) *projectIncidentActionProps {
	p.projectIncident = projectIncident
	return p
}

// setNew updates projectIncidentActionProps's new
//
// This function is auto-generated.
func (p *projectIncidentActionProps) setNew(new *types.ProjectIncident) *projectIncidentActionProps {
	p.new = new
	return p
}

// setUpdate updates projectIncidentActionProps's update
//
// This function is auto-generated.
func (p *projectIncidentActionProps) setUpdate(update *types.ProjectIncident) *projectIncidentActionProps {
	p.update = update
	return p
}

// setFilter updates projectIncidentActionProps's filter
//
// This function is auto-generated.
func (p *projectIncidentActionProps) setFilter(filter *types.ProjectIncidentFilter) *projectIncidentActionProps {
	p.filter = filter
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *projectIncidentActionProps) setDiff(diff []*revisions.Change) *projectIncidentActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *projectIncidentActionProps) setOld(v any) *projectIncidentActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts projectIncidentActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p projectIncidentActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.projectIncident != nil {
		m.Set("projectIncident.title", p.projectIncident.Title, true)
		m.Set("projectIncident.ID", p.projectIncident.ID, true)
	}
	if p.new != nil {
		m.Set("new.title", p.new.Title, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.title", p.update.Title, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p projectIncidentActionProps) Format(in string, err error) string {
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

	if p.projectIncident != nil {
		// replacement for "{{projectIncident}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{projectIncident}}",
			fns(
				p.projectIncident.Title,
				p.projectIncident.ID,
			),
		)
		pairs = append(pairs, "{{projectIncident.title}}", fns(p.projectIncident.Title))
		pairs = append(pairs, "{{projectIncident.ID}}", fns(p.projectIncident.ID))
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

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
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
func (a *projectIncidentAction) String() string {
	var props = &projectIncidentActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *projectIncidentAction) ToAction() *actionlog.Action {
	resource := e.resource
	var resourceProjectID uint64
	var resourceRevisionID uint64
	if e.props != nil && e.props.projectIncident != nil {
		if r, ok := any(e.props.projectIncident).(actionlog.RbacResourcer); ok {
			resource = r.RbacResource()
		}

		// Attribute the event to the project owning the affected resource. This is
		// independent of the request scope, so the log stays filterable per project
		// on routes that carry no project.
		if r, ok := any(e.props.projectIncident).(actionlog.ProjectResourcer); ok {
			resourceProjectID = r.ProjectRef()
		}

		// Attribute the event to the revision owning (or assigned to) the affected
		// resource. Work items (RevisionResourcer) file against the chain root and
		// carry their own, independent revision assignment. Everything else already
		// denormalises the revision onto ProjectRef() -- every projects row is its
		// own revision, so whichever revision project row a resource lives under
		// already is one.
		if r, ok := any(e.props.projectIncident).(actionlog.RevisionResourcer); ok {
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

// ProjectIncidentActionSearch returns "system:project-incident.search" action
//
// This function is auto-generated.
func ProjectIncidentActionSearch(props ...*projectIncidentActionProps) *projectIncidentAction {
	a := &projectIncidentAction{
		timestamp: time.Now(),
		resource:  "system:project-incident",
		action:    "search",
		log:       "searched for project incidents",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectIncidentActionLookup returns "system:project-incident.lookup" action
//
// This function is auto-generated.
func ProjectIncidentActionLookup(props ...*projectIncidentActionProps) *projectIncidentAction {
	a := &projectIncidentAction{
		timestamp: time.Now(),
		resource:  "system:project-incident",
		action:    "lookup",
		log:       "looked-up for a {{projectIncident}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectIncidentActionCreate returns "system:project-incident.create" action
//
// This function is auto-generated.
func ProjectIncidentActionCreate(props ...*projectIncidentActionProps) *projectIncidentAction {
	a := &projectIncidentAction{
		timestamp: time.Now(),
		resource:  "system:project-incident",
		action:    "create",
		log:       "created {{projectIncident}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectIncidentActionUpdate returns "system:project-incident.update" action
//
// This function is auto-generated.
func ProjectIncidentActionUpdate(props ...*projectIncidentActionProps) *projectIncidentAction {
	a := &projectIncidentAction{
		timestamp: time.Now(),
		resource:  "system:project-incident",
		action:    "update",
		log:       "updated {{projectIncident}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectIncidentActionDelete returns "system:project-incident.delete" action
//
// This function is auto-generated.
func ProjectIncidentActionDelete(props ...*projectIncidentActionProps) *projectIncidentAction {
	a := &projectIncidentAction{
		timestamp: time.Now(),
		resource:  "system:project-incident",
		action:    "delete",
		log:       "deleted {{projectIncident}}",
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

// ProjectIncidentErrGeneric returns "system:project-incident.generic" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrGeneric(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:project-incident"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectIncidentLogMetaKey{}, "{err}"),
		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrNotFound returns "system:project-incident.notFound" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrNotFound(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project incident not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:project-incident"),

		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrInvalidID returns "system:project-incident.invalidID" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrInvalidID(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:project-incident"),

		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrMissingProject returns "system:project-incident.missingProject" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrMissingProject(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project is required", nil),

		errors.Meta("type", "missingProject"),
		errors.Meta("resource", "system:project-incident"),

		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.missingProject"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrInvalidRevision returns "system:project-incident.invalidRevision" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrInvalidRevision(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("revision does not belong to this project's revision chain", nil),

		errors.Meta("type", "invalidRevision"),
		errors.Meta("resource", "system:project-incident"),

		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.invalidRevision"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrStaleData returns "system:project-incident.staleData" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrStaleData(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project incident was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:project-incident"),

		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrNotAllowedToCreate returns "system:project-incident.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrNotAllowedToCreate(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a project incident", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:project-incident"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectIncidentLogMetaKey{}, "failed to create a project incident; insufficient permissions"),
		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrNotAllowedToRead returns "system:project-incident.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrNotAllowedToRead(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this project incident", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:project-incident"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectIncidentLogMetaKey{}, "failed to read {{projectIncident.title}}; insufficient permissions"),
		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrNotAllowedToSearch returns "system:project-incident.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrNotAllowedToSearch(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list project incidents", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:project-incident"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectIncidentLogMetaKey{}, "failed to search or list project incidents; insufficient permissions"),
		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrNotAllowedToUpdate returns "system:project-incident.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrNotAllowedToUpdate(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this project incident", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:project-incident"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectIncidentLogMetaKey{}, "failed to update {{projectIncident.title}}; insufficient permissions"),
		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectIncidentErrNotAllowedToDelete returns "system:project-incident.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ProjectIncidentErrNotAllowedToDelete(mm ...*projectIncidentActionProps) *errors.Error {
	var p = &projectIncidentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this project incident", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:project-incident"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectIncidentLogMetaKey{}, "failed to delete {{projectIncident.title}}; insufficient permissions"),
		errors.Meta(projectIncidentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-incident.errors.notAllowedToDelete"),

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
func (svc projectIncident) recordAction(ctx context.Context, props *projectIncidentActionProps, actionFn func(...*projectIncidentActionProps) *projectIncidentAction, err error, diff ...any) error {
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
		a.Description = props.Format(m.AsString(projectIncidentLogMetaKey{}), err)

		if p, has := m[projectIncidentPropsMetaKey{}]; has {
			a.Meta = p.(*projectIncidentActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
