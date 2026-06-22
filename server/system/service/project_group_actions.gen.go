package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/project_group_actions.yaml

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
	projectGroupActionProps struct {
		projectGroup *types.ProjectGroup
		new          *types.ProjectGroup
		update       *types.ProjectGroup
		search       *types.ProjectGroupFilter
		diff         []*revisions.Change
		old          json.RawMessage
	}

	projectGroupAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *projectGroupActionProps
	}

	projectGroupLogMetaKey   struct{}
	projectGroupPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setProjectGroup updates projectGroupActionProps's projectGroup
//
// This function is auto-generated.
func (p *projectGroupActionProps) setProjectGroup(projectGroup *types.ProjectGroup) *projectGroupActionProps {
	p.projectGroup = projectGroup
	return p
}

// setNew updates projectGroupActionProps's new
//
// This function is auto-generated.
func (p *projectGroupActionProps) setNew(new *types.ProjectGroup) *projectGroupActionProps {
	p.new = new
	return p
}

// setUpdate updates projectGroupActionProps's update
//
// This function is auto-generated.
func (p *projectGroupActionProps) setUpdate(update *types.ProjectGroup) *projectGroupActionProps {
	p.update = update
	return p
}

// setSearch updates projectGroupActionProps's search
//
// This function is auto-generated.
func (p *projectGroupActionProps) setSearch(search *types.ProjectGroupFilter) *projectGroupActionProps {
	p.search = search
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *projectGroupActionProps) setDiff(diff []*revisions.Change) *projectGroupActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *projectGroupActionProps) setOld(v any) *projectGroupActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts projectGroupActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p projectGroupActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.projectGroup != nil {
		m.Set("projectGroup.handle", p.projectGroup.Handle, true)
		m.Set("projectGroup.ID", p.projectGroup.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.search != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p projectGroupActionProps) Format(in string, err error) string {
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

	if p.projectGroup != nil {
		// replacement for "{{projectGroup}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{projectGroup}}",
			fns(
				p.projectGroup.Handle,
				p.projectGroup.ID,
			),
		)
		pairs = append(pairs, "{{projectGroup.handle}}", fns(p.projectGroup.Handle))
		pairs = append(pairs, "{{projectGroup.ID}}", fns(p.projectGroup.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Handle,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.handle}}", fns(p.new.Handle))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Handle,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.handle}}", fns(p.update.Handle))
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
func (a *projectGroupAction) String() string {
	var props = &projectGroupActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *projectGroupAction) ToAction() *actionlog.Action {
	resource := e.resource
	if e.props != nil && e.props.projectGroup != nil {
		if r, ok := any(e.props.projectGroup).(actionlog.RbacResourcer); ok {
			resource = r.RbacResource()
		}
	}
	return &actionlog.Action{
		Resource:    resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.Serialize(),
		Delta:       actionlog.Delta(e.props.diff),
		OldState:    actionlog.OldState(e.props.old),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// ProjectGroupActionSearch returns "system:project-group.search" action
//
// This function is auto-generated.
func ProjectGroupActionSearch(props ...*projectGroupActionProps) *projectGroupAction {
	a := &projectGroupAction{
		timestamp: time.Now(),
		resource:  "system:project-group",
		action:    "search",
		log:       "searched for project groups",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectGroupActionLookup returns "system:project-group.lookup" action
//
// This function is auto-generated.
func ProjectGroupActionLookup(props ...*projectGroupActionProps) *projectGroupAction {
	a := &projectGroupAction{
		timestamp: time.Now(),
		resource:  "system:project-group",
		action:    "lookup",
		log:       "looked-up for a {{projectGroup}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectGroupActionCreate returns "system:project-group.create" action
//
// This function is auto-generated.
func ProjectGroupActionCreate(props ...*projectGroupActionProps) *projectGroupAction {
	a := &projectGroupAction{
		timestamp: time.Now(),
		resource:  "system:project-group",
		action:    "create",
		log:       "created {{projectGroup}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectGroupActionUpdate returns "system:project-group.update" action
//
// This function is auto-generated.
func ProjectGroupActionUpdate(props ...*projectGroupActionProps) *projectGroupAction {
	a := &projectGroupAction{
		timestamp: time.Now(),
		resource:  "system:project-group",
		action:    "update",
		log:       "updated {{projectGroup}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectGroupActionDelete returns "system:project-group.delete" action
//
// This function is auto-generated.
func ProjectGroupActionDelete(props ...*projectGroupActionProps) *projectGroupAction {
	a := &projectGroupAction{
		timestamp: time.Now(),
		resource:  "system:project-group",
		action:    "delete",
		log:       "deleted {{projectGroup}}",
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

// ProjectGroupErrGeneric returns "system:project-group.generic" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrGeneric(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:project-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectGroupLogMetaKey{}, "{err}"),
		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrNotFound returns "system:project-group.notFound" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrNotFound(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project group not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:project-group"),

		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrInvalidID returns "system:project-group.invalidID" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrInvalidID(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:project-group"),

		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrInvalidHandle returns "system:project-group.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrInvalidHandle(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:project-group"),

		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrHandleNotUnique returns "system:project-group.handleNotUnique" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrHandleNotUnique(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "system:project-group"),

		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrStaleData returns "system:project-group.staleData" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrStaleData(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project group was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:project-group"),

		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrNotAllowedToCreate returns "system:project-group.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrNotAllowedToCreate(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a project group", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:project-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectGroupLogMetaKey{}, "failed to create a project group; insufficient permissions"),
		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrNotAllowedToRead returns "system:project-group.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrNotAllowedToRead(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this project group", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:project-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectGroupLogMetaKey{}, "failed to read {{projectGroup.handle}}; insufficient permissions"),
		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrNotAllowedToSearch returns "system:project-group.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrNotAllowedToSearch(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list project groups", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:project-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectGroupLogMetaKey{}, "failed to search or list project groups; insufficient permissions"),
		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrNotAllowedToUpdate returns "system:project-group.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrNotAllowedToUpdate(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this project group", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:project-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectGroupLogMetaKey{}, "failed to update {{projectGroup.handle}}; insufficient permissions"),
		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrNotAllowedToDelete returns "system:project-group.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrNotAllowedToDelete(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this project group", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:project-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectGroupLogMetaKey{}, "failed to delete {{projectGroup.handle}}; insufficient permissions"),
		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectGroupErrNotAllowedToManageMembers returns "system:project-group.notAllowedToManageMembers" as *errors.Error
//
// This function is auto-generated.
func ProjectGroupErrNotAllowedToManageMembers(mm ...*projectGroupActionProps) *errors.Error {
	var p = &projectGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage members of this project group", nil),

		errors.Meta("type", "notAllowedToManageMembers"),
		errors.Meta("resource", "system:project-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectGroupLogMetaKey{}, "failed to manage members of {{projectGroup.handle}}; insufficient permissions"),
		errors.Meta(projectGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-group.errors.notAllowedToManageMembers"),

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
func (svc projectGroup) recordAction(ctx context.Context, props *projectGroupActionProps, actionFn func(...*projectGroupActionProps) *projectGroupAction, err error, diff ...any) error {
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
		a.Description = props.Format(m.AsString(projectGroupLogMetaKey{}), err)

		if p, has := m[projectGroupPropsMetaKey{}]; has {
			a.Meta = p.(*projectGroupActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
