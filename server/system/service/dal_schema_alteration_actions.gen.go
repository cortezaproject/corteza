package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/dal_schema_alteration_actions.yaml

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
	dalSchemaAlterationActionProps struct {
		dalSchemaAlteration *types.DalSchemaAlteration
		new                 *types.DalSchemaAlteration
		update              *types.DalSchemaAlteration
		existing            *types.DalSchemaAlteration
		filter              *types.DalSchemaAlterationFilter
	}

	dalSchemaAlterationAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *dalSchemaAlterationActionProps
	}

	dalSchemaAlterationLogMetaKey   struct{}
	dalSchemaAlterationPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setDalSchemaAlteration updates dalSchemaAlterationActionProps's dalSchemaAlteration
//
// This function is auto-generated.
//
func (p *dalSchemaAlterationActionProps) setDalSchemaAlteration(dalSchemaAlteration *types.DalSchemaAlteration) *dalSchemaAlterationActionProps {
	p.dalSchemaAlteration = dalSchemaAlteration
	return p
}

// setNew updates dalSchemaAlterationActionProps's new
//
// This function is auto-generated.
//
func (p *dalSchemaAlterationActionProps) setNew(new *types.DalSchemaAlteration) *dalSchemaAlterationActionProps {
	p.new = new
	return p
}

// setUpdate updates dalSchemaAlterationActionProps's update
//
// This function is auto-generated.
//
func (p *dalSchemaAlterationActionProps) setUpdate(update *types.DalSchemaAlteration) *dalSchemaAlterationActionProps {
	p.update = update
	return p
}

// setExisting updates dalSchemaAlterationActionProps's existing
//
// This function is auto-generated.
//
func (p *dalSchemaAlterationActionProps) setExisting(existing *types.DalSchemaAlteration) *dalSchemaAlterationActionProps {
	p.existing = existing
	return p
}

// setFilter updates dalSchemaAlterationActionProps's filter
//
// This function is auto-generated.
//
func (p *dalSchemaAlterationActionProps) setFilter(filter *types.DalSchemaAlterationFilter) *dalSchemaAlterationActionProps {
	p.filter = filter
	return p
}

