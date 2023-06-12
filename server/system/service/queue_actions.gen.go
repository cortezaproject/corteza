package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/queue_actions.yaml

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
	queueActionProps struct {
		queue  *types.Queue
		new    *types.Queue
		update *types.Queue
		search *types.QueueFilter
	}

	queueAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *queueActionProps
	}

	queueLogMetaKey   struct{}
	queuePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setQueue updates queueActionProps's queue
//
// This function is auto-generated.
//
func (p *queueActionProps) setQueue(queue *types.Queue) *queueActionProps {
	p.queue = queue
	return p
}

// setNew updates queueActionProps's new
//
// This function is auto-generated.
//
func (p *queueActionProps) setNew(new *types.Queue) *queueActionProps {
	p.new = new
	return p
}

// setUpdate updates queueActionProps's update
//
// This function is auto-generated.
//
func (p *queueActionProps) setUpdate(update *types.Queue) *queueActionProps {
	p.update = update
	return p
}

// setSearch updates queueActionProps's search
//
// This function is auto-generated.
//
func (p *queueActionProps) setSearch(search *types.QueueFilter) *queueActionProps {
	p.search = search
	return p
}

// Serialize converts queueActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p queueActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.queue != nil {
		m.Set("queue.queue", p.queue.Queue, true)
		m.Set("queue.ID", p.queue.ID, true)
	}
	if p.new != nil {
		m.Set("new.queue", p.new.Queue, true)
		m.Set("new.consumer", p.new.Consumer, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.queue", p.update.Queue, true)
		m.Set("update.consumer", p.update.Consumer, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.search != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p queueActionProps) Format(in string, err error) string {
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

	if p.queue != nil {
		// replacement for "{{queue}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{queue}}",
			fns(
				p.queue.Queue,
				p.queue.ID,
			),
		)
		pairs = append(pairs, "{{queue.queue}}", fns(p.queue.Queue))
		pairs = append(pairs, "{{queue.ID}}", fns(p.queue.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Queue,
				p.new.Consumer,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.queue}}", fns(p.new.Queue))
		pairs = append(pairs, "{{new.consumer}}", fns(p.new.Consumer))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Queue,
				p.update.Consumer,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.queue}}", fns(p.update.Queue))
		pairs = append(pairs, "{{update.consumer}}", fns(p.update.Consumer))
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.search != nil {
		// replacement for "{{search}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{search}}",
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
func (a *queueAction) String() string {
	var props = &queueActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *queueAction) ToAction() *actionlog.Action {
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

// QueueActionSearch returns "system:queue.search" action
//
// This function is auto-generated.
//
func QueueActionSearch(props ...*queueActionProps) *queueAction {
	a := &queueAction{
		timestamp: time.Now(),
		resource:  "system:queue",
		action:    "search",
		log:       "searched for queues",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// QueueActionLookup returns "system:queue.lookup" action
//
// This function is auto-generated.
//
func QueueActionLookup(props ...*queueActionProps) *queueAction {
	a := &queueAction{
		timestamp: time.Now(),
		resource:  "system:queue",
		action:    "lookup",
		log:       "looked-up for a {{queue}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// QueueActionCreate returns "system:queue.create" action
//
// This function is auto-generated.
//
func QueueActionCreate(props ...*queueActionProps) *queueAction {
	a := &queueAction{
		timestamp: time.Now(),
		resource:  "system:queue",
		action:    "create",
		log:       "created {{queue}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// QueueActionUpdate returns "system:queue.update" action
//
// This function is auto-generated.
//
func QueueActionUpdate(props ...*queueActionProps) *queueAction {
	a := &queueAction{
		timestamp: time.Now(),
		resource:  "system:queue",
		action:    "update",
		log:       "updated {{queue}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// QueueActionDelete returns "system:queue.delete" action
//
// This function is auto-generated.
//
func QueueActionDelete(props ...*queueActionProps) *queueAction {
	a := &queueAction{
		timestamp: time.Now(),
		resource:  "system:queue",
		action:    "delete",
		log:       "deleted {{queue}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// QueueActionUndelete returns "system:queue.undelete" action
//
// This function is auto-generated.
//
func QueueActionUndelete(props ...*queueActionProps) *queueAction {
	a := &queueAction{
		timestamp: time.Now(),
		resource:  "system:queue",
		action:    "undelete",
		log:       "undeleted {{queue}}",
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

// QueueErrGeneric returns "system:queue.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrGeneric(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:queue"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(queueLogMetaKey{}, "{err}"),
		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrNotFound returns "system:queue.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrNotFound(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("queue not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:queue"),

		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrInvalidID returns "system:queue.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrInvalidID(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:queue"),

		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrInvalidConsumer returns "system:queue.invalidConsumer" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrInvalidConsumer(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid consumer", nil),

		errors.Meta("type", "invalidConsumer"),
		errors.Meta("resource", "system:queue"),

		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.invalidConsumer"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrStaleData returns "system:queue.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrStaleData(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:queue"),

		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrAlreadyExists returns "system:queue.alreadyExists" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrAlreadyExists(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("queue by that name already exists", nil),

		errors.Meta("type", "alreadyExists"),
		errors.Meta("resource", "system:queue"),

		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.alreadyExists"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrNotAllowedToCreate returns "system:queue.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrNotAllowedToCreate(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a queue", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:queue"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(queueLogMetaKey{}, "failed to create a queue; insufficient permissions"),
		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrNotAllowedToRead returns "system:queue.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrNotAllowedToRead(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this queue", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:queue"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(queueLogMetaKey{}, "failed to read {{queue.queue}}; insufficient permissions"),
		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrNotAllowedToSearch returns "system:queue.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrNotAllowedToSearch(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list queues", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:queue"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(queueLogMetaKey{}, "failed to search or list; insufficient permissions"),
		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrNotAllowedToUpdate returns "system:queue.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrNotAllowedToUpdate(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this queue", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:queue"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(queueLogMetaKey{}, "failed to update {{queue.queue}}; insufficient permissions"),
		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrNotAllowedToDelete returns "system:queue.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrNotAllowedToDelete(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this queue", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:queue"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(queueLogMetaKey{}, "failed to delete {{queue.queue}}; insufficient permissions"),
		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrNotAllowedToUndelete returns "system:queue.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrNotAllowedToUndelete(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this queue", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:queue"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(queueLogMetaKey{}, "failed to undelete {{queue.queue}}; insufficient permissions"),
		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrNotAllowedToWriteTo returns "system:queue.notAllowedToWriteTo" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrNotAllowedToWriteTo(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to add messages to this queue", nil),

		errors.Meta("type", "notAllowedToWriteTo"),
		errors.Meta("resource", "system:queue"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(queueLogMetaKey{}, "failed to add message to {{queue.queue}}; insufficient permissions"),
		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.notAllowedToWriteTo"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// QueueErrNotAllowedToReadFrom returns "system:queue.notAllowedToReadFrom" as *errors.Error
//
//
// This function is auto-generated.
//
func QueueErrNotAllowedToReadFrom(mm ...*queueActionProps) *errors.Error {
	var p = &queueActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read messages from this queue", nil),

		errors.Meta("type", "notAllowedToReadFrom"),
		errors.Meta("resource", "system:queue"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(queueLogMetaKey{}, "failed to read message from {{queue.queue}}; insufficient permissions"),
		errors.Meta(queuePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "queue.errors.notAllowedToReadFrom"),

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
func (svc queue) recordAction(ctx context.Context, props *queueActionProps, actionFn func(...*queueActionProps) *queueAction, err error) error {
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
		a.Description = props.Format(m.AsString(queueLogMetaKey{}), err)

		if p, has := m[queuePropsMetaKey{}]; has {
			a.Meta = p.(*queueActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
