package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/service/workflow_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"strings"
	"time"
)

type (
	workflowActionProps struct {
		workflow *types.Workflow
		new      *types.Workflow
		update   *types.Workflow
		trigger  *types.Trigger
		filter   *types.WorkflowFilter
	}

	workflowAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *workflowActionProps
	}

	workflowLogMetaKey   struct{}
	workflowPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setWorkflow updates workflowActionProps's workflow
//
// This function is auto-generated.
//
func (p *workflowActionProps) setWorkflow(workflow *types.Workflow) *workflowActionProps {
	p.workflow = workflow
	return p
}

// setNew updates workflowActionProps's new
//
// This function is auto-generated.
//
func (p *workflowActionProps) setNew(new *types.Workflow) *workflowActionProps {
	p.new = new
	return p
}

// setUpdate updates workflowActionProps's update
//
// This function is auto-generated.
//
func (p *workflowActionProps) setUpdate(update *types.Workflow) *workflowActionProps {
	p.update = update
	return p
}

// setTrigger updates workflowActionProps's trigger
//
// This function is auto-generated.
//
func (p *workflowActionProps) setTrigger(trigger *types.Trigger) *workflowActionProps {
	p.trigger = trigger
	return p
}

// setFilter updates workflowActionProps's filter
//
// This function is auto-generated.
//
func (p *workflowActionProps) setFilter(filter *types.WorkflowFilter) *workflowActionProps {
	p.filter = filter
	return p
}

