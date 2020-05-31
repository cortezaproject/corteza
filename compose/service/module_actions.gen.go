package service

// This file is auto-generated from compose/service/module_actions.yaml
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

	moduleError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *moduleActionProps
	}
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
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleActionProps) setModule(module *types.Module) *moduleActionProps {
	p.module = module
	return p
}

// setChanged updates moduleActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleActionProps) setChanged(changed *types.Module) *moduleActionProps {
	p.changed = changed
	return p
}

// setFilter updates moduleActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleActionProps) setFilter(filter *types.ModuleFilter) *moduleActionProps {
	p.filter = filter
	return p
}

// setNamespace updates moduleActionProps's namespace
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleActionProps) setNamespace(namespace *types.Namespace) *moduleActionProps {
	p.namespace = namespace
	return p
}

// serialize converts moduleActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p moduleActionProps) serialize() actionlog.Meta {
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
func (p moduleActionProps) tr(in string, err error) string {
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

	if p.module != nil {
		// replacement for "{module}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{module}",
			fns(
				p.module.Name,
				p.module.Handle,
				p.module.ID,
				p.module.NamespaceID,
			),
		)
		pairs = append(pairs, "{module.name}", fns(p.module.Name))
		pairs = append(pairs, "{module.handle}", fns(p.module.Handle))
		pairs = append(pairs, "{module.ID}", fns(p.module.ID))
		pairs = append(pairs, "{module.namespaceID}", fns(p.module.NamespaceID))
	}

	if p.changed != nil {
		// replacement for "{changed}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{changed}",
			fns(
				p.changed.Name,
				p.changed.Handle,
				p.changed.ID,
				p.changed.NamespaceID,
				p.changed.Meta,
				p.changed.Fields,
			),
		)
		pairs = append(pairs, "{changed.name}", fns(p.changed.Name))
		pairs = append(pairs, "{changed.handle}", fns(p.changed.Handle))
		pairs = append(pairs, "{changed.ID}", fns(p.changed.ID))
		pairs = append(pairs, "{changed.namespaceID}", fns(p.changed.NamespaceID))
		pairs = append(pairs, "{changed.meta}", fns(p.changed.Meta))
		pairs = append(pairs, "{changed.fields}", fns(p.changed.Fields))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.Name,
				p.filter.Handle,
				p.filter.Name,
				p.filter.NamespaceID,
				p.filter.Sort,
				p.filter.Limit,
				p.filter.Offset,
				p.filter.Page,
				p.filter.PerPage,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.name}", fns(p.filter.Name))
		pairs = append(pairs, "{filter.handle}", fns(p.filter.Handle))
		pairs = append(pairs, "{filter.name}", fns(p.filter.Name))
		pairs = append(pairs, "{filter.namespaceID}", fns(p.filter.NamespaceID))
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
func (a *moduleAction) String() string {
	var props = &moduleActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *moduleAction) LoggableAction() *actionlog.Action {
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
func (e *moduleError) String() string {
	var props = &moduleActionProps{}

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
func (e *moduleError) Error() string {
	var props = &moduleActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *moduleError) Is(Resource error) bool {
	t, ok := Resource.(*moduleError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps moduleError around another error
//
// This function is auto-generated.
//
func (e *moduleError) Wrap(err error) *moduleError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *moduleError) Unwrap() error {
	return e.wrap
}

func (e *moduleError) LoggableAction() *actionlog.Action {
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

// ModuleActionSearch returns "compose:module.search" error
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

// ModuleActionLookup returns "compose:module.lookup" error
//
// This function is auto-generated.
//
func ModuleActionLookup(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "lookup",
		log:       "looked-up for a {module}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleActionCreate returns "compose:module.create" error
//
// This function is auto-generated.
//
func ModuleActionCreate(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "create",
		log:       "created {module}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleActionUpdate returns "compose:module.update" error
//
// This function is auto-generated.
//
func ModuleActionUpdate(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "update",
		log:       "updated {module}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleActionDelete returns "compose:module.delete" error
//
// This function is auto-generated.
//
func ModuleActionDelete(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "delete",
		log:       "deleted {module}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleActionUndelete returns "compose:module.undelete" error
//
// This function is auto-generated.
//
func ModuleActionUndelete(props ...*moduleActionProps) *moduleAction {
	a := &moduleAction{
		timestamp: time.Now(),
		resource:  "compose:module",
		action:    "undelete",
		log:       "undeleted {module}",
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

// ModuleErrGeneric returns "compose:module.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleErrGeneric(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *moduleActionProps {
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

// ModuleErrNotFound returns "compose:module.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleErrNotFound(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "notFound",
		action:    "error",
		message:   "module does not exist",
		log:       "module does not exist",
		severity:  actionlog.Warning,
		props: func() *moduleActionProps {
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

// ModuleErrNamespaceNotFound returns "compose:module.namespaceNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleErrNamespaceNotFound(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "namespaceNotFound",
		action:    "error",
		message:   "namespace does not exist",
		log:       "namespace does not exist",
		severity:  actionlog.Warning,
		props: func() *moduleActionProps {
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

// ModuleErrInvalidID returns "compose:module.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleErrInvalidID(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *moduleActionProps {
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

// ModuleErrInvalidHandle returns "compose:module.invalidHandle" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleErrInvalidHandle(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "invalidHandle",
		action:    "error",
		message:   "invalid handle",
		log:       "invalid handle",
		severity:  actionlog.Warning,
		props: func() *moduleActionProps {
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

// ModuleErrHandleNotUnique returns "compose:module.handleNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleErrHandleNotUnique(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "handleNotUnique",
		action:    "error",
		message:   "handle not unique",
		log:       "used duplicate handle ({module.handle}) for module",
		severity:  actionlog.Warning,
		props: func() *moduleActionProps {
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

// ModuleErrNameNotUnique returns "compose:module.nameNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleErrNameNotUnique(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "nameNotUnique",
		action:    "error",
		message:   "name not unique",
		log:       "used duplicate username ({module.name}) for module",
		severity:  actionlog.Warning,
		props: func() *moduleActionProps {
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

// ModuleErrStaleData returns "compose:module.staleData" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleErrStaleData(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "staleData",
		action:    "error",
		message:   "stale data",
		log:       "stale data",
		severity:  actionlog.Warning,
		props: func() *moduleActionProps {
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

// ModuleErrInvalidNamespaceID returns "compose:module.invalidNamespaceID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleErrInvalidNamespaceID(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "invalidNamespaceID",
		action:    "error",
		message:   "invalid or missing namespace ID",
		log:       "invalid or missing namespace ID",
		severity:  actionlog.Warning,
		props: func() *moduleActionProps {
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

// ModuleErrNotAllowedToRead returns "compose:module.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToRead(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this module",
		log:       "could not read {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleActionProps {
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

// ModuleErrNotAllowedToReadNamespace returns "compose:module.notAllowedToReadNamespace" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToReadNamespace(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "notAllowedToReadNamespace",
		action:    "error",
		message:   "not allowed to read this namespace",
		log:       "could not read namespace {namespace}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleActionProps {
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

// ModuleErrNotAllowedToListModules returns "compose:module.notAllowedToListModules" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToListModules(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "notAllowedToListModules",
		action:    "error",
		message:   "not allowed to list modules",
		log:       "could not list modules; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleActionProps {
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

// ModuleErrNotAllowedToCreate returns "compose:module.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToCreate(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create modules",
		log:       "could not create modules; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleActionProps {
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

// ModuleErrNotAllowedToUpdate returns "compose:module.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToUpdate(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this module",
		log:       "could not update {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleActionProps {
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

// ModuleErrNotAllowedToDelete returns "compose:module.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToDelete(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this module",
		log:       "could not delete {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleActionProps {
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

// ModuleErrNotAllowedToUndelete returns "compose:module.notAllowedToUndelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleErrNotAllowedToUndelete(props ...*moduleActionProps) *moduleError {
	var e = &moduleError{
		timestamp: time.Now(),
		resource:  "compose:module",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this module",
		log:       "could not undelete {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleActionProps {
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
// action (optional) fn will be used to construct moduleAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc module) recordAction(ctx context.Context, props *moduleActionProps, action func(...*moduleActionProps) *moduleAction, err error) error {
	var (
		ok bool

		// Return error
		retError *moduleError

		// Recorder error
		recError *moduleError
	)

	if err != nil {
		if retError, ok = err.(*moduleError); !ok {
			// got non-module error, wrap it with ModuleErrGeneric
			retError = ModuleErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use ModuleErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type moduleError
				if unwrappedSinkError, ok := unwrappedError.(*moduleError); ok {
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
