package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/apigw_route_actions.yaml

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
	apigwRouteActionProps struct {
		route  *types.ApigwRoute
		new    *types.ApigwRoute
		update *types.ApigwRoute
		search *types.ApigwRouteFilter
	}

	apigwRouteAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *apigwRouteActionProps
	}

	apigwRouteLogMetaKey   struct{}
	apigwRoutePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setRoute updates apigwRouteActionProps's route
//
// This function is auto-generated.
//
func (p *apigwRouteActionProps) setRoute(route *types.ApigwRoute) *apigwRouteActionProps {
	p.route = route
	return p
}

// setNew updates apigwRouteActionProps's new
//
// This function is auto-generated.
//
func (p *apigwRouteActionProps) setNew(new *types.ApigwRoute) *apigwRouteActionProps {
	p.new = new
	return p
}

// setUpdate updates apigwRouteActionProps's update
//
// This function is auto-generated.
//
func (p *apigwRouteActionProps) setUpdate(update *types.ApigwRoute) *apigwRouteActionProps {
	p.update = update
	return p
}

// setSearch updates apigwRouteActionProps's search
//
// This function is auto-generated.
//
func (p *apigwRouteActionProps) setSearch(search *types.ApigwRouteFilter) *apigwRouteActionProps {
	p.search = search
	return p
}

