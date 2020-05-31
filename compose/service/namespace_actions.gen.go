package service

// This file is auto-generated from compose/service/namespace_actions.yaml
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
	namespaceActionProps struct {
		namespace *types.Namespace
		changed   *types.Namespace
		filter    *types.NamespaceFilter
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

	namespaceError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *namespaceActionProps
	}
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

// serialize converts namespaceActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p namespaceActionProps) serialize() actionlog.Meta {
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
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.slug", p.filter.Slug, true)
		m.Set("filter.sort", p.filter.Sort, true)
		m.Set("filter.limit", p.filter.Limit, true)
		m.Set("filter.offset", p.filter.Offset, true)
		m.Set("filter.page", p.filter.Page, true)
		m.Set("filter.perPage", p.filter.PerPage, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p namespaceActionProps) tr(in string, err error) string {
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

	if p.namespace != nil {
		// replacement for "{namespace}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{namespace}",
			fns(
				p.namespace.Name,
				p.namespace.Slug,
				p.namespace.ID,
				p.namespace.Enabled,
			),
		)
		pairs = append(pairs, "{namespace.name}", fns(p.namespace.Name))
		pairs = append(pairs, "{namespace.slug}", fns(p.namespace.Slug))
		pairs = append(pairs, "{namespace.ID}", fns(p.namespace.ID))
		pairs = append(pairs, "{namespace.enabled}", fns(p.namespace.Enabled))
	}

	if p.changed != nil {
		// replacement for "{changed}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{changed}",
			fns(
				p.changed.Name,
				p.changed.Slug,
				p.changed.ID,
				p.changed.Meta,
				p.changed.Enabled,
			),
		)
		pairs = append(pairs, "{changed.name}", fns(p.changed.Name))
		pairs = append(pairs, "{changed.slug}", fns(p.changed.Slug))
		pairs = append(pairs, "{changed.ID}", fns(p.changed.ID))
		pairs = append(pairs, "{changed.meta}", fns(p.changed.Meta))
		pairs = append(pairs, "{changed.enabled}", fns(p.changed.Enabled))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.Slug,
				p.filter.Sort,
				p.filter.Limit,
				p.filter.Offset,
				p.filter.Page,
				p.filter.PerPage,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.slug}", fns(p.filter.Slug))
		pairs = append(pairs, "{filter.sort}", fns(p.filter.Sort))
		pairs = append(pairs, "{filter.limit}", fns(p.filter.Limit))
		pairs = append(pairs, "{filter.offset}", fns(p.filter.Offset))
		pairs = append(pairs, "{filter.page}", fns(p.filter.Page))
		pairs = append(pairs, "{filter.perPage}", fns(p.filter.PerPage))
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

	return props.tr(a.log, nil)
}

