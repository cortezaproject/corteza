package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// federation/service/node_actions.yaml

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
)

type (
	nodeActionProps struct {
		node       *types.Node
		pairingURI string
		filter     *types.NodeFilter
	}

	nodeAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *nodeActionProps
	}

	nodeError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *nodeActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setNode updates nodeActionProps's node
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *nodeActionProps) setNode(node *types.Node) *nodeActionProps {
	p.node = node
	return p
}

// setPairingURI updates nodeActionProps's pairingURI
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *nodeActionProps) setPairingURI(pairingURI string) *nodeActionProps {
	p.pairingURI = pairingURI
	return p
}

// setFilter updates nodeActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *nodeActionProps) setFilter(filter *types.NodeFilter) *nodeActionProps {
	p.filter = filter
	return p
}

// serialize converts nodeActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p nodeActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.node != nil {
		m.Set("node.Name", p.node.Name, true)
		m.Set("node.BaseURL", p.node.BaseURL, true)
		m.Set("node.ID", p.node.ID, true)
		m.Set("node.Status", p.node.Status, true)
	}
	m.Set("pairingURI", p.pairingURI, true)
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.status", p.filter.Status, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p nodeActionProps) tr(in string, err error) string {
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

	if p.node != nil {
		// replacement for "{node}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{node}",
			fns(
				p.node.Name,
				p.node.BaseURL,
				p.node.ID,
				p.node.Status,
			),
		)
		pairs = append(pairs, "{node.Name}", fns(p.node.Name))
		pairs = append(pairs, "{node.BaseURL}", fns(p.node.BaseURL))
		pairs = append(pairs, "{node.ID}", fns(p.node.ID))
		pairs = append(pairs, "{node.Status}", fns(p.node.Status))
	}
	pairs = append(pairs, "{pairingURI}", fns(p.pairingURI))

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.Status,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.status}", fns(p.filter.Status))
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
func (a *nodeAction) String() string {
	var props = &nodeActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *nodeAction) LoggableAction() *actionlog.Action {
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
func (e *nodeError) String() string {
	var props = &nodeActionProps{}

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
func (e *nodeError) Error() string {
	var props = &nodeActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *nodeError) Is(err error) bool {
	t, ok := err.(*nodeError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *nodeError) IsGeneric() bool {
	return e.error == "generic"
}

// Wrap wraps nodeError around another error
//
// This function is auto-generated.
//
func (e *nodeError) Wrap(err error) *nodeError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *nodeError) Unwrap() error {
	return e.wrap
}

func (e *nodeError) LoggableAction() *actionlog.Action {
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

// NodeActionSearch returns "federation:node.search" error
//
// This function is auto-generated.
//
func NodeActionSearch(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "search",
		log:       "searched for nodes",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionLookup returns "federation:node.lookup" error
//
// This function is auto-generated.
//
func NodeActionLookup(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "lookup",
		log:       "looked-up for a {node}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionCreate returns "federation:node.create" error
//
// This function is auto-generated.
//
func NodeActionCreate(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "create",
		log:       "created {node}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionCreateFromPairingURI returns "federation:node.createFromPairingURI" error
//
// This function is auto-generated.
//
func NodeActionCreateFromPairingURI(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "createFromPairingURI",
		log:       "created {node} from pairing URI",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionRecreateFromPairingURI returns "federation:node.recreateFromPairingURI" error
//
// This function is auto-generated.
//
func NodeActionRecreateFromPairingURI(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "recreateFromPairingURI",
		log:       "recreate {node} from pairing URI",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionUpdate returns "federation:node.update" error
//
// This function is auto-generated.
//
func NodeActionUpdate(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "update",
		log:       "updated {node}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionDelete returns "federation:node.delete" error
//
// This function is auto-generated.
//
func NodeActionDelete(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "delete",
		log:       "deleted {node}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionUndelete returns "federation:node.undelete" error
//
// This function is auto-generated.
//
func NodeActionUndelete(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "undelete",
		log:       "undeleted {node}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionOttRegenerated returns "federation:node.ottRegenerated" error
//
// This function is auto-generated.
//
func NodeActionOttRegenerated(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "ottRegenerated",
		log:       "regenerated one-time-token for {node}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionPair returns "federation:node.pair" error
//
// This function is auto-generated.
//
func NodeActionPair(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "pair",
		log:       "{node} pairing started",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionHandshakeInit returns "federation:node.handshakeInit" error
//
// This function is auto-generated.
//
func NodeActionHandshakeInit(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "handshakeInit",
		log:       "{node} handshake initialized",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionHandshakeConfirm returns "federation:node.handshakeConfirm" error
//
// This function is auto-generated.
//
func NodeActionHandshakeConfirm(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "handshakeConfirm",
		log:       "{node} handshake confirmed",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionHandshakeComplete returns "federation:node.handshakeComplete" error
//
// This function is auto-generated.
//
func NodeActionHandshakeComplete(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "handshakeComplete",
		log:       "{node} handshake completed",
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

// NodeErrGeneric returns "federation:node.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NodeErrGeneric(props ...*nodeActionProps) *nodeError {
	var e = &nodeError{
		timestamp: time.Now(),
		resource:  "federation:node",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *nodeActionProps {
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

// NodeErrNotFound returns "federation:node.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func NodeErrNotFound(props ...*nodeActionProps) *nodeError {
	var e = &nodeError{
		timestamp: time.Now(),
		resource:  "federation:node",
		error:     "notFound",
		action:    "error",
		message:   "node does not exist",
		log:       "node does not exist",
		severity:  actionlog.Warning,
		props: func() *nodeActionProps {
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

// NodeErrPairingURIInvalid returns "federation:node.pairingURIInvalid" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NodeErrPairingURIInvalid(props ...*nodeActionProps) *nodeError {
	var e = &nodeError{
		timestamp: time.Now(),
		resource:  "federation:node",
		error:     "pairingURIInvalid",
		action:    "error",
		message:   "pairing URI invalid: {err}",
		log:       "pairing URI invalid: {err}",
		severity:  actionlog.Error,
		props: func() *nodeActionProps {
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

// NodeErrPairingURITokenInvalid returns "federation:node.pairingURITokenInvalid" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NodeErrPairingURITokenInvalid(props ...*nodeActionProps) *nodeError {
	var e = &nodeError{
		timestamp: time.Now(),
		resource:  "federation:node",
		error:     "pairingURITokenInvalid",
		action:    "error",
		message:   "pairing URI with invalid pairing token",
		log:       "pairing URI with invalid pairing token",
		severity:  actionlog.Error,
		props: func() *nodeActionProps {
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

// NodeErrPairingURISourceIDInvalid returns "federation:node.pairingURISourceIDInvalid" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NodeErrPairingURISourceIDInvalid(props ...*nodeActionProps) *nodeError {
	var e = &nodeError{
		timestamp: time.Now(),
		resource:  "federation:node",
		error:     "pairingURISourceIDInvalid",
		action:    "error",
		message:   "pairing URI without source node ID",
		log:       "pairing URI without source node ID",
		severity:  actionlog.Error,
		props: func() *nodeActionProps {
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

// NodeErrPairingTokenInvalid returns "federation:node.pairingTokenInvalid" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NodeErrPairingTokenInvalid(props ...*nodeActionProps) *nodeError {
	var e = &nodeError{
		timestamp: time.Now(),
		resource:  "federation:node",
		error:     "pairingTokenInvalid",
		action:    "error",
		message:   "pairing token invalid",
		log:       "pairing token invalid",
		severity:  actionlog.Error,
		props: func() *nodeActionProps {
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
// action (optional) fn will be used to construct nodeAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc node) recordAction(ctx context.Context, props *nodeActionProps, action func(...*nodeActionProps) *nodeAction, err error) error {
	var (
		ok bool

		// Return error
		retError *nodeError

		// Recorder error
		recError *nodeError
	)

	if err != nil {
		if retError, ok = err.(*nodeError); !ok {
			// got non-node error, wrap it with NodeErrGeneric
			retError = NodeErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use NodeErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type nodeError
				if unwrappedSinkError, ok := unwrappedError.(*nodeError); ok {
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
