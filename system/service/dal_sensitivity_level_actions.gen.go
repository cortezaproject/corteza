package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/dal_sensitivity_level_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
	"time"
)

type (
	dalSensitivityLevelActionProps struct {
		sensitivityLevel *types.DalSensitivityLevel
		new              *types.DalSensitivityLevel
		update           *types.DalSensitivityLevel
		search           *types.DalSensitivityLevelFilter
	}

	dalSensitivityLevelAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *dalSensitivityLevelActionProps
	}

	dalSensitivityLevelLogMetaKey   struct{}
	dalSensitivityLevelPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setSensitivityLevel updates dalSensitivityLevelActionProps's sensitivityLevel
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *dalSensitivityLevelActionProps) setSensitivityLevel(sensitivityLevel *types.DalSensitivityLevel) *dalSensitivityLevelActionProps {
	p.sensitivityLevel = sensitivityLevel
	return p
}

// setNew updates dalSensitivityLevelActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *dalSensitivityLevelActionProps) setNew(new *types.DalSensitivityLevel) *dalSensitivityLevelActionProps {
	p.new = new
	return p
}

// setUpdate updates dalSensitivityLevelActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *dalSensitivityLevelActionProps) setUpdate(update *types.DalSensitivityLevel) *dalSensitivityLevelActionProps {
	p.update = update
	return p
}

// setSearch updates dalSensitivityLevelActionProps's search
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *dalSensitivityLevelActionProps) setSearch(search *types.DalSensitivityLevelFilter) *dalSensitivityLevelActionProps {
	p.search = search
	return p
}

