package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/service/trigger_definition_actions.yaml

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
	triggerDefinitionActionProps struct {
		triggerDefinition *types.TriggerDefinition
		new               *types.TriggerDefinition
		update            *types.TriggerDefinition
		filter            *types.TriggerDefinitionFilter
	}

	triggerDefinitionAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *triggerDefinitionActionProps
	}

	triggerDefinitionLogMetaKey   struct{}
	triggerDefinitionPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setTriggerDefinition updates triggerDefinitionActionProps's triggerDefinition
//
// This function is auto-generated.
func (p *triggerDefinitionActionProps) setTriggerDefinition(triggerDefinition *types.TriggerDefinition) *triggerDefinitionActionProps {
	p.triggerDefinition = triggerDefinition
	return p
}

// setNew updates triggerDefinitionActionProps's new
//
// This function is auto-generated.
func (p *triggerDefinitionActionProps) setNew(new *types.TriggerDefinition) *triggerDefinitionActionProps {
	p.new = new
	return p
}

// setUpdate updates triggerDefinitionActionProps's update
//
// This function is auto-generated.
func (p *triggerDefinitionActionProps) setUpdate(update *types.TriggerDefinition) *triggerDefinitionActionProps {
	p.update = update
	return p
}

// setFilter updates triggerDefinitionActionProps's filter
//
// This function is auto-generated.
func (p *triggerDefinitionActionProps) setFilter(filter *types.TriggerDefinitionFilter) *triggerDefinitionActionProps {
	p.filter = filter
	return p
}

