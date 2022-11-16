package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// federation/service/exposed_module_actions.yaml

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
	exposedModuleActionProps struct {
		module *types.ExposedModule
		update *types.ExposedModule
		create *types.ExposedModule
		delete *types.ExposedModule
		filter *types.ExposedModuleFilter
		node   *types.Node
	}

	exposedModuleAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *exposedModuleActionProps
	}

	exposedModuleLogMetaKey   struct{}
	exposedModulePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setModule updates exposedModuleActionProps's module
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setModule(module *types.ExposedModule) *exposedModuleActionProps {
	p.module = module
	return p
}

// setUpdate updates exposedModuleActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setUpdate(update *types.ExposedModule) *exposedModuleActionProps {
	p.update = update
	return p
}

// setCreate updates exposedModuleActionProps's create
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setCreate(create *types.ExposedModule) *exposedModuleActionProps {
	p.create = create
	return p
}

// setDelete updates exposedModuleActionProps's delete
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setDelete(delete *types.ExposedModule) *exposedModuleActionProps {
	p.delete = delete
	return p
}

// setFilter updates exposedModuleActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setFilter(filter *types.ExposedModuleFilter) *exposedModuleActionProps {
	p.filter = filter
	return p
}

// setNode updates exposedModuleActionProps's node
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *exposedModuleActionProps) setNode(node *types.Node) *exposedModuleActionProps {
	p.node = node
	return p
}

