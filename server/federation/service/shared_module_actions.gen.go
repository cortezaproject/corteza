package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// federation/service/shared_module_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"strings"
	"time"
)

type (
	sharedModuleActionProps struct {
		module  *types.SharedModule
		changed *types.SharedModule
		filter  *types.SharedModuleFilter
		node    *types.Node
	}

	sharedModuleAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *sharedModuleActionProps
	}

	sharedModuleLogMetaKey   struct{}
	sharedModulePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setModule updates sharedModuleActionProps's module
//
// This function is auto-generated.
//
func (p *sharedModuleActionProps) setModule(module *types.SharedModule) *sharedModuleActionProps {
	p.module = module
	return p
}

// setChanged updates sharedModuleActionProps's changed
//
// This function is auto-generated.
//
func (p *sharedModuleActionProps) setChanged(changed *types.SharedModule) *sharedModuleActionProps {
	p.changed = changed
	return p
}

// setFilter updates sharedModuleActionProps's filter
//
// This function is auto-generated.
//
func (p *sharedModuleActionProps) setFilter(filter *types.SharedModuleFilter) *sharedModuleActionProps {
	p.filter = filter
	return p
}

// setNode updates sharedModuleActionProps's node
//
// This function is auto-generated.
//
func (p *sharedModuleActionProps) setNode(node *types.Node) *sharedModuleActionProps {
	p.node = node
	return p
}

// Serialize converts sharedModuleActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p sharedModuleActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.module != nil {
		m.Set("module.ID", p.module.ID, true)
	}
	if p.changed != nil {
		m.Set("changed.ID", p.changed.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.sort", p.filter.Sort, true)
		m.Set("filter.limit", p.filter.Limit, true)
	}
	if p.node != nil {
		m.Set("node.ID", p.node.ID, true)
		m.Set("node.Name", p.node.Name, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p sharedModuleActionProps) Format(in string, err error) string {
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

	if p.module != nil {
		// replacement for "{{module}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{module}}",
			fns(
				p.module.ID,
			),
		)
		pairs = append(pairs, "{{module.ID}}", fns(p.module.ID))
	}

	if p.changed != nil {
		// replacement for "{{changed}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{changed}}",
			fns(
				p.changed.ID,
			),
		)
		pairs = append(pairs, "{{changed.ID}}", fns(p.changed.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Query,
				p.filter.Sort,
				p.filter.Limit,
			),
		)
		pairs = append(pairs, "{{filter.query}}", fns(p.filter.Query))
		pairs = append(pairs, "{{filter.sort}}", fns(p.filter.Sort))
		pairs = append(pairs, "{{filter.limit}}", fns(p.filter.Limit))
	}

	if p.node != nil {
		// replacement for "{{node}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{node}}",
			fns(
				p.node.ID,
				p.node.Name,
			),
		)
		pairs = append(pairs, "{{node.ID}}", fns(p.node.ID))
		pairs = append(pairs, "{{node.Name}}", fns(p.node.Name))
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
func (a *sharedModuleAction) String() string {
	var props = &sharedModuleActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *sharedModuleAction) ToAction() *actionlog.Action {
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

// SharedModuleActionSearch returns "federation:shared_module.search" action
//
// This function is auto-generated.
//
func SharedModuleActionSearch(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "search",
		log:       "searched for modules",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionLookup returns "federation:shared_module.lookup" action
//
// This function is auto-generated.
//
func SharedModuleActionLookup(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "lookup",
		log:       "looked-up for a {{module}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionCreate returns "federation:shared_module.create" action
//
// This function is auto-generated.
//
func SharedModuleActionCreate(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "create",
		log:       "created {{module}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionUpdate returns "federation:shared_module.update" action
//
// This function is auto-generated.
//
func SharedModuleActionUpdate(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "update",
		log:       "updated {{module}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionDelete returns "federation:shared_module.delete" action
//
// This function is auto-generated.
//
func SharedModuleActionDelete(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "delete",
		log:       "deleted {{module}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SharedModuleActionUndelete returns "federation:shared_module.undelete" action
//
// This function is auto-generated.
//
func SharedModuleActionUndelete(props ...*sharedModuleActionProps) *sharedModuleAction {
	a := &sharedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:shared_module",
		action:    "undelete",
		log:       "undeleted {{module}}",
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

// SharedModuleErrGeneric returns "federation:shared_module.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrGeneric(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "federation:shared_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sharedModuleLogMetaKey{}, "{err}"),
		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SharedModuleErrNotFound returns "federation:shared_module.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotFound(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "federation:shared_module"),

		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SharedModuleErrInvalidID returns "federation:shared_module.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrInvalidID(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "federation:shared_module"),

		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SharedModuleErrStaleData returns "federation:shared_module.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrStaleData(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "federation:shared_module"),

		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SharedModuleErrFederationSyncStructureChanged returns "federation:shared_module.federationSyncStructureChanged" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrFederationSyncStructureChanged(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module structure changed", nil),

		errors.Meta("type", "federationSyncStructureChanged"),
		errors.Meta("resource", "federation:shared_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sharedModuleLogMetaKey{}, "could not update shared module, structure different"),
		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.federationSyncStructureChanged"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SharedModuleErrNotUnique returns "federation:shared_module.notUnique" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotUnique(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("node not unique", nil),

		errors.Meta("type", "notUnique"),
		errors.Meta("resource", "federation:shared_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sharedModuleLogMetaKey{}, "used duplicate node TODO"),
		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.notUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SharedModuleErrNodeNotFound returns "federation:shared_module.nodeNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNodeNotFound(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("node does not exist", nil),

		errors.Meta("type", "nodeNotFound"),
		errors.Meta("resource", "federation:shared_module"),

		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.nodeNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SharedModuleErrNotAllowedToCreate returns "federation:shared_module.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotAllowedToCreate(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create modules", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "federation:shared_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sharedModuleLogMetaKey{}, "could not create modules; insufficient permissions"),
		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SharedModuleErrNotAllowedToManage returns "federation:shared_module.notAllowedToManage" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotAllowedToManage(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage this module", nil),

		errors.Meta("type", "notAllowedToManage"),
		errors.Meta("resource", "federation:shared_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sharedModuleLogMetaKey{}, "could not read {{module}}; insufficient permissions"),
		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.notAllowedToManage"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SharedModuleErrNotAllowedToMap returns "federation:shared_module.notAllowedToMap" as *errors.Error
//
//
// This function is auto-generated.
//
func SharedModuleErrNotAllowedToMap(mm ...*sharedModuleActionProps) *errors.Error {
	var p = &sharedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to map this module", nil),

		errors.Meta("type", "notAllowedToMap"),
		errors.Meta("resource", "federation:shared_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sharedModuleLogMetaKey{}, "could not map {{module}}; insufficient permissions"),
		errors.Meta(sharedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "shared-module.errors.notAllowedToMap"),

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
func (svc sharedModule) recordAction(ctx context.Context, props *sharedModuleActionProps, actionFn func(...*sharedModuleActionProps) *sharedModuleAction, err error) error {
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
		a.Description = props.Format(m.AsString(sharedModuleLogMetaKey{}), err)

		if p, has := m[sharedModulePropsMetaKey{}]; has {
			a.Meta = p.(*sharedModuleActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