// Serialize converts triggerDefinitionActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p triggerDefinitionActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.triggerDefinition != nil {
		m.Set("triggerDefinition.handle", p.triggerDefinition.Handle, true)
		m.Set("triggerDefinition.ID", p.triggerDefinition.ID, true)
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
func (p triggerDefinitionActionProps) Format(in string, err error) string {
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

	if p.triggerDefinition != nil {
		// replacement for "{{triggerDefinition}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{triggerDefinition}}",
			fns(
				p.triggerDefinition.Handle,
				p.triggerDefinition.ID,
			),
		)
		pairs = append(pairs, "{{triggerDefinition.handle}}", fns(p.triggerDefinition.Handle))
		pairs = append(pairs, "{{triggerDefinition.ID}}", fns(p.triggerDefinition.ID))
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
func (a *triggerDefinitionAction) String() string {
	var props = &triggerDefinitionActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *triggerDefinitionAction) ToAction() *actionlog.Action {
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

// TriggerDefinitionActionSearch returns "automation:trigger-definition.search" action
//
// This function is auto-generated.
func TriggerDefinitionActionSearch(props ...*triggerDefinitionActionProps) *triggerDefinitionAction {
	a := &triggerDefinitionAction{
		timestamp: time.Now(),
		resource:  "automation:trigger-definition",
		action:    "search",
		log:       "searched for matching trigger definitions",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerDefinitionActionLookup returns "automation:trigger-definition.lookup" action
//
// This function is auto-generated.
func TriggerDefinitionActionLookup(props ...*triggerDefinitionActionProps) *triggerDefinitionAction {
	a := &triggerDefinitionAction{
		timestamp: time.Now(),
		resource:  "automation:trigger-definition",
		action:    "lookup",
		log:       "looked-up for a {{triggerDefinition}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerDefinitionActionCreate returns "automation:trigger-definition.create" action
//
// This function is auto-generated.
func TriggerDefinitionActionCreate(props ...*triggerDefinitionActionProps) *triggerDefinitionAction {
	a := &triggerDefinitionAction{
		timestamp: time.Now(),
		resource:  "automation:trigger-definition",
		action:    "create",
		log:       "created {{triggerDefinition}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerDefinitionActionUpdate returns "automation:trigger-definition.update" action
//
// This function is auto-generated.
func TriggerDefinitionActionUpdate(props ...*triggerDefinitionActionProps) *triggerDefinitionAction {
	a := &triggerDefinitionAction{
		timestamp: time.Now(),
		resource:  "automation:trigger-definition",
		action:    "update",
		log:       "updated {{triggerDefinition}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerDefinitionActionDelete returns "automation:trigger-definition.delete" action
//
// This function is auto-generated.
func TriggerDefinitionActionDelete(props ...*triggerDefinitionActionProps) *triggerDefinitionAction {
	a := &triggerDefinitionAction{
		timestamp: time.Now(),
		resource:  "automation:trigger-definition",
		action:    "delete",
		log:       "deleted {{triggerDefinition}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerDefinitionActionUndelete returns "automation:trigger-definition.undelete" action
//
// This function is auto-generated.
func TriggerDefinitionActionUndelete(props ...*triggerDefinitionActionProps) *triggerDefinitionAction {
	a := &triggerDefinitionAction{
		timestamp: time.Now(),
		resource:  "automation:trigger-definition",
		action:    "undelete",
		log:       "undeleted {{triggerDefinition}}",
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

// TriggerDefinitionErrGeneric returns "automation:trigger-definition.generic" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrGeneric(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "automation:trigger-definition"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerDefinitionLogMetaKey{}, "{err}"),
		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrNotFound returns "automation:trigger-definition.notFound" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrNotFound(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("trigger definition not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "automation:trigger-definition"),

		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrInvalidID returns "automation:trigger-definition.invalidID" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrInvalidID(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "automation:trigger-definition"),

		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrInvalidHandle returns "automation:trigger-definition.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrInvalidHandle(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "automation:trigger-definition"),

		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrStaleData returns "automation:trigger-definition.staleData" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrStaleData(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "automation:trigger-definition"),

		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrHandleNotUnique returns "automation:trigger-definition.handleNotUnique" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrHandleNotUnique(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("trigger definition handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "automation:trigger-definition"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerDefinitionLogMetaKey{}, "duplicate handle used for trigger definition ({{triggerDefinition}})"),
		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrNotAllowedToRead returns "automation:trigger-definition.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrNotAllowedToRead(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this trigger definition", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "automation:trigger-definition"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerDefinitionLogMetaKey{}, "failed to read {{triggerDefinition.handle}}; insufficient permissions"),
		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrNotAllowedToSearch returns "automation:trigger-definition.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrNotAllowedToSearch(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list trigger definitions", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "automation:trigger-definition"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerDefinitionLogMetaKey{}, "failed to search or list trigger definitions; insufficient permissions"),
		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrNotAllowedToCreate returns "automation:trigger-definition.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrNotAllowedToCreate(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create trigger definitions", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "automation:trigger-definition"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerDefinitionLogMetaKey{}, "failed to create trigger definition; insufficient permissions"),
		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrNotAllowedToUpdate returns "automation:trigger-definition.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrNotAllowedToUpdate(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this trigger definition", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "automation:trigger-definition"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerDefinitionLogMetaKey{}, "failed to update {{triggerDefinition}}; insufficient permissions"),
		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrNotAllowedToDelete returns "automation:trigger-definition.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrNotAllowedToDelete(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this trigger definition", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "automation:trigger-definition"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerDefinitionLogMetaKey{}, "failed to delete {{triggerDefinition}}; insufficient permissions"),
		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerDefinitionErrNotAllowedToUndelete returns "automation:trigger-definition.notAllowedToUndelete" as *errors.Error
//
// This function is auto-generated.
func TriggerDefinitionErrNotAllowedToUndelete(mm ...*triggerDefinitionActionProps) *errors.Error {
	var p = &triggerDefinitionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this trigger definition", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "automation:trigger-definition"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerDefinitionLogMetaKey{}, "failed to undelete {{triggerDefinition}}; insufficient permissions"),
		errors.Meta(triggerDefinitionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "automation"),
		errors.Meta(locale.ErrorMetaKey{}, "trigger-definition.errors.notAllowedToUndelete"),

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
func (svc triggerDefinition) recordAction(ctx context.Context, props *triggerDefinitionActionProps, actionFn func(...*triggerDefinitionActionProps) *triggerDefinitionAction, err error) error {
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
		a.Description = props.Format(m.AsString(triggerDefinitionLogMetaKey{}), err)

		if p, has := m[triggerDefinitionPropsMetaKey{}]; has {
			a.Meta = p.(*triggerDefinitionActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
