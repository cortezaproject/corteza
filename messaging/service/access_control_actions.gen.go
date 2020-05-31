package service

// This file is auto-generated from messaging/service/access_control_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	accessControlActionProps struct {
		rule *permissions.Rule
	}

	accessControlAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *accessControlActionProps
	}

	accessControlError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *accessControlActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setRule updates accessControlActionProps's rule
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *accessControlActionProps) setRule(rule *permissions.Rule) *accessControlActionProps {
	p.rule = rule
	return p
}

// serialize converts accessControlActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p accessControlActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.rule != nil {
		m.Set("rule.operation", p.rule.Operation, true)
		m.Set("rule.roleID", p.rule.RoleID, true)
		m.Set("rule.access", p.rule.Access, true)
		m.Set("rule.resource", p.rule.Resource, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p accessControlActionProps) tr(in string, err error) string {
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

	if p.rule != nil {
		// replacement for "{rule}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{rule}",
			fns(
				p.rule.Operation,
				p.rule.RoleID,
				p.rule.Access,
				p.rule.Resource,
			),
		)
		pairs = append(pairs, "{rule.operation}", fns(p.rule.Operation))
		pairs = append(pairs, "{rule.roleID}", fns(p.rule.RoleID))
		pairs = append(pairs, "{rule.access}", fns(p.rule.Access))
		pairs = append(pairs, "{rule.resource}", fns(p.rule.Resource))
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
func (a *accessControlAction) String() string {
	var props = &accessControlActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *accessControlAction) LoggableAction() *actionlog.Action {
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
func (e *accessControlError) String() string {
	var props = &accessControlActionProps{}

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
func (e *accessControlError) Error() string {
	var props = &accessControlActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *accessControlError) Is(Resource error) bool {
	t, ok := Resource.(*accessControlError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps accessControlError around another error
//
// This function is auto-generated.
//
func (e *accessControlError) Wrap(err error) *accessControlError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *accessControlError) Unwrap() error {
	return e.wrap
}

func (e *accessControlError) LoggableAction() *actionlog.Action {
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

// AccessControlActionGrant returns "messaging:access_control.grant" error
//
// This function is auto-generated.
//
func AccessControlActionGrant(props ...*accessControlActionProps) *accessControlAction {
	a := &accessControlAction{
		timestamp: time.Now(),
		resource:  "messaging:access_control",
		action:    "grant",
		log:       "grant",
		severity:  actionlog.Error,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// AccessControlErrGeneric returns "messaging:access_control.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func AccessControlErrGeneric(props ...*accessControlActionProps) *accessControlError {
	var e = &accessControlError{
		timestamp: time.Now(),
		resource:  "messaging:access_control",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *accessControlActionProps {
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

// AccessControlErrNotAllowedToSetPermissions returns "messaging:access_control.notAllowedToSetPermissions" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AccessControlErrNotAllowedToSetPermissions(props ...*accessControlActionProps) *accessControlError {
	var e = &accessControlError{
		timestamp: time.Now(),
		resource:  "messaging:access_control",
		error:     "notAllowedToSetPermissions",
		action:    "error",
		message:   "not allowed to set permissions",
		log:       "not allowed to set permissions",
		severity:  actionlog.Alert,
		props: func() *accessControlActionProps {
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
// action (optional) fn will be used to construct accessControlAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc accessControl) recordAction(ctx context.Context, props *accessControlActionProps, action func(...*accessControlActionProps) *accessControlAction, err error) error {
	var (
		ok bool

		// Return error
		retError *accessControlError

		// Recorder error
		recError *accessControlError
	)

	if err != nil {
		if retError, ok = err.(*accessControlError); !ok {
			// got non-accessControl error, wrap it with AccessControlErrGeneric
			retError = AccessControlErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use AccessControlErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type accessControlError
				if unwrappedSinkError, ok := unwrappedError.(*accessControlError); ok {
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
