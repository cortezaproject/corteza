package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/project_review_actions.yaml

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
	projectReviewActionProps struct {
		projectReview *types.ProjectReview
		new           *types.ProjectReview
		update        *types.ProjectReview
		filter        *types.ProjectReviewFilter
		diff          []*revisions.Change
		old           json.RawMessage
	}

	projectReviewAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *projectReviewActionProps
	}

	projectReviewLogMetaKey   struct{}
	projectReviewPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setProjectReview updates projectReviewActionProps's projectReview
//
// This function is auto-generated.
func (p *projectReviewActionProps) setProjectReview(projectReview *types.ProjectReview) *projectReviewActionProps {
	p.projectReview = projectReview
	return p
}

// setNew updates projectReviewActionProps's new
//
// This function is auto-generated.
func (p *projectReviewActionProps) setNew(new *types.ProjectReview) *projectReviewActionProps {
	p.new = new
	return p
}

// setUpdate updates projectReviewActionProps's update
//
// This function is auto-generated.
func (p *projectReviewActionProps) setUpdate(update *types.ProjectReview) *projectReviewActionProps {
	p.update = update
	return p
}

// setFilter updates projectReviewActionProps's filter
//
// This function is auto-generated.
func (p *projectReviewActionProps) setFilter(filter *types.ProjectReviewFilter) *projectReviewActionProps {
	p.filter = filter
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *projectReviewActionProps) setDiff(diff []*revisions.Change) *projectReviewActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *projectReviewActionProps) setOld(v any) *projectReviewActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts projectReviewActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p projectReviewActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.projectReview != nil {
		m.Set("projectReview.title", p.projectReview.Title, true)
		m.Set("projectReview.ID", p.projectReview.ID, true)
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
func (p projectReviewActionProps) Format(in string, err error) string {
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

	if p.projectReview != nil {
		// replacement for "{{projectReview}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{projectReview}}",
			fns(
				p.projectReview.Title,
				p.projectReview.ID,
			),
		)
		pairs = append(pairs, "{{projectReview.title}}", fns(p.projectReview.Title))
		pairs = append(pairs, "{{projectReview.ID}}", fns(p.projectReview.ID))
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
func (a *projectReviewAction) String() string {
	var props = &projectReviewActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *projectReviewAction) ToAction() *actionlog.Action {
	resource := e.resource
	var resourceProjectID uint64
	if e.props != nil && e.props.projectReview != nil {
		if r, ok := any(e.props.projectReview).(actionlog.RbacResourcer); ok {
			resource = r.RbacResource()
		}

		// Attribute the event to the project owning the affected resource. This is
		// independent of the request scope, so the log stays filterable per project
		// on routes that carry no project.
		if r, ok := any(e.props.projectReview).(actionlog.ProjectResourcer); ok {
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

// ProjectReviewActionSearch returns "system:project-review.search" action
//
// This function is auto-generated.
func ProjectReviewActionSearch(props ...*projectReviewActionProps) *projectReviewAction {
	a := &projectReviewAction{
		timestamp: time.Now(),
		resource:  "system:project-review",
		action:    "search",
		log:       "searched for project reviews",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectReviewActionLookup returns "system:project-review.lookup" action
//
// This function is auto-generated.
func ProjectReviewActionLookup(props ...*projectReviewActionProps) *projectReviewAction {
	a := &projectReviewAction{
		timestamp: time.Now(),
		resource:  "system:project-review",
		action:    "lookup",
		log:       "looked-up for a {{projectReview}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectReviewActionCreate returns "system:project-review.create" action
//
// This function is auto-generated.
func ProjectReviewActionCreate(props ...*projectReviewActionProps) *projectReviewAction {
	a := &projectReviewAction{
		timestamp: time.Now(),
		resource:  "system:project-review",
		action:    "create",
		log:       "created {{projectReview}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectReviewActionUpdate returns "system:project-review.update" action
//
// This function is auto-generated.
func ProjectReviewActionUpdate(props ...*projectReviewActionProps) *projectReviewAction {
	a := &projectReviewAction{
		timestamp: time.Now(),
		resource:  "system:project-review",
		action:    "update",
		log:       "updated {{projectReview}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ProjectReviewActionDelete returns "system:project-review.delete" action
//
// This function is auto-generated.
func ProjectReviewActionDelete(props ...*projectReviewActionProps) *projectReviewAction {
	a := &projectReviewAction{
		timestamp: time.Now(),
		resource:  "system:project-review",
		action:    "delete",
		log:       "deleted {{projectReview}}",
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

// ProjectReviewErrGeneric returns "system:project-review.generic" as *errors.Error
//
// This function is auto-generated.
func ProjectReviewErrGeneric(mm ...*projectReviewActionProps) *errors.Error {
	var p = &projectReviewActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:project-review"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectReviewLogMetaKey{}, "{err}"),
		errors.Meta(projectReviewPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-review.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectReviewErrNotFound returns "system:project-review.notFound" as *errors.Error
//
// This function is auto-generated.
func ProjectReviewErrNotFound(mm ...*projectReviewActionProps) *errors.Error {
	var p = &projectReviewActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project review not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:project-review"),

		errors.Meta(projectReviewPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-review.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectReviewErrInvalidID returns "system:project-review.invalidID" as *errors.Error
//
// This function is auto-generated.
func ProjectReviewErrInvalidID(mm ...*projectReviewActionProps) *errors.Error {
	var p = &projectReviewActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:project-review"),

		errors.Meta(projectReviewPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-review.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectReviewErrStaleData returns "system:project-review.staleData" as *errors.Error
//
// This function is auto-generated.
func ProjectReviewErrStaleData(mm ...*projectReviewActionProps) *errors.Error {
	var p = &projectReviewActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("project review was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:project-review"),

		errors.Meta(projectReviewPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-review.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectReviewErrNotAllowedToCreate returns "system:project-review.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ProjectReviewErrNotAllowedToCreate(mm ...*projectReviewActionProps) *errors.Error {
	var p = &projectReviewActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a project review", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:project-review"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectReviewLogMetaKey{}, "failed to create a project review; insufficient permissions"),
		errors.Meta(projectReviewPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-review.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectReviewErrNotAllowedToRead returns "system:project-review.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ProjectReviewErrNotAllowedToRead(mm ...*projectReviewActionProps) *errors.Error {
	var p = &projectReviewActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this project review", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:project-review"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectReviewLogMetaKey{}, "failed to read {{projectReview.title}}; insufficient permissions"),
		errors.Meta(projectReviewPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-review.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectReviewErrNotAllowedToSearch returns "system:project-review.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ProjectReviewErrNotAllowedToSearch(mm ...*projectReviewActionProps) *errors.Error {
	var p = &projectReviewActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list project reviews", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:project-review"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectReviewLogMetaKey{}, "failed to search or list project reviews; insufficient permissions"),
		errors.Meta(projectReviewPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-review.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectReviewErrNotAllowedToUpdate returns "system:project-review.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ProjectReviewErrNotAllowedToUpdate(mm ...*projectReviewActionProps) *errors.Error {
	var p = &projectReviewActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this project review", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:project-review"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectReviewLogMetaKey{}, "failed to update {{projectReview.title}}; insufficient permissions"),
		errors.Meta(projectReviewPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-review.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ProjectReviewErrNotAllowedToDelete returns "system:project-review.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ProjectReviewErrNotAllowedToDelete(mm ...*projectReviewActionProps) *errors.Error {
	var p = &projectReviewActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this project review", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:project-review"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(projectReviewLogMetaKey{}, "failed to delete {{projectReview.title}}; insufficient permissions"),
		errors.Meta(projectReviewPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "project-review.errors.notAllowedToDelete"),

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
func (svc projectReview) recordAction(ctx context.Context, props *projectReviewActionProps, actionFn func(...*projectReviewActionProps) *projectReviewAction, err error, diff ...any) error {
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
		a.Description = props.Format(m.AsString(projectReviewLogMetaKey{}), err)

		if p, has := m[projectReviewPropsMetaKey{}]; has {
			a.Meta = p.(*projectReviewActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
