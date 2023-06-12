package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/application_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/system/types"
	"strings"
	"time"
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

	applicationLogMetaKey   struct{}
	applicationPropsMetaKey struct{}
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
// This function is auto-generated.
//
func (p *applicationActionProps) setApplication(application *types.Application) *applicationActionProps {
	p.application = application
	return p
}

// setNew updates applicationActionProps's new
//
// This function is auto-generated.
//
func (p *applicationActionProps) setNew(new *types.Application) *applicationActionProps {
	p.new = new
	return p
}

// setUpdate updates applicationActionProps's update
//
// This function is auto-generated.
//
func (p *applicationActionProps) setUpdate(update *types.Application) *applicationActionProps {
	p.update = update
	return p
}

// setFilter updates applicationActionProps's filter
//
// This function is auto-generated.
//
func (p *applicationActionProps) setFilter(filter *types.ApplicationFilter) *applicationActionProps {
	p.filter = filter
	return p
}

// Serialize converts applicationActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p applicationActionProps) Serialize() actionlog.Meta {
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
func (p applicationActionProps) Format(in string, err error) string {
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

	if p.application != nil {
		// replacement for "{{application}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{application}}",
			fns(
				p.application.Name,
				p.application.ID,
			),
		)
		pairs = append(pairs, "{{application.name}}", fns(p.application.Name))
		pairs = append(pairs, "{{application.ID}}", fns(p.application.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Name,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.name}}", fns(p.new.Name))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Name,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.name}}", fns(p.update.Name))
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Query,
				p.filter.Name,
				p.filter.Deleted,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{{filter.query}}", fns(p.filter.Query))
		pairs = append(pairs, "{{filter.name}}", fns(p.filter.Name))
		pairs = append(pairs, "{{filter.deleted}}", fns(p.filter.Deleted))
		pairs = append(pairs, "{{filter.sort}}", fns(p.filter.Sort))
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

	return props.Format(a.log, nil)
}

func (e *applicationAction) ToAction() *actionlog.Action {
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

// ApplicationActionSearch returns "system:application.search" action
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

// ApplicationActionReorder returns "system:application.reorder" action
//
// This function is auto-generated.
//
func ApplicationActionReorder(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "reorder",
		log:       "reordered applications",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionLookup returns "system:application.lookup" action
//
// This function is auto-generated.
//
func ApplicationActionLookup(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "lookup",
		log:       "looked-up for a {{application}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionCreate returns "system:application.create" action
//
// This function is auto-generated.
//
func ApplicationActionCreate(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "create",
		log:       "created {{application}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionUpdate returns "system:application.update" action
//
// This function is auto-generated.
//
func ApplicationActionUpdate(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "update",
		log:       "updated {{application}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionDelete returns "system:application.delete" action
//
// This function is auto-generated.
//
func ApplicationActionDelete(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "delete",
		log:       "deleted {{application}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionUndelete returns "system:application.undelete" action
//
// This function is auto-generated.
//
func ApplicationActionUndelete(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "undelete",
		log:       "undeleted {{application}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionFlagManage returns "system:application.flagManage" action
//
// This function is auto-generated.
//
func ApplicationActionFlagManage(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "flagManage",
		log:       "managed flags for application {{application}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ApplicationActionFlagManageGlobal returns "system:application.flagManageGlobal" action
//
// This function is auto-generated.
//
func ApplicationActionFlagManageGlobal(props ...*applicationActionProps) *applicationAction {
	a := &applicationAction{
		timestamp: time.Now(),
		resource:  "system:application",
		action:    "flagManageGlobal",
		log:       "managed global flags for application {{application}}",
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

// ApplicationErrGeneric returns "system:application.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrGeneric(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:application"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(applicationLogMetaKey{}, "{err}"),
		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrNotFound returns "system:application.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotFound(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("application not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:application"),

		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrInvalidID returns "system:application.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrInvalidID(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:application"),

		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrStaleData returns "system:application.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrStaleData(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:application"),

		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrNotAllowedToRead returns "system:application.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToRead(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this application", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:application"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(applicationLogMetaKey{}, "failed to read {{application.name}}; insufficient permissions"),
		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrNotAllowedToSearch returns "system:application.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToSearch(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list applications", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:application"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(applicationLogMetaKey{}, "failed to search or list applications; insufficient permissions"),
		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrNotAllowedToCreate returns "system:application.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToCreate(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create applications", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:application"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(applicationLogMetaKey{}, "failed to create application; insufficient permissions"),
		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrNotAllowedToUpdate returns "system:application.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToUpdate(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this application", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:application"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(applicationLogMetaKey{}, "failed to update {{application.name}}; insufficient permissions"),
		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrNotAllowedToDelete returns "system:application.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToDelete(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this application", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:application"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(applicationLogMetaKey{}, "failed to delete {{application.name}}; insufficient permissions"),
		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrNotAllowedToUndelete returns "system:application.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToUndelete(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this application", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:application"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(applicationLogMetaKey{}, "failed to undelete {{application.name}}; insufficient permissions"),
		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrNotAllowedToManageFlag returns "system:application.notAllowedToManageFlag" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToManageFlag(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage flags for applications", nil),

		errors.Meta("type", "notAllowedToManageFlag"),
		errors.Meta("resource", "system:application"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(applicationLogMetaKey{}, "failed to manage flags {{application.name}}; insufficient permissions"),
		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.notAllowedToManageFlag"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ApplicationErrNotAllowedToManageFlagGlobal returns "system:application.notAllowedToManageFlagGlobal" as *errors.Error
//
//
// This function is auto-generated.
//
func ApplicationErrNotAllowedToManageFlagGlobal(mm ...*applicationActionProps) *errors.Error {
	var p = &applicationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage global flags for applications", nil),

		errors.Meta("type", "notAllowedToManageFlagGlobal"),
		errors.Meta("resource", "system:application"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(applicationLogMetaKey{}, "failed to manage global flags {{application.name}}; insufficient permissions"),
		errors.Meta(applicationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "application.errors.notAllowedToManageFlagGlobal"),

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
func (svc application) recordAction(ctx context.Context, props *applicationActionProps, actionFn func(...*applicationActionProps) *applicationAction, err error) error {
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
		a.Description = props.Format(m.AsString(applicationLogMetaKey{}), err)

		if p, has := m[applicationPropsMetaKey{}]; has {
			a.Meta = p.(*applicationActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
