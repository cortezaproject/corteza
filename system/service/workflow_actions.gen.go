package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/workflow_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
	"time"
)

type (
	workflowServiceActionProps struct {
		workflow *types.Workflow
		new      *types.Workflow
		update   *types.Workflow
		filter   *types.WorkflowFilter
	}

	workflowServiceAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *workflowServiceActionProps
	}

	workflowServiceLogMetaKey   struct{}
	workflowServicePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setWorkflow updates workflowServiceActionProps's workflow
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *workflowServiceActionProps) setWorkflow(workflow *types.Workflow) *workflowServiceActionProps {
	p.workflow = workflow
	return p
}

// setNew updates workflowServiceActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *workflowServiceActionProps) setNew(new *types.Workflow) *workflowServiceActionProps {
	p.new = new
	return p
}

// setUpdate updates workflowServiceActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *workflowServiceActionProps) setUpdate(update *types.Workflow) *workflowServiceActionProps {
	p.update = update
	return p
}

// setFilter updates workflowServiceActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *workflowServiceActionProps) setFilter(filter *types.WorkflowFilter) *workflowServiceActionProps {
	p.filter = filter
	return p
}

// Serialize converts workflowServiceActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p workflowServiceActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.workflow != nil {
		m.Set("workflow.handle", p.workflow.Handle, true)
		m.Set("workflow.ID", p.workflow.ID, true)
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
//
func (p workflowServiceActionProps) Format(in string, err error) string {
	var (
		pairs = []string{"{err}"}
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

	if p.workflow != nil {
		// replacement for "{workflow}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{workflow}",
			fns(
				p.workflow.Handle,
				p.workflow.ID,
			),
		)
		pairs = append(pairs, "{workflow.handle}", fns(p.workflow.Handle))
		pairs = append(pairs, "{workflow.ID}", fns(p.workflow.ID))
	}

	if p.new != nil {
		// replacement for "{new}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{new}",
			fns(
				p.new.Handle,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{new.handle}", fns(p.new.Handle))
		pairs = append(pairs, "{new.ID}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{update}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{update}",
			fns(
				p.update.Handle,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{update.handle}", fns(p.update.Handle))
		pairs = append(pairs, "{update.ID}", fns(p.update.ID))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
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
//
func (a *workflowServiceAction) String() string {
	var props = &workflowServiceActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *workflowServiceAction) ToAction() *actionlog.Action {
	return &actionlog.Action{
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.Serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// WorkflowServiceActionSearch returns "system:workflow.search" action
//
// This function is auto-generated.
//
func WorkflowServiceActionSearch(props ...*workflowServiceActionProps) *workflowServiceAction {
	a := &workflowServiceAction{
		timestamp: time.Now(),
		resource:  "system:workflow",
		action:    "search",
		log:       "searched for matching workflows",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowServiceActionLookup returns "system:workflow.lookup" action
//
// This function is auto-generated.
//
func WorkflowServiceActionLookup(props ...*workflowServiceActionProps) *workflowServiceAction {
	a := &workflowServiceAction{
		timestamp: time.Now(),
		resource:  "system:workflow",
		action:    "lookup",
		log:       "looked-up for a {workflow}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowServiceActionCreate returns "system:workflow.create" action
//
// This function is auto-generated.
//
func WorkflowServiceActionCreate(props ...*workflowServiceActionProps) *workflowServiceAction {
	a := &workflowServiceAction{
		timestamp: time.Now(),
		resource:  "system:workflow",
		action:    "create",
		log:       "created {workflow}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowServiceActionUpdate returns "system:workflow.update" action
//
// This function is auto-generated.
//
func WorkflowServiceActionUpdate(props ...*workflowServiceActionProps) *workflowServiceAction {
	a := &workflowServiceAction{
		timestamp: time.Now(),
		resource:  "system:workflow",
		action:    "update",
		log:       "updated {workflow}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowServiceActionDelete returns "system:workflow.delete" action
//
// This function is auto-generated.
//
func WorkflowServiceActionDelete(props ...*workflowServiceActionProps) *workflowServiceAction {
	a := &workflowServiceAction{
		timestamp: time.Now(),
		resource:  "system:workflow",
		action:    "delete",
		log:       "deleted {workflow}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowServiceActionUndelete returns "system:workflow.undelete" action
//
// This function is auto-generated.
//
func WorkflowServiceActionUndelete(props ...*workflowServiceActionProps) *workflowServiceAction {
	a := &workflowServiceAction{
		timestamp: time.Now(),
		resource:  "system:workflow",
		action:    "undelete",
		log:       "undeleted {workflow}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// WorkflowServiceErrGeneric returns "system:workflow.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrGeneric(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowServiceLogMetaKey{}, "{err}"),
		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrNotFound returns "system:workflow.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrNotFound(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("workflow not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:workflow"),

		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrInvalidID returns "system:workflow.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrInvalidID(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:workflow"),

		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrInvalidHandle returns "system:workflow.invalidHandle" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrInvalidHandle(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:workflow"),

		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrStaleData returns "system:workflow.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrStaleData(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:workflow"),

		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrNotAllowedToRead returns "system:workflow.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrNotAllowedToRead(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this workflow", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowServiceLogMetaKey{}, "failed to read {workflow.handle}; insufficient permissions"),
		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrNotAllowedToCreate returns "system:workflow.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrNotAllowedToCreate(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create workflows", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowServiceLogMetaKey{}, "failed to create workflow; insufficient permissions"),
		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrNotAllowedToUpdate returns "system:workflow.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrNotAllowedToUpdate(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this workflow", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowServiceLogMetaKey{}, "failed to update {workflow}; insufficient permissions"),
		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrNotAllowedToDelete returns "system:workflow.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrNotAllowedToDelete(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this workflow", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowServiceLogMetaKey{}, "failed to delete {workflow}; insufficient permissions"),
		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrNotAllowedToUndelete returns "system:workflow.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrNotAllowedToUndelete(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this workflow", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowServiceLogMetaKey{}, "failed to undelete {workflow}; insufficient permissions"),
		errors.Meta(workflowServicePropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowServiceErrHandleNotUnique returns "system:workflow.handleNotUnique" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowServiceErrHandleNotUnique(mm ...*workflowServiceActionProps) *errors.Error {
	var p = &workflowServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "system:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowServiceLogMetaKey{}, "used duplicate handle ({workflow}) for workflow"),
		errors.Meta(workflowServicePropsMetaKey{}, p),

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
//
func (svc workflowService) recordAction(ctx context.Context, props *workflowServiceActionProps, actionFn func(...*workflowServiceActionProps) *workflowServiceAction, err error) error {
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
		a.Description = props.Format(m.AsString(workflowServiceLogMetaKey{}), err)

		if p, has := m[workflowServicePropsMetaKey{}]; has {
			a.Meta = p.(*workflowServiceActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
