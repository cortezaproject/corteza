package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/service/chart_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"strings"
	"time"
)

type (
	chartActionProps struct {
		chart     *types.Chart
		changed   *types.Chart
		filter    *types.ChartFilter
		namespace *types.Namespace
	}

	chartAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *chartActionProps
	}

	chartLogMetaKey   struct{}
	chartPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setChart updates chartActionProps's chart
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *chartActionProps) setChart(chart *types.Chart) *chartActionProps {
	p.chart = chart
	return p
}

// setChanged updates chartActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *chartActionProps) setChanged(changed *types.Chart) *chartActionProps {
	p.changed = changed
	return p
}

// setFilter updates chartActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *chartActionProps) setFilter(filter *types.ChartFilter) *chartActionProps {
	p.filter = filter
	return p
}

// setNamespace updates chartActionProps's namespace
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *chartActionProps) setNamespace(namespace *types.Namespace) *chartActionProps {
	p.namespace = namespace
	return p
}

// Serialize converts chartActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p chartActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.chart != nil {
		m.Set("chart.name", p.chart.Name, true)
		m.Set("chart.handle", p.chart.Handle, true)
		m.Set("chart.ID", p.chart.ID, true)
		m.Set("chart.namespaceID", p.chart.NamespaceID, true)
	}
	if p.changed != nil {
		m.Set("changed.name", p.changed.Name, true)
		m.Set("changed.handle", p.changed.Handle, true)
		m.Set("changed.ID", p.changed.ID, true)
		m.Set("changed.namespaceID", p.changed.NamespaceID, true)
		m.Set("changed.config", p.changed.Config, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.handle", p.filter.Handle, true)
		m.Set("filter.namespaceID", p.filter.NamespaceID, true)
		m.Set("filter.sort", p.filter.Sort, true)
		m.Set("filter.limit", p.filter.Limit, true)
	}
	if p.namespace != nil {
		m.Set("namespace.name", p.namespace.Name, true)
		m.Set("namespace.slug", p.namespace.Slug, true)
		m.Set("namespace.ID", p.namespace.ID, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p chartActionProps) Format(in string, err error) string {
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

	if p.chart != nil {
		// replacement for "{{chart}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{chart}}",
			fns(
				p.chart.Name,
				p.chart.Handle,
				p.chart.ID,
				p.chart.NamespaceID,
			),
		)
		pairs = append(pairs, "{{chart.name}}", fns(p.chart.Name))
		pairs = append(pairs, "{{chart.handle}}", fns(p.chart.Handle))
		pairs = append(pairs, "{{chart.ID}}", fns(p.chart.ID))
		pairs = append(pairs, "{{chart.namespaceID}}", fns(p.chart.NamespaceID))
	}

	if p.changed != nil {
		// replacement for "{{changed}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{changed}}",
			fns(
				p.changed.Name,
				p.changed.Handle,
				p.changed.ID,
				p.changed.NamespaceID,
				p.changed.Config,
			),
		)
		pairs = append(pairs, "{{changed.name}}", fns(p.changed.Name))
		pairs = append(pairs, "{{changed.handle}}", fns(p.changed.Handle))
		pairs = append(pairs, "{{changed.ID}}", fns(p.changed.ID))
		pairs = append(pairs, "{{changed.namespaceID}}", fns(p.changed.NamespaceID))
		pairs = append(pairs, "{{changed.config}}", fns(p.changed.Config))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Query,
				p.filter.Handle,
				p.filter.NamespaceID,
				p.filter.Sort,
				p.filter.Limit,
			),
		)
		pairs = append(pairs, "{{filter.query}}", fns(p.filter.Query))
		pairs = append(pairs, "{{filter.handle}}", fns(p.filter.Handle))
		pairs = append(pairs, "{{filter.namespaceID}}", fns(p.filter.NamespaceID))
		pairs = append(pairs, "{{filter.sort}}", fns(p.filter.Sort))
		pairs = append(pairs, "{{filter.limit}}", fns(p.filter.Limit))
	}

	if p.namespace != nil {
		// replacement for "{{namespace}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{namespace}}",
			fns(
				p.namespace.Name,
				p.namespace.Slug,
				p.namespace.ID,
			),
		)
		pairs = append(pairs, "{{namespace.name}}", fns(p.namespace.Name))
		pairs = append(pairs, "{{namespace.slug}}", fns(p.namespace.Slug))
		pairs = append(pairs, "{{namespace.ID}}", fns(p.namespace.ID))
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
func (a *chartAction) String() string {
	var props = &chartActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *chartAction) ToAction() *actionlog.Action {
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

// ChartActionSearch returns "compose:chart.search" action
//
// This function is auto-generated.
//
func ChartActionSearch(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "search",
		log:       "searched for charts",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionLookup returns "compose:chart.lookup" action
//
// This function is auto-generated.
//
func ChartActionLookup(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "lookup",
		log:       "looked-up for a {{chart}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionCreate returns "compose:chart.create" action
//
// This function is auto-generated.
//
func ChartActionCreate(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "create",
		log:       "created {{chart}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionUpdate returns "compose:chart.update" action
//
// This function is auto-generated.
//
func ChartActionUpdate(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "update",
		log:       "updated {{chart}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionDelete returns "compose:chart.delete" action
//
// This function is auto-generated.
//
func ChartActionDelete(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "delete",
		log:       "deleted {{chart}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionUndelete returns "compose:chart.undelete" action
//
// This function is auto-generated.
//
func ChartActionUndelete(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "undelete",
		log:       "undeleted {{chart}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionReorder returns "compose:chart.reorder" action
//
// This function is auto-generated.
//
func ChartActionReorder(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "reorder",
		log:       "reordered {{chart}}",
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

// ChartErrGeneric returns "compose:chart.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrGeneric(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "compose:chart"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chartLogMetaKey{}, "{err}"),
		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrNotFound returns "compose:chart.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrNotFound(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("chart does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "compose:chart"),

		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrNamespaceNotFound returns "compose:chart.namespaceNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrNamespaceNotFound(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("namespace does not exist", nil),

		errors.Meta("type", "namespaceNotFound"),
		errors.Meta("resource", "compose:chart"),

		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.namespaceNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrModuleNotFound returns "compose:chart.moduleNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrModuleNotFound(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module does not exist", nil),

		errors.Meta("type", "moduleNotFound"),
		errors.Meta("resource", "compose:chart"),

		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.moduleNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrInvalidID returns "compose:chart.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrInvalidID(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "compose:chart"),

		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrInvalidHandle returns "compose:chart.invalidHandle" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrInvalidHandle(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "compose:chart"),

		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrHandleNotUnique returns "compose:chart.handleNotUnique" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrHandleNotUnique(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "compose:chart"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chartLogMetaKey{}, "used duplicate handle ({{chart.handle}}) for chart"),
		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrStaleData returns "compose:chart.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrStaleData(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "compose:chart"),

		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrInvalidNamespaceID returns "compose:chart.invalidNamespaceID" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrInvalidNamespaceID(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid or missing namespace ID", nil),

		errors.Meta("type", "invalidNamespaceID"),
		errors.Meta("resource", "compose:chart"),

		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.invalidNamespaceID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrNotAllowedToRead returns "compose:chart.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToRead(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this chart", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "compose:chart"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chartLogMetaKey{}, "could not read {{chart}}; insufficient permissions"),
		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrNotAllowedToSearch returns "compose:chart.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToSearch(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list charts", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "compose:chart"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chartLogMetaKey{}, "could not search or list charts; insufficient permissions"),
		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrNotAllowedToReadNamespace returns "compose:chart.notAllowedToReadNamespace" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToReadNamespace(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this namespace", nil),

		errors.Meta("type", "notAllowedToReadNamespace"),
		errors.Meta("resource", "compose:chart"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chartLogMetaKey{}, "could not read namespace {{namespace}}; insufficient permissions"),
		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.notAllowedToReadNamespace"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrNotAllowedToCreate returns "compose:chart.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToCreate(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create charts", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "compose:chart"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chartLogMetaKey{}, "could not create charts; insufficient permissions"),
		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrNotAllowedToUpdate returns "compose:chart.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToUpdate(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this chart", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "compose:chart"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chartLogMetaKey{}, "could not update {{chart}}; insufficient permissions"),
		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrNotAllowedToDelete returns "compose:chart.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToDelete(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this chart", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "compose:chart"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chartLogMetaKey{}, "could not delete {{chart}}; insufficient permissions"),
		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChartErrNotAllowedToUndelete returns "compose:chart.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToUndelete(mm ...*chartActionProps) *errors.Error {
	var p = &chartActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this chart", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "compose:chart"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chartLogMetaKey{}, "could not undelete {{chart}}; insufficient permissions"),
		errors.Meta(chartPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "chart.errors.notAllowedToUndelete"),

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
func (svc chart) recordAction(ctx context.Context, props *chartActionProps, actionFn func(...*chartActionProps) *chartAction, err error) error {
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
		a.Description = props.Format(m.AsString(chartLogMetaKey{}), err)

		if p, has := m[chartPropsMetaKey{}]; has {
			a.Meta = p.(*chartActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
