package service

// This file is auto-generated from compose/service/page_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
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

	pageError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *pageActionProps
	}
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

// serialize converts pageActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p pageActionProps) serialize() actionlog.Meta {
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
		m.Set("filter.offset", p.filter.Offset, true)
		m.Set("filter.page", p.filter.Page, true)
		m.Set("filter.perPage", p.filter.PerPage, true)
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
func (p pageActionProps) tr(in string, err error) string {
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
		for {
			// Unwrap errors
			ue := errors.Unwrap(err)
			if ue == nil {
				break
			}

			err = ue
		}

		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}

	if p.page != nil {
		// replacement for "{page}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{page}",
			fns(
				p.page.Title,
				p.page.Handle,
				p.page.ID,
				p.page.NamespaceID,
				p.page.ModuleID,
			),
		)
		pairs = append(pairs, "{page.title}", fns(p.page.Title))
		pairs = append(pairs, "{page.handle}", fns(p.page.Handle))
		pairs = append(pairs, "{page.ID}", fns(p.page.ID))
		pairs = append(pairs, "{page.namespaceID}", fns(p.page.NamespaceID))
		pairs = append(pairs, "{page.moduleID}", fns(p.page.ModuleID))
	}

	if p.changed != nil {
		// replacement for "{changed}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{changed}",
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
		pairs = append(pairs, "{changed.title}", fns(p.changed.Title))
		pairs = append(pairs, "{changed.handle}", fns(p.changed.Handle))
		pairs = append(pairs, "{changed.ID}", fns(p.changed.ID))
		pairs = append(pairs, "{changed.namespaceID}", fns(p.changed.NamespaceID))
		pairs = append(pairs, "{changed.description}", fns(p.changed.Description))
		pairs = append(pairs, "{changed.moduleID}", fns(p.changed.ModuleID))
		pairs = append(pairs, "{changed.blocks}", fns(p.changed.Blocks))
		pairs = append(pairs, "{changed.visible}", fns(p.changed.Visible))
		pairs = append(pairs, "{changed.weight}", fns(p.changed.Weight))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.Handle,
				p.filter.Root,
				p.filter.NamespaceID,
				p.filter.ParentID,
				p.filter.Sort,
				p.filter.Limit,
				p.filter.Offset,
				p.filter.Page,
				p.filter.PerPage,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.handle}", fns(p.filter.Handle))
		pairs = append(pairs, "{filter.root}", fns(p.filter.Root))
		pairs = append(pairs, "{filter.namespaceID}", fns(p.filter.NamespaceID))
		pairs = append(pairs, "{filter.parentID}", fns(p.filter.ParentID))
		pairs = append(pairs, "{filter.sort}", fns(p.filter.Sort))
		pairs = append(pairs, "{filter.limit}", fns(p.filter.Limit))
		pairs = append(pairs, "{filter.offset}", fns(p.filter.Offset))
		pairs = append(pairs, "{filter.page}", fns(p.filter.Page))
		pairs = append(pairs, "{filter.perPage}", fns(p.filter.PerPage))
	}

	if p.namespace != nil {
		// replacement for "{namespace}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{namespace}",
			fns(
				p.namespace.Name,
				p.namespace.Slug,
				p.namespace.ID,
			),
		)
		pairs = append(pairs, "{namespace.name}", fns(p.namespace.Name))
		pairs = append(pairs, "{namespace.slug}", fns(p.namespace.Slug))
		pairs = append(pairs, "{namespace.ID}", fns(p.namespace.ID))
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

	return props.tr(a.log, nil)
}

func (e *pageAction) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error methods

// String returns loggable description as string
//
// It falls back to message if log is not set
//
// This function is auto-generated.
//
func (e *pageError) String() string {
	var props = &pageActionProps{}

	if e.props != nil {
		props = e.props
	}

	if e.wrap != nil && !strings.Contains(e.log, "{err}") {
		// Suffix error log with {err} to ensure
		// we log the cause for this error
		e.log += ": {err}"
	}

	return props.tr(e.log, e.wrap)
}

