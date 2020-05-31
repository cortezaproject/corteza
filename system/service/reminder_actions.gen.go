package service

// This file is auto-generated from system/service/reminder_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/system/types"
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

	reminderError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *reminderActionProps
	}
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
// Allows method chaining
//
// This function is auto-generated.
//
func (p *reminderActionProps) setReminder(reminder *types.Reminder) *reminderActionProps {
	p.reminder = reminder
	return p
}

// setNew updates reminderActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *reminderActionProps) setNew(new *types.Reminder) *reminderActionProps {
	p.new = new
	return p
}

// setUpdated updates reminderActionProps's updated
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *reminderActionProps) setUpdated(updated *types.Reminder) *reminderActionProps {
	p.updated = updated
	return p
}

// setFilter updates reminderActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *reminderActionProps) setFilter(filter *types.ReminderFilter) *reminderActionProps {
	p.filter = filter
	return p
}

// serialize converts reminderActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p reminderActionProps) serialize() actionlog.Meta {
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
//
func (p reminderActionProps) tr(in string, err error) string {
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
		for {
			// Unwrap errors
			ue := errors.Unwrap(err)
			if ue == nil {
				break
			}

			err = ue
		}

		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}

	if p.reminder != nil {
		// replacement for "{reminder}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{reminder}",
			fns(
				p.reminder.Resource,
				p.reminder.ID,
				p.reminder.AssignedTo,
				p.reminder.AssignedBy,
				p.reminder.RemindAt,
			),
		)
		pairs = append(pairs, "{reminder.resource}", fns(p.reminder.Resource))
		pairs = append(pairs, "{reminder.ID}", fns(p.reminder.ID))
		pairs = append(pairs, "{reminder.assignedTo}", fns(p.reminder.AssignedTo))
		pairs = append(pairs, "{reminder.assignedBy}", fns(p.reminder.AssignedBy))
		pairs = append(pairs, "{reminder.remindAt}", fns(p.reminder.RemindAt))
	}

	if p.new != nil {
		// replacement for "{new}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{new}",
			fns(
				p.new.Resource,
				p.new.ID,
				p.new.AssignedTo,
				p.new.AssignedBy,
				p.new.RemindAt,
			),
		)
		pairs = append(pairs, "{new.resource}", fns(p.new.Resource))
		pairs = append(pairs, "{new.ID}", fns(p.new.ID))
		pairs = append(pairs, "{new.assignedTo}", fns(p.new.AssignedTo))
		pairs = append(pairs, "{new.assignedBy}", fns(p.new.AssignedBy))
		pairs = append(pairs, "{new.remindAt}", fns(p.new.RemindAt))
	}

	if p.updated != nil {
		// replacement for "{updated}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{updated}",
			fns(
				p.updated.Resource,
				p.updated.ID,
				p.updated.AssignedTo,
				p.updated.AssignedBy,
				p.updated.RemindAt,
			),
		)
		pairs = append(pairs, "{updated.resource}", fns(p.updated.Resource))
		pairs = append(pairs, "{updated.ID}", fns(p.updated.ID))
		pairs = append(pairs, "{updated.assignedTo}", fns(p.updated.AssignedTo))
		pairs = append(pairs, "{updated.assignedBy}", fns(p.updated.AssignedBy))
		pairs = append(pairs, "{updated.remindAt}", fns(p.updated.RemindAt))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
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
		pairs = append(pairs, "{filter.resource}", fns(p.filter.Resource))
		pairs = append(pairs, "{filter.reminderID}", fns(p.filter.ReminderID))
		pairs = append(pairs, "{filter.assignedTo}", fns(p.filter.AssignedTo))
		pairs = append(pairs, "{filter.scheduledFrom}", fns(p.filter.ScheduledFrom))
		pairs = append(pairs, "{filter.scheduledUntil}", fns(p.filter.ScheduledUntil))
		pairs = append(pairs, "{filter.excludeDismissed}", fns(p.filter.ExcludeDismissed))
		pairs = append(pairs, "{filter.scheduledOnly}", fns(p.filter.ScheduledOnly))
		pairs = append(pairs, "{filter.sort}", fns(p.filter.Sort))
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
func (a *reminderAction) String() string {
	var props = &reminderActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *reminderAction) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error methods

// String returns loggable description as string
//
// It falls back to message if log is not set
//
// This function is auto-generated.
//
func (e *reminderError) String() string {
	var props = &reminderActionProps{}

	if e.props != nil {
		props = e.props
	}

	if e.wrap != nil && !strings.Contains(e.log, "{err}") {
		// Suffix error log with {err} to ensure
		// we log the cause for this error
		e.log += ": {err}"
	}

	return props.tr(e.log, e.wrap)
}

