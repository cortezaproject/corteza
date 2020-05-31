package service

// This file is auto-generated from system/service/statistics_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
)

type (
	statisticsActionProps struct {
	}

	statisticsAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *statisticsActionProps
	}

	statisticsError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *statisticsActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods

// serialize converts statisticsActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p statisticsActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p statisticsActionProps) tr(in string, err error) string {
	var (
		pairs = []string{"{err}"}
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
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
//
func (a *statisticsAction) String() string {
	var props = &statisticsActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *statisticsAction) LoggableAction() *actionlog.Action {
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
func (e *statisticsError) String() string {
	var props = &statisticsActionProps{}

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
func (e *statisticsError) Error() string {
	var props = &statisticsActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *statisticsError) Is(Resource error) bool {
	t, ok := Resource.(*statisticsError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps statisticsError around another error
//
// This function is auto-generated.
//
func (e *statisticsError) Wrap(err error) *statisticsError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *statisticsError) Unwrap() error {
	return e.wrap
}

func (e *statisticsError) LoggableAction() *actionlog.Action {
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

// StatisticsActionServe returns "system:statistics.serve" error
//
// This function is auto-generated.
//
func StatisticsActionServe(props ...*statisticsActionProps) *statisticsAction {
	a := &statisticsAction{
		timestamp: time.Now(),
		resource:  "system:statistics",
		action:    "serve",
		log:       "metrics served",
		severity:  actionlog.Debug,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// StatisticsErrGeneric returns "system:statistics.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func StatisticsErrGeneric(props ...*statisticsActionProps) *statisticsError {
	var e = &statisticsError{
		timestamp: time.Now(),
		resource:  "system:statistics",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *statisticsActionProps {
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

// StatisticsErrNotAllowedToReadStatistics returns "system:statistics.notAllowedToReadStatistics" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func StatisticsErrNotAllowedToReadStatistics(props ...*statisticsActionProps) *statisticsError {
	var e = &statisticsError{
		timestamp: time.Now(),
		resource:  "system:statistics",
		error:     "notAllowedToReadStatistics",
		action:    "error",
		message:   "not allowed to read statistics",
		log:       "not allowed to read statistics",
		severity:  actionlog.Warning,
		props: func() *statisticsActionProps {
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
// action (optional) fn will be used to construct statisticsAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc statistics) recordAction(ctx context.Context, props *statisticsActionProps, action func(...*statisticsActionProps) *statisticsAction, err error) error {
	var (
		ok bool

		// Return error
		retError *statisticsError

		// Recorder error
		recError *statisticsError
	)

	if err != nil {
		if retError, ok = err.(*statisticsError); !ok {
			// got non-statistics error, wrap it with StatisticsErrGeneric
			retError = StatisticsErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use StatisticsErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type statisticsError
				if unwrappedSinkError, ok := unwrappedError.(*statisticsError); ok {
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
