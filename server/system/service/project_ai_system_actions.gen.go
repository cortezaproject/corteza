package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/project_ai_system_actions.yaml

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
	projectAiSystemActionProps struct {
		projectAiSystem *types.ProjectAiSystem
		new             *types.ProjectAiSystem
		update          *types.ProjectAiSystem
		search          *types.ProjectAiSystemFilter
		diff            []*revisions.Change
		old             json.RawMessage
	}

	projectAiSystemAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *projectAiSystemActionProps
	}

	projectAiSystemLogMetaKey   struct{}
	projectAiSystemPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setProjectAiSystem updates projectAiSystemActionProps's projectAiSystem
//
// This function is auto-generated.
func (p *projectAiSystemActionProps) setProjectAiSystem(projectAiSystem *types.ProjectAiSystem) *projectAiSystemActionProps {
	p.projectAiSystem = projectAiSystem
	return p
}

// setNew updates projectAiSystemActionProps's new
//
// This function is auto-generated.
func (p *projectAiSystemActionProps) setNew(new *types.ProjectAiSystem) *projectAiSystemActionProps {
	p.new = new
	return p
}

// setUpdate updates projectAiSystemActionProps's update
//
// This function is auto-generated.
func (p *projectAiSystemActionProps) setUpdate(update *types.ProjectAiSystem) *projectAiSystemActionProps {
	p.update = update
	return p
}

