package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// federation/service/exposed_module_actions.yaml

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
	exposedModuleActionProps struct {
		module  *types.ExposedModule
		changed *types.ExposedModule
		filter  *types.ExposedModuleFilter
		node    *types.Node
	}

	exposedModuleAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *exposedModuleActionProps
	}

	exposedModuleError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *exposedModuleActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setModule updates exposedModuleActionProps's module
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setModule(module *types.ExposedModule) *exposedModuleActionProps {
	p.module = module
	return p
}

// setChanged updates exposedModuleActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setChanged(changed *types.ExposedModule) *exposedModuleActionProps {
	p.changed = changed
	return p
}

// setFilter updates exposedModuleActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setFilter(filter *types.ExposedModuleFilter) *exposedModuleActionProps {
	p.filter = filter
	return p
}

// setNode updates exposedModuleActionProps's node
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setNode(node *types.Node) *exposedModuleActionProps {
	p.node = node
	return p
}

// serialize converts exposedModuleActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p exposedModuleActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.module != nil {
		m.Set("module.ID", p.module.ID, true)
		m.Set("module.ComposeModuleID", p.module.ComposeModuleID, true)
	}
	if p.changed != nil {
		m.Set("changed.ID", p.changed.ID, true)
		m.Set("changed.ComposeModuleID", p.changed.ComposeModuleID, true)
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
func (p exposedModuleActionProps) tr(in string, err error) string {
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
				p.module.ComposeModuleID,
			),
		)
		pairs = append(pairs, "{module.ID}", fns(p.module.ID))
		pairs = append(pairs, "{module.ComposeModuleID}", fns(p.module.ComposeModuleID))
	}

	if p.changed != nil {
		// replacement for "{changed}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{changed}",
			fns(
				p.changed.ID,
				p.changed.ComposeModuleID,
			),
		)
		pairs = append(pairs, "{changed.ID}", fns(p.changed.ID))
		pairs = append(pairs, "{changed.ComposeModuleID}", fns(p.changed.ComposeModuleID))
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
func (a *exposedModuleAction) String() string {
	var props = &exposedModuleActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *exposedModuleAction) LoggableAction() *actionlog.Action {
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
func (e *exposedModuleError) String() string {
	var props = &exposedModuleActionProps{}

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
func (e *exposedModuleError) Error() string {
	var props = &exposedModuleActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *exposedModuleError) Is(err error) bool {
	t, ok := err.(*exposedModuleError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *exposedModuleError) IsGeneric() bool {
	return e.error == "generic"
}

// Wrap wraps exposedModuleError around another error
//
// This function is auto-generated.
//
func (e *exposedModuleError) Wrap(err error) *exposedModuleError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *exposedModuleError) Unwrap() error {
	return e.wrap
}

func (e *exposedModuleError) LoggableAction() *actionlog.Action {
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

// ExposedModuleActionSearch returns "federation:exposed_module.search" error
//
// This function is auto-generated.
//
func ExposedModuleActionSearch(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "search",
		log:       "searched for modules",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionLookup returns "federation:exposed_module.lookup" error
//
// This function is auto-generated.
//
func ExposedModuleActionLookup(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "lookup",
		log:       "looked-up for a {module}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionCreate returns "federation:exposed_module.create" error
//
// This function is auto-generated.
//
func ExposedModuleActionCreate(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "create",
		log:       "created {module}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionUpdate returns "federation:exposed_module.update" error
//
// This function is auto-generated.
//
func ExposedModuleActionUpdate(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "update",
		log:       "updated {module}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionDelete returns "federation:exposed_module.delete" error
//
// This function is auto-generated.
//
func ExposedModuleActionDelete(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "delete",
		log:       "deleted {module}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionUndelete returns "federation:exposed_module.undelete" error
//
// This function is auto-generated.
//
func ExposedModuleActionUndelete(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
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

// ExposedModuleErrGeneric returns "federation:exposed_module.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrGeneric(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrNotFound returns "federation:exposed_module.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotFound(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "notFound",
		action:    "error",
		message:   "module does not exist",
		log:       "module does not exist",
		severity:  actionlog.Warning,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrInvalidID returns "federation:exposed_module.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ExposedModuleErrInvalidID(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrStaleData returns "federation:exposed_module.staleData" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ExposedModuleErrStaleData(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "staleData",
		action:    "error",
		message:   "stale data",
		log:       "stale data",
		severity:  actionlog.Warning,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrNotUnique returns "federation:exposed_module.notUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotUnique(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "notUnique",
		action:    "error",
		message:   "node not unique",
		log:       "used duplicate node TODO - module.node_id for this compose module TODO - module.rel_compose_module",
		severity:  actionlog.Warning,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrComposeModuleNotFound returns "federation:exposed_module.composeModuleNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ExposedModuleErrComposeModuleNotFound(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "composeModuleNotFound",
		action:    "error",
		message:   "compose module not found",
		log:       "compose module not found",
		severity:  actionlog.Warning,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrNotAllowedToRead returns "federation:exposed_module.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotAllowedToRead(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this module",
		log:       "could not read {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrNotAllowedToListModules returns "federation:exposed_module.notAllowedToListModules" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotAllowedToListModules(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "notAllowedToListModules",
		action:    "error",
		message:   "not allowed to list modules",
		log:       "could not list modules; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrNotAllowedToCreate returns "federation:exposed_module.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotAllowedToCreate(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create modules",
		log:       "could not create modules; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrNotAllowedToUpdate returns "federation:exposed_module.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotAllowedToUpdate(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this module",
		log:       "could not update {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrNotAllowedToDelete returns "federation:exposed_module.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotAllowedToDelete(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this module",
		log:       "could not delete {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *exposedModuleActionProps {
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

// ExposedModuleErrNotAllowedToUndelete returns "federation:exposed_module.notAllowedToUndelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotAllowedToUndelete(props ...*exposedModuleActionProps) *exposedModuleError {
	var e = &exposedModuleError{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this module",
		log:       "could not undelete {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *exposedModuleActionProps {
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
// action (optional) fn will be used to construct exposedModuleAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc exposedModule) recordAction(ctx context.Context, props *exposedModuleActionProps, action func(...*exposedModuleActionProps) *exposedModuleAction, err error) error {
	var (
		ok bool

		// Return error
		retError *exposedModuleError

		// Recorder error
		recError *exposedModuleError
	)

	if err != nil {
		if retError, ok = err.(*exposedModuleError); !ok {
			// got non-exposedModule error, wrap it with ExposedModuleErrGeneric
			retError = ExposedModuleErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use ExposedModuleErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type exposedModuleError
				if unwrappedSinkError, ok := unwrappedError.(*exposedModuleError); ok {
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