// Serialize converts workflowActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p workflowActionProps) Serialize() actionlog.Meta {
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
	if p.trigger != nil {
		m.Set("trigger.eventType", p.trigger.EventType, true)
		m.Set("trigger.resourceType", p.trigger.ResourceType, true)
		m.Set("trigger.ID", p.trigger.ID, true)
		m.Set("trigger.stepID", p.trigger.StepID, true)
	}
	if p.filter != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p workflowActionProps) Format(in string, err error) string {
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

	if p.workflow != nil {
		// replacement for "{{workflow}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{workflow}}",
			fns(
				p.workflow.Handle,
				p.workflow.ID,
			),
		)
		pairs = append(pairs, "{{workflow.handle}}", fns(p.workflow.Handle))
		pairs = append(pairs, "{{workflow.ID}}", fns(p.workflow.ID))
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

	if p.trigger != nil {
		// replacement for "{{trigger}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{trigger}}",
			fns(
				p.trigger.EventType,
				p.trigger.ResourceType,
				p.trigger.ID,
				p.trigger.StepID,
			),
		)
		pairs = append(pairs, "{{trigger.eventType}}", fns(p.trigger.EventType))
		pairs = append(pairs, "{{trigger.resourceType}}", fns(p.trigger.ResourceType))
		pairs = append(pairs, "{{trigger.ID}}", fns(p.trigger.ID))
		pairs = append(pairs, "{{trigger.stepID}}", fns(p.trigger.StepID))
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
//
func (a *workflowAction) String() string {
	var props = &workflowActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *workflowAction) ToAction() *actionlog.Action {
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

// WorkflowActionSearch returns "automation:workflow.search" action
//
// This function is auto-generated.
//
func WorkflowActionSearch(props ...*workflowActionProps) *workflowAction {
	a := &workflowAction{
		timestamp: time.Now(),
		resource:  "automation:workflow",
		action:    "search",
		log:       "searched for matching workflows",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowActionLookup returns "automation:workflow.lookup" action
//
// This function is auto-generated.
//
func WorkflowActionLookup(props ...*workflowActionProps) *workflowAction {
	a := &workflowAction{
		timestamp: time.Now(),
		resource:  "automation:workflow",
		action:    "lookup",
		log:       "looked-up for a {{workflow}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowActionCreate returns "automation:workflow.create" action
//
// This function is auto-generated.
//
func WorkflowActionCreate(props ...*workflowActionProps) *workflowAction {
	a := &workflowAction{
		timestamp: time.Now(),
		resource:  "automation:workflow",
		action:    "create",
		log:       "created {{workflow}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowActionUpdate returns "automation:workflow.update" action
//
// This function is auto-generated.
//
func WorkflowActionUpdate(props ...*workflowActionProps) *workflowAction {
	a := &workflowAction{
		timestamp: time.Now(),
		resource:  "automation:workflow",
		action:    "update",
		log:       "updated {{workflow}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowActionDelete returns "automation:workflow.delete" action
//
// This function is auto-generated.
//
func WorkflowActionDelete(props ...*workflowActionProps) *workflowAction {
	a := &workflowAction{
		timestamp: time.Now(),
		resource:  "automation:workflow",
		action:    "delete",
		log:       "deleted {{workflow}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowActionUndelete returns "automation:workflow.undelete" action
//
// This function is auto-generated.
//
func WorkflowActionUndelete(props ...*workflowActionProps) *workflowAction {
	a := &workflowAction{
		timestamp: time.Now(),
		resource:  "automation:workflow",
		action:    "undelete",
		log:       "undeleted {{workflow}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// WorkflowActionExecute returns "automation:workflow.execute" action
//
// This function is auto-generated.
//
func WorkflowActionExecute(props ...*workflowActionProps) *workflowAction {
	a := &workflowAction{
		timestamp: time.Now(),
		resource:  "automation:workflow",
		action:    "execute",
		log:       "{{workflow}} executed",
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

// WorkflowErrGeneric returns "automation:workflow.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrGeneric(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "{err}"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrNotFound returns "automation:workflow.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrNotFound(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("workflow not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "automation:workflow"),

		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrInvalidID returns "automation:workflow.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrInvalidID(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "automation:workflow"),

		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrDisabled returns "automation:workflow.disabled" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrDisabled(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("disabled workflow or trigger", nil),

		errors.Meta("type", "disabled"),
		errors.Meta("resource", "automation:workflow"),

		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.disabled"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrInvalidHandle returns "automation:workflow.invalidHandle" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrInvalidHandle(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "automation:workflow"),

		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrStaleData returns "automation:workflow.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrStaleData(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "automation:workflow"),

		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrNotAllowedToRead returns "automation:workflow.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrNotAllowedToRead(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this workflow", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "failed to read {{workflow.handle}}; insufficient permissions"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrNotAllowedToSearch returns "automation:workflow.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrNotAllowedToSearch(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list workflows", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "failed to search or list workflow; insufficient permissions"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrNotAllowedToCreate returns "automation:workflow.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrNotAllowedToCreate(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create workflows", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "failed to create workflow; insufficient permissions"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrNotAllowedToUpdate returns "automation:workflow.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrNotAllowedToUpdate(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this workflow", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "failed to update {{workflow}}; insufficient permissions"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrNotAllowedToDelete returns "automation:workflow.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrNotAllowedToDelete(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this workflow", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "failed to delete {{workflow}}; insufficient permissions"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrNotAllowedToUndelete returns "automation:workflow.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrNotAllowedToUndelete(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this workflow", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "failed to undelete {{workflow}}; insufficient permissions"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrNotAllowedToExecute returns "automation:workflow.notAllowedToExecute" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrNotAllowedToExecute(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to execute this workflow", nil),

		errors.Meta("type", "notAllowedToExecute"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "failed to execute {{workflow}}; insufficient permissions"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.notAllowedToExecute"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrUnknownWorkflowStep returns "automation:workflow.unknownWorkflowStep" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrUnknownWorkflowStep(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("unknown workflow step", nil),

		errors.Meta("type", "unknownWorkflowStep"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "failed to execute {{workflow}}; unknown workflow step"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.unknownWorkflowStep"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrHandleNotUnique returns "automation:workflow.handleNotUnique" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrHandleNotUnique(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("workflow handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "duplicate handle used for workflow ({{workflow}})"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrNotAllowedToExecuteCorredorStep returns "automation:workflow.notAllowedToExecuteCorredorStep" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrNotAllowedToExecuteCorredorStep(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to run corredorExec function, corredor is disabled", nil),

		errors.Meta("type", "notAllowedToExecuteCorredorStep"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "failed to execute {{workflow}} with corredorExec function step; corredor is disabled"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.notAllowedToExecuteCorredorStep"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// WorkflowErrMaximumCallStackSizeExceeded returns "automation:workflow.maximumCallStackSizeExceeded" as *errors.Error
//
//
// This function is auto-generated.
//
func WorkflowErrMaximumCallStackSizeExceeded(mm ...*workflowActionProps) *errors.Error {
	var p = &workflowActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("maximum call stack size exceeded", nil),

		errors.Meta("type", "maximumCallStackSizeExceeded"),
		errors.Meta("resource", "automation:workflow"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(workflowLogMetaKey{}, "maximum call stack size exceeded"),
		errors.Meta(workflowPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "workflow.errors.maximumCallStackSizeExceeded"),

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
func (svc workflow) recordAction(ctx context.Context, props *workflowActionProps, actionFn func(...*workflowActionProps) *workflowAction, err error) error {
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
		a.Description = props.Format(m.AsString(workflowLogMetaKey{}), err)

		if p, has := m[workflowPropsMetaKey{}]; has {
			a.Meta = p.(*workflowActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
