package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/service/ng_automation_actions.yaml

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
	ngAutomationActionProps struct {
		ngAutomation *types.NgAutomation
		new          *types.NgAutomation
		update       *types.NgAutomation
		trigger      *types.Trigger
		filter       *types.NgAutomationFilter
	}

	ngAutomationAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *ngAutomationActionProps
	}

	ngAutomationLogMetaKey   struct{}
	ngAutomationPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setNgAutomation updates ngAutomationActionProps's ngAutomation
//
// This function is auto-generated.
func (p *ngAutomationActionProps) setNgAutomation(ngAutomation *types.NgAutomation) *ngAutomationActionProps {
	p.ngAutomation = ngAutomation
	return p
}

// setNew updates ngAutomationActionProps's new
//
// This function is auto-generated.
func (p *ngAutomationActionProps) setNew(new *types.NgAutomation) *ngAutomationActionProps {
	p.new = new
	return p
}

// setUpdate updates ngAutomationActionProps's update
//
// This function is auto-generated.
func (p *ngAutomationActionProps) setUpdate(update *types.NgAutomation) *ngAutomationActionProps {
	p.update = update
	return p
}

// setTrigger updates ngAutomationActionProps's trigger
//
// This function is auto-generated.
func (p *ngAutomationActionProps) setTrigger(trigger *types.Trigger) *ngAutomationActionProps {
	p.trigger = trigger
	return p
}

// setFilter updates ngAutomationActionProps's filter
//
// This function is auto-generated.
func (p *ngAutomationActionProps) setFilter(filter *types.NgAutomationFilter) *ngAutomationActionProps {
	p.filter = filter
	return p
}

