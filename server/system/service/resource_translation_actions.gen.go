package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/resource_translation_actions.yaml

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
	resourceTranslationActionProps struct {
		resourceTranslation *types.ResourceTranslation
		new                 *types.ResourceTranslation
		update              *types.ResourceTranslation
		filter              *types.ResourceTranslationFilter
	}

	resourceTranslationAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *resourceTranslationActionProps
	}

	resourceTranslationLogMetaKey   struct{}
	resourceTranslationPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setResourceTranslation updates resourceTranslationActionProps's resourceTranslation
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *resourceTranslationActionProps) setResourceTranslation(resourceTranslation *types.ResourceTranslation) *resourceTranslationActionProps {
	p.resourceTranslation = resourceTranslation
	return p
}

// setNew updates resourceTranslationActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *resourceTranslationActionProps) setNew(new *types.ResourceTranslation) *resourceTranslationActionProps {
	p.new = new
	return p
}

// setUpdate updates resourceTranslationActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *resourceTranslationActionProps) setUpdate(update *types.ResourceTranslation) *resourceTranslationActionProps {
	p.update = update
	return p
}

// setFilter updates resourceTranslationActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *resourceTranslationActionProps) setFilter(filter *types.ResourceTranslationFilter) *resourceTranslationActionProps {
	p.filter = filter
	return p
}

