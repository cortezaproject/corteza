package service

// This file is auto-generated from compose/service/record_actions.yaml
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
	recordActionProps struct {
		record        *types.Record
		changed       *types.Record
		filter        *types.RecordFilter
		namespace     *types.Namespace
		module        *types.Module
		bulkOperation string
		field         string
		value         string
		valueErrors   *types.RecordValueErrorSet
	}

	recordAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *recordActionProps
	}

	recordError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *recordActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setRecord updates recordActionProps's record
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *recordActionProps) setRecord(record *types.Record) *recordActionProps {
	p.record = record
	return p
}

// setChanged updates recordActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *recordActionProps) setChanged(changed *types.Record) *recordActionProps {
	p.changed = changed
	return p
}

// setFilter updates recordActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *recordActionProps) setFilter(filter *types.RecordFilter) *recordActionProps {
	p.filter = filter
	return p
}

// setNamespace updates recordActionProps's namespace
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *recordActionProps) setNamespace(namespace *types.Namespace) *recordActionProps {
	p.namespace = namespace
	return p
}

// setModule updates recordActionProps's module
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *recordActionProps) setModule(module *types.Module) *recordActionProps {
	p.module = module
	return p
}

// setBulkOperation updates recordActionProps's bulkOperation
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *recordActionProps) setBulkOperation(bulkOperation string) *recordActionProps {
	p.bulkOperation = bulkOperation
	return p
}

// setField updates recordActionProps's field
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *recordActionProps) setField(field string) *recordActionProps {
	p.field = field
	return p
}

// setValue updates recordActionProps's value
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *recordActionProps) setValue(value string) *recordActionProps {
	p.value = value
	return p
}

// setValueErrors updates recordActionProps's valueErrors
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *recordActionProps) setValueErrors(valueErrors *types.RecordValueErrorSet) *recordActionProps {
	p.valueErrors = valueErrors
	return p
}

