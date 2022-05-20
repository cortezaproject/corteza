package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/connection_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
	"time"
)

type (
	connectionActionProps struct {
		connection *types.Connection
		new        *types.Connection
		update     *types.Connection
		search     *types.ConnectionFilter
	}

	connectionAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *connectionActionProps
	}

	connectionLogMetaKey   struct{}
	connectionPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setConnection updates connectionActionProps's connection
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *connectionActionProps) setConnection(connection *types.Connection) *connectionActionProps {
	p.connection = connection
	return p
}

// setNew updates connectionActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *connectionActionProps) setNew(new *types.Connection) *connectionActionProps {
	p.new = new
	return p
}

// setUpdate updates connectionActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *connectionActionProps) setUpdate(update *types.Connection) *connectionActionProps {
	p.update = update
	return p
}

// setSearch updates connectionActionProps's search
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *connectionActionProps) setSearch(search *types.ConnectionFilter) *connectionActionProps {
	p.search = search
	return p
}

// Serialize converts connectionActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p connectionActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.connection != nil {
		m.Set("connection.DSN", p.connection.DSN, true)
		m.Set("connection.ID", p.connection.ID, true)
	}
	if p.new != nil {
		m.Set("new.DSN", p.new.DSN, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.DSN", p.update.DSN, true)
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
func (p connectionActionProps) Format(in string, err error) string {
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

	if p.connection != nil {
		// replacement for "{{connection}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{connection}}",
			fns(
				p.connection.DSN,
				p.connection.ID,
			),
		)
		pairs = append(pairs, "{{connection.DSN}}", fns(p.connection.DSN))
		pairs = append(pairs, "{{connection.ID}}", fns(p.connection.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.DSN,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.DSN}}", fns(p.new.DSN))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.DSN,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.DSN}}", fns(p.update.DSN))
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
func (a *connectionAction) String() string {
	var props = &connectionActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *connectionAction) ToAction() *actionlog.Action {
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

// ConnectionActionSearch returns "system:connection.search" action
//
// This function is auto-generated.
//
func ConnectionActionSearch(props ...*connectionActionProps) *connectionAction {
	a := &connectionAction{
		timestamp: time.Now(),
		resource:  "system:connection",
		action:    "search",
		log:       "searched for connection",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ConnectionActionLookup returns "system:connection.lookup" action
//
// This function is auto-generated.
//
func ConnectionActionLookup(props ...*connectionActionProps) *connectionAction {
	a := &connectionAction{
		timestamp: time.Now(),
		resource:  "system:connection",
		action:    "lookup",
		log:       "looked-up for a {{connection}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ConnectionActionCreate returns "system:connection.create" action
//
// This function is auto-generated.
//
func ConnectionActionCreate(props ...*connectionActionProps) *connectionAction {
	a := &connectionAction{
		timestamp: time.Now(),
		resource:  "system:connection",
		action:    "create",
		log:       "created {{connection}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ConnectionActionUpdate returns "system:connection.update" action
//
// This function is auto-generated.
//
func ConnectionActionUpdate(props ...*connectionActionProps) *connectionAction {
	a := &connectionAction{
		timestamp: time.Now(),
		resource:  "system:connection",
		action:    "update",
		log:       "updated {{connection}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ConnectionActionDelete returns "system:connection.delete" action
//
// This function is auto-generated.
//
func ConnectionActionDelete(props ...*connectionActionProps) *connectionAction {
	a := &connectionAction{
		timestamp: time.Now(),
		resource:  "system:connection",
		action:    "delete",
		log:       "deleted {{connection}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ConnectionActionUndelete returns "system:connection.undelete" action
//
// This function is auto-generated.
//
func ConnectionActionUndelete(props ...*connectionActionProps) *connectionAction {
	a := &connectionAction{
		timestamp: time.Now(),
		resource:  "system:connection",
		action:    "undelete",
		log:       "undeleted {{connection}}",
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

// ConnectionErrGeneric returns "system:connection.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrGeneric(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(connectionLogMetaKey{}, "{err}"),
		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrNotFound returns "system:connection.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrNotFound(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("connection not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:connection"),

		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrInvalidID returns "system:connection.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrInvalidID(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:connection"),

		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrInvalidEndpoint returns "system:connection.invalidEndpoint" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrInvalidEndpoint(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid DSN", nil),

		errors.Meta("type", "invalidEndpoint"),
		errors.Meta("resource", "system:connection"),

		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.invalidEndpoint"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrExistsEndpoint returns "system:connection.existsEndpoint" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrExistsEndpoint(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("connection with this DSN already exists", nil),

		errors.Meta("type", "existsEndpoint"),
		errors.Meta("resource", "system:connection"),

		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.existsEndpoint"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrAlreadyExists returns "system:connection.alreadyExists" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrAlreadyExists(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("connection by that DSN already exists", nil),

		errors.Meta("type", "alreadyExists"),
		errors.Meta("resource", "system:connection"),

		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.alreadyExists"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrNotAllowedToCreate returns "system:connection.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrNotAllowedToCreate(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a connection", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(connectionLogMetaKey{}, "failed to create a connection; insufficient permissions"),
		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrNotAllowedToRead returns "system:connection.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrNotAllowedToRead(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this connection", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(connectionLogMetaKey{}, "failed to read {{connection.DSN}}; insufficient permissions"),
		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrNotAllowedToSearch returns "system:connection.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrNotAllowedToSearch(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list or search connections", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(connectionLogMetaKey{}, "failed to search for connections; insufficient permissions"),
		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrNotAllowedToUpdate returns "system:connection.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrNotAllowedToUpdate(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this connection", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(connectionLogMetaKey{}, "failed to update {{connection.DSN}}; insufficient permissions"),
		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrNotAllowedToDelete returns "system:connection.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrNotAllowedToDelete(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this connection", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(connectionLogMetaKey{}, "failed to delete {{connection.DSN}}; insufficient permissions"),
		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrNotAllowedToUndelete returns "system:connection.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrNotAllowedToUndelete(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this connection", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(connectionLogMetaKey{}, "failed to undelete {{connection.DSN}}; insufficient permissions"),
		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConnectionErrNotAllowedToExec returns "system:connection.notAllowedToExec" as *errors.Error
//
//
// This function is auto-generated.
//
func ConnectionErrNotAllowedToExec(mm ...*connectionActionProps) *errors.Error {
	var p = &connectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to execute this connection", nil),

		errors.Meta("type", "notAllowedToExec"),
		errors.Meta("resource", "system:connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(connectionLogMetaKey{}, "failed to exec {{connection.DSN}}; insufficient permissions"),
		errors.Meta(connectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "connection.errors.notAllowedToExec"),

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
func (svc connection) recordAction(ctx context.Context, props *connectionActionProps, actionFn func(...*connectionActionProps) *connectionAction, err error) error {
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
		a.Description = props.Format(m.AsString(connectionLogMetaKey{}), err)

		if p, has := m[connectionPropsMetaKey{}]; has {
			a.Meta = p.(*connectionActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
