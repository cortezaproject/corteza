package service

// This file is auto-generated from system/service/user_actions.yaml
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
	userActionProps struct {
		user     *types.User
		new      *types.User
		update   *types.User
		existing *types.User
		filter   *types.UserFilter
	}

	userAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *userActionProps
	}

	userError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *userActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setUser updates userActionProps's user
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *userActionProps) setUser(user *types.User) *userActionProps {
	p.user = user
	return p
}

// setNew updates userActionProps's new
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *userActionProps) setNew(new *types.User) *userActionProps {
	p.new = new
	return p
}

// setUpdate updates userActionProps's update
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *userActionProps) setUpdate(update *types.User) *userActionProps {
	p.update = update
	return p
}

// setExisting updates userActionProps's existing
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *userActionProps) setExisting(existing *types.User) *userActionProps {
	p.existing = existing
	return p
}

// setFilter updates userActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *userActionProps) setFilter(filter *types.UserFilter) *userActionProps {
	p.filter = filter
	return p
}

// serialize converts userActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p userActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.user != nil {
		m.Set("user.handle", p.user.Handle, true)
		m.Set("user.email", p.user.Email, true)
		m.Set("user.name", p.user.Name, true)
		m.Set("user.username", p.user.Username, true)
		m.Set("user.ID", p.user.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.email", p.new.Email, true)
		m.Set("new.name", p.new.Name, true)
		m.Set("new.username", p.new.Username, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.email", p.update.Email, true)
		m.Set("update.name", p.update.Name, true)
		m.Set("update.username", p.update.Username, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.existing != nil {
		m.Set("existing.handle", p.existing.Handle, true)
		m.Set("existing.email", p.existing.Email, true)
		m.Set("existing.name", p.existing.Name, true)
		m.Set("existing.username", p.existing.Username, true)
		m.Set("existing.ID", p.existing.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.query", p.filter.Query, true)
		m.Set("filter.userID", p.filter.UserID, true)
		m.Set("filter.roleID", p.filter.RoleID, true)
		m.Set("filter.handle", p.filter.Handle, true)
		m.Set("filter.email", p.filter.Email, true)
		m.Set("filter.username", p.filter.Username, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
		m.Set("filter.suspended", p.filter.Suspended, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p userActionProps) tr(in string, err error) string {
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

	if p.user != nil {
		// replacement for "{user}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{user}",
			fns(
				p.user.Handle,
				p.user.Email,
				p.user.Name,
				p.user.Username,
				p.user.ID,
			),
		)
		pairs = append(pairs, "{user.handle}", fns(p.user.Handle))
		pairs = append(pairs, "{user.email}", fns(p.user.Email))
		pairs = append(pairs, "{user.name}", fns(p.user.Name))
		pairs = append(pairs, "{user.username}", fns(p.user.Username))
		pairs = append(pairs, "{user.ID}", fns(p.user.ID))
	}

	if p.new != nil {
		// replacement for "{new}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{new}",
			fns(
				p.new.Handle,
				p.new.Email,
				p.new.Name,
				p.new.Username,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{new.handle}", fns(p.new.Handle))
		pairs = append(pairs, "{new.email}", fns(p.new.Email))
		pairs = append(pairs, "{new.name}", fns(p.new.Name))
		pairs = append(pairs, "{new.username}", fns(p.new.Username))
		pairs = append(pairs, "{new.ID}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{update}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{update}",
			fns(
				p.update.Handle,
				p.update.Email,
				p.update.Name,
				p.update.Username,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{update.handle}", fns(p.update.Handle))
		pairs = append(pairs, "{update.email}", fns(p.update.Email))
		pairs = append(pairs, "{update.name}", fns(p.update.Name))
		pairs = append(pairs, "{update.username}", fns(p.update.Username))
		pairs = append(pairs, "{update.ID}", fns(p.update.ID))
	}

	if p.existing != nil {
		// replacement for "{existing}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{existing}",
			fns(
				p.existing.Handle,
				p.existing.Email,
				p.existing.Name,
				p.existing.Username,
				p.existing.ID,
			),
		)
		pairs = append(pairs, "{existing.handle}", fns(p.existing.Handle))
		pairs = append(pairs, "{existing.email}", fns(p.existing.Email))
		pairs = append(pairs, "{existing.name}", fns(p.existing.Name))
		pairs = append(pairs, "{existing.username}", fns(p.existing.Username))
		pairs = append(pairs, "{existing.ID}", fns(p.existing.ID))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Query,
				p.filter.UserID,
				p.filter.RoleID,
				p.filter.Handle,
				p.filter.Email,
				p.filter.Username,
				p.filter.Deleted,
				p.filter.Suspended,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{filter.query}", fns(p.filter.Query))
		pairs = append(pairs, "{filter.userID}", fns(p.filter.UserID))
		pairs = append(pairs, "{filter.roleID}", fns(p.filter.RoleID))
		pairs = append(pairs, "{filter.handle}", fns(p.filter.Handle))
		pairs = append(pairs, "{filter.email}", fns(p.filter.Email))
		pairs = append(pairs, "{filter.username}", fns(p.filter.Username))
		pairs = append(pairs, "{filter.deleted}", fns(p.filter.Deleted))
		pairs = append(pairs, "{filter.suspended}", fns(p.filter.Suspended))
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
func (a *userAction) String() string {
	var props = &userActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *userAction) LoggableAction() *actionlog.Action {
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
func (e *userError) String() string {
	var props = &userActionProps{}

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
func (e *userError) Error() string {
	var props = &userActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *userError) Is(Resource error) bool {
	t, ok := Resource.(*userError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps userError around another error
//
// This function is auto-generated.
//
func (e *userError) Wrap(err error) *userError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *userError) Unwrap() error {
	return e.wrap
}

func (e *userError) LoggableAction() *actionlog.Action {
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

// UserActionSearch returns "system:user.search" error
//
// This function is auto-generated.
//
func UserActionSearch(props ...*userActionProps) *userAction {
	a := &userAction{
		timestamp: time.Now(),
		resource:  "system:user",
		action:    "search",
		log:       "searched for matching users",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserActionLookup returns "system:user.lookup" error
//
// This function is auto-generated.
//
func UserActionLookup(props ...*userActionProps) *userAction {
	a := &userAction{
		timestamp: time.Now(),
		resource:  "system:user",
		action:    "lookup",
		log:       "looked-up for a {user}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserActionCreate returns "system:user.create" error
//
// This function is auto-generated.
//
func UserActionCreate(props ...*userActionProps) *userAction {
	a := &userAction{
		timestamp: time.Now(),
		resource:  "system:user",
		action:    "create",
		log:       "created {user}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserActionUpdate returns "system:user.update" error
//
// This function is auto-generated.
//
func UserActionUpdate(props ...*userActionProps) *userAction {
	a := &userAction{
		timestamp: time.Now(),
		resource:  "system:user",
		action:    "update",
		log:       "updated {user}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserActionDelete returns "system:user.delete" error
//
// This function is auto-generated.
//
func UserActionDelete(props ...*userActionProps) *userAction {
	a := &userAction{
		timestamp: time.Now(),
		resource:  "system:user",
		action:    "delete",
		log:       "deleted {user}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserActionUndelete returns "system:user.undelete" error
//
// This function is auto-generated.
//
func UserActionUndelete(props ...*userActionProps) *userAction {
	a := &userAction{
		timestamp: time.Now(),
		resource:  "system:user",
		action:    "undelete",
		log:       "undeleted {user}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserActionSuspend returns "system:user.suspend" error
//
// This function is auto-generated.
//
func UserActionSuspend(props ...*userActionProps) *userAction {
	a := &userAction{
		timestamp: time.Now(),
		resource:  "system:user",
		action:    "suspend",
		log:       "suspended {user}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserActionUnsuspend returns "system:user.unsuspend" error
//
// This function is auto-generated.
//
func UserActionUnsuspend(props ...*userActionProps) *userAction {
	a := &userAction{
		timestamp: time.Now(),
		resource:  "system:user",
		action:    "unsuspend",
		log:       "unsuspended {user}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// UserActionSetPassword returns "system:user.setPassword" error
//
// This function is auto-generated.
//
func UserActionSetPassword(props ...*userActionProps) *userAction {
	a := &userAction{
		timestamp: time.Now(),
		resource:  "system:user",
		action:    "setPassword",
		log:       "password changed for {user}",
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

// UserErrGeneric returns "system:user.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func UserErrGeneric(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *userActionProps {
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

// UserErrNotFound returns "system:user.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func UserErrNotFound(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "notFound",
		action:    "error",
		message:   "user not found",
		log:       "user not found",
		severity:  actionlog.Warning,
		props: func() *userActionProps {
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

// UserErrInvalidID returns "system:user.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func UserErrInvalidID(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *userActionProps {
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

// UserErrInvalidHandle returns "system:user.invalidHandle" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func UserErrInvalidHandle(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "invalidHandle",
		action:    "error",
		message:   "invalid handle",
		log:       "invalid handle",
		severity:  actionlog.Warning,
		props: func() *userActionProps {
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

// UserErrInvalidEmail returns "system:user.invalidEmail" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func UserErrInvalidEmail(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "invalidEmail",
		action:    "error",
		message:   "invalid email",
		log:       "invalid email",
		severity:  actionlog.Warning,
		props: func() *userActionProps {
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

// UserErrNotAllowedToRead returns "system:user.notAllowedToRead" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func UserErrNotAllowedToRead(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "notAllowedToRead",
		action:    "error",
		message:   "not allowed to read this user",
		log:       "failed to read {user.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *userActionProps {
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

// UserErrNotAllowedToListUsers returns "system:user.notAllowedToListUsers" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func UserErrNotAllowedToListUsers(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "notAllowedToListUsers",
		action:    "error",
		message:   "not allowed to list users",
		log:       "failed to list user; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *userActionProps {
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

// UserErrNotAllowedToCreate returns "system:user.notAllowedToCreate" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func UserErrNotAllowedToCreate(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create users",
		log:       "failed to create users; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *userActionProps {
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

// UserErrNotAllowedToUpdate returns "system:user.notAllowedToUpdate" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func UserErrNotAllowedToUpdate(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "notAllowedToUpdate",
		action:    "error",
		message:   "not allowed to update this user",
		log:       "failed to update {user.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *userActionProps {
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

// UserErrNotAllowedToDelete returns "system:user.notAllowedToDelete" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func UserErrNotAllowedToDelete(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "notAllowedToDelete",
		action:    "error",
		message:   "not allowed to delete this user",
		log:       "failed to delete {user.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *userActionProps {
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

// UserErrNotAllowedToUndelete returns "system:user.notAllowedToUndelete" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func UserErrNotAllowedToUndelete(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "notAllowedToUndelete",
		action:    "error",
		message:   "not allowed to undelete this user",
		log:       "failed to undelete {user.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *userActionProps {
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

// UserErrNotAllowedToSuspend returns "system:user.notAllowedToSuspend" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func UserErrNotAllowedToSuspend(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "notAllowedToSuspend",
		action:    "error",
		message:   "not allowed to suspend this user",
		log:       "failed to suspend {user.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *userActionProps {
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

// UserErrNotAllowedToUnsuspend returns "system:user.notAllowedToUnsuspend" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func UserErrNotAllowedToUnsuspend(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "notAllowedToUnsuspend",
		action:    "error",
		message:   "not allowed to unsuspend this user",
		log:       "failed to unsuspend {user.handle}; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *userActionProps {
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

// UserErrHandleNotUnique returns "system:user.handleNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func UserErrHandleNotUnique(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "handleNotUnique",
		action:    "error",
		message:   "handle not unique",
		log:       "used duplicate handle ({user.handle}) for user",
		severity:  actionlog.Warning,
		props: func() *userActionProps {
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

// UserErrEmailNotUnique returns "system:user.emailNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func UserErrEmailNotUnique(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "emailNotUnique",
		action:    "error",
		message:   "email not unique",
		log:       "used duplicate email ({user.email}) for user",
		severity:  actionlog.Warning,
		props: func() *userActionProps {
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

// UserErrUsernameNotUnique returns "system:user.usernameNotUnique" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func UserErrUsernameNotUnique(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "usernameNotUnique",
		action:    "error",
		message:   "username not unique",
		log:       "used duplicate username ({user.username}) for user",
		severity:  actionlog.Warning,
		props: func() *userActionProps {
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

// UserErrPasswordNotSecure returns "system:user.passwordNotSecure" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func UserErrPasswordNotSecure(props ...*userActionProps) *userError {
	var e = &userError{
		timestamp: time.Now(),
		resource:  "system:user",
		error:     "passwordNotSecure",
		action:    "error",
		message:   "provided password is not secure; use longer password with more non-alphanumeric character",
		log:       "provided password is not secure; use longer password with more non-alphanumeric character",
		severity:  actionlog.Alert,
		props: func() *userActionProps {
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
// action (optional) fn will be used to construct userAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc user) recordAction(ctx context.Context, props *userActionProps, action func(...*userActionProps) *userAction, err error) error {
	var (
		ok bool

		// Return error
		retError *userError

		// Recorder error
		recError *userError
	)

	if err != nil {
		if retError, ok = err.(*userError); !ok {
			// got non-user error, wrap it with UserErrGeneric
			retError = UserErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use UserErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type userError
				if unwrappedSinkError, ok := unwrappedError.(*userError); ok {
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
