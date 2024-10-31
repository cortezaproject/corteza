package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/service/record_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"strings"
	"time"
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
		positionField *types.ModuleField
		groupField    *types.ModuleField
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

	recordLogMetaKey   struct{}
	recordPropsMetaKey struct{}
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
// This function is auto-generated.
func (p *recordActionProps) setRecord(record *types.Record) *recordActionProps {
	p.record = record
	return p
}

// setChanged updates recordActionProps's changed
//
// This function is auto-generated.
func (p *recordActionProps) setChanged(changed *types.Record) *recordActionProps {
	p.changed = changed
	return p
}

// setFilter updates recordActionProps's filter
//
// This function is auto-generated.
func (p *recordActionProps) setFilter(filter *types.RecordFilter) *recordActionProps {
	p.filter = filter
	return p
}

// setNamespace updates recordActionProps's namespace
//
// This function is auto-generated.
func (p *recordActionProps) setNamespace(namespace *types.Namespace) *recordActionProps {
	p.namespace = namespace
	return p
}

// setModule updates recordActionProps's module
//
// This function is auto-generated.
func (p *recordActionProps) setModule(module *types.Module) *recordActionProps {
	p.module = module
	return p
}

// setBulkOperation updates recordActionProps's bulkOperation
//
// This function is auto-generated.
func (p *recordActionProps) setBulkOperation(bulkOperation string) *recordActionProps {
	p.bulkOperation = bulkOperation
	return p
}

// setField updates recordActionProps's field
//
// This function is auto-generated.
func (p *recordActionProps) setField(field string) *recordActionProps {
	p.field = field
	return p
}

// setPositionField updates recordActionProps's positionField
//
// This function is auto-generated.
func (p *recordActionProps) setPositionField(positionField *types.ModuleField) *recordActionProps {
	p.positionField = positionField
	return p
}

// setGroupField updates recordActionProps's groupField
//
// This function is auto-generated.
func (p *recordActionProps) setGroupField(groupField *types.ModuleField) *recordActionProps {
	p.groupField = groupField
	return p
}

// setValue updates recordActionProps's value
//
// This function is auto-generated.
func (p *recordActionProps) setValue(value string) *recordActionProps {
	p.value = value
	return p
}

// setValueErrors updates recordActionProps's valueErrors
//
// This function is auto-generated.
func (p *recordActionProps) setValueErrors(valueErrors *types.RecordValueErrorSet) *recordActionProps {
	p.valueErrors = valueErrors
	return p
}