// Serialize converts exposedModuleActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p exposedModuleActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.module != nil {
		m.Set("module.ID", p.module.ID, true)
		m.Set("module.ComposeNamespaceID", p.module.ComposeNamespaceID, true)
		m.Set("module.ComposeModuleID", p.module.ComposeModuleID, true)
		m.Set("module.NodeID", p.module.NodeID, true)
	}
	if p.update != nil {
		m.Set("update.ID", p.update.ID, true)
		m.Set("update.ComposeNamespaceID", p.update.ComposeNamespaceID, true)
		m.Set("update.ComposeModuleID", p.update.ComposeModuleID, true)
		m.Set("update.NodeID", p.update.NodeID, true)
	}
	if p.create != nil {
		m.Set("create.ID", p.create.ID, true)
		m.Set("create.ComposeNamespaceID", p.create.ComposeNamespaceID, true)
		m.Set("create.ComposeModuleID", p.create.ComposeModuleID, true)
		m.Set("create.NodeID", p.create.NodeID, true)
	}
	if p.delete != nil {
		m.Set("delete.ID", p.delete.ID, true)
		m.Set("delete.ComposeNamespaceID", p.delete.ComposeNamespaceID, true)
		m.Set("delete.ComposeModuleID", p.delete.ComposeModuleID, true)
		m.Set("delete.NodeID", p.delete.NodeID, true)
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
func (p exposedModuleActionProps) Format(in string, err error) string {
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
				p.module.ComposeNamespaceID,
				p.module.ComposeModuleID,
				p.module.NodeID,
			),
		)
		pairs = append(pairs, "{{module.ID}}", fns(p.module.ID))
		pairs = append(pairs, "{{module.ComposeNamespaceID}}", fns(p.module.ComposeNamespaceID))
		pairs = append(pairs, "{{module.ComposeModuleID}}", fns(p.module.ComposeModuleID))
		pairs = append(pairs, "{{module.NodeID}}", fns(p.module.NodeID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.ID,
				p.update.ComposeNamespaceID,
				p.update.ComposeModuleID,
				p.update.NodeID,
			),
		)
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
		pairs = append(pairs, "{{update.ComposeNamespaceID}}", fns(p.update.ComposeNamespaceID))
		pairs = append(pairs, "{{update.ComposeModuleID}}", fns(p.update.ComposeModuleID))
		pairs = append(pairs, "{{update.NodeID}}", fns(p.update.NodeID))
	}

	if p.create != nil {
		// replacement for "{{create}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{create}}",
			fns(
				p.create.ID,
				p.create.ComposeNamespaceID,
				p.create.ComposeModuleID,
				p.create.NodeID,
			),
		)
		pairs = append(pairs, "{{create.ID}}", fns(p.create.ID))
		pairs = append(pairs, "{{create.ComposeNamespaceID}}", fns(p.create.ComposeNamespaceID))
		pairs = append(pairs, "{{create.ComposeModuleID}}", fns(p.create.ComposeModuleID))
		pairs = append(pairs, "{{create.NodeID}}", fns(p.create.NodeID))
	}

	if p.delete != nil {
		// replacement for "{{delete}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{delete}}",
			fns(
				p.delete.ID,
				p.delete.ComposeNamespaceID,
				p.delete.ComposeModuleID,
				p.delete.NodeID,
			),
		)
		pairs = append(pairs, "{{delete.ID}}", fns(p.delete.ID))
		pairs = append(pairs, "{{delete.ComposeNamespaceID}}", fns(p.delete.ComposeNamespaceID))
		pairs = append(pairs, "{{delete.ComposeModuleID}}", fns(p.delete.ComposeModuleID))
		pairs = append(pairs, "{{delete.NodeID}}", fns(p.delete.NodeID))
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
func (a *exposedModuleAction) String() string {
	var props = &exposedModuleActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *exposedModuleAction) ToAction() *actionlog.Action {
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

// ExposedModuleActionSearch returns "federation:exposed_module.search" action
//
// This function is auto-generated.
//
func ExposedModuleActionSearch(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "search",
		log:       "searched for modules",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionLookup returns "federation:exposed_module.lookup" action
//
// This function is auto-generated.
//
func ExposedModuleActionLookup(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "lookup",
		log:       "looked-up for a {{module}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionCreate returns "federation:exposed_module.create" action
//
// This function is auto-generated.
//
func ExposedModuleActionCreate(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "create",
		log:       "created {{module}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionUpdate returns "federation:exposed_module.update" action
//
// This function is auto-generated.
//
func ExposedModuleActionUpdate(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "update",
		log:       "updated {{module}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionDelete returns "federation:exposed_module.delete" action
//
// This function is auto-generated.
//
func ExposedModuleActionDelete(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
		action:    "delete",
		log:       "deleted {{module}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ExposedModuleActionUndelete returns "federation:exposed_module.undelete" action
//
// This function is auto-generated.
//
func ExposedModuleActionUndelete(props ...*exposedModuleActionProps) *exposedModuleAction {
	a := &exposedModuleAction{
		timestamp: time.Now(),
		resource:  "federation:exposed_module",
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

// ExposedModuleErrGeneric returns "federation:exposed_module.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrGeneric(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "federation:exposed_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(exposedModuleLogMetaKey{}, "{err}"),
		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrNotFound returns "federation:exposed_module.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotFound(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "federation:exposed_module"),

		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrInvalidID returns "federation:exposed_module.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrInvalidID(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "federation:exposed_module"),

		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrStaleData returns "federation:exposed_module.staleData" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrStaleData(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "federation:exposed_module"),

		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrNotUnique returns "federation:exposed_module.notUnique" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotUnique(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("node not unique", nil),

		errors.Meta("type", "notUnique"),
		errors.Meta("resource", "federation:exposed_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(exposedModuleLogMetaKey{}, "used duplicate node TODO - {{module.NodeID}} for this compose module TODO - module.rel_compose_module"),
		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.notUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrNodeNotFound returns "federation:exposed_module.nodeNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNodeNotFound(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("node does not exist", nil),

		errors.Meta("type", "nodeNotFound"),
		errors.Meta("resource", "federation:exposed_module"),

		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.nodeNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrComposeModuleNotFound returns "federation:exposed_module.composeModuleNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrComposeModuleNotFound(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("compose module not found", nil),

		errors.Meta("type", "composeModuleNotFound"),
		errors.Meta("resource", "federation:exposed_module"),

		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.composeModuleNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrComposeNamespaceNotFound returns "federation:exposed_module.composeNamespaceNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrComposeNamespaceNotFound(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("compose namespace not found", nil),

		errors.Meta("type", "composeNamespaceNotFound"),
		errors.Meta("resource", "federation:exposed_module"),

		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.composeNamespaceNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrRequestParametersInvalid returns "federation:exposed_module.requestParametersInvalid" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrRequestParametersInvalid(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("request parameters invalid", nil),

		errors.Meta("type", "requestParametersInvalid"),
		errors.Meta("resource", "federation:exposed_module"),

		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.requestParametersInvalid"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrNotAllowedToCreate returns "federation:exposed_module.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotAllowedToCreate(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create modules", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "federation:exposed_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(exposedModuleLogMetaKey{}, "could not create modules; insufficient permissions"),
		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ExposedModuleErrNotAllowedToManage returns "federation:exposed_module.notAllowedToManage" as *errors.Error
//
//
// This function is auto-generated.
//
func ExposedModuleErrNotAllowedToManage(mm ...*exposedModuleActionProps) *errors.Error {
	var p = &exposedModuleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage this module", nil),

		errors.Meta("type", "notAllowedToManage"),
		errors.Meta("resource", "federation:exposed_module"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(exposedModuleLogMetaKey{}, "could not manage {{module}}; insufficient permissions"),
		errors.Meta(exposedModulePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "exposedModule.errors.notAllowedToManage"),

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
func (svc exposedModule) recordAction(ctx context.Context, props *exposedModuleActionProps, actionFn func(...*exposedModuleActionProps) *exposedModuleAction, err error) error {
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
		a.Description = props.Format(m.AsString(exposedModuleLogMetaKey{}), err)

		if p, has := m[exposedModulePropsMetaKey{}]; has {
			a.Meta = p.(*exposedModuleActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
