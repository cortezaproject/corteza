package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/apigw_function_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
	"time"
)

type (
	apigwFunctionActionProps struct {
		function *types.ApigwFunction
		search   *types.ApigwFunctionFilter
	}

	apigwFunctionAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *apigwFunctionActionProps
	}

	apigwFunctionLogMetaKey   struct{}
	apigwFunctionPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setFunction updates apigwFunctionActionProps's function
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *apigwFunctionActionProps) setFunction(function *types.ApigwFunction) *apigwFunctionActionProps {
	p.function = function
	return p
}

// setSearch updates apigwFunctionActionProps's search
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *apigwFunctionActionProps) setSearch(search *types.ApigwFunctionFilter) *apigwFunctionActionProps {
	p.search = search
	return p
}

// Serialize converts apigwFunctionActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p apigwFunctionActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.function != nil {
		m.Set("function.ID", p.function.ID, true)
		m.Set("function.ref", p.function.Ref, true)
	}
	if p.search != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p apigwFunctionActionProps) Format(in string, err error) string {
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

	if p.function != nil {
		// replacement for "{function}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{function}",
			fns(
				p.function.ID,
				p.function.Ref,
			),
		)
		pairs = append(pairs, "{function.ID}", fns(p.function.ID))
		pairs = append(pairs, "{function.ref}", fns(p.function.Ref))
	}

	if p.search != nil {
		// replacement for "{search}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{search}",
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
func (a *apigwFunctionAction) String() string {
	var props = &apigwFunctionActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *apigwFunctionAction) ToAction() *actionlog.Action {
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

// ApigwFunctionActionSearch returns "system:function.search" action
//
// This function is auto-generated.
//
func ApigwFunctionActionSearch(props ...*apigwFunctionActionProps) *apigwFunctionAction {
	a := &apigwFunctionAction{
		timestamp: time.Now(),
		resource:  "system:function",
		action:    "search",
		log:       "searched for function",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFunctionActionLookup returns "system:function.lookup" action
//
// This function is auto-generated.
//
func ApigwFunctionActionLookup(props ...*apigwFunctionActionProps) *apigwFunctionAction {
	a := &apigwFunctionAction{
		timestamp: time.Now(),
		resource:  "system:function",
		action:    "lookup",
		log:       "looked-up for a {function}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFunctionActionCreate returns "system:function.create" action
//
// This function is auto-generated.
//
func ApigwFunctionActionCreate(props ...*apigwFunctionActionProps) *apigwFunctionAction {
	a := &apigwFunctionAction{
		timestamp: time.Now(),
		resource:  "system:function",
		action:    "create",
		log:       "created {function}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFunctionActionUpdate returns "system:function.update" action
//
// This function is auto-generated.
//
func ApigwFunctionActionUpdate(props ...*apigwFunctionActionProps) *apigwFunctionAction {
	a := &apigwFunctionAction{
		timestamp: time.Now(),
		resource:  "system:function",
		action:    "update",
		log:       "updated {function}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFunctionActionDelete returns "system:function.delete" action
//
// This function is auto-generated.
//
func ApigwFunctionActionDelete(props ...*apigwFunctionActionProps) *apigwFunctionAction {
	a := &apigwFunctionAction{
		timestamp: time.Now(),
		resource:  "system:function",
		action:    "delete",
		log:       "deleted {function}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFunctionActionUndelete returns "system:function.undelete" action
//
// This function is auto-generated.
//
func ApigwFunctionActionUndelete(props ...*apigwFunctionActionProps) *apigwFunctionAction {
	a := &apigwFunctionAction{
		timestamp: time.Now(),
		resource:  "system:function",
		action:    "undelete",
		log:       "undeleted {function}",
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

// ApigwFunctionErrGeneric returns "system:function.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFunctionErrGeneric(mm ...*apigwFunctionActionProps) *errors.Error {
	var p = &apigwFunctionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:function"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFunctionLogMetaKey{}, "{err}"),
		errors.Meta(apigwFunctionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFunctionErrNotFound returns "system:function.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFunctionErrNotFound(mm ...*apigwFunctionActionProps) *errors.Error {
	var p = &apigwFunctionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("function not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:function"),

		errors.Meta(apigwFunctionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFunctionErrInvalidID returns "system:function.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFunctionErrInvalidID(mm ...*apigwFunctionActionProps) *errors.Error {
	var p = &apigwFunctionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:function"),

		errors.Meta(apigwFunctionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFunctionErrInvalidRoute returns "system:function.invalidRoute" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFunctionErrInvalidRoute(mm ...*apigwFunctionActionProps) *errors.Error {
	var p = &apigwFunctionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid route", nil),

		errors.Meta("type", "invalidRoute"),
		errors.Meta("resource", "system:function"),

		errors.Meta(apigwFunctionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFunctionErrNotAllowedToCreate returns "system:function.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFunctionErrNotAllowedToCreate(mm ...*apigwFunctionActionProps) *errors.Error {
	var p = &apigwFunctionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a function", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:function"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFunctionLogMetaKey{}, "failed to create a route; insufficient permissions"),
		errors.Meta(apigwFunctionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFunctionErrNotAllowedToRead returns "system:function.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFunctionErrNotAllowedToRead(mm ...*apigwFunctionActionProps) *errors.Error {
	var p = &apigwFunctionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this function", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:function"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFunctionLogMetaKey{}, "failed to read {function}; insufficient permissions"),
		errors.Meta(apigwFunctionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFunctionErrNotAllowedToUpdate returns "system:function.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFunctionErrNotAllowedToUpdate(mm ...*apigwFunctionActionProps) *errors.Error {
	var p = &apigwFunctionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this function", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:function"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFunctionLogMetaKey{}, "failed to update {function}; insufficient permissions"),
		errors.Meta(apigwFunctionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFunctionErrNotAllowedToDelete returns "system:function.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFunctionErrNotAllowedToDelete(mm ...*apigwFunctionActionProps) *errors.Error {
	var p = &apigwFunctionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this function", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:function"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFunctionLogMetaKey{}, "failed to delete {function}; insufficient permissions"),
		errors.Meta(apigwFunctionPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFunctionErrNotAllowedToUndelete returns "system:function.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFunctionErrNotAllowedToUndelete(mm ...*apigwFunctionActionProps) *errors.Error {
	var p = &apigwFunctionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this function", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:function"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFunctionLogMetaKey{}, "failed to undelete {function}; insufficient permissions"),
		errors.Meta(apigwFunctionPropsMetaKey{}, p),

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
func (svc apigwFunction) recordAction(ctx context.Context, props *apigwFunctionActionProps, actionFn func(...*apigwFunctionActionProps) *apigwFunctionAction, err error) error {
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
		a.Description = props.Format(m.AsString(apigwFunctionLogMetaKey{}), err)

		if p, has := m[apigwFunctionPropsMetaKey{}]; has {
			a.Meta = p.(*apigwFunctionActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
