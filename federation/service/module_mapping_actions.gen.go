package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// federation/service/module_mapping_actions.yaml

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
	moduleMappingActionProps struct {
		created *types.ModuleMapping
		mapping *types.ModuleMapping
		changed *types.ModuleMapping
		filter  *types.ModuleMappingFilter
	}

	moduleMappingAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *moduleMappingActionProps
	}

	moduleMappingError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *moduleMappingActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setCreated updates moduleMappingActionProps's created
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleMappingActionProps) setCreated(created *types.ModuleMapping) *moduleMappingActionProps {
	p.created = created
	return p
}

// setMapping updates moduleMappingActionProps's mapping
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleMappingActionProps) setMapping(mapping *types.ModuleMapping) *moduleMappingActionProps {
	p.mapping = mapping
	return p
}

// setChanged updates moduleMappingActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleMappingActionProps) setChanged(changed *types.ModuleMapping) *moduleMappingActionProps {
	p.changed = changed
	return p
}

// setFilter updates moduleMappingActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleMappingActionProps) setFilter(filter *types.ModuleMappingFilter) *moduleMappingActionProps {
	p.filter = filter
	return p
}

// serialize converts moduleMappingActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p moduleMappingActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.created != nil {
		m.Set("created.FederationModuleID", p.created.FederationModuleID, true)
		m.Set("created.ComposeModuleID", p.created.ComposeModuleID, true)
	}
	if p.mapping != nil {
		m.Set("mapping.FederationModuleID", p.mapping.FederationModuleID, true)
		m.Set("mapping.ComposeModuleID", p.mapping.ComposeModuleID, true)
	}
	if p.changed != nil {
		m.Set("changed.FederationModuleID", p.changed.FederationModuleID, true)
		m.Set("changed.ComposeModuleID", p.changed.ComposeModuleID, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.sort", p.filter.Sort, true)
		m.Set("filter.limit", p.filter.Limit, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p moduleMappingActionProps) tr(in string, err error) string {
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

	if p.created != nil {
		// replacement for "{created}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{created}",
			fns(
				p.created.FederationModuleID,
				p.created.ComposeModuleID,
			),
		)
		pairs = append(pairs, "{created.FederationModuleID}", fns(p.created.FederationModuleID))
		pairs = append(pairs, "{created.ComposeModuleID}", fns(p.created.ComposeModuleID))
	}

	if p.mapping != nil {
		// replacement for "{mapping}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{mapping}",
			fns(
				p.mapping.FederationModuleID,
				p.mapping.ComposeModuleID,
			),
		)
		pairs = append(pairs, "{mapping.FederationModuleID}", fns(p.mapping.FederationModuleID))
		pairs = append(pairs, "{mapping.ComposeModuleID}", fns(p.mapping.ComposeModuleID))
	}

	if p.changed != nil {
		// replacement for "{changed}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{changed}",
			fns(
				p.changed.FederationModuleID,
				p.changed.ComposeModuleID,
			),
		)
		pairs = append(pairs, "{changed.FederationModuleID}", fns(p.changed.FederationModuleID))
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
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
//
func (a *moduleMappingAction) String() string {
	var props = &moduleMappingActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *moduleMappingAction) LoggableAction() *actionlog.Action {
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
func (e *moduleMappingError) String() string {
	var props = &moduleMappingActionProps{}

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
func (e *moduleMappingError) Error() string {
	var props = &moduleMappingActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *moduleMappingError) Is(err error) bool {
	t, ok := err.(*moduleMappingError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *moduleMappingError) IsGeneric() bool {
	return e.error == "generic"
}

// Wrap wraps moduleMappingError around another error
//
// This function is auto-generated.
//
func (e *moduleMappingError) Wrap(err error) *moduleMappingError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *moduleMappingError) Unwrap() error {
	return e.wrap
}

func (e *moduleMappingError) LoggableAction() *actionlog.Action {
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

// ModuleMappingActionSearch returns "federation:module_mapping.search" error
//
// This function is auto-generated.
//
func ModuleMappingActionSearch(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "search",
		log:       "searched for modules",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleMappingActionLookup returns "federation:module_mapping.lookup" error
//
// This function is auto-generated.
//
func ModuleMappingActionLookup(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "lookup",
		log:       "looked-up for a module",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleMappingActionCreate returns "federation:module_mapping.create" error
//
// This function is auto-generated.
//
func ModuleMappingActionCreate(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "create",
		log:       "created module",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleMappingActionUpdate returns "federation:module_mapping.update" error
//
// This function is auto-generated.
//
func ModuleMappingActionUpdate(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "update",
		log:       "updated module",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleMappingActionDelete returns "federation:module_mapping.delete" error
//
// This function is auto-generated.
//
func ModuleMappingActionDelete(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "delete",
		log:       "deleted module",
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

// ModuleMappingErrGeneric returns "federation:module_mapping.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrGeneric(props ...*moduleMappingActionProps) *moduleMappingError {
	var e = &moduleMappingError{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *moduleMappingActionProps {
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

// ModuleMappingErrNotFound returns "federation:module_mapping.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleMappingErrNotFound(props ...*moduleMappingActionProps) *moduleMappingError {
	var e = &moduleMappingError{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		error:     "notFound",
		action:    "error",
		message:   "module mapping does not exist",
		log:       "module mapping does not exist",
		severity:  actionlog.Warning,
		props: func() *moduleMappingActionProps {
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

// ModuleMappingErrComposeModuleNotFound returns "federation:module_mapping.composeModuleNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleMappingErrComposeModuleNotFound(props ...*moduleMappingActionProps) *moduleMappingError {
	var e = &moduleMappingError{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		error:     "composeModuleNotFound",
		action:    "error",
		message:   "compose module not found",
		log:       "compose module not found",
		severity:  actionlog.Warning,
		props: func() *moduleMappingActionProps {
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

// ModuleMappingErrFederationModuleNotFound returns "federation:module_mapping.federationModuleNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ModuleMappingErrFederationModuleNotFound(props ...*moduleMappingActionProps) *moduleMappingError {
	var e = &moduleMappingError{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		error:     "federationModuleNotFound",
		action:    "error",
		message:   "federation module not found",
		log:       "federation module not found",
		severity:  actionlog.Warning,
		props: func() *moduleMappingActionProps {
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

// ModuleMappingErrNotAllowedToRead returns "federation:module_mapping.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrNotAllowedToRead(props ...*moduleMappingActionProps) *moduleMappingError {
	var e = &moduleMappingError{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this module",
		log:       "could not read module; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleMappingActionProps {
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

// ModuleMappingErrNotAllowedToListModules returns "federation:module_mapping.notAllowedToListModules" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrNotAllowedToListModules(props ...*moduleMappingActionProps) *moduleMappingError {
	var e = &moduleMappingError{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		error:     "notAllowedToListModules",
		action:    "error",
		message:   "not allowed to list modules",
		log:       "could not list modules; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleMappingActionProps {
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

// ModuleMappingErrNotAllowedToCreate returns "federation:module_mapping.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrNotAllowedToCreate(props ...*moduleMappingActionProps) *moduleMappingError {
	var e = &moduleMappingError{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create modules",
		log:       "could not create modules; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleMappingActionProps {
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

// ModuleMappingErrNotAllowedToUpdate returns "federation:module_mapping.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrNotAllowedToUpdate(props ...*moduleMappingActionProps) *moduleMappingError {
	var e = &moduleMappingError{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this module",
		log:       "could not update module; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleMappingActionProps {
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

// ModuleMappingErrNotAllowedToDelete returns "federation:module_mapping.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrNotAllowedToDelete(props ...*moduleMappingActionProps) *moduleMappingError {
	var e = &moduleMappingError{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this module",
		log:       "could not delete module; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *moduleMappingActionProps {
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
// action (optional) fn will be used to construct moduleMappingAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc moduleMapping) recordAction(ctx context.Context, props *moduleMappingActionProps, action func(...*moduleMappingActionProps) *moduleMappingAction, err error) error {
	var (
		ok bool

		// Return error
		retError *moduleMappingError

		// Recorder error
		recError *moduleMappingError
	)

	if err != nil {
		if retError, ok = err.(*moduleMappingError); !ok {
			// got non-moduleMapping error, wrap it with ModuleMappingErrGeneric
			retError = ModuleMappingErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use ModuleMappingErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type moduleMappingError
				if unwrappedSinkError, ok := unwrappedError.(*moduleMappingError); ok {
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
