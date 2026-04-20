package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/knowledge_base_actions.yaml

import (
	"context"
	"fmt"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/system/types"
	"strings"
	"time"
)

type (
	knowledgeBaseActionProps struct {
		knowledgeBase *types.KnowledgeBase
		new           *types.KnowledgeBase
		update        *types.KnowledgeBase
		search        *types.KnowledgeBaseFilter
	}

	knowledgeBaseAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *knowledgeBaseActionProps
	}

	knowledgeBaseLogMetaKey   struct{}
	knowledgeBasePropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setKnowledgeBase updates knowledgeBaseActionProps's knowledgeBase
//
// This function is auto-generated.
func (p *knowledgeBaseActionProps) setKnowledgeBase(knowledgeBase *types.KnowledgeBase) *knowledgeBaseActionProps {
	p.knowledgeBase = knowledgeBase
	return p
}

// setNew updates knowledgeBaseActionProps's new
//
// This function is auto-generated.
func (p *knowledgeBaseActionProps) setNew(new *types.KnowledgeBase) *knowledgeBaseActionProps {
	p.new = new
	return p
}

// setUpdate updates knowledgeBaseActionProps's update
//
// This function is auto-generated.
func (p *knowledgeBaseActionProps) setUpdate(update *types.KnowledgeBase) *knowledgeBaseActionProps {
	p.update = update
	return p
}

// setSearch updates knowledgeBaseActionProps's search
//
// This function is auto-generated.
func (p *knowledgeBaseActionProps) setSearch(search *types.KnowledgeBaseFilter) *knowledgeBaseActionProps {
	p.search = search
	return p
}

