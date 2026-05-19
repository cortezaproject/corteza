package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/crusttech/human/server/system/agentic/observability"
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
		mu    sync.Mutex
		order []*ChatbotPreviewSession
		byID  map[string]*ChatbotPreviewSession
		cap   int
		bus   *observability.Bus
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
