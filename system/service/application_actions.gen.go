package service

// This file is auto-generated from system/service/application_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	applicationActionProps struct {
		application *types.Application
		new         *types.Application
		update      *types.Application
		filter      *types.ApplicationFilter
	}

	applicationAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *applicationActionProps
	}

	applicationError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *applicationActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setApplication updates applicationActionProps's application
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *applicationActionProps) setApplication(application *types.Application) *applicationActionProps {
	p.application = application
	return p
}

// setNew updates applicationActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *applicationActionProps) setNew(new *types.Application) *applicationActionProps {
	p.new = new
	return p
}

// setUpdate updates applicationActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *applicationActionProps) setUpdate(update *types.Application) *applicationActionProps {
	p.update = update
	return p
}

// setFilter updates applicationActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *applicationActionProps) setFilter(filter *types.ApplicationFilter) *applicationActionProps {
	p.filter = filter
	return p
}

// serialize converts applicationActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p applicationActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.application != nil {
		m.Set("application.name", p.application.Name, true)
		m.Set("application.ID", p.application.ID, true)
	}
	if p.new != nil {
		m.Set("new.name", p.new.Name, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.name", p.update.Name, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.name", p.filter.Name, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p applicationActionProps) tr(in string, err error) string {
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

	if p.application != nil {
		// replacement for "{application}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{application}",
			fns(
				p.application.Name,
				p.application.ID,
			),
		)
		pairs = append(pairs, "{application.name}", fns(p.application.Name))
		pairs = append(pairs, "{application.ID}", fns(p.application.ID))
	}

	if p.new != nil {
		// replacement for "{new}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{new}",
			fns(
				p.new.Name,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{new.name}", fns(p.new.Name))
		pairs = append(pairs, "{new.ID}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{update}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{update}",
			fns(
				p.update.Name,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{update.name}", fns(p.update.Name))
		pairs = append(pairs, "{update.ID}", fns(p.update.ID))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.Name,
				p.filter.Deleted,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.name}", fns(p.filter.Name))
		pairs = append(pairs, "{filter.deleted}", fns(p.filter.Deleted))
		pairs = append(pairs, "{filter.sort}", fns(p.filter.Sort))
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
func (a *applicationAction) String() string {
	var props = &applicationActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *applicationAction) LoggableAction() *actionlog.Action {
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
func (e *applicationError) String() string {
	var props = &applicationActionProps{}

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
func (e *applicationError) Error() string {
	var props = &applicationActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *applicationError) Is(Resource error) bool {
	t, ok := Resource.(*applicationError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps applicationError around another error
//
// This function is auto-generated.
//
func (e *applicationError) Wrap(err error) *applicationError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *applicationError) Unwrap() error {
	return e.wrap
}

func (e *applicationError) LoggableAction() *actionlog.Action {
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

// ApplicationActionSearch returns "system:application.search" error
//
// This function is auto-generated.
//
func ApplicationActionSearch(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "search",
		log:       "searched for applications",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionLookup returns "system:application.lookup" error
//
// This function is auto-generated.
//
func ApplicationActionLookup(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "lookup",
		log:       "looked-up for a {application}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionCreate returns "system:application.create" error
//
// This function is auto-generated.
//
func ApplicationActionCreate(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "create",
		log:       "created {application}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionUpdate returns "system:application.update" error
//
// This function is auto-generated.
//
func ApplicationActionUpdate(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "update",
		log:       "updated {application}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionDelete returns "system:application.delete" error
//
// This function is auto-generated.
//
func ApplicationActionDelete(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "delete",
		log:       "deleted {application}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionUndelete returns "system:application.undelete" error
//
// This function is auto-generated.
//
func ApplicationActionUndelete(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "undelete",
		log:       "undeleted {application}",
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

// ApplicationErrGeneric returns "system:application.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ApplicationErrGeneric(props ...*applicationActionProps) *applicationError {
	var e = &applicationError{
		timestamp: time.Now(),
		resource:  "system:application",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *applicationActionProps {
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

// ApplicationErrNotFound returns "system:application.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ApplicationErrNotFound(props ...*applicationActionProps) *applicationError {
	var e = &applicationError{
		timestamp: time.Now(),
		resource:  "system:application",
		error:     "notFound",
		action:    "error",
		message:   "application not found",
		log:       "application not found",
		severity:  actionlog.Warning,
		props: func() *applicationActionProps {
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

// ApplicationErrInvalidID returns "system:application.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ApplicationErrInvalidID(props ...*applicationActionProps) *applicationError {
	var e = &applicationError{
		timestamp: time.Now(),
		resource:  "system:application",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *applicationActionProps {
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

// ApplicationErrNotAllowedToRead returns "system:application.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToRead(props ...*applicationActionProps) *applicationError {
	var e = &applicationError{
		timestamp: time.Now(),
		resource:  "system:application",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this application",
		log:       "failed to read {application.name}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *applicationActionProps {
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

// ApplicationErrNotAllowedToListApplications returns "system:application.notAllowedToListApplications" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToListApplications(props ...*applicationActionProps) *applicationError {
	var e = &applicationError{
		timestamp: time.Now(),
		resource:  "system:application",
		error:     "notAllowedToListApplications",
		action:    "error",
		message:   "not allowed to list applications",
		log:       "failed to list application; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *applicationActionProps {
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

// ApplicationErrNotAllowedToCreate returns "system:application.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToCreate(props ...*applicationActionProps) *applicationError {
	var e = &applicationError{
		timestamp: time.Now(),
		resource:  "system:application",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create applications",
		log:       "failed to create application; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *applicationActionProps {
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

// ApplicationErrNotAllowedToUpdate returns "system:application.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToUpdate(props ...*applicationActionProps) *applicationError {
	var e = &applicationError{
		timestamp: time.Now(),
		resource:  "system:application",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this application",
		log:       "failed to update {application.name}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *applicationActionProps {
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

// ApplicationErrNotAllowedToDelete returns "system:application.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToDelete(props ...*applicationActionProps) *applicationError {
	var e = &applicationError{
		timestamp: time.Now(),
		resource:  "system:application",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this application",
		log:       "failed to delete {application.name}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *applicationActionProps {
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

// ApplicationErrNotAllowedToUndelete returns "system:application.notAllowedToUndelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToUndelete(props ...*applicationActionProps) *applicationError {
	var e = &applicationError{
		timestamp: time.Now(),
		resource:  "system:application",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this application",
		log:       "failed to undelete {application.name}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *applicationActionProps {
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
// action (optional) fn will be used to construct applicationAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc application) recordAction(ctx context.Context, props *applicationActionProps, action func(...*applicationActionProps) *applicationAction, err error) error {
	var (
		ok bool

		// Return error
		retError *applicationError

		// Recorder error
		recError *applicationError
	)

	if err != nil {
		if retError, ok = err.(*applicationError); !ok {
			// got non-application error, wrap it with ApplicationErrGeneric
			retError = ApplicationErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use ApplicationErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type applicationError
				if unwrappedSinkError, ok := unwrappedError.(*applicationError); ok {
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
