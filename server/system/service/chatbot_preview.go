package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	pkgAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/types"
)

// previewCap caps the in-memory preview session count. Once exceeded, the
// oldest entries are dropped FIFO. Restart-safe: state is never persisted.
const previewCap = 64

type (
	// ChatbotPreviewSession is one transient admin-test session. Mirrors the
	// shape of types.ChatbotSession + ChatbotSessionStep + ChatbotSessionHandoff
	// kept in-memory. The Chatbot snapshot is the unsaved editor draft.
	ChatbotPreviewSession struct {
		ID             string
		Chatbot        types.Chatbot
		ConversationID uint64
		Status         string // active|handoff_requested|handoff_active|closed
		CurrentStep    int
		Steps          []*ChatbotPreviewStep
		Handoff        *ChatbotPreviewHandoff
		CreatedAt      time.Time
		LastSeen       time.Time

		mu sync.Mutex
	}

	ChatbotPreviewStep struct {
		ScenarioIndex int
		ScenarioID    string
		Status        string // active|complete|failed
		HookLog       []string
	}

	ChatbotPreviewHandoff struct {
		ID       string
		Status   string // requested|active|closed
		Started  time.Time
		Closed   *time.Time
	}

	chatbotPreview struct {
		mu      sync.Mutex
		order   []*ChatbotPreviewSession
		byID    map[string]*ChatbotPreviewSession
		cap     int
		bus     *observability.Bus
		runtime previewAgenticRunner
		agent   previewAgentLookup
		conv    previewConvService
		store   store.Storer
	}

	previewAgenticRunner interface {
		Run(ctx context.Context, req *runtime.AgentRequest) (*runtime.AgentResponse, error)
	}

	previewAgentLookup interface {
		FindByID(ctx context.Context, ID uint64) (*types.Agent, error)
	}

	previewConvService interface {
		Create(ctx context.Context, new *types.AiConversation) (*types.AiConversation, error)
		FindByID(ctx context.Context, ID uint64) (*types.AiConversation, error)
		Update(ctx context.Context, upd *types.AiConversation) (*types.AiConversation, error)
	}
)

// ChatbotPreview wires the in-memory preview store with the shared SSE bus.
// The REST controller owns AiConversation creation separately and passes the
// resulting conversationID to Open().
func ChatbotPreview() *chatbotPreview {
	return &chatbotPreview{
		byID: make(map[string]*ChatbotPreviewSession),
		cap:  previewCap,
		bus:  DefaultObsBus,
	}
}

// WithDeps wires the agentic runtime, agent svc, conv svc and store. Called
// once after all dependent services are constructed.
func (s *chatbotPreview) WithDeps(rt previewAgenticRunner, ag previewAgentLookup, cv previewConvService, st store.Storer) *chatbotPreview {
	s.runtime = rt
	s.agent = ag
	s.conv = cv
	s.store = st
	return s
}

// OpenWithConversation creates a fresh AiConversation under service identity
// and opens a preview session bound to it.
func (s *chatbotPreview) OpenWithConversation(ctx context.Context, cb types.Chatbot) (*ChatbotPreviewSession, *types.AiConversation, error) {
	if len(cb.Scenarios) == 0 {
		return nil, nil, ChatbotSessionErrChatbotNoScenarios()
	}
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())
	conv, err := s.conv.Create(svcCtx, &types.AiConversation{})
	if err != nil {
		return nil, nil, err
	}
	ps := s.Open(cb, conv.ID)
	return ps, conv, nil
}

// Start is idempotent: starts step 0 only if no step has been started yet.
func (s *chatbotPreview) Start(ps *ChatbotPreviewSession) {
	if s.CurrentStep(ps) == nil {
		s.StartStep(ps, 0)
	}
}