// Serialize converts resourceTranslationActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p resourceTranslationActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.resourceTranslation != nil {
		m.Set("resourceTranslation.lang", p.resourceTranslation.Lang, true)
		m.Set("resourceTranslation.resource", p.resourceTranslation.Resource, true)
		m.Set("resourceTranslation.K", p.resourceTranslation.K, true)
		m.Set("resourceTranslation.ID", p.resourceTranslation.ID, true)
	}
	if p.new != nil {
		m.Set("new.lang", p.new.Lang, true)
		m.Set("new.resource", p.new.Resource, true)
		m.Set("new.K", p.new.K, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.lang", p.update.Lang, true)
		m.Set("update.resource", p.update.Resource, true)
		m.Set("update.K", p.update.K, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.translationID", p.filter.TranslationID, true)
		m.Set("filter.lang", p.filter.Lang, true)
		m.Set("filter.resource", p.filter.Resource, true)
		m.Set("filter.resourceType", p.filter.ResourceType, true)
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
func (p resourceTranslationActionProps) Format(in string, err error) string {
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

	if p.resourceTranslation != nil {
		// replacement for "{{resourceTranslation}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{resourceTranslation}}",
			fns(
				p.resourceTranslation.Lang,
				p.resourceTranslation.Resource,
				p.resourceTranslation.K,
				p.resourceTranslation.ID,
			),
		)
		pairs = append(pairs, "{{resourceTranslation.lang}}", fns(p.resourceTranslation.Lang))
		pairs = append(pairs, "{{resourceTranslation.resource}}", fns(p.resourceTranslation.Resource))
		pairs = append(pairs, "{{resourceTranslation.K}}", fns(p.resourceTranslation.K))
		pairs = append(pairs, "{{resourceTranslation.ID}}", fns(p.resourceTranslation.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Lang,
				p.new.Resource,
				p.new.K,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.lang}}", fns(p.new.Lang))
		pairs = append(pairs, "{{new.resource}}", fns(p.new.Resource))
		pairs = append(pairs, "{{new.K}}", fns(p.new.K))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Lang,
				p.update.Resource,
				p.update.K,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.lang}}", fns(p.update.Lang))
		pairs = append(pairs, "{{update.resource}}", fns(p.update.Resource))
		pairs = append(pairs, "{{update.K}}", fns(p.update.K))
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.TranslationID,
				p.filter.Lang,
				p.filter.Resource,
				p.filter.ResourceType,
				p.filter.OwnerID,
				p.filter.Deleted,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{{filter.translationID}}", fns(p.filter.TranslationID))
		pairs = append(pairs, "{{filter.lang}}", fns(p.filter.Lang))
		pairs = append(pairs, "{{filter.resource}}", fns(p.filter.Resource))
		pairs = append(pairs, "{{filter.resourceType}}", fns(p.filter.ResourceType))
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
func (a *resourceTranslationAction) String() string {
	var props = &resourceTranslationActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *resourceTranslationAction) ToAction() *actionlog.Action {
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

// ResourceTranslationActionSearch returns "system:resource-translation.search" action
//
// This function is auto-generated.
//
func ResourceTranslationActionSearch(props ...*resourceTranslationActionProps) *resourceTranslationAction {
	a := &resourceTranslationAction{
		timestamp: time.Now(),
		resource:  "system:resource-translation",
		action:    "search",
		log:       "searched for resourceTranslations",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ResourceTranslationActionLookup returns "system:resource-translation.lookup" action
//
// This function is auto-generated.
//
func ResourceTranslationActionLookup(props ...*resourceTranslationActionProps) *resourceTranslationAction {
	a := &resourceTranslationAction{
		timestamp: time.Now(),
		resource:  "system:resource-translation",
		action:    "lookup",
		log:       "looked-up for a {{resourceTranslation}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ResourceTranslationActionCreate returns "system:resource-translation.create" action
//
// This function is auto-generated.
//
func ResourceTranslationActionCreate(props ...*resourceTranslationActionProps) *resourceTranslationAction {
	a := &resourceTranslationAction{
		timestamp: time.Now(),
		resource:  "system:resource-translation",
		action:    "create",
		log:       "created {{resourceTranslation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ResourceTranslationActionUpdate returns "system:resource-translation.update" action
//
// This function is auto-generated.
//
func ResourceTranslationActionUpdate(props ...*resourceTranslationActionProps) *resourceTranslationAction {
	a := &resourceTranslationAction{
		timestamp: time.Now(),
		resource:  "system:resource-translation",
		action:    "update",
		log:       "updated {{resourceTranslation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ResourceTranslationActionDelete returns "system:resource-translation.delete" action
//
// This function is auto-generated.
//
func ResourceTranslationActionDelete(props ...*resourceTranslationActionProps) *resourceTranslationAction {
	a := &resourceTranslationAction{
		timestamp: time.Now(),
		resource:  "system:resource-translation",
		action:    "delete",
		log:       "deleted {{resourceTranslation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ResourceTranslationActionUndelete returns "system:resource-translation.undelete" action
//
// This function is auto-generated.
//
func ResourceTranslationActionUndelete(props ...*resourceTranslationActionProps) *resourceTranslationAction {
	a := &resourceTranslationAction{
		timestamp: time.Now(),
		resource:  "system:resource-translation",
		action:    "undelete",
		log:       "undeleted {{resourceTranslation}}",
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

// ResourceTranslationErrGeneric returns "system:resource-translation.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ResourceTranslationErrGeneric(mm ...*resourceTranslationActionProps) *errors.Error {
	var p = &resourceTranslationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:resource-translation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(resourceTranslationLogMetaKey{}, "{err}"),
		errors.Meta(resourceTranslationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "resourceTranslation.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ResourceTranslationErrNotFound returns "system:resource-translation.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ResourceTranslationErrNotFound(mm ...*resourceTranslationActionProps) *errors.Error {
	var p = &resourceTranslationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("resource translation not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:resource-translation"),

		errors.Meta(resourceTranslationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "resourceTranslation.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ResourceTranslationErrInvalidID returns "system:resource-translation.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ResourceTranslationErrInvalidID(mm ...*resourceTranslationActionProps) *errors.Error {
	var p = &resourceTranslationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:resource-translation"),

		errors.Meta(resourceTranslationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "resourceTranslation.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ResourceTranslationErrNotAllowedToRead returns "system:resource-translation.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func ResourceTranslationErrNotAllowedToRead(mm ...*resourceTranslationActionProps) *errors.Error {
	var p = &resourceTranslationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this resource translation", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:resource-translation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(resourceTranslationLogMetaKey{}, "failed to read {{resourceTranslation}}; insufficient permissions"),
		errors.Meta(resourceTranslationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "resourceTranslation.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ResourceTranslationErrNotAllowedToSearch returns "system:resource-translation.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func ResourceTranslationErrNotAllowedToSearch(mm ...*resourceTranslationActionProps) *errors.Error {
	var p = &resourceTranslationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list resource translations", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:resource-translation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(resourceTranslationLogMetaKey{}, "failed to search or list resource translations; insufficient permissions"),
		errors.Meta(resourceTranslationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "resourceTranslation.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ResourceTranslationErrNotAllowedToManage returns "system:resource-translation.notAllowedToManage" as *errors.Error
//
//
// This function is auto-generated.
//
func ResourceTranslationErrNotAllowedToManage(mm ...*resourceTranslationActionProps) *errors.Error {
	var p = &resourceTranslationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage resource translations", nil),

		errors.Meta("type", "notAllowedToManage"),
		errors.Meta("resource", "system:resource-translation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(resourceTranslationLogMetaKey{}, "failed to manage resource translation; insufficient permissions"),
		errors.Meta(resourceTranslationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "resourceTranslation.errors.notAllowedToManage"),

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
func (svc resourceTranslation) recordAction(ctx context.Context, props *resourceTranslationActionProps, actionFn func(...*resourceTranslationActionProps) *resourceTranslationAction, err error) error {
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
		a.Description = props.Format(m.AsString(resourceTranslationLogMetaKey{}), err)

		if p, has := m[resourceTranslationPropsMetaKey{}]; has {
			a.Meta = p.(*resourceTranslationActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