// Error satisfies
//
// This function is auto-generated.
//
func (e *reminderError) Error() string {
	var props = &reminderActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *reminderError) Is(Resource error) bool {
	t, ok := Resource.(*reminderError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps reminderError around another error
//
// This function is auto-generated.
//
func (e *reminderError) Wrap(err error) *reminderError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *reminderError) Unwrap() error {
	return e.wrap
}

func (e *reminderError) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Error:       e.Error(),
		Meta:        e.props.serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// ReminderActionSearch returns "system:reminder.search" error
//
// This function is auto-generated.
//
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

// ReminderActionLookup returns "system:reminder.lookup" error
//
// This function is auto-generated.
//
func ReminderActionLookup(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "lookup",
		log:       "looked-up for a {reminder}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionCreate returns "system:reminder.create" error
//
// This function is auto-generated.
//
func ReminderActionCreate(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "create",
		log:       "created {reminder}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionUpdate returns "system:reminder.update" error
//
// This function is auto-generated.
//
func ReminderActionUpdate(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "update",
		log:       "updated {reminder}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionDelete returns "system:reminder.delete" error
//
// This function is auto-generated.
//
func ReminderActionDelete(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "delete",
		log:       "deleted {reminder}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionDismiss returns "system:reminder.dismiss" error
//
// This function is auto-generated.
//
func ReminderActionDismiss(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "dismiss",
		log:       "deleted {reminder}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReminderActionSnooze returns "system:reminder.snooze" error
//
// This function is auto-generated.
//
func ReminderActionSnooze(props ...*reminderActionProps) *reminderAction {
	a := &reminderAction{
		timestamp: time.Now(),
		resource:  "system:reminder",
		action:    "snooze",
		log:       "deleted {reminder}",
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

// ReminderErrGeneric returns "system:reminder.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ReminderErrGeneric(props ...*reminderActionProps) *reminderError {
	var e = &reminderError{
		timestamp: time.Now(),
		resource:  "system:reminder",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *reminderActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ReminderErrNotFound returns "system:reminder.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ReminderErrNotFound(props ...*reminderActionProps) *reminderError {
	var e = &reminderError{
		timestamp: time.Now(),
		resource:  "system:reminder",
		error:     "notFound",
		action:    "error",
		message:   "reminder not found",
		log:       "reminder not found",
		severity:  actionlog.Warning,
		props: func() *reminderActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ReminderErrInvalidID returns "system:reminder.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ReminderErrInvalidID(props ...*reminderActionProps) *reminderError {
	var e = &reminderError{
		timestamp: time.Now(),
		resource:  "system:reminder",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *reminderActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// ReminderErrNotAllowedToAssign returns "system:reminder.notAllowedToAssign" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ReminderErrNotAllowedToAssign(props ...*reminderActionProps) *reminderError {
	var e = &reminderError{
		timestamp: time.Now(),
		resource:  "system:reminder",
		error:     "notAllowedToAssign",
		action:    "error",
		message:   "not allowed to assign reminders to other users",
		log:       "not allowed to assign reminders to other users",
		severity:  actionlog.Error,
		props: func() *reminderActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// context is used to enrich audit log entry with current user info, request ID, IP address...
// props are collected action/error properties
// action (optional) fn will be used to construct reminderAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc reminder) recordAction(ctx context.Context, props *reminderActionProps, action func(...*reminderActionProps) *reminderAction, err error) error {
	var (
		ok bool

		// Return error
		retError *reminderError

		// Recorder error
		recError *reminderError
	)

	if err != nil {
		if retError, ok = err.(*reminderError); !ok {
			// got non-reminder error, wrap it with ReminderErrGeneric
			retError = ReminderErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use ReminderErrGeneric for recording too
			// because it can hold more info
			recError = retError
		} else if retError != nil {
			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}
			// start with copy of return error for recording
			// this will be updated with tha root cause as we try and
			// unwrap the error
			recError = retError

			// find the original recError for this error
			// for the purpose of logging
			var unwrappedError error = retError
			for {
				if unwrappedError = errors.Unwrap(unwrappedError); unwrappedError == nil {
					// nothing wrapped
					break
				}

				// update recError ONLY of wrapped error is of type reminderError
				if unwrappedSinkError, ok := unwrappedError.(*reminderError); ok {
					recError = unwrappedSinkError
				}
			}

			if retError.props == nil {
				// set props on returning error if empty
				retError.props = props
			}

			if recError.props == nil {
				// set props on recording error if empty
				recError.props = props
			}
		}
	}

	if svc.actionlog != nil {
		if retError != nil {
			// failed action, log error
			svc.actionlog.Record(ctx, recError)
		} else if action != nil {
			// successful
			svc.actionlog.Record(ctx, action(props))
		}
	}

	if err == nil {
		// retError not an interface and that WILL (!!) cause issues
		// with nil check (== nil) when it is not explicitly returned
		return nil
	}

	return retError
}
