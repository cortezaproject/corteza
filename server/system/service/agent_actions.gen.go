package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/agent_actions.yaml

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
	agentActionProps struct {
		agent  *types.Agent
		new    *types.Agent
		update *types.Agent
		search *types.AgentFilter
	}

	agentAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *agentActionProps
	}

	agentLogMetaKey   struct{}
	agentPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setAgent updates agentActionProps's agent
//
// This function is auto-generated.
func (p *agentActionProps) setAgent(agent *types.Agent) *agentActionProps {
	p.agent = agent
	return p
}

// setNew updates agentActionProps's new
//
// This function is auto-generated.
func (p *agentActionProps) setNew(new *types.Agent) *agentActionProps {
	p.new = new
	return p
}

// setUpdate updates agentActionProps's update
//
// This function is auto-generated.
func (p *agentActionProps) setUpdate(update *types.Agent) *agentActionProps {
	p.update = update
	return p
}

// setSearch updates agentActionProps's search
//
// This function is auto-generated.
func (p *agentActionProps) setSearch(search *types.AgentFilter) *agentActionProps {
	p.search = search
	return p
}

// Serialize converts agentActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p agentActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.agent != nil {
		m.Set("agent.handle", p.agent.Handle, true)
		m.Set("agent.ID", p.agent.ID, true)
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
func (p agentActionProps) Format(in string, err error) string {
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

	if p.agent != nil {
		// replacement for "{{agent}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{agent}}",
			fns(
				p.agent.Handle,
				p.agent.ID,
			),
		)
		pairs = append(pairs, "{{agent.handle}}", fns(p.agent.Handle))
		pairs = append(pairs, "{{agent.ID}}", fns(p.agent.ID))
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
func (a *agentAction) String() string {
	var props = &agentActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *agentAction) ToAction() *actionlog.Action {
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

// AgentActionSearch returns "system:agent.search" action
//
// This function is auto-generated.
func AgentActionSearch(props ...*agentActionProps) *agentAction {
	a := &agentAction{
		timestamp: time.Now(),
		resource:  "system:agent",
		action:    "search",
		log:       "searched for agents",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AgentActionLookup returns "system:agent.lookup" action
//
// This function is auto-generated.
func AgentActionLookup(props ...*agentActionProps) *agentAction {
	a := &agentAction{
		timestamp: time.Now(),
		resource:  "system:agent",
		action:    "lookup",
		log:       "looked-up for a {{agent}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AgentActionCreate returns "system:agent.create" action
//
// This function is auto-generated.
func AgentActionCreate(props ...*agentActionProps) *agentAction {
	a := &agentAction{
		timestamp: time.Now(),
		resource:  "system:agent",
		action:    "create",
		log:       "created {{agent}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AgentActionUpdate returns "system:agent.update" action
//
// This function is auto-generated.
func AgentActionUpdate(props ...*agentActionProps) *agentAction {
	a := &agentAction{
		timestamp: time.Now(),
		resource:  "system:agent",
		action:    "update",
		log:       "updated {{agent}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AgentActionDelete returns "system:agent.delete" action
//
// This function is auto-generated.
func AgentActionDelete(props ...*agentActionProps) *agentAction {
	a := &agentAction{
		timestamp: time.Now(),
		resource:  "system:agent",
		action:    "delete",
		log:       "deleted {{agent}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AgentActionUndelete returns "system:agent.undelete" action
//
// This function is auto-generated.
func AgentActionUndelete(props ...*agentActionProps) *agentAction {
	a := &agentAction{
		timestamp: time.Now(),
		resource:  "system:agent",
		action:    "undelete",
		log:       "undeleted {{agent}}",
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

// AgentErrGeneric returns "system:agent.generic" as *errors.Error
//
// This function is auto-generated.
func AgentErrGeneric(mm ...*agentActionProps) *errors.Error {
	var p = &agentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:agent"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(agentLogMetaKey{}, "{err}"),
		errors.Meta(agentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "agent.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AgentErrNotFound returns "system:agent.notFound" as *errors.Error
//
// This function is auto-generated.
func AgentErrNotFound(mm ...*agentActionProps) *errors.Error {
	var p = &agentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("agent not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:agent"),

		errors.Meta(agentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "agent.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AgentErrInvalidID returns "system:agent.invalidID" as *errors.Error
//
// This function is auto-generated.
func AgentErrInvalidID(mm ...*agentActionProps) *errors.Error {
	var p = &agentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:agent"),

		errors.Meta(agentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "agent.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AgentErrStaleData returns "system:agent.staleData" as *errors.Error
//
// This function is auto-generated.
func AgentErrStaleData(mm ...*agentActionProps) *errors.Error {
	var p = &agentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("agent was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:agent"),

		errors.Meta(agentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "agent.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AgentErrNotAllowedToCreate returns "system:agent.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func AgentErrNotAllowedToCreate(mm ...*agentActionProps) *errors.Error {
	var p = &agentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create an agent", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:agent"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(agentLogMetaKey{}, "failed to create an agent; insufficient permissions"),
		errors.Meta(agentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "agent.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AgentErrNotAllowedToRead returns "system:agent.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func AgentErrNotAllowedToRead(mm ...*agentActionProps) *errors.Error {
	var p = &agentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this agent", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:agent"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(agentLogMetaKey{}, "failed to read {{agent.handle}}; insufficient permissions"),
		errors.Meta(agentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "agent.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AgentErrNotAllowedToSearch returns "system:agent.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func AgentErrNotAllowedToSearch(mm ...*agentActionProps) *errors.Error {
	var p = &agentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list agents", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:agent"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(agentLogMetaKey{}, "failed to search or list agents; insufficient permissions"),
		errors.Meta(agentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "agent.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AgentErrNotAllowedToUpdate returns "system:agent.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func AgentErrNotAllowedToUpdate(mm ...*agentActionProps) *errors.Error {
	var p = &agentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this agent", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:agent"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(agentLogMetaKey{}, "failed to update {{agent.handle}}; insufficient permissions"),
		errors.Meta(agentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "agent.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AgentErrNotAllowedToDelete returns "system:agent.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func AgentErrNotAllowedToDelete(mm ...*agentActionProps) *errors.Error {
	var p = &agentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this agent", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:agent"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(agentLogMetaKey{}, "failed to delete {{agent.handle}}; insufficient permissions"),
		errors.Meta(agentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "agent.errors.notAllowedToDelete"),

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
func (svc agent) recordAction(ctx context.Context, props *agentActionProps, actionFn func(...*agentActionProps) *agentAction, err error) error {
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
		a.Description = props.Format(m.AsString(agentLogMetaKey{}), err)

		if p, has := m[agentPropsMetaKey{}]; has {
			a.Meta = p.(*agentActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
