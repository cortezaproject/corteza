package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/project_task_actions.yaml

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
	projectTaskActionProps struct {
		projectTask *types.ProjectTask
		new         *types.ProjectTask
		update      *types.ProjectTask
		filter      *types.ProjectTaskFilter
		diff        []*revisions.Change
		old         json.RawMessage
	}

	projectTaskAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *projectTaskActionProps
	}

	projectTaskLogMetaKey   struct{}
	projectTaskPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setProjectTask updates projectTaskActionProps's projectTask
//
// This function is auto-generated.
func (p *projectTaskActionProps) setProjectTask(projectTask *types.ProjectTask) *projectTaskActionProps {
	p.projectTask = projectTask
	return p
}

// setNew updates projectTaskActionProps's new
//
// This function is auto-generated.
func (p *projectTaskActionProps) setNew(new *types.ProjectTask) *projectTaskActionProps {
	p.new = new
	return p
}

// setUpdate updates projectTaskActionProps's update
//
// This function is auto-generated.
func (p *projectTaskActionProps) setUpdate(update *types.ProjectTask) *projectTaskActionProps {
	p.update = update
	return p
}

// setFilter updates projectTaskActionProps's filter
//
// This function is auto-generated.
func (p *projectTaskActionProps) setFilter(filter *types.ProjectTaskFilter) *projectTaskActionProps {
	p.filter = filter
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *projectTaskActionProps) setDiff(diff []*revisions.Change) *projectTaskActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *projectTaskActionProps) setOld(v any) *projectTaskActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts projectTaskActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p projectTaskActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.projectTask != nil {
		m.Set("projectTask.title", p.projectTask.Title, true)
		m.Set("projectTask.ID", p.projectTask.ID, true)
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
func (p projectTaskActionProps) Format(in string, err error) string {
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

	if p.projectTask != nil {
		// replacement for "{{projectTask}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{projectTask}}",
			fns(
				p.projectTask.Title,
				p.projectTask.ID,
			),
		)
		pairs = append(pairs, "{{projectTask.title}}", fns(p.projectTask.Title))
		pairs = append(pairs, "{{projectTask.ID}}", fns(p.projectTask.ID))
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
func (a *projectTaskAction) String() string {
	var props = &projectTaskActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *projectTaskAction) ToAction() *actionlog.Action {
	resource := e.resource
	var resourceProjectID uint64
	var resourceRevisionID uint64
	if e.props != nil && e.props.projectTask != nil {
		if r, ok := any(e.props.projectTask).(actionlog.RbacResourcer); ok {
			resource = r.RbacResource()
		}

		// Attribute the event to the project owning the affected resource. This is
		// independent of the request scope, so the log stays filterable per project
		// on routes that carry no project.
		if r, ok := any(e.props.projectTask).(actionlog.ProjectResourcer); ok {
			resourceProjectID = r.ProjectRef()
		}

		// Attribute the event to the revision owning (or assigned to) the affected
		// resource. Work items (RevisionResourcer) file against the chain root and
		// carry their own, independent revision assignment. Everything else already
		// denormalises the revision onto ProjectRef() -- every projects row is its
		// own revision, so whichever revision project row a resource lives under
		// already is one.
		if r, ok := any(e.props.projectTask).(actionlog.RevisionResourcer); ok {
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

// ProjectTaskActionSearch returns "system:project-task.search" action
//
// This function is auto-generated.
func ProjectTaskActionSearch(props ...*projectTaskActionProps) *projectTaskAction {
	a := &projectTaskAction{
		timestamp: time.Now(),
		resource:  "system:project-task",
		action:    "search",
		log:       "searched for project tasks",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectTaskActionLookup returns "system:project-task.lookup" action
//
// This function is auto-generated.
func ProjectTaskActionLookup(props ...*projectTaskActionProps) *projectTaskAction {
	a := &projectTaskAction{
		timestamp: time.Now(),
		resource:  "system:project-task",
		action:    "lookup",
		log:       "looked-up for a {{projectTask}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectTaskActionCreate returns "system:project-task.create" action
//
// This function is auto-generated.
func ProjectTaskActionCreate(props ...*projectTaskActionProps) *projectTaskAction {
	a := &projectTaskAction{
		timestamp: time.Now(),
		resource:  "system:project-task",
		action:    "create",
		log:       "created {{projectTask}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectTaskActionUpdate returns "system:project-task.update" action
//
// This function is auto-generated.
func ProjectTaskActionUpdate(props ...*projectTaskActionProps) *projectTaskAction {
	a := &projectTaskAction{
		timestamp: time.Now(),
		resource:  "system:project-task",
		action:    "update",
		log:       "updated {{projectTask}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectTaskActionDelete returns "system:project-task.delete" action
//
// This function is auto-generated.
func ProjectTaskActionDelete(props ...*projectTaskActionProps) *projectTaskAction {
	a := &projectTaskAction{
		timestamp: time.Now(),
		resource:  "system:project-task",
		action:    "delete",
		log:       "deleted {{projectTask}}",
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

// ProjectTaskErrGeneric returns "system:project-task.generic" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrGeneric(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:project-task"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectTaskLogMetaKey{}, "{err}"),
		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrNotFound returns "system:project-task.notFound" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrNotFound(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project task not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:project-task"),

		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrInvalidID returns "system:project-task.invalidID" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrInvalidID(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:project-task"),

		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrMissingProject returns "system:project-task.missingProject" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrMissingProject(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project is required", nil),

		errors.Meta("type", "missingProject"),
		errors.Meta("resource", "system:project-task"),

		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.missingProject"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrInvalidRevision returns "system:project-task.invalidRevision" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrInvalidRevision(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("revision does not belong to this project's revision chain", nil),

		errors.Meta("type", "invalidRevision"),
		errors.Meta("resource", "system:project-task"),

		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.invalidRevision"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrStaleData returns "system:project-task.staleData" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrStaleData(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project task was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:project-task"),

		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrNotAllowedToCreate returns "system:project-task.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrNotAllowedToCreate(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a project task", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:project-task"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectTaskLogMetaKey{}, "failed to create a project task; insufficient permissions"),
		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrNotAllowedToRead returns "system:project-task.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrNotAllowedToRead(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this project task", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:project-task"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectTaskLogMetaKey{}, "failed to read {{projectTask.title}}; insufficient permissions"),
		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrNotAllowedToSearch returns "system:project-task.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrNotAllowedToSearch(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list project tasks", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:project-task"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectTaskLogMetaKey{}, "failed to search or list project tasks; insufficient permissions"),
		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrNotAllowedToUpdate returns "system:project-task.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrNotAllowedToUpdate(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this project task", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:project-task"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectTaskLogMetaKey{}, "failed to update {{projectTask.title}}; insufficient permissions"),
		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectTaskErrNotAllowedToDelete returns "system:project-task.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ProjectTaskErrNotAllowedToDelete(mm ...*projectTaskActionProps) *errors.Error {
	var p = &projectTaskActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this project task", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:project-task"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectTaskLogMetaKey{}, "failed to delete {{projectTask.title}}; insufficient permissions"),
		errors.Meta(projectTaskPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-task.errors.notAllowedToDelete"),

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
func (svc projectTask) recordAction(ctx context.Context, props *projectTaskActionProps, actionFn func(...*projectTaskActionProps) *projectTaskAction, err error, diff ...any) error {
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
		a.Description = props.Format(m.AsString(projectTaskLogMetaKey{}), err)

		if p, has := m[projectTaskPropsMetaKey{}]; has {
			a.Meta = p.(*projectTaskActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
