package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/reminder_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/system/types"
	"strings"
	"time"
)

type (
	reminderActionProps struct {
		reminder *types.Reminder
		new      *types.Reminder
		updated  *types.Reminder
		filter   *types.ReminderFilter
	}

	reminderAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *reminderActionProps
	}

	reminderLogMetaKey   struct{}
	reminderPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setReminder updates reminderActionProps's reminder
//
// This function is auto-generated.
func (p *reminderActionProps) setReminder(reminder *types.Reminder) *reminderActionProps {
	p.reminder = reminder
	return p
}

// setNew updates reminderActionProps's new
//
// This function is auto-generated.
func (p *reminderActionProps) setNew(new *types.Reminder) *reminderActionProps {
	p.new = new
	return p
}

// setUpdated updates reminderActionProps's updated
//
// This function is auto-generated.
func (p *reminderActionProps) setUpdated(updated *types.Reminder) *reminderActionProps {
	p.updated = updated
	return p
}

// setFilter updates reminderActionProps's filter
//
// This function is auto-generated.
func (p *reminderActionProps) setFilter(filter *types.ReminderFilter) *reminderActionProps {
	p.filter = filter
	return p
}

// Serialize converts reminderActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p reminderActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.reminder != nil {
		m.Set("reminder.resource", p.reminder.Resource, true)
		m.Set("reminder.ID", p.reminder.ID, true)
		m.Set("reminder.assignedTo", p.reminder.AssignedTo, true)
		m.Set("reminder.assignedBy", p.reminder.AssignedBy, true)
		m.Set("reminder.remindAt", p.reminder.RemindAt, true)
	}
	if p.new != nil {
		m.Set("new.resource", p.new.Resource, true)
		m.Set("new.ID", p.new.ID, true)
		m.Set("new.assignedTo", p.new.AssignedTo, true)
		m.Set("new.assignedBy", p.new.AssignedBy, true)
		m.Set("new.remindAt", p.new.RemindAt, true)
	}
	if p.updated != nil {
		m.Set("updated.resource", p.updated.Resource, true)
		m.Set("updated.ID", p.updated.ID, true)
		m.Set("updated.assignedTo", p.updated.AssignedTo, true)
		m.Set("updated.assignedBy", p.updated.AssignedBy, true)
		m.Set("updated.remindAt", p.updated.RemindAt, true)
	}
	if p.filter != nil {
		m.Set("filter.resource", p.filter.Resource, true)
		m.Set("filter.reminderID", p.filter.ReminderID, true)
		m.Set("filter.assignedTo", p.filter.AssignedTo, true)
		m.Set("filter.scheduledFrom", p.filter.ScheduledFrom, true)
		m.Set("filter.scheduledUntil", p.filter.ScheduledUntil, true)
		m.Set("filter.excludeDismissed", p.filter.ExcludeDismissed, true)
		m.Set("filter.scheduledOnly", p.filter.ScheduledOnly, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p reminderActionProps) Format(in string, err error) string {
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

	if p.reminder != nil {
		// replacement for "{{reminder}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{reminder}}",
			fns(
				p.reminder.Resource,
				p.reminder.ID,
				p.reminder.AssignedTo,
				p.reminder.AssignedBy,
				p.reminder.RemindAt,
			),
		)
		pairs = append(pairs, "{{reminder.resource}}", fns(p.reminder.Resource))
		pairs = append(pairs, "{{reminder.ID}}", fns(p.reminder.ID))
		pairs = append(pairs, "{{reminder.assignedTo}}", fns(p.reminder.AssignedTo))
		pairs = append(pairs, "{{reminder.assignedBy}}", fns(p.reminder.AssignedBy))
		pairs = append(pairs, "{{reminder.remindAt}}", fns(p.reminder.RemindAt))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Resource,
				p.new.ID,
				p.new.AssignedTo,
				p.new.AssignedBy,
				p.new.RemindAt,
			),
		)
		pairs = append(pairs, "{{new.resource}}", fns(p.new.Resource))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
		pairs = append(pairs, "{{new.assignedTo}}", fns(p.new.AssignedTo))
		pairs = append(pairs, "{{new.assignedBy}}", fns(p.new.AssignedBy))
		pairs = append(pairs, "{{new.remindAt}}", fns(p.new.RemindAt))
	}

	if p.updated != nil {
		// replacement for "{{updated}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{updated}}",
			fns(
				p.updated.Resource,
				p.updated.ID,
				p.updated.AssignedTo,
				p.updated.AssignedBy,
				p.updated.RemindAt,
			),
		)
		pairs = append(pairs, "{{updated.resource}}", fns(p.updated.Resource))
		pairs = append(pairs, "{{updated.ID}}", fns(p.updated.ID))
		pairs = append(pairs, "{{updated.assignedTo}}", fns(p.updated.AssignedTo))
		pairs = append(pairs, "{{updated.assignedBy}}", fns(p.updated.AssignedBy))
		pairs = append(pairs, "{{updated.remindAt}}", fns(p.updated.RemindAt))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Resource,
				p.filter.ReminderID,
				p.filter.AssignedTo,
				p.filter.ScheduledFrom,
				p.filter.ScheduledUntil,
				p.filter.ExcludeDismissed,
				p.filter.ScheduledOnly,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{{filter.resource}}", fns(p.filter.Resource))
		pairs = append(pairs, "{{filter.reminderID}}", fns(p.filter.ReminderID))
		pairs = append(pairs, "{{filter.assignedTo}}", fns(p.filter.AssignedTo))
		pairs = append(pairs, "{{filter.scheduledFrom}}", fns(p.filter.ScheduledFrom))
		pairs = append(pairs, "{{filter.scheduledUntil}}", fns(p.filter.ScheduledUntil))
		pairs = append(pairs, "{{filter.excludeDismissed}}", fns(p.filter.ExcludeDismissed))
		pairs = append(pairs, "{{filter.scheduledOnly}}", fns(p.filter.ScheduledOnly))
		pairs = append(pairs, "{{filter.sort}}", fns(p.filter.Sort))
	}
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
func (a *reminderAction) String() string {
	var props = &reminderActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *reminderAction) ToAction() *actionlog.Action {
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

// ReminderActionSearch returns "system:reminder.search" action
//
// This function is auto-generated.
func ReminderActionSearch(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "search",
		log:       "searched for reminders",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionLookup returns "system:reminder.lookup" action
//
// This function is auto-generated.
func ReminderActionLookup(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "lookup",
		log:       "looked-up for a {{reminder}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionCreate returns "system:reminder.create" action
//
// This function is auto-generated.
func ReminderActionCreate(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "create",
		log:       "created {{reminder}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionUpdate returns "system:reminder.update" action
//
// This function is auto-generated.
func ReminderActionUpdate(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "update",
		log:       "updated {{reminder}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionDelete returns "system:reminder.delete" action
//
// This function is auto-generated.
func ReminderActionDelete(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "delete",
		log:       "deleted {{reminder}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionDismiss returns "system:reminder.dismiss" action
//
// This function is auto-generated.
func ReminderActionDismiss(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "dismiss",
		log:       "deleted {{reminder}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionSnooze returns "system:reminder.snooze" action
//
// This function is auto-generated.
func ReminderActionSnooze(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "snooze",
		log:       "deleted {{reminder}}",
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

// ReminderErrGeneric returns "system:reminder.generic" as *errors.Error
//
// This function is auto-generated.
func ReminderErrGeneric(mm ...*reminderActionProps) *errors.Error {
	var p = &reminderActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:reminder"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(reminderLogMetaKey{}, "{err}"),
		errors.Meta(reminderPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "reminder.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReminderErrNotFound returns "system:reminder.notFound" as *errors.Error
//
// This function is auto-generated.
func ReminderErrNotFound(mm ...*reminderActionProps) *errors.Error {
	var p = &reminderActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("reminder not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:reminder"),

		errors.Meta(reminderPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "reminder.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReminderErrInvalidID returns "system:reminder.invalidID" as *errors.Error
//
// This function is auto-generated.
func ReminderErrInvalidID(mm ...*reminderActionProps) *errors.Error {
	var p = &reminderActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:reminder"),

		errors.Meta(reminderPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "reminder.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReminderErrNotAllowedToAssign returns "system:reminder.notAllowedToAssign" as *errors.Error
//
// This function is auto-generated.
func ReminderErrNotAllowedToAssign(mm ...*reminderActionProps) *errors.Error {
	var p = &reminderActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to assign reminders to other users", nil),

		errors.Meta("type", "notAllowedToAssign"),
		errors.Meta("resource", "system:reminder"),

		errors.Meta(reminderPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "reminder.errors.notAllowedToAssign"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReminderErrNotAllowedToDismiss returns "system:reminder.notAllowedToDismiss" as *errors.Error
//
// This function is auto-generated.
func ReminderErrNotAllowedToDismiss(mm ...*reminderActionProps) *errors.Error {
	var p = &reminderActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to dismiss reminders of other users", nil),

		errors.Meta("type", "notAllowedToDismiss"),
		errors.Meta("resource", "system:reminder"),

		errors.Meta(reminderPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "reminder.errors.notAllowedToDismiss"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReminderErrNotAllowedToRead returns "system:reminder.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ReminderErrNotAllowedToRead(mm ...*reminderActionProps) *errors.Error {
	var p = &reminderActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read reminders of other users", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:reminder"),

		errors.Meta(reminderPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "reminder.errors.notAllowedToRead"),

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
func (svc reminder) recordAction(ctx context.Context, props *reminderActionProps, actionFn func(...*reminderActionProps) *reminderAction, err error) error {
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
		a.Description = props.Format(m.AsString(reminderLogMetaKey{}), err)

		if p, has := m[reminderPropsMetaKey{}]; has {
			a.Meta = p.(*reminderActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
