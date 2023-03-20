package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/statistics_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"strings"
	"time"
)

type (
	statisticsActionProps struct {
	}

	statisticsAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *statisticsActionProps
	}

	statisticsLogMetaKey   struct{}
	statisticsPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods

// Serialize converts statisticsActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p statisticsActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p statisticsActionProps) Format(in string, err error) string {
	var (
		pairs = []string{"{{err}}"}
	)

	if err != nil {
		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
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
func (a *statisticsAction) String() string {
	var props = &statisticsActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *statisticsAction) ToAction() *actionlog.Action {
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

// StatisticsActionServe returns "system:statistics.serve" action
//
// This function is auto-generated.
//
func StatisticsActionServe(props ...*statisticsActionProps) *statisticsAction {
	a := &statisticsAction{
		timestamp: time.Now(),
		resource:  "system:statistics",
		action:    "serve",
		log:       "metrics served",
		severity:  actionlog.Debug,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// StatisticsErrGeneric returns "system:statistics.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func StatisticsErrGeneric(mm ...*statisticsActionProps) *errors.Error {
	var p = &statisticsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:statistics"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(statisticsLogMetaKey{}, "{err}"),
		errors.Meta(statisticsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "statistics.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// StatisticsErrNotAllowedToReadStatistics returns "system:statistics.notAllowedToReadStatistics" as *errors.Error
//
//
// This function is auto-generated.
//
func StatisticsErrNotAllowedToReadStatistics(mm ...*statisticsActionProps) *errors.Error {
	var p = &statisticsActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read statistics", nil),

		errors.Meta("type", "notAllowedToReadStatistics"),
		errors.Meta("resource", "system:statistics"),

		errors.Meta(statisticsPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "statistics.errors.notAllowedToReadStatistics"),

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
func (svc statistics) recordAction(ctx context.Context, props *statisticsActionProps, actionFn func(...*statisticsActionProps) *statisticsAction, err error) error {
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
		a.Description = props.Format(m.AsString(statisticsLogMetaKey{}), err)

		if p, has := m[statisticsPropsMetaKey{}]; has {
			a.Meta = p.(*statisticsActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