// SubmitMessage handles a user message turn. Returns typed errors for state
// mismatches; runtime execution is fire-and-forget via SSE.
func (s *chatbotPreview) SubmitMessage(ctx context.Context, ps *ChatbotPreviewSession, input string) error {
	switch ps.Status {
	case "handoff_requested", "handoff_active":
		s.appendUserMessage(ps, input)
		s.EmitUserMessage(ps, input)
		return nil
	case "active":
		// fall through
	default:
		return ChatbotSessionErrSessionNotActive()
	}

	step := s.CurrentStep(ps)
	if step == nil {
		return ChatbotSessionErrNoActiveStep()
	}
	scenario := s.scenarioAt(ps, step.ScenarioIndex)
	if scenario == nil {
		return ChatbotSessionErrScenarioOutOfRange()
	}
	if scenario.Type != "conversation" || scenario.AgentID == 0 {
		return ChatbotSessionErrScenarioNotMessage()
	}

	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())
	agent, err := s.agent.FindByID(svcCtx, scenario.AgentID)
	if err != nil || agent == nil {
		return ChatbotSessionErrAgentUnavailable()
	}
	if !agent.Invocation.System.Enabled || agent.Invocation.System.ServiceAccount == 0 {
		return ChatbotSessionErrAgentNotConfigured()
	}

	if conv, err := s.conv.FindByID(svcCtx, ps.ConversationID); err == nil && conv != nil && conv.AgentID == 0 {
		conv.AgentID = agent.ID
		_, _ = s.conv.Update(svcCtx, conv)
	}

	// Broadcast the user turn so live SSE subscribers (operator inbox) see
	// the visitor's input during the bot phase, not just the agent's reply.
	// Mirrors the widget-side behavior in chatbot_session.SubmitMessage.
	s.EmitUserMessage(ps, input)

	saCtx := ImpersonateServiceAccount(context.Background(), s.store, agent.Invocation.System.ServiceAccount)
	cid := ps.ConversationID
	agentID := agent.ID
	go func() {
		resp, err := s.runtime.Run(saCtx, &runtime.AgentRequest{
			AgentID:        agentID,
			Input:          input,
			ConversationID: cid,
		})
		if err != nil {
			s.emit(cid, "agent_error", map[string]any{"error": err.Error()})
		} else if resp != nil && resp.Output != "" {
			s.emit(cid, "token", map[string]any{"text": resp.Output})
		}
		s.emit(cid, "done", map[string]any{})
	}()
	return nil
}

// SubmitForm validates and processes a form scenario submission. The first
// return is per-field validation errors (non-empty means the form was rejected
// without state change); the second is any unexpected error.
func (s *chatbotPreview) SubmitForm(ctx context.Context, ps *ChatbotPreviewSession, fields map[string]string) (map[string]string, error) {
	if ps.Status != "active" {
		return nil, ChatbotSessionErrSessionNotActive()
	}
	step := s.CurrentStep(ps)
	if step == nil {
		return nil, ChatbotSessionErrNoActiveStep()
	}
	scenario := s.scenarioAt(ps, step.ScenarioIndex)
	if scenario == nil {
		return nil, ChatbotSessionErrScenarioOutOfRange()
	}
	if scenario.Type != "form" {
		return nil, ChatbotSessionErrScenarioNotForm()
	}

	errs := ValidateChatbotFormFields(scenario.Config, fields)
	if len(errs) > 0 {
		s.EmitFormError(ps, scenario.ID, errs)
		return errs, nil
	}

	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())
	s.appendConvMessage(svcCtx, ps.ConversationID, types.AiConversationMessage{
		Role:    "user",
		Content: EncodeChatbotFormSubmission(fields),
	})

	s.FinalizeStep(ps)
	s.StartStep(ps, step.ScenarioIndex+1)
	return nil, nil
}

// SubmitConsent records the visitor's accept/reject decision on a consent
// step. Accept advances the preview; reject closes the preview session. The
// decision itself is not persisted (preview is in-memory).
func (s *chatbotPreview) SubmitConsent(ps *ChatbotPreviewSession, accepted bool) error {
	if ps.Status != "active" {
		return ChatbotSessionErrSessionNotActive()
	}
	step := s.CurrentStep(ps)
	if step == nil {
		return ChatbotSessionErrNoActiveStep()
	}
	scenario := s.scenarioAt(ps, step.ScenarioIndex)
	if scenario == nil {
		return ChatbotSessionErrScenarioOutOfRange()
	}
	if scenario.Type != "consent" {
		return ChatbotSessionErrScenarioNotConsent()
	}

	if !accepted {
		return s.CloseSession(ps)
	}
	s.FinalizeStep(ps)
	s.StartStep(ps, step.ScenarioIndex+1)
	return nil
}

