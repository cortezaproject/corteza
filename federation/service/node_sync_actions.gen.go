package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// federation/service/node_sync_actions.yaml

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
	nodeSyncActionProps struct {
		nodeSync       *types.NodeSync
		nodeSyncFilter *types.NodeSyncFilter
	}

	nodeSyncAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *nodeSyncActionProps
	}

	nodeSyncError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *nodeSyncActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setNodeSync updates nodeSyncActionProps's nodeSync
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *nodeSyncActionProps) setNodeSync(nodeSync *types.NodeSync) *nodeSyncActionProps {
	p.nodeSync = nodeSync
	return p
}

// setNodeSyncFilter updates nodeSyncActionProps's nodeSyncFilter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *nodeSyncActionProps) setNodeSyncFilter(nodeSyncFilter *types.NodeSyncFilter) *nodeSyncActionProps {
	p.nodeSyncFilter = nodeSyncFilter
	return p
}

// serialize converts nodeSyncActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p nodeSyncActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.nodeSync != nil {
		m.Set("nodeSync.NodeID", p.nodeSync.NodeID, true)
		m.Set("nodeSync.SyncStatus", p.nodeSync.SyncStatus, true)
		m.Set("nodeSync.SyncType", p.nodeSync.SyncType, true)
		m.Set("nodeSync.TimeOfAction", p.nodeSync.TimeOfAction, true)
	}
	if p.nodeSyncFilter != nil {
		m.Set("nodeSyncFilter.Query", p.nodeSyncFilter.Query, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p nodeSyncActionProps) tr(in string, err error) string {
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

	if p.nodeSync != nil {
		// replacement for "{nodeSync}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{nodeSync}",
			fns(
				p.nodeSync.NodeID,
				p.nodeSync.SyncStatus,
				p.nodeSync.SyncType,
				p.nodeSync.TimeOfAction,
			),
		)
		pairs = append(pairs, "{nodeSync.NodeID}", fns(p.nodeSync.NodeID))
		pairs = append(pairs, "{nodeSync.SyncStatus}", fns(p.nodeSync.SyncStatus))
		pairs = append(pairs, "{nodeSync.SyncType}", fns(p.nodeSync.SyncType))
		pairs = append(pairs, "{nodeSync.TimeOfAction}", fns(p.nodeSync.TimeOfAction))
	}

	if p.nodeSyncFilter != nil {
		// replacement for "{nodeSyncFilter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{nodeSyncFilter}",
			fns(
				p.nodeSyncFilter.Query,
			),
		)
		pairs = append(pairs, "{nodeSyncFilter.Query}", fns(p.nodeSyncFilter.Query))
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
func (a *nodeSyncAction) String() string {
	var props = &nodeSyncActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *nodeSyncAction) LoggableAction() *actionlog.Action {
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
func (e *nodeSyncError) String() string {
	var props = &nodeSyncActionProps{}

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
func (e *nodeSyncError) Error() string {
	var props = &nodeSyncActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *nodeSyncError) Is(err error) bool {
	t, ok := err.(*nodeSyncError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *nodeSyncError) IsGeneric() bool {
	return e.error == "generic"
}

// Wrap wraps nodeSyncError around another error
//
// This function is auto-generated.
//
func (e *nodeSyncError) Wrap(err error) *nodeSyncError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *nodeSyncError) Unwrap() error {
	return e.wrap
}

func (e *nodeSyncError) LoggableAction() *actionlog.Action {
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

// NodeSyncActionLookup returns "federation:node_sync.lookup" error
//
// This function is auto-generated.
//
func NodeSyncActionLookup(props ...*nodeSyncActionProps) *nodeSyncAction {
	a := &nodeSyncAction{
		timestamp: time.Now(),
		resource:  "federation:node_sync",
		action:    "lookup",
		log:       "looked-up for the last successful sync",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeSyncActionCreate returns "federation:node_sync.create" error
//
// This function is auto-generated.
//
func NodeSyncActionCreate(props ...*nodeSyncActionProps) *nodeSyncAction {
	a := &nodeSyncAction{
		timestamp: time.Now(),
		resource:  "federation:node_sync",
		action:    "create",
		log:       "created node_sync",
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

// NodeSyncErrGeneric returns "federation:node_sync.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NodeSyncErrGeneric(props ...*nodeSyncActionProps) *nodeSyncError {
	var e = &nodeSyncError{
		timestamp: time.Now(),
		resource:  "federation:node_sync",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *nodeSyncActionProps {
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

// NodeSyncErrNotFound returns "federation:node_sync.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func NodeSyncErrNotFound(props ...*nodeSyncActionProps) *nodeSyncError {
	var e = &nodeSyncError{
		timestamp: time.Now(),
		resource:  "federation:node_sync",
		error:     "notFound",
		action:    "error",
		message:   "node_sync does not exist",
		log:       "node_sync does not exist",
		severity:  actionlog.Warning,
		props: func() *nodeSyncActionProps {
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

// NodeSyncErrNodeNotFound returns "federation:node_sync.nodeNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func NodeSyncErrNodeNotFound(props ...*nodeSyncActionProps) *nodeSyncError {
	var e = &nodeSyncError{
		timestamp: time.Now(),
		resource:  "federation:node_sync",
		error:     "nodeNotFound",
		action:    "error",
		message:   "node does not exist",
		log:       "node does not exist",
		severity:  actionlog.Warning,
		props: func() *nodeSyncActionProps {
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
// action (optional) fn will be used to construct nodeSyncAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc nodeSync) recordAction(ctx context.Context, props *nodeSyncActionProps, action func(...*nodeSyncActionProps) *nodeSyncAction, err error) error {
	var (
		ok bool

		// Return error
		retError *nodeSyncError

		// Recorder error
		recError *nodeSyncError
	)

	if err != nil {
		if retError, ok = err.(*nodeSyncError); !ok {
			// got non-nodeSync error, wrap it with NodeSyncErrGeneric
			retError = NodeSyncErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use NodeSyncErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type nodeSyncError
				if unwrappedSinkError, ok := unwrappedError.(*nodeSyncError); ok {
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
