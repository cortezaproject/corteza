package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/role_actions.yaml

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
	roleActionProps struct {
		member   *types.User
		role     *types.Role
		new      *types.Role
		update   *types.Role
		existing *types.Role
		target   *types.Role
		filter   *types.RoleFilter
	}

	roleAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *roleActionProps
	}

	roleLogMetaKey   struct{}
	rolePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setMember updates roleActionProps's member
//
// This function is auto-generated.
func (p *roleActionProps) setMember(member *types.User) *roleActionProps {
	p.member = member
	return p
}

// setRole updates roleActionProps's role
//
// This function is auto-generated.
func (p *roleActionProps) setRole(role *types.Role) *roleActionProps {
	p.role = role
	return p
}

// setNew updates roleActionProps's new
//
// This function is auto-generated.
func (p *roleActionProps) setNew(new *types.Role) *roleActionProps {
	p.new = new
	return p
}

// setUpdate updates roleActionProps's update
//
// This function is auto-generated.
func (p *roleActionProps) setUpdate(update *types.Role) *roleActionProps {
	p.update = update
	return p
}

// setExisting updates roleActionProps's existing
//
// This function is auto-generated.
func (p *roleActionProps) setExisting(existing *types.Role) *roleActionProps {
	p.existing = existing
	return p
}

// setTarget updates roleActionProps's target
//
// This function is auto-generated.
func (p *roleActionProps) setTarget(target *types.Role) *roleActionProps {
	p.target = target
	return p
}

// setFilter updates roleActionProps's filter
//
// This function is auto-generated.
func (p *roleActionProps) setFilter(filter *types.RoleFilter) *roleActionProps {
	p.filter = filter
	return p
}

