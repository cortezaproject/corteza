package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/chatbot_actions.yaml

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/system/types"
	"strings"
	"time"
)

type (
	chatbotActionProps struct {
		chatbot *types.Chatbot
		new     *types.Chatbot
		update  *types.Chatbot
		filter  *types.ChatbotFilter
		diff    []*revisions.Change
		old     json.RawMessage
	}

	chatbotAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *chatbotActionProps
	}

	chatbotLogMetaKey   struct{}
	chatbotPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setChatbot updates chatbotActionProps's chatbot
//
// This function is auto-generated.
func (p *chatbotActionProps) setChatbot(chatbot *types.Chatbot) *chatbotActionProps {
	p.chatbot = chatbot
	return p
}

// setNew updates chatbotActionProps's new
//
// This function is auto-generated.
func (p *chatbotActionProps) setNew(new *types.Chatbot) *chatbotActionProps {
	p.new = new
	return p
}

// setUpdate updates chatbotActionProps's update
//
// This function is auto-generated.
func (p *chatbotActionProps) setUpdate(update *types.Chatbot) *chatbotActionProps {
	p.update = update
	return p
}

// setFilter updates chatbotActionProps's filter
//
// This function is auto-generated.
func (p *chatbotActionProps) setFilter(filter *types.ChatbotFilter) *chatbotActionProps {
	p.filter = filter
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *chatbotActionProps) setDiff(diff []*revisions.Change) *chatbotActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *chatbotActionProps) setOld(v any) *chatbotActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts chatbotActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p chatbotActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.chatbot != nil {
		m.Set("chatbot.handle", p.chatbot.Handle, true)
		m.Set("chatbot.ID", p.chatbot.ID, true)
	}
	if p.new != nil {
		m.Set("new.handle", p.new.Handle, true)
		m.Set("new.ID", p.new.ID, true)
	}
	if p.update != nil {
		m.Set("update.handle", p.update.Handle, true)
		m.Set("update.ID", p.update.ID, true)
	}
	if p.filter != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p chatbotActionProps) Format(in string, err error) string {
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

	if p.chatbot != nil {
		// replacement for "{{chatbot}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{chatbot}}",
			fns(
				p.chatbot.Handle,
				p.chatbot.ID,
			),
		)
		pairs = append(pairs, "{{chatbot.handle}}", fns(p.chatbot.Handle))
		pairs = append(pairs, "{{chatbot.ID}}", fns(p.chatbot.ID))
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

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(),
		)
	}
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
func (a *chatbotAction) String() string {
	var props = &chatbotActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *chatbotAction) ToAction() *actionlog.Action {
	resource := e.resource
	if e.props != nil && e.props.chatbot != nil {
		if r, ok := any(e.props.chatbot).(actionlog.RbacResourcer); ok {
			resource = r.RbacResource()
		}
	}
	return &actionlog.Action{
		Resource:    resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.Serialize(),
		Delta:       actionlog.Delta(e.props.diff),
		OldState:    actionlog.OldState(e.props.old),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// ChatbotActionSearch returns "system:chatbot.search" action
//
// This function is auto-generated.
func ChatbotActionSearch(props ...*chatbotActionProps) *chatbotAction {
	a := &chatbotAction{
		timestamp: time.Now(),
		resource:  "system:chatbot",
		action:    "search",
		log:       "searched for chatbots",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChatbotActionLookup returns "system:chatbot.lookup" action
//
// This function is auto-generated.
func ChatbotActionLookup(props ...*chatbotActionProps) *chatbotAction {
	a := &chatbotAction{
		timestamp: time.Now(),
		resource:  "system:chatbot",
		action:    "lookup",
		log:       "looked-up for a {{chatbot}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChatbotActionCreate returns "system:chatbot.create" action
//
// This function is auto-generated.
func ChatbotActionCreate(props ...*chatbotActionProps) *chatbotAction {
	a := &chatbotAction{
		timestamp: time.Now(),
		resource:  "system:chatbot",
		action:    "create",
		log:       "created {{chatbot}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChatbotActionUpdate returns "system:chatbot.update" action
//
// This function is auto-generated.
func ChatbotActionUpdate(props ...*chatbotActionProps) *chatbotAction {
	a := &chatbotAction{
		timestamp: time.Now(),
		resource:  "system:chatbot",
		action:    "update",
		log:       "updated {{chatbot}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChatbotActionDelete returns "system:chatbot.delete" action
//
// This function is auto-generated.
func ChatbotActionDelete(props ...*chatbotActionProps) *chatbotAction {
	a := &chatbotAction{
		timestamp: time.Now(),
		resource:  "system:chatbot",
		action:    "delete",
		log:       "deleted {{chatbot}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChatbotActionUndelete returns "system:chatbot.undelete" action
//
// This function is auto-generated.
func ChatbotActionUndelete(props ...*chatbotActionProps) *chatbotAction {
	a := &chatbotAction{
		timestamp: time.Now(),
		resource:  "system:chatbot",
		action:    "undelete",
		log:       "undeleted {{chatbot}}",
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

// ChatbotErrGeneric returns "system:chatbot.generic" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrGeneric(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:chatbot"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chatbotLogMetaKey{}, "{err}"),
		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrNotFound returns "system:chatbot.notFound" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrNotFound(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("chatbot not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:chatbot"),

		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrInvalidID returns "system:chatbot.invalidID" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrInvalidID(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:chatbot"),

		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrStaleData returns "system:chatbot.staleData" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrStaleData(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("chatbot was modified by someone else after you've opened it. Please refresh to see the latest updated version", nil),

		errors.Meta("type", "staleData"),
		errors.Meta("resource", "system:chatbot"),

		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.staleData"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrNotAllowedToCreate returns "system:chatbot.notAllowedToCreate" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrNotAllowedToCreate(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create a chatbot", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:chatbot"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chatbotLogMetaKey{}, "failed to create a chatbot; insufficient permissions"),
		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrNotAllowedToRead returns "system:chatbot.notAllowedToRead" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrNotAllowedToRead(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this chatbot", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "system:chatbot"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chatbotLogMetaKey{}, "failed to read {{chatbot.handle}}; insufficient permissions"),
		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrNotAllowedToSearch returns "system:chatbot.notAllowedToSearch" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrNotAllowedToSearch(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list chatbots", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "system:chatbot"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chatbotLogMetaKey{}, "failed to search or list chatbots; insufficient permissions"),
		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrNotAllowedToUpdate returns "system:chatbot.notAllowedToUpdate" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrNotAllowedToUpdate(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this chatbot", nil),

		errors.Meta("type", "notAllowedToUpdate"),
		errors.Meta("resource", "system:chatbot"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chatbotLogMetaKey{}, "failed to update {{chatbot.handle}}; insufficient permissions"),
		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.notAllowedToUpdate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrNotAllowedToDelete returns "system:chatbot.notAllowedToDelete" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrNotAllowedToDelete(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to delete this chatbot", nil),

		errors.Meta("type", "notAllowedToDelete"),
		errors.Meta("resource", "system:chatbot"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chatbotLogMetaKey{}, "failed to delete {{chatbot.handle}}; insufficient permissions"),
		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.notAllowedToDelete"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrConversationScenarioMissingAgent returns "system:chatbot.conversationScenarioMissingAgent" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrConversationScenarioMissingAgent(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("conversation scenario is missing an agent", nil),

		errors.Meta("type", "conversationScenarioMissingAgent"),
		errors.Meta("resource", "system:chatbot"),

		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.conversationScenarioMissingAgent"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotErrInvalidHandle returns "system:chatbot.invalidHandle" as *errors.Error
//
// This function is auto-generated.
func ChatbotErrInvalidHandle(mm ...*chatbotActionProps) *errors.Error {
	var p = &chatbotActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid handle", nil),

		errors.Meta("type", "invalidHandle"),
		errors.Meta("resource", "system:chatbot"),

		errors.Meta(chatbotPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot.errors.invalidHandle"),

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
func (svc chatbot) recordAction(ctx context.Context, props *chatbotActionProps, actionFn func(...*chatbotActionProps) *chatbotAction, err error, diff ...any) error {
	if len(diff) == 2 && diff[0] != nil {
		props.setDiff(actionlog.DiffResourceState(diff[0], diff[1]))
		props.setOld(diff[0])
	}
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
		a.Description = props.Format(m.AsString(chatbotLogMetaKey{}), err)

		if p, has := m[chatbotPropsMetaKey{}]; has {
			a.Meta = p.(*chatbotActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
