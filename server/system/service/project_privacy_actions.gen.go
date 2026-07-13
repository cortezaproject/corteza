package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/project_privacy_actions.yaml

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
	projectPrivacyActionProps struct {
		projectPrivacy *types.ProjectPrivacy
		new            *types.ProjectPrivacy
		update         *types.ProjectPrivacy
		filter         *types.ProjectPrivacyFilter
		diff           []*revisions.Change
		old            json.RawMessage
	}

	projectPrivacyAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *projectPrivacyActionProps
	}

	projectPrivacyLogMetaKey   struct{}
	projectPrivacyPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setProjectPrivacy updates projectPrivacyActionProps's projectPrivacy
//
// This function is auto-generated.
func (p *projectPrivacyActionProps) setProjectPrivacy(projectPrivacy *types.ProjectPrivacy) *projectPrivacyActionProps {
	p.projectPrivacy = projectPrivacy
	return p
}

// setNew updates projectPrivacyActionProps's new
//
// This function is auto-generated.
func (p *projectPrivacyActionProps) setNew(new *types.ProjectPrivacy) *projectPrivacyActionProps {
	p.new = new
	return p
}

// setUpdate updates projectPrivacyActionProps's update
//
// This function is auto-generated.
func (p *projectPrivacyActionProps) setUpdate(update *types.ProjectPrivacy) *projectPrivacyActionProps {
	p.update = update
	return p
}

