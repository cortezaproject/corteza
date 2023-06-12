package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/auth_actions.yaml

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
	authActionProps struct {
		email       string
		provider    string
		credentials *types.Credential
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

	authLogMetaKey   struct{}
	authPropsMetaKey struct{}
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
// This function is auto-generated.
//
func (p *authActionProps) setEmail(email string) *authActionProps {
	p.email = email
	return p
}

// setProvider updates authActionProps's provider
//
// This function is auto-generated.
//
func (p *authActionProps) setProvider(provider string) *authActionProps {
	p.provider = provider
	return p
}

// setCredentials updates authActionProps's credentials
//
// This function is auto-generated.
//
func (p *authActionProps) setCredentials(credentials *types.Credential) *authActionProps {
	p.credentials = credentials
	return p
}

// setRole updates authActionProps's role
//
// This function is auto-generated.
//
func (p *authActionProps) setRole(role *types.Role) *authActionProps {
	p.role = role
	return p
}

// setUser updates authActionProps's user
//
// This function is auto-generated.
//
func (p *authActionProps) setUser(user *types.User) *authActionProps {
	p.user = user
	return p
}

// Serialize converts authActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p authActionProps) Serialize() actionlog.Meta {
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
func (p authActionProps) Format(in string, err error) string {
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
	pairs = append(pairs, "{{email}}", fns(p.email))
	pairs = append(pairs, "{{provider}}", fns(p.provider))

	if p.credentials != nil {
		// replacement for "{{credentials}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{credentials}}",
			fns(
				p.credentials.Kind,
				p.credentials.Label,
				p.credentials.ID,
			),
		)
		pairs = append(pairs, "{{credentials.kind}}", fns(p.credentials.Kind))
		pairs = append(pairs, "{{credentials.label}}", fns(p.credentials.Label))
		pairs = append(pairs, "{{credentials.ID}}", fns(p.credentials.ID))
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

	if p.user != nil {
		// replacement for "{{user}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{user}}",
			fns(
				p.user.Handle,
				p.user.Name,
				p.user.ID,
				p.user.Email,
				p.user.SuspendedAt,
				p.user.DeletedAt,
			),
		)
		pairs = append(pairs, "{{user.handle}}", fns(p.user.Handle))
		pairs = append(pairs, "{{user.name}}", fns(p.user.Name))
		pairs = append(pairs, "{{user.ID}}", fns(p.user.ID))
		pairs = append(pairs, "{{user.email}}", fns(p.user.Email))
		pairs = append(pairs, "{{user.suspendedAt}}", fns(p.user.SuspendedAt))
		pairs = append(pairs, "{{user.deletedAt}}", fns(p.user.DeletedAt))
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

	return props.Format(a.log, nil)
}

func (e *authAction) ToAction() *actionlog.Action {
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

// AuthActionAuthenticate returns "system:auth.authenticate" action
//
// This function is auto-generated.
//
func AuthActionAuthenticate(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "authenticate",
		log:       "successfully authenticated with {{credentials.kind}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionIssueToken returns "system:auth.issueToken" action
//
// This function is auto-generated.
//
func AuthActionIssueToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "issueToken",
		log:       "token '{{credentials.kind}}' issued",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionValidateToken returns "system:auth.validateToken" action
//
// This function is auto-generated.
//
func AuthActionValidateToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "validateToken",
		log:       "token '{{credentials.kind}}' validated",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionChangePassword returns "system:auth.changePassword" action
//
// This function is auto-generated.
//
func AuthActionChangePassword(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "changePassword",
		log:       "password changed",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionInternalSignup returns "system:auth.internalSignup" action
//
// This function is auto-generated.
//
func AuthActionInternalSignup(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "internalSignup",
		log:       "{{user.email}} signed-up",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionConfirmEmail returns "system:auth.confirmEmail" action
//
// This function is auto-generated.
//
func AuthActionConfirmEmail(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "confirmEmail",
		log:       "email {{user.email}} confirmed",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionExternalSignup returns "system:auth.externalSignup" action
//
// This function is auto-generated.
//
func AuthActionExternalSignup(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "externalSignup",
		log:       "{{user.email}} signed-up after successful external authentication via {{credentials.kind}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionSendEmailConfirmationToken returns "system:auth.sendEmailConfirmationToken" action
//
// This function is auto-generated.
//
func AuthActionSendEmailConfirmationToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "sendEmailConfirmationToken",
		log:       "confirmation notification sent to {{email}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionSendPasswordResetToken returns "system:auth.sendPasswordResetToken" action
//
// This function is auto-generated.
//
func AuthActionSendPasswordResetToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "sendPasswordResetToken",
		log:       "password reset token sent to {{email}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionExchangePasswordResetToken returns "system:auth.exchangePasswordResetToken" action
//
// This function is auto-generated.
//
func AuthActionExchangePasswordResetToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "exchangePasswordResetToken",
		log:       "password reset token exchanged",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionGeneratePasswordCreateToken returns "system:auth.generatePasswordCreateToken" action
//
// This function is auto-generated.
//
func AuthActionGeneratePasswordCreateToken(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "generatePasswordCreateToken",
		log:       "password create token generated for {{email}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionAutoPromote returns "system:auth.autoPromote" action
//
// This function is auto-generated.
//
func AuthActionAutoPromote(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "autoPromote",
		log:       "auto-promoted to {{role}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionUpdateCredentials returns "system:auth.updateCredentials" action
//
// This function is auto-generated.
//
func AuthActionUpdateCredentials(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "updateCredentials",
		log:       "credentials {{credentials.kind}} updated",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionCreateCredentials returns "system:auth.createCredentials" action
//
// This function is auto-generated.
//
func AuthActionCreateCredentials(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "createCredentials",
		log:       "new credentials {{credentials.kind}} created",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionImpersonate returns "system:auth.impersonate" action
//
// This function is auto-generated.
//
func AuthActionImpersonate(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "impersonate",
		log:       "impersonating {{user}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionTotpConfigure returns "system:auth.totpConfigure" action
//
// This function is auto-generated.
//
func AuthActionTotpConfigure(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "totpConfigure",
		log:       "time-based one-time-password for {{user}} configured",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionTotpRemove returns "system:auth.totpRemove" action
//
// This function is auto-generated.
//
func AuthActionTotpRemove(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "totpRemove",
		log:       "time-based one-time-password for {{user}} removed",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionTotpValidate returns "system:auth.totpValidate" action
//
// This function is auto-generated.
//
func AuthActionTotpValidate(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "totpValidate",
		log:       "time-based one-time-password for {{user}} validated",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionEmailOtpVerify returns "system:auth.emailOtpVerify" action
//
// This function is auto-generated.
//
func AuthActionEmailOtpVerify(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "emailOtpVerify",
		log:       "email one-time-password for {{user}} verified",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionAccessTokensRemoved returns "system:auth.accessTokensRemoved" action
//
// This function is auto-generated.
//
func AuthActionAccessTokensRemoved(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "accessTokensRemoved",
		log:       "access tokens for {{user}} removed",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AuthActionSendInviteEMail returns "system:auth.sendInviteEMail" action
//
// This function is auto-generated.
//
func AuthActionSendInviteEMail(props ...*authActionProps) *authAction {
	a := &authAction{
		timestamp: time.Now(),
		resource:  "system:auth",
		action:    "sendInviteEMail",
		log:       "invite email sent to {{email}}",
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

// AuthErrGeneric returns "system:auth.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrGeneric(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "{err}"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrInvalidCredentials returns "system:auth.invalidCredentials" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrInvalidCredentials(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid username and password combination", nil),

		errors.Meta("type", "invalidCredentials"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "{{email}} failed to authenticate with {{credentials.kind}}"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.invalidCredentials"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrInvalidEmailFormat returns "system:auth.invalidEmailFormat" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrInvalidEmailFormat(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid email", nil),

		errors.Meta("type", "invalidEmailFormat"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.invalidEmailFormat"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrInvalidHandle returns "system:auth.invalidHandle" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrInvalidHandle(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.invalidHandle"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrFailedForUnknownUser returns "system:auth.failedForUnknownUser" as *errors.Error
//
// Note: This error will be wrapped with safe (system:auth.invalidCredentials) error!
//
// This function is auto-generated.
//
func AuthErrFailedForUnknownUser(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		"failedForUnknownUser",

		errors.Meta("type", "failedForUnknownUser"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "unknown user {{email}} tried to log-in with {{credentials.kind}}"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.failedForUnknownUser"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	// Wrap with safe error
	e = AuthErrInvalidCredentials().Wrap(e)

	return e
}

// AuthErrFailedForDeletedUser returns "system:auth.failedForDeletedUser" as *errors.Error
//
// Note: This error will be wrapped with safe (system:auth.invalidCredentials) error!
//
// This function is auto-generated.
//
func AuthErrFailedForDeletedUser(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		"failedForDeletedUser",

		errors.Meta("type", "failedForDeletedUser"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "deleted user {{user}} tried to log-in with {{credentials.kind}}"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.failedForDeletedUser"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	// Wrap with safe error
	e = AuthErrInvalidCredentials().Wrap(e)

	return e
}

// AuthErrFailedForSuspendedUser returns "system:auth.failedForSuspendedUser" as *errors.Error
//
// Note: This error will be wrapped with safe (system:auth.invalidCredentials) error!
//
// This function is auto-generated.
//
func AuthErrFailedForSuspendedUser(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		"failedForSuspendedUser",

		errors.Meta("type", "failedForSuspendedUser"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "suspended user {{user}} tried to log-in with {{credentials.kind}}"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.failedForSuspendedUser"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	// Wrap with safe error
	e = AuthErrInvalidCredentials().Wrap(e)

	return e
}

// AuthErrFailedForSystemUser returns "system:auth.failedForSystemUser" as *errors.Error
//
// Note: This error will be wrapped with safe (system:auth.invalidCredentials) error!
//
// This function is auto-generated.
//
func AuthErrFailedForSystemUser(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		"failedForSystemUser",

		errors.Meta("type", "failedForSystemUser"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "system user {{user}} tried to log-in with {{credentials.kind}}"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.failedForSystemUser"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	// Wrap with safe error
	e = AuthErrInvalidCredentials().Wrap(e)

	return e
}

// AuthErrFailedUnconfirmedEmail returns "system:auth.failedUnconfirmedEmail" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrFailedUnconfirmedEmail(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("system requires confirmed email before logging in", nil),

		errors.Meta("type", "failedUnconfirmedEmail"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "failed to log-in with with unconfirmed email"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.failedUnconfirmedEmail"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrInternalLoginDisabledByConfig returns "system:auth.internalLoginDisabledByConfig" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrInternalLoginDisabledByConfig(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("internal login (username/password) is disabled", nil),

		errors.Meta("type", "internalLoginDisabledByConfig"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.internalLoginDisabledByConfig"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrInternalSignupDisabledByConfig returns "system:auth.internalSignupDisabledByConfig" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrInternalSignupDisabledByConfig(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("internal sign-up (username/password) is disabled", nil),

		errors.Meta("type", "internalSignupDisabledByConfig"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.internalSignupDisabledByConfig"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrPasswordChangeFailedForUnknownUser returns "system:auth.passwordChangeFailedForUnknownUser" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrPasswordChangeFailedForUnknownUser(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to change password for the unknown user", nil),

		errors.Meta("type", "passwordChangeFailedForUnknownUser"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.passwordChangeFailedForUnknownUser"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrPasswordResetFailedOldPasswordCheckFailed returns "system:auth.passwordResetFailedOldPasswordCheckFailed" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrPasswordResetFailedOldPasswordCheckFailed(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to change password, old password does not match", nil),

		errors.Meta("type", "passwordResetFailedOldPasswordCheckFailed"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.passwordResetFailedOldPasswordCheckFailed"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrPasswordSetFailedReusedPasswordCheckFailed returns "system:auth.passwordSetFailedReusedPasswordCheckFailed" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrPasswordSetFailedReusedPasswordCheckFailed(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to set password, already used", nil),

		errors.Meta("type", "passwordSetFailedReusedPasswordCheckFailed"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.passwordSetFailedReusedPasswordCheckFailed"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrPasswordCreateFailedForUnknownUser returns "system:auth.passwordCreateFailedForUnknownUser" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrPasswordCreateFailedForUnknownUser(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to create password for the unknown user", nil),

		errors.Meta("type", "passwordCreateFailedForUnknownUser"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.passwordCreateFailedForUnknownUser"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrPasswordResetDisabledByConfig returns "system:auth.passwordResetDisabledByConfig" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrPasswordResetDisabledByConfig(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("password reset is disabled", nil),

		errors.Meta("type", "passwordResetDisabledByConfig"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.passwordResetDisabledByConfig"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrPasswordCreateDisabledByConfig returns "system:auth.passwordCreateDisabledByConfig" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrPasswordCreateDisabledByConfig(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("password create is disabled", nil),

		errors.Meta("type", "passwordCreateDisabledByConfig"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.passwordCreateDisabledByConfig"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrPasswordNotSecure returns "system:auth.passwordNotSecure" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrPasswordNotSecure(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("provided password is not secure; use longer password to follow the password policy", nil),

		errors.Meta("type", "passwordNotSecure"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.passwordNotSecure"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrExternalDisabledByConfig returns "system:auth.externalDisabledByConfig" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrExternalDisabledByConfig(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("external authentication (using external authentication provider) is disabled", nil),

		errors.Meta("type", "externalDisabledByConfig"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "external authentication is disabled"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.externalDisabledByConfig"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrProfileWithoutValidEmail returns "system:auth.profileWithoutValidEmail" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrProfileWithoutValidEmail(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("external authentication provider returned profile without valid email", nil),

		errors.Meta("type", "profileWithoutValidEmail"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "external authentication provider {{credentials.kind}} returned profile without valid email"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.profileWithoutValidEmail"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrCredentialsLinkedToInvalidUser returns "system:auth.credentialsLinkedToInvalidUser" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrCredentialsLinkedToInvalidUser(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("credentials {{credentials.kind}} linked to disabled or deleted user {{user}}", nil),

		errors.Meta("type", "credentialsLinkedToInvalidUser"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.credentialsLinkedToInvalidUser"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrInvalidToken returns "system:auth.invalidToken" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrInvalidToken(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid token", nil),

		errors.Meta("type", "invalidToken"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.invalidToken"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrNotAllowedToImpersonate returns "system:auth.notAllowedToImpersonate" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrNotAllowedToImpersonate(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to impersonate this user", nil),

		errors.Meta("type", "notAllowedToImpersonate"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.notAllowedToImpersonate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrNotAllowedToRemoveTOTP returns "system:auth.notAllowedToRemoveTOTP" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrNotAllowedToRemoveTOTP(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to remove TOTP", nil),

		errors.Meta("type", "notAllowedToRemoveTOTP"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.notAllowedToRemoveTOTP"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrUnconfiguredTOTP returns "system:auth.unconfiguredTOTP" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrUnconfiguredTOTP(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("TOTP not configured", nil),

		errors.Meta("type", "unconfiguredTOTP"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.unconfiguredTOTP"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrNotAllowedToConfigureTOTP returns "system:auth.notAllowedToConfigureTOTP" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrNotAllowedToConfigureTOTP(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to configure TOTP", nil),

		errors.Meta("type", "notAllowedToConfigureTOTP"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.notAllowedToConfigureTOTP"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrEnforcedMFAWithTOTP returns "system:auth.enforcedMFAWithTOTP" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrEnforcedMFAWithTOTP(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("TOTP is enforced and cannot be disabled", nil),

		errors.Meta("type", "enforcedMFAWithTOTP"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.enforcedMFAWithTOTP"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrInvalidTOTP returns "system:auth.invalidTOTP" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrInvalidTOTP(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid code", nil),

		errors.Meta("type", "invalidTOTP"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.invalidTOTP"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrDisabledMFAWithTOTP returns "system:auth.disabledMFAWithTOTP" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrDisabledMFAWithTOTP(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("multi factor authentication with TOTP is disabled", nil),

		errors.Meta("type", "disabledMFAWithTOTP"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.disabledMFAWithTOTP"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrDisabledMFAWithEmailOTP returns "system:auth.disabledMFAWithEmailOTP" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrDisabledMFAWithEmailOTP(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("multi factor authentication with email OTP is disabled", nil),

		errors.Meta("type", "disabledMFAWithEmailOTP"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.disabledMFAWithEmailOTP"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrEnforcedMFAWithEmailOTP returns "system:auth.enforcedMFAWithEmailOTP" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrEnforcedMFAWithEmailOTP(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("OTP over email is enforced and cannot be disabled", nil),

		errors.Meta("type", "enforcedMFAWithEmailOTP"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.enforcedMFAWithEmailOTP"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrInvalidEmailOTP returns "system:auth.invalidEmailOTP" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrInvalidEmailOTP(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid code", nil),

		errors.Meta("type", "invalidEmailOTP"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.invalidEmailOTP"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrRateLimitExceeded returns "system:auth.rateLimitExceeded" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrRateLimitExceeded(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("rate limit exceeded", nil),

		errors.Meta("type", "rateLimitExceeded"),
		errors.Meta("resource", "system:auth"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(authLogMetaKey{}, "rate limit exceeded for {{user}}"),
		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.rateLimitExceeded"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrMaxUserLimitReached returns "system:auth.maxUserLimitReached" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrMaxUserLimitReached(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("you have reached your user limit, contact your Corteza administrator", nil),

		errors.Meta("type", "maxUserLimitReached"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.maxUserLimitReached"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AuthErrDisabledSendUserInviteEmail returns "system:auth.disabledSendUserInviteEmail" as *errors.Error
//
//
// This function is auto-generated.
//
func AuthErrDisabledSendUserInviteEmail(mm ...*authActionProps) *errors.Error {
	var p = &authActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("sending user invite email is disabled", nil),

		errors.Meta("type", "disabledSendUserInviteEmail"),
		errors.Meta("resource", "system:auth"),

		errors.Meta(authPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.disabledSendUserInviteEmail"),

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
func (svc auth) recordAction(ctx context.Context, props *authActionProps, actionFn func(...*authActionProps) *authAction, err error) error {
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
		a.Description = props.Format(m.AsString(authLogMetaKey{}), err)

		if p, has := m[authPropsMetaKey{}]; has {
			a.Meta = p.(*authActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
