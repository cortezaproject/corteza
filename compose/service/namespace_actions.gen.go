package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/service/namespace_actions.yaml

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
	namespaceActionProps struct {
		namespace     *types.Namespace
		changed       *types.Namespace
		archiveFormat string
		filter        *types.NamespaceFilter
	}

	namespaceAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *namespaceActionProps
	}

	namespaceLogMetaKey   struct{}
	namespacePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setNamespace updates namespaceActionProps's namespace
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *namespaceActionProps) setNamespace(namespace *types.Namespace) *namespaceActionProps {
	p.namespace = namespace
	return p
}

// setChanged updates namespaceActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *namespaceActionProps) setChanged(changed *types.Namespace) *namespaceActionProps {
	p.changed = changed
	return p
}

// setArchiveFormat updates namespaceActionProps's archiveFormat
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *namespaceActionProps) setArchiveFormat(archiveFormat string) *namespaceActionProps {
	p.archiveFormat = archiveFormat
	return p
}

// setFilter updates namespaceActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *namespaceActionProps) setFilter(filter *types.NamespaceFilter) *namespaceActionProps {
	p.filter = filter
	return p
}

// Serialize converts namespaceActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p namespaceActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.namespace != nil {
		m.Set("namespace.name", p.namespace.Name, true)
		m.Set("namespace.slug", p.namespace.Slug, true)
		m.Set("namespace.ID", p.namespace.ID, true)
		m.Set("namespace.enabled", p.namespace.Enabled, true)
	}
	if p.changed != nil {
		m.Set("changed.name", p.changed.Name, true)
		m.Set("changed.slug", p.changed.Slug, true)
		m.Set("changed.ID", p.changed.ID, true)
		m.Set("changed.meta", p.changed.Meta, true)
		m.Set("changed.enabled", p.changed.Enabled, true)
	}
	m.Set("archiveFormat", p.archiveFormat, true)
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.slug", p.filter.Slug, true)
		m.Set("filter.sort", p.filter.Sort, true)
		m.Set("filter.limit", p.filter.Limit, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p namespaceActionProps) Format(in string, err error) string {
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

	if p.namespace != nil {
		// replacement for "{{namespace}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{namespace}}",
			fns(
				p.namespace.Name,
				p.namespace.Slug,
				p.namespace.ID,
				p.namespace.Enabled,
			),
		)
		pairs = append(pairs, "{{namespace.name}}", fns(p.namespace.Name))
		pairs = append(pairs, "{{namespace.slug}}", fns(p.namespace.Slug))
		pairs = append(pairs, "{{namespace.ID}}", fns(p.namespace.ID))
		pairs = append(pairs, "{{namespace.enabled}}", fns(p.namespace.Enabled))
	}

	if p.changed != nil {
		// replacement for "{{changed}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{changed}}",
			fns(
				p.changed.Name,
				p.changed.Slug,
				p.changed.ID,
				p.changed.Meta,
				p.changed.Enabled,
			),
		)
		pairs = append(pairs, "{{changed.name}}", fns(p.changed.Name))
		pairs = append(pairs, "{{changed.slug}}", fns(p.changed.Slug))
		pairs = append(pairs, "{{changed.ID}}", fns(p.changed.ID))
		pairs = append(pairs, "{{changed.meta}}", fns(p.changed.Meta))
		pairs = append(pairs, "{{changed.enabled}}", fns(p.changed.Enabled))
	}
	pairs = append(pairs, "{{archiveFormat}}", fns(p.archiveFormat))

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Query,
				p.filter.Slug,
				p.filter.Sort,
				p.filter.Limit,
			),
		)
		pairs = append(pairs, "{{filter.query}}", fns(p.filter.Query))
		pairs = append(pairs, "{{filter.slug}}", fns(p.filter.Slug))
		pairs = append(pairs, "{{filter.sort}}", fns(p.filter.Sort))
		pairs = append(pairs, "{{filter.limit}}", fns(p.filter.Limit))
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
func (a *namespaceAction) String() string {
	var props = &namespaceActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *namespaceAction) ToAction() *actionlog.Action {
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

// NamespaceActionSearch returns "compose:namespace.search" action
//
// This function is auto-generated.
//
func NamespaceActionSearch(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "search",
		log:       "searched for namespaces",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionLookup returns "compose:namespace.lookup" action
//
// This function is auto-generated.
//
func NamespaceActionLookup(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "lookup",
		log:       "looked-up for a {{namespace}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionCreate returns "compose:namespace.create" action
//
// This function is auto-generated.
//
func NamespaceActionCreate(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "create",
		log:       "created {{namespace}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionUpdate returns "compose:namespace.update" action
//
// This function is auto-generated.
//
func NamespaceActionUpdate(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "update",
		log:       "updated {{namespace}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionClone returns "compose:namespace.clone" action
//
// This function is auto-generated.
//
func NamespaceActionClone(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "clone",
		log:       "cloned {namespace}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionExport returns "compose:namespace.export" action
//
// This function is auto-generated.
//
func NamespaceActionExport(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "export",
		log:       "exported {namespace}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionImportInit returns "compose:namespace.importInit" action
//
// This function is auto-generated.
//
func NamespaceActionImportInit(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "importInit",
		log:       "import initialized for {namespace}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionImportRun returns "compose:namespace.importRun" action
//
// This function is auto-generated.
//
func NamespaceActionImportRun(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "importRun",
		log:       "imported {namespace}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionDelete returns "compose:namespace.delete" action
//
// This function is auto-generated.
//
func NamespaceActionDelete(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "delete",
		log:       "deleted {{namespace}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionUndelete returns "compose:namespace.undelete" action
//
// This function is auto-generated.
//
func NamespaceActionUndelete(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "undelete",
		log:       "undeleted {{namespace}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionReorder returns "compose:namespace.reorder" action
//
// This function is auto-generated.
//
func NamespaceActionReorder(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "reorder",
		log:       "reordered {{namespace}}",
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

// NamespaceErrGeneric returns "compose:namespace.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrGeneric(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "{err}"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrNotFound returns "compose:namespace.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotFound(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("namespace does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "compose:namespace"),

		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrInvalidID returns "compose:namespace.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrInvalidID(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "compose:namespace"),

		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrInvalidHandle returns "compose:namespace.invalidHandle" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrInvalidHandle(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "compose:namespace"),

		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrHandleNotUnique returns "compose:namespace.handleNotUnique" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrHandleNotUnique(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "used duplicate handle ({{namespace.slug}}) for namespace"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrStaleData returns "compose:namespace.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrStaleData(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "compose:namespace"),

		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrUnsupportedExportFormat returns "compose:namespace.unsupportedExportFormat" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrUnsupportedExportFormat(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("unsupported export format", nil),

		errors.Meta("type", "unsupportedExportFormat"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not export namespace {{namespace}}; unsupported format {{archiveFormat}}"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.unsupportedExportFormat"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrUnsupportedImportFormat returns "compose:namespace.unsupportedImportFormat" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrUnsupportedImportFormat(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("unsupported import format", nil),

		errors.Meta("type", "unsupportedImportFormat"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not import namespace {{namespace}}; unsupported format {{archiveFormat}}"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.unsupportedImportFormat"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrImportMissingNamespace returns "compose:namespace.importMissingNamespace" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrImportMissingNamespace(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("import source does not contain a namespace", nil),

		errors.Meta("type", "importMissingNamespace"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not import namespace; import source does not contain a namespace"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.importMissingNamespace"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrImportSessionNotFound returns "compose:namespace.importSessionNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrImportSessionNotFound(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("the import session does not exist", nil),

		errors.Meta("type", "importSessionNotFound"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not import namespace {{namespace}}; the import session does not exist"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.importSessionNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrCloneMultiple returns "compose:namespace.cloneMultiple" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrCloneMultiple(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to clone multiple namespaces at once", nil),

		errors.Meta("type", "cloneMultiple"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not clone namespaces; multiple duplications requested at once"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.cloneMultiple"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrNotAllowedToRead returns "compose:namespace.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToRead(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this namespace", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not read {{namespace}}; insufficient permissions"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrNotAllowedToSearch returns "compose:namespace.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToSearch(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list namespaces", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not search or list namespaces; insufficient permissions"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrNotAllowedToCreate returns "compose:namespace.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToCreate(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create namespaces", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not create namespaces; insufficient permissions"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrNotAllowedToUpdate returns "compose:namespace.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToUpdate(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this namespace", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not update {{namespace}}; insufficient permissions"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrNotAllowedToDelete returns "compose:namespace.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToDelete(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this namespace", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not delete {{namespace}}; insufficient permissions"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NamespaceErrNotAllowedToUndelete returns "compose:namespace.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToUndelete(mm ...*namespaceActionProps) *errors.Error {
	var p = &namespaceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this namespace", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "compose:namespace"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(namespaceLogMetaKey{}, "could not undelete {{namespace}}; insufficient permissions"),
		errors.Meta(namespacePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "namespace.errors.notAllowedToUndelete"),

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
func (svc namespace) recordAction(ctx context.Context, props *namespaceActionProps, actionFn func(...*namespaceActionProps) *namespaceAction, err error) error {
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
		a.Description = props.Format(m.AsString(namespaceLogMetaKey{}), err)

		if p, has := m[namespacePropsMetaKey{}]; has {
			a.Meta = p.(*namespaceActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
