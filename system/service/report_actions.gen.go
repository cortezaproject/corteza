package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/report_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
	"time"
)

type (
	reportActionProps struct {
		report *types.Report
		new    *types.Report
		update *types.Report
		filter *types.ReportFilter
	}

	reportAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *reportActionProps
	}

	reportLogMetaKey   struct{}
	reportPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setReport updates reportActionProps's report
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *reportActionProps) setReport(report *types.Report) *reportActionProps {
	p.report = report
	return p
}

// setNew updates reportActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *reportActionProps) setNew(new *types.Report) *reportActionProps {
	p.new = new
	return p
}

// setUpdate updates reportActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *reportActionProps) setUpdate(update *types.Report) *reportActionProps {
	p.update = update
	return p
}

// setFilter updates reportActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *reportActionProps) setFilter(filter *types.ReportFilter) *reportActionProps {
	p.filter = filter
	return p
}

// Serialize converts reportActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p reportActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.report != nil {
		m.Set("report.handle", p.report.Handle, true)
		m.Set("report.ID", p.report.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.handle", p.filter.Handle, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p reportActionProps) Format(in string, err error) string {
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
		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}

	if p.report != nil {
		// replacement for "{report}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{report}",
			fns(
				p.report.Handle,
				p.report.ID,
			),
		)
		pairs = append(pairs, "{report.handle}", fns(p.report.Handle))
		pairs = append(pairs, "{report.ID}", fns(p.report.ID))
	}

	if p.new != nil {
		// replacement for "{new}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{new}",
			fns(
				p.new.Handle,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{new.handle}", fns(p.new.Handle))
		pairs = append(pairs, "{new.ID}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{update}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{update}",
			fns(
				p.update.Handle,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{update.handle}", fns(p.update.Handle))
		pairs = append(pairs, "{update.ID}", fns(p.update.ID))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Handle,
				p.filter.Deleted,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{filter.handle}", fns(p.filter.Handle))
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
func (a *reportAction) String() string {
	var props = &reportActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *reportAction) ToAction() *actionlog.Action {
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

// ReportActionSearch returns "system:report.search" action
//
// This function is auto-generated.
//
func ReportActionSearch(props ...*reportActionProps) *reportAction {
	a := &reportAction{
		timestamp: time.Now(),
		resource:  "system:report",
		action:    "search",
		log:       "searched for reports",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReportActionLookup returns "system:report.lookup" action
//
// This function is auto-generated.
//
func ReportActionLookup(props ...*reportActionProps) *reportAction {
	a := &reportAction{
		timestamp: time.Now(),
		resource:  "system:report",
		action:    "lookup",
		log:       "looked-up for a {report}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReportActionCreate returns "system:report.create" action
//
// This function is auto-generated.
//
func ReportActionCreate(props ...*reportActionProps) *reportAction {
	a := &reportAction{
		timestamp: time.Now(),
		resource:  "system:report",
		action:    "create",
		log:       "created {report}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReportActionUpdate returns "system:report.update" action
//
// This function is auto-generated.
//
func ReportActionUpdate(props ...*reportActionProps) *reportAction {
	a := &reportAction{
		timestamp: time.Now(),
		resource:  "system:report",
		action:    "update",
		log:       "updated {report}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReportActionDelete returns "system:report.delete" action
//
// This function is auto-generated.
//
func ReportActionDelete(props ...*reportActionProps) *reportAction {
	a := &reportAction{
		timestamp: time.Now(),
		resource:  "system:report",
		action:    "delete",
		log:       "deleted {report}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReportActionUndelete returns "system:report.undelete" action
//
// This function is auto-generated.
//
func ReportActionUndelete(props ...*reportActionProps) *reportAction {
	a := &reportAction{
		timestamp: time.Now(),
		resource:  "system:report",
		action:    "undelete",
		log:       "undeleted {report}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ReportActionRun returns "system:report.run" action
//
// This function is auto-generated.
//
func ReportActionRun(props ...*reportActionProps) *reportAction {
	a := &reportAction{
		timestamp: time.Now(),
		resource:  "system:report",
		action:    "run",
		log:       "report ran",
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

// ReportErrGeneric returns "system:report.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrGeneric(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:report"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(reportLogMetaKey{}, "{err}"),
		errors.Meta(reportPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReportErrNotFound returns "system:report.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrNotFound(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("report not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:report"),

		errors.Meta(reportPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReportErrInvalidID returns "system:report.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrInvalidID(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:report"),

		errors.Meta(reportPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReportErrNotAllowedToRead returns "system:report.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrNotAllowedToRead(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this report", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:report"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(reportLogMetaKey{}, "failed to read {report}; insufficient permissions"),
		errors.Meta(reportPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReportErrNotAllowedToListReports returns "system:report.notAllowedToListReports" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrNotAllowedToListReports(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list reports", nil),

		errors.Meta("type", "notAllowedToListReports"),
		errors.Meta("resource", "system:report"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(reportLogMetaKey{}, "failed to list report; insufficient permissions"),
		errors.Meta(reportPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReportErrNotAllowedToCreate returns "system:report.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrNotAllowedToCreate(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create reports", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:report"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(reportLogMetaKey{}, "failed to create report; insufficient permissions"),
		errors.Meta(reportPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReportErrNotAllowedToUpdate returns "system:report.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrNotAllowedToUpdate(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this report", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:report"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(reportLogMetaKey{}, "failed to update {report}; insufficient permissions"),
		errors.Meta(reportPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReportErrNotAllowedToDelete returns "system:report.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrNotAllowedToDelete(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this report", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:report"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(reportLogMetaKey{}, "failed to delete {report}; insufficient permissions"),
		errors.Meta(reportPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReportErrNotAllowedToUndelete returns "system:report.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrNotAllowedToUndelete(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this report", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:report"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(reportLogMetaKey{}, "failed to undelete {report}; insufficient permissions"),
		errors.Meta(reportPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ReportErrNotAllowedToRun returns "system:report.notAllowedToRun" as *errors.Error
//
//
// This function is auto-generated.
//
func ReportErrNotAllowedToRun(mm ...*reportActionProps) *errors.Error {
	var p = &reportActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to run this report", nil),

		errors.Meta("type", "notAllowedToRun"),
		errors.Meta("resource", "system:report"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(reportLogMetaKey{}, "failed to run {report}; insufficient permissions"),
		errors.Meta(reportPropsMetaKey{}, p),

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
//
func (svc report) recordAction(ctx context.Context, props *reportActionProps, actionFn func(...*reportActionProps) *reportAction, err error) error {
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
		a.Description = props.Format(m.AsString(reportLogMetaKey{}), err)

		if p, has := m[reportPropsMetaKey{}]; has {
			a.Meta = p.(*reportActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
