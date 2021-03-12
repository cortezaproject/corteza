package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/service/session_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"strings"
	"time"
)

type (
	sessionActionProps struct {
		session *types.Session
		new     *types.Session
		update  *types.Session
		filter  *types.SessionFilter
	}

	sessionAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *sessionActionProps
	}

	sessionLogMetaKey   struct{}
	sessionPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setSession updates sessionActionProps's session
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sessionActionProps) setSession(session *types.Session) *sessionActionProps {
	p.session = session
	return p
}

// setNew updates sessionActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sessionActionProps) setNew(new *types.Session) *sessionActionProps {
	p.new = new
	return p
}

// setUpdate updates sessionActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sessionActionProps) setUpdate(update *types.Session) *sessionActionProps {
	p.update = update
	return p
}

// setFilter updates sessionActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sessionActionProps) setFilter(filter *types.SessionFilter) *sessionActionProps {
	p.filter = filter
	return p
}

// Serialize converts sessionActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p sessionActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.session != nil {
		m.Set("session.ID", p.session.ID, true)
	}
	if p.new != nil {
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p sessionActionProps) Format(in string, err error) string {
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
		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}

	if p.session != nil {
		// replacement for "{session}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{session}",
			fns(
				p.session.ID,
			),
		)
		pairs = append(pairs, "{session.ID}", fns(p.session.ID))
	}

	if p.new != nil {
		// replacement for "{new}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{new}",
			fns(
				p.new.ID,
			),
		)
		pairs = append(pairs, "{new.ID}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{update}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{update}",
			fns(
				p.update.ID,
			),
		)
		pairs = append(pairs, "{update.ID}", fns(p.update.ID))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
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
func (a *sessionAction) String() string {
	var props = &sessionActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *sessionAction) ToAction() *actionlog.Action {
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

// SessionActionSearch returns "automation:session.search" action
//
// This function is auto-generated.
//
func SessionActionSearch(props ...*sessionActionProps) *sessionAction {
	a := &sessionAction{
		timestamp: time.Now(),
		resource:  "automation:session",
		action:    "search",
		log:       "searched for matching sessions",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SessionActionLookup returns "automation:session.lookup" action
//
// This function is auto-generated.
//
func SessionActionLookup(props ...*sessionActionProps) *sessionAction {
	a := &sessionAction{
		timestamp: time.Now(),
		resource:  "automation:session",
		action:    "lookup",
		log:       "looked-up for a {session}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SessionActionCreate returns "automation:session.create" action
//
// This function is auto-generated.
//
func SessionActionCreate(props ...*sessionActionProps) *sessionAction {
	a := &sessionAction{
		timestamp: time.Now(),
		resource:  "automation:session",
		action:    "create",
		log:       "created {session}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SessionActionUpdate returns "automation:session.update" action
//
// This function is auto-generated.
//
func SessionActionUpdate(props ...*sessionActionProps) *sessionAction {
	a := &sessionAction{
		timestamp: time.Now(),
		resource:  "automation:session",
		action:    "update",
		log:       "updated {session}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SessionActionDelete returns "automation:session.delete" action
//
// This function is auto-generated.
//
func SessionActionDelete(props ...*sessionActionProps) *sessionAction {
	a := &sessionAction{
		timestamp: time.Now(),
		resource:  "automation:session",
		action:    "delete",
		log:       "deleted {session}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SessionActionUndelete returns "automation:session.undelete" action
//
// This function is auto-generated.
//
func SessionActionUndelete(props ...*sessionActionProps) *sessionAction {
	a := &sessionAction{
		timestamp: time.Now(),
		resource:  "automation:session",
		action:    "undelete",
		log:       "undeleted {session}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// SessionErrGeneric returns "automation:session.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func SessionErrGeneric(mm ...*sessionActionProps) *errors.Error {
	var p = &sessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "automation:session"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sessionLogMetaKey{}, "{err}"),
		errors.Meta(sessionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SessionErrNotFound returns "automation:session.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func SessionErrNotFound(mm ...*sessionActionProps) *errors.Error {
	var p = &sessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("session not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "automation:session"),

		errors.Meta(sessionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SessionErrInvalidID returns "automation:session.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func SessionErrInvalidID(mm ...*sessionActionProps) *errors.Error {
	var p = &sessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "automation:session"),

		errors.Meta(sessionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SessionErrStaleData returns "automation:session.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func SessionErrStaleData(mm ...*sessionActionProps) *errors.Error {
	var p = &sessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "automation:session"),

		errors.Meta(sessionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SessionErrNotAllowedToRead returns "automation:session.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func SessionErrNotAllowedToRead(mm ...*sessionActionProps) *errors.Error {
	var p = &sessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this session", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "automation:session"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sessionLogMetaKey{}, "failed to read {session}; insufficient permissions"),
		errors.Meta(sessionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SessionErrNotAllowedToSearch returns "automation:session.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func SessionErrNotAllowedToSearch(mm ...*sessionActionProps) *errors.Error {
	var p = &sessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search sessions", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "automation:session"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sessionLogMetaKey{}, "failed to list session; insufficient permissions"),
		errors.Meta(sessionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SessionErrNotAllowedToDelete returns "automation:session.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func SessionErrNotAllowedToDelete(mm ...*sessionActionProps) *errors.Error {
	var p = &sessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this session", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "automation:session"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sessionLogMetaKey{}, "failed to delete {session}; insufficient permissions"),
		errors.Meta(sessionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SessionErrNotAllowedToManage returns "automation:session.notAllowedToManage" as *errors.Error
//
//
// This function is auto-generated.
//
func SessionErrNotAllowedToManage(mm ...*sessionActionProps) *errors.Error {
	var p = &sessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage session's workflow", nil),

		errors.Meta("type", "notAllowedToManage"),
		errors.Meta("resource", "automation:session"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sessionLogMetaKey{}, "failed to manage {session}; insufficient permissions"),
		errors.Meta(sessionPropsMetaKey{}, p),

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
func (svc session) recordAction(ctx context.Context, props *sessionActionProps, actionFn func(...*sessionActionProps) *sessionAction, err error) error {
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
		a.Description = props.Format(m.AsString(sessionLogMetaKey{}), err)

		if p, has := m[sessionPropsMetaKey{}]; has {
			a.Meta = p.(*sessionActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