// setSearch updates projectAiSystemActionProps's search
//
// This function is auto-generated.
func (p *projectAiSystemActionProps) setSearch(search *types.ProjectAiSystemFilter) *projectAiSystemActionProps {
	p.search = search
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *projectAiSystemActionProps) setDiff(diff []*revisions.Change) *projectAiSystemActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *projectAiSystemActionProps) setOld(v any) *projectAiSystemActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts projectAiSystemActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p projectAiSystemActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.projectAiSystem != nil {
		m.Set("projectAiSystem.handle", p.projectAiSystem.Handle, true)
		m.Set("projectAiSystem.ID", p.projectAiSystem.ID, true)
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
func (p projectAiSystemActionProps) Format(in string, err error) string {
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

	if p.projectAiSystem != nil {
		// replacement for "{{projectAiSystem}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{projectAiSystem}}",
			fns(
				p.projectAiSystem.Handle,
				p.projectAiSystem.ID,
			),
		)
		pairs = append(pairs, "{{projectAiSystem.handle}}", fns(p.projectAiSystem.Handle))
		pairs = append(pairs, "{{projectAiSystem.ID}}", fns(p.projectAiSystem.ID))
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
func (a *projectAiSystemAction) String() string {
	var props = &projectAiSystemActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *projectAiSystemAction) ToAction() *actionlog.Action {
	resource := e.resource
	var resourceProjectID uint64
	var resourceRevisionID uint64
	if e.props != nil && e.props.projectAiSystem != nil {
		if r, ok := any(e.props.projectAiSystem).(actionlog.RbacResourcer); ok {
			resource = r.RbacResource()
		}

		// Attribute the event to the project owning the affected resource. This is
		// independent of the request scope, so the log stays filterable per project
		// on routes that carry no project.
		if r, ok := any(e.props.projectAiSystem).(actionlog.ProjectResourcer); ok {
			resourceProjectID = r.ProjectRef()
		}

		// Attribute the event to the revision owning (or assigned to) the affected
		// resource. Work items (RevisionResourcer) file against the chain root and
		// carry their own, independent revision assignment. Everything else already
		// denormalises the revision onto ProjectRef() -- every projects row is its
		// own revision, so whichever revision project row a resource lives under
		// already is one.
		if r, ok := any(e.props.projectAiSystem).(actionlog.RevisionResourcer); ok {
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

// ProjectAiSystemActionSearch returns "system:project-ai-system.search" action
//
// This function is auto-generated.
func ProjectAiSystemActionSearch(props ...*projectAiSystemActionProps) *projectAiSystemAction {
	a := &projectAiSystemAction{
		timestamp: time.Now(),
		resource:  "system:project-ai-system",
		action:    "search",
		log:       "searched for AI systems",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectAiSystemActionLookup returns "system:project-ai-system.lookup" action
//
// This function is auto-generated.
func ProjectAiSystemActionLookup(props ...*projectAiSystemActionProps) *projectAiSystemAction {
	a := &projectAiSystemAction{
		timestamp: time.Now(),
		resource:  "system:project-ai-system",
		action:    "lookup",
		log:       "looked-up for a {{projectAiSystem}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectAiSystemActionCreate returns "system:project-ai-system.create" action
//
// This function is auto-generated.
func ProjectAiSystemActionCreate(props ...*projectAiSystemActionProps) *projectAiSystemAction {
	a := &projectAiSystemAction{
		timestamp: time.Now(),
		resource:  "system:project-ai-system",
		action:    "create",
		log:       "created {{projectAiSystem}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectAiSystemActionUpdate returns "system:project-ai-system.update" action
//
// This function is auto-generated.
func ProjectAiSystemActionUpdate(props ...*projectAiSystemActionProps) *projectAiSystemAction {
	a := &projectAiSystemAction{
		timestamp: time.Now(),
		resource:  "system:project-ai-system",
		action:    "update",
		log:       "updated {{projectAiSystem}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectAiSystemActionDelete returns "system:project-ai-system.delete" action
//
// This function is auto-generated.
func ProjectAiSystemActionDelete(props ...*projectAiSystemActionProps) *projectAiSystemAction {
	a := &projectAiSystemAction{
		timestamp: time.Now(),
		resource:  "system:project-ai-system",
		action:    "delete",
		log:       "deleted {{projectAiSystem}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectAiSystemActionMemberList returns "system:project-ai-system.memberList" action
//
// This function is auto-generated.
func ProjectAiSystemActionMemberList(props ...*projectAiSystemActionProps) *projectAiSystemAction {
	a := &projectAiSystemAction{
		timestamp: time.Now(),
		resource:  "system:project-ai-system",
		action:    "memberList",
		log:       "listed resources of {{projectAiSystem}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectAiSystemActionMemberAdd returns "system:project-ai-system.memberAdd" action
//
// This function is auto-generated.
func ProjectAiSystemActionMemberAdd(props ...*projectAiSystemActionProps) *projectAiSystemAction {
	a := &projectAiSystemAction{
		timestamp: time.Now(),
		resource:  "system:project-ai-system",
		action:    "memberAdd",
		log:       "added resource to {{projectAiSystem}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectAiSystemActionMemberRemove returns "system:project-ai-system.memberRemove" action
//
// This function is auto-generated.
func ProjectAiSystemActionMemberRemove(props ...*projectAiSystemActionProps) *projectAiSystemAction {
	a := &projectAiSystemAction{
		timestamp: time.Now(),
		resource:  "system:project-ai-system",
		action:    "memberRemove",
		log:       "removed resource from {{projectAiSystem}}",
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

// ProjectAiSystemErrGeneric returns "system:project-ai-system.generic" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrGeneric(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:project-ai-system"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectAiSystemLogMetaKey{}, "{err}"),
		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrNotFound returns "system:project-ai-system.notFound" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrNotFound(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("AI system not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:project-ai-system"),

		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrInvalidID returns "system:project-ai-system.invalidID" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrInvalidID(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:project-ai-system"),

		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrInvalidHandle returns "system:project-ai-system.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrInvalidHandle(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:project-ai-system"),

		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrHandleNotUnique returns "system:project-ai-system.handleNotUnique" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrHandleNotUnique(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "system:project-ai-system"),

		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrStaleData returns "system:project-ai-system.staleData" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrStaleData(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("AI system was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:project-ai-system"),

		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrNotAllowedToCreate returns "system:project-ai-system.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrNotAllowedToCreate(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create an AI system", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:project-ai-system"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectAiSystemLogMetaKey{}, "failed to create an AI system; insufficient permissions"),
		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrNotAllowedToRead returns "system:project-ai-system.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrNotAllowedToRead(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this AI system", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:project-ai-system"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectAiSystemLogMetaKey{}, "failed to read {{projectAiSystem.handle}}; insufficient permissions"),
		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrNotAllowedToSearch returns "system:project-ai-system.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrNotAllowedToSearch(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list AI systems", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:project-ai-system"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectAiSystemLogMetaKey{}, "failed to search or list AI systems; insufficient permissions"),
		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrNotAllowedToUpdate returns "system:project-ai-system.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrNotAllowedToUpdate(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this AI system", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:project-ai-system"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectAiSystemLogMetaKey{}, "failed to update {{projectAiSystem.handle}}; insufficient permissions"),
		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrNotAllowedToDelete returns "system:project-ai-system.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrNotAllowedToDelete(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this AI system", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:project-ai-system"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectAiSystemLogMetaKey{}, "failed to delete {{projectAiSystem.handle}}; insufficient permissions"),
		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectAiSystemErrNotAllowedToManageResources returns "system:project-ai-system.notAllowedToManageResources" as *errors.Error
//
// This function is auto-generated.
func ProjectAiSystemErrNotAllowedToManageResources(mm ...*projectAiSystemActionProps) *errors.Error {
	var p = &projectAiSystemActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage resources of this AI system", nil),

		errors.Meta("type", "notAllowedToManageResources"),
		errors.Meta("resource", "system:project-ai-system"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectAiSystemLogMetaKey{}, "failed to manage resources of {{projectAiSystem.handle}}; insufficient permissions"),
		errors.Meta(projectAiSystemPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-ai-system.errors.notAllowedToManageResources"),

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
func (svc projectAiSystem) recordAction(ctx context.Context, props *projectAiSystemActionProps, actionFn func(...*projectAiSystemActionProps) *projectAiSystemAction, err error, diff ...any) error {
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
		a.Description = props.Format(m.AsString(projectAiSystemLogMetaKey{}), err)

		if p, has := m[projectAiSystemPropsMetaKey{}]; has {
			a.Meta = p.(*projectAiSystemActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
