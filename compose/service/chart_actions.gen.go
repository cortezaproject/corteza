package service

// This file is auto-generated from compose/service/chart_actions.yaml
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
	chartActionProps struct {
		chart     *types.Chart
		changed   *types.Chart
		filter    *types.ChartFilter
		namespace *types.Namespace
	}

	chartAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *chartActionProps
	}

	chartError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *chartActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setChart updates chartActionProps's chart
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *chartActionProps) setChart(chart *types.Chart) *chartActionProps {
	p.chart = chart
	return p
}

// setChanged updates chartActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *chartActionProps) setChanged(changed *types.Chart) *chartActionProps {
	p.changed = changed
	return p
}

// setFilter updates chartActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *chartActionProps) setFilter(filter *types.ChartFilter) *chartActionProps {
	p.filter = filter
	return p
}

// setNamespace updates chartActionProps's namespace
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *chartActionProps) setNamespace(namespace *types.Namespace) *chartActionProps {
	p.namespace = namespace
	return p
}

// serialize converts chartActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p chartActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.chart != nil {
		m.Set("chart.name", p.chart.Name, true)
		m.Set("chart.handle", p.chart.Handle, true)
		m.Set("chart.ID", p.chart.ID, true)
		m.Set("chart.namespaceID", p.chart.NamespaceID, true)
	}
	if p.changed != nil {
		m.Set("changed.name", p.changed.Name, true)
		m.Set("changed.handle", p.changed.Handle, true)
		m.Set("changed.ID", p.changed.ID, true)
		m.Set("changed.namespaceID", p.changed.NamespaceID, true)
		m.Set("changed.config", p.changed.Config, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.handle", p.filter.Handle, true)
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
func (p chartActionProps) tr(in string, err error) string {
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

	if p.chart != nil {
		// replacement for "{chart}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{chart}",
			fns(
				p.chart.Name,
				p.chart.Handle,
				p.chart.ID,
				p.chart.NamespaceID,
			),
		)
		pairs = append(pairs, "{chart.name}", fns(p.chart.Name))
		pairs = append(pairs, "{chart.handle}", fns(p.chart.Handle))
		pairs = append(pairs, "{chart.ID}", fns(p.chart.ID))
		pairs = append(pairs, "{chart.namespaceID}", fns(p.chart.NamespaceID))
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
				p.changed.Config,
			),
		)
		pairs = append(pairs, "{changed.name}", fns(p.changed.Name))
		pairs = append(pairs, "{changed.handle}", fns(p.changed.Handle))
		pairs = append(pairs, "{changed.ID}", fns(p.changed.ID))
		pairs = append(pairs, "{changed.namespaceID}", fns(p.changed.NamespaceID))
		pairs = append(pairs, "{changed.config}", fns(p.changed.Config))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.Handle,
				p.filter.NamespaceID,
				p.filter.Sort,
				p.filter.Limit,
				p.filter.Offset,
				p.filter.Page,
				p.filter.PerPage,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.handle}", fns(p.filter.Handle))
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
func (a *chartAction) String() string {
	var props = &chartActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *chartAction) LoggableAction() *actionlog.Action {
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
func (e *chartError) String() string {
	var props = &chartActionProps{}

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
func (e *chartError) Error() string {
	var props = &chartActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *chartError) Is(Resource error) bool {
	t, ok := Resource.(*chartError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps chartError around another error
//
// This function is auto-generated.
//
func (e *chartError) Wrap(err error) *chartError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *chartError) Unwrap() error {
	return e.wrap
}

func (e *chartError) LoggableAction() *actionlog.Action {
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

// ChartActionSearch returns "compose:chart.search" error
//
// This function is auto-generated.
//
func ChartActionSearch(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "search",
		log:       "searched for charts",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionLookup returns "compose:chart.lookup" error
//
// This function is auto-generated.
//
func ChartActionLookup(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "lookup",
		log:       "looked-up for a {chart}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionCreate returns "compose:chart.create" error
//
// This function is auto-generated.
//
func ChartActionCreate(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "create",
		log:       "created {chart}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionUpdate returns "compose:chart.update" error
//
// This function is auto-generated.
//
func ChartActionUpdate(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "update",
		log:       "updated {chart}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionDelete returns "compose:chart.delete" error
//
// This function is auto-generated.
//
func ChartActionDelete(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "delete",
		log:       "deleted {chart}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionUndelete returns "compose:chart.undelete" error
//
// This function is auto-generated.
//
func ChartActionUndelete(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "undelete",
		log:       "undeleted {chart}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChartActionReorder returns "compose:chart.reorder" error
//
// This function is auto-generated.
//
func ChartActionReorder(props ...*chartActionProps) *chartAction {
	a := &chartAction{
		timestamp: time.Now(),
		resource:  "compose:chart",
		action:    "reorder",
		log:       "reordered {chart}",
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

// ChartErrGeneric returns "compose:chart.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChartErrGeneric(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *chartActionProps {
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

// ChartErrNotFound returns "compose:chart.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ChartErrNotFound(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "notFound",
		action:    "error",
		message:   "chart does not exist",
		log:       "chart does not exist",
		severity:  actionlog.Warning,
		props: func() *chartActionProps {
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

// ChartErrNamespaceNotFound returns "compose:chart.namespaceNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ChartErrNamespaceNotFound(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "namespaceNotFound",
		action:    "error",
		message:   "namespace does not exist",
		log:       "namespace does not exist",
		severity:  actionlog.Warning,
		props: func() *chartActionProps {
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

// ChartErrModuleNotFound returns "compose:chart.moduleNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ChartErrModuleNotFound(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "moduleNotFound",
		action:    "error",
		message:   "module does not exist",
		log:       "module does not exist",
		severity:  actionlog.Warning,
		props: func() *chartActionProps {
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

// ChartErrInvalidID returns "compose:chart.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ChartErrInvalidID(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *chartActionProps {
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

// ChartErrInvalidHandle returns "compose:chart.invalidHandle" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ChartErrInvalidHandle(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "invalidHandle",
		action:    "error",
		message:   "invalid handle",
		log:       "invalid handle",
		severity:  actionlog.Warning,
		props: func() *chartActionProps {
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

// ChartErrHandleNotUnique returns "compose:chart.handleNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ChartErrHandleNotUnique(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "handleNotUnique",
		action:    "error",
		message:   "handle not unique",
		log:       "used duplicate handle ({chart.handle}) for chart",
		severity:  actionlog.Warning,
		props: func() *chartActionProps {
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

// ChartErrStaleData returns "compose:chart.staleData" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ChartErrStaleData(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "staleData",
		action:    "error",
		message:   "stale data",
		log:       "stale data",
		severity:  actionlog.Warning,
		props: func() *chartActionProps {
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

// ChartErrInvalidNamespaceID returns "compose:chart.invalidNamespaceID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func ChartErrInvalidNamespaceID(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "invalidNamespaceID",
		action:    "error",
		message:   "invalid or missing namespace ID",
		log:       "invalid or missing namespace ID",
		severity:  actionlog.Warning,
		props: func() *chartActionProps {
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

// ChartErrNotAllowedToRead returns "compose:chart.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToRead(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this chart",
		log:       "could not read {chart}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *chartActionProps {
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

// ChartErrNotAllowedToReadNamespace returns "compose:chart.notAllowedToReadNamespace" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToReadNamespace(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "notAllowedToReadNamespace",
		action:    "error",
		message:   "not allowed to read this namespace",
		log:       "could not read namespace {namespace}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *chartActionProps {
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

// ChartErrNotAllowedToListCharts returns "compose:chart.notAllowedToListCharts" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToListCharts(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "notAllowedToListCharts",
		action:    "error",
		message:   "not allowed to list charts",
		log:       "could not list charts; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *chartActionProps {
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

// ChartErrNotAllowedToCreate returns "compose:chart.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToCreate(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create charts",
		log:       "could not create charts; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *chartActionProps {
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

// ChartErrNotAllowedToUpdate returns "compose:chart.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToUpdate(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this chart",
		log:       "could not update {chart}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *chartActionProps {
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

// ChartErrNotAllowedToDelete returns "compose:chart.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToDelete(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this chart",
		log:       "could not delete {chart}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *chartActionProps {
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

// ChartErrNotAllowedToUndelete returns "compose:chart.notAllowedToUndelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func ChartErrNotAllowedToUndelete(props ...*chartActionProps) *chartError {
	var e = &chartError{
		timestamp: time.Now(),
		resource:  "compose:chart",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this chart",
		log:       "could not undelete {chart}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *chartActionProps {
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
// action (optional) fn will be used to construct chartAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc chart) recordAction(ctx context.Context, props *chartActionProps, action func(...*chartActionProps) *chartAction, err error) error {
	var (
		ok bool

		// Return error
		retError *chartError

		// Recorder error
		recError *chartError
	)

	if err != nil {
		if retError, ok = err.(*chartError); !ok {
			// got non-chart error, wrap it with ChartErrGeneric
			retError = ChartErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use ChartErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type chartError
				if unwrappedSinkError, ok := unwrappedError.(*chartError); ok {
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
