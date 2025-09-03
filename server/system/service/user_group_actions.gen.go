package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/user_group_actions.yaml

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
	userGroupActionProps struct {
		member    *types.User
		userGroup *types.UserGroup
		new       *types.UserGroup
		update    *types.UserGroup
		existing  *types.UserGroup
		filter    *types.UserGroupFilter
	}

	userGroupAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *userGroupActionProps
	}

	userGroupLogMetaKey   struct{}
	userGroupPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setMember updates userGroupActionProps's member
//
// This function is auto-generated.
func (p *userGroupActionProps) setMember(member *types.User) *userGroupActionProps {
	p.member = member
	return p
}

// setUserGroup updates userGroupActionProps's userGroup
//
// This function is auto-generated.
func (p *userGroupActionProps) setUserGroup(userGroup *types.UserGroup) *userGroupActionProps {
	p.userGroup = userGroup
	return p
}

// setNew updates userGroupActionProps's new
//
// This function is auto-generated.
func (p *userGroupActionProps) setNew(new *types.UserGroup) *userGroupActionProps {
	p.new = new
	return p
}

// setUpdate updates userGroupActionProps's update
//
// This function is auto-generated.
func (p *userGroupActionProps) setUpdate(update *types.UserGroup) *userGroupActionProps {
	p.update = update
	return p
}

// setExisting updates userGroupActionProps's existing
//
// This function is auto-generated.
func (p *userGroupActionProps) setExisting(existing *types.UserGroup) *userGroupActionProps {
	p.existing = existing
	return p
}

// setFilter updates userGroupActionProps's filter
//
// This function is auto-generated.
func (p *userGroupActionProps) setFilter(filter *types.UserGroupFilter) *userGroupActionProps {
	p.filter = filter
	return p
}