// Serialize converts dalSensitivityLevelActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p dalSensitivityLevelActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.sensitivityLevel != nil {
		m.Set("sensitivityLevel.handle", p.sensitivityLevel.Handle, true)
		m.Set("sensitivityLevel.ID", p.sensitivityLevel.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.search != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p dalSensitivityLevelActionProps) Format(in string, err error) string {
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

	if p.sensitivityLevel != nil {
		// replacement for "{{sensitivityLevel}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{sensitivityLevel}}",
			fns(
				p.sensitivityLevel.Handle,
				p.sensitivityLevel.ID,
			),
		)
		pairs = append(pairs, "{{sensitivityLevel.handle}}", fns(p.sensitivityLevel.Handle))
		pairs = append(pairs, "{{sensitivityLevel.ID}}", fns(p.sensitivityLevel.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Handle,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.handle}}", fns(p.new.Handle))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Handle,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.handle}}", fns(p.update.Handle))
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.search != nil {
		// replacement for "{{search}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{search}}",
			fns(),
		)
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
func (a *dalSensitivityLevelAction) String() string {
	var props = &dalSensitivityLevelActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *dalSensitivityLevelAction) ToAction() *actionlog.Action {
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

// DalSensitivityLevelActionSearch returns "system:dal-sensitivity-level.search" action
//
// This function is auto-generated.
//
func DalSensitivityLevelActionSearch(props ...*dalSensitivityLevelActionProps) *dalSensitivityLevelAction {
	a := &dalSensitivityLevelAction{
		timestamp: time.Now(),
		resource:  "system:dal-sensitivity-level",
		action:    "search",
		log:       "searched for sensitivityLevel",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalSensitivityLevelActionLookup returns "system:dal-sensitivity-level.lookup" action
//
// This function is auto-generated.
//
func DalSensitivityLevelActionLookup(props ...*dalSensitivityLevelActionProps) *dalSensitivityLevelAction {
	a := &dalSensitivityLevelAction{
		timestamp: time.Now(),
		resource:  "system:dal-sensitivity-level",
		action:    "lookup",
		log:       "looked-up for a {{sensitivityLevel}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalSensitivityLevelActionCreate returns "system:dal-sensitivity-level.create" action
//
// This function is auto-generated.
//
func DalSensitivityLevelActionCreate(props ...*dalSensitivityLevelActionProps) *dalSensitivityLevelAction {
	a := &dalSensitivityLevelAction{
		timestamp: time.Now(),
		resource:  "system:dal-sensitivity-level",
		action:    "create",
		log:       "created {{sensitivityLevel}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalSensitivityLevelActionUpdate returns "system:dal-sensitivity-level.update" action
//
// This function is auto-generated.
//
func DalSensitivityLevelActionUpdate(props ...*dalSensitivityLevelActionProps) *dalSensitivityLevelAction {
	a := &dalSensitivityLevelAction{
		timestamp: time.Now(),
		resource:  "system:dal-sensitivity-level",
		action:    "update",
		log:       "updated {{sensitivityLevel}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalSensitivityLevelActionDelete returns "system:dal-sensitivity-level.delete" action
//
// This function is auto-generated.
//
func DalSensitivityLevelActionDelete(props ...*dalSensitivityLevelActionProps) *dalSensitivityLevelAction {
	a := &dalSensitivityLevelAction{
		timestamp: time.Now(),
		resource:  "system:dal-sensitivity-level",
		action:    "delete",
		log:       "deleted {{sensitivityLevel}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// DalSensitivityLevelActionUndelete returns "system:dal-sensitivity-level.undelete" action
//
// This function is auto-generated.
//
func DalSensitivityLevelActionUndelete(props ...*dalSensitivityLevelActionProps) *dalSensitivityLevelAction {
	a := &dalSensitivityLevelAction{
		timestamp: time.Now(),
		resource:  "system:dal-sensitivity-level",
		action:    "undelete",
		log:       "undeleted {{sensitivityLevel}}",
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

// DalSensitivityLevelErrGeneric returns "system:dal-sensitivity-level.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSensitivityLevelErrGeneric(mm ...*dalSensitivityLevelActionProps) *errors.Error {
	var p = &dalSensitivityLevelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:dal-sensitivity-level"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalSensitivityLevelLogMetaKey{}, "{err}"),
		errors.Meta(dalSensitivityLevelPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dalSensitivityLevel.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSensitivityLevelErrNotFound returns "system:dal-sensitivity-level.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSensitivityLevelErrNotFound(mm ...*dalSensitivityLevelActionProps) *errors.Error {
	var p = &dalSensitivityLevelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("sensitivityLevel not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:dal-sensitivity-level"),

		errors.Meta(dalSensitivityLevelPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dalSensitivityLevel.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSensitivityLevelErrInvalidID returns "system:dal-sensitivity-level.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSensitivityLevelErrInvalidID(mm ...*dalSensitivityLevelActionProps) *errors.Error {
	var p = &dalSensitivityLevelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:dal-sensitivity-level"),

		errors.Meta(dalSensitivityLevelPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dalSensitivityLevel.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSensitivityLevelErrInvalidEndpoint returns "system:dal-sensitivity-level.invalidEndpoint" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSensitivityLevelErrInvalidEndpoint(mm ...*dalSensitivityLevelActionProps) *errors.Error {
	var p = &dalSensitivityLevelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid DSN", nil),

		errors.Meta("type", "invalidEndpoint"),
		errors.Meta("resource", "system:dal-sensitivity-level"),

		errors.Meta(dalSensitivityLevelPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dalSensitivityLevel.errors.invalidEndpoint"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSensitivityLevelErrExistsEndpoint returns "system:dal-sensitivity-level.existsEndpoint" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSensitivityLevelErrExistsEndpoint(mm ...*dalSensitivityLevelActionProps) *errors.Error {
	var p = &dalSensitivityLevelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("sensitivityLevel with this DSN already exists", nil),

		errors.Meta("type", "existsEndpoint"),
		errors.Meta("resource", "system:dal-sensitivity-level"),

		errors.Meta(dalSensitivityLevelPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dalSensitivityLevel.errors.existsEndpoint"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSensitivityLevelErrAlreadyExists returns "system:dal-sensitivity-level.alreadyExists" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSensitivityLevelErrAlreadyExists(mm ...*dalSensitivityLevelActionProps) *errors.Error {
	var p = &dalSensitivityLevelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("sensitivityLevel by that DSN already exists", nil),

		errors.Meta("type", "alreadyExists"),
		errors.Meta("resource", "system:dal-sensitivity-level"),

		errors.Meta(dalSensitivityLevelPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dalSensitivityLevel.errors.alreadyExists"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// DalSensitivityLevelErrNotAllowedToManage returns "system:dal-sensitivity-level.notAllowedToManage" as *errors.Error
//
//
// This function is auto-generated.
//
func DalSensitivityLevelErrNotAllowedToManage(mm ...*dalSensitivityLevelActionProps) *errors.Error {
	var p = &dalSensitivityLevelActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to Manage a sensitivityLevel", nil),

		errors.Meta("type", "notAllowedToManage"),
		errors.Meta("resource", "system:dal-sensitivity-level"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(dalSensitivityLevelLogMetaKey{}, "failed to Manage a sensitivityLevel; insufficient permissions"),
		errors.Meta(dalSensitivityLevelPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "dalSensitivityLevel.errors.notAllowedToManage"),

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
func (svc dalSensitivityLevel) recordAction(ctx context.Context, props *dalSensitivityLevelActionProps, actionFn func(...*dalSensitivityLevelActionProps) *dalSensitivityLevelAction, err error) error {
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
		a.Description = props.Format(m.AsString(dalSensitivityLevelLogMetaKey{}), err)

		if p, has := m[dalSensitivityLevelPropsMetaKey{}]; has {
			a.Meta = p.(*dalSensitivityLevelActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
