package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/chatbot_session_actions.yaml

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
	chatbotSessionActionProps struct {
		session *types.ChatbotSession
		handoff *types.ChatbotSessionHandoff
		search  *types.ChatbotSessionFilter
		diff    []*revisions.Change
		old     json.RawMessage
	}

	chatbotSessionAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *chatbotSessionActionProps
	}

	chatbotSessionLogMetaKey   struct{}
	chatbotSessionPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setSession updates chatbotSessionActionProps's session
//
// This function is auto-generated.
func (p *chatbotSessionActionProps) setSession(session *types.ChatbotSession) *chatbotSessionActionProps {
	p.session = session
	return p
}

// setHandoff updates chatbotSessionActionProps's handoff
//
// This function is auto-generated.
func (p *chatbotSessionActionProps) setHandoff(handoff *types.ChatbotSessionHandoff) *chatbotSessionActionProps {
	p.handoff = handoff
	return p
}

// setSearch updates chatbotSessionActionProps's search
//
// This function is auto-generated.
func (p *chatbotSessionActionProps) setSearch(search *types.ChatbotSessionFilter) *chatbotSessionActionProps {
	p.search = search
	return p
}

// setDiff stores the field-level delta of the changed resource for the action log.
//
// This function is auto-generated.
func (p *chatbotSessionActionProps) setDiff(diff []*revisions.Change) *chatbotSessionActionProps {
	p.diff = diff
	return p
}

// setOld stores the full JSON snapshot of the resource before the update.
//
// This function is auto-generated.
func (p *chatbotSessionActionProps) setOld(v any) *chatbotSessionActionProps {
	p.old, _ = json.Marshal(v)
	return p
}