// setFilter updates projectPrivacyActionProps's filter
//
// This function is auto-generated.
func (p *projectPrivacyActionProps) setFilter(filter *types.ProjectPrivacyFilter) *projectPrivacyActionProps {
	p.filter = filter
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *projectPrivacyActionProps) setDiff(diff []*revisions.Change) *projectPrivacyActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *projectPrivacyActionProps) setOld(v any) *projectPrivacyActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts projectPrivacyActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p projectPrivacyActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.projectPrivacy != nil {
		m.Set("projectPrivacy.title", p.projectPrivacy.Title, true)
		m.Set("projectPrivacy.ID", p.projectPrivacy.ID, true)
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
func (p projectPrivacyActionProps) Format(in string, err error) string {
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

	if p.projectPrivacy != nil {
		// replacement for "{{projectPrivacy}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{projectPrivacy}}",
			fns(
				p.projectPrivacy.Title,
				p.projectPrivacy.ID,
			),
		)
		pairs = append(pairs, "{{projectPrivacy.title}}", fns(p.projectPrivacy.Title))
		pairs = append(pairs, "{{projectPrivacy.ID}}", fns(p.projectPrivacy.ID))
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
func (a *projectPrivacyAction) String() string {
	var props = &projectPrivacyActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *projectPrivacyAction) ToAction() *actionlog.Action {
	resource := e.resource
	if e.props != nil && e.props.projectPrivacy != nil {
		if r, ok := any(e.props.projectPrivacy).(actionlog.RbacResourcer); ok {
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

// ProjectPrivacyActionSearch returns "system:project-privacy.search" action
//
// This function is auto-generated.
func ProjectPrivacyActionSearch(props ...*projectPrivacyActionProps) *projectPrivacyAction {
	a := &projectPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:project-privacy",
		action:    "search",
		log:       "searched for project privacys",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectPrivacyActionLookup returns "system:project-privacy.lookup" action
//
// This function is auto-generated.
func ProjectPrivacyActionLookup(props ...*projectPrivacyActionProps) *projectPrivacyAction {
	a := &projectPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:project-privacy",
		action:    "lookup",
		log:       "looked-up for a {{projectPrivacy}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectPrivacyActionCreate returns "system:project-privacy.create" action
//
// This function is auto-generated.
func ProjectPrivacyActionCreate(props ...*projectPrivacyActionProps) *projectPrivacyAction {
	a := &projectPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:project-privacy",
		action:    "create",
		log:       "created {{projectPrivacy}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectPrivacyActionUpdate returns "system:project-privacy.update" action
//
// This function is auto-generated.
func ProjectPrivacyActionUpdate(props ...*projectPrivacyActionProps) *projectPrivacyAction {
	a := &projectPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:project-privacy",
		action:    "update",
		log:       "updated {{projectPrivacy}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectPrivacyActionDelete returns "system:project-privacy.delete" action
//
// This function is auto-generated.
func ProjectPrivacyActionDelete(props ...*projectPrivacyActionProps) *projectPrivacyAction {
	a := &projectPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:project-privacy",
		action:    "delete",
		log:       "deleted {{projectPrivacy}}",
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

// ProjectPrivacyErrGeneric returns "system:project-privacy.generic" as *errors.Error
//
// This function is auto-generated.
func ProjectPrivacyErrGeneric(mm ...*projectPrivacyActionProps) *errors.Error {
	var p = &projectPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:project-privacy"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectPrivacyLogMetaKey{}, "{err}"),
		errors.Meta(projectPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-privacy.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectPrivacyErrNotFound returns "system:project-privacy.notFound" as *errors.Error
//
// This function is auto-generated.
func ProjectPrivacyErrNotFound(mm ...*projectPrivacyActionProps) *errors.Error {
	var p = &projectPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project privacy not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:project-privacy"),

		errors.Meta(projectPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-privacy.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectPrivacyErrInvalidID returns "system:project-privacy.invalidID" as *errors.Error
//
// This function is auto-generated.
func ProjectPrivacyErrInvalidID(mm ...*projectPrivacyActionProps) *errors.Error {
	var p = &projectPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:project-privacy"),

		errors.Meta(projectPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-privacy.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectPrivacyErrStaleData returns "system:project-privacy.staleData" as *errors.Error
//
// This function is auto-generated.
func ProjectPrivacyErrStaleData(mm ...*projectPrivacyActionProps) *errors.Error {
	var p = &projectPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project privacy was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:project-privacy"),

		errors.Meta(projectPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-privacy.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectPrivacyErrNotAllowedToCreate returns "system:project-privacy.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ProjectPrivacyErrNotAllowedToCreate(mm ...*projectPrivacyActionProps) *errors.Error {
	var p = &projectPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a project privacy", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:project-privacy"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectPrivacyLogMetaKey{}, "failed to create a project privacy; insufficient permissions"),
		errors.Meta(projectPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-privacy.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectPrivacyErrNotAllowedToRead returns "system:project-privacy.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ProjectPrivacyErrNotAllowedToRead(mm ...*projectPrivacyActionProps) *errors.Error {
	var p = &projectPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this project privacy", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:project-privacy"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectPrivacyLogMetaKey{}, "failed to read {{projectPrivacy.title}}; insufficient permissions"),
		errors.Meta(projectPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-privacy.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectPrivacyErrNotAllowedToSearch returns "system:project-privacy.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ProjectPrivacyErrNotAllowedToSearch(mm ...*projectPrivacyActionProps) *errors.Error {
	var p = &projectPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list project privacys", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:project-privacy"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectPrivacyLogMetaKey{}, "failed to search or list project privacys; insufficient permissions"),
		errors.Meta(projectPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-privacy.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectPrivacyErrNotAllowedToUpdate returns "system:project-privacy.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ProjectPrivacyErrNotAllowedToUpdate(mm ...*projectPrivacyActionProps) *errors.Error {
	var p = &projectPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this project privacy", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:project-privacy"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectPrivacyLogMetaKey{}, "failed to update {{projectPrivacy.title}}; insufficient permissions"),
		errors.Meta(projectPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-privacy.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectPrivacyErrNotAllowedToDelete returns "system:project-privacy.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ProjectPrivacyErrNotAllowedToDelete(mm ...*projectPrivacyActionProps) *errors.Error {
	var p = &projectPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this project privacy", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:project-privacy"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectPrivacyLogMetaKey{}, "failed to delete {{projectPrivacy.title}}; insufficient permissions"),
		errors.Meta(projectPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-privacy.errors.notAllowedToDelete"),

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
func (svc projectPrivacy) recordAction(ctx context.Context, props *projectPrivacyActionProps, actionFn func(...*projectPrivacyActionProps) *projectPrivacyAction, err error, diff ...any) error {
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
		a.Description = props.Format(m.AsString(projectPrivacyLogMetaKey{}), err)

		if p, has := m[projectPrivacyPropsMetaKey{}]; has {
			a.Meta = p.(*projectPrivacyActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