// Serialize converts dalSchemaAlterationActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p dalSchemaAlterationActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.dalSchemaAlteration != nil {
		m.Set("dalSchemaAlteration.ID", p.dalSchemaAlteration.ID, true)
	}
	if p.new != nil {
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.ID", p.update.ID, true)
	}
	if p.existing != nil {
		m.Set("existing.ID", p.existing.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.alterationID", p.filter.AlterationID, true)
		m.Set("filter.batchID", p.filter.BatchID, true)
		m.Set("filter.kind", p.filter.Kind, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
		m.Set("filter.completed", p.filter.Completed, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p dalSchemaAlterationActionProps) Format(in string, err error) string {
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

	if p.dalSchemaAlteration != nil {
		// replacement for "{{dalSchemaAlteration}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{dalSchemaAlteration}}",
			fns(
				p.dalSchemaAlteration.ID,
			),
		)
		pairs = append(pairs, "{{dalSchemaAlteration.ID}}", fns(p.dalSchemaAlteration.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.existing != nil {
		// replacement for "{{existing}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{existing}}",
			fns(
				p.existing.ID,
			),
		)
		pairs = append(pairs, "{{existing.ID}}", fns(p.existing.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.AlterationID,
				p.filter.BatchID,
				p.filter.Kind,
				p.filter.Deleted,
				p.filter.Completed,
			),
		)
		pairs = append(pairs, "{{filter.alterationID}}", fns(p.filter.AlterationID))
		pairs = append(pairs, "{{filter.batchID}}", fns(p.filter.BatchID))
		pairs = append(pairs, "{{filter.kind}}", fns(p.filter.Kind))
		pairs = append(pairs, "{{filter.deleted}}", fns(p.filter.Deleted))
		pairs = append(pairs, "{{filter.completed}}", fns(p.filter.Completed))
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
func (a *dalSchemaAlterationAction) String() string {
	var props = &dalSchemaAlterationActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *dalSchemaAlterationAction) ToAction() *actionlog.Action {
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

// DalSchemaAlterationActionSearch returns "system:dal-schema-alteration.search" action
//
// This function is auto-generated.
//
func DalSchemaAlterationActionSearch(props ...*dalSchemaAlterationActionProps) *dalSchemaAlterationAction {
	a := &dalSchemaAlterationAction{
		timestamp: time.Now(),
		resource:  "system:dal-schema-alteration",
		action:    "search",
		log:       "searched for matching dalSchemaAlterations",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalSchemaAlterationActionLookup returns "system:dal-schema-alteration.lookup" action
//
// This function is auto-generated.
//
func DalSchemaAlterationActionLookup(props ...*dalSchemaAlterationActionProps) *dalSchemaAlterationAction {
	a := &dalSchemaAlterationAction{
		timestamp: time.Now(),
		resource:  "system:dal-schema-alteration",
		action:    "lookup",
		log:       "looked-up for a {{dalSchemaAlteration}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalSchemaAlterationActionDelete returns "system:dal-schema-alteration.delete" action
//
// This function is auto-generated.
//
func DalSchemaAlterationActionDelete(props ...*dalSchemaAlterationActionProps) *dalSchemaAlterationAction {
	a := &dalSchemaAlterationAction{
		timestamp: time.Now(),
		resource:  "system:dal-schema-alteration",
		action:    "delete",
		log:       "deleted {{dalSchemaAlteration}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalSchemaAlterationActionUndelete returns "system:dal-schema-alteration.undelete" action
//
// This function is auto-generated.
//
func DalSchemaAlterationActionUndelete(props ...*dalSchemaAlterationActionProps) *dalSchemaAlterationAction {
	a := &dalSchemaAlterationAction{
		timestamp: time.Now(),
		resource:  "system:dal-schema-alteration",
		action:    "undelete",
		log:       "undeleted {{dalSchemaAlteration}}",
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

// DalSchemaAlterationErrGeneric returns "system:dal-schema-alteration.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSchemaAlterationErrGeneric(mm ...*dalSchemaAlterationActionProps) *errors.Error {
	var p = &dalSchemaAlterationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:dal-schema-alteration"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalSchemaAlterationLogMetaKey{}, "{err}"),
		errors.Meta(dalSchemaAlterationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-schema-alteration.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSchemaAlterationErrNotFound returns "system:dal-schema-alteration.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSchemaAlterationErrNotFound(mm ...*dalSchemaAlterationActionProps) *errors.Error {
	var p = &dalSchemaAlterationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("dalSchemaAlteration not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:dal-schema-alteration"),

		errors.Meta(dalSchemaAlterationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-schema-alteration.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSchemaAlterationErrInvalidID returns "system:dal-schema-alteration.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSchemaAlterationErrInvalidID(mm ...*dalSchemaAlterationActionProps) *errors.Error {
	var p = &dalSchemaAlterationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:dal-schema-alteration"),

		errors.Meta(dalSchemaAlterationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-schema-alteration.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSchemaAlterationErrNotAllowedToRead returns "system:dal-schema-alteration.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSchemaAlterationErrNotAllowedToRead(mm ...*dalSchemaAlterationActionProps) *errors.Error {
	var p = &dalSchemaAlterationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this dal schema alteration", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:dal-schema-alteration"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalSchemaAlterationLogMetaKey{}, "failed to read {{dalSchemaAlteration.ID}}; insufficient permissions"),
		errors.Meta(dalSchemaAlterationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-schema-alteration.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSchemaAlterationErrNotAllowedToSearch returns "system:dal-schema-alteration.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSchemaAlterationErrNotAllowedToSearch(mm ...*dalSchemaAlterationActionProps) *errors.Error {
	var p = &dalSchemaAlterationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list or search dal schema alterations", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:dal-schema-alteration"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalSchemaAlterationLogMetaKey{}, "failed to search for dal schema alterations; insufficient permissions"),
		errors.Meta(dalSchemaAlterationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-schema-alteration.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSchemaAlterationErrNotAllowedToListUsers returns "system:dal-schema-alteration.notAllowedToListUsers" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSchemaAlterationErrNotAllowedToListUsers(mm ...*dalSchemaAlterationActionProps) *errors.Error {
	var p = &dalSchemaAlterationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list dal schema alterations", nil),

		errors.Meta("type", "notAllowedToListUsers"),
		errors.Meta("resource", "system:dal-schema-alteration"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalSchemaAlterationLogMetaKey{}, "failed to list dalSchemaAlteration; insufficient permissions"),
		errors.Meta(dalSchemaAlterationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-schema-alteration.errors.notAllowedToListUsers"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSchemaAlterationErrNotAllowedToDelete returns "system:dal-schema-alteration.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSchemaAlterationErrNotAllowedToDelete(mm ...*dalSchemaAlterationActionProps) *errors.Error {
	var p = &dalSchemaAlterationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this dal schema alteration", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:dal-schema-alteration"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalSchemaAlterationLogMetaKey{}, "failed to delete {{dalSchemaAlteration.ID}}; insufficient permissions"),
		errors.Meta(dalSchemaAlterationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-schema-alteration.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSchemaAlterationErrNotAllowedToUndelete returns "system:dal-schema-alteration.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSchemaAlterationErrNotAllowedToUndelete(mm ...*dalSchemaAlterationActionProps) *errors.Error {
	var p = &dalSchemaAlterationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this dal schema alteration", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:dal-schema-alteration"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalSchemaAlterationLogMetaKey{}, "failed to undelete {{dalSchemaAlteration.ID}}; insufficient permissions"),
		errors.Meta(dalSchemaAlterationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dal-schema-alteration.errors.notAllowedToUndelete"),

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
func (svc dalSchemaAlteration) recordAction(ctx context.Context, props *dalSchemaAlterationActionProps, actionFn func(...*dalSchemaAlterationActionProps) *dalSchemaAlterationAction, err error) error {
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
		a.Description = props.Format(m.AsString(dalSchemaAlterationLogMetaKey{}), err)

		if p, has := m[dalSchemaAlterationPropsMetaKey{}]; has {
			a.Meta = p.(*dalSchemaAlterationActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
