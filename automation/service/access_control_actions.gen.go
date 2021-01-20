package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/service/access_control_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"strings"
	"time"
)

type (
	accessControlActionProps struct {
		rule *rbac.Rule
	}

	accessControlAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *accessControlActionProps
	}

	accessControlLogMetaKey   struct{}
	accessControlPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setRule updates accessControlActionProps's rule
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *accessControlActionProps) setRule(rule *rbac.Rule) *accessControlActionProps {
	p.rule = rule
	return p
}

// Serialize converts accessControlActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p accessControlActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.rule != nil {
		m.Set("rule.operation", p.rule.Operation, true)
		m.Set("rule.roleID", p.rule.RoleID, true)
		m.Set("rule.access", p.rule.Access, true)
		m.Set("rule.resource", p.rule.Resource, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p accessControlActionProps) Format(in string, err error) string {
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
		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}

	if p.rule != nil {
		// replacement for "{rule}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{rule}",
			fns(
				p.rule.Operation,
				p.rule.RoleID,
				p.rule.Access,
				p.rule.Resource,
			),
		)
		pairs = append(pairs, "{rule.operation}", fns(p.rule.Operation))
		pairs = append(pairs, "{rule.roleID}", fns(p.rule.RoleID))
		pairs = append(pairs, "{rule.access}", fns(p.rule.Access))
		pairs = append(pairs, "{rule.resource}", fns(p.rule.Resource))
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
func (a *accessControlAction) String() string {
	var props = &accessControlActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *accessControlAction) ToAction() *actionlog.Action {
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

// AccessControlActionGrant returns "automation:access_control.grant" action
//
// This function is auto-generated.
//
func AccessControlActionGrant(props ...*accessControlActionProps) *accessControlAction {
	a := &accessControlAction{
		timestamp: time.Now(),
		resource:  "automation:access_control",
		action:    "grant",
		log:       "grant",
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

// AccessControlErrGeneric returns "automation:access_control.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func AccessControlErrGeneric(mm ...*accessControlActionProps) *errors.Error {
	var p = &accessControlActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "automation:access_control"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(accessControlLogMetaKey{}, "{err}"),
		errors.Meta(accessControlPropsMetaKey{}, p),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AccessControlErrNotAllowedToSetPermissions returns "automation:access_control.notAllowedToSetPermissions" as *errors.Error
//
//
// This function is auto-generated.
//
func AccessControlErrNotAllowedToSetPermissions(mm ...*accessControlActionProps) *errors.Error {
	var p = &accessControlActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to set permissions", nil),

		errors.Meta("type", "notAllowedToSetPermissions"),
		errors.Meta("resource", "automation:access_control"),

		errors.Meta(accessControlPropsMetaKey{}, p),

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
func (svc accessControl) recordAction(ctx context.Context, props *accessControlActionProps, actionFn func(...*accessControlActionProps) *accessControlAction, err error) error {
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
		a.Description = props.Format(m.AsString(accessControlLogMetaKey{}), err)

		if p, has := m[accessControlPropsMetaKey{}]; has {
			a.Meta = p.(*accessControlActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
