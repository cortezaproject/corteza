package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/configured_connection_actions.yaml

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
	configuredConnectionActionProps struct {
		connection *types.ConfiguredConnection
		filter     *types.ConfiguredConnectionFilter
	}

	configuredConnectionAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *configuredConnectionActionProps
	}

	configuredConnectionLogMetaKey   struct{}
	configuredConnectionPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setConnection updates configuredConnectionActionProps's connection
//
// This function is auto-generated.
func (p *configuredConnectionActionProps) setConnection(connection *types.ConfiguredConnection) *configuredConnectionActionProps {
	p.connection = connection
	return p
}

// setFilter updates configuredConnectionActionProps's filter
//
// This function is auto-generated.
func (p *configuredConnectionActionProps) setFilter(filter *types.ConfiguredConnectionFilter) *configuredConnectionActionProps {
	p.filter = filter
	return p
}

// Serialize converts configuredConnectionActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p configuredConnectionActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.connection != nil {
		m.Set("connection.name", p.connection.Name, true)
		m.Set("connection.ID", p.connection.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.connectionID", p.filter.ConnectionID, true)
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.status", p.filter.Status, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p configuredConnectionActionProps) Format(in string, err error) string {
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
				p.connection.Name,
				p.connection.ID,
			),
		)
		pairs = append(pairs, "{{connection.name}}", fns(p.connection.Name))
		pairs = append(pairs, "{{connection.ID}}", fns(p.connection.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.ConnectionID,
				p.filter.Query,
				p.filter.Status,
			),
		)
		pairs = append(pairs, "{{filter.connectionID}}", fns(p.filter.ConnectionID))
		pairs = append(pairs, "{{filter.query}}", fns(p.filter.Query))
		pairs = append(pairs, "{{filter.status}}", fns(p.filter.Status))
	}
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
func (a *configuredConnectionAction) String() string {
	var props = &configuredConnectionActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *configuredConnectionAction) ToAction() *actionlog.Action {
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

// ConfiguredConnectionActionSearch returns "system:configured-connection.search" action
//
// This function is auto-generated.
func ConfiguredConnectionActionSearch(props ...*configuredConnectionActionProps) *configuredConnectionAction {
	a := &configuredConnectionAction{
		timestamp: time.Now(),
		resource:  "system:configured-connection",
		action:    "search",
		log:       "searched for matching connections",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ConfiguredConnectionActionLookup returns "system:configured-connection.lookup" action
//
// This function is auto-generated.
func ConfiguredConnectionActionLookup(props ...*configuredConnectionActionProps) *configuredConnectionAction {
	a := &configuredConnectionAction{
		timestamp: time.Now(),
		resource:  "system:configured-connection",
		action:    "lookup",
		log:       "looked-up for a {{connection}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ConfiguredConnectionActionInstall returns "system:configured-connection.install" action
//
// This function is auto-generated.
func ConfiguredConnectionActionInstall(props ...*configuredConnectionActionProps) *configuredConnectionAction {
	a := &configuredConnectionAction{
		timestamp: time.Now(),
		resource:  "system:configured-connection",
		action:    "install",
		log:       "installed {{connection}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ConfiguredConnectionActionDelete returns "system:configured-connection.delete" action
//
// This function is auto-generated.
func ConfiguredConnectionActionDelete(props ...*configuredConnectionActionProps) *configuredConnectionAction {
	a := &configuredConnectionAction{
		timestamp: time.Now(),
		resource:  "system:configured-connection",
		action:    "delete",
		log:       "deleted {{connection}}",
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

// ConfiguredConnectionErrGeneric returns "system:configured-connection.generic" as *errors.Error
//
// This function is auto-generated.
func ConfiguredConnectionErrGeneric(mm ...*configuredConnectionActionProps) *errors.Error {
	var p = &configuredConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:configured-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(configuredConnectionLogMetaKey{}, "{err}"),
		errors.Meta(configuredConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "configured-connection.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConfiguredConnectionErrNotFound returns "system:configured-connection.notFound" as *errors.Error
//
// This function is auto-generated.
func ConfiguredConnectionErrNotFound(mm ...*configuredConnectionActionProps) *errors.Error {
	var p = &configuredConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("connection not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:configured-connection"),

		errors.Meta(configuredConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "configured-connection.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConfiguredConnectionErrInvalidID returns "system:configured-connection.invalidID" as *errors.Error
//
// This function is auto-generated.
func ConfiguredConnectionErrInvalidID(mm ...*configuredConnectionActionProps) *errors.Error {
	var p = &configuredConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:configured-connection"),

		errors.Meta(configuredConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "configured-connection.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConfiguredConnectionErrDeletionNotSupported returns "system:configured-connection.deletionNotSupported" as *errors.Error
//
// This function is auto-generated.
func ConfiguredConnectionErrDeletionNotSupported(mm ...*configuredConnectionActionProps) *errors.Error {
	var p = &configuredConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("connection deletion is not yet supported", nil),

		errors.Meta("type", "deletionNotSupported"),
		errors.Meta("resource", "system:configured-connection"),

		errors.Meta(configuredConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "configured-connection.errors.deletionNotSupported"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConfiguredConnectionErrNotAllowedToRead returns "system:configured-connection.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ConfiguredConnectionErrNotAllowedToRead(mm ...*configuredConnectionActionProps) *errors.Error {
	var p = &configuredConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this connection", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:configured-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(configuredConnectionLogMetaKey{}, "failed to read {{connection.name}}; insufficient permissions"),
		errors.Meta(configuredConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "configured-connection.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConfiguredConnectionErrNotAllowedToSearch returns "system:configured-connection.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ConfiguredConnectionErrNotAllowedToSearch(mm ...*configuredConnectionActionProps) *errors.Error {
	var p = &configuredConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search connections", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:configured-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(configuredConnectionLogMetaKey{}, "failed to search for connections; insufficient permissions"),
		errors.Meta(configuredConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "configured-connection.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConfiguredConnectionErrNotAllowedToInstall returns "system:configured-connection.notAllowedToInstall" as *errors.Error
//
// This function is auto-generated.
func ConfiguredConnectionErrNotAllowedToInstall(mm ...*configuredConnectionActionProps) *errors.Error {
	var p = &configuredConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to install connections", nil),

		errors.Meta("type", "notAllowedToInstall"),
		errors.Meta("resource", "system:configured-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(configuredConnectionLogMetaKey{}, "failed to install connection; insufficient permissions"),
		errors.Meta(configuredConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "configured-connection.errors.notAllowedToInstall"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ConfiguredConnectionErrNotAllowedToDelete returns "system:configured-connection.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ConfiguredConnectionErrNotAllowedToDelete(mm ...*configuredConnectionActionProps) *errors.Error {
	var p = &configuredConnectionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this connection", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:configured-connection"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(configuredConnectionLogMetaKey{}, "failed to delete {{connection.name}}; insufficient permissions"),
		errors.Meta(configuredConnectionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "configured-connection.errors.notAllowedToDelete"),

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
func (svc configuredConnection) recordAction(ctx context.Context, props *configuredConnectionActionProps, actionFn func(...*configuredConnectionActionProps) *configuredConnectionAction, err error) error {
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
		a.Description = props.Format(m.AsString(configuredConnectionLogMetaKey{}), err)

		if p, has := m[configuredConnectionPropsMetaKey{}]; has {
			a.Meta = p.(*configuredConnectionActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
