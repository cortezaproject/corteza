package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/service/module_actions.yaml

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
	moduleActionProps struct {
		module    *types.Module
		changed   *types.Module
		filter    *types.ModuleFilter
		namespace *types.Namespace
	}

	moduleAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *moduleActionProps
	}

	moduleLogMetaKey   struct{}
	modulePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setModule updates moduleActionProps's module
//
// This function is auto-generated.
//
func (p *moduleActionProps) setModule(module *types.Module) *moduleActionProps {
	p.module = module
	return p
}

// setChanged updates moduleActionProps's changed
//
// This function is auto-generated.
//
func (p *moduleActionProps) setChanged(changed *types.Module) *moduleActionProps {
	p.changed = changed
	return p
}

// setFilter updates moduleActionProps's filter
//
// This function is auto-generated.
//
func (p *moduleActionProps) setFilter(filter *types.ModuleFilter) *moduleActionProps {
	p.filter = filter
	return p
}

// setNamespace updates moduleActionProps's namespace
//
// This function is auto-generated.
//
func (p *moduleActionProps) setNamespace(namespace *types.Namespace) *moduleActionProps {
	p.namespace = namespace
	return p
}

// Serialize converts moduleActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p moduleActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.module != nil {
		m.Set("module.name", p.module.Name, true)
		m.Set("module.handle", p.module.Handle, true)
		m.Set("module.ID", p.module.ID, true)
		m.Set("module.namespaceID", p.module.NamespaceID, true)
	}
	if p.changed != nil {
		m.Set("changed.name", p.changed.Name, true)
		m.Set("changed.handle", p.changed.Handle, true)
		m.Set("changed.ID", p.changed.ID, true)
		m.Set("changed.namespaceID", p.changed.NamespaceID, true)
		m.Set("changed.meta", p.changed.Meta, true)
		m.Set("changed.fields", p.changed.Fields, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.name", p.filter.Name, true)
		m.Set("filter.handle", p.filter.Handle, true)
		m.Set("filter.name", p.filter.Name, true)
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
func (p moduleActionProps) Format(in string, err error) string {
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

	if p.module != nil {
		// replacement for "{{module}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{module}}",
			fns(
				p.module.Name,
				p.module.Handle,
				p.module.ID,
				p.module.NamespaceID,
			),
		)
		pairs = append(pairs, "{{module.name}}", fns(p.module.Name))
		pairs = append(pairs, "{{module.handle}}", fns(p.module.Handle))
		pairs = append(pairs, "{{module.ID}}", fns(p.module.ID))
		pairs = append(pairs, "{{module.namespaceID}}", fns(p.module.NamespaceID))
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
				p.changed.Meta,
				p.changed.Fields,
			),
		)
		pairs = append(pairs, "{{changed.name}}", fns(p.changed.Name))
		pairs = append(pairs, "{{changed.handle}}", fns(p.changed.Handle))
		pairs = append(pairs, "{{changed.ID}}", fns(p.changed.ID))
		pairs = append(pairs, "{{changed.namespaceID}}", fns(p.changed.NamespaceID))
		pairs = append(pairs, "{{changed.meta}}", fns(p.changed.Meta))
		pairs = append(pairs, "{{changed.fields}}", fns(p.changed.Fields))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Query,
				p.filter.Name,
				p.filter.Handle,
				p.filter.Name,
				p.filter.NamespaceID,
				p.filter.Sort,
				p.filter.Limit,
			),
		)
		pairs = append(pairs, "{{filter.query}}", fns(p.filter.Query))
		pairs = append(pairs, "{{filter.name}}", fns(p.filter.Name))
		pairs = append(pairs, "{{filter.handle}}", fns(p.filter.Handle))
		pairs = append(pairs, "{{filter.name}}", fns(p.filter.Name))
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
func (a *moduleAction) String() string {
	var props = &moduleActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *moduleAction) ToAction() *actionlog.Action {
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

// ModuleActionSearch returns "compose:module.search" action
//
// This function is auto-generated.
//
func ModuleActionSearch(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "search",
		log:       "searched for modules",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleActionLookup returns "compose:module.lookup" action
//
// This function is auto-generated.
//
func ModuleActionLookup(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "lookup",
		log:       "looked-up for a {{module}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleActionCreate returns "compose:module.create" action
//
// This function is auto-generated.
//
func ModuleActionCreate(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "create",
		log:       "created {{module}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleActionUpdate returns "compose:module.update" action
//
// This function is auto-generated.
//
func ModuleActionUpdate(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "update",
		log:       "updated {{module}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleActionDelete returns "compose:module.delete" action
//
// This function is auto-generated.
//
func ModuleActionDelete(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "delete",
		log:       "deleted {{module}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleActionUndelete returns "compose:module.undelete" action
//
// This function is auto-generated.
//
func ModuleActionUndelete(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "undelete",
		log:       "undeleted {{module}}",
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

// ModuleErrGeneric returns "compose:module.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrGeneric(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "{err}"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNotFound returns "compose:module.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotFound(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "compose:module"),

		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNamespaceNotFound returns "compose:module.namespaceNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNamespaceNotFound(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("namespace does not exist", nil),

		errors.Meta("type", "namespaceNotFound"),
		errors.Meta("resource", "compose:module"),

		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.namespaceNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrInvalidID returns "compose:module.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrInvalidID(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "compose:module"),

		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrInvalidHandle returns "compose:module.invalidHandle" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrInvalidHandle(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "compose:module"),

		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrHandleNotUnique returns "compose:module.handleNotUnique" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrHandleNotUnique(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "used duplicate handle ({{module.handle}}) for module"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNameNotUnique returns "compose:module.nameNotUnique" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNameNotUnique(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("name not unique", nil),

		errors.Meta("type", "nameNotUnique"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "used duplicate name ({{module.name}}) for module"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.nameNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrFieldNameReserved returns "compose:module.fieldNameReserved" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrFieldNameReserved(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("field name is reserved for system fields", nil),

		errors.Meta("type", "fieldNameReserved"),
		errors.Meta("resource", "compose:module"),

		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.fieldNameReserved"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrStaleData returns "compose:module.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrStaleData(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "compose:module"),

		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrInvalidNamespaceID returns "compose:module.invalidNamespaceID" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrInvalidNamespaceID(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid or missing namespace ID", nil),

		errors.Meta("type", "invalidNamespaceID"),
		errors.Meta("resource", "compose:module"),

		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.invalidNamespaceID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNotAllowedToRead returns "compose:module.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToRead(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this module", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "could not read {{module}}; insufficient permissions"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNotAllowedToSearch returns "compose:module.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToSearch(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list modules", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "could not search or list modules; insufficient permissions"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNotAllowedToReadNamespace returns "compose:module.notAllowedToReadNamespace" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToReadNamespace(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this namespace", nil),

		errors.Meta("type", "notAllowedToReadNamespace"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "could not read namespace {{namespace}}; insufficient permissions"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.notAllowedToReadNamespace"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNotAllowedToListModules returns "compose:module.notAllowedToListModules" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToListModules(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list modules", nil),

		errors.Meta("type", "notAllowedToListModules"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "could not list modules; insufficient permissions"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.notAllowedToListModules"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNotAllowedToCreate returns "compose:module.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToCreate(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create modules", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "could not create modules; insufficient permissions"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNotAllowedToUpdate returns "compose:module.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToUpdate(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this module", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "could not update {{module}}; insufficient permissions"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNotAllowedToDelete returns "compose:module.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToDelete(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this module", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "could not delete {{module}}; insufficient permissions"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleErrNotAllowedToUndelete returns "compose:module.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToUndelete(mm ...*moduleActionProps) *errors.Error {
	var p = &moduleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this module", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "compose:module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleLogMetaKey{}, "could not undelete {{module}}; insufficient permissions"),
		errors.Meta(modulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "module.errors.notAllowedToUndelete"),

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
func (svc module) recordAction(ctx context.Context, props *moduleActionProps, actionFn func(...*moduleActionProps) *moduleAction, err error) error {
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
		a.Description = props.Format(m.AsString(moduleLogMetaKey{}), err)

		if p, has := m[modulePropsMetaKey{}]; has {
			a.Meta = p.(*moduleActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
