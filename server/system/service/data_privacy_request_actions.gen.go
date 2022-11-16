package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/data_privacy_request_actions.yaml

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
	dataPrivacyActionProps struct {
		dataPrivacyRequest *types.DataPrivacyRequest
		new                *types.DataPrivacyRequest
		update             *types.DataPrivacyRequest
		filter             *types.DataPrivacyRequestFilter
	}

	dataPrivacyAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *dataPrivacyActionProps
	}

	dataPrivacyLogMetaKey   struct{}
	dataPrivacyPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setDataPrivacyRequest updates dataPrivacyActionProps's dataPrivacyRequest
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *dataPrivacyActionProps) setDataPrivacyRequest(dataPrivacyRequest *types.DataPrivacyRequest) *dataPrivacyActionProps {
	p.dataPrivacyRequest = dataPrivacyRequest
	return p
}

// setNew updates dataPrivacyActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *dataPrivacyActionProps) setNew(new *types.DataPrivacyRequest) *dataPrivacyActionProps {
	p.new = new
	return p
}

// setUpdate updates dataPrivacyActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *dataPrivacyActionProps) setUpdate(update *types.DataPrivacyRequest) *dataPrivacyActionProps {
	p.update = update
	return p
}

// setFilter updates dataPrivacyActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *dataPrivacyActionProps) setFilter(filter *types.DataPrivacyRequestFilter) *dataPrivacyActionProps {
	p.filter = filter
	return p
}

