package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// federation/service/shared_module_actions.yaml

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
)

type (
	sharedModuleActionProps struct {
		module  *types.SharedModule
		changed *types.SharedModule
		filter  *types.SharedModuleFilter
		node    *types.Node
	}

	sharedModuleAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *sharedModuleActionProps
	}

	sharedModuleError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *sharedModuleActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setModule updates sharedModuleActionProps's module
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sharedModuleActionProps) setModule(module *types.SharedModule) *sharedModuleActionProps {
	p.module = module
	return p
}

// setChanged updates sharedModuleActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sharedModuleActionProps) setChanged(changed *types.SharedModule) *sharedModuleActionProps {
	p.changed = changed
	return p
}

// setFilter updates sharedModuleActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sharedModuleActionProps) setFilter(filter *types.SharedModuleFilter) *sharedModuleActionProps {
	p.filter = filter
	return p
}

// setNode updates sharedModuleActionProps's node
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sharedModuleActionProps) setNode(node *types.Node) *sharedModuleActionProps {
	p.node = node
	return p
}

// serialize converts sharedModuleActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p sharedModuleActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.module != nil {
		m.Set("module.ID", p.module.ID, true)
	}
	if p.changed != nil {
		m.Set("changed.ID", p.changed.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.sort", p.filter.Sort, true)
		m.Set("filter.limit", p.filter.Limit, true)
	}
	if p.node != nil {
		m.Set("node.ID", p.node.ID, true)
		m.Set("node.Name", p.node.Name, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p sharedModuleActionProps) tr(in string, err error) string {
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
				p.module.ID,
			),
		)
		pairs = append(pairs, "{module.ID}", fns(p.module.ID))
	}

	if p.changed != nil {
		// replacement for "{changed}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{changed}",
			fns(
				p.changed.ID,
			),
		)
		pairs = append(pairs, "{changed.ID}", fns(p.changed.ID))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.Sort,
				p.filter.Limit,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.sort}", fns(p.filter.Sort))
		pairs = append(pairs, "{filter.limit}", fns(p.filter.Limit))
	}

	if p.node != nil {
		// replacement for "{node}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{node}",
			fns(
				p.node.ID,
				p.node.Name,
			),
		)
		pairs = append(pairs, "{node.ID}", fns(p.node.ID))
		pairs = append(pairs, "{node.Name}", fns(p.node.Name))
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
func (a *sharedModuleAction) String() string {
	var props = &sharedModuleActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *sharedModuleAction) LoggableAction() *actionlog.Action {
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
func (e *sharedModuleError) String() string {
	var props = &sharedModuleActionProps{}

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
func (e *sharedModuleError) Error() string {
	var props = &sharedModuleActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *sharedModuleError) Is(err error) bool {
	t, ok := err.(*sharedModuleError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *sharedModuleError) IsGeneric() bool {
	return e.error == "generic"
}

// Wrap wraps sharedModuleError around another error
//
// This function is auto-generated.
//
func (e *sharedModuleError) Wrap(err error) *sharedModuleError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *sharedModuleError) Unwrap() error {
	return e.wrap
}

func (e *sharedModuleError) LoggableAction() *actionlog.Action {
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

// SharedModuleActionSearch returns "federation:shared_module.search" error
//
// This function is auto-generated.
//
func SharedModuleActionSearch(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "search",
		log:       "searched for modules",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionLookup returns "federation:shared_module.lookup" error
//
// This function is auto-generated.
//
func SharedModuleActionLookup(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "lookup",
		log:       "looked-up for a {module}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionCreate returns "federation:shared_module.create" error
//
// This function is auto-generated.
//
func SharedModuleActionCreate(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "create",
		log:       "created {module}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionUpdate returns "federation:shared_module.update" error
//
// This function is auto-generated.
//
func SharedModuleActionUpdate(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "update",
		log:       "updated {module}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionDelete returns "federation:shared_module.delete" error
//
// This function is auto-generated.
//
func SharedModuleActionDelete(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "delete",
		log:       "deleted {module}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionUndelete returns "federation:shared_module.undelete" error
//
// This function is auto-generated.
//
func SharedModuleActionUndelete(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "undelete",
		log:       "undeleted {module}",
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

// SharedModuleErrGeneric returns "federation:shared_module.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrGeneric(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *sharedModuleActionProps {
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

// SharedModuleErrNotFound returns "federation:shared_module.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func SharedModuleErrNotFound(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "notFound",
		action:    "error",
		message:   "module does not exist",
		log:       "module does not exist",
		severity:  actionlog.Warning,
		props: func() *sharedModuleActionProps {
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

// SharedModuleErrInvalidID returns "federation:shared_module.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func SharedModuleErrInvalidID(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *sharedModuleActionProps {
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

// SharedModuleErrStaleData returns "federation:shared_module.staleData" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func SharedModuleErrStaleData(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "staleData",
		action:    "error",
		message:   "stale data",
		log:       "stale data",
		severity:  actionlog.Warning,
		props: func() *sharedModuleActionProps {
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

// SharedModuleErrNotAllowedToRead returns "federation:shared_module.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotAllowedToRead(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this module",
		log:       "could not read {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *sharedModuleActionProps {
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

// SharedModuleErrNotAllowedToListModules returns "federation:shared_module.notAllowedToListModules" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotAllowedToListModules(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "notAllowedToListModules",
		action:    "error",
		message:   "not allowed to list modules",
		log:       "could not list modules; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *sharedModuleActionProps {
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

// SharedModuleErrNotAllowedToCreate returns "federation:shared_module.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotAllowedToCreate(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create modules",
		log:       "could not create modules; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *sharedModuleActionProps {
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

// SharedModuleErrNotAllowedToUpdate returns "federation:shared_module.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotAllowedToUpdate(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this module",
		log:       "could not update {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *sharedModuleActionProps {
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

// SharedModuleErrNotAllowedToDelete returns "federation:shared_module.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotAllowedToDelete(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this module",
		log:       "could not delete {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *sharedModuleActionProps {
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

// SharedModuleErrNotAllowedToUndelete returns "federation:shared_module.notAllowedToUndelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotAllowedToUndelete(props ...*sharedModuleActionProps) *sharedModuleError {
	var e = &sharedModuleError{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this module",
		log:       "could not undelete {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *sharedModuleActionProps {
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
// action (optional) fn will be used to construct sharedModuleAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc sharedModule) recordAction(ctx context.Context, props *sharedModuleActionProps, action func(...*sharedModuleActionProps) *sharedModuleAction, err error) error {
	var (
		ok bool

		// Return error
		retError *sharedModuleError

		// Recorder error
		recError *sharedModuleError
	)

	if err != nil {
		if retError, ok = err.(*sharedModuleError); !ok {
			// got non-sharedModule error, wrap it with SharedModuleErrGeneric
			retError = SharedModuleErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use SharedModuleErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type sharedModuleError
				if unwrappedSinkError, ok := unwrappedError.(*sharedModuleError); ok {
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
