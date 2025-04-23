package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/notification_actions.yaml

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
	notificationActionProps struct {
		notification *types.Notification
		new          *types.Notification
		updated      *types.Notification
		filter       *types.NotificationFilter
	}

	notificationAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *notificationActionProps
	}

	notificationLogMetaKey   struct{}
	notificationPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setNotification updates notificationActionProps's notification
//
// This function is auto-generated.
func (p *notificationActionProps) setNotification(notification *types.Notification) *notificationActionProps {
	p.notification = notification
	return p
}

// setNew updates notificationActionProps's new
//
// This function is auto-generated.
func (p *notificationActionProps) setNew(new *types.Notification) *notificationActionProps {
	p.new = new
	return p
}

// setUpdated updates notificationActionProps's updated
//
// This function is auto-generated.
func (p *notificationActionProps) setUpdated(updated *types.Notification) *notificationActionProps {
	p.updated = updated
	return p
}

// setFilter updates notificationActionProps's filter
//
// This function is auto-generated.
func (p *notificationActionProps) setFilter(filter *types.NotificationFilter) *notificationActionProps {
	p.filter = filter
	return p
}

// Serialize converts notificationActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p notificationActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.notification != nil {
		m.Set("notification.ID", p.notification.ID, true)
		m.Set("notification.recipient", p.notification.Recipient, true)
		m.Set("notification.createdBy", p.notification.CreatedBy, true)
		m.Set("notification.kind", p.notification.Kind, true)
		m.Set("notification.config", p.notification.Config, true)
	}
	if p.new != nil {
		m.Set("new.ID", p.new.ID, true)
		m.Set("new.recipient", p.new.Recipient, true)
		m.Set("new.createdBy", p.new.CreatedBy, true)
		m.Set("new.kind", p.new.Kind, true)
		m.Set("new.config", p.new.Config, true)
	}
	if p.updated != nil {
		m.Set("updated.ID", p.updated.ID, true)
		m.Set("updated.recipient", p.updated.Recipient, true)
		m.Set("updated.createdBy", p.updated.CreatedBy, true)
		m.Set("updated.kind", p.updated.Kind, true)
		m.Set("updated.config", p.updated.Config, true)
	}
	if p.filter != nil {
		m.Set("filter.notificationID", p.filter.NotificationID, true)
		m.Set("filter.recipient", p.filter.Recipient, true)
		m.Set("filter.kind", p.filter.Kind, true)
		m.Set("filter.read", p.filter.Read, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p notificationActionProps) Format(in string, err error) string {
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

	if p.notification != nil {
		// replacement for "{{notification}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{notification}}",
			fns(
				p.notification.ID,
				p.notification.Recipient,
				p.notification.CreatedBy,
				p.notification.Kind,
				p.notification.Config,
			),
		)
		pairs = append(pairs, "{{notification.ID}}", fns(p.notification.ID))
		pairs = append(pairs, "{{notification.recipient}}", fns(p.notification.Recipient))
		pairs = append(pairs, "{{notification.createdBy}}", fns(p.notification.CreatedBy))
		pairs = append(pairs, "{{notification.kind}}", fns(p.notification.Kind))
		pairs = append(pairs, "{{notification.config}}", fns(p.notification.Config))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.ID,
				p.new.Recipient,
				p.new.CreatedBy,
				p.new.Kind,
				p.new.Config,
			),
		)
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
		pairs = append(pairs, "{{new.recipient}}", fns(p.new.Recipient))
		pairs = append(pairs, "{{new.createdBy}}", fns(p.new.CreatedBy))
		pairs = append(pairs, "{{new.kind}}", fns(p.new.Kind))
		pairs = append(pairs, "{{new.config}}", fns(p.new.Config))
	}

	if p.updated != nil {
		// replacement for "{{updated}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{updated}}",
			fns(
				p.updated.ID,
				p.updated.Recipient,
				p.updated.CreatedBy,
				p.updated.Kind,
				p.updated.Config,
			),
		)
		pairs = append(pairs, "{{updated.ID}}", fns(p.updated.ID))
		pairs = append(pairs, "{{updated.recipient}}", fns(p.updated.Recipient))
		pairs = append(pairs, "{{updated.createdBy}}", fns(p.updated.CreatedBy))
		pairs = append(pairs, "{{updated.kind}}", fns(p.updated.Kind))
		pairs = append(pairs, "{{updated.config}}", fns(p.updated.Config))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.NotificationID,
				p.filter.Recipient,
				p.filter.Kind,
				p.filter.Read,
				p.filter.Deleted,
			),
		)
		pairs = append(pairs, "{{filter.notificationID}}", fns(p.filter.NotificationID))
		pairs = append(pairs, "{{filter.recipient}}", fns(p.filter.Recipient))
		pairs = append(pairs, "{{filter.kind}}", fns(p.filter.Kind))
		pairs = append(pairs, "{{filter.read}}", fns(p.filter.Read))
		pairs = append(pairs, "{{filter.deleted}}", fns(p.filter.Deleted))
	}
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
func (a *notificationAction) String() string {
	var props = &notificationActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *notificationAction) ToAction() *actionlog.Action {
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

// NotificationActionSearch returns "system:notification.search" action
//
// This function is auto-generated.
func NotificationActionSearch(props ...*notificationActionProps) *notificationAction {
	a := &notificationAction{
		timestamp: time.Now(),
		resource:  "system:notification",
		action:    "search",
		log:       "searched for notifications",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NotificationActionLookup returns "system:notification.lookup" action
//
// This function is auto-generated.
func NotificationActionLookup(props ...*notificationActionProps) *notificationAction {
	a := &notificationAction{
		timestamp: time.Now(),
		resource:  "system:notification",
		action:    "lookup",
		log:       "looked-up for a {{notification}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NotificationActionCreate returns "system:notification.create" action
//
// This function is auto-generated.
func NotificationActionCreate(props ...*notificationActionProps) *notificationAction {
	a := &notificationAction{
		timestamp: time.Now(),
		resource:  "system:notification",
		action:    "create",
		log:       "created {{notification}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NotificationActionUpdate returns "system:notification.update" action
//
// This function is auto-generated.
func NotificationActionUpdate(props ...*notificationActionProps) *notificationAction {
	a := &notificationAction{
		timestamp: time.Now(),
		resource:  "system:notification",
		action:    "update",
		log:       "updated {{notification}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NotificationActionDelete returns "system:notification.delete" action
//
// This function is auto-generated.
func NotificationActionDelete(props ...*notificationActionProps) *notificationAction {
	a := &notificationAction{
		timestamp: time.Now(),
		resource:  "system:notification",
		action:    "delete",
		log:       "deleted {{notification}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NotificationActionMarkAsRead returns "system:notification.markAsRead" action
//
// This function is auto-generated.
func NotificationActionMarkAsRead(props ...*notificationActionProps) *notificationAction {
	a := &notificationAction{
		timestamp: time.Now(),
		resource:  "system:notification",
		action:    "markAsRead",
		log:       "marked {{notification}} as read",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NotificationActionMarkAllAsRead returns "system:notification.markAllAsRead" action
//
// This function is auto-generated.
func NotificationActionMarkAllAsRead(props ...*notificationActionProps) *notificationAction {
	a := &notificationAction{
		timestamp: time.Now(),
		resource:  "system:notification",
		action:    "markAllAsRead",
		log:       "marked all notifications as read for current user",
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

// NotificationErrGeneric returns "system:notification.generic" as *errors.Error
//
// This function is auto-generated.
func NotificationErrGeneric(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:notification"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(notificationLogMetaKey{}, "{err}"),
		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrNotFound returns "system:notification.notFound" as *errors.Error
//
// This function is auto-generated.
func NotificationErrNotFound(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("notification not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrInvalidID returns "system:notification.invalidID" as *errors.Error
//
// This function is auto-generated.
func NotificationErrInvalidID(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrNotAllowedToRead returns "system:notification.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func NotificationErrNotAllowedToRead(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read notifications of other users", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrNotAllowedToCreate returns "system:notification.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func NotificationErrNotAllowedToCreate(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create notifications", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrNotAllowedToUpdate returns "system:notification.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func NotificationErrNotAllowedToUpdate(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update notifications", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrNotAllowedToDelete returns "system:notification.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func NotificationErrNotAllowedToDelete(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete notifications", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrNotAllowedToAssign returns "system:notification.notAllowedToAssign" as *errors.Error
//
// This function is auto-generated.
func NotificationErrNotAllowedToAssign(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to assign notifications to other users", nil),

		errors.Meta("type", "notAllowedToAssign"),
		errors.Meta("resource", "system:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.notAllowedToAssign"),

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
func (svc notification) recordAction(ctx context.Context, props *notificationActionProps, actionFn func(...*notificationActionProps) *notificationAction, err error) error {
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
		a.Description = props.Format(m.AsString(notificationLogMetaKey{}), err)

		if p, has := m[notificationPropsMetaKey{}]; has {
			a.Meta = p.(*notificationActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