// Serialize converts dataPrivacyActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p dataPrivacyActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.dataPrivacyRequest != nil {
		m.Set("dataPrivacyRequest.kind", p.dataPrivacyRequest.Kind, true)
		m.Set("dataPrivacyRequest.ID", p.dataPrivacyRequest.ID, true)
	}
	if p.new != nil {
		m.Set("new.kind", p.new.Kind, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.kind", p.update.Kind, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.kind", p.filter.Kind, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p dataPrivacyActionProps) Format(in string, err error) string {
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

	if p.dataPrivacyRequest != nil {
		// replacement for "{{dataPrivacyRequest}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{dataPrivacyRequest}}",
			fns(
				p.dataPrivacyRequest.Kind,
				p.dataPrivacyRequest.ID,
			),
		)
		pairs = append(pairs, "{{dataPrivacyRequest.kind}}", fns(p.dataPrivacyRequest.Kind))
		pairs = append(pairs, "{{dataPrivacyRequest.ID}}", fns(p.dataPrivacyRequest.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Kind,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.kind}}", fns(p.new.Kind))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Kind,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.kind}}", fns(p.update.Kind))
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Kind,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{{filter.kind}}", fns(p.filter.Kind))
		pairs = append(pairs, "{{filter.sort}}", fns(p.filter.Sort))
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
func (a *dataPrivacyAction) String() string {
	var props = &dataPrivacyActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *dataPrivacyAction) ToAction() *actionlog.Action {
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

// DataPrivacyActionSearch returns "system:data-privacy-request.search" action
//
// This function is auto-generated.
//
func DataPrivacyActionSearch(props ...*dataPrivacyActionProps) *dataPrivacyAction {
	a := &dataPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:data-privacy-request",
		action:    "search",
		log:       "searched for data privacy requests",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DataPrivacyActionLookup returns "system:data-privacy-request.lookup" action
//
// This function is auto-generated.
//
func DataPrivacyActionLookup(props ...*dataPrivacyActionProps) *dataPrivacyAction {
	a := &dataPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:data-privacy-request",
		action:    "lookup",
		log:       "looked-up for a {{dataPrivacyRequest}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DataPrivacyActionCreate returns "system:data-privacy-request.create" action
//
// This function is auto-generated.
//
func DataPrivacyActionCreate(props ...*dataPrivacyActionProps) *dataPrivacyAction {
	a := &dataPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:data-privacy-request",
		action:    "create",
		log:       "created {{dataPrivacyRequest}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DataPrivacyActionUpdate returns "system:data-privacy-request.update" action
//
// This function is auto-generated.
//
func DataPrivacyActionUpdate(props ...*dataPrivacyActionProps) *dataPrivacyAction {
	a := &dataPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:data-privacy-request",
		action:    "update",
		log:       "updated {{dataPrivacyRequest}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DataPrivacyActionApprove returns "system:data-privacy-request.approve" action
//
// This function is auto-generated.
//
func DataPrivacyActionApprove(props ...*dataPrivacyActionProps) *dataPrivacyAction {
	a := &dataPrivacyAction{
		timestamp: time.Now(),
		resource:  "system:data-privacy-request",
		action:    "approve",
		log:       "approved {{dataPrivacyRequest}}",
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

// DataPrivacyErrGeneric returns "system:data-privacy-request.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func DataPrivacyErrGeneric(mm ...*dataPrivacyActionProps) *errors.Error {
	var p = &dataPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:data-privacy-request"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dataPrivacyLogMetaKey{}, "{err}"),
		errors.Meta(dataPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "data-privacy.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DataPrivacyErrNotFound returns "system:data-privacy-request.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func DataPrivacyErrNotFound(mm ...*dataPrivacyActionProps) *errors.Error {
	var p = &dataPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("data privacy request not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:data-privacy-request"),

		errors.Meta(dataPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "data-privacy.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DataPrivacyErrInvalidID returns "system:data-privacy-request.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func DataPrivacyErrInvalidID(mm ...*dataPrivacyActionProps) *errors.Error {
	var p = &dataPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:data-privacy-request"),

		errors.Meta(dataPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "data-privacy.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DataPrivacyErrInvalidKind returns "system:data-privacy-request.invalidKind" as *errors.Error
//
//
// This function is auto-generated.
//
func DataPrivacyErrInvalidKind(mm ...*dataPrivacyActionProps) *errors.Error {
	var p = &dataPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid Kind", nil),

		errors.Meta("type", "invalidKind"),
		errors.Meta("resource", "system:data-privacy-request"),

		errors.Meta(dataPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "data-privacy.errors.invalidKind"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DataPrivacyErrInvalidStatus returns "system:data-privacy-request.invalidStatus" as *errors.Error
//
//
// This function is auto-generated.
//
func DataPrivacyErrInvalidStatus(mm ...*dataPrivacyActionProps) *errors.Error {
	var p = &dataPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid Status", nil),

		errors.Meta("type", "invalidStatus"),
		errors.Meta("resource", "system:data-privacy-request"),

		errors.Meta(dataPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "data-privacy.errors.invalidStatus"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DataPrivacyErrNotAllowedToRead returns "system:data-privacy-request.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func DataPrivacyErrNotAllowedToRead(mm ...*dataPrivacyActionProps) *errors.Error {
	var p = &dataPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this data privacy request", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:data-privacy-request"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dataPrivacyLogMetaKey{}, "failed to read data privacy request; insufficient permissions"),
		errors.Meta(dataPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "data-privacy.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DataPrivacyErrNotAllowedToSearch returns "system:data-privacy-request.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func DataPrivacyErrNotAllowedToSearch(mm ...*dataPrivacyActionProps) *errors.Error {
	var p = &dataPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list data privacy request", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:data-privacy-request"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dataPrivacyLogMetaKey{}, "failed to search or list data privacy requests; insufficient permissions"),
		errors.Meta(dataPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "data-privacy.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DataPrivacyErrNotAllowedToCreate returns "system:data-privacy-request.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func DataPrivacyErrNotAllowedToCreate(mm ...*dataPrivacyActionProps) *errors.Error {
	var p = &dataPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create data privacy request", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:data-privacy-request"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dataPrivacyLogMetaKey{}, "failed to create data privacy request; insufficient permissions"),
		errors.Meta(dataPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "data-privacy.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DataPrivacyErrNotAllowedToApprove returns "system:data-privacy-request.notAllowedToApprove" as *errors.Error
//
//
// This function is auto-generated.
//
func DataPrivacyErrNotAllowedToApprove(mm ...*dataPrivacyActionProps) *errors.Error {
	var p = &dataPrivacyActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to approve/reject data privacy request", nil),

		errors.Meta("type", "notAllowedToApprove"),
		errors.Meta("resource", "system:data-privacy-request"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dataPrivacyLogMetaKey{}, "failed to approve/reject data privacy request; insufficient permissions"),
		errors.Meta(dataPrivacyPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "data-privacy.errors.notAllowedToApprove"),

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
func (svc dataPrivacy) recordAction(ctx context.Context, props *dataPrivacyActionProps, actionFn func(...*dataPrivacyActionProps) *dataPrivacyAction, err error) error {
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
		a.Description = props.Format(m.AsString(dataPrivacyLogMetaKey{}), err)

		if p, has := m[dataPrivacyPropsMetaKey{}]; has {
			a.Meta = p.(*dataPrivacyActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
