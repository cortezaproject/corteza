package service

// This file is auto-generated from system/service/auth_actions.yaml
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
	authActionProps struct {
		email       string
		provider    string
		credentials *types.Credentials
		role        *types.Role
		user        *types.User
	}

	authAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *authActionProps
	}

	authError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *authActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setEmail updates authActionProps's email
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *authActionProps) setEmail(email string) *authActionProps {
	p.email = email
	return p
}

// setProvider updates authActionProps's provider
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *authActionProps) setProvider(provider string) *authActionProps {
	p.provider = provider
	return p
}

// setCredentials updates authActionProps's credentials
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *authActionProps) setCredentials(credentials *types.Credentials) *authActionProps {
	p.credentials = credentials
	return p
}

// setRole updates authActionProps's role
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *authActionProps) setRole(role *types.Role) *authActionProps {
	p.role = role
	return p
}

// setUser updates authActionProps's user
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *authActionProps) setUser(user *types.User) *authActionProps {
	p.user = user
	return p
}

// serialize converts authActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p authActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	m.Set("email", p.email, true)
	m.Set("provider", p.provider, true)
	if p.credentials != nil {
		m.Set("credentials.kind", p.credentials.Kind, true)
		m.Set("credentials.label", p.credentials.Label, true)
		m.Set("credentials.ID", p.credentials.ID, true)
	}
	if p.role != nil {
		m.Set("role.handle", p.role.Handle, true)
		m.Set("role.name", p.role.Name, true)
		m.Set("role.ID", p.role.ID, true)
	}
	if p.user != nil {
		m.Set("user.handle", p.user.Handle, true)
		m.Set("user.name", p.user.Name, true)
		m.Set("user.ID", p.user.ID, true)
		m.Set("user.email", p.user.Email, true)
		m.Set("user.suspendedAt", p.user.SuspendedAt, true)
		m.Set("user.deletedAt", p.user.DeletedAt, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p authActionProps) tr(in string, err error) string {
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
	pairs = append(pairs, "{email}", fns(p.email))
	pairs = append(pairs, "{provider}", fns(p.provider))

	if p.credentials != nil {
		// replacement for "{credentials}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{credentials}",
			fns(
				p.credentials.Kind,
				p.credentials.Label,
				p.credentials.ID,
			),
		)
		pairs = append(pairs, "{credentials.kind}", fns(p.credentials.Kind))
		pairs = append(pairs, "{credentials.label}", fns(p.credentials.Label))
		pairs = append(pairs, "{credentials.ID}", fns(p.credentials.ID))
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

	if p.user != nil {
		// replacement for "{user}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{user}",
			fns(
				p.user.Handle,
				p.user.Name,
				p.user.ID,
				p.user.Email,
				p.user.SuspendedAt,
				p.user.DeletedAt,
			),
		)
		pairs = append(pairs, "{user.handle}", fns(p.user.Handle))
		pairs = append(pairs, "{user.name}", fns(p.user.Name))
		pairs = append(pairs, "{user.ID}", fns(p.user.ID))
		pairs = append(pairs, "{user.email}", fns(p.user.Email))
		pairs = append(pairs, "{user.suspendedAt}", fns(p.user.SuspendedAt))
		pairs = append(pairs, "{user.deletedAt}", fns(p.user.DeletedAt))
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
func (a *authAction) String() string {
	var props = &authActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *authAction) LoggableAction() *actionlog.Action {
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
func (e *authError) String() string {
	var props = &authActionProps{}

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
func (e *authError) Error() string {
	var props = &authActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *authError) Is(Resource error) bool {
	t, ok := Resource.(*authError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps authError around another error
//
// This function is auto-generated.
//
func (e *authError) Wrap(err error) *authError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *authError) Unwrap() error {
	return e.wrap
}

func (e *authError) LoggableAction() *actionlog.Action {
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

// AuthActionAuthenticate returns "system:auth.authenticate" error
//
// This function is auto-generated.
//
func AuthActionAuthenticate(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "authenticate",
		log:       "successfully authenticated with {credentials.kind}",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionIssueToken returns "system:auth.issueToken" error
//
// This function is auto-generated.
//
func AuthActionIssueToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "issueToken",
		log:       "token '{credentials.kind}' issued",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionValidateToken returns "system:auth.validateToken" error
//
// This function is auto-generated.
//
func AuthActionValidateToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "validateToken",
		log:       "token '{credentials.kind}' validated",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionChangePassword returns "system:auth.changePassword" error
//
// This function is auto-generated.
//
func AuthActionChangePassword(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "changePassword",
		log:       "password changed",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionInternalSignup returns "system:auth.internalSignup" error
//
// This function is auto-generated.
//
func AuthActionInternalSignup(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "internalSignup",
		log:       "{user.email} signed-up",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionConfirmEmail returns "system:auth.confirmEmail" error
//
// This function is auto-generated.
//
func AuthActionConfirmEmail(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "confirmEmail",
		log:       "email {user.email} confirmed",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionExternalSignup returns "system:auth.externalSignup" error
//
// This function is auto-generated.
//
func AuthActionExternalSignup(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "externalSignup",
		log:       "{user.email} signed-up after successful external authentication via {credentials.kind}",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionSendEmailConfirmationToken returns "system:auth.sendEmailConfirmationToken" error
//
// This function is auto-generated.
//
func AuthActionSendEmailConfirmationToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "sendEmailConfirmationToken",
		log:       "confirmation notification sent to {email}",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionSendPasswordResetToken returns "system:auth.sendPasswordResetToken" error
//
// This function is auto-generated.
//
func AuthActionSendPasswordResetToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "sendPasswordResetToken",
		log:       "password reset token sent to {email}",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionExchangePasswordResetToken returns "system:auth.exchangePasswordResetToken" error
//
// This function is auto-generated.
//
func AuthActionExchangePasswordResetToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "exchangePasswordResetToken",
		log:       "password reset token exchanged",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionAutoPromote returns "system:auth.autoPromote" error
//
// This function is auto-generated.
//
func AuthActionAutoPromote(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "autoPromote",
		log:       "auto-promoted to {role}",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionUpdateCredentials returns "system:auth.updateCredentials" error
//
// This function is auto-generated.
//
func AuthActionUpdateCredentials(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "updateCredentials",
		log:       "credentials {credentials.kind} updated",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionCreateCredentials returns "system:auth.createCredentials" error
//
// This function is auto-generated.
//
func AuthActionCreateCredentials(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "createCredentials",
		log:       "new credentials {credentials.kind} created",
		severity:  actionlog.Warning,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// AuthErrGeneric returns "system:auth.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func AuthErrGeneric(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *authActionProps {
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

// AuthErrSubscription returns "system:auth.subscription" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AuthErrSubscription(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "subscription",
		action:    "error",
		message:   "{err}",
		log:       "{err}",
		severity:  actionlog.Warning,
		props: func() *authActionProps {
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

// AuthErrInvalidCredentials returns "system:auth.invalidCredentials" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AuthErrInvalidCredentials(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "invalidCredentials",
		action:    "error",
		message:   "invalid username and password combination",
		log:       "{email} failed to authenticate with {credentials.kind}",
		severity:  actionlog.Warning,
		props: func() *authActionProps {
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

// AuthErrInvalidEmailFormat returns "system:auth.invalidEmailFormat" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AuthErrInvalidEmailFormat(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "invalidEmailFormat",
		action:    "error",
		message:   "invalid email",
		log:       "invalid email",
		severity:  actionlog.Alert,
		props: func() *authActionProps {
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

// AuthErrInvalidHandle returns "system:auth.invalidHandle" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AuthErrInvalidHandle(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "invalidHandle",
		action:    "error",
		message:   "invalid handle",
		log:       "invalid handle",
		severity:  actionlog.Alert,
		props: func() *authActionProps {
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

// AuthErrFailedForUnknownUser returns "system:auth.invalidCredentials" audit event as actionlog.Warning
//
// Note: This error will be wrapped with safe (invalidCredentials) error!
//
// This function is auto-generated.
//
func AuthErrFailedForUnknownUser(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "failedForUnknownUser",
		action:    "error",
		message:   "failedForUnknownUser",
		log:       "unknown user {email} tried to log-in with {credentials.kind}",
		severity:  actionlog.Warning,
		props: func() *authActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	// Wrap with safe error
	return AuthErrInvalidCredentials().Wrap(e)

}

// AuthErrFailedForDisabledUser returns "system:auth.invalidCredentials" audit event as actionlog.Warning
//
// Note: This error will be wrapped with safe (invalidCredentials) error!
//
// This function is auto-generated.
//
func AuthErrFailedForDisabledUser(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "failedForDisabledUser",
		action:    "error",
		message:   "failedForDisabledUser",
		log:       "disabled user {user} tried to log-in with {credentials.kind}",
		severity:  actionlog.Warning,
		props: func() *authActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	// Wrap with safe error
	return AuthErrInvalidCredentials().Wrap(e)

}

// AuthErrFailedUnconfirmedEmail returns "system:auth.failedUnconfirmedEmail" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AuthErrFailedUnconfirmedEmail(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "failedUnconfirmedEmail",
		action:    "error",
		message:   "system requires confirmed email before logging in",
		log:       "failed to log-in with with unconfirmed email",
		severity:  actionlog.Alert,
		props: func() *authActionProps {
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

// AuthErrInteralLoginDisabledByConfig returns "system:auth.interalLoginDisabledByConfig" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AuthErrInteralLoginDisabledByConfig(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "interalLoginDisabledByConfig",
		action:    "error",
		message:   "internal login (username/password) is disabled",
		log:       "internal login (username/password) is disabled",
		severity:  actionlog.Alert,
		props: func() *authActionProps {
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

// AuthErrInternalSignupDisabledByConfig returns "system:auth.internalSignupDisabledByConfig" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AuthErrInternalSignupDisabledByConfig(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "internalSignupDisabledByConfig",
		action:    "error",
		message:   "internal sign-up (username/password) is disabled",
		log:       "internal sign-up (username/password) is disabled",
		severity:  actionlog.Alert,
		props: func() *authActionProps {
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

// AuthErrPasswordChangeFailedForUnknownUser returns "system:auth.passwordChangeFailedForUnknownUser" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AuthErrPasswordChangeFailedForUnknownUser(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "passwordChangeFailedForUnknownUser",
		action:    "error",
		message:   "failed to change password for the unknown user",
		log:       "failed to change password for the unknown user",
		severity:  actionlog.Alert,
		props: func() *authActionProps {
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

// AuthErrPasswodResetFailedOldPasswordCheckFailed returns "system:auth.passwodResetFailedOldPasswordCheckFailed" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AuthErrPasswodResetFailedOldPasswordCheckFailed(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "passwodResetFailedOldPasswordCheckFailed",
		action:    "error",
		message:   "failed to change password, old password does not match",
		log:       "failed to change password, old password does not match",
		severity:  actionlog.Alert,
		props: func() *authActionProps {
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

// AuthErrPasswordResetDisabledByConfig returns "system:auth.passwordResetDisabledByConfig" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AuthErrPasswordResetDisabledByConfig(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "passwordResetDisabledByConfig",
		action:    "error",
		message:   "password reset is disabled",
		log:       "password reset is disabled",
		severity:  actionlog.Alert,
		props: func() *authActionProps {
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

// AuthErrPasswordNotSecure returns "system:auth.passwordNotSecure" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AuthErrPasswordNotSecure(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "passwordNotSecure",
		action:    "error",
		message:   "provided password is not secure; use longer password with more non-alphanumeric character",
		log:       "provided password is not secure; use longer password with more non-alphanumeric character",
		severity:  actionlog.Alert,
		props: func() *authActionProps {
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

// AuthErrExternalDisabledByConfig returns "system:auth.externalDisabledByConfig" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AuthErrExternalDisabledByConfig(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "externalDisabledByConfig",
		action:    "error",
		message:   "external authentication (using external authentication provider) is disabled",
		log:       "external authentication is disabled",
		severity:  actionlog.Warning,
		props: func() *authActionProps {
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

// AuthErrProfileWithoutValidEmail returns "system:auth.profileWithoutValidEmail" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AuthErrProfileWithoutValidEmail(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "profileWithoutValidEmail",
		action:    "error",
		message:   "external authentication provider returned profile without valid email",
		log:       "external authentication provider {credentials.kind} returned profile without valid email",
		severity:  actionlog.Warning,
		props: func() *authActionProps {
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

// AuthErrCredentialsLinkedToInvalidUser returns "system:auth.credentialsLinkedToInvalidUser" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AuthErrCredentialsLinkedToInvalidUser(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "credentialsLinkedToInvalidUser",
		action:    "error",
		message:   "credentials {credentials.kind} linked to disabled or deleted user {user}",
		log:       "credentials {credentials.kind} linked to disabled or deleted user {user}",
		severity:  actionlog.Warning,
		props: func() *authActionProps {
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

// AuthErrInvalidToken returns "system:auth.invalidToken" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AuthErrInvalidToken(props ...*authActionProps) *authError {
	var e = &authError{
		timestamp: time.Now(),
		resource:  "system:auth",
		error:     "invalidToken",
		action:    "error",
		message:   "invalid token",
		log:       "invalid token",
		severity:  actionlog.Warning,
		props: func() *authActionProps {
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
// action (optional) fn will be used to construct authAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc auth) recordAction(ctx context.Context, props *authActionProps, action func(...*authActionProps) *authAction, err error) error {
	var (
		ok bool

		// Return error
		retError *authError

		// Recorder error
		recError *authError
	)

	if err != nil {
		if retError, ok = err.(*authError); !ok {
			// got non-auth error, wrap it with AuthErrGeneric
			retError = AuthErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use AuthErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type authError
				if unwrappedSinkError, ok := unwrappedError.(*authError); ok {
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