// Serialize converts apigwRouteActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p apigwRouteActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.route != nil {
		m.Set("route.endpoint", p.route.Endpoint, true)
		m.Set("route.ID", p.route.ID, true)
	}
	if p.new != nil {
		m.Set("new.endpoint", p.new.Endpoint, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.endpoint", p.update.Endpoint, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.search != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p apigwRouteActionProps) Format(in string, err error) string {
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

	if p.route != nil {
		// replacement for "{{route}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{route}}",
			fns(
				p.route.Endpoint,
				p.route.ID,
			),
		)
		pairs = append(pairs, "{{route.endpoint}}", fns(p.route.Endpoint))
		pairs = append(pairs, "{{route.ID}}", fns(p.route.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Endpoint,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.endpoint}}", fns(p.new.Endpoint))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Endpoint,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.endpoint}}", fns(p.update.Endpoint))
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
//
func (a *apigwRouteAction) String() string {
	var props = &apigwRouteActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *apigwRouteAction) ToAction() *actionlog.Action {
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

// ApigwRouteActionSearch returns "system:apigw-route.search" action
//
// This function is auto-generated.
//
func ApigwRouteActionSearch(props ...*apigwRouteActionProps) *apigwRouteAction {
	a := &apigwRouteAction{
		timestamp: time.Now(),
		resource:  "system:apigw-route",
		action:    "search",
		log:       "searched for route",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwRouteActionLookup returns "system:apigw-route.lookup" action
//
// This function is auto-generated.
//
func ApigwRouteActionLookup(props ...*apigwRouteActionProps) *apigwRouteAction {
	a := &apigwRouteAction{
		timestamp: time.Now(),
		resource:  "system:apigw-route",
		action:    "lookup",
		log:       "looked-up for a {{route}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwRouteActionCreate returns "system:apigw-route.create" action
//
// This function is auto-generated.
//
func ApigwRouteActionCreate(props ...*apigwRouteActionProps) *apigwRouteAction {
	a := &apigwRouteAction{
		timestamp: time.Now(),
		resource:  "system:apigw-route",
		action:    "create",
		log:       "created {{route}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwRouteActionUpdate returns "system:apigw-route.update" action
//
// This function is auto-generated.
//
func ApigwRouteActionUpdate(props ...*apigwRouteActionProps) *apigwRouteAction {
	a := &apigwRouteAction{
		timestamp: time.Now(),
		resource:  "system:apigw-route",
		action:    "update",
		log:       "updated {{route}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwRouteActionDelete returns "system:apigw-route.delete" action
//
// This function is auto-generated.
//
func ApigwRouteActionDelete(props ...*apigwRouteActionProps) *apigwRouteAction {
	a := &apigwRouteAction{
		timestamp: time.Now(),
		resource:  "system:apigw-route",
		action:    "delete",
		log:       "deleted {{route}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApigwRouteActionUndelete returns "system:apigw-route.undelete" action
//
// This function is auto-generated.
//
func ApigwRouteActionUndelete(props ...*apigwRouteActionProps) *apigwRouteAction {
	a := &apigwRouteAction{
		timestamp: time.Now(),
		resource:  "system:apigw-route",
		action:    "undelete",
		log:       "undeleted {{route}}",
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

// ApigwRouteErrGeneric returns "system:apigw-route.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrGeneric(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:apigw-route"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwRouteLogMetaKey{}, "{err}"),
		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrNotFound returns "system:apigw-route.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrNotFound(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("route not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:apigw-route"),

		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrInvalidID returns "system:apigw-route.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrInvalidID(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:apigw-route"),

		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrInvalidEndpoint returns "system:apigw-route.invalidEndpoint" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrInvalidEndpoint(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid endpoint", nil),

		errors.Meta("type", "invalidEndpoint"),
		errors.Meta("resource", "system:apigw-route"),

		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.invalidEndpoint"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrExistsEndpoint returns "system:apigw-route.existsEndpoint" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrExistsEndpoint(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("route with this endpoint already exists", nil),

		errors.Meta("type", "existsEndpoint"),
		errors.Meta("resource", "system:apigw-route"),

		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.existsEndpoint"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrAlreadyExists returns "system:apigw-route.alreadyExists" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrAlreadyExists(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("route by that endpoint already exists", nil),

		errors.Meta("type", "alreadyExists"),
		errors.Meta("resource", "system:apigw-route"),

		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.alreadyExists"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrNotAllowedToCreate returns "system:apigw-route.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrNotAllowedToCreate(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a route", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:apigw-route"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwRouteLogMetaKey{}, "failed to create a route; insufficient permissions"),
		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrNotAllowedToRead returns "system:apigw-route.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrNotAllowedToRead(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this route", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:apigw-route"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwRouteLogMetaKey{}, "failed to read {{route.endpoint}}; insufficient permissions"),
		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrNotAllowedToSearch returns "system:apigw-route.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrNotAllowedToSearch(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list or search routes", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:apigw-route"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwRouteLogMetaKey{}, "failed to search for routes; insufficient permissions"),
		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrNotAllowedToUpdate returns "system:apigw-route.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrNotAllowedToUpdate(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this route", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:apigw-route"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwRouteLogMetaKey{}, "failed to update {{route.endpoint}}; insufficient permissions"),
		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrNotAllowedToDelete returns "system:apigw-route.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrNotAllowedToDelete(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this route", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:apigw-route"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwRouteLogMetaKey{}, "failed to delete {{route.endpoint}}; insufficient permissions"),
		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrNotAllowedToUndelete returns "system:apigw-route.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrNotAllowedToUndelete(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this route", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:apigw-route"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwRouteLogMetaKey{}, "failed to undelete {{route.endpoint}}; insufficient permissions"),
		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApigwRouteErrNotAllowedToExec returns "system:apigw-route.notAllowedToExec" as *errors.Error
//
//
// This function is auto-generated.
//
func ApigwRouteErrNotAllowedToExec(mm ...*apigwRouteActionProps) *errors.Error {
	var p = &apigwRouteActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to execute this route", nil),

		errors.Meta("type", "notAllowedToExec"),
		errors.Meta("resource", "system:apigw-route"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(apigwRouteLogMetaKey{}, "failed to exec {{route.endpoint}}; insufficient permissions"),
		errors.Meta(apigwRoutePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "apigw-route.errors.notAllowedToExec"),

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
func (svc apigwRoute) recordAction(ctx context.Context, props *apigwRouteActionProps, actionFn func(...*apigwRouteActionProps) *apigwRouteAction, err error) error {
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
		a.Description = props.Format(m.AsString(apigwRouteLogMetaKey{}), err)

		if p, has := m[apigwRoutePropsMetaKey{}]; has {
			a.Meta = p.(*apigwRouteActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
