package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/tenant_actions.yaml

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/system/types"
	"strings"
	"time"
)

type (
	tenantActionProps struct {
		tenant *types.Tenant
		new    *types.Tenant
		update *types.Tenant
		search *types.TenantFilter
		diff   []*revisions.Change
		old    json.RawMessage
	}

	tenantAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *tenantActionProps
	}

	tenantLogMetaKey   struct{}
	tenantPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setTenant updates tenantActionProps's tenant
//
// This function is auto-generated.
func (p *tenantActionProps) setTenant(tenant *types.Tenant) *tenantActionProps {
	p.tenant = tenant
	return p
}

// setNew updates tenantActionProps's new
//
// This function is auto-generated.
func (p *tenantActionProps) setNew(new *types.Tenant) *tenantActionProps {
	p.new = new
	return p
}

// setUpdate updates tenantActionProps's update
//
// This function is auto-generated.
func (p *tenantActionProps) setUpdate(update *types.Tenant) *tenantActionProps {
	p.update = update
	return p
}

// setSearch updates tenantActionProps's search
//
// This function is auto-generated.
func (p *tenantActionProps) setSearch(search *types.TenantFilter) *tenantActionProps {
	p.search = search
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *tenantActionProps) setDiff(diff []*revisions.Change) *tenantActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *tenantActionProps) setOld(v any) *tenantActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts tenantActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p tenantActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.tenant != nil {
		m.Set("tenant.handle", p.tenant.Handle, true)
		m.Set("tenant.ID", p.tenant.ID, true)
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
func (p tenantActionProps) Format(in string, err error) string {
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

	if p.tenant != nil {
		// replacement for "{{tenant}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{tenant}}",
			fns(
				p.tenant.Handle,
				p.tenant.ID,
			),
		)
		pairs = append(pairs, "{{tenant.handle}}", fns(p.tenant.Handle))
		pairs = append(pairs, "{{tenant.ID}}", fns(p.tenant.ID))
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
func (a *tenantAction) String() string {
	var props = &tenantActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *tenantAction) ToAction() *actionlog.Action {
	resource := e.resource
	var resourceProjectID uint64
	var resourceRevisionID uint64
	if e.props != nil && e.props.tenant != nil {
		if r, ok := any(e.props.tenant).(actionlog.RbacResourcer); ok {
			resource = r.RbacResource()
		}

		// Attribute the event to the project owning the affected resource. This is
		// independent of the request scope, so the log stays filterable per project
		// on routes that carry no project.
		if r, ok := any(e.props.tenant).(actionlog.ProjectResourcer); ok {
			resourceProjectID = r.ProjectRef()
		}

		// Attribute the event to the revision owning (or assigned to) the affected
		// resource. Work items (RevisionResourcer) file against the chain root and
		// carry their own, independent revision assignment. Everything else already
		// denormalises the revision onto ProjectRef() -- every projects row is its
		// own revision, so whichever revision project row a resource lives under
		// already is one.
		if r, ok := any(e.props.tenant).(actionlog.RevisionResourcer); ok {
			resourceRevisionID = r.RevisionRef()
		} else {
			resourceRevisionID = resourceProjectID
		}
	}
	return &actionlog.Action{
		Resource:           resource,
		ResourceProjectID:  resourceProjectID,
		ResourceRevisionID: resourceRevisionID,
		Action:             e.action,
		Severity:           e.severity,
		Description:        e.String(),
		Meta:               e.props.Serialize(),
		Delta:              actionlog.Delta(e.props.diff),
		OldState:           actionlog.OldState(e.props.old),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// TenantActionSearch returns "system:tenant.search" action
//
// This function is auto-generated.
func TenantActionSearch(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "search",
		log:       "searched for tenants",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionLookup returns "system:tenant.lookup" action
//
// This function is auto-generated.
func TenantActionLookup(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "lookup",
		log:       "looked-up for a {{tenant}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionCreate returns "system:tenant.create" action
//
// This function is auto-generated.
func TenantActionCreate(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "create",
		log:       "created {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionUpdate returns "system:tenant.update" action
//
// This function is auto-generated.
func TenantActionUpdate(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "update",
		log:       "updated {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionDelete returns "system:tenant.delete" action
//
// This function is auto-generated.
func TenantActionDelete(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "delete",
		log:       "deleted {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionUndelete returns "system:tenant.undelete" action
//
// This function is auto-generated.
func TenantActionUndelete(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "undelete",
		log:       "undeleted {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionSuspend returns "system:tenant.suspend" action
//
// This function is auto-generated.
func TenantActionSuspend(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "suspend",
		log:       "suspended {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionActivate returns "system:tenant.activate" action
//
// This function is auto-generated.
func TenantActionActivate(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "activate",
		log:       "activated {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionArchive returns "system:tenant.archive" action
//
// This function is auto-generated.
func TenantActionArchive(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "archive",
		log:       "archived {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionSearchMembers returns "system:tenant.searchMembers" action
//
// This function is auto-generated.
func TenantActionSearchMembers(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "searchMembers",
		log:       "searched members of {{tenant}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionInvite returns "system:tenant.invite" action
//
// This function is auto-generated.
func TenantActionInvite(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "invite",
		log:       "invited member to {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionAcceptInvite returns "system:tenant.acceptInvite" action
//
// This function is auto-generated.
func TenantActionAcceptInvite(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "acceptInvite",
		log:       "accepted invite to {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionUpdateMember returns "system:tenant.updateMember" action
//
// This function is auto-generated.
func TenantActionUpdateMember(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "updateMember",
		log:       "updated member of {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionRemoveMember returns "system:tenant.removeMember" action
//
// This function is auto-generated.
func TenantActionRemoveMember(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "removeMember",
		log:       "removed member from {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionSuspendMember returns "system:tenant.suspendMember" action
//
// This function is auto-generated.
func TenantActionSuspendMember(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "suspendMember",
		log:       "suspended member of {{tenant}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// TenantActionActivateMember returns "system:tenant.activateMember" action
//
// This function is auto-generated.
func TenantActionActivateMember(props ...*tenantActionProps) *tenantAction {
	a := &tenantAction{
		timestamp: time.Now(),
		resource:  "system:tenant",
		action:    "activateMember",
		log:       "activated member of {{tenant}}",
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

// TenantErrGeneric returns "system:tenant.generic" as *errors.Error
//
// This function is auto-generated.
func TenantErrGeneric(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:tenant"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(tenantLogMetaKey{}, "{err}"),
		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrNotFound returns "system:tenant.notFound" as *errors.Error
//
// This function is auto-generated.
func TenantErrNotFound(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("tenant not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrInvalidID returns "system:tenant.invalidID" as *errors.Error
//
// This function is auto-generated.
func TenantErrInvalidID(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrInvalidHandle returns "system:tenant.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func TenantErrInvalidHandle(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrHandleNotUnique returns "system:tenant.handleNotUnique" as *errors.Error
//
// This function is auto-generated.
func TenantErrHandleNotUnique(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrInvalidStatus returns "system:tenant.invalidStatus" as *errors.Error
//
// This function is auto-generated.
func TenantErrInvalidStatus(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid tenant status", nil),

		errors.Meta("type", "invalidStatus"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.invalidStatus"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrStaleData returns "system:tenant.staleData" as *errors.Error
//
// This function is auto-generated.
func TenantErrStaleData(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("tenant was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrNotAllowedToCreate returns "system:tenant.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func TenantErrNotAllowedToCreate(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a tenant", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:tenant"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(tenantLogMetaKey{}, "failed to create a tenant; insufficient permissions"),
		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrNotAllowedToRead returns "system:tenant.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func TenantErrNotAllowedToRead(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this tenant", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:tenant"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(tenantLogMetaKey{}, "failed to read {{tenant.handle}}; insufficient permissions"),
		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrNotAllowedToSearch returns "system:tenant.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func TenantErrNotAllowedToSearch(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list tenants", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:tenant"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(tenantLogMetaKey{}, "failed to search or list tenants; insufficient permissions"),
		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrNotAllowedToUpdate returns "system:tenant.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func TenantErrNotAllowedToUpdate(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this tenant", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:tenant"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(tenantLogMetaKey{}, "failed to update {{tenant.handle}}; insufficient permissions"),
		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrNotAllowedToDelete returns "system:tenant.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func TenantErrNotAllowedToDelete(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this tenant", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:tenant"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(tenantLogMetaKey{}, "failed to delete {{tenant.handle}}; insufficient permissions"),
		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrNotAllowedToSuspend returns "system:tenant.notAllowedToSuspend" as *errors.Error
//
// This function is auto-generated.
func TenantErrNotAllowedToSuspend(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to suspend this tenant", nil),

		errors.Meta("type", "notAllowedToSuspend"),
		errors.Meta("resource", "system:tenant"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(tenantLogMetaKey{}, "failed to suspend {{tenant.handle}}; insufficient permissions"),
		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.notAllowedToSuspend"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrNotAllowedToManageMembers returns "system:tenant.notAllowedToManageMembers" as *errors.Error
//
// This function is auto-generated.
func TenantErrNotAllowedToManageMembers(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage members of this tenant", nil),

		errors.Meta("type", "notAllowedToManageMembers"),
		errors.Meta("resource", "system:tenant"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(tenantLogMetaKey{}, "failed to manage members of {{tenant.handle}}; insufficient permissions"),
		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.notAllowedToManageMembers"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrInvalidRole returns "system:tenant.invalidRole" as *errors.Error
//
// This function is auto-generated.
func TenantErrInvalidRole(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid tenant member role", nil),

		errors.Meta("type", "invalidRole"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.invalidRole"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrInvalidMemberStatus returns "system:tenant.invalidMemberStatus" as *errors.Error
//
// This function is auto-generated.
func TenantErrInvalidMemberStatus(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid tenant member status", nil),

		errors.Meta("type", "invalidMemberStatus"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.invalidMemberStatus"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrMemberNotFound returns "system:tenant.memberNotFound" as *errors.Error
//
// This function is auto-generated.
func TenantErrMemberNotFound(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("tenant member not found", nil),

		errors.Meta("type", "memberNotFound"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.memberNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrMemberAlreadyExists returns "system:tenant.memberAlreadyExists" as *errors.Error
//
// This function is auto-generated.
func TenantErrMemberAlreadyExists(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("user is already a member of this tenant", nil),

		errors.Meta("type", "memberAlreadyExists"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.memberAlreadyExists"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// TenantErrUserAlreadyInTenant returns "system:tenant.userAlreadyInTenant" as *errors.Error
//
// This function is auto-generated.
func TenantErrUserAlreadyInTenant(mm ...*tenantActionProps) *errors.Error {
	var p = &tenantActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("user already holds an active membership in another tenant", nil),

		errors.Meta("type", "userAlreadyInTenant"),
		errors.Meta("resource", "system:tenant"),

		errors.Meta(tenantPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "tenant.errors.userAlreadyInTenant"),

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
func (svc tenant) recordAction(ctx context.Context, props *tenantActionProps, actionFn func(...*tenantActionProps) *tenantAction, err error, diff ...any) error {
	if len(diff) == 2 && diff[0] != nil {
		props.setDiff(actionlog.DiffResourceState(diff[0], diff[1]))
		props.setOld(diff[0])
	}
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
		a.Description = props.Format(m.AsString(tenantLogMetaKey{}), err)

		if p, has := m[tenantPropsMetaKey{}]; has {
			a.Meta = p.(*tenantActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
