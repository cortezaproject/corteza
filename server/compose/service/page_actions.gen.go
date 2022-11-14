package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/service/page_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"strings"
	"time"
)

type (
	pageActionProps struct {
		page      *types.Page
		changed   *types.Page
		filter    *types.PageFilter
		namespace *types.Namespace
	}

	pageAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *pageActionProps
	}

	pageLogMetaKey   struct{}
	pagePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setPage updates pageActionProps's page
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *pageActionProps) setPage(page *types.Page) *pageActionProps {
	p.page = page
	return p
}

// setChanged updates pageActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *pageActionProps) setChanged(changed *types.Page) *pageActionProps {
	p.changed = changed
	return p
}

// setFilter updates pageActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *pageActionProps) setFilter(filter *types.PageFilter) *pageActionProps {
	p.filter = filter
	return p
}

// setNamespace updates pageActionProps's namespace
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *pageActionProps) setNamespace(namespace *types.Namespace) *pageActionProps {
	p.namespace = namespace
	return p
}

// Serialize converts pageActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p pageActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.page != nil {
		m.Set("page.title", p.page.Title, true)
		m.Set("page.handle", p.page.Handle, true)
		m.Set("page.ID", p.page.ID, true)
		m.Set("page.namespaceID", p.page.NamespaceID, true)
		m.Set("page.moduleID", p.page.ModuleID, true)
	}
	if p.changed != nil {
		m.Set("changed.title", p.changed.Title, true)
		m.Set("changed.handle", p.changed.Handle, true)
		m.Set("changed.ID", p.changed.ID, true)
		m.Set("changed.namespaceID", p.changed.NamespaceID, true)
		m.Set("changed.description", p.changed.Description, true)
		m.Set("changed.moduleID", p.changed.ModuleID, true)
		m.Set("changed.blocks", p.changed.Blocks, true)
		m.Set("changed.visible", p.changed.Visible, true)
		m.Set("changed.weight", p.changed.Weight, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.handle", p.filter.Handle, true)
		m.Set("filter.root", p.filter.Root, true)
		m.Set("filter.namespaceID", p.filter.NamespaceID, true)
		m.Set("filter.parentID", p.filter.ParentID, true)
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
func (p pageActionProps) Format(in string, err error) string {
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

	if p.page != nil {
		// replacement for "{{page}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{page}}",
			fns(
				p.page.Title,
				p.page.Handle,
				p.page.ID,
				p.page.NamespaceID,
				p.page.ModuleID,
			),
		)
		pairs = append(pairs, "{{page.title}}", fns(p.page.Title))
		pairs = append(pairs, "{{page.handle}}", fns(p.page.Handle))
		pairs = append(pairs, "{{page.ID}}", fns(p.page.ID))
		pairs = append(pairs, "{{page.namespaceID}}", fns(p.page.NamespaceID))
		pairs = append(pairs, "{{page.moduleID}}", fns(p.page.ModuleID))
	}

	if p.changed != nil {
		// replacement for "{{changed}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{changed}}",
			fns(
				p.changed.Title,
				p.changed.Handle,
				p.changed.ID,
				p.changed.NamespaceID,
				p.changed.Description,
				p.changed.ModuleID,
				p.changed.Blocks,
				p.changed.Visible,
				p.changed.Weight,
			),
		)
		pairs = append(pairs, "{{changed.title}}", fns(p.changed.Title))
		pairs = append(pairs, "{{changed.handle}}", fns(p.changed.Handle))
		pairs = append(pairs, "{{changed.ID}}", fns(p.changed.ID))
		pairs = append(pairs, "{{changed.namespaceID}}", fns(p.changed.NamespaceID))
		pairs = append(pairs, "{{changed.description}}", fns(p.changed.Description))
		pairs = append(pairs, "{{changed.moduleID}}", fns(p.changed.ModuleID))
		pairs = append(pairs, "{{changed.blocks}}", fns(p.changed.Blocks))
		pairs = append(pairs, "{{changed.visible}}", fns(p.changed.Visible))
		pairs = append(pairs, "{{changed.weight}}", fns(p.changed.Weight))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Query,
				p.filter.Handle,
				p.filter.Root,
				p.filter.NamespaceID,
				p.filter.ParentID,
				p.filter.Sort,
				p.filter.Limit,
			),
		)
		pairs = append(pairs, "{{filter.query}}", fns(p.filter.Query))
		pairs = append(pairs, "{{filter.handle}}", fns(p.filter.Handle))
		pairs = append(pairs, "{{filter.root}}", fns(p.filter.Root))
		pairs = append(pairs, "{{filter.namespaceID}}", fns(p.filter.NamespaceID))
		pairs = append(pairs, "{{filter.parentID}}", fns(p.filter.ParentID))
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
func (a *pageAction) String() string {
	var props = &pageActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *pageAction) ToAction() *actionlog.Action {
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

// PageActionSearch returns "compose:page.search" action
//
// This function is auto-generated.
//
func PageActionSearch(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "search",
		log:       "searched for pages",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionLookup returns "compose:page.lookup" action
//
// This function is auto-generated.
//
func PageActionLookup(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "lookup",
		log:       "looked-up for a {{page}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionCreate returns "compose:page.create" action
//
// This function is auto-generated.
//
func PageActionCreate(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "create",
		log:       "created {{page}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionUpdate returns "compose:page.update" action
//
// This function is auto-generated.
//
func PageActionUpdate(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "update",
		log:       "updated {{page}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionDelete returns "compose:page.delete" action
//
// This function is auto-generated.
//
func PageActionDelete(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "delete",
		log:       "deleted {{page}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionUndelete returns "compose:page.undelete" action
//
// This function is auto-generated.
//
func PageActionUndelete(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "undelete",
		log:       "undeleted {{page}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionReorder returns "compose:page.reorder" action
//
// This function is auto-generated.
//
func PageActionReorder(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "reorder",
		log:       "reordered {{page}}",
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

// PageErrGeneric returns "compose:page.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrGeneric(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "{err}"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNotFound returns "compose:page.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNotFound(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("page does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "compose:page"),

		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNamespaceNotFound returns "compose:page.namespaceNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNamespaceNotFound(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("namespace does not exist", nil),

		errors.Meta("type", "namespaceNotFound"),
		errors.Meta("resource", "compose:page"),

		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.namespaceNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrModuleNotFound returns "compose:page.moduleNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrModuleNotFound(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module does not exist", nil),

		errors.Meta("type", "moduleNotFound"),
		errors.Meta("resource", "compose:page"),

		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.moduleNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrInvalidID returns "compose:page.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrInvalidID(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "compose:page"),

		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrInvalidHandle returns "compose:page.invalidHandle" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrInvalidHandle(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "compose:page"),

		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrHandleNotUnique returns "compose:page.handleNotUnique" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrHandleNotUnique(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "used duplicate handle ({{page.handle}}) for page"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrStaleData returns "compose:page.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrStaleData(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "compose:page"),

		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrInvalidNamespaceID returns "compose:page.invalidNamespaceID" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrInvalidNamespaceID(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid or missing namespace ID", nil),

		errors.Meta("type", "invalidNamespaceID"),
		errors.Meta("resource", "compose:page"),

		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.invalidNamespaceID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrDeleteAbortedForPageWithSubpages returns "compose:page.deleteAbortedForPageWithSubpages" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrDeleteAbortedForPageWithSubpages(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("removal of page with subpages aborted", nil),

		errors.Meta("type", "deleteAbortedForPageWithSubpages"),
		errors.Meta("resource", "compose:page"),

		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.deleteAbortedForPageWithSubpages"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrUnknownDeleteStrategy returns "compose:page.unknownDeleteStrategy" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrUnknownDeleteStrategy(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("unknown delete strategy", nil),

		errors.Meta("type", "unknownDeleteStrategy"),
		errors.Meta("resource", "compose:page"),

		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.unknownDeleteStrategy"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNotAllowedToRead returns "compose:page.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToRead(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this page", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "could not read {{page}}; insufficient permissions"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNotAllowedToSearch returns "compose:page.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToSearch(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list pages", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "could not search pages; insufficient permissions"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNotAllowedToReadNamespace returns "compose:page.notAllowedToReadNamespace" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToReadNamespace(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this namespace", nil),

		errors.Meta("type", "notAllowedToReadNamespace"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "could not read namespace {{namespace}}; insufficient permissions"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.notAllowedToReadNamespace"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNotAllowedToListPages returns "compose:page.notAllowedToListPages" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToListPages(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list pages", nil),

		errors.Meta("type", "notAllowedToListPages"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "could not list pages; insufficient permissions"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.notAllowedToListPages"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNotAllowedToCreate returns "compose:page.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToCreate(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create pages", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "could not create pages; insufficient permissions"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNotAllowedToUpdate returns "compose:page.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToUpdate(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this page", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "could not update {{page}}; insufficient permissions"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNotAllowedToDelete returns "compose:page.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToDelete(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this page", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "could not delete {{page}}; insufficient permissions"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// PageErrNotAllowedToUndelete returns "compose:page.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToUndelete(mm ...*pageActionProps) *errors.Error {
	var p = &pageActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this page", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "compose:page"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(pageLogMetaKey{}, "could not undelete {{page}}; insufficient permissions"),
		errors.Meta(pagePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "page.errors.notAllowedToUndelete"),

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
func (svc page) recordAction(ctx context.Context, props *pageActionProps, actionFn func(...*pageActionProps) *pageAction, err error) error {
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
		a.Description = props.Format(m.AsString(pageLogMetaKey{}), err)

		if p, has := m[pagePropsMetaKey{}]; has {
			a.Meta = p.(*pageActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
