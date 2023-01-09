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
	"fmt"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"strings"
	"time"
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

	nodeLogMetaKey   struct{}
	nodePropsMetaKey struct{}
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
// This function is auto-generated.
//
func (p *nodeActionProps) setNode(node *types.Node) *nodeActionProps {
	p.node = node
	return p
}

// setPairingURI updates nodeActionProps's pairingURI
//
// This function is auto-generated.
//
func (p *nodeActionProps) setPairingURI(pairingURI string) *nodeActionProps {
	p.pairingURI = pairingURI
	return p
}

// setFilter updates nodeActionProps's filter
//
// This function is auto-generated.
//
func (p *nodeActionProps) setFilter(filter *types.NodeFilter) *nodeActionProps {
	p.filter = filter
	return p
}

// Serialize converts nodeActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p nodeActionProps) Serialize() actionlog.Meta {
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
func (p nodeActionProps) Format(in string, err error) string {
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

	if p.node != nil {
		// replacement for "{{node}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{node}}",
			fns(
				p.node.Name,
				p.node.BaseURL,
				p.node.ID,
				p.node.Status,
			),
		)
		pairs = append(pairs, "{{node.Name}}", fns(p.node.Name))
		pairs = append(pairs, "{{node.BaseURL}}", fns(p.node.BaseURL))
		pairs = append(pairs, "{{node.ID}}", fns(p.node.ID))
		pairs = append(pairs, "{{node.Status}}", fns(p.node.Status))
	}
	pairs = append(pairs, "{{pairingURI}}", fns(p.pairingURI))

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Query,
				p.filter.Status,
			),
		)
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
//
func (a *nodeAction) String() string {
	var props = &nodeActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *nodeAction) ToAction() *actionlog.Action {
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

// NodeActionSearch returns "federation:node.search" action
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

// NodeActionLookup returns "federation:node.lookup" action
//
// This function is auto-generated.
//
func NodeActionLookup(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "lookup",
		log:       "looked-up for a {{node}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionCreate returns "federation:node.create" action
//
// This function is auto-generated.
//
func NodeActionCreate(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "create",
		log:       "created {{node}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionCreateFromPairingURI returns "federation:node.createFromPairingURI" action
//
// This function is auto-generated.
//
func NodeActionCreateFromPairingURI(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "createFromPairingURI",
		log:       "created {{node}} from pairing URI",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionRecreateFromPairingURI returns "federation:node.recreateFromPairingURI" action
//
// This function is auto-generated.
//
func NodeActionRecreateFromPairingURI(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "recreateFromPairingURI",
		log:       "recreate {{node}} from pairing URI",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionUpdate returns "federation:node.update" action
//
// This function is auto-generated.
//
func NodeActionUpdate(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "update",
		log:       "updated {{node}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionDelete returns "federation:node.delete" action
//
// This function is auto-generated.
//
func NodeActionDelete(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "delete",
		log:       "deleted {{node}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionUndelete returns "federation:node.undelete" action
//
// This function is auto-generated.
//
func NodeActionUndelete(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "undelete",
		log:       "undeleted {{node}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionOttRegenerated returns "federation:node.ottRegenerated" action
//
// This function is auto-generated.
//
func NodeActionOttRegenerated(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "ottRegenerated",
		log:       "regenerated one-time-token for {{node}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionPair returns "federation:node.pair" action
//
// This function is auto-generated.
//
func NodeActionPair(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "pair",
		log:       "{{node}} pairing started",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionHandshakeInit returns "federation:node.handshakeInit" action
//
// This function is auto-generated.
//
func NodeActionHandshakeInit(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "handshakeInit",
		log:       "{{node}} handshake initialized",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionHandshakeConfirm returns "federation:node.handshakeConfirm" action
//
// This function is auto-generated.
//
func NodeActionHandshakeConfirm(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "handshakeConfirm",
		log:       "{{node}} handshake confirmed",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NodeActionHandshakeComplete returns "federation:node.handshakeComplete" action
//
// This function is auto-generated.
//
func NodeActionHandshakeComplete(props ...*nodeActionProps) *nodeAction {
	a := &nodeAction{
		timestamp: time.Now(),
		resource:  "federation:node",
		action:    "handshakeComplete",
		log:       "{{node}} handshake completed",
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

// NodeErrGeneric returns "federation:node.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrGeneric(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "federation:node"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(nodeLogMetaKey{}, "{err}"),
		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrNotFound returns "federation:node.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrNotFound(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("node does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "federation:node"),

		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrInvalidID returns "federation:node.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrInvalidID(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "federation:node"),

		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrPairingURIInvalid returns "federation:node.pairingURIInvalid" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrPairingURIInvalid(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("pairing URI invalid: {{err}}", nil),

		errors.Meta("type", "pairingURIInvalid"),
		errors.Meta("resource", "federation:node"),

		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.pairingURIInvalid"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrPairingURITokenInvalid returns "federation:node.pairingURITokenInvalid" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrPairingURITokenInvalid(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("pairing URI with invalid pairing token", nil),

		errors.Meta("type", "pairingURITokenInvalid"),
		errors.Meta("resource", "federation:node"),

		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.pairingURITokenInvalid"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrPairingURISourceIDInvalid returns "federation:node.pairingURISourceIDInvalid" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrPairingURISourceIDInvalid(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("pairing URI without source node ID", nil),

		errors.Meta("type", "pairingURISourceIDInvalid"),
		errors.Meta("resource", "federation:node"),

		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.pairingURISourceIDInvalid"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrPairingTokenInvalid returns "federation:node.pairingTokenInvalid" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrPairingTokenInvalid(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("pairing token invalid", nil),

		errors.Meta("type", "pairingTokenInvalid"),
		errors.Meta("resource", "federation:node"),

		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.pairingTokenInvalid"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrNotAllowedToCreate returns "federation:node.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrNotAllowedToCreate(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create nodes", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "federation:node"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(nodeLogMetaKey{}, "could not create nodes; insufficient permissions"),
		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrNotAllowedToSearch returns "federation:node.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrNotAllowedToSearch(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list nodes", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "federation:node"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(nodeLogMetaKey{}, "could not search or list nodes; insufficient permissions"),
		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrNotAllowedToManage returns "federation:node.notAllowedToManage" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrNotAllowedToManage(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage this node", nil),

		errors.Meta("type", "notAllowedToManage"),
		errors.Meta("resource", "federation:node"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(nodeLogMetaKey{}, "could not manage {{node}}; insufficient permissions"),
		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.notAllowedToManage"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NodeErrNotAllowedToPair returns "federation:node.notAllowedToPair" as *errors.Error
//
//
// This function is auto-generated.
//
func NodeErrNotAllowedToPair(mm ...*nodeActionProps) *errors.Error {
	var p = &nodeActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to pair this node", nil),

		errors.Meta("type", "notAllowedToPair"),
		errors.Meta("resource", "federation:node"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(nodeLogMetaKey{}, "could not pair {{node}}; insufficient permissions"),
		errors.Meta(nodePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "node.errors.notAllowedToPair"),

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
func (svc node) recordAction(ctx context.Context, props *nodeActionProps, actionFn func(...*nodeActionProps) *nodeAction, err error) error {
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
		a.Description = props.Format(m.AsString(nodeLogMetaKey{}), err)

		if p, has := m[nodePropsMetaKey{}]; has {
			a.Meta = p.(*nodeActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