// AdvanceStep finalizes the current step and starts the next. Closes any
// open handoff first.
func (s *chatbotPreview) AdvanceStep(ps *ChatbotPreviewSession) error {
	if ps.Handoff != nil && ps.Handoff.Status != "closed" {
		s.CloseHandoff(ps)
	}
	step := s.CurrentStep(ps)
	if step == nil {
		return ChatbotSessionErrNoActiveStep()
	}
	s.FinalizeStep(ps)
	s.StartStep(ps, step.ScenarioIndex+1)
	return nil
}

// CloseSession finalizes any active step and closes the session.
func (s *chatbotPreview) CloseSession(ps *ChatbotPreviewSession) error {
	if step := s.CurrentStep(ps); step != nil {
		s.FinalizeStep(ps)
	}
	s.Close(ps, "")
	return nil
}

// AcceptHandoff activates a requested handoff. Idempotent if already active.
func (s *chatbotPreview) AcceptHandoff(ps *ChatbotPreviewSession, operator string) error {
	if ps.Handoff == nil {
		return ChatbotSessionErrHandoffNotFound()
	}
	if ps.Handoff.Status == "active" {
		return nil
	}
	if ps.Handoff.Status != "requested" {
		return ChatbotSessionErrHandoffNotRequested()
	}
	s.ActivateHandoff(ps, operator)
	return nil
}

// SendOperatorMessage persists an operator message into the conversation and
// emits operator_message. Requires the handoff to be active.
func (s *chatbotPreview) SendOperatorMessage(ctx context.Context, ps *ChatbotPreviewSession, msg, operator string) error {
	if ps.Handoff == nil || ps.Handoff.Status != "active" {
		return ChatbotSessionErrHandoffNotActive()
	}
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())
	// Persist operator attribution as a structured field; the historical
	// load path reconstructs the same shape the SSE `operator_message` event
	// already carries (content + operator name).
	s.appendConvMessage(svcCtx, ps.ConversationID, types.AiConversationMessage{
		Role:     "assistant",
		Content:  msg,
		Operator: operator,
	})
	s.EmitOperatorMessage(ps, msg, operator)
	return nil
}

func (s *chatbotPreview) scenarioAt(ps *ChatbotPreviewSession, idx int) *types.ChatbotScenario {
	if idx < 0 || idx >= len(ps.Chatbot.Scenarios) {
		return nil
	}
	return &ps.Chatbot.Scenarios[idx]
}

func (s *chatbotPreview) appendConvMessage(ctx context.Context, convID uint64, msg types.AiConversationMessage) {
	if s.conv == nil {
		return
	}
	conv, err := s.conv.FindByID(ctx, convID)
	if err != nil || conv == nil {
		return
	}
	conv.Messages = append(conv.Messages, msg)
	_, _ = s.conv.Update(ctx, conv)
}

func (s *chatbotPreview) appendUserMessage(ps *ChatbotPreviewSession, content string) {
	svcCtx := pkgAuth.SetIdentityToContext(context.Background(), pkgAuth.ServiceUser())
	s.appendConvMessage(svcCtx, ps.ConversationID, types.AiConversationMessage{
		Role:    "user",
		Content: content,
	})
}

// Open allocates a new preview session. The Chatbot snapshot is taken by
// value and stored unchanged. Caller is responsible for calling StartStep
// after the SSE stream has been registered so the first step_start event
// isn't dropped before the listener attaches.
func (s *chatbotPreview) Open(cb types.Chatbot, convID uint64) *ChatbotPreviewSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	ps := &ChatbotPreviewSession{
		ID:             previewID(),
		Chatbot:        cb,
		ConversationID: convID,
		Status:         "active",
		CreatedAt:      time.Now(),
		LastSeen:       time.Now(),
	}
	s.order = append(s.order, ps)
	s.byID[ps.ID] = ps

	// FIFO trim: drop the oldest entries until at cap.
	for len(s.order) > s.cap {
		dropped := s.order[0]
		s.order = s.order[1:]
		delete(s.byID, dropped.ID)
		// Notify any listeners; they can clean up downstream resources.
		s.emit(dropped.ConversationID, "session_closed", map[string]any{
			"reason": "evicted",
		})
	}
	return ps
}