// Serialize converts roleActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p roleActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.member != nil {
		m.Set("member.handle", p.member.Handle, true)
		m.Set("member.email", p.member.Email, true)
		m.Set("member.name", p.member.Name, true)
		m.Set("member.ID", p.member.ID, true)
	}
	if p.role != nil {
		m.Set("role.handle", p.role.Handle, true)
		m.Set("role.name", p.role.Name, true)
		m.Set("role.ID", p.role.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.name", p.new.Name, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.name", p.update.Name, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.existing != nil {
		m.Set("existing.handle", p.existing.Handle, true)
		m.Set("existing.name", p.existing.Name, true)
		m.Set("existing.ID", p.existing.ID, true)
	}
	if p.target != nil {
		m.Set("target.handle", p.target.Handle, true)
		m.Set("target.name", p.target.Name, true)
		m.Set("target.ID", p.target.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.roleID", p.filter.RoleID, true)
		m.Set("filter.memberID", p.filter.MemberID, true)
		m.Set("filter.handle", p.filter.Handle, true)
		m.Set("filter.name", p.filter.Name, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
		m.Set("filter.archived", p.filter.Archived, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p roleActionProps) Format(in string, err error) string {
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

	if p.member != nil {
		// replacement for "{{member}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{member}}",
			fns(
				p.member.Handle,
				p.member.Email,
				p.member.Name,
				p.member.ID,
			),
		)
		pairs = append(pairs, "{{member.handle}}", fns(p.member.Handle))
		pairs = append(pairs, "{{member.email}}", fns(p.member.Email))
		pairs = append(pairs, "{{member.name}}", fns(p.member.Name))
		pairs = append(pairs, "{{member.ID}}", fns(p.member.ID))
	}

	if p.role != nil {
		// replacement for "{{role}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{role}}",
			fns(
				p.role.Handle,
				p.role.Name,
				p.role.ID,
			),
		)
		pairs = append(pairs, "{{role.handle}}", fns(p.role.Handle))
		pairs = append(pairs, "{{role.name}}", fns(p.role.Name))
		pairs = append(pairs, "{{role.ID}}", fns(p.role.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.Handle,
				p.new.Name,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.handle}}", fns(p.new.Handle))
		pairs = append(pairs, "{{new.name}}", fns(p.new.Name))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.Handle,
				p.update.Name,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.handle}}", fns(p.update.Handle))
		pairs = append(pairs, "{{update.name}}", fns(p.update.Name))
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.existing != nil {
		// replacement for "{{existing}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{existing}}",
			fns(
				p.existing.Handle,
				p.existing.Name,
				p.existing.ID,
			),
		)
		pairs = append(pairs, "{{existing.handle}}", fns(p.existing.Handle))
		pairs = append(pairs, "{{existing.name}}", fns(p.existing.Name))
		pairs = append(pairs, "{{existing.ID}}", fns(p.existing.ID))
	}

	if p.target != nil {
		// replacement for "{{target}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{target}}",
			fns(
				p.target.Handle,
				p.target.Name,
				p.target.ID,
			),
		)
		pairs = append(pairs, "{{target.handle}}", fns(p.target.Handle))
		pairs = append(pairs, "{{target.name}}", fns(p.target.Name))
		pairs = append(pairs, "{{target.ID}}", fns(p.target.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Query,
				p.filter.RoleID,
				p.filter.MemberID,
				p.filter.Handle,
				p.filter.Name,
				p.filter.Deleted,
				p.filter.Archived,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{{filter.query}}", fns(p.filter.Query))
		pairs = append(pairs, "{{filter.roleID}}", fns(p.filter.RoleID))
		pairs = append(pairs, "{{filter.memberID}}", fns(p.filter.MemberID))
		pairs = append(pairs, "{{filter.handle}}", fns(p.filter.Handle))
		pairs = append(pairs, "{{filter.name}}", fns(p.filter.Name))
		pairs = append(pairs, "{{filter.deleted}}", fns(p.filter.Deleted))
		pairs = append(pairs, "{{filter.archived}}", fns(p.filter.Archived))
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
func (a *roleAction) String() string {
	var props = &roleActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *roleAction) ToAction() *actionlog.Action {
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

// RoleActionSearch returns "system:role.search" action
//
// This function is auto-generated.
func RoleActionSearch(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "search",
		log:       "searched for roles",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionLookup returns "system:role.lookup" action
//
// This function is auto-generated.
func RoleActionLookup(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "lookup",
		log:       "looked-up for a {{role}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionCreate returns "system:role.create" action
//
// This function is auto-generated.
func RoleActionCreate(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "create",
		log:       "created {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionUpdate returns "system:role.update" action
//
// This function is auto-generated.
func RoleActionUpdate(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "update",
		log:       "updated {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionDelete returns "system:role.delete" action
//
// This function is auto-generated.
func RoleActionDelete(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "delete",
		log:       "deleted {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionUndelete returns "system:role.undelete" action
//
// This function is auto-generated.
func RoleActionUndelete(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "undelete",
		log:       "undeleted {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionArchive returns "system:role.archive" action
//
// This function is auto-generated.
func RoleActionArchive(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "archive",
		log:       "archived {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionUnarchive returns "system:role.unarchive" action
//
// This function is auto-generated.
func RoleActionUnarchive(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "unarchive",
		log:       "unarchived {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionMerge returns "system:role.merge" action
//
// This function is auto-generated.
func RoleActionMerge(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "merge",
		log:       "merged {{target}} with {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionMembers returns "system:role.members" action
//
// This function is auto-generated.
func RoleActionMembers(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "members",
		log:       "searched for members of {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionMemberAdd returns "system:role.memberAdd" action
//
// This function is auto-generated.
func RoleActionMemberAdd(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "memberAdd",
		log:       "added {{member.email}} to {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionMemberRemove returns "system:role.memberRemove" action
//
// This function is auto-generated.
func RoleActionMemberRemove(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "memberRemove",
		log:       "removed {{member.email}} from {{role}}",
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

// RoleErrGeneric returns "system:role.generic" as *errors.Error
//
// This function is auto-generated.
func RoleErrGeneric(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "{err}"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotFound returns "system:role.notFound" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotFound(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("role not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:role"),

		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrInvalidID returns "system:role.invalidID" as *errors.Error
//
// This function is auto-generated.
func RoleErrInvalidID(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:role"),

		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrInvalidHandle returns "system:role.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func RoleErrInvalidHandle(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:role"),

		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToRead returns "system:role.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToRead(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this role", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to read {{role.handle}}; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToSearch returns "system:role.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToSearch(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list roles", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to search or list roles; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToCreate returns "system:role.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToCreate(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create roles", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to create role; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToUpdate returns "system:role.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToUpdate(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this role", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to update {{role.handle}}; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToDelete returns "system:role.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToDelete(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this role", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to delete {{role.handle}}; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToUndelete returns "system:role.notAllowedToUndelete" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToUndelete(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this role", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to undelete {{role.handle}}; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToArchive returns "system:role.notAllowedToArchive" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToArchive(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to archive this role", nil),

		errors.Meta("type", "notAllowedToArchive"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to archive {{role.handle}}; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToArchive"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToUnarchive returns "system:role.notAllowedToUnarchive" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToUnarchive(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to unarchive this role", nil),

		errors.Meta("type", "notAllowedToUnarchive"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to unarchive {{role.handle}}; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToUnarchive"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToCloneRules returns "system:role.notAllowedToCloneRules" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToCloneRules(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to clone rules of this role", nil),

		errors.Meta("type", "notAllowedToCloneRules"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to clone rules of {{role.handle}}; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToCloneRules"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNotAllowedToManageMembers returns "system:role.notAllowedToManageMembers" as *errors.Error
//
// This function is auto-generated.
func RoleErrNotAllowedToManageMembers(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage role members", nil),

		errors.Meta("type", "notAllowedToManageMembers"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "failed to manage {{role.handle}} members; insufficient permissions"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.notAllowedToManageMembers"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrHandleNotUnique returns "system:role.handleNotUnique" as *errors.Error
//
// This function is auto-generated.
func RoleErrHandleNotUnique(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("role handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "used duplicate handle ({{role.handle}}) for role"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// RoleErrNameNotUnique returns "system:role.nameNotUnique" as *errors.Error
//
// This function is auto-generated.
func RoleErrNameNotUnique(mm ...*roleActionProps) *errors.Error {
	var p = &roleActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("role name not unique", nil),

		errors.Meta("type", "nameNotUnique"),
		errors.Meta("resource", "system:role"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(roleLogMetaKey{}, "used duplicate name ({{role.name}}) for role"),
		errors.Meta(rolePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "role.errors.nameNotUnique"),

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
func (svc role) recordAction(ctx context.Context, props *roleActionProps, actionFn func(...*roleActionProps) *roleAction, err error) error {
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
		a.Description = props.Format(m.AsString(roleLogMetaKey{}), err)

		if p, has := m[rolePropsMetaKey{}]; has {
			a.Meta = p.(*roleActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
