package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/dal_connection_actions.yaml

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
	dalConnectionActionProps struct {
		connection *types.DalConnection
		new        *types.DalConnection
		update     *types.DalConnection
		search     *types.DalConnectionFilter
	}

	dalConnectionAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *dalConnectionActionProps
	}

	dalConnectionLogMetaKey   struct{}
	dalConnectionPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setConnection updates dalConnectionActionProps's connection
//
// This function is auto-generated.
//
func (p *dalConnectionActionProps) setConnection(connection *types.DalConnection) *dalConnectionActionProps {
	p.connection = connection
	return p
}

// setNew updates dalConnectionActionProps's new
//
// This function is auto-generated.
//
func (p *dalConnectionActionProps) setNew(new *types.DalConnection) *dalConnectionActionProps {
	p.new = new
	return p
}

// setUpdate updates dalConnectionActionProps's update
//
// This function is auto-generated.
//
func (p *dalConnectionActionProps) setUpdate(update *types.DalConnection) *dalConnectionActionProps {
	p.update = update
	return p
}

// setSearch updates dalConnectionActionProps's search
//
// This function is auto-generated.
//
func (p *dalConnectionActionProps) setSearch(search *types.DalConnectionFilter) *dalConnectionActionProps {
	p.search = search
	return p
}

