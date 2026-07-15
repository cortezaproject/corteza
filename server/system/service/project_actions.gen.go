package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/project_actions.yaml

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
	projectActionProps struct {
		project *types.Project
		new     *types.Project
		update  *types.Project
		filter  *types.ProjectFilter
		diff    []*revisions.Change
		old     json.RawMessage
	}

	projectAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *projectActionProps
	}

	projectLogMetaKey   struct{}
	projectPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setProject updates projectActionProps's project
//
// This function is auto-generated.
func (p *projectActionProps) setProject(project *types.Project) *projectActionProps {
	p.project = project
	return p
}

// setNew updates projectActionProps's new
//
// This function is auto-generated.
func (p *projectActionProps) setNew(new *types.Project) *projectActionProps {
	p.new = new
	return p
}

// setUpdate updates projectActionProps's update
//
// This function is auto-generated.
func (p *projectActionProps) setUpdate(update *types.Project) *projectActionProps {
	p.update = update
	return p
}

// setFilter updates projectActionProps's filter
//
// This function is auto-generated.
func (p *projectActionProps) setFilter(filter *types.ProjectFilter) *projectActionProps {
	p.filter = filter
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *projectActionProps) setDiff(diff []*revisions.Change) *projectActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *projectActionProps) setOld(v any) *projectActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts projectActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p projectActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.project != nil {
		m.Set("project.handle", p.project.Handle, true)
		m.Set("project.ID", p.project.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p projectActionProps) Format(in string, err error) string {
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

	if p.project != nil {
		// replacement for "{{project}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{project}}",
			fns(
				p.project.Handle,
				p.project.ID,
			),
		)
		pairs = append(pairs, "{{project.handle}}", fns(p.project.Handle))
		pairs = append(pairs, "{{project.ID}}", fns(p.project.ID))
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
func (a *projectAction) String() string {
	var props = &projectActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *projectAction) ToAction() *actionlog.Action {
	resource := e.resource
	var resourceProjectID uint64
	if e.props != nil && e.props.project != nil {
		if r, ok := any(e.props.project).(actionlog.RbacResourcer); ok {
			resource = r.RbacResource()
		}

		// Attribute the event to the project owning the affected resource. This is
		// independent of the request scope, so the log stays filterable per project
		// on routes that carry no project.
		if r, ok := any(e.props.project).(actionlog.ProjectResourcer); ok {
			resourceProjectID = r.ProjectRef()
		}
	}
	return &actionlog.Action{
		Resource:          resource,
		ResourceProjectID: resourceProjectID,
		Action:            e.action,
		Severity:          e.severity,
		Description:       e.String(),
		Meta:              e.props.Serialize(),
		Delta:             actionlog.Delta(e.props.diff),
		OldState:          actionlog.OldState(e.props.old),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// ProjectActionSearch returns "system:project.search" action
//
// This function is auto-generated.
func ProjectActionSearch(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "search",
		log:       "searched for projects",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionLookup returns "system:project.lookup" action
//
// This function is auto-generated.
func ProjectActionLookup(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "lookup",
		log:       "looked-up for a {{project}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionCreate returns "system:project.create" action
//
// This function is auto-generated.
func ProjectActionCreate(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "create",
		log:       "created {{project}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionUpdate returns "system:project.update" action
//
// This function is auto-generated.
func ProjectActionUpdate(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "update",
		log:       "updated {{project}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionDelete returns "system:project.delete" action
//
// This function is auto-generated.
func ProjectActionDelete(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "delete",
		log:       "deleted {{project}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionUndelete returns "system:project.undelete" action
//
// This function is auto-generated.
func ProjectActionUndelete(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "undelete",
		log:       "undeleted {{project}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionGovernanceSave returns "system:project.governanceSave" action
//
// This function is auto-generated.
func ProjectActionGovernanceSave(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "governanceSave",
		log:       "saved governance step on {{project}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionGovernanceTransition returns "system:project.governanceTransition" action
//
// This function is auto-generated.
func ProjectActionGovernanceTransition(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "governanceTransition",
		log:       "transitioned governance step on {{project}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionSearchMembers returns "system:project.searchMembers" action
//
// This function is auto-generated.
func ProjectActionSearchMembers(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "searchMembers",
		log:       "searched members of {{project}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionAddMember returns "system:project.addMember" action
//
// This function is auto-generated.
func ProjectActionAddMember(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "addMember",
		log:       "added member to {{project}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionUpdateMember returns "system:project.updateMember" action
//
// This function is auto-generated.
func ProjectActionUpdateMember(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "updateMember",
		log:       "updated member of {{project}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectActionRemoveMember returns "system:project.removeMember" action
//
// This function is auto-generated.
func ProjectActionRemoveMember(props ...*projectActionProps) *projectAction {
	a := &projectAction{
		timestamp: time.Now(),
		resource:  "system:project",
		action:    "removeMember",
		log:       "removed member from {{project}}",
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

// ProjectErrGeneric returns "system:project.generic" as *errors.Error
//
// This function is auto-generated.
func ProjectErrGeneric(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:project"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectLogMetaKey{}, "{err}"),
		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrNotFound returns "system:project.notFound" as *errors.Error
//
// This function is auto-generated.
func ProjectErrNotFound(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrInvalidID returns "system:project.invalidID" as *errors.Error
//
// This function is auto-generated.
func ProjectErrInvalidID(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrInvalidHandle returns "system:project.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func ProjectErrInvalidHandle(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrHandleNotUnique returns "system:project.handleNotUnique" as *errors.Error
//
// This function is auto-generated.
func ProjectErrHandleNotUnique(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrInvalidStatus returns "system:project.invalidStatus" as *errors.Error
//
// This function is auto-generated.
func ProjectErrInvalidStatus(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid project status", nil),

		errors.Meta("type", "invalidStatus"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.invalidStatus"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrStaleData returns "system:project.staleData" as *errors.Error
//
// This function is auto-generated.
func ProjectErrStaleData(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrNotAllowedToCreate returns "system:project.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ProjectErrNotAllowedToCreate(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a project", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:project"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectLogMetaKey{}, "failed to create a project; insufficient permissions"),
		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrNotAllowedToRead returns "system:project.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ProjectErrNotAllowedToRead(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this project", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:project"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectLogMetaKey{}, "failed to read {{project.handle}}; insufficient permissions"),
		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrNotAllowedToSearch returns "system:project.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ProjectErrNotAllowedToSearch(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list projects", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:project"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectLogMetaKey{}, "failed to search or list projects; insufficient permissions"),
		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrNotAllowedToUpdate returns "system:project.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ProjectErrNotAllowedToUpdate(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this project", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:project"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectLogMetaKey{}, "failed to update {{project.handle}}; insufficient permissions"),
		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrNotAllowedToDelete returns "system:project.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ProjectErrNotAllowedToDelete(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this project", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:project"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectLogMetaKey{}, "failed to delete {{project.handle}}; insufficient permissions"),
		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrNotAllowedToManageMembers returns "system:project.notAllowedToManageMembers" as *errors.Error
//
// This function is auto-generated.
func ProjectErrNotAllowedToManageMembers(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage members of this project", nil),

		errors.Meta("type", "notAllowedToManageMembers"),
		errors.Meta("resource", "system:project"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectLogMetaKey{}, "failed to manage members of {{project.handle}}; insufficient permissions"),
		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.notAllowedToManageMembers"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrInvalidRolePreset returns "system:project.invalidRolePreset" as *errors.Error
//
// This function is auto-generated.
func ProjectErrInvalidRolePreset(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid project member role preset", nil),

		errors.Meta("type", "invalidRolePreset"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.invalidRolePreset"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrMemberNotFound returns "system:project.memberNotFound" as *errors.Error
//
// This function is auto-generated.
func ProjectErrMemberNotFound(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project member not found", nil),

		errors.Meta("type", "memberNotFound"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.memberNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrMemberAlreadyExists returns "system:project.memberAlreadyExists" as *errors.Error
//
// This function is auto-generated.
func ProjectErrMemberAlreadyExists(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("user is already a member of this project", nil),

		errors.Meta("type", "memberAlreadyExists"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.memberAlreadyExists"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrInvalidMode returns "system:project.invalidMode" as *errors.Error
//
// This function is auto-generated.
func ProjectErrInvalidMode(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid project build mode", nil),

		errors.Meta("type", "invalidMode"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.invalidMode"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrNotAllowedToEditGovernance returns "system:project.notAllowedToEditGovernance" as *errors.Error
//
// This function is auto-generated.
func ProjectErrNotAllowedToEditGovernance(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to edit governance state of this project", nil),

		errors.Meta("type", "notAllowedToEditGovernance"),
		errors.Meta("resource", "system:project"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectLogMetaKey{}, "failed to edit governance of {{project.handle}}; insufficient capabilities"),
		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.notAllowedToEditGovernance"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrGovernanceStepLocked returns "system:project.governanceStepLocked" as *errors.Error
//
// This function is auto-generated.
func ProjectErrGovernanceStepLocked(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("governance step is locked; it was submitted or approved", nil),

		errors.Meta("type", "governanceStepLocked"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.governanceStepLocked"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectErrInvalidGovernanceTransition returns "system:project.invalidGovernanceTransition" as *errors.Error
//
// This function is auto-generated.
func ProjectErrInvalidGovernanceTransition(mm ...*projectActionProps) *errors.Error {
	var p = &projectActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid governance step transition", nil),

		errors.Meta("type", "invalidGovernanceTransition"),
		errors.Meta("resource", "system:project"),

		errors.Meta(projectPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project.errors.invalidGovernanceTransition"),

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
func (svc project) recordAction(ctx context.Context, props *projectActionProps, actionFn func(...*projectActionProps) *projectAction, err error, diff ...any) error {
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
		a.Description = props.Format(m.AsString(projectLogMetaKey{}), err)

		if p, has := m[projectPropsMetaKey{}]; has {
			a.Meta = p.(*projectActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