func (e *namespaceAction) LoggableAction() *actionlog.Action {
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
func (e *namespaceError) String() string {
	var props = &namespaceActionProps{}

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
func (e *namespaceError) Error() string {
	var props = &namespaceActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *namespaceError) Is(Resource error) bool {
	t, ok := Resource.(*namespaceError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps namespaceError around another error
//
// This function is auto-generated.
//
func (e *namespaceError) Wrap(err error) *namespaceError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *namespaceError) Unwrap() error {
	return e.wrap
}

func (e *namespaceError) LoggableAction() *actionlog.Action {
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

// NamespaceActionSearch returns "compose:namespace.search" error
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

// NamespaceActionLookup returns "compose:namespace.lookup" error
//
// This function is auto-generated.
//
func NamespaceActionLookup(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "lookup",
		log:       "looked-up for a {namespace}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionCreate returns "compose:namespace.create" error
//
// This function is auto-generated.
//
func NamespaceActionCreate(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "create",
		log:       "created {namespace}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionUpdate returns "compose:namespace.update" error
//
// This function is auto-generated.
//
func NamespaceActionUpdate(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "update",
		log:       "updated {namespace}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionDelete returns "compose:namespace.delete" error
//
// This function is auto-generated.
//
func NamespaceActionDelete(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "delete",
		log:       "deleted {namespace}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionUndelete returns "compose:namespace.undelete" error
//
// This function is auto-generated.
//
func NamespaceActionUndelete(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "undelete",
		log:       "undeleted {namespace}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NamespaceActionReorder returns "compose:namespace.reorder" error
//
// This function is auto-generated.
//
func NamespaceActionReorder(props ...*namespaceActionProps) *namespaceAction {
	a := &namespaceAction{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		action:    "reorder",
		log:       "reordered {namespace}",
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

// NamespaceErrGeneric returns "compose:namespace.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NamespaceErrGeneric(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *namespaceActionProps {
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

// NamespaceErrNotFound returns "compose:namespace.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func NamespaceErrNotFound(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "notFound",
		action:    "error",
		message:   "namespace does not exist",
		log:       "namespace does not exist",
		severity:  actionlog.Warning,
		props: func() *namespaceActionProps {
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

// NamespaceErrInvalidID returns "compose:namespace.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func NamespaceErrInvalidID(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *namespaceActionProps {
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

// NamespaceErrInvalidHandle returns "compose:namespace.invalidHandle" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func NamespaceErrInvalidHandle(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "invalidHandle",
		action:    "error",
		message:   "invalid handle",
		log:       "invalid handle",
		severity:  actionlog.Warning,
		props: func() *namespaceActionProps {
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

// NamespaceErrHandleNotUnique returns "compose:namespace.handleNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func NamespaceErrHandleNotUnique(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "handleNotUnique",
		action:    "error",
		message:   "handle not unique",
		log:       "used duplicate handle ({namespace.slug}) for namespace",
		severity:  actionlog.Warning,
		props: func() *namespaceActionProps {
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

// NamespaceErrStaleData returns "compose:namespace.staleData" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func NamespaceErrStaleData(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "staleData",
		action:    "error",
		message:   "stale data",
		log:       "stale data",
		severity:  actionlog.Warning,
		props: func() *namespaceActionProps {
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

// NamespaceErrNotAllowedToRead returns "compose:namespace.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToRead(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this namespace",
		log:       "could not read {namespace}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *namespaceActionProps {
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

// NamespaceErrNotAllowedToListNamespaces returns "compose:namespace.notAllowedToListNamespaces" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToListNamespaces(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "notAllowedToListNamespaces",
		action:    "error",
		message:   "not allowed to list this namespaces",
		log:       "could not list namespaces; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *namespaceActionProps {
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

// NamespaceErrNotAllowedToCreate returns "compose:namespace.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToCreate(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create namespaces",
		log:       "could not create namespaces; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *namespaceActionProps {
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

// NamespaceErrNotAllowedToUpdate returns "compose:namespace.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToUpdate(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this namespace",
		log:       "could not update {namespace}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *namespaceActionProps {
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

// NamespaceErrNotAllowedToDelete returns "compose:namespace.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToDelete(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this namespace",
		log:       "could not delete {namespace}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *namespaceActionProps {
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

// NamespaceErrNotAllowedToUndelete returns "compose:namespace.notAllowedToUndelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NamespaceErrNotAllowedToUndelete(props ...*namespaceActionProps) *namespaceError {
	var e = &namespaceError{
		timestamp: time.Now(),
		resource:  "compose:namespace",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this namespace",
		log:       "could not undelete {namespace}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *namespaceActionProps {
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
// action (optional) fn will be used to construct namespaceAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc namespace) recordAction(ctx context.Context, props *namespaceActionProps, action func(...*namespaceActionProps) *namespaceAction, err error) error {
	var (
		ok bool

		// Return error
		retError *namespaceError

		// Recorder error
		recError *namespaceError
	)

	if err != nil {
		if retError, ok = err.(*namespaceError); !ok {
			// got non-namespace error, wrap it with NamespaceErrGeneric
			retError = NamespaceErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use NamespaceErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type namespaceError
				if unwrappedSinkError, ok := unwrappedError.(*namespaceError); ok {
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
