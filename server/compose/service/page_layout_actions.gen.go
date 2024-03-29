package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/service/page_layout_actions.yaml

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
	pageLayoutActionProps struct {
		pageLayout *types.PageLayout
		changed    *types.PageLayout
		filter     *types.PageLayoutFilter
		namespace  *types.Namespace
	}

	pageLayoutAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *pageLayoutActionProps
	}

	pageLayoutLogMetaKey   struct{}
	pageLayoutPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setPageLayout updates pageLayoutActionProps's pageLayout
//
// This function is auto-generated.
func (p *pageLayoutActionProps) setPageLayout(pageLayout *types.PageLayout) *pageLayoutActionProps {
	p.pageLayout = pageLayout
	return p
}

// setChanged updates pageLayoutActionProps's changed
//
// This function is auto-generated.
func (p *pageLayoutActionProps) setChanged(changed *types.PageLayout) *pageLayoutActionProps {
	p.changed = changed
	return p
}

// setFilter updates pageLayoutActionProps's filter
//
// This function is auto-generated.
func (p *pageLayoutActionProps) setFilter(filter *types.PageLayoutFilter) *pageLayoutActionProps {
	p.filter = filter
	return p
}

// setNamespace updates pageLayoutActionProps's namespace
//
// This function is auto-generated.
func (p *pageLayoutActionProps) setNamespace(namespace *types.Namespace) *pageLayoutActionProps {
	p.namespace = namespace
	return p
}

