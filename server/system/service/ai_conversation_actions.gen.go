package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/ai_conversation_actions.yaml

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
	aiConversationActionProps struct {
		aiConversation *types.AiConversation
		new            *types.AiConversation
		update         *types.AiConversation
		filter         *types.AiConversationFilter
	}

	aiConversationAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *aiConversationActionProps
	}

	aiConversationLogMetaKey   struct{}
	aiConversationPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setAiConversation updates aiConversationActionProps's aiConversation
//
// This function is auto-generated.
func (p *aiConversationActionProps) setAiConversation(aiConversation *types.AiConversation) *aiConversationActionProps {
	p.aiConversation = aiConversation
	return p
}

// setNew updates aiConversationActionProps's new
//
// This function is auto-generated.
func (p *aiConversationActionProps) setNew(new *types.AiConversation) *aiConversationActionProps {
	p.new = new
	return p
}

// setUpdate updates aiConversationActionProps's update
//
// This function is auto-generated.
func (p *aiConversationActionProps) setUpdate(update *types.AiConversation) *aiConversationActionProps {
	p.update = update
	return p
}

// setFilter updates aiConversationActionProps's filter
//
// This function is auto-generated.
func (p *aiConversationActionProps) setFilter(filter *types.AiConversationFilter) *aiConversationActionProps {
	p.filter = filter
	return p
}

