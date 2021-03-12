package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/service/trigger_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"strings"
	"time"
)

type (
	triggerActionProps struct {
		trigger *types.Trigger
		new     *types.Trigger
		update  *types.Trigger
		filter  *types.TriggerFilter
	}

	triggerAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *triggerActionProps
	}

	triggerLogMetaKey   struct{}
	triggerPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setTrigger updates triggerActionProps's trigger
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *triggerActionProps) setTrigger(trigger *types.Trigger) *triggerActionProps {
	p.trigger = trigger
	return p
}

// setNew updates triggerActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *triggerActionProps) setNew(new *types.Trigger) *triggerActionProps {
	p.new = new
	return p
}

// setUpdate updates triggerActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *triggerActionProps) setUpdate(update *types.Trigger) *triggerActionProps {
	p.update = update
	return p
}

// setFilter updates triggerActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *triggerActionProps) setFilter(filter *types.TriggerFilter) *triggerActionProps {
	p.filter = filter
	return p
}

// Serialize converts triggerActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p triggerActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.trigger != nil {
		m.Set("trigger.ID", p.trigger.ID, true)
	}
	if p.new != nil {
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
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
func (p triggerActionProps) Format(in string, err error) string {
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

	if p.trigger != nil {
		// replacement for "{trigger}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{trigger}",
			fns(
				p.trigger.ID,
			),
		)
		pairs = append(pairs, "{trigger.ID}", fns(p.trigger.ID))
	}

	if p.new != nil {
		// replacement for "{new}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{new}",
			fns(
				p.new.ID,
			),
		)
		pairs = append(pairs, "{new.ID}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{update}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{update}",
			fns(
				p.update.ID,
			),
		)
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
func (a *triggerAction) String() string {
	var props = &triggerActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *triggerAction) ToAction() *actionlog.Action {
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

// TriggerActionSearch returns "automation:trigger.search" action
//
// This function is auto-generated.
//
func TriggerActionSearch(props ...*triggerActionProps) *triggerAction {
	a := &triggerAction{
		timestamp: time.Now(),
		resource:  "automation:trigger",
		action:    "search",
		log:       "searched for matching triggers",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerActionLookup returns "automation:trigger.lookup" action
//
// This function is auto-generated.
//
func TriggerActionLookup(props ...*triggerActionProps) *triggerAction {
	a := &triggerAction{
		timestamp: time.Now(),
		resource:  "automation:trigger",
		action:    "lookup",
		log:       "looked-up for a {trigger}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerActionCreate returns "automation:trigger.create" action
//
// This function is auto-generated.
//
func TriggerActionCreate(props ...*triggerActionProps) *triggerAction {
	a := &triggerAction{
		timestamp: time.Now(),
		resource:  "automation:trigger",
		action:    "create",
		log:       "created {trigger}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerActionUpdate returns "automation:trigger.update" action
//
// This function is auto-generated.
//
func TriggerActionUpdate(props ...*triggerActionProps) *triggerAction {
	a := &triggerAction{
		timestamp: time.Now(),
		resource:  "automation:trigger",
		action:    "update",
		log:       "updated {trigger}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerActionDelete returns "automation:trigger.delete" action
//
// This function is auto-generated.
//
func TriggerActionDelete(props ...*triggerActionProps) *triggerAction {
	a := &triggerAction{
		timestamp: time.Now(),
		resource:  "automation:trigger",
		action:    "delete",
		log:       "deleted {trigger}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TriggerActionUndelete returns "automation:trigger.undelete" action
//
// This function is auto-generated.
//
func TriggerActionUndelete(props ...*triggerActionProps) *triggerAction {
	a := &triggerAction{
		timestamp: time.Now(),
		resource:  "automation:trigger",
		action:    "undelete",
		log:       "undeleted {trigger}",
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

// TriggerErrGeneric returns "automation:trigger.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrGeneric(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "automation:trigger"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerLogMetaKey{}, "{err}"),
		errors.Meta(triggerPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerErrNotFound returns "automation:trigger.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrNotFound(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("trigger not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "automation:trigger"),

		errors.Meta(triggerPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerErrInvalidID returns "automation:trigger.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrInvalidID(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "automation:trigger"),

		errors.Meta(triggerPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerErrStaleData returns "automation:trigger.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrStaleData(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "automation:trigger"),

		errors.Meta(triggerPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerErrNotAllowedToRead returns "automation:trigger.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrNotAllowedToRead(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this trigger", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "automation:trigger"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerLogMetaKey{}, "failed to read {trigger.ID}; insufficient permissions"),
		errors.Meta(triggerPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerErrNotAllowedToSearch returns "automation:trigger.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrNotAllowedToSearch(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search triggers", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "automation:trigger"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerLogMetaKey{}, "failed to list trigger; insufficient permissions"),
		errors.Meta(triggerPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerErrNotAllowedToCreate returns "automation:trigger.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrNotAllowedToCreate(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create triggers", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "automation:trigger"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerLogMetaKey{}, "failed to create trigger; insufficient permissions"),
		errors.Meta(triggerPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerErrNotAllowedToUpdate returns "automation:trigger.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrNotAllowedToUpdate(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this trigger", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "automation:trigger"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerLogMetaKey{}, "failed to update {trigger.ID}; insufficient permissions"),
		errors.Meta(triggerPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerErrNotAllowedToDelete returns "automation:trigger.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrNotAllowedToDelete(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this trigger", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "automation:trigger"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerLogMetaKey{}, "failed to delete {trigger.ID}; insufficient permissions"),
		errors.Meta(triggerPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TriggerErrNotAllowedToUndelete returns "automation:trigger.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func TriggerErrNotAllowedToUndelete(mm ...*triggerActionProps) *errors.Error {
	var p = &triggerActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this trigger", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "automation:trigger"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(triggerLogMetaKey{}, "failed to undelete {trigger.ID}; insufficient permissions"),
		errors.Meta(triggerPropsMetaKey{}, p),

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
func (svc trigger) recordAction(ctx context.Context, props *triggerActionProps, actionFn func(...*triggerActionProps) *triggerAction, err error) error {
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
		a.Description = props.Format(m.AsString(triggerLogMetaKey{}), err)

		if p, has := m[triggerPropsMetaKey{}]; has {
			a.Meta = p.(*triggerActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
