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
		m   = make(actionlog.Meta)
		str = func(i interface{}) string { return fmt.Sprintf("%v", i) }
	)

	// avoiding declared but not used
	_ = str

	if p.member != nil {
		m["member.handle"] = str(p.member.Handle)
		m["member.email"] = str(p.member.Email)
		m["member.name"] = str(p.member.Name)
		m["member.ID"] = str(p.member.ID)
	}
	if p.role != nil {
		m["role.handle"] = str(p.role.Handle)
		m["role.name"] = str(p.role.Name)
		m["role.ID"] = str(p.role.ID)
	}
	if p.new != nil {
		m["new.handle"] = str(p.new.Handle)
		m["new.name"] = str(p.new.Name)
		m["new.ID"] = str(p.new.ID)
	}
	if p.update != nil {
		m["update.handle"] = str(p.update.Handle)
		m["update.name"] = str(p.update.Name)
		m["update.ID"] = str(p.update.ID)
	}
	if p.existing != nil {
		m["existing.handle"] = str(p.existing.Handle)
		m["existing.name"] = str(p.existing.Name)
		m["existing.ID"] = str(p.existing.ID)
	}
	if p.target != nil {
		m["target.handle"] = str(p.target.Handle)
		m["target.name"] = str(p.target.Name)
		m["target.ID"] = str(p.target.ID)
	}
	if p.filter != nil {
		m["filter.query"] = str(p.filter.Query)
		m["filter.roleID"] = str(p.filter.RoleID)
		m["filter.memberID"] = str(p.filter.MemberID)
		m["filter.handle"] = str(p.filter.Handle)
		m["filter.name"] = str(p.filter.Name)
		m["filter.deleted"] = str(p.filter.Deleted)
		m["filter.archived"] = str(p.filter.Archived)
		m["filter.sort"] = str(p.filter.Sort)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p roleActionProps) tr(in string, err error) string {
	var pairs = []string{"{err}"}

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
		pairs = append(pairs, "{member}", fmt.Sprintf("%v", p.member.Handle))
		pairs = append(pairs, "{member.handle}", fmt.Sprintf("%v", p.member.Handle))
		pairs = append(pairs, "{member.email}", fmt.Sprintf("%v", p.member.Email))
		pairs = append(pairs, "{member.name}", fmt.Sprintf("%v", p.member.Name))
		pairs = append(pairs, "{member.ID}", fmt.Sprintf("%v", p.member.ID))
	}

	if p.role != nil {
		pairs = append(pairs, "{role}", fmt.Sprintf("%v", p.role.Handle))
		pairs = append(pairs, "{role.handle}", fmt.Sprintf("%v", p.role.Handle))
		pairs = append(pairs, "{role.name}", fmt.Sprintf("%v", p.role.Name))
		pairs = append(pairs, "{role.ID}", fmt.Sprintf("%v", p.role.ID))
	}

	if p.new != nil {
		pairs = append(pairs, "{new}", fmt.Sprintf("%v", p.new.Handle))
		pairs = append(pairs, "{new.handle}", fmt.Sprintf("%v", p.new.Handle))
		pairs = append(pairs, "{new.name}", fmt.Sprintf("%v", p.new.Name))
		pairs = append(pairs, "{new.ID}", fmt.Sprintf("%v", p.new.ID))
	}

	if p.update != nil {
		pairs = append(pairs, "{update}", fmt.Sprintf("%v", p.update.Handle))
		pairs = append(pairs, "{update.handle}", fmt.Sprintf("%v", p.update.Handle))
		pairs = append(pairs, "{update.name}", fmt.Sprintf("%v", p.update.Name))
		pairs = append(pairs, "{update.ID}", fmt.Sprintf("%v", p.update.ID))
	}

	if p.existing != nil {
		pairs = append(pairs, "{existing}", fmt.Sprintf("%v", p.existing.Handle))
		pairs = append(pairs, "{existing.handle}", fmt.Sprintf("%v", p.existing.Handle))
		pairs = append(pairs, "{existing.name}", fmt.Sprintf("%v", p.existing.Name))
		pairs = append(pairs, "{existing.ID}", fmt.Sprintf("%v", p.existing.ID))
	}

	if p.target != nil {
		pairs = append(pairs, "{target}", fmt.Sprintf("%v", p.target.Handle))
		pairs = append(pairs, "{target.handle}", fmt.Sprintf("%v", p.target.Handle))
		pairs = append(pairs, "{target.name}", fmt.Sprintf("%v", p.target.Name))
		pairs = append(pairs, "{target.ID}", fmt.Sprintf("%v", p.target.ID))
	}

	if p.filter != nil {
		pairs = append(pairs, "{filter}", fmt.Sprintf("%v", p.filter.Query))
		pairs = append(pairs, "{filter.query}", fmt.Sprintf("%v", p.filter.Query))
		pairs = append(pairs, "{filter.roleID}", fmt.Sprintf("%v", p.filter.RoleID))
		pairs = append(pairs, "{filter.memberID}", fmt.Sprintf("%v", p.filter.MemberID))
		pairs = append(pairs, "{filter.handle}", fmt.Sprintf("%v", p.filter.Handle))
		pairs = append(pairs, "{filter.name}", fmt.Sprintf("%v", p.filter.Name))
		pairs = append(pairs, "{filter.deleted}", fmt.Sprintf("%v", p.filter.Deleted))
		pairs = append(pairs, "{filter.archived}", fmt.Sprintf("%v", p.filter.Archived))
		pairs = append(pairs, "{filter.sort}", fmt.Sprintf("%v", p.filter.Sort))
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

// RoleErrNonexistent returns "system:role.nonexistent" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func RoleErrNonexistent(props ...*roleActionProps) *roleError {
	var e = &roleError{
		timestamp: time.Now(),
		resource:  "system:role",
		error:     "nonexistent",
		action:    "error",
		message:   "role does not exist",
		log:       "role does not exist",
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
		message:   "not allowed to read role",
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
		message:   "not allowed to create role",
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
		message:   "not allowed to update role",
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
		message:   "not allowed to delete role",
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
		message:   "not allowed to undelete role",
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
		message:   "not allowed to archive role",
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
		message:   "not allowed to unarchive role",
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

			// copy action to returning and recording error
			retError.action = action().action

			// we'll use RoleErrGeneric for recording too
			// because it can hold more info
			recError = retError
		} else if retError != nil {
			// copy action to returning and recording error
			retError.action = action().action
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