// Serialize converts aiConversationActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p aiConversationActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.aiConversation != nil {
		m.Set("aiConversation.agentID", p.aiConversation.AgentID, true)
		m.Set("aiConversation.ID", p.aiConversation.ID, true)
	}
	if p.new != nil {
		m.Set("new.agentID", p.new.AgentID, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.agentID", p.update.AgentID, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.aiConversationID", p.filter.AiConversationID, true)
		m.Set("filter.agentID", p.filter.AgentID, true)
		m.Set("filter.deleted", p.filter.Deleted, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p aiConversationActionProps) Format(in string, err error) string {
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

	if p.aiConversation != nil {
		// replacement for "{{aiConversation}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{aiConversation}}",
			fns(
				p.aiConversation.AgentID,
				p.aiConversation.ID,
			),
		)
		pairs = append(pairs, "{{aiConversation.agentID}}", fns(p.aiConversation.AgentID))
		pairs = append(pairs, "{{aiConversation.ID}}", fns(p.aiConversation.ID))
	}

	if p.new != nil {
		// replacement for "{{new}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{new}}",
			fns(
				p.new.AgentID,
				p.new.ID,
			),
		)
		pairs = append(pairs, "{{new.agentID}}", fns(p.new.AgentID))
		pairs = append(pairs, "{{new.ID}}", fns(p.new.ID))
	}

	if p.update != nil {
		// replacement for "{{update}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{update}}",
			fns(
				p.update.AgentID,
				p.update.ID,
			),
		)
		pairs = append(pairs, "{{update.agentID}}", fns(p.update.AgentID))
		pairs = append(pairs, "{{update.ID}}", fns(p.update.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.AiConversationID,
				p.filter.AgentID,
				p.filter.Deleted,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{{filter.aiConversationID}}", fns(p.filter.AiConversationID))
		pairs = append(pairs, "{{filter.agentID}}", fns(p.filter.AgentID))
		pairs = append(pairs, "{{filter.deleted}}", fns(p.filter.Deleted))
		pairs = append(pairs, "{{filter.sort}}", fns(p.filter.Sort))
	}
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
func (a *aiConversationAction) String() string {
	var props = &aiConversationActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *aiConversationAction) ToAction() *actionlog.Action {
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

// AiConversationActionSearch returns "system:ai-conversation.search" action
//
// This function is auto-generated.
func AiConversationActionSearch(props ...*aiConversationActionProps) *aiConversationAction {
	a := &aiConversationAction{
		timestamp: time.Now(),
		resource:  "system:ai-conversation",
		action:    "search",
		log:       "searched for AI conversations",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AiConversationActionLookup returns "system:ai-conversation.lookup" action
//
// This function is auto-generated.
func AiConversationActionLookup(props ...*aiConversationActionProps) *aiConversationAction {
	a := &aiConversationAction{
		timestamp: time.Now(),
		resource:  "system:ai-conversation",
		action:    "lookup",
		log:       "looked-up for a {{aiConversation}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AiConversationActionCreate returns "system:ai-conversation.create" action
//
// This function is auto-generated.
func AiConversationActionCreate(props ...*aiConversationActionProps) *aiConversationAction {
	a := &aiConversationAction{
		timestamp: time.Now(),
		resource:  "system:ai-conversation",
		action:    "create",
		log:       "created {{aiConversation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AiConversationActionUpdate returns "system:ai-conversation.update" action
//
// This function is auto-generated.
func AiConversationActionUpdate(props ...*aiConversationActionProps) *aiConversationAction {
	a := &aiConversationAction{
		timestamp: time.Now(),
		resource:  "system:ai-conversation",
		action:    "update",
		log:       "updated {{aiConversation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AiConversationActionDelete returns "system:ai-conversation.delete" action
//
// This function is auto-generated.
func AiConversationActionDelete(props ...*aiConversationActionProps) *aiConversationAction {
	a := &aiConversationAction{
		timestamp: time.Now(),
		resource:  "system:ai-conversation",
		action:    "delete",
		log:       "deleted {{aiConversation}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AiConversationActionUndelete returns "system:ai-conversation.undelete" action
//
// This function is auto-generated.
func AiConversationActionUndelete(props ...*aiConversationActionProps) *aiConversationAction {
	a := &aiConversationAction{
		timestamp: time.Now(),
		resource:  "system:ai-conversation",
		action:    "undelete",
		log:       "undeleted {{aiConversation}}",
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

// AiConversationErrGeneric returns "system:ai-conversation.generic" as *errors.Error
//
// This function is auto-generated.
func AiConversationErrGeneric(mm ...*aiConversationActionProps) *errors.Error {
	var p = &aiConversationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:ai-conversation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(aiConversationLogMetaKey{}, "{err}"),
		errors.Meta(aiConversationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "ai-conversation.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AiConversationErrNotFound returns "system:ai-conversation.notFound" as *errors.Error
//
// This function is auto-generated.
func AiConversationErrNotFound(mm ...*aiConversationActionProps) *errors.Error {
	var p = &aiConversationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("AI conversation not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:ai-conversation"),

		errors.Meta(aiConversationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "ai-conversation.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AiConversationErrInvalidID returns "system:ai-conversation.invalidID" as *errors.Error
//
// This function is auto-generated.
func AiConversationErrInvalidID(mm ...*aiConversationActionProps) *errors.Error {
	var p = &aiConversationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:ai-conversation"),

		errors.Meta(aiConversationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "ai-conversation.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AiConversationErrStaleData returns "system:ai-conversation.staleData" as *errors.Error
//
// This function is auto-generated.
func AiConversationErrStaleData(mm ...*aiConversationActionProps) *errors.Error {
	var p = &aiConversationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("stale data", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:ai-conversation"),

		errors.Meta(aiConversationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "ai-conversation.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AiConversationErrNotAllowedToCreate returns "system:ai-conversation.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func AiConversationErrNotAllowedToCreate(mm ...*aiConversationActionProps) *errors.Error {
	var p = &aiConversationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create AI conversations", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:ai-conversation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(aiConversationLogMetaKey{}, "could not create AI conversations; insufficient permissions"),
		errors.Meta(aiConversationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "ai-conversation.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AiConversationErrNotAllowedToSearch returns "system:ai-conversation.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func AiConversationErrNotAllowedToSearch(mm ...*aiConversationActionProps) *errors.Error {
	var p = &aiConversationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list AI conversations", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:ai-conversation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(aiConversationLogMetaKey{}, "could not search or list AI conversations; insufficient permissions"),
		errors.Meta(aiConversationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "ai-conversation.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AiConversationErrNotAllowedToRead returns "system:ai-conversation.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func AiConversationErrNotAllowedToRead(mm ...*aiConversationActionProps) *errors.Error {
	var p = &aiConversationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this AI conversation", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:ai-conversation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(aiConversationLogMetaKey{}, "could not read {{aiConversation}}; insufficient permissions"),
		errors.Meta(aiConversationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "ai-conversation.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AiConversationErrNotAllowedToUpdate returns "system:ai-conversation.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func AiConversationErrNotAllowedToUpdate(mm ...*aiConversationActionProps) *errors.Error {
	var p = &aiConversationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this AI conversation", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:ai-conversation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(aiConversationLogMetaKey{}, "could not update {{aiConversation}}; insufficient permissions"),
		errors.Meta(aiConversationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "ai-conversation.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AiConversationErrNotAllowedToDelete returns "system:ai-conversation.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func AiConversationErrNotAllowedToDelete(mm ...*aiConversationActionProps) *errors.Error {
	var p = &aiConversationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this AI conversation", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:ai-conversation"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(aiConversationLogMetaKey{}, "could not delete {{aiConversation}}; insufficient permissions"),
		errors.Meta(aiConversationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "ai-conversation.errors.notAllowedToDelete"),

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
func (svc aiConversation) recordAction(ctx context.Context, props *aiConversationActionProps, actionFn func(...*aiConversationActionProps) *aiConversationAction, err error) error {
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
		a.Description = props.Format(m.AsString(aiConversationLogMetaKey{}), err)

		if p, has := m[aiConversationPropsMetaKey{}]; has {
			a.Meta = p.(*aiConversationActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
