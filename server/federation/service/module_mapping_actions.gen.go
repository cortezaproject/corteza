package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// federation/service/module_mapping_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"strings"
	"time"
)

type (
	moduleMappingActionProps struct {
		created *types.ModuleMapping
		mapping *types.ModuleMapping
		changed *types.ModuleMapping
		filter  *types.ModuleMappingFilter
	}

	moduleMappingAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *moduleMappingActionProps
	}

	moduleMappingLogMetaKey   struct{}
	moduleMappingPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setCreated updates moduleMappingActionProps's created
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleMappingActionProps) setCreated(created *types.ModuleMapping) *moduleMappingActionProps {
	p.created = created
	return p
}

// setMapping updates moduleMappingActionProps's mapping
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleMappingActionProps) setMapping(mapping *types.ModuleMapping) *moduleMappingActionProps {
	p.mapping = mapping
	return p
}

// setChanged updates moduleMappingActionProps's changed
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleMappingActionProps) setChanged(changed *types.ModuleMapping) *moduleMappingActionProps {
	p.changed = changed
	return p
}

// setFilter updates moduleMappingActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *moduleMappingActionProps) setFilter(filter *types.ModuleMappingFilter) *moduleMappingActionProps {
	p.filter = filter
	return p
}

// Serialize converts moduleMappingActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p moduleMappingActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.created != nil {
		m.Set("created.FederationModuleID", p.created.FederationModuleID, true)
		m.Set("created.ComposeModuleID", p.created.ComposeModuleID, true)
	}
	if p.mapping != nil {
		m.Set("mapping.FederationModuleID", p.mapping.FederationModuleID, true)
		m.Set("mapping.ComposeModuleID", p.mapping.ComposeModuleID, true)
	}
	if p.changed != nil {
		m.Set("changed.FederationModuleID", p.changed.FederationModuleID, true)
		m.Set("changed.ComposeModuleID", p.changed.ComposeModuleID, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.sort", p.filter.Sort, true)
		m.Set("filter.limit", p.filter.Limit, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p moduleMappingActionProps) Format(in string, err error) string {
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

	if p.created != nil {
		// replacement for "{{created}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{created}}",
			fns(
				p.created.FederationModuleID,
				p.created.ComposeModuleID,
			),
		)
		pairs = append(pairs, "{{created.FederationModuleID}}", fns(p.created.FederationModuleID))
		pairs = append(pairs, "{{created.ComposeModuleID}}", fns(p.created.ComposeModuleID))
	}

	if p.mapping != nil {
		// replacement for "{{mapping}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{mapping}}",
			fns(
				p.mapping.FederationModuleID,
				p.mapping.ComposeModuleID,
			),
		)
		pairs = append(pairs, "{{mapping.FederationModuleID}}", fns(p.mapping.FederationModuleID))
		pairs = append(pairs, "{{mapping.ComposeModuleID}}", fns(p.mapping.ComposeModuleID))
	}

	if p.changed != nil {
		// replacement for "{{changed}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{changed}}",
			fns(
				p.changed.FederationModuleID,
				p.changed.ComposeModuleID,
			),
		)
		pairs = append(pairs, "{{changed.FederationModuleID}}", fns(p.changed.FederationModuleID))
		pairs = append(pairs, "{{changed.ComposeModuleID}}", fns(p.changed.ComposeModuleID))
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
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
//
func (a *moduleMappingAction) String() string {
	var props = &moduleMappingActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *moduleMappingAction) ToAction() *actionlog.Action {
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

// ModuleMappingActionSearch returns "federation:module_mapping.search" action
//
// This function is auto-generated.
//
func ModuleMappingActionSearch(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "search",
		log:       "searched for modules",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleMappingActionLookup returns "federation:module_mapping.lookup" action
//
// This function is auto-generated.
//
func ModuleMappingActionLookup(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "lookup",
		log:       "looked-up for a module",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleMappingActionCreate returns "federation:module_mapping.create" action
//
// This function is auto-generated.
//
func ModuleMappingActionCreate(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "create",
		log:       "created module",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleMappingActionUpdate returns "federation:module_mapping.update" action
//
// This function is auto-generated.
//
func ModuleMappingActionUpdate(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "update",
		log:       "updated module",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ModuleMappingActionDelete returns "federation:module_mapping.delete" action
//
// This function is auto-generated.
//
func ModuleMappingActionDelete(props ...*moduleMappingActionProps) *moduleMappingAction {
	a := &moduleMappingAction{
		timestamp: time.Now(),
		resource:  "federation:module_mapping",
		action:    "delete",
		log:       "deleted module",
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

// ModuleMappingErrGeneric returns "federation:module_mapping.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrGeneric(mm ...*moduleMappingActionProps) *errors.Error {
	var p = &moduleMappingActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "federation:module_mapping"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleMappingLogMetaKey{}, "{err}"),
		errors.Meta(moduleMappingPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "module-mapping.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleMappingErrNotFound returns "federation:module_mapping.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrNotFound(mm ...*moduleMappingActionProps) *errors.Error {
	var p = &moduleMappingActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module mapping does not exist", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "federation:module_mapping"),

		errors.Meta(moduleMappingPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "module-mapping.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleMappingErrComposeModuleNotFound returns "federation:module_mapping.composeModuleNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrComposeModuleNotFound(mm ...*moduleMappingActionProps) *errors.Error {
	var p = &moduleMappingActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("compose module not found", nil),

		errors.Meta("type", "composeModuleNotFound"),
		errors.Meta("resource", "federation:module_mapping"),

		errors.Meta(moduleMappingPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "module-mapping.errors.composeModuleNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleMappingErrComposeNamespaceNotFound returns "federation:module_mapping.composeNamespaceNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrComposeNamespaceNotFound(mm ...*moduleMappingActionProps) *errors.Error {
	var p = &moduleMappingActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("compose namespace not found", nil),

		errors.Meta("type", "composeNamespaceNotFound"),
		errors.Meta("resource", "federation:module_mapping"),

		errors.Meta(moduleMappingPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "module-mapping.errors.composeNamespaceNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleMappingErrFederationModuleNotFound returns "federation:module_mapping.federationModuleNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrFederationModuleNotFound(mm ...*moduleMappingActionProps) *errors.Error {
	var p = &moduleMappingActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("federation module not found", nil),

		errors.Meta("type", "federationModuleNotFound"),
		errors.Meta("resource", "federation:module_mapping"),

		errors.Meta(moduleMappingPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "module-mapping.errors.federationModuleNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleMappingErrNodeNotFound returns "federation:module_mapping.nodeNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrNodeNotFound(mm ...*moduleMappingActionProps) *errors.Error {
	var p = &moduleMappingActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("node does not exist", nil),

		errors.Meta("type", "nodeNotFound"),
		errors.Meta("resource", "federation:module_mapping"),

		errors.Meta(moduleMappingPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "module-mapping.errors.nodeNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleMappingErrModuleMappingExists returns "federation:module_mapping.moduleMappingExists" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrModuleMappingExists(mm ...*moduleMappingActionProps) *errors.Error {
	var p = &moduleMappingActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module mapping already exists", nil),

		errors.Meta("type", "moduleMappingExists"),
		errors.Meta("resource", "federation:module_mapping"),

		errors.Meta(moduleMappingPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "module-mapping.errors.moduleMappingExists"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ModuleMappingErrNotAllowedToMap returns "federation:module_mapping.notAllowedToMap" as *errors.Error
//
//
// This function is auto-generated.
//
func ModuleMappingErrNotAllowedToMap(mm ...*moduleMappingActionProps) *errors.Error {
	var p = &moduleMappingActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to map this module", nil),

		errors.Meta("type", "notAllowedToMap"),
		errors.Meta("resource", "federation:module_mapping"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(moduleMappingLogMetaKey{}, "could not manage mapping; insufficient permissions"),
		errors.Meta(moduleMappingPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "federation"),
		errors.Meta(locale.ErrorMetaKey{}, "module-mapping.errors.notAllowedToMap"),

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
func (svc moduleMapping) recordAction(ctx context.Context, props *moduleMappingActionProps, actionFn func(...*moduleMappingActionProps) *moduleMappingAction, err error) error {
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
		a.Description = props.Format(m.AsString(moduleMappingLogMetaKey{}), err)

		if p, has := m[moduleMappingPropsMetaKey{}]; has {
			a.Meta = p.(*moduleMappingActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