// Serialize converts userGroupActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p userGroupActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.member != nil {
		m.Set("member.handle", p.member.Handle, true)
		m.Set("member.email", p.member.Email, true)
		m.Set("member.ID", p.member.ID, true)
	}
	if p.userGroup != nil {
		m.Set("userGroup.handle", p.userGroup.Handle, true)
		m.Set("userGroup.ID", p.userGroup.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.existing != nil {
		m.Set("existing.handle", p.existing.Handle, true)
		m.Set("existing.ID", p.existing.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.userGroupID", p.filter.UserGroupID, true)
		m.Set("filter.handle", p.filter.Handle, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p userGroupActionProps) Format(in string, err error) string {
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
				p.member.ID,
			),
		)
		pairs = append(pairs, "{{member.handle}}", fns(p.member.Handle))
		pairs = append(pairs, "{{member.email}}", fns(p.member.Email))
		pairs = append(pairs, "{{member.ID}}", fns(p.member.ID))
	}

	if p.userGroup != nil {
		// replacement for "{{userGroup}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{userGroup}}",
			fns(
				p.userGroup.Handle,
				p.userGroup.ID,
			),
		)
		pairs = append(pairs, "{{userGroup.handle}}", fns(p.userGroup.Handle))
		pairs = append(pairs, "{{userGroup.ID}}", fns(p.userGroup.ID))
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

	if p.existing != nil {
		// replacement for "{{existing}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{existing}}",
			fns(
				p.existing.Handle,
				p.existing.ID,
			),
		)
		pairs = append(pairs, "{{existing.handle}}", fns(p.existing.Handle))
		pairs = append(pairs, "{{existing.ID}}", fns(p.existing.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.UserGroupID,
				p.filter.Handle,
				p.filter.Deleted,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{{filter.userGroupID}}", fns(p.filter.UserGroupID))
		pairs = append(pairs, "{{filter.handle}}", fns(p.filter.Handle))
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
func (a *userGroupAction) String() string {
	var props = &userGroupActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *userGroupAction) ToAction() *actionlog.Action {
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

// UserGroupActionSearch returns "system:user-group.search" action
//
// This function is auto-generated.
func UserGroupActionSearch(props ...*userGroupActionProps) *userGroupAction {
	a := &userGroupAction{
		timestamp: time.Now(),
		resource:  "system:user-group",
		action:    "search",
		log:       "searched for matching userGroups",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserGroupActionLookup returns "system:user-group.lookup" action
//
// This function is auto-generated.
func UserGroupActionLookup(props ...*userGroupActionProps) *userGroupAction {
	a := &userGroupAction{
		timestamp: time.Now(),
		resource:  "system:user-group",
		action:    "lookup",
		log:       "looked-up for a {{userGroup}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserGroupActionCreate returns "system:user-group.create" action
//
// This function is auto-generated.
func UserGroupActionCreate(props ...*userGroupActionProps) *userGroupAction {
	a := &userGroupAction{
		timestamp: time.Now(),
		resource:  "system:user-group",
		action:    "create",
		log:       "created {{userGroup}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserGroupActionUpdate returns "system:user-group.update" action
//
// This function is auto-generated.
func UserGroupActionUpdate(props ...*userGroupActionProps) *userGroupAction {
	a := &userGroupAction{
		timestamp: time.Now(),
		resource:  "system:user-group",
		action:    "update",
		log:       "updated {{userGroup}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserGroupActionDelete returns "system:user-group.delete" action
//
// This function is auto-generated.
func UserGroupActionDelete(props ...*userGroupActionProps) *userGroupAction {
	a := &userGroupAction{
		timestamp: time.Now(),
		resource:  "system:user-group",
		action:    "delete",
		log:       "deleted {{userGroup}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserGroupActionUndelete returns "system:user-group.undelete" action
//
// This function is auto-generated.
func UserGroupActionUndelete(props ...*userGroupActionProps) *userGroupAction {
	a := &userGroupAction{
		timestamp: time.Now(),
		resource:  "system:user-group",
		action:    "undelete",
		log:       "undeleted {{userGroup}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserGroupActionMembers returns "system:user-group.members" action
//
// This function is auto-generated.
func UserGroupActionMembers(props ...*userGroupActionProps) *userGroupAction {
	a := &userGroupAction{
		timestamp: time.Now(),
		resource:  "system:user-group",
		action:    "members",
		log:       "searched for members of {{userGroup}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserGroupActionMemberAdd returns "system:user-group.memberAdd" action
//
// This function is auto-generated.
func UserGroupActionMemberAdd(props ...*userGroupActionProps) *userGroupAction {
	a := &userGroupAction{
		timestamp: time.Now(),
		resource:  "system:user-group",
		action:    "memberAdd",
		log:       "added {{member.email}} to {{userGroup}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserGroupActionMemberRemove returns "system:user-group.memberRemove" action
//
// This function is auto-generated.
func UserGroupActionMemberRemove(props ...*userGroupActionProps) *userGroupAction {
	a := &userGroupAction{
		timestamp: time.Now(),
		resource:  "system:user-group",
		action:    "memberRemove",
		log:       "removed {{member.email}} from {{userGroup}}",
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

// UserGroupErrGeneric returns "system:user-group.generic" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrGeneric(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "{err}"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNotFound returns "system:user-group.notFound" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNotFound(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("userGroup not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:user-group"),

		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrInvalidID returns "system:user-group.invalidID" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrInvalidID(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:user-group"),

		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrInvalidHandle returns "system:user-group.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrInvalidHandle(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:user-group"),

		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrMissingSelfID returns "system:user-group.missingSelfID" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrMissingSelfID(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("missing self ID", nil),

		errors.Meta("type", "missingSelfID"),
		errors.Meta("resource", "system:user-group"),

		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.missingSelfID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrInvalidSelfID returns "system:user-group.invalidSelfID" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrInvalidSelfID(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid self ID", nil),

		errors.Meta("type", "invalidSelfID"),
		errors.Meta("resource", "system:user-group"),

		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.invalidSelfID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrInvalidUpdateStructure returns "system:user-group.invalidUpdateStructure" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrInvalidUpdateStructure(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid structure after update", nil),

		errors.Meta("type", "invalidUpdateStructure"),
		errors.Meta("resource", "system:user-group"),

		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.invalidUpdateStructure"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrStaleData returns "system:user-group.staleData" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrStaleData(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:user-group"),

		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNotAllowedToRead returns "system:user-group.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNotAllowedToRead(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this userGroup", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "failed to read {{userGroup.handle}}; insufficient permissions"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNotAllowedToSearch returns "system:user-group.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNotAllowedToSearch(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list or search userGroups", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "failed to search for userGroups; insufficient permissions"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNotAllowedToListUserGroups returns "system:user-group.notAllowedToListUserGroups" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNotAllowedToListUserGroups(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list userGroups", nil),

		errors.Meta("type", "notAllowedToListUserGroups"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "failed to list userGroup; insufficient permissions"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.notAllowedToListUserGroups"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNotAllowedToCreate returns "system:user-group.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNotAllowedToCreate(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create userGroups", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "failed to create userGroups; insufficient permissions"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNotAllowedToUpdate returns "system:user-group.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNotAllowedToUpdate(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this userGroup", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "failed to update {{userGroup.handle}}; insufficient permissions"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNotAllowedToDelete returns "system:user-group.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNotAllowedToDelete(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this userGroup", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "failed to delete {{userGroup.handle}}; insufficient permissions"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNotAllowedToUndelete returns "system:user-group.notAllowedToUndelete" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNotAllowedToUndelete(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to undelete this userGroup", nil),

		errors.Meta("type", "notAllowedToUndelete"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "failed to undelete {{userGroup.handle}}; insufficient permissions"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.notAllowedToUndelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrHandleNotUnique returns "system:user-group.handleNotUnique" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrHandleNotUnique(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handle not unique", nil),

		errors.Meta("type", "handleNotUnique"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "used duplicate handle ({{userGroup.handle}}) for userGroup"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.handleNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNameNotUnique returns "system:user-group.nameNotUnique" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNameNotUnique(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("name not unique", nil),

		errors.Meta("type", "nameNotUnique"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "used duplicate name ({{userGroup.handle}}) for userGroup"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.nameNotUnique"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// UserGroupErrNotAllowedToManageMembers returns "system:user-group.notAllowedToManageMembers" as *errors.Error
//
// This function is auto-generated.
func UserGroupErrNotAllowedToManageMembers(mm ...*userGroupActionProps) *errors.Error {
	var p = &userGroupActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage user group members", nil),

		errors.Meta("type", "notAllowedToManageMembers"),
		errors.Meta("resource", "system:user-group"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(userGroupLogMetaKey{}, "failed to manage {{userGroup.handle}} members; insufficient permissions"),
		errors.Meta(userGroupPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "user-group.errors.notAllowedToManageMembers"),

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
func (svc userGroup) recordAction(ctx context.Context, props *userGroupActionProps, actionFn func(...*userGroupActionProps) *userGroupAction, err error) error {
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
		a.Description = props.Format(m.AsString(userGroupLogMetaKey{}), err)

		if p, has := m[userGroupPropsMetaKey{}]; has {
			a.Meta = p.(*userGroupActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
