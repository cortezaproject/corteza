package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/credentials_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
	"time"
)

type (
	credentialsActionProps struct {
		user        *types.User
		credentials *types.Credential
	}

	credentialsAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *credentialsActionProps
	}

	credentialsLogMetaKey   struct{}
	credentialsPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setUser updates credentialsActionProps's user
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *credentialsActionProps) setUser(user *types.User) *credentialsActionProps {
	p.user = user
	return p
}

// setCredentials updates credentialsActionProps's credentials
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *credentialsActionProps) setCredentials(credentials *types.Credential) *credentialsActionProps {
	p.credentials = credentials
	return p
}

// Serialize converts credentialsActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p credentialsActionProps) Serialize() actionlog.Meta {
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
	if p.credentials != nil {
		m.Set("credentials.kind", p.credentials.Kind, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p credentialsActionProps) Format(in string, err error) string {
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

	if p.user != nil {
		// replacement for "{{user}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{user}}",
			fns(
				p.user.Handle,
				p.user.Email,
				p.user.Name,
				p.user.Username,
				p.user.ID,
			),
		)
		pairs = append(pairs, "{{user.handle}}", fns(p.user.Handle))
		pairs = append(pairs, "{{user.email}}", fns(p.user.Email))
		pairs = append(pairs, "{{user.name}}", fns(p.user.Name))
		pairs = append(pairs, "{{user.username}}", fns(p.user.Username))
		pairs = append(pairs, "{{user.ID}}", fns(p.user.ID))
	}

	if p.credentials != nil {
		// replacement for "{{credentials}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{credentials}}",
			fns(
				p.credentials.Kind,
			),
		)
		pairs = append(pairs, "{{credentials.kind}}", fns(p.credentials.Kind))
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
func (a *credentialsAction) String() string {
	var props = &credentialsActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *credentialsAction) ToAction() *actionlog.Action {
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

// CredentialsActionSearch returns "system:credentials.search" action
//
// This function is auto-generated.
//
func CredentialsActionSearch(props ...*credentialsActionProps) *credentialsAction {
	a := &credentialsAction{
		timestamp: time.Now(),
		resource:  "system:credentials",
		action:    "search",
		log:       "searched for matching users",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// CredentialsActionDelete returns "system:credentials.delete" action
//
// This function is auto-generated.
//
func CredentialsActionDelete(props ...*credentialsActionProps) *credentialsAction {
	a := &credentialsAction{
		timestamp: time.Now(),
		resource:  "system:credentials",
		action:    "delete",
		log:       "deleted {{user}}",
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

// CredentialsErrGeneric returns "system:credentials.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func CredentialsErrGeneric(mm ...*credentialsActionProps) *errors.Error {
	var p = &credentialsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:credentials"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(credentialsLogMetaKey{}, "{err}"),
		errors.Meta(credentialsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "credentials.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// CredentialsErrNotFound returns "system:credentials.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func CredentialsErrNotFound(mm ...*credentialsActionProps) *errors.Error {
	var p = &credentialsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("credentials not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:credentials"),

		errors.Meta(credentialsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "credentials.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// CredentialsErrInvalidID returns "system:credentials.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func CredentialsErrInvalidID(mm ...*credentialsActionProps) *errors.Error {
	var p = &credentialsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:credentials"),

		errors.Meta(credentialsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "credentials.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// CredentialsErrNotAllowedToManage returns "system:credentials.notAllowedToManage" as *errors.Error
//
//
// This function is auto-generated.
//
func CredentialsErrNotAllowedToManage(mm ...*credentialsActionProps) *errors.Error {
	var p = &credentialsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage credentials for this user", nil),

		errors.Meta("type", "notAllowedToManage"),
		errors.Meta("resource", "system:credentials"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(credentialsLogMetaKey{}, "failed to manage credentials for {{user.handle}}; insufficient permissions"),
		errors.Meta(credentialsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "credentials.errors.notAllowedToManage"),

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
func (svc credentials) recordAction(ctx context.Context, props *credentialsActionProps, actionFn func(...*credentialsActionProps) *credentialsAction, err error) error {
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
		a.Description = props.Format(m.AsString(credentialsLogMetaKey{}), err)

		if p, has := m[credentialsPropsMetaKey{}]; has {
			a.Meta = p.(*credentialsActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