// Serialize converts recordActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p recordActionProps) Serialize() actionlog.Meta {
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
	if p.positionField != nil {
		m.Set("positionField.name", p.positionField.Name, true)
		m.Set("positionField.label", p.positionField.Label, true)
	}
	if p.groupField != nil {
		m.Set("groupField.name", p.groupField.Name, true)
		m.Set("groupField.label", p.groupField.Label, true)
	}
	m.Set("value", p.value, true)
	if p.valueErrors != nil {
		m.Set("valueErrors.set", p.valueErrors.Set, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p recordActionProps) Format(in string, err error) string {
	var (
		pairs = []string{"{{err}}"}
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
		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}

	if p.record != nil {
		// replacement for "{{record}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{record}}",
			fns(
				p.record.ID,
				p.record.ModuleID,
				p.record.NamespaceID,
				p.record.OwnedBy,
			),
		)
		pairs = append(pairs, "{{record.ID}}", fns(p.record.ID))
		pairs = append(pairs, "{{record.moduleID}}", fns(p.record.ModuleID))
		pairs = append(pairs, "{{record.namespaceID}}", fns(p.record.NamespaceID))
		pairs = append(pairs, "{{record.ownedBy}}", fns(p.record.OwnedBy))
	}

	if p.changed != nil {
		// replacement for "{{changed}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{changed}}",
			fns(
				p.changed.ID,
				p.changed.ModuleID,
				p.changed.NamespaceID,
				p.changed.OwnedBy,
			),
		)
		pairs = append(pairs, "{{changed.ID}}", fns(p.changed.ID))
		pairs = append(pairs, "{{changed.moduleID}}", fns(p.changed.ModuleID))
		pairs = append(pairs, "{{changed.namespaceID}}", fns(p.changed.NamespaceID))
		pairs = append(pairs, "{{changed.ownedBy}}", fns(p.changed.OwnedBy))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Query,
				p.filter.NamespaceID,
				p.filter.ModuleID,
				p.filter.Deleted,
				p.filter.Sort,
				p.filter.Limit,
			),
		)
		pairs = append(pairs, "{{filter.query}}", fns(p.filter.Query))
		pairs = append(pairs, "{{filter.namespaceID}}", fns(p.filter.NamespaceID))
		pairs = append(pairs, "{{filter.moduleID}}", fns(p.filter.ModuleID))
		pairs = append(pairs, "{{filter.deleted}}", fns(p.filter.Deleted))
		pairs = append(pairs, "{{filter.sort}}", fns(p.filter.Sort))
		pairs = append(pairs, "{{filter.limit}}", fns(p.filter.Limit))
	}

	if p.namespace != nil {
		// replacement for "{{namespace}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{namespace}}",
			fns(
				p.namespace.Name,
				p.namespace.Slug,
				p.namespace.ID,
			),
		)
		pairs = append(pairs, "{{namespace.name}}", fns(p.namespace.Name))
		pairs = append(pairs, "{{namespace.slug}}", fns(p.namespace.Slug))
		pairs = append(pairs, "{{namespace.ID}}", fns(p.namespace.ID))
	}

	if p.module != nil {
		// replacement for "{{module}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{module}}",
			fns(
				p.module.Name,
				p.module.Handle,
				p.module.ID,
				p.module.NamespaceID,
			),
		)
		pairs = append(pairs, "{{module.name}}", fns(p.module.Name))
		pairs = append(pairs, "{{module.handle}}", fns(p.module.Handle))
		pairs = append(pairs, "{{module.ID}}", fns(p.module.ID))
		pairs = append(pairs, "{{module.namespaceID}}", fns(p.module.NamespaceID))
	}
	pairs = append(pairs, "{{bulkOperation}}", fns(p.bulkOperation))
	pairs = append(pairs, "{{field}}", fns(p.field))

	if p.positionField != nil {
		// replacement for "{{positionField}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{positionField}}",
			fns(
				p.positionField.Name,
				p.positionField.Label,
			),
		)
		pairs = append(pairs, "{{positionField.name}}", fns(p.positionField.Name))
		pairs = append(pairs, "{{positionField.label}}", fns(p.positionField.Label))
	}

	if p.groupField != nil {
		// replacement for "{{groupField}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{groupField}}",
			fns(
				p.groupField.Name,
				p.groupField.Label,
			),
		)
		pairs = append(pairs, "{{groupField.name}}", fns(p.groupField.Name))
		pairs = append(pairs, "{{groupField.label}}", fns(p.groupField.Label))
	}
	pairs = append(pairs, "{{value}}", fns(p.value))

	if p.valueErrors != nil {
		// replacement for "{{valueErrors}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{valueErrors}}",
			fns(
				p.valueErrors.Set,
			),
		)
		pairs = append(pairs, "{{valueErrors.set}}", fns(p.valueErrors.Set))
	}
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
func (a *recordAction) String() string {
	var props = &recordActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *recordAction) ToAction() *actionlog.Action {
	return &actionlog.Action{
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.Serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// RecordActionSearch returns "compose:record.search" action
//
// This function is auto-generated.
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

// RecordActionSearchSensitive returns "compose:record.searchSensitive" action
//
// This function is auto-generated.
func RecordActionSearchSensitive(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "searchSensitive",
		log:       "searched for records with sensitive data",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionLookup returns "compose:record.lookup" action
//
// This function is auto-generated.
func RecordActionLookup(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "lookup",
		log:       "looked-up for a {{record}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionReport returns "compose:record.report" action
//
// This function is auto-generated.
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

// RecordActionBulk returns "compose:record.bulk" action
//
// This function is auto-generated.
func RecordActionBulk(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "bulk",
		log:       "bulk record operation",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionCreate returns "compose:record.create" action
//
// This function is auto-generated.
func RecordActionCreate(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "create",
		log:       "created {{record}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionUpdate returns "compose:record.update" action
//
// This function is auto-generated.
func RecordActionUpdate(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "update",
		log:       "updated {{record}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionDelete returns "compose:record.delete" action
//
// This function is auto-generated.
func RecordActionDelete(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "delete",
		log:       "deleted {{record}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionPatch returns "compose:record.patch" action
//
// This function is auto-generated.
func RecordActionPatch(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "patch",
		log:       "patched {{record}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionUndelete returns "compose:record.undelete" action
//
// This function is auto-generated.
func RecordActionUndelete(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "undelete",
		log:       "undeleted {{record}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionImport returns "compose:record.import" action
//
// This function is auto-generated.
func RecordActionImport(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "import",
		log:       "records imported",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionSearchRevisions returns "compose:record.searchRevisions" action
//
// This function is auto-generated.
func RecordActionSearchRevisions(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "searchRevisions",
		log:       "record revisions searched",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionExport returns "compose:record.export" action
//
// This function is auto-generated.
func RecordActionExport(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "export",
		log:       "records exported",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionOrganize returns "compose:record.organize" action
//
// This function is auto-generated.
func RecordActionOrganize(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "organize",
		log:       "records organized",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorInvoked returns "compose:record.iteratorInvoked" action
//
// This function is auto-generated.
func RecordActionIteratorInvoked(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorInvoked",
		log:       "iterator invoked",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorIteration returns "compose:record.iteratorIteration" action
//
// This function is auto-generated.
func RecordActionIteratorIteration(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorIteration",
		log:       "processed record iteration",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorClone returns "compose:record.iteratorClone" action
//
// This function is auto-generated.
func RecordActionIteratorClone(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorClone",
		log:       "cloned record in iteration",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorUpdate returns "compose:record.iteratorUpdate" action
//
// This function is auto-generated.
func RecordActionIteratorUpdate(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorUpdate",
		log:       "updated record in iteration",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorDelete returns "compose:record.iteratorDelete" action
//
// This function is auto-generated.
func RecordActionIteratorDelete(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorDelete",
		log:       "deleted record in iteration",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RecordActionIteratorUndelete returns "compose:record.iteratorUndelete" action
//
// This function is auto-generated.
func RecordActionIteratorUndelete(props ...*recordActionProps) *recordAction {
	a := &recordAction{
		timestamp: time.Now(),
		resource:  "compose:record",
		action:    "iteratorUndelete",
		log:       "undeleted record in iteration",
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

// RecordErrGeneric returns "compose:record.generic" as *errors.Error
//
// This function is auto-generated.
func RecordErrGeneric(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "{err}"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotFound returns "compose:record.notFound" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotFound(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("record not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNamespaceNotFound returns "compose:record.namespaceNotFound" as *errors.Error
//
// This function is auto-generated.
func RecordErrNamespaceNotFound(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("namespace not found", nil),

		errors.Meta("type", "namespaceNotFound"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.namespaceNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrModuleNotFoundModule returns "compose:record.moduleNotFoundModule" as *errors.Error
//
// This function is auto-generated.
func RecordErrModuleNotFoundModule(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module not found", nil),

		errors.Meta("type", "moduleNotFoundModule"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.moduleNotFoundModule"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrInvalidID returns "compose:record.invalidID" as *errors.Error
//
// This function is auto-generated.
func RecordErrInvalidID(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrInvalidNamespaceID returns "compose:record.invalidNamespaceID" as *errors.Error
//
// This function is auto-generated.
func RecordErrInvalidNamespaceID(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid or missing namespace ID", nil),

		errors.Meta("type", "invalidNamespaceID"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.invalidNamespaceID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrInvalidModuleID returns "compose:record.invalidModuleID" as *errors.Error
//
// This function is auto-generated.
func RecordErrInvalidModuleID(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid or missing module ID", nil),

		errors.Meta("type", "invalidModuleID"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.invalidModuleID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrStaleData returns "compose:record.staleData" as *errors.Error
//
// This function is auto-generated.
func RecordErrStaleData(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToRead returns "compose:record.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToRead(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this record", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to read {{record}}; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToSearch returns "compose:record.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToSearch(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list records", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to search or list records; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToSearchRevisions returns "compose:record.notAllowedToSearchRevisions" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToSearchRevisions(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list record revisions", nil),

		errors.Meta("type", "notAllowedToSearchRevisions"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to search or list record revisions; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToSearchRevisions"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrRevisionsDisabledOnModule returns "compose:record.revisionsDisabledOnModule" as *errors.Error
//
// This function is auto-generated.
func RecordErrRevisionsDisabledOnModule(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("revisions are disabled on module", nil),

		errors.Meta("type", "revisionsDisabledOnModule"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to search or list record revisions; disabled on module"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.revisionsDisabledOnModule"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToReadNamespace returns "compose:record.notAllowedToReadNamespace" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToReadNamespace(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this namespace", nil),

		errors.Meta("type", "notAllowedToReadNamespace"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to read namespace {{namespace}}; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToReadNamespace"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToReadModule returns "compose:record.notAllowedToReadModule" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToReadModule(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read module", nil),

		errors.Meta("type", "notAllowedToReadModule"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to read module {{module}}; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToReadModule"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToListRecords returns "compose:record.notAllowedToListRecords" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToListRecords(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list records", nil),

		errors.Meta("type", "notAllowedToListRecords"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to list record; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToListRecords"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToCreate returns "compose:record.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToCreate(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create records", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to create record; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToUpdate returns "compose:record.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToUpdate(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this record", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to update {{record}}; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToDelete returns "compose:record.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToDelete(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this record", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to delete {{record}}; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToUndelete returns "compose:record.notAllowedToUndelete" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToUndelete(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this record", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to undelete {{record}}; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrNotAllowedToChangeFieldValue returns "compose:record.notAllowedToChangeFieldValue" as *errors.Error
//
// This function is auto-generated.
func RecordErrNotAllowedToChangeFieldValue(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to change value of field {{field}}", nil),

		errors.Meta("type", "notAllowedToChangeFieldValue"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to change value of field {{field}}; insufficient permissions"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.notAllowedToChangeFieldValue"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrMaxRecordsReached returns "compose:record.maxRecordsReached" as *errors.Error
//
// This function is auto-generated.
func RecordErrMaxRecordsReached(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("maximum number of records per namespace reached", nil),

		errors.Meta("type", "maxRecordsReached"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "maximum number of records per namespace reached"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.maxRecordsReached"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrMissingPositionField returns "compose:record.missingPositionField" as *errors.Error
//
// This function is auto-generated.
func RecordErrMissingPositionField(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("position module field not found", nil),

		errors.Meta("type", "missingPositionField"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "position module field not found"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.missingPositionField"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrInvalidPositionFieldKind returns "compose:record.invalidPositionFieldKind" as *errors.Error
//
// This function is auto-generated.
func RecordErrInvalidPositionFieldKind(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid position field {{positionField}} kind; kind must be 'Number'", nil),

		errors.Meta("type", "invalidPositionFieldKind"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "invalid position field {{positionField}} kind; kind must be 'Number'"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.invalidPositionFieldKind"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrInvalidPositionFieldConfigMultiValue returns "compose:record.invalidPositionFieldConfigMultiValue" as *errors.Error
//
// This function is auto-generated.
func RecordErrInvalidPositionFieldConfigMultiValue(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid position field {{positionField}} configuration; field must not be multi-value", nil),

		errors.Meta("type", "invalidPositionFieldConfigMultiValue"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "invalid position field {{positionField}} configuration; field must not be multi-value"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.invalidPositionFieldConfigMultiValue"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrInvalidPositionValueType returns "compose:record.invalidPositionValueType" as *errors.Error
//
// This function is auto-generated.
func RecordErrInvalidPositionValueType(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid position value data type; value must be numeric", nil),

		errors.Meta("type", "invalidPositionValueType"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "invalid position value data type; value must be numeric"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.invalidPositionValueType"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrMissingGroupField returns "compose:record.missingGroupField" as *errors.Error
//
// This function is auto-generated.
func RecordErrMissingGroupField(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("group module field not found", nil),

		errors.Meta("type", "missingGroupField"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "group module field not found"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.missingGroupField"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrInvalidGroupFieldConfigMultiValue returns "compose:record.invalidGroupFieldConfigMultiValue" as *errors.Error
//
// This function is auto-generated.
func RecordErrInvalidGroupFieldConfigMultiValue(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid group field {{groupField}} configuration; field must not be multi-value", nil),

		errors.Meta("type", "invalidGroupFieldConfigMultiValue"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "invalid group field {{groupField}} configuration; field must not be multi-value"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.invalidGroupFieldConfigMultiValue"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrImportSessionAlreadActive returns "compose:record.importSessionAlreadActive" as *errors.Error
//
// This function is auto-generated.
func RecordErrImportSessionAlreadActive(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("import session already active", nil),

		errors.Meta("type", "importSessionAlreadActive"),
		errors.Meta("resource", "compose:record"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(recordLogMetaKey{}, "failed to start import session"),
		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.importSessionAlreadActive"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrFieldNotFound returns "compose:record.fieldNotFound" as *errors.Error
//
// This function is auto-generated.
func RecordErrFieldNotFound(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("no such field {{field}}", nil),

		errors.Meta("type", "fieldNotFound"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.fieldNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrInvalidValueStructure returns "compose:record.invalidValueStructure" as *errors.Error
//
// This function is auto-generated.
func RecordErrInvalidValueStructure(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("more than one value for a single-value field {{field}}", nil),

		errors.Meta("type", "invalidValueStructure"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.invalidValueStructure"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrUnknownBulkOperation returns "compose:record.unknownBulkOperation" as *errors.Error
//
// This function is auto-generated.
func RecordErrUnknownBulkOperation(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("unknown bulk operation {{bulkOperation}}", nil),

		errors.Meta("type", "unknownBulkOperation"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.unknownBulkOperation"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrInvalidReferenceFormat returns "compose:record.invalidReferenceFormat" as *errors.Error
//
// This function is auto-generated.
func RecordErrInvalidReferenceFormat(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid reference format", nil),

		errors.Meta("type", "invalidReferenceFormat"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.invalidReferenceFormat"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RecordErrValueInput returns "compose:record.valueInput" as *errors.Error
//
// This function is auto-generated.
func RecordErrValueInput(mm ...*recordActionProps) *errors.Error {
	var p = &recordActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid record value input", nil),

		errors.Meta("type", "valueInput"),
		errors.Meta("resource", "compose:record"),

		errors.Meta(recordPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "record.errors.valueInput"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// It will wrap unrecognized/internal errors with generic errors.
//
// This function is auto-generated.
func (svc record) recordAction(ctx context.Context, props *recordActionProps, actionFn func(...*recordActionProps) *recordAction, err error) error {
	if svc.actionlog == nil || actionFn == nil {
		// action log disabled or no action fn passed, return error as-is
		return err
	} else if err == nil {
		// action completed w/o error, record it
		svc.actionlog.Record(ctx, actionFn(props).ToAction())
		return nil
	}

	a := actionFn(props).ToAction()

	// Extracting error information and recording it as action
	a.Error = err.Error()

	switch c := err.(type) {
	case *errors.Error:
		m := c.Meta()

		a.Error = err.Error()
		a.Severity = actionlog.Severity(m.AsInt("severity"))
		a.Description = props.Format(m.AsString(recordLogMetaKey{}), err)

		if p, has := m[recordPropsMetaKey{}]; has {
			a.Meta = p.(*recordActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