// Serialize converts pageLayoutActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p pageLayoutActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.pageLayout != nil {
		m.Set("pageLayout.handle", p.pageLayout.Handle, true)
		m.Set("pageLayout.ID", p.pageLayout.ID, true)
		m.Set("pageLayout.namespaceID", p.pageLayout.NamespaceID, true)
	}
	if p.changed != nil {
		m.Set("changed.handle", p.changed.Handle, true)
		m.Set("changed.ID", p.changed.ID, true)
		m.Set("changed.namespaceID", p.changed.NamespaceID, true)
		m.Set("changed.blocks", p.changed.Blocks, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.handle", p.filter.Handle, true)
		m.Set("filter.namespaceID", p.filter.NamespaceID, true)
		m.Set("filter.sort", p.filter.Sort, true)
		m.Set("filter.limit", p.filter.Limit, true)
	}
	if p.namespace != nil {
		m.Set("namespace.slug", p.namespace.Slug, true)
		m.Set("namespace.ID", p.namespace.ID, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p pageLayoutActionProps) Format(in string, err error) string {
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

	if p.pageLayout != nil {
		// replacement for "{{pageLayout}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{pageLayout}}",
			fns(
				p.pageLayout.Handle,
				p.pageLayout.ID,
				p.pageLayout.NamespaceID,
			),
		)
		pairs = append(pairs, "{{pageLayout.handle}}", fns(p.pageLayout.Handle))
		pairs = append(pairs, "{{pageLayout.ID}}", fns(p.pageLayout.ID))
		pairs = append(pairs, "{{pageLayout.namespaceID}}", fns(p.pageLayout.NamespaceID))
	}

	if p.changed != nil {
		// replacement for "{{changed}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{changed}}",
			fns(
				p.changed.Handle,
				p.changed.ID,
				p.changed.NamespaceID,
				p.changed.Blocks,
			),
		)
		pairs = append(pairs, "{{changed.handle}}", fns(p.changed.Handle))
		pairs = append(pairs, "{{changed.ID}}", fns(p.changed.ID))
		pairs = append(pairs, "{{changed.namespaceID}}", fns(p.changed.NamespaceID))
		pairs = append(pairs, "{{changed.blocks}}", fns(p.changed.Blocks))
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
				p.namespace.Slug,
				p.namespace.ID,
			),
		)
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
func (a *pageLayoutAction) String() string {
	var props = &pageLayoutActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *pageLayoutAction) ToAction() *actionlog.Action {
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

// PageLayoutActionSearch returns "compose:page-layout.search" action
//
// This function is auto-generated.
func PageLayoutActionSearch(props ...*pageLayoutActionProps) *pageLayoutAction {
	a := &pageLayoutAction{
		timestamp: time.Now(),
		resource:  "compose:page-layout",
		action:    "search",
		log:       "searched for pageLayouts",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageLayoutActionLookup returns "compose:page-layout.lookup" action
//
// This function is auto-generated.
func PageLayoutActionLookup(props ...*pageLayoutActionProps) *pageLayoutAction {
	a := &pageLayoutAction{
		timestamp: time.Now(),
		resource:  "compose:page-layout",
		action:    "lookup",
		log:       "looked-up for a {{pageLayout}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageLayoutActionCreate returns "compose:page-layout.create" action
//
// This function is auto-generated.
func PageLayoutActionCreate(props ...*pageLayoutActionProps) *pageLayoutAction {
	a := &pageLayoutAction{
		timestamp: time.Now(),
		resource:  "compose:page-layout",
		action:    "create",
		log:       "created {{pageLayout}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageLayoutActionUpdate returns "compose:page-layout.update" action
//
// This function is auto-generated.
func PageLayoutActionUpdate(props ...*pageLayoutActionProps) *pageLayoutAction {
	a := &pageLayoutAction{
		timestamp: time.Now(),
		resource:  "compose:page-layout",
		action:    "update",
		log:       "updated {{pageLayout}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageLayoutActionReorder returns "compose:page-layout.reorder" action
//
// This function is auto-generated.
func PageLayoutActionReorder(props ...*pageLayoutActionProps) *pageLayoutAction {
	a := &pageLayoutAction{
		timestamp: time.Now(),
		resource:  "compose:page-layout",
		action:    "reorder",
		log:       "reordered {{pageLayout}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageLayoutActionDelete returns "compose:page-layout.delete" action
//
// This function is auto-generated.
func PageLayoutActionDelete(props ...*pageLayoutActionProps) *pageLayoutAction {
	a := &pageLayoutAction{
		timestamp: time.Now(),
		resource:  "compose:page-layout",
		action:    "delete",
		log:       "deleted {{pageLayout}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageLayoutActionUndelete returns "compose:page-layout.undelete" action
//
// This function is auto-generated.
func PageLayoutActionUndelete(props ...*pageLayoutActionProps) *pageLayoutAction {
	a := &pageLayoutAction{
		timestamp: time.Now(),
		resource:  "compose:page-layout",
		action:    "undelete",
		log:       "undeleted {{pageLayout}}",
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

// PageLayoutErrGeneric returns "compose:page-layout.generic" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrGeneric(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "compose:page-layout"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLayoutLogMetaKey{}, "{err}"),
		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrNotFound returns "compose:page-layout.notFound" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrNotFound(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("pageLayout does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "compose:page-layout"),

		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrNamespaceNotFound returns "compose:page-layout.namespaceNotFound" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrNamespaceNotFound(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("namespace does not exist", nil),

		errors.Meta("type", "namespaceNotFound"),
		errors.Meta("resource", "compose:page-layout"),

		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.namespaceNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrModuleNotFound returns "compose:page-layout.moduleNotFound" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrModuleNotFound(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module does not exist", nil),

		errors.Meta("type", "moduleNotFound"),
		errors.Meta("resource", "compose:page-layout"),

		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.moduleNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrInvalidID returns "compose:page-layout.invalidID" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrInvalidID(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "compose:page-layout"),

		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrInvalidHandle returns "compose:page-layout.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrInvalidHandle(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "compose:page-layout"),

		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrHandleNotUnique returns "compose:page-layout.handleNotUnique" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrHandleNotUnique(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "compose:page-layout"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLayoutLogMetaKey{}, "used duplicate handle ({{pageLayout.handle}}) for pageLayout"),
		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrStaleData returns "compose:page-layout.staleData" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrStaleData(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "compose:page-layout"),

		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrInvalidNamespaceID returns "compose:page-layout.invalidNamespaceID" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrInvalidNamespaceID(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid or missing namespace ID", nil),

		errors.Meta("type", "invalidNamespaceID"),
		errors.Meta("resource", "compose:page-layout"),

		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.invalidNamespaceID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrNotAllowedToRead returns "compose:page-layout.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrNotAllowedToRead(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this pageLayout", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "compose:page-layout"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLayoutLogMetaKey{}, "could not read {{pageLayout}}; insufficient permissions"),
		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrNotAllowedToSearch returns "compose:page-layout.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrNotAllowedToSearch(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list pageLayouts", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "compose:page-layout"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLayoutLogMetaKey{}, "could not search pageLayouts; insufficient permissions"),
		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrNotAllowedToListPageLayouts returns "compose:page-layout.notAllowedToListPageLayouts" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrNotAllowedToListPageLayouts(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list pageLayouts", nil),

		errors.Meta("type", "notAllowedToListPageLayouts"),
		errors.Meta("resource", "compose:page-layout"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLayoutLogMetaKey{}, "could not list pageLayouts; insufficient permissions"),
		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.notAllowedToListPageLayouts"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrNotAllowedToCreate returns "compose:page-layout.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrNotAllowedToCreate(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create pageLayouts", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "compose:page-layout"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLayoutLogMetaKey{}, "could not create pageLayouts; insufficient permissions"),
		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrNotAllowedToUpdate returns "compose:page-layout.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrNotAllowedToUpdate(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this pageLayout", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "compose:page-layout"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLayoutLogMetaKey{}, "could not update {{pageLayout}}; insufficient permissions"),
		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrNotAllowedToDelete returns "compose:page-layout.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrNotAllowedToDelete(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this pageLayout", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "compose:page-layout"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLayoutLogMetaKey{}, "could not delete {{pageLayout}}; insufficient permissions"),
		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageLayoutErrNotAllowedToUndelete returns "compose:page-layout.notAllowedToUndelete" as *errors.Error
//
// This function is auto-generated.
func PageLayoutErrNotAllowedToUndelete(mm ...*pageLayoutActionProps) *errors.Error {
	var p = &pageLayoutActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this pageLayout", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "compose:page-layout"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLayoutLogMetaKey{}, "could not undelete {{pageLayout}}; insufficient permissions"),
		errors.Meta(pageLayoutPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page-layout.errors.notAllowedToUndelete"),

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
func (svc pageLayout) recordAction(ctx context.Context, props *pageLayoutActionProps, actionFn func(...*pageLayoutActionProps) *pageLayoutAction, err error) error {
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
		a.Description = props.Format(m.AsString(pageLayoutLogMetaKey{}), err)

		if p, has := m[pageLayoutPropsMetaKey{}]; has {
			a.Meta = p.(*pageLayoutActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