// Serialize converts knowledgeBaseActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p knowledgeBaseActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.knowledgeBase != nil {
		m.Set("knowledgeBase.handle", p.knowledgeBase.Handle, true)
		m.Set("knowledgeBase.ID", p.knowledgeBase.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.search != nil {
		m.Set("search.knowledgeBaseID", p.search.KnowledgeBaseID, true)
		m.Set("search.handle", p.search.Handle, true)
		m.Set("search.deleted", p.search.Deleted, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p knowledgeBaseActionProps) Format(in string, err error) string {
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

	if p.knowledgeBase != nil {
		// replacement for "{{knowledgeBase}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{knowledgeBase}}",
			fns(
				p.knowledgeBase.Handle,
				p.knowledgeBase.ID,
			),
		)
		pairs = append(pairs, "{{knowledgeBase.handle}}", fns(p.knowledgeBase.Handle))
		pairs = append(pairs, "{{knowledgeBase.ID}}", fns(p.knowledgeBase.ID))
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

	if p.search != nil {
		// replacement for "{{search}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{search}}",
			fns(
				p.search.KnowledgeBaseID,
				p.search.Handle,
				p.search.Deleted,
			),
		)
		pairs = append(pairs, "{{search.knowledgeBaseID}}", fns(p.search.KnowledgeBaseID))
		pairs = append(pairs, "{{search.handle}}", fns(p.search.Handle))
		pairs = append(pairs, "{{search.deleted}}", fns(p.search.Deleted))
	}
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
func (a *knowledgeBaseAction) String() string {
	var props = &knowledgeBaseActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *knowledgeBaseAction) ToAction() *actionlog.Action {
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

// KnowledgeBaseActionSearch returns "system:knowledge-base.search" action
//
// This function is auto-generated.
func KnowledgeBaseActionSearch(props ...*knowledgeBaseActionProps) *knowledgeBaseAction {
	a := &knowledgeBaseAction{
		timestamp: time.Now(),
		resource:  "system:knowledge-base",
		action:    "search",
		log:       "searched for knowledge bases",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// KnowledgeBaseActionLookup returns "system:knowledge-base.lookup" action
//
// This function is auto-generated.
func KnowledgeBaseActionLookup(props ...*knowledgeBaseActionProps) *knowledgeBaseAction {
	a := &knowledgeBaseAction{
		timestamp: time.Now(),
		resource:  "system:knowledge-base",
		action:    "lookup",
		log:       "looked-up for a {{knowledgeBase}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// KnowledgeBaseActionCreate returns "system:knowledge-base.create" action
//
// This function is auto-generated.
func KnowledgeBaseActionCreate(props ...*knowledgeBaseActionProps) *knowledgeBaseAction {
	a := &knowledgeBaseAction{
		timestamp: time.Now(),
		resource:  "system:knowledge-base",
		action:    "create",
		log:       "created {{knowledgeBase}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// KnowledgeBaseActionUpdate returns "system:knowledge-base.update" action
//
// This function is auto-generated.
func KnowledgeBaseActionUpdate(props ...*knowledgeBaseActionProps) *knowledgeBaseAction {
	a := &knowledgeBaseAction{
		timestamp: time.Now(),
		resource:  "system:knowledge-base",
		action:    "update",
		log:       "updated {{knowledgeBase}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// KnowledgeBaseActionDelete returns "system:knowledge-base.delete" action
//
// This function is auto-generated.
func KnowledgeBaseActionDelete(props ...*knowledgeBaseActionProps) *knowledgeBaseAction {
	a := &knowledgeBaseAction{
		timestamp: time.Now(),
		resource:  "system:knowledge-base",
		action:    "delete",
		log:       "deleted {{knowledgeBase}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// KnowledgeBaseActionUndelete returns "system:knowledge-base.undelete" action
//
// This function is auto-generated.
func KnowledgeBaseActionUndelete(props ...*knowledgeBaseActionProps) *knowledgeBaseAction {
	a := &knowledgeBaseAction{
		timestamp: time.Now(),
		resource:  "system:knowledge-base",
		action:    "undelete",
		log:       "undeleted {{knowledgeBase}}",
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

// KnowledgeBaseErrGeneric returns "system:knowledge-base.generic" as *errors.Error
//
// This function is auto-generated.
func KnowledgeBaseErrGeneric(mm ...*knowledgeBaseActionProps) *errors.Error {
	var p = &knowledgeBaseActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:knowledge-base"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(knowledgeBaseLogMetaKey{}, "{err}"),
		errors.Meta(knowledgeBasePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "knowledge-base.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// KnowledgeBaseErrNotFound returns "system:knowledge-base.notFound" as *errors.Error
//
// This function is auto-generated.
func KnowledgeBaseErrNotFound(mm ...*knowledgeBaseActionProps) *errors.Error {
	var p = &knowledgeBaseActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("knowledge base not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:knowledge-base"),

		errors.Meta(knowledgeBasePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "knowledge-base.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// KnowledgeBaseErrInvalidID returns "system:knowledge-base.invalidID" as *errors.Error
//
// This function is auto-generated.
func KnowledgeBaseErrInvalidID(mm ...*knowledgeBaseActionProps) *errors.Error {
	var p = &knowledgeBaseActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:knowledge-base"),

		errors.Meta(knowledgeBasePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "knowledge-base.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// KnowledgeBaseErrStaleData returns "system:knowledge-base.staleData" as *errors.Error
//
// This function is auto-generated.
func KnowledgeBaseErrStaleData(mm ...*knowledgeBaseActionProps) *errors.Error {
	var p = &knowledgeBaseActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("knowledge base was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:knowledge-base"),

		errors.Meta(knowledgeBasePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "knowledge-base.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// KnowledgeBaseErrNotAllowedToCreate returns "system:knowledge-base.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func KnowledgeBaseErrNotAllowedToCreate(mm ...*knowledgeBaseActionProps) *errors.Error {
	var p = &knowledgeBaseActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a knowledge base", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:knowledge-base"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(knowledgeBaseLogMetaKey{}, "failed to create a knowledge base; insufficient permissions"),
		errors.Meta(knowledgeBasePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "knowledge-base.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// KnowledgeBaseErrNotAllowedToSearch returns "system:knowledge-base.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func KnowledgeBaseErrNotAllowedToSearch(mm ...*knowledgeBaseActionProps) *errors.Error {
	var p = &knowledgeBaseActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list knowledge bases", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:knowledge-base"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(knowledgeBaseLogMetaKey{}, "failed to search or list knowledge bases; insufficient permissions"),
		errors.Meta(knowledgeBasePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "knowledge-base.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// KnowledgeBaseErrNotAllowedToRead returns "system:knowledge-base.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func KnowledgeBaseErrNotAllowedToRead(mm ...*knowledgeBaseActionProps) *errors.Error {
	var p = &knowledgeBaseActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this knowledge base", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:knowledge-base"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(knowledgeBaseLogMetaKey{}, "failed to read {{knowledgeBase.handle}}; insufficient permissions"),
		errors.Meta(knowledgeBasePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "knowledge-base.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// KnowledgeBaseErrNotAllowedToUpdate returns "system:knowledge-base.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func KnowledgeBaseErrNotAllowedToUpdate(mm ...*knowledgeBaseActionProps) *errors.Error {
	var p = &knowledgeBaseActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this knowledge base", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:knowledge-base"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(knowledgeBaseLogMetaKey{}, "failed to update {{knowledgeBase.handle}}; insufficient permissions"),
		errors.Meta(knowledgeBasePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "knowledge-base.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// KnowledgeBaseErrNotAllowedToDelete returns "system:knowledge-base.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func KnowledgeBaseErrNotAllowedToDelete(mm ...*knowledgeBaseActionProps) *errors.Error {
	var p = &knowledgeBaseActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this knowledge base", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:knowledge-base"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(knowledgeBaseLogMetaKey{}, "failed to delete {{knowledgeBase.handle}}; insufficient permissions"),
		errors.Meta(knowledgeBasePropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "knowledge-base.errors.notAllowedToDelete"),

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
func (svc knowledgeBase) recordAction(ctx context.Context, props *knowledgeBaseActionProps, actionFn func(...*knowledgeBaseActionProps) *knowledgeBaseAction, err error) error {
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
		a.Description = props.Format(m.AsString(knowledgeBaseLogMetaKey{}), err)

		if p, has := m[knowledgeBasePropsMetaKey{}]; has {
			a.Meta = p.(*knowledgeBaseActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