// Serialize converts chatbotSessionActionProps to actionlog.Meta
//
// This function is auto-generated.
func (p chatbotSessionActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.session != nil {
		m.Set("session.ID", p.session.ID, true)
		m.Set("session.ChatbotID", p.session.ChatbotID, true)
		m.Set("session.Status", p.session.Status, true)
	}
	if p.handoff != nil {
		m.Set("handoff.ID", p.handoff.ID, true)
		m.Set("handoff.SessionID", p.handoff.SessionID, true)
		m.Set("handoff.Status", p.handoff.Status, true)
	}
	if p.search != nil {
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
func (p chatbotSessionActionProps) Format(in string, err error) string {
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

	if p.session != nil {
		// replacement for "{{session}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{session}}",
			fns(
				p.session.ID,
				p.session.ChatbotID,
				p.session.Status,
			),
		)
		pairs = append(pairs, "{{session.ID}}", fns(p.session.ID))
		pairs = append(pairs, "{{session.ChatbotID}}", fns(p.session.ChatbotID))
		pairs = append(pairs, "{{session.Status}}", fns(p.session.Status))
	}

	if p.handoff != nil {
		// replacement for "{{handoff}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{handoff}}",
			fns(
				p.handoff.ID,
				p.handoff.SessionID,
				p.handoff.Status,
			),
		)
		pairs = append(pairs, "{{handoff.ID}}", fns(p.handoff.ID))
		pairs = append(pairs, "{{handoff.SessionID}}", fns(p.handoff.SessionID))
		pairs = append(pairs, "{{handoff.Status}}", fns(p.handoff.Status))
	}

	if p.search != nil {
		// replacement for "{{search}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{search}}",
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
func (a *chatbotSessionAction) String() string {
	var props = &chatbotSessionActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *chatbotSessionAction) ToAction() *actionlog.Action {
	resource := e.resource
	var resourceProjectID uint64

	return &actionlog.Action{
		Resource:          resource,
		ResourceProjectID: resourceProjectID,
		Action:            e.action,
		Severity:          e.severity,
		Description:       e.String(),
		Meta:              e.props.Serialize(),
		Delta:             actionlog.Delta(e.props.diff),
		OldState:          actionlog.OldState(e.props.old),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// ChatbotSessionActionSearch returns "system:chatbot-session.search" action
//
// This function is auto-generated.
func ChatbotSessionActionSearch(props ...*chatbotSessionActionProps) *chatbotSessionAction {
	a := &chatbotSessionAction{
		timestamp: time.Now(),
		resource:  "system:chatbot-session",
		action:    "search",
		log:       "searched for chatbot sessions",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChatbotSessionActionLookup returns "system:chatbot-session.lookup" action
//
// This function is auto-generated.
func ChatbotSessionActionLookup(props ...*chatbotSessionActionProps) *chatbotSessionAction {
	a := &chatbotSessionAction{
		timestamp: time.Now(),
		resource:  "system:chatbot-session",
		action:    "lookup",
		log:       "looked-up for a {{session}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChatbotSessionActionCreate returns "system:chatbot-session.create" action
//
// This function is auto-generated.
func ChatbotSessionActionCreate(props ...*chatbotSessionActionProps) *chatbotSessionAction {
	a := &chatbotSessionAction{
		timestamp: time.Now(),
		resource:  "system:chatbot-session",
		action:    "create",
		log:       "created {{session}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChatbotSessionActionUpdate returns "system:chatbot-session.update" action
//
// This function is auto-generated.
func ChatbotSessionActionUpdate(props ...*chatbotSessionActionProps) *chatbotSessionAction {
	a := &chatbotSessionAction{
		timestamp: time.Now(),
		resource:  "system:chatbot-session",
		action:    "update",
		log:       "updated {{session}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// ChatbotSessionActionDelete returns "system:chatbot-session.delete" action
//
// This function is auto-generated.
func ChatbotSessionActionDelete(props ...*chatbotSessionActionProps) *chatbotSessionAction {
	a := &chatbotSessionAction{
		timestamp: time.Now(),
		resource:  "system:chatbot-session",
		action:    "delete",
		log:       "deleted {{session}}",
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

// ChatbotSessionErrGeneric returns "system:chatbot-session.generic" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrGeneric(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:chatbot-session"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(chatbotSessionLogMetaKey{}, "{err}"),
		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrChatbotNoScenarios returns "system:chatbot-session.chatbotNoScenarios" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrChatbotNoScenarios(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("chatbot has no scenarios", nil),

		errors.Meta("type", "chatbotNoScenarios"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.chatbotNoScenarios"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrSessionNotActive returns "system:chatbot-session.sessionNotActive" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrSessionNotActive(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("session not active", nil),

		errors.Meta("type", "sessionNotActive"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.sessionNotActive"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrSessionNotFound returns "system:chatbot-session.sessionNotFound" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrSessionNotFound(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("session not found", nil),

		errors.Meta("type", "sessionNotFound"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.sessionNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrNoActiveStep returns "system:chatbot-session.noActiveStep" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrNoActiveStep(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("no active step", nil),

		errors.Meta("type", "noActiveStep"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.noActiveStep"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrScenarioOutOfRange returns "system:chatbot-session.scenarioOutOfRange" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrScenarioOutOfRange(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("scenario out of range", nil),

		errors.Meta("type", "scenarioOutOfRange"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.scenarioOutOfRange"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrScenarioNotMessage returns "system:chatbot-session.scenarioNotMessage" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrScenarioNotMessage(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("current step does not accept messages", nil),

		errors.Meta("type", "scenarioNotMessage"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.scenarioNotMessage"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrScenarioNotForm returns "system:chatbot-session.scenarioNotForm" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrScenarioNotForm(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("current step is not a form", nil),

		errors.Meta("type", "scenarioNotForm"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.scenarioNotForm"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrScenarioNotConsent returns "system:chatbot-session.scenarioNotConsent" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrScenarioNotConsent(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("current step is not a consent step", nil),

		errors.Meta("type", "scenarioNotConsent"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.scenarioNotConsent"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrAgentUnavailable returns "system:chatbot-session.agentUnavailable" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrAgentUnavailable(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("agent unavailable", nil),

		errors.Meta("type", "agentUnavailable"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.agentUnavailable"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrAgentNotConfigured returns "system:chatbot-session.agentNotConfigured" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrAgentNotConfigured(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("agent not configured for system invocation", nil),

		errors.Meta("type", "agentNotConfigured"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.agentNotConfigured"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrHandoffNotFound returns "system:chatbot-session.handoffNotFound" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrHandoffNotFound(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handoff not found", nil),

		errors.Meta("type", "handoffNotFound"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.handoffNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrHandoffNotRequested returns "system:chatbot-session.handoffNotRequested" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrHandoffNotRequested(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handoff not in requested state", nil),

		errors.Meta("type", "handoffNotRequested"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.handoffNotRequested"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrHandoffNotActive returns "system:chatbot-session.handoffNotActive" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrHandoffNotActive(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("handoff not active", nil),

		errors.Meta("type", "handoffNotActive"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.handoffNotActive"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// ChatbotSessionErrBadHandoffID returns "system:chatbot-session.badHandoffID" as *errors.Error
//
// This function is auto-generated.
func ChatbotSessionErrBadHandoffID(mm ...*chatbotSessionActionProps) *errors.Error {
	var p = &chatbotSessionActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("bad handoff id", nil),

		errors.Meta("type", "badHandoffID"),
		errors.Meta("resource", "system:chatbot-session"),

		errors.Meta(chatbotSessionPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "chatbot-session.errors.badHandoffID"),

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
func (svc chatbotSession) recordAction(ctx context.Context, props *chatbotSessionActionProps, actionFn func(...*chatbotSessionActionProps) *chatbotSessionAction, err error, diff ...any) error {
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
		a.Description = props.Format(m.AsString(chatbotSessionLogMetaKey{}), err)

		if p, has := m[chatbotSessionPropsMetaKey{}]; has {
			a.Meta = p.(*chatbotSessionActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