// Serialize converts dalConnectionActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p dalConnectionActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.connection != nil {
		m.Set("connection.handle", p.connection.Handle, true)
		m.Set("connection.ID", p.connection.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
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
func (p dalConnectionActionProps) Format(in string, err error) string {
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
				p.connection.Handle,
				p.connection.ID,
			),
		)
		pairs = append(pairs, "{{connection.handle}}", fns(p.connection.Handle))
		pairs = append(pairs, "{{connection.ID}}", fns(p.connection.ID))
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
func (a *dalConnectionAction) String() string {
	var props = &dalConnectionActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *dalConnectionAction) ToAction() *actionlog.Action {
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

// DalConnectionActionSearch returns "system:dal-connection.search" action
//
// This function is auto-generated.
//
func DalConnectionActionSearch(props ...*dalConnectionActionProps) *dalConnectionAction {
	a := &dalConnectionAction{
		timestamp: time.Now(),
		resource:  "system:dal-connection",
		action:    "search",
		log:       "searched for connection",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalConnectionActionLookup returns "system:dal-connection.lookup" action
//
// This function is auto-generated.
//
func DalConnectionActionLookup(props ...*dalConnectionActionProps) *dalConnectionAction {
	a := &dalConnectionAction{
		timestamp: time.Now(),
		resource:  "system:dal-connection",
		action:    "lookup",
		log:       "looked-up for a {{connection}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalConnectionActionCreate returns "system:dal-connection.create" action
//
// This function is auto-generated.
//
func DalConnectionActionCreate(props ...*dalConnectionActionProps) *dalConnectionAction {
	a := &dalConnectionAction{
		timestamp: time.Now(),
		resource:  "system:dal-connection",
		action:    "create",
		log:       "created {{connection}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalConnectionActionUpdate returns "system:dal-connection.update" action
//
// This function is auto-generated.
//
func DalConnectionActionUpdate(props ...*dalConnectionActionProps) *dalConnectionAction {
	a := &dalConnectionAction{
		timestamp: time.Now(),
		resource:  "system:dal-connection",
		action:    "update",
		log:       "updated {{connection}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalConnectionActionDelete returns "system:dal-connection.delete" action
//
// This function is auto-generated.
//
func DalConnectionActionDelete(props ...*dalConnectionActionProps) *dalConnectionAction {
	a := &dalConnectionAction{
		timestamp: time.Now(),
		resource:  "system:dal-connection",
		action:    "delete",
		log:       "deleted {{connection}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalConnectionActionUndelete returns "system:dal-connection.undelete" action
//
// This function is auto-generated.
//
func DalConnectionActionUndelete(props ...*dalConnectionActionProps) *dalConnectionAction {
	a := &dalConnectionAction{
		timestamp: time.Now(),
		resource:  "system:dal-connection",
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

// DalConnectionErrGeneric returns "system:dal-connection.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrGeneric(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:dal-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalConnectionLogMetaKey{}, "{err}"),
		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrNotFound returns "system:dal-connection.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrNotFound(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("connection not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:dal-connection"),

		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrInvalidID returns "system:dal-connection.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrInvalidID(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:dal-connection"),

		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrMissingName returns "system:dal-connection.missingName" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrMissingName(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("missing name", nil),

		errors.Meta("type", "missingName"),
		errors.Meta("resource", "system:dal-connection"),

		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.missingName"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrInvalidEndpoint returns "system:dal-connection.invalidEndpoint" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrInvalidEndpoint(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid DSN", nil),

		errors.Meta("type", "invalidEndpoint"),
		errors.Meta("resource", "system:dal-connection"),

		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.invalidEndpoint"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrStaleData returns "system:dal-connection.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrStaleData(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:dal-connection"),

		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrExistsEndpoint returns "system:dal-connection.existsEndpoint" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrExistsEndpoint(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("connection with this DSN already exists", nil),

		errors.Meta("type", "existsEndpoint"),
		errors.Meta("resource", "system:dal-connection"),

		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.existsEndpoint"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrAlreadyExists returns "system:dal-connection.alreadyExists" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrAlreadyExists(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("connection by that DSN already exists", nil),

		errors.Meta("type", "alreadyExists"),
		errors.Meta("resource", "system:dal-connection"),

		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.alreadyExists"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrNotAllowedToCreate returns "system:dal-connection.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrNotAllowedToCreate(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a connection", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:dal-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalConnectionLogMetaKey{}, "failed to create a connection; insufficient permissions"),
		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrNotAllowedToRead returns "system:dal-connection.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrNotAllowedToRead(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this connection", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:dal-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalConnectionLogMetaKey{}, "failed to read {{connection.handle}}; insufficient permissions"),
		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrNotAllowedToSearch returns "system:dal-connection.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrNotAllowedToSearch(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list or search connections", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:dal-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalConnectionLogMetaKey{}, "failed to search for connections; insufficient permissions"),
		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrNotAllowedToUpdate returns "system:dal-connection.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrNotAllowedToUpdate(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this connection", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:dal-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalConnectionLogMetaKey{}, "failed to update {{connection.handle}}; insufficient permissions"),
		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrNotAllowedToDelete returns "system:dal-connection.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrNotAllowedToDelete(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this connection", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:dal-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalConnectionLogMetaKey{}, "failed to delete {{connection.handle}}; insufficient permissions"),
		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrNotAllowedToUndelete returns "system:dal-connection.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrNotAllowedToUndelete(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this connection", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:dal-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalConnectionLogMetaKey{}, "failed to undelete {{connection.handle}}; insufficient permissions"),
		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalConnectionErrNotAllowedToExec returns "system:dal-connection.notAllowedToExec" as *errors.Error
//
//
// This function is auto-generated.
//
func DalConnectionErrNotAllowedToExec(mm ...*dalConnectionActionProps) *errors.Error {
	var p = &dalConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to execute this connection", nil),

		errors.Meta("type", "notAllowedToExec"),
		errors.Meta("resource", "system:dal-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalConnectionLogMetaKey{}, "failed to exec {{connection.handle}}; insufficient permissions"),
		errors.Meta(dalConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-connection.errors.notAllowedToExec"),

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
func (svc dalConnection) recordAction(ctx context.Context, props *dalConnectionActionProps, actionFn func(...*dalConnectionActionProps) *dalConnectionAction, err error) error {
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
		a.Description = props.Format(m.AsString(dalConnectionLogMetaKey{}), err)

		if p, has := m[dalConnectionPropsMetaKey{}]; has {
			a.Meta = p.(*dalConnectionActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
