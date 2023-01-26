package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/template_actions.yaml

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
	templateActionProps struct {
		template *types.Template
		new      *types.Template
		update   *types.Template
		filter   *types.TemplateFilter
	}

	templateAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *templateActionProps
	}

	templateLogMetaKey   struct{}
	templatePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setTemplate updates templateActionProps's template
//
// This function is auto-generated.
//
func (p *templateActionProps) setTemplate(template *types.Template) *templateActionProps {
	p.template = template
	return p
}

// setNew updates templateActionProps's new
//
// This function is auto-generated.
//
func (p *templateActionProps) setNew(new *types.Template) *templateActionProps {
	p.new = new
	return p
}

// setUpdate updates templateActionProps's update
//
// This function is auto-generated.
//
func (p *templateActionProps) setUpdate(update *types.Template) *templateActionProps {
	p.update = update
	return p
}

// setFilter updates templateActionProps's filter
//
// This function is auto-generated.
//
func (p *templateActionProps) setFilter(filter *types.TemplateFilter) *templateActionProps {
	p.filter = filter
	return p
}

// Serialize converts templateActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p templateActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.template != nil {
		m.Set("template.handle", p.template.Handle, true)
		m.Set("template.type", p.template.Type, true)
		m.Set("template.ID", p.template.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.type", p.new.Type, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.type", p.update.Type, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.templateID", p.filter.TemplateID, true)
		m.Set("filter.handle", p.filter.Handle, true)
		m.Set("filter.type", p.filter.Type, true)
		m.Set("filter.ownerID", p.filter.OwnerID, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p templateActionProps) Format(in string, err error) string {
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

	if p.template != nil {
		// replacement for "{{template}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{template}}",
			fns(
				p.template.Handle,
				p.template.Type,
				p.template.ID,
			),
		)
		pairs = append(pairs, "{{template.handle}}", fns(p.template.Handle))
		pairs = append(pairs, "{{template.type}}", fns(p.template.Type))
		pairs = append(pairs, "{{template.ID}}", fns(p.template.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Handle,
				p.new.Type,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.handle}}", fns(p.new.Handle))
		pairs = append(pairs, "{{new.type}}", fns(p.new.Type))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Handle,
				p.update.Type,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.handle}}", fns(p.update.Handle))
		pairs = append(pairs, "{{update.type}}", fns(p.update.Type))
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.TemplateID,
				p.filter.Handle,
				p.filter.Type,
				p.filter.OwnerID,
				p.filter.Deleted,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{{filter.templateID}}", fns(p.filter.TemplateID))
		pairs = append(pairs, "{{filter.handle}}", fns(p.filter.Handle))
		pairs = append(pairs, "{{filter.type}}", fns(p.filter.Type))
		pairs = append(pairs, "{{filter.ownerID}}", fns(p.filter.OwnerID))
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
func (a *templateAction) String() string {
	var props = &templateActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *templateAction) ToAction() *actionlog.Action {
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

// TemplateActionSearch returns "system:template.search" action
//
// This function is auto-generated.
//
func TemplateActionSearch(props ...*templateActionProps) *templateAction {
	a := &templateAction{
		timestamp: time.Now(),
		resource:  "system:template",
		action:    "search",
		log:       "searched for templates",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TemplateActionLookup returns "system:template.lookup" action
//
// This function is auto-generated.
//
func TemplateActionLookup(props ...*templateActionProps) *templateAction {
	a := &templateAction{
		timestamp: time.Now(),
		resource:  "system:template",
		action:    "lookup",
		log:       "looked-up for a {{template}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TemplateActionCreate returns "system:template.create" action
//
// This function is auto-generated.
//
func TemplateActionCreate(props ...*templateActionProps) *templateAction {
	a := &templateAction{
		timestamp: time.Now(),
		resource:  "system:template",
		action:    "create",
		log:       "created {{template}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TemplateActionUpdate returns "system:template.update" action
//
// This function is auto-generated.
//
func TemplateActionUpdate(props ...*templateActionProps) *templateAction {
	a := &templateAction{
		timestamp: time.Now(),
		resource:  "system:template",
		action:    "update",
		log:       "updated {{template}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TemplateActionDelete returns "system:template.delete" action
//
// This function is auto-generated.
//
func TemplateActionDelete(props ...*templateActionProps) *templateAction {
	a := &templateAction{
		timestamp: time.Now(),
		resource:  "system:template",
		action:    "delete",
		log:       "deleted {{template}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TemplateActionUndelete returns "system:template.undelete" action
//
// This function is auto-generated.
//
func TemplateActionUndelete(props ...*templateActionProps) *templateAction {
	a := &templateAction{
		timestamp: time.Now(),
		resource:  "system:template",
		action:    "undelete",
		log:       "undeleted {{template}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TemplateActionRender returns "system:template.render" action
//
// This function is auto-generated.
//
func TemplateActionRender(props ...*templateActionProps) *templateAction {
	a := &templateAction{
		timestamp: time.Now(),
		resource:  "system:template",
		action:    "render",
		log:       "rendered {{template}}",
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

// TemplateErrGeneric returns "system:template.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrGeneric(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:template"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(templateLogMetaKey{}, "{err}"),
		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrNotFound returns "system:template.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrNotFound(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("template not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:template"),

		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrInvalidID returns "system:template.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrInvalidID(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:template"),

		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrInvalidHandle returns "system:template.invalidHandle" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrInvalidHandle(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:template"),

		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrMissingShort returns "system:template.missingShort" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrMissingShort(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("missing short name", nil),

		errors.Meta("type", "missingShort"),
		errors.Meta("resource", "system:template"),

		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.missingShort"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrCannotRenderPartial returns "system:template.cannotRenderPartial" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrCannotRenderPartial(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("cannot render partial templates", nil),

		errors.Meta("type", "cannotRenderPartial"),
		errors.Meta("resource", "system:template"),

		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.cannotRenderPartial"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrNotAllowedToRead returns "system:template.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrNotAllowedToRead(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this template", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:template"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(templateLogMetaKey{}, "failed to read {{template.handle}}; insufficient permissions"),
		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrNotAllowedToSearch returns "system:template.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrNotAllowedToSearch(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list templates", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:template"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(templateLogMetaKey{}, "failed to search or list templates; insufficient permissions"),
		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrNotAllowedToCreate returns "system:template.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrNotAllowedToCreate(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create templates", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:template"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(templateLogMetaKey{}, "failed to create template; insufficient permissions"),
		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrNotAllowedToUpdate returns "system:template.notAllowedToUpdate" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrNotAllowedToUpdate(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this template", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:template"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(templateLogMetaKey{}, "failed to update {{template.handle}}; insufficient permissions"),
		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrNotAllowedToDelete returns "system:template.notAllowedToDelete" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrNotAllowedToDelete(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this template", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:template"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(templateLogMetaKey{}, "failed to delete {{template.handle}}; insufficient permissions"),
		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrNotAllowedToUndelete returns "system:template.notAllowedToUndelete" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrNotAllowedToUndelete(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this template", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:template"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(templateLogMetaKey{}, "failed to undelete {{template.handle}}; insufficient permissions"),
		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TemplateErrNotAllowedToRender returns "system:template.notAllowedToRender" as *errors.Error
//
//
// This function is auto-generated.
//
func TemplateErrNotAllowedToRender(mm ...*templateActionProps) *errors.Error {
	var p = &templateActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to render this template", nil),

		errors.Meta("type", "notAllowedToRender"),
		errors.Meta("resource", "system:template"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(templateLogMetaKey{}, "failed to render {{template.handle}}; insufficient permissions"),
		errors.Meta(templatePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "template.errors.notAllowedToRender"),

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
func (svc template) recordAction(ctx context.Context, props *templateActionProps, actionFn func(...*templateActionProps) *templateAction, err error) error {
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
		a.Description = props.Format(m.AsString(templateLogMetaKey{}), err)

		if p, has := m[templatePropsMetaKey{}]; has {
			a.Meta = p.(*templateActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
