package service

// This file is auto-generated from system/service/role_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/system/types"
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

	roleError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *roleActionProps
	}
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
// Allows method chaining
//
// This function is auto-generated.
//
func (p *roleActionProps) setMember(member *types.User) *roleActionProps {
	p.member = member
	return p
}

// setRole updates roleActionProps's role
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *roleActionProps) setRole(role *types.Role) *roleActionProps {
	p.role = role
	return p
}

// setNew updates roleActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *roleActionProps) setNew(new *types.Role) *roleActionProps {
	p.new = new
	return p
}

// setUpdate updates roleActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *roleActionProps) setUpdate(update *types.Role) *roleActionProps {
	p.update = update
	return p
}

// setExisting updates roleActionProps's existing
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *roleActionProps) setExisting(existing *types.Role) *roleActionProps {
	p.existing = existing
	return p
}

// setTarget updates roleActionProps's target
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *roleActionProps) setTarget(target *types.Role) *roleActionProps {
	p.target = target
	return p
}

// setFilter updates roleActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *roleActionProps) setFilter(filter *types.RoleFilter) *roleActionProps {
	p.filter = filter
	return p
}

// serialize converts roleActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p roleActionProps) serialize() actionlog.Meta {
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
//
func (p roleActionProps) tr(in string, err error) string {
	var (
		pairs = []string{"{err}"}
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
		for {
			// Unwrap errors
			ue := errors.Unwrap(err)
			if ue == nil {
				break
			}

			err = ue
		}

		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}

	if p.member != nil {
		// replacement for "{member}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{member}",
			fns(
				p.member.Handle,
				p.member.Email,
				p.member.Name,
				p.member.ID,
			),
		)
		pairs = append(pairs, "{member.handle}", fns(p.member.Handle))
		pairs = append(pairs, "{member.email}", fns(p.member.Email))
		pairs = append(pairs, "{member.name}", fns(p.member.Name))
		pairs = append(pairs, "{member.ID}", fns(p.member.ID))
	}

	if p.role != nil {
		// replacement for "{role}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{role}",
			fns(
				p.role.Handle,
				p.role.Name,
				p.role.ID,
			),
		)
		pairs = append(pairs, "{role.handle}", fns(p.role.Handle))
		pairs = append(pairs, "{role.name}", fns(p.role.Name))
		pairs = append(pairs, "{role.ID}", fns(p.role.ID))
	}

	if p.new != nil {
		// replacement for "{new}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{new}",
			fns(
				p.new.Handle,
				p.new.Name,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{new.handle}", fns(p.new.Handle))
		pairs = append(pairs, "{new.name}", fns(p.new.Name))
		pairs = append(pairs, "{new.ID}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{update}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{update}",
			fns(
				p.update.Handle,
				p.update.Name,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{update.handle}", fns(p.update.Handle))
		pairs = append(pairs, "{update.name}", fns(p.update.Name))
		pairs = append(pairs, "{update.ID}", fns(p.update.ID))
	}

	if p.existing != nil {
		// replacement for "{existing}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{existing}",
			fns(
				p.existing.Handle,
				p.existing.Name,
				p.existing.ID,
			),
		)
		pairs = append(pairs, "{existing.handle}", fns(p.existing.Handle))
		pairs = append(pairs, "{existing.name}", fns(p.existing.Name))
		pairs = append(pairs, "{existing.ID}", fns(p.existing.ID))
	}

	if p.target != nil {
		// replacement for "{target}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{target}",
			fns(
				p.target.Handle,
				p.target.Name,
				p.target.ID,
			),
		)
		pairs = append(pairs, "{target.handle}", fns(p.target.Handle))
		pairs = append(pairs, "{target.name}", fns(p.target.Name))
		pairs = append(pairs, "{target.ID}", fns(p.target.ID))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
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
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.roleID}", fns(p.filter.RoleID))
		pairs = append(pairs, "{filter.memberID}", fns(p.filter.MemberID))
		pairs = append(pairs, "{filter.handle}", fns(p.filter.Handle))
		pairs = append(pairs, "{filter.name}", fns(p.filter.Name))
		pairs = append(pairs, "{filter.deleted}", fns(p.filter.Deleted))
		pairs = append(pairs, "{filter.archived}", fns(p.filter.Archived))
		pairs = append(pairs, "{filter.sort}", fns(p.filter.Sort))
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
func (a *roleAction) String() string {
	var props = &roleActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *roleAction) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error methods

// String returns loggable description as string
//
// It falls back to message if log is not set
//
// This function is auto-generated.
//
func (e *roleError) String() string {
	var props = &roleActionProps{}

	if e.props != nil {
		props = e.props
	}

	if e.wrap != nil && !strings.Contains(e.log, "{err}") {
		// Suffix error log with {err} to ensure
		// we log the cause for this error
		e.log += ": {err}"
	}

	return props.tr(e.log, e.wrap)
}

// Error satisfies
//
// This function is auto-generated.
//
func (e *roleError) Error() string {
	var props = &roleActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *roleError) Is(Resource error) bool {
	t, ok := Resource.(*roleError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps roleError around another error
//
// This function is auto-generated.
//
func (e *roleError) Wrap(err error) *roleError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *roleError) Unwrap() error {
	return e.wrap
}

func (e *roleError) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Error:       e.Error(),
		Meta:        e.props.serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// RoleActionSearch returns "system:role.search" error
//
// This function is auto-generated.
//
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

// RoleActionLookup returns "system:role.lookup" error
//
// This function is auto-generated.
//
func RoleActionLookup(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "lookup",
		log:       "looked-up for a {role}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionCreate returns "system:role.create" error
//
// This function is auto-generated.
//
func RoleActionCreate(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "create",
		log:       "created {role}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionUpdate returns "system:role.update" error
//
// This function is auto-generated.
//
func RoleActionUpdate(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "update",
		log:       "updated {role}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionDelete returns "system:role.delete" error
//
// This function is auto-generated.
//
func RoleActionDelete(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "delete",
		log:       "deleted {role}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionUndelete returns "system:role.undelete" error
//
// This function is auto-generated.
//
func RoleActionUndelete(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "undelete",
		log:       "undeleted {role}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionArchive returns "system:role.archive" error
//
// This function is auto-generated.
//
func RoleActionArchive(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "archive",
		log:       "archived {role}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionUnarchive returns "system:role.unarchive" error
//
// This function is auto-generated.
//
func RoleActionUnarchive(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "unarchive",
		log:       "unarchived {role}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionMerge returns "system:role.merge" error
//
// This function is auto-generated.
//
func RoleActionMerge(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "merge",
		log:       "merged {target} with {role}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionMembers returns "system:role.members" error
//
// This function is auto-generated.
//
func RoleActionMembers(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "members",
		log:       "searched for members of {role}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionMemberAdd returns "system:role.memberAdd" error
//
// This function is auto-generated.
//
func RoleActionMemberAdd(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "memberAdd",
		log:       "added {member.email} to {role}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// RoleActionMemberRemove returns "system:role.memberRemove" error
//
// This function is auto-generated.
//
func RoleActionMemberRemove(props ...*roleActionProps) *roleAction {
	a := &roleAction{
		timestamp: time.Now(),
		resource:  "system:role",
		action:    "memberRemove",
		log:       "removed {member.email} from {role}",
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

// RoleErrGeneric returns "system:role.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func RoleErrGeneric(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotFOund returns "system:role.notFOund" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RoleErrNotFOund(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notFOund",
		action:    "error",
		message:   "role not found",
		log:       "role not found",
		severity:  actionlog.Warning,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrInvalidID returns "system:role.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RoleErrInvalidID(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrInvalidHandle returns "system:role.invalidHandle" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RoleErrInvalidHandle(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "invalidHandle",
		action:    "error",
		message:   "invalid handle",
		log:       "invalid handle",
		severity:  actionlog.Warning,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotAllowedToRead returns "system:role.notAllowedToRead" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func RoleErrNotAllowedToRead(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this role",
		log:       "failed to read {role.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotAllowedToListRoles returns "system:role.notAllowedToListRoles" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func RoleErrNotAllowedToListRoles(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notAllowedToListRoles",
		action:    "error",
		message:   "not allowed to list roles",
		log:       "failed to list role; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotAllowedToCreate returns "system:role.notAllowedToCreate" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func RoleErrNotAllowedToCreate(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create roles",
		log:       "failed to create role; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotAllowedToUpdate returns "system:role.notAllowedToUpdate" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func RoleErrNotAllowedToUpdate(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this role",
		log:       "failed to update {role.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotAllowedToDelete returns "system:role.notAllowedToDelete" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func RoleErrNotAllowedToDelete(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this role",
		log:       "failed to delete {role.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotAllowedToUndelete returns "system:role.notAllowedToUndelete" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func RoleErrNotAllowedToUndelete(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this role",
		log:       "failed to undelete {role.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotAllowedToArchive returns "system:role.notAllowedToArchive" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func RoleErrNotAllowedToArchive(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notAllowedToArchive",
		action:    "error",
		message:   "not allowed to archive this role",
		log:       "failed to archive {role.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotAllowedToUnarchive returns "system:role.notAllowedToUnarchive" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func RoleErrNotAllowedToUnarchive(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notAllowedToUnarchive",
		action:    "error",
		message:   "not allowed to unarchive this role",
		log:       "failed to unarchive {role.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNotAllowedToManageMembers returns "system:role.notAllowedToManageMembers" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func RoleErrNotAllowedToManageMembers(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "notAllowedToManageMembers",
		action:    "error",
		message:   "not allowed to manage role members",
		log:       "failed to manage {role.handle} members; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrHandleNotUnique returns "system:role.handleNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RoleErrHandleNotUnique(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "handleNotUnique",
		action:    "error",
		message:   "role handle not unique",
		log:       "used duplicate handle ({role.handle}) for role",
		severity:  actionlog.Warning,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// RoleErrNameNotUnique returns "system:role.nameNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RoleErrNameNotUnique(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "nameNotUnique",
		action:    "error",
		message:   "role name not unique",
		log:       "used duplicate name ({role.name}) for role",
		severity:  actionlog.Warning,
		props: func() *roleActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// context is used to enrich audit log entry with current user info, request ID, IP address...
// props are collected action/error properties
// action (optional) fn will be used to construct roleAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc role) recordAction(ctx context.Context, props *roleActionProps, action func(...*roleActionProps) *roleAction, err error) error {
	var (
		ok bool

		// Return error
		retError *roleError

		// Recorder error
		recError *roleError
	)

	if err != nil {
		if retError, ok = err.(*roleError); !ok {
			// got non-role error, wrap it with RoleErrGeneric
			retError = RoleErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use RoleErrGeneric for recording too
			// because it can hold more info
			recError = retError
		} else if retError != nil {
			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}
			// start with copy of return error for recording
			// this will be updated with tha root cause as we try and
			// unwrap the error
			recError = retError

			// find the original recError for this error
			// for the purpose of logging
			var unwrappedError error = retError
			for {
				if unwrappedError = errors.Unwrap(unwrappedError); unwrappedError == nil {
					// nothing wrapped
					break
				}

				// update recError ONLY of wrapped error is of type roleError
				if unwrappedSinkError, ok := unwrappedError.(*roleError); ok {
					recError = unwrappedSinkError
				}
			}

			if retError.props == nil {
				// set props on returning error if empty
				retError.props = props
			}

			if recError.props == nil {
				// set props on recording error if empty
				recError.props = props
			}
		}
	}

	if svc.actionlog != nil {
		if retError != nil {
			// failed action, log error
			svc.actionlog.Record(ctx, recError)
		} else if action != nil {
			// successful
			svc.actionlog.Record(ctx, action(props))
		}
	}

	if err == nil {
		// retError not an interface and that WILL (!!) cause issues
		// with nil check (== nil) when it is not explicitly returned
		return nil
	}

	return retError
}
