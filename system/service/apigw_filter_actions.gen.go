package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/apigw_filter_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
	"time"
)

type (
	apigwFilterActionProps struct {
		filter *types.ApigwFilter
		search *types.ApigwFilterFilter
	}

	apigwFilterAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *apigwFilterActionProps
	}

	apigwFilterLogMetaKey   struct{}
	apigwFilterPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setFilter updates apigwFilterActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *apigwFilterActionProps) setFilter(filter *types.ApigwFilter) *apigwFilterActionProps {
	p.filter = filter
	return p
}

// setSearch updates apigwFilterActionProps's search
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *apigwFilterActionProps) setSearch(search *types.ApigwFilterFilter) *apigwFilterActionProps {
	p.search = search
	return p
}

// Serialize converts apigwFilterActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p apigwFilterActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.filter != nil {
		m.Set("filter.ID", p.filter.ID, true)
		m.Set("filter.ref", p.filter.Ref, true)
	}
	if p.search != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p apigwFilterActionProps) Format(in string, err error) string {
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

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.ID,
				p.filter.Ref,
			),
		)
		pairs = append(pairs, "{{filter.ID}}", fns(p.filter.ID))
		pairs = append(pairs, "{{filter.ref}}", fns(p.filter.Ref))
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
func (a *apigwFilterAction) String() string {
	var props = &apigwFilterActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *apigwFilterAction) ToAction() *actionlog.Action {
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

// ApigwFilterActionSearch returns "system:filter.search" action
//
// This function is auto-generated.
//
func ApigwFilterActionSearch(props ...*apigwFilterActionProps) *apigwFilterAction {
	a := &apigwFilterAction{
		timestamp: time.Now(),
		resource:  "system:filter",
		action:    "search",
		log:       "searched for filter",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFilterActionLookup returns "system:filter.lookup" action
//
// This function is auto-generated.
//
func ApigwFilterActionLookup(props ...*apigwFilterActionProps) *apigwFilterAction {
	a := &apigwFilterAction{
		timestamp: time.Now(),
		resource:  "system:filter",
		action:    "lookup",
		log:       "looked-up for a {{filter}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFilterActionCreate returns "system:filter.create" action
//
// This function is auto-generated.
//
func ApigwFilterActionCreate(props ...*apigwFilterActionProps) *apigwFilterAction {
	a := &apigwFilterAction{
		timestamp: time.Now(),
		resource:  "system:filter",
		action:    "create",
		log:       "created {{filter}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFilterActionUpdate returns "system:filter.update" action
//
// This function is auto-generated.
//
func ApigwFilterActionUpdate(props ...*apigwFilterActionProps) *apigwFilterAction {
	a := &apigwFilterAction{
		timestamp: time.Now(),
		resource:  "system:filter",
		action:    "update",
		log:       "updated {{filter}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFilterActionDelete returns "system:filter.delete" action
//
// This function is auto-generated.
//
func ApigwFilterActionDelete(props ...*apigwFilterActionProps) *apigwFilterAction {
	a := &apigwFilterAction{
		timestamp: time.Now(),
		resource:  "system:filter",
		action:    "delete",
		log:       "deleted {{filter}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwFilterActionUndelete returns "system:filter.undelete" action
//
// This function is auto-generated.
//
func ApigwFilterActionUndelete(props ...*apigwFilterActionProps) *apigwFilterAction {
	a := &apigwFilterAction{
		timestamp: time.Now(),
		resource:  "system:filter",
		action:    "undelete",
		log:       "undeleted {{filter}}",
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

// ApigwFilterErrGeneric returns "system:filter.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrGeneric(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:filter"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFilterLogMetaKey{}, "{err}"),
		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrNotFound returns "system:filter.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrNotFound(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("filter not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:filter"),

		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrInvalidID returns "system:filter.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrInvalidID(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:filter"),

		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrInvalidRoute returns "system:filter.invalidRoute" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrInvalidRoute(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid route", nil),

		errors.Meta("type", "invalidRoute"),
		errors.Meta("resource", "system:filter"),

		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.invalidRoute"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrNotAllowedToCreate returns "system:filter.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrNotAllowedToCreate(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a filter", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:filter"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFilterLogMetaKey{}, "failed to create a route; insufficient permissions"),
		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrNotAllowedToRead returns "system:filter.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrNotAllowedToRead(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this filter", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:filter"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFilterLogMetaKey{}, "failed to read {{filter}}; insufficient permissions"),
		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrNotAllowedToUpdate returns "system:filter.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrNotAllowedToUpdate(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this filter", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:filter"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFilterLogMetaKey{}, "failed to update {{filter}}; insufficient permissions"),
		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrNotAllowedToDelete returns "system:filter.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrNotAllowedToDelete(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this filter", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:filter"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFilterLogMetaKey{}, "failed to delete {{filter}}; insufficient permissions"),
		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrNotAllowedToUndelete returns "system:filter.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrNotAllowedToUndelete(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this filter", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:filter"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFilterLogMetaKey{}, "failed to undelete {{filter}}; insufficient permissions"),
		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrAsyncRouteTooManyProcessers returns "system:filter.asyncRouteTooManyProcessers" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrAsyncRouteTooManyProcessers(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("processer already exists for this async route", nil),

		errors.Meta("type", "asyncRouteTooManyProcessers"),
		errors.Meta("resource", "system:filter"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFilterLogMetaKey{}, "failed to add {{filter}}; too many processers, async route"),
		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.asyncRouteTooManyProcessers"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwFilterErrAsyncRouteTooManyAfterFilters returns "system:filter.asyncRouteTooManyAfterFilters" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwFilterErrAsyncRouteTooManyAfterFilters(mm ...*apigwFilterActionProps) *errors.Error {
	var p = &apigwFilterActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("no after filters are allowd for this async route", nil),

		errors.Meta("type", "asyncRouteTooManyAfterFilters"),
		errors.Meta("resource", "system:filter"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwFilterLogMetaKey{}, "failed to add {{filter}}; too many afterfilters, async route"),
		errors.Meta(apigwFilterPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigwFilter.errors.asyncRouteTooManyAfterFilters"),

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
func (svc apigwFilter) recordAction(ctx context.Context, props *apigwFilterActionProps, actionFn func(...*apigwFilterActionProps) *apigwFilterAction, err error) error {
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
		a.Description = props.Format(m.AsString(apigwFilterLogMetaKey{}), err)

		if p, has := m[apigwFilterPropsMetaKey{}]; has {
			a.Meta = p.(*apigwFilterActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
