package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/setting_actions.yaml

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
	settingsActionProps struct {
		settings *types.SettingValue
	}

	settingsAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *settingsActionProps
	}

	settingsLogMetaKey   struct{}
	settingsPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setSettings updates settingsActionProps's settings
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *settingsActionProps) setSettings(settings *types.SettingValue) *settingsActionProps {
	p.settings = settings
	return p
}

// Serialize converts settingsActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p settingsActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.settings != nil {
		m.Set("settings.name", p.settings.Name, true)
		m.Set("settings.value", p.settings.Value, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p settingsActionProps) Format(in string, err error) string {
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

	if p.settings != nil {
		// replacement for "{{settings}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{settings}}",
			fns(
				p.settings.Name,
				p.settings.Value,
			),
		)
		pairs = append(pairs, "{{settings.name}}", fns(p.settings.Name))
		pairs = append(pairs, "{{settings.value}}", fns(p.settings.Value))
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
func (a *settingsAction) String() string {
	var props = &settingsActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *settingsAction) ToAction() *actionlog.Action {
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

// SettingsActionLookup returns "system:setting.lookup" action
//
// This function is auto-generated.
//
func SettingsActionLookup(props ...*settingsActionProps) *settingsAction {
	a := &settingsAction{
		timestamp: time.Now(),
		resource:  "system:setting",
		action:    "lookup",
		log:       "looked-up for a setting",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// SettingsErrGeneric returns "system:setting.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func SettingsErrGeneric(mm ...*settingsActionProps) *errors.Error {
	var p = &settingsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:setting"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(settingsLogMetaKey{}, "{err}"),
		errors.Meta(settingsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "settings.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SettingsErrNotAllowedToRead returns "system:setting.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func SettingsErrNotAllowedToRead(mm ...*settingsActionProps) *errors.Error {
	var p = &settingsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this setting", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:setting"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(settingsLogMetaKey{}, "failed to read {{settings.name}}; insufficient permissions"),
		errors.Meta(settingsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "settings.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SettingsErrNotAllowedToManage returns "system:setting.notAllowedToManage" as *errors.Error
//
//
// This function is auto-generated.
//
func SettingsErrNotAllowedToManage(mm ...*settingsActionProps) *errors.Error {
	var p = &settingsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to manage this setting", nil),

		errors.Meta("type", "notAllowedToManage"),
		errors.Meta("resource", "system:setting"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(settingsLogMetaKey{}, "failed to manage {{settings.name}}; insufficient permissions"),
		errors.Meta(settingsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "settings.errors.notAllowedToManage"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SettingsErrInvalidPasswordMinLength returns "system:setting.invalidPasswordMinLength" as *errors.Error
//
//
// This function is auto-generated.
//
func SettingsErrInvalidPasswordMinLength(mm ...*settingsActionProps) *errors.Error {
	var p = &settingsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("password constraint minimum length should be at least 8 characters", nil),

		errors.Meta("type", "invalidPasswordMinLength"),
		errors.Meta("resource", "system:setting"),

		errors.Meta(settingsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "settings.errors.invalidPasswordMinLength"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SettingsErrInvalidPasswordMinUpperCase returns "system:setting.invalidPasswordMinUpperCase" as *errors.Error
//
//
// This function is auto-generated.
//
func SettingsErrInvalidPasswordMinUpperCase(mm ...*settingsActionProps) *errors.Error {
	var p = &settingsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("password constraint minimum upper case count should not be a negative number", nil),

		errors.Meta("type", "invalidPasswordMinUpperCase"),
		errors.Meta("resource", "system:setting"),

		errors.Meta(settingsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "settings.errors.invalidPasswordMinUpperCase"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SettingsErrInvalidPasswordMinLowerCase returns "system:setting.invalidPasswordMinLowerCase" as *errors.Error
//
//
// This function is auto-generated.
//
func SettingsErrInvalidPasswordMinLowerCase(mm ...*settingsActionProps) *errors.Error {
	var p = &settingsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("password constraint minimum lower case count should not be a negative number", nil),

		errors.Meta("type", "invalidPasswordMinLowerCase"),
		errors.Meta("resource", "system:setting"),

		errors.Meta(settingsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "settings.errors.invalidPasswordMinLowerCase"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SettingsErrInvalidPasswordMinNumCount returns "system:setting.invalidPasswordMinNumCount" as *errors.Error
//
//
// This function is auto-generated.
//
func SettingsErrInvalidPasswordMinNumCount(mm ...*settingsActionProps) *errors.Error {
	var p = &settingsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("password constraint minimum number count should not be a negative number", nil),

		errors.Meta("type", "invalidPasswordMinNumCount"),
		errors.Meta("resource", "system:setting"),

		errors.Meta(settingsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "settings.errors.invalidPasswordMinNumCount"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SettingsErrInvalidPasswordMinSpecialCharCount returns "system:setting.invalidPasswordMinSpecialCharCount" as *errors.Error
//
//
// This function is auto-generated.
//
func SettingsErrInvalidPasswordMinSpecialCharCount(mm ...*settingsActionProps) *errors.Error {
	var p = &settingsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("password constraint minimum special character count should not be a negative number", nil),

		errors.Meta("type", "invalidPasswordMinSpecialCharCount"),
		errors.Meta("resource", "system:setting"),

		errors.Meta(settingsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "settings.errors.invalidPasswordMinSpecialCharCount"),

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
func (svc settings) recordAction(ctx context.Context, props *settingsActionProps, actionFn func(...*settingsActionProps) *settingsAction, err error) error {
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
		a.Description = props.Format(m.AsString(settingsLogMetaKey{}), err)

		if p, has := m[settingsPropsMetaKey{}]; has {
			a.Meta = p.(*settingsActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