// Serialize converts ngAutomationActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p ngAutomationActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.ngAutomation != nil {
		m.Set("ngAutomation.handle", p.ngAutomation.Handle, true)
		m.Set("ngAutomation.ID", p.ngAutomation.ID, true)
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
func (p ngAutomationActionProps) Format(in string, err error) string {
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

	if p.ngAutomation != nil {
		// replacement for "{{ngAutomation}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{ngAutomation}}",
			fns(
				p.ngAutomation.Handle,
				p.ngAutomation.ID,
			),
		)
		pairs = append(pairs, "{{ngAutomation.handle}}", fns(p.ngAutomation.Handle))
		pairs = append(pairs, "{{ngAutomation.ID}}", fns(p.ngAutomation.ID))
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
func (a *ngAutomationAction) String() string {
	var props = &ngAutomationActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *ngAutomationAction) ToAction() *actionlog.Action {
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

// NgAutomationActionSearch returns "automation:ng-automation.search" action
//
// This function is auto-generated.
func NgAutomationActionSearch(props ...*ngAutomationActionProps) *ngAutomationAction {
	a := &ngAutomationAction{
		timestamp: time.Now(),
		resource:  "automation:ng-automation",
		action:    "search",
		log:       "searched for matching ngAutomations",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NgAutomationActionLookup returns "automation:ng-automation.lookup" action
//
// This function is auto-generated.
func NgAutomationActionLookup(props ...*ngAutomationActionProps) *ngAutomationAction {
	a := &ngAutomationAction{
		timestamp: time.Now(),
		resource:  "automation:ng-automation",
		action:    "lookup",
		log:       "looked-up for a {{ngAutomation}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NgAutomationActionCreate returns "automation:ng-automation.create" action
//
// This function is auto-generated.
func NgAutomationActionCreate(props ...*ngAutomationActionProps) *ngAutomationAction {
	a := &ngAutomationAction{
		timestamp: time.Now(),
		resource:  "automation:ng-automation",
		action:    "create",
		log:       "created {{ngAutomation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NgAutomationActionUpdate returns "automation:ng-automation.update" action
//
// This function is auto-generated.
func NgAutomationActionUpdate(props ...*ngAutomationActionProps) *ngAutomationAction {
	a := &ngAutomationAction{
		timestamp: time.Now(),
		resource:  "automation:ng-automation",
		action:    "update",
		log:       "updated {{ngAutomation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NgAutomationActionDelete returns "automation:ng-automation.delete" action
//
// This function is auto-generated.
func NgAutomationActionDelete(props ...*ngAutomationActionProps) *ngAutomationAction {
	a := &ngAutomationAction{
		timestamp: time.Now(),
		resource:  "automation:ng-automation",
		action:    "delete",
		log:       "deleted {{ngAutomation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NgAutomationActionUndelete returns "automation:ng-automation.undelete" action
//
// This function is auto-generated.
func NgAutomationActionUndelete(props ...*ngAutomationActionProps) *ngAutomationAction {
	a := &ngAutomationAction{
		timestamp: time.Now(),
		resource:  "automation:ng-automation",
		action:    "undelete",
		log:       "undeleted {{ngAutomation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NgAutomationActionExecute returns "automation:ng-automation.execute" action
//
// This function is auto-generated.
func NgAutomationActionExecute(props ...*ngAutomationActionProps) *ngAutomationAction {
	a := &ngAutomationAction{
		timestamp: time.Now(),
		resource:  "automation:ng-automation",
		action:    "execute",
		log:       "{{ngAutomation}} executed",
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

// NgAutomationErrGeneric returns "automation:ng-automation.generic" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrGeneric(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "{err}"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrNotFound returns "automation:ng-automation.notFound" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrNotFound(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("ngAutomation not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "automation:ng-automation"),

		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrInvalidID returns "automation:ng-automation.invalidID" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrInvalidID(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "automation:ng-automation"),

		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrDisabled returns "automation:ng-automation.disabled" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrDisabled(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("disabled ngAutomation or trigger", nil),

		errors.Meta("type", "disabled"),
		errors.Meta("resource", "automation:ng-automation"),

		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.disabled"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrInvalidHandle returns "automation:ng-automation.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrInvalidHandle(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "automation:ng-automation"),

		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrMissingName returns "automation:ng-automation.missingName" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrMissingName(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("missing name", nil),

		errors.Meta("type", "missingName"),
		errors.Meta("resource", "automation:ng-automation"),

		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.missingName"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrStaleData returns "automation:ng-automation.staleData" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrStaleData(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("ngAutomation was modified by someone else or by a ngAutomation after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "automation:ng-automation"),

		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrNotAllowedToRead returns "automation:ng-automation.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrNotAllowedToRead(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this ngAutomation", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "failed to read {{ngAutomation.handle}}; insufficient permissions"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrNotAllowedToSearch returns "automation:ng-automation.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrNotAllowedToSearch(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list ngAutomations", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "failed to search or list ngAutomation; insufficient permissions"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrNotAllowedToCreate returns "automation:ng-automation.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrNotAllowedToCreate(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create ngAutomations", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "failed to create ngAutomation; insufficient permissions"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrNotAllowedToUpdate returns "automation:ng-automation.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrNotAllowedToUpdate(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this ngAutomation", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "failed to update {{ngAutomation}}; insufficient permissions"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrNotAllowedToDelete returns "automation:ng-automation.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrNotAllowedToDelete(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this ngAutomation", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "failed to delete {{ngAutomation}}; insufficient permissions"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrNotAllowedToUndelete returns "automation:ng-automation.notAllowedToUndelete" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrNotAllowedToUndelete(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this ngAutomation", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "failed to undelete {{ngAutomation}}; insufficient permissions"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrNotAllowedToExecute returns "automation:ng-automation.notAllowedToExecute" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrNotAllowedToExecute(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to execute this ngAutomation", nil),

		errors.Meta("type", "notAllowedToExecute"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "failed to execute {{ngAutomation}}; insufficient permissions"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.notAllowedToExecute"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrUnknownNgAutomationStep returns "automation:ng-automation.unknownNgAutomationStep" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrUnknownNgAutomationStep(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("unknown ngAutomation step", nil),

		errors.Meta("type", "unknownNgAutomationStep"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "failed to execute {{ngAutomation}}; unknown ngAutomation step"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.unknownNgAutomationStep"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrHandleNotUnique returns "automation:ng-automation.handleNotUnique" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrHandleNotUnique(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("ngAutomation handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "duplicate handle used for ngAutomation ({{ngAutomation}})"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrNotAllowedToExecuteCorredorStep returns "automation:ng-automation.notAllowedToExecuteCorredorStep" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrNotAllowedToExecuteCorredorStep(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to run corredorExec function, corredor is disabled", nil),

		errors.Meta("type", "notAllowedToExecuteCorredorStep"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "failed to execute {{ngAutomation}} with corredorExec function step; corredor is disabled"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.notAllowedToExecuteCorredorStep"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NgAutomationErrMaximumCallStackSizeExceeded returns "automation:ng-automation.maximumCallStackSizeExceeded" as *errors.Error
//
// This function is auto-generated.
func NgAutomationErrMaximumCallStackSizeExceeded(mm ...*ngAutomationActionProps) *errors.Error {
	var p = &ngAutomationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("maximum call stack size exceeded", nil),

		errors.Meta("type", "maximumCallStackSizeExceeded"),
		errors.Meta("resource", "automation:ng-automation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(ngAutomationLogMetaKey{}, "maximum call stack size exceeded"),
		errors.Meta(ngAutomationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "ng-automation.errors.maximumCallStackSizeExceeded"),

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
func (svc ngAutomation) recordAction(ctx context.Context, props *ngAutomationActionProps, actionFn func(...*ngAutomationActionProps) *ngAutomationAction, err error) error {
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
		a.Description = props.Format(m.AsString(ngAutomationLogMetaKey{}), err)

		if p, has := m[ngAutomationPropsMetaKey{}]; has {
			a.Meta = p.(*ngAutomationActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