// serialize converts recordActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p recordActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.record != nil {
		m.Set("record.ID", p.record.ID, true)
		m.Set("record.moduleID", p.record.ModuleID, true)
		m.Set("record.namespaceID", p.record.NamespaceID, true)
		m.Set("record.ownedBy", p.record.OwnedBy, true)
	}
	if p.changed != nil {
		m.Set("changed.ID", p.changed.ID, true)
		m.Set("changed.moduleID", p.changed.ModuleID, true)
		m.Set("changed.namespaceID", p.changed.NamespaceID, true)
		m.Set("changed.ownedBy", p.changed.OwnedBy, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.namespaceID", p.filter.NamespaceID, true)
		m.Set("filter.moduleID", p.filter.ModuleID, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
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
	if p.module != nil {
		m.Set("module.name", p.module.Name, true)
		m.Set("module.handle", p.module.Handle, true)
		m.Set("module.ID", p.module.ID, true)
		m.Set("module.namespaceID", p.module.NamespaceID, true)
	}
	m.Set("bulkOperation", p.bulkOperation, true)
	m.Set("field", p.field, true)
	m.Set("value", p.value, true)
	if p.valueErrors != nil {
		m.Set("valueErrors.set", p.valueErrors.Set, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p recordActionProps) tr(in string, err error) string {
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

	if p.record != nil {
		// replacement for "{record}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{record}",
			fns(
				p.record.ID,
				p.record.ModuleID,
				p.record.NamespaceID,
				p.record.OwnedBy,
			),
		)
		pairs = append(pairs, "{record.ID}", fns(p.record.ID))
		pairs = append(pairs, "{record.moduleID}", fns(p.record.ModuleID))
		pairs = append(pairs, "{record.namespaceID}", fns(p.record.NamespaceID))
		pairs = append(pairs, "{record.ownedBy}", fns(p.record.OwnedBy))
	}

	if p.changed != nil {
		// replacement for "{changed}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{changed}",
			fns(
				p.changed.ID,
				p.changed.ModuleID,
				p.changed.NamespaceID,
				p.changed.OwnedBy,
			),
		)
		pairs = append(pairs, "{changed.ID}", fns(p.changed.ID))
		pairs = append(pairs, "{changed.moduleID}", fns(p.changed.ModuleID))
		pairs = append(pairs, "{changed.namespaceID}", fns(p.changed.NamespaceID))
		pairs = append(pairs, "{changed.ownedBy}", fns(p.changed.OwnedBy))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.NamespaceID,
				p.filter.ModuleID,
				p.filter.Deleted,
				p.filter.Sort,
				p.filter.Limit,
				p.filter.Offset,
				p.filter.Page,
				p.filter.PerPage,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.namespaceID}", fns(p.filter.NamespaceID))
		pairs = append(pairs, "{filter.moduleID}", fns(p.filter.ModuleID))
		pairs = append(pairs, "{filter.deleted}", fns(p.filter.Deleted))
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
	pairs = append(pairs, "{bulkOperation}", fns(p.bulkOperation))
	pairs = append(pairs, "{field}", fns(p.field))
	pairs = append(pairs, "{value}", fns(p.value))

	if p.valueErrors != nil {
		// replacement for "{valueErrors}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{valueErrors}",
			fns(
				p.valueErrors.Set,
			),
		)
		pairs = append(pairs, "{valueErrors.set}", fns(p.valueErrors.Set))
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
func (a *recordAction) String() string {
	var props = &recordActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *recordAction) LoggableAction() *actionlog.Action {
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
func (e *recordError) String() string {
	var props = &recordActionProps{}

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
func (e *recordError) Error() string {
	var props = &recordActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *recordError) Is(Resource error) bool {
	t, ok := Resource.(*recordError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps recordError around another error
//
// This function is auto-generated.
//
func (e *recordError) Wrap(err error) *recordError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *recordError) Unwrap() error {
	return e.wrap
}

func (e *recordError) LoggableAction() *actionlog.Action {
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

// RecordActionSearch returns "compose:record.search" error
//
// This function is auto-generated.
//
func RecordActionSearch(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "search",
		log:       "searched for records",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionLookup returns "compose:record.lookup" error
//
// This function is auto-generated.
//
func RecordActionLookup(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "lookup",
		log:       "looked-up for a {record}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionReport returns "compose:record.report" error
//
// This function is auto-generated.
//
func RecordActionReport(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "report",
		log:       "report generated",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionBulk returns "compose:record.bulk" error
//
// This function is auto-generated.
//
func RecordActionBulk(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "bulk",
		log:       "bulk record operation",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionCreate returns "compose:record.create" error
//
// This function is auto-generated.
//
func RecordActionCreate(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "create",
		log:       "created {record}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionUpdate returns "compose:record.update" error
//
// This function is auto-generated.
//
func RecordActionUpdate(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "update",
		log:       "updated {record}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionDelete returns "compose:record.delete" error
//
// This function is auto-generated.
//
func RecordActionDelete(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "delete",
		log:       "deleted {record}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionUndelete returns "compose:record.undelete" error
//
// This function is auto-generated.
//
func RecordActionUndelete(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "undelete",
		log:       "undeleted {record}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionImport returns "compose:record.import" error
//
// This function is auto-generated.
//
func RecordActionImport(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "import",
		log:       "records imported",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionExport returns "compose:record.export" error
//
// This function is auto-generated.
//
func RecordActionExport(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "export",
		log:       "records exported",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionOrganize returns "compose:record.organize" error
//
// This function is auto-generated.
//
func RecordActionOrganize(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "organize",
		log:       "records organized",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorInvoked returns "compose:record.iteratorInvoked" error
//
// This function is auto-generated.
//
func RecordActionIteratorInvoked(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorInvoked",
		log:       "iterator invoked",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorIteration returns "compose:record.iteratorIteration" error
//
// This function is auto-generated.
//
func RecordActionIteratorIteration(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorIteration",
		log:       "processed record iteration",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorClone returns "compose:record.iteratorClone" error
//
// This function is auto-generated.
//
func RecordActionIteratorClone(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorClone",
		log:       "cloned record in iteration",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorUpdate returns "compose:record.iteratorUpdate" error
//
// This function is auto-generated.
//
func RecordActionIteratorUpdate(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorUpdate",
		log:       "updated record in iteration",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorDelete returns "compose:record.iteratorDelete" error
//
// This function is auto-generated.
//
func RecordActionIteratorDelete(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorDelete",
		log:       "deleted record in iteration",
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

// RecordErrGeneric returns "compose:record.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrGeneric(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrNotFound returns "compose:record.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RecordErrNotFound(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notFound",
		action:    "error",
		message:   "record not found",
		log:       "record not found",
		severity:  actionlog.Warning,
		props: func() *recordActionProps {
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

// RecordErrNamespaceNotFound returns "compose:record.namespaceNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RecordErrNamespaceNotFound(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "namespaceNotFound",
		action:    "error",
		message:   "namespace not found",
		log:       "namespace not found",
		severity:  actionlog.Warning,
		props: func() *recordActionProps {
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

// RecordErrModuleNotFoundModule returns "compose:record.moduleNotFoundModule" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RecordErrModuleNotFoundModule(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "moduleNotFoundModule",
		action:    "error",
		message:   "module not found",
		log:       "module not found",
		severity:  actionlog.Warning,
		props: func() *recordActionProps {
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

// RecordErrInvalidID returns "compose:record.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RecordErrInvalidID(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *recordActionProps {
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

// RecordErrInvalidNamespaceID returns "compose:record.invalidNamespaceID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RecordErrInvalidNamespaceID(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "invalidNamespaceID",
		action:    "error",
		message:   "invalid or missing namespace ID",
		log:       "invalid or missing namespace ID",
		severity:  actionlog.Warning,
		props: func() *recordActionProps {
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

// RecordErrInvalidModuleID returns "compose:record.invalidModuleID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RecordErrInvalidModuleID(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "invalidModuleID",
		action:    "error",
		message:   "invalid or missing module ID",
		log:       "invalid or missing module ID",
		severity:  actionlog.Warning,
		props: func() *recordActionProps {
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

// RecordErrStaleData returns "compose:record.staleData" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RecordErrStaleData(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "staleData",
		action:    "error",
		message:   "stale data",
		log:       "stale data",
		severity:  actionlog.Warning,
		props: func() *recordActionProps {
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

// RecordErrNotAllowedToRead returns "compose:record.notAllowedToRead" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrNotAllowedToRead(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this record",
		log:       "failed to read {record}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrNotAllowedToReadNamespace returns "compose:record.notAllowedToReadNamespace" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrNotAllowedToReadNamespace(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notAllowedToReadNamespace",
		action:    "error",
		message:   "not allowed to read this namespace",
		log:       "failed to read namespace {namespace}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrNotAllowedToReadModule returns "compose:record.notAllowedToReadModule" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrNotAllowedToReadModule(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notAllowedToReadModule",
		action:    "error",
		message:   "not allowed to read module",
		log:       "failed to read module {module}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrNotAllowedToListRecords returns "compose:record.notAllowedToListRecords" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrNotAllowedToListRecords(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notAllowedToListRecords",
		action:    "error",
		message:   "not allowed to list records",
		log:       "failed to list record; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrNotAllowedToCreate returns "compose:record.notAllowedToCreate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrNotAllowedToCreate(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create records",
		log:       "failed to create record; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrNotAllowedToUpdate returns "compose:record.notAllowedToUpdate" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrNotAllowedToUpdate(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this record",
		log:       "failed to update {record}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrNotAllowedToDelete returns "compose:record.notAllowedToDelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrNotAllowedToDelete(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this record",
		log:       "failed to delete {record}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrNotAllowedToUndelete returns "compose:record.notAllowedToUndelete" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrNotAllowedToUndelete(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this record",
		log:       "failed to undelete {record}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrNotAllowedToChangeFieldValue returns "compose:record.notAllowedToChangeFieldValue" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrNotAllowedToChangeFieldValue(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "notAllowedToChangeFieldValue",
		action:    "error",
		message:   "not allowed to change value of field {field}",
		log:       "failed to change value of field {field}; insufficient permissions",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrImportSessionAlreadActive returns "compose:record.importSessionAlreadActive" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrImportSessionAlreadActive(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "importSessionAlreadActive",
		action:    "error",
		message:   "import session already active",
		log:       "failed to start import session",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrFieldNotFound returns "compose:record.fieldNotFound" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrFieldNotFound(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "fieldNotFound",
		action:    "error",
		message:   "no such field {field}",
		log:       "no such field {field}",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrInvalidValueStructure returns "compose:record.invalidValueStructure" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrInvalidValueStructure(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "invalidValueStructure",
		action:    "error",
		message:   "more than one value for a single-value field {field}",
		log:       "more than one value for a single-value field {field}",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrUnknownBulkOperation returns "compose:record.unknownBulkOperation" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrUnknownBulkOperation(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "unknownBulkOperation",
		action:    "error",
		message:   "unknown bulk operation {bulkOperation}",
		log:       "unknown bulk operation {bulkOperation}",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrInvalidReferenceFormat returns "compose:record.invalidReferenceFormat" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrInvalidReferenceFormat(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "invalidReferenceFormat",
		action:    "error",
		message:   "invalid reference format",
		log:       "invalid reference format",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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

// RecordErrValueInput returns "compose:record.valueInput" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RecordErrValueInput(props ...*recordActionProps) *recordError {
	var e = &recordError{
		timestamp: time.Now(),
		resource:  "compose:record",
		error:     "valueInput",
		action:    "error",
		message:   "invalid record value input: {err}",
		log:       "invalid record value input: {err}",
		severity:  actionlog.Error,
		props: func() *recordActionProps {
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
// action (optional) fn will be used to construct recordAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc record) recordAction(ctx context.Context, props *recordActionProps, action func(...*recordActionProps) *recordAction, err error) error {
	var (
		ok bool

		// Return error
		retError *recordError

		// Recorder error
		recError *recordError
	)

	if err != nil {
		if retError, ok = err.(*recordError); !ok {
			// got non-record error, wrap it with RecordErrGeneric
			retError = RecordErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use RecordErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type recordError
				if unwrappedSinkError, ok := unwrappedError.(*recordError); ok {
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
