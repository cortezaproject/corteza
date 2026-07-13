package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/project_feature_actions.yaml

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
	projectFeatureActionProps struct {
		projectFeature *types.ProjectFeature
		new            *types.ProjectFeature
		update         *types.ProjectFeature
		filter         *types.ProjectFeatureFilter
		diff           []*revisions.Change
		old            json.RawMessage
	}

	projectFeatureAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *projectFeatureActionProps
	}

	projectFeatureLogMetaKey   struct{}
	projectFeaturePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setProjectFeature updates projectFeatureActionProps's projectFeature
//
// This function is auto-generated.
func (p *projectFeatureActionProps) setProjectFeature(projectFeature *types.ProjectFeature) *projectFeatureActionProps {
	p.projectFeature = projectFeature
	return p
}

// setNew updates projectFeatureActionProps's new
//
// This function is auto-generated.
func (p *projectFeatureActionProps) setNew(new *types.ProjectFeature) *projectFeatureActionProps {
	p.new = new
	return p
}

// setUpdate updates projectFeatureActionProps's update
//
// This function is auto-generated.
func (p *projectFeatureActionProps) setUpdate(update *types.ProjectFeature) *projectFeatureActionProps {
	p.update = update
	return p
}

// setFilter updates projectFeatureActionProps's filter
//
// This function is auto-generated.
func (p *projectFeatureActionProps) setFilter(filter *types.ProjectFeatureFilter) *projectFeatureActionProps {
	p.filter = filter
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *projectFeatureActionProps) setDiff(diff []*revisions.Change) *projectFeatureActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *projectFeatureActionProps) setOld(v any) *projectFeatureActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts projectFeatureActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p projectFeatureActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.projectFeature != nil {
		m.Set("projectFeature.title", p.projectFeature.Title, true)
		m.Set("projectFeature.ID", p.projectFeature.ID, true)
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
func (p projectFeatureActionProps) Format(in string, err error) string {
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

	if p.projectFeature != nil {
		// replacement for "{{projectFeature}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{projectFeature}}",
			fns(
				p.projectFeature.Title,
				p.projectFeature.ID,
			),
		)
		pairs = append(pairs, "{{projectFeature.title}}", fns(p.projectFeature.Title))
		pairs = append(pairs, "{{projectFeature.ID}}", fns(p.projectFeature.ID))
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
func (a *projectFeatureAction) String() string {
	var props = &projectFeatureActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *projectFeatureAction) ToAction() *actionlog.Action {
	resource := e.resource
	if e.props != nil && e.props.projectFeature != nil {
		if r, ok := any(e.props.projectFeature).(actionlog.RbacResourcer); ok {
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

// ProjectFeatureActionSearch returns "system:project-feature.search" action
//
// This function is auto-generated.
func ProjectFeatureActionSearch(props ...*projectFeatureActionProps) *projectFeatureAction {
	a := &projectFeatureAction{
		timestamp: time.Now(),
		resource:  "system:project-feature",
		action:    "search",
		log:       "searched for project features",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectFeatureActionLookup returns "system:project-feature.lookup" action
//
// This function is auto-generated.
func ProjectFeatureActionLookup(props ...*projectFeatureActionProps) *projectFeatureAction {
	a := &projectFeatureAction{
		timestamp: time.Now(),
		resource:  "system:project-feature",
		action:    "lookup",
		log:       "looked-up for a {{projectFeature}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectFeatureActionCreate returns "system:project-feature.create" action
//
// This function is auto-generated.
func ProjectFeatureActionCreate(props ...*projectFeatureActionProps) *projectFeatureAction {
	a := &projectFeatureAction{
		timestamp: time.Now(),
		resource:  "system:project-feature",
		action:    "create",
		log:       "created {{projectFeature}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectFeatureActionUpdate returns "system:project-feature.update" action
//
// This function is auto-generated.
func ProjectFeatureActionUpdate(props ...*projectFeatureActionProps) *projectFeatureAction {
	a := &projectFeatureAction{
		timestamp: time.Now(),
		resource:  "system:project-feature",
		action:    "update",
		log:       "updated {{projectFeature}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectFeatureActionDelete returns "system:project-feature.delete" action
//
// This function is auto-generated.
func ProjectFeatureActionDelete(props ...*projectFeatureActionProps) *projectFeatureAction {
	a := &projectFeatureAction{
		timestamp: time.Now(),
		resource:  "system:project-feature",
		action:    "delete",
		log:       "deleted {{projectFeature}}",
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

// ProjectFeatureErrGeneric returns "system:project-feature.generic" as *errors.Error
//
// This function is auto-generated.
func ProjectFeatureErrGeneric(mm ...*projectFeatureActionProps) *errors.Error {
	var p = &projectFeatureActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:project-feature"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFeatureLogMetaKey{}, "{err}"),
		errors.Meta(projectFeaturePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-feature.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFeatureErrNotFound returns "system:project-feature.notFound" as *errors.Error
//
// This function is auto-generated.
func ProjectFeatureErrNotFound(mm ...*projectFeatureActionProps) *errors.Error {
	var p = &projectFeatureActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project feature not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:project-feature"),

		errors.Meta(projectFeaturePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-feature.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFeatureErrInvalidID returns "system:project-feature.invalidID" as *errors.Error
//
// This function is auto-generated.
func ProjectFeatureErrInvalidID(mm ...*projectFeatureActionProps) *errors.Error {
	var p = &projectFeatureActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:project-feature"),

		errors.Meta(projectFeaturePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-feature.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFeatureErrStaleData returns "system:project-feature.staleData" as *errors.Error
//
// This function is auto-generated.
func ProjectFeatureErrStaleData(mm ...*projectFeatureActionProps) *errors.Error {
	var p = &projectFeatureActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project feature was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:project-feature"),

		errors.Meta(projectFeaturePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-feature.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFeatureErrNotAllowedToCreate returns "system:project-feature.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ProjectFeatureErrNotAllowedToCreate(mm ...*projectFeatureActionProps) *errors.Error {
	var p = &projectFeatureActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a project feature", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:project-feature"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFeatureLogMetaKey{}, "failed to create a project feature; insufficient permissions"),
		errors.Meta(projectFeaturePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-feature.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFeatureErrNotAllowedToRead returns "system:project-feature.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ProjectFeatureErrNotAllowedToRead(mm ...*projectFeatureActionProps) *errors.Error {
	var p = &projectFeatureActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this project feature", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:project-feature"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFeatureLogMetaKey{}, "failed to read {{projectFeature.title}}; insufficient permissions"),
		errors.Meta(projectFeaturePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-feature.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFeatureErrNotAllowedToSearch returns "system:project-feature.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ProjectFeatureErrNotAllowedToSearch(mm ...*projectFeatureActionProps) *errors.Error {
	var p = &projectFeatureActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list project features", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:project-feature"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFeatureLogMetaKey{}, "failed to search or list project features; insufficient permissions"),
		errors.Meta(projectFeaturePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-feature.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFeatureErrNotAllowedToUpdate returns "system:project-feature.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ProjectFeatureErrNotAllowedToUpdate(mm ...*projectFeatureActionProps) *errors.Error {
	var p = &projectFeatureActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this project feature", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:project-feature"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFeatureLogMetaKey{}, "failed to update {{projectFeature.title}}; insufficient permissions"),
		errors.Meta(projectFeaturePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-feature.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectFeatureErrNotAllowedToDelete returns "system:project-feature.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ProjectFeatureErrNotAllowedToDelete(mm ...*projectFeatureActionProps) *errors.Error {
	var p = &projectFeatureActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this project feature", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:project-feature"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectFeatureLogMetaKey{}, "failed to delete {{projectFeature.title}}; insufficient permissions"),
		errors.Meta(projectFeaturePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-feature.errors.notAllowedToDelete"),

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
func (svc projectFeature) recordAction(ctx context.Context, props *projectFeatureActionProps, actionFn func(...*projectFeatureActionProps) *projectFeatureAction, err error, diff ...any) error {
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
		a.Description = props.Format(m.AsString(projectFeatureLogMetaKey{}), err)

		if p, has := m[projectFeaturePropsMetaKey{}]; has {
			a.Meta = p.(*projectFeatureActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
