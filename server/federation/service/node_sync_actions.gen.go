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
	"fmt"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"strings"
	"time"
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

	nodeSyncLogMetaKey   struct{}
	nodeSyncPropsMetaKey struct{}
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

// Serialize converts nodeSyncActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p nodeSyncActionProps) Serialize() actionlog.Meta {
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
func (p nodeSyncActionProps) Format(in string, err error) string {
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

	if p.nodeSync != nil {
		// replacement for "{{nodeSync}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{nodeSync}}",
			fns(
				p.nodeSync.NodeID,
				p.nodeSync.SyncStatus,
				p.nodeSync.SyncType,
				p.nodeSync.TimeOfAction,
			),
		)
		pairs = append(pairs, "{{nodeSync.NodeID}}", fns(p.nodeSync.NodeID))
		pairs = append(pairs, "{{nodeSync.SyncStatus}}", fns(p.nodeSync.SyncStatus))
		pairs = append(pairs, "{{nodeSync.SyncType}}", fns(p.nodeSync.SyncType))
		pairs = append(pairs, "{{nodeSync.TimeOfAction}}", fns(p.nodeSync.TimeOfAction))
	}

	if p.nodeSyncFilter != nil {
		// replacement for "{{nodeSyncFilter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{nodeSyncFilter}}",
			fns(
				p.nodeSyncFilter.Query,
			),
		)
		pairs = append(pairs, "{{nodeSyncFilter.Query}}", fns(p.nodeSyncFilter.Query))
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

	return props.Format(a.log, nil)
}

func (e *nodeSyncAction) ToAction() *actionlog.Action {
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

// NodeSyncActionLookup returns "federation:node_sync.lookup" action
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

// NodeSyncActionCreate returns "federation:node_sync.create" action
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

// NodeSyncErrGeneric returns "federation:node_sync.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeSyncErrGeneric(mm ...*nodeSyncActionProps) *errors.Error {
	var p = &nodeSyncActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "federation:node_sync"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(nodeSyncLogMetaKey{}, "{err}"),
		errors.Meta(nodeSyncPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node-sync.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeSyncErrNotFound returns "federation:node_sync.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeSyncErrNotFound(mm ...*nodeSyncActionProps) *errors.Error {
	var p = &nodeSyncActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("node_sync does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "federation:node_sync"),

		errors.Meta(nodeSyncPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node-sync.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeSyncErrNodeNotFound returns "federation:node_sync.nodeNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeSyncErrNodeNotFound(mm ...*nodeSyncActionProps) *errors.Error {
	var p = &nodeSyncActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("node does not exist", nil),

		errors.Meta("type", "nodeNotFound"),
		errors.Meta("resource", "federation:node_sync"),

		errors.Meta(nodeSyncPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node-sync.errors.nodeNotFound"),

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
func (svc nodeSync) recordAction(ctx context.Context, props *nodeSyncActionProps, actionFn func(...*nodeSyncActionProps) *nodeSyncAction, err error) error {
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
		a.Description = props.Format(m.AsString(nodeSyncLogMetaKey{}), err)

		if p, has := m[nodeSyncPropsMetaKey{}]; has {
			a.Meta = p.(*nodeSyncActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
