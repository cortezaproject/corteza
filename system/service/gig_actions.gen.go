package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/gig_actions.yaml

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
	gigServiceActionProps struct {
		gig *types.Gig
		new *types.Gig
	}

	gigServiceAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *gigServiceActionProps
	}

	gigServiceLogMetaKey   struct{}
	gigServicePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setGig updates gigServiceActionProps's gig
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *gigServiceActionProps) setGig(gig *types.Gig) *gigServiceActionProps {
	p.gig = gig
	return p
}

// setNew updates gigServiceActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *gigServiceActionProps) setNew(new *types.Gig) *gigServiceActionProps {
	p.new = new
	return p
}

// Serialize converts gigServiceActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p gigServiceActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.gig != nil {
		m.Set("gig.ID", p.gig.ID, true)
	}
	if p.new != nil {
		m.Set("new.ID", p.new.ID, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p gigServiceActionProps) Format(in string, err error) string {
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

	if p.gig != nil {
		// replacement for "{{gig}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{gig}}",
			fns(
				p.gig.ID,
			),
		)
		pairs = append(pairs, "{{gig.ID}}", fns(p.gig.ID))
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
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
//
func (a *gigServiceAction) String() string {
	var props = &gigServiceActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *gigServiceAction) ToAction() *actionlog.Action {
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

// GigServiceActionLookup returns "system:gig.lookup" action
//
// This function is auto-generated.
//
func GigServiceActionLookup(props ...*gigServiceActionProps) *gigServiceAction {
	a := &gigServiceAction{
		timestamp: time.Now(),
		resource:  "system:gig",
		action:    "lookup",
		log:       "looked-up for a {{gig}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// GigServiceActionCreate returns "system:gig.create" action
//
// This function is auto-generated.
//
func GigServiceActionCreate(props ...*gigServiceActionProps) *gigServiceAction {
	a := &gigServiceAction{
		timestamp: time.Now(),
		resource:  "system:gig",
		action:    "create",
		log:       "created {{gig}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// GigServiceActionUpdate returns "system:gig.update" action
//
// This function is auto-generated.
//
func GigServiceActionUpdate(props ...*gigServiceActionProps) *gigServiceAction {
	a := &gigServiceAction{
		timestamp: time.Now(),
		resource:  "system:gig",
		action:    "update",
		log:       "updated {{gig}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// GigServiceActionDelete returns "system:gig.delete" action
//
// This function is auto-generated.
//
func GigServiceActionDelete(props ...*gigServiceActionProps) *gigServiceAction {
	a := &gigServiceAction{
		timestamp: time.Now(),
		resource:  "system:gig",
		action:    "delete",
		log:       "deleted {{gig}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// GigServiceActionUndelete returns "system:gig.undelete" action
//
// This function is auto-generated.
//
func GigServiceActionUndelete(props ...*gigServiceActionProps) *gigServiceAction {
	a := &gigServiceAction{
		timestamp: time.Now(),
		resource:  "system:gig",
		action:    "undelete",
		log:       "undeleted {{gig}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// GigServiceActionExec returns "system:gig.exec" action
//
// This function is auto-generated.
//
func GigServiceActionExec(props ...*gigServiceActionProps) *gigServiceAction {
	a := &gigServiceAction{
		timestamp: time.Now(),
		resource:  "system:gig",
		action:    "exec",
		log:       "exec {{gig}}",
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

// GigServiceErrGeneric returns "system:gig.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func GigServiceErrGeneric(mm ...*gigServiceActionProps) *errors.Error {
	var p = &gigServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:gig"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(gigServiceLogMetaKey{}, "{err}"),
		errors.Meta(gigServicePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "gigService.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// GigServiceErrNotFound returns "system:gig.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func GigServiceErrNotFound(mm ...*gigServiceActionProps) *errors.Error {
	var p = &gigServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("gig not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:gig"),

		errors.Meta(gigServicePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "gigService.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// GigServiceErrInvalidID returns "system:gig.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func GigServiceErrInvalidID(mm ...*gigServiceActionProps) *errors.Error {
	var p = &gigServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:gig"),

		errors.Meta(gigServicePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "gigService.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// GigServiceErrNotAllowedToCreate returns "system:gig.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func GigServiceErrNotAllowedToCreate(mm ...*gigServiceActionProps) *errors.Error {
	var p = &gigServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a gig", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:gig"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(gigServiceLogMetaKey{}, "failed to create a gig; insufficient permissions"),
		errors.Meta(gigServicePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "gigService.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// GigServiceErrNotAllowedToRead returns "system:gig.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func GigServiceErrNotAllowedToRead(mm ...*gigServiceActionProps) *errors.Error {
	var p = &gigServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this gig", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:gig"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(gigServiceLogMetaKey{}, "failed to read {{gig.ID}}; insufficient permissions"),
		errors.Meta(gigServicePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "gigService.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// GigServiceErrNotAllowedToUpdate returns "system:gig.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func GigServiceErrNotAllowedToUpdate(mm ...*gigServiceActionProps) *errors.Error {
	var p = &gigServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this gig", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:gig"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(gigServiceLogMetaKey{}, "failed to update {{gig.ID}}; insufficient permissions"),
		errors.Meta(gigServicePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "gigService.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// GigServiceErrNotAllowedToDelete returns "system:gig.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func GigServiceErrNotAllowedToDelete(mm ...*gigServiceActionProps) *errors.Error {
	var p = &gigServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this gig", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:gig"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(gigServiceLogMetaKey{}, "failed to delete {{gig.ID}}; insufficient permissions"),
		errors.Meta(gigServicePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "gigService.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// GigServiceErrNotAllowedToUndelete returns "system:gig.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func GigServiceErrNotAllowedToUndelete(mm ...*gigServiceActionProps) *errors.Error {
	var p = &gigServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this gig", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:gig"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(gigServiceLogMetaKey{}, "failed to undelete {{gig.ID}}; insufficient permissions"),
		errors.Meta(gigServicePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "gigService.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// GigServiceErrNotAllowedToExec returns "system:gig.notAllowedToExec" as *errors.Error
//
//
// This function is auto-generated.
//
func GigServiceErrNotAllowedToExec(mm ...*gigServiceActionProps) *errors.Error {
	var p = &gigServiceActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to execute this gig", nil),

		errors.Meta("type", "notAllowedToExec"),
		errors.Meta("resource", "system:gig"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(gigServiceLogMetaKey{}, "failed to exec {{gig.ID}}; insufficient permissions"),
		errors.Meta(gigServicePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "gigService.errors.notAllowedToExec"),

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
func (svc gigService) recordAction(ctx context.Context, props *gigServiceActionProps, actionFn func(...*gigServiceActionProps) *gigServiceAction, err error) error {
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
		a.Description = props.Format(m.AsString(gigServiceLogMetaKey{}), err)

		if p, has := m[gigServicePropsMetaKey{}]; has {
			a.Meta = p.(*gigServiceActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