// Error satisfies
//
// This function is auto-generated.
//
func (e *pageError) Error() string {
	var props = &pageActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *pageError) Is(Resource error) bool {
	t, ok := Resource.(*pageError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps pageError around another error
//
// This function is auto-generated.
//
func (e *pageError) Wrap(err error) *pageError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *pageError) Unwrap() error {
	return e.wrap
}

func (e *pageError) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Error:       e.Error(),
		Meta:        e.props.serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// PageActionSearch returns "compose:page.search" error
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

// PageActionLookup returns "compose:page.lookup" error
//
// This function is auto-generated.
//
func PageActionLookup(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "lookup",
		log:       "looked-up for a {page}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionCreate returns "compose:page.create" error
//
// This function is auto-generated.
//
func PageActionCreate(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "create",
		log:       "created {page}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionUpdate returns "compose:page.update" error
//
// This function is auto-generated.
//
func PageActionUpdate(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "update",
		log:       "updated {page}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionDelete returns "compose:page.delete" error
//
// This function is auto-generated.
//
func PageActionDelete(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "delete",
		log:       "deleted {page}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionUndelete returns "compose:page.undelete" error
//
// This function is auto-generated.
//
func PageActionUndelete(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "undelete",
		log:       "undeleted {page}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// PageActionReorder returns "compose:page.reorder" error
//
// This function is auto-generated.
//
func PageActionReorder(props ...*pageActionProps) *pageAction {
	a := &pageAction{
		timestamp: time.Now(),
		resource:  "compose:page",
		action:    "reorder",
		log:       "reordered {page}",
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

// PageErrGeneric returns "compose:page.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func PageErrGeneric(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrNotFound returns "compose:page.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func PageErrNotFound(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "notFound",
		action:    "error",
		message:   "page does not exist",
		log:       "page does not exist",
		severity:  actionlog.Warning,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrNamespaceNotFound returns "compose:page.namespaceNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func PageErrNamespaceNotFound(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "namespaceNotFound",
		action:    "error",
		message:   "namespace does not exist",
		log:       "namespace does not exist",
		severity:  actionlog.Warning,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrModuleNotFound returns "compose:page.moduleNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func PageErrModuleNotFound(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "moduleNotFound",
		action:    "error",
		message:   "module does not exist",
		log:       "module does not exist",
		severity:  actionlog.Warning,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrInvalidID returns "compose:page.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func PageErrInvalidID(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrInvalidHandle returns "compose:page.invalidHandle" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func PageErrInvalidHandle(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "invalidHandle",
		action:    "error",
		message:   "invalid handle",
		log:       "invalid handle",
		severity:  actionlog.Warning,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrHandleNotUnique returns "compose:page.handleNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func PageErrHandleNotUnique(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "handleNotUnique",
		action:    "error",
		message:   "handle not unique",
		log:       "used duplicate handle ({page.handle}) for page",
		severity:  actionlog.Warning,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrStaleData returns "compose:page.staleData" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func PageErrStaleData(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "staleData",
		action:    "error",
		message:   "stale data",
		log:       "stale data",
		severity:  actionlog.Warning,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrInvalidNamespaceID returns "compose:page.invalidNamespaceID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func PageErrInvalidNamespaceID(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "invalidNamespaceID",
		action:    "error",
		message:   "invalid or missing namespace ID",
		log:       "invalid or missing namespace ID",
		severity:  actionlog.Warning,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrNotAllowedToRead returns "compose:page.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToRead(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this page",
		log:       "could not read {page}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrNotAllowedToReadNamespace returns "compose:page.notAllowedToReadNamespace" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToReadNamespace(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "notAllowedToReadNamespace",
		action:    "error",
		message:   "not allowed to read this namespace",
		log:       "could not read namespace {namespace}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrNotAllowedToListPages returns "compose:page.notAllowedToListPages" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToListPages(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "notAllowedToListPages",
		action:    "error",
		message:   "not allowed to list pages",
		log:       "could not list pages; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrNotAllowedToCreate returns "compose:page.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToCreate(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create pages",
		log:       "could not create pages; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrNotAllowedToUpdate returns "compose:page.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToUpdate(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this page",
		log:       "could not update {page}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrNotAllowedToDelete returns "compose:page.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToDelete(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this page",
		log:       "could not delete {page}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// PageErrNotAllowedToUndelete returns "compose:page.notAllowedToUndelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func PageErrNotAllowedToUndelete(props ...*pageActionProps) *pageError {
	var e = &pageError{
		timestamp: time.Now(),
		resource:  "compose:page",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this page",
		log:       "could not undelete {page}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *pageActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// context is used to enrich audit log entry with current user info, request ID, IP address...
// props are collected action/error properties
// action (optional) fn will be used to construct pageAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc page) recordAction(ctx context.Context, props *pageActionProps, action func(...*pageActionProps) *pageAction, err error) error {
	var (
		ok bool

		// Return error
		retError *pageError

		// Recorder error
		recError *pageError
	)

	if err != nil {
		if retError, ok = err.(*pageError); !ok {
			// got non-page error, wrap it with PageErrGeneric
			retError = PageErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use PageErrGeneric for recording too
			// because it can hold more info
			recError = retError
		} else if retError != nil {
			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}
			// start with copy of return error for recording
			// this will be updated with tha root cause as we try and
			// unwrap the error
			recError = retError

			// find the original recError for this error
			// for the purpose of logging
			var unwrappedError error = retError
			for {
				if unwrappedError = errors.Unwrap(unwrappedError); unwrappedError == nil {
					// nothing wrapped
					break
				}

				// update recError ONLY of wrapped error is of type pageError
				if unwrappedSinkError, ok := unwrappedError.(*pageError); ok {
					recError = unwrappedSinkError
				}
			}

			if retError.props == nil {
				// set props on returning error if empty
				retError.props = props
			}

			if recError.props == nil {
				// set props on recording error if empty
				recError.props = props
			}
		}
	}

	if svc.actionlog != nil {
		if retError != nil {
			// failed action, log error
			svc.actionlog.Record(ctx, recError)
		} else if action != nil {
			// successful
			svc.actionlog.Record(ctx, action(props))
		}
	}

	if err == nil {
		// retError not an interface and that WILL (!!) cause issues
		// with nil check (== nil) when it is not explicitly returned
		return nil
	}

	return retError
}