// Find returns the session by ID and refreshes LastSeen. nil if missing.
func (s *chatbotPreview) Find(id string) *ChatbotPreviewSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	ps, ok := s.byID[id]
	if !ok {
		return nil
	}
	ps.LastSeen = time.Now()
	return ps
}

// List returns preview sessions filtered by chatbotID (0 = all) and status
// (nil/empty = all). Result ordered by CreatedAt desc.
func (s *chatbotPreview) List(chatbotID uint64, status []string) []*ChatbotPreviewSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]*ChatbotPreviewSession, 0, len(s.order))
	for _, ps := range s.order {
		if chatbotID != 0 && ps.Chatbot.ID != chatbotID {
			continue
		}
		if len(status) > 0 {
			matched := false
			for _, st := range status {
				if ps.Status == st {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		out = append(out, ps)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

// StartStep records the active step, logs any automation hook that *would*
// fire (without executing it), and broadcasts step_start. If the index is
// past the last scenario, the session is closed.
func (s *chatbotPreview) StartStep(ps *ChatbotPreviewSession, scenarioIndex int) {
	if scenarioIndex >= len(ps.Chatbot.Scenarios) {
		s.Close(ps, "completed")
		return
	}
	scenario := &ps.Chatbot.Scenarios[scenarioIndex]

	ps.mu.Lock()
	step := &ChatbotPreviewStep{
		ScenarioIndex: scenarioIndex,
		ScenarioID:    scenario.ID,
		Status:        "active",
	}
	if scenario.Automation.Before.Automation != "" {
		entry := fmt.Sprintf("before: would invoke %s (async=%v)", scenario.Automation.Before.Automation, scenario.Automation.Before.Async)
		step.HookLog = append(step.HookLog, entry)
	}
	ps.Steps = append(ps.Steps, step)
	ps.CurrentStep = scenarioIndex
	ps.mu.Unlock()

	if scenario.Automation.Before.Automation != "" {
		s.emit(ps.ConversationID, "step_hook_skipped", map[string]any{
			"phase":    "before",
			"resource": scenario.Automation.Before.Automation,
			"async":    scenario.Automation.Before.Async,
		})
	}

	s.emit(ps.ConversationID, "step_start", map[string]any{
		"scenarioID":    scenario.ID,
		"scenarioIndex": scenarioIndex,
		"type":          scenario.Type,
		"config":        ChatbotScenarioConfig(scenario),
		"stepID":        fmt.Sprintf("preview-%s-%d", ps.ID, scenarioIndex),
	})
}

// FinalizeStep marks the current step complete, logs the after-hook intent,
// and emits step_complete. No-op if there is no active step.
func (s *chatbotPreview) FinalizeStep(ps *ChatbotPreviewSession) {
	ps.mu.Lock()
	if len(ps.Steps) == 0 {
		ps.mu.Unlock()
		return
	}
	step := ps.Steps[len(ps.Steps)-1]
	step.Status = "complete"
	idx := step.ScenarioIndex
	id := step.ScenarioID
	var afterHook *types.ChatbotAutomationHook
	if idx >= 0 && idx < len(ps.Chatbot.Scenarios) {
		afterHook = &ps.Chatbot.Scenarios[idx].Automation.After
	}
	if afterHook != nil && afterHook.Automation != "" {
		step.HookLog = append(step.HookLog,
			fmt.Sprintf("after: would invoke %s (async=%v)", afterHook.Automation, afterHook.Async))
	}
	ps.mu.Unlock()

	if afterHook != nil && afterHook.Automation != "" {
		s.emit(ps.ConversationID, "step_hook_skipped", map[string]any{
			"phase":    "after",
			"resource": afterHook.Automation,
			"async":    afterHook.Async,
		})
	}

	s.emit(ps.ConversationID, "step_complete", map[string]any{
		"scenarioID":    id,
		"scenarioIndex": idx,
	})
}

// FailStep marks the current step as failed.
func (s *chatbotPreview) FailStep(ps *ChatbotPreviewSession, reason string) {
	ps.mu.Lock()
	if len(ps.Steps) > 0 {
		step := ps.Steps[len(ps.Steps)-1]
		step.Status = "failed"
		step.HookLog = append(step.HookLog, "failed: "+reason)
	}
	ps.mu.Unlock()
}

// CurrentStep returns the active step (last entry in the audit list) or nil.
func (s *chatbotPreview) CurrentStep(ps *ChatbotPreviewSession) *ChatbotPreviewStep {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if len(ps.Steps) == 0 {
		return nil
	}
	return ps.Steps[len(ps.Steps)-1]
}

// RequestHandoff transitions to handoff_requested.
func (s *chatbotPreview) RequestHandoff(ps *ChatbotPreviewSession) *ChatbotPreviewHandoff {
	ps.mu.Lock()
	h := &ChatbotPreviewHandoff{ID: previewID(), Status: "requested", Started: time.Now()}
	ps.Handoff = h
	ps.Status = "handoff_requested"
	ps.mu.Unlock()
	s.emit(ps.ConversationID, "handoff_requested", map[string]any{"handoffID": h.ID})
	return h
}

// ActivateHandoff transitions to handoff_active.
func (s *chatbotPreview) ActivateHandoff(ps *ChatbotPreviewSession, operator string) {
	ps.mu.Lock()
	if ps.Handoff != nil {
		ps.Handoff.Status = "active"
	}
	ps.Status = "handoff_active"
	hid := ""
	if ps.Handoff != nil {
		hid = ps.Handoff.ID
	}
	ps.mu.Unlock()
	s.emit(ps.ConversationID, "handoff_active", map[string]any{
		"handoffID": hid,
		"operator":  operator,
	})
}

// CloseHandoff resumes the session at the same step.
func (s *chatbotPreview) CloseHandoff(ps *ChatbotPreviewSession) {
	ps.mu.Lock()
	hid := ""
	if ps.Handoff != nil {
		n := time.Now()
		ps.Handoff.Status = "closed"
		ps.Handoff.Closed = &n
		hid = ps.Handoff.ID
	}
	ps.Status = "active"
	ps.mu.Unlock()
	s.emit(ps.ConversationID, "handoff_complete", map[string]any{"handoffID": hid})
}

// Close terminates the session.
func (s *chatbotPreview) Close(ps *ChatbotPreviewSession, reason string) {
	ps.mu.Lock()
	ps.Status = "closed"
	ps.mu.Unlock()
	details := map[string]any{}
	if reason != "" {
		details["reason"] = reason
	}
	s.emit(ps.ConversationID, "session_closed", details)
}

// EmitUserMessage emits a user_message event on the SSE bus for operator
// console visibility while a handoff is in progress.
func (s *chatbotPreview) EmitUserMessage(ps *ChatbotPreviewSession, content string) {
	s.emit(ps.ConversationID, "user_message", map[string]any{"content": content})
}

// EmitOperatorMessage emits an operator_message event for the user widget.
func (s *chatbotPreview) EmitOperatorMessage(ps *ChatbotPreviewSession, content, operator string) {
	s.emit(ps.ConversationID, "operator_message", map[string]any{
		"content":  content,
		"operator": operator,
	})
}

// EmitFormError emits a form_error event for the current scenario.
func (s *chatbotPreview) EmitFormError(ps *ChatbotPreviewSession, scenarioID string, errs map[string]string) {
	s.emit(ps.ConversationID, "form_error", map[string]any{
		"scenarioID": scenarioID,
		"errors":     errs,
	})
}

func (s *chatbotPreview) emit(convID uint64, name string, details map[string]any) {
	if s.bus == nil {
		return
	}
	s.bus.EmitEvent(observability.AgentEvent{
		ConversationID: strconv.FormatUint(convID, 10),
		Event:          name,
		Details:        details,
		Timestamp:      time.Now(),
	})
}

// previewID returns a URL-safe random session identifier.
func previewID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
