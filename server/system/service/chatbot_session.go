package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	automationService "github.com/crusttech/human/server/automation/service"
	automationTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	pkgAuth "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/agentic/runtime"
	"github.com/crusttech/human/server/system/types"
)

type (
	chatbotSession struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        chatbotSessionAccessController

		// Orchestration deps; wired post-init via WithDeps. nil-safe — methods
		// that need them no-op or error when unset.
		bus     *observability.Bus
		runtime previewAgenticRunner
		agent   previewAgentLookup
		conv    previewConvService
	}

	chatbotSessionAccessController interface {
		CanCreateChatbotSession(ctx context.Context) bool
		CanReadChatbotSession(ctx context.Context, s *types.ChatbotSession) bool
		CanUpdateChatbotSession(ctx context.Context, s *types.ChatbotSession) bool
		CanDeleteChatbotSession(ctx context.Context, s *types.ChatbotSession) bool
	}
)

func ChatbotSession() *chatbotSession {
	return &chatbotSession{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

// Create new chatbot session
func (svc *chatbotSession) Create(ctx context.Context, new *types.ChatbotSession) (*types.ChatbotSession, error) {
	err := func() error {
		if !svc.ac.CanCreateChatbotSession(ctx) {
			return fmt.Errorf("not allowed to create chatbot session")
		}

		if new.ID == 0 {
			new.ID = nextID()
		}
		new.CreatedAt = *now()
		new.Status = "active"
		new.CurrentStep = 0

		if err := store.CreateChatbotSession(ctx, svc.store, new); err != nil {
			return err
		}

		svc.actionlog.Record(ctx, &actionlog.Action{
			Resource: "chatbot-session",
			Action:   "create",
		})

		return nil
	}()

	return new, err
}

// FindByID load session by ID
func (svc *chatbotSession) FindByID(ctx context.Context, ID uint64) (*types.ChatbotSession, error) {
	return store.LookupChatbotSessionByID(ctx, svc.store, ID)
}

// FindByChatbotID find session by chatbot ID
func (svc *chatbotSession) FindByChatbotID(ctx context.Context, chatbotID uint64) (*types.ChatbotSession, error) {
	return store.LookupChatbotSessionByChatbotID(ctx, svc.store, chatbotID)
}

// Search sessions by filter
func (svc *chatbotSession) Search(ctx context.Context, f types.ChatbotSessionFilter) (types.ChatbotSessionSet, types.ChatbotSessionFilter, error) {
	return store.SearchChatbotSessions(ctx, svc.store, f)
}

// FindStepsBySession returns all steps for a session in order
func (svc *chatbotSession) FindStepsBySession(ctx context.Context, sessionID uint64) (types.ChatbotSessionStepSet, error) {
	set, _, err := store.SearchChatbotSessionSteps(ctx, svc.store, types.ChatbotSessionStepFilter{
		SessionID: sessionID,
	})
	return set, err
}

// UpdateStatus update session status
func (svc *chatbotSession) UpdateStatus(ctx context.Context, id uint64, status string) error {
	s, err := svc.FindByID(ctx, id)
	if err != nil {
		return err
	}

	s.Status = status
	s.UpdatedAt = now()

	return store.UpdateChatbotSession(ctx, svc.store, s)
}

// UpdateCurrentStep sets the session pointer to the given scenario index.
func (svc *chatbotSession) UpdateCurrentStep(ctx context.Context, id uint64, idx int) error {
	s, err := svc.FindByID(ctx, id)
	if err != nil {
		return err
	}
	s.CurrentStep = idx
	s.UpdatedAt = now()
	return store.UpdateChatbotSession(ctx, svc.store, s)
}

// CreateStep create new session step
func (svc *chatbotSession) CreateStep(ctx context.Context, step *types.ChatbotSessionStep) (*types.ChatbotSessionStep, error) {
	if step.ID == 0 {
		step.ID = nextID()
	}
	step.CreatedAt = *now()
	step.Status = "active"

	if err := store.CreateChatbotSessionStep(ctx, svc.store, step); err != nil {
		return nil, err
	}

	return step, nil
}

// FindStepByID find step by ID
func (svc *chatbotSession) FindStepByID(ctx context.Context, id uint64) (*types.ChatbotSessionStep, error) {
	return store.LookupChatbotSessionStepByID(ctx, svc.store, id)
}

// FindStepBySession find current active step in session (most recent active row).
func (svc *chatbotSession) FindStepBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionStep, error) {
	steps, _, err := store.SearchChatbotSessionSteps(ctx, svc.store, types.ChatbotSessionStepFilter{
		SessionID: sessionID,
		Status:    "active",
		Paging:    filter.Paging{Limit: 1},
	})
	if err != nil {
		return nil, err
	}
	if len(steps) == 0 {
		return store.LookupChatbotSessionStepBySessionID(ctx, svc.store, sessionID)
	}
	return steps[0], nil
}

// FindStepBySessionAndIndex returns the step row for a (session, scenarioIndex) pair if one exists.
func (svc *chatbotSession) FindStepBySessionAndIndex(ctx context.Context, sessionID uint64, scenarioIndex int) (*types.ChatbotSessionStep, error) {
	steps, _, err := store.SearchChatbotSessionSteps(ctx, svc.store, types.ChatbotSessionStepFilter{
		SessionID: sessionID,
		Check: func(s *types.ChatbotSessionStep) (bool, error) {
			return s.ScenarioIndex == scenarioIndex, nil
		},
	})
	if err != nil {
		return nil, err
	}
	if len(steps) == 0 {
		return nil, nil
	}
	return steps[0], nil
}

// CompleteStep mark step as complete
func (svc *chatbotSession) CompleteStep(ctx context.Context, stepID uint64) error {
	step, err := svc.FindStepByID(ctx, stepID)
	if err != nil {
		return err
	}

	step.Status = "complete"
	step.UpdatedAt = now()

	return store.UpdateChatbotSessionStep(ctx, svc.store, step)
}

// StartStep creates a new step row (if not already present) for the given scenario
// index and runs the scenario's before-automation. Returns the step in "active"
// state on success, or "failed" if a synchronous before-automation errored.
func (svc *chatbotSession) StartStep(
	ctx context.Context,
	sessionID uint64,
	scenario *types.ChatbotScenario,
	conversationID uint64,
	scenarioIndex int,
	vars *expr.Vars,
) (*types.ChatbotSessionStep, error) {
	step, err := svc.FindStepBySessionAndIndex(ctx, sessionID, scenarioIndex)
	if err != nil {
		return nil, fmt.Errorf("lookup step: %w", err)
	}

	if step == nil {
		step, err = svc.CreateStep(ctx, &types.ChatbotSessionStep{
			SessionID:      sessionID,
			ScenarioIndex:  scenarioIndex,
			ConversationID: conversationID,
			Status:         "active",
		})
		if err != nil {
			return nil, fmt.Errorf("create step: %w", err)
		}
	}

	_ = svc.UpdateCurrentStep(ctx, sessionID, scenarioIndex)

	if scenario != nil {
		if err := svc.runHook(ctx, scenario.Automation.Before, vars, "before"); err != nil {
			step.Status = "failed"
			_ = store.UpdateChatbotSessionStep(ctx, svc.store, step)
			return step, fmt.Errorf("before automation: %w", err)
		}
	}

	return step, nil
}

// FinalizeStep marks the step complete and runs the scenario's after-automation.
// After-automation errors are recorded but do not bubble up.
func (svc *chatbotSession) FinalizeStep(
	ctx context.Context,
	stepID uint64,
	scenario *types.ChatbotScenario,
	vars *expr.Vars,
) error {
	if err := svc.CompleteStep(ctx, stepID); err != nil {
		return fmt.Errorf("complete step: %w", err)
	}

	if scenario != nil {
		if err := svc.runHook(ctx, scenario.Automation.After, vars, "after"); err != nil {
			svc.actionlog.Record(ctx, &actionlog.Action{
				Resource: "chatbot-session-step",
				Action:   "automation-error",
				Error:    err.Error(),
			})
		}
	}
	return nil
}

// runHook dispatches a single automation hook. Sync hooks block and surface
// errors to the caller; async hooks fire-and-forget and always return nil.
// An empty hook (no automation resource) is a no-op.
func (svc *chatbotSession) runHook(ctx context.Context, hook types.ChatbotAutomationHook, vars *expr.Vars, phase string) error {
	if hook.Automation == "" {
		return nil
	}
	if hook.Async {
		svc.invokeAutomationAsync(ctx, hook.Automation, vars, phase)
		return nil
	}
	return svc.invokeAutomation(ctx, hook.Automation, vars)
}

// RequestHandoff request handoff from AI to human
func (svc *chatbotSession) RequestHandoff(ctx context.Context, sessionID, stepID uint64) (*types.ChatbotSessionHandoff, error) {
	if err := svc.UpdateStatus(ctx, sessionID, "handoff_requested"); err != nil {
		return nil, err
	}

	h := &types.ChatbotSessionHandoff{
		ID:          nextID(),
		SessionID:   sessionID,
		StepID:      stepID,
		Status:      "requested",
		InitiatedAt: *now(),
		CreatedAt:   *now(),
	}

	if err := store.CreateChatbotSessionHandoff(ctx, svc.store, h); err != nil {
		return nil, err
	}

	svc.actionlog.Record(ctx, &actionlog.Action{
		Resource: "chatbot-session-handoff",
		Action:   "request",
	})

	svc.runHandoffHook(ctx, sessionID, h, "onRequested", func(hd *types.ChatbotHandoff) types.ChatbotAutomationHook {
		return hd.Automation.OnRequested
	})

	return h, nil
}

// ActivateHandoff mark handoff as active
func (svc *chatbotSession) ActivateHandoff(ctx context.Context, handoffID uint64) error {
	h, err := svc.FindHandoffByID(ctx, handoffID)
	if err != nil {
		return err
	}

	h.Status = "active"
	h.UpdatedAt = now()

	if err := store.UpdateChatbotSessionHandoff(ctx, svc.store, h); err != nil {
		return err
	}

	if err := svc.UpdateStatus(ctx, h.SessionID, "handoff_active"); err != nil {
		return err
	}

	svc.runHandoffHook(ctx, h.SessionID, h, "onAccepted", func(hd *types.ChatbotHandoff) types.ChatbotAutomationHook {
		return hd.Automation.OnAccepted
	})

	return nil
}

// runHandoffHook resolves the chatbot's handoff automation config for the
// session and dispatches the requested hook. No-op when the chatbot, session,
// or hook is missing, or when handoff is disabled on the chatbot.
func (svc *chatbotSession) runHandoffHook(
	ctx context.Context,
	sessionID uint64,
	h *types.ChatbotSessionHandoff,
	phase string,
	pick func(*types.ChatbotHandoff) types.ChatbotAutomationHook,
) {
	sess, err := svc.FindByID(ctx, sessionID)
	if err != nil || sess == nil {
		return
	}
	cb, err := store.LookupChatbotByID(ctx, svc.store, sess.ChatbotID)
	if err != nil || cb == nil || !cb.Handoff.Enabled {
		return
	}
	hook := pick(&cb.Handoff)
	if hook.Automation == "" {
		return
	}

	vars := &expr.Vars{}
	_ = vars.Set("chatbotID", strconv.FormatUint(cb.ID, 10))
	_ = vars.Set("sessionID", strconv.FormatUint(sessionID, 10))
	if h != nil {
		_ = vars.Set("handoffID", strconv.FormatUint(h.ID, 10))
		_ = vars.Set("stepID", strconv.FormatUint(h.StepID, 10))
	}

	if err := svc.runHook(ctx, hook, vars, phase); err != nil {
		svc.actionlog.Record(ctx, &actionlog.Action{
			Resource: "chatbot-session-handoff",
			Action:   "automation-error",
			Error:    fmt.Sprintf("%s: %v", phase, err),
		})
	}
}

// CloseHandoff terminates the handoff and resumes the session at the same step.
// Status flips back to "active" so the AI loop can pick up further user messages.
func (svc *chatbotSession) CloseHandoff(ctx context.Context, handoffID uint64) error {
	h, err := svc.FindHandoffByID(ctx, handoffID)
	if err != nil {
		return err
	}

	n := now()
	h.Status = "closed"
	h.ClosedAt = n
	h.UpdatedAt = n

	if err := store.UpdateChatbotSessionHandoff(ctx, svc.store, h); err != nil {
		return err
	}

	return svc.UpdateStatus(ctx, h.SessionID, "active")
}

// FindHandoffByID find handoff by ID
func (svc *chatbotSession) FindHandoffByID(ctx context.Context, id uint64) (*types.ChatbotSessionHandoff, error) {
	return store.LookupChatbotSessionHandoffByID(ctx, svc.store, id)
}

// FindHandoffBySession find active handoff for session
func (svc *chatbotSession) FindHandoffBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionHandoff, error) {
	return store.LookupChatbotSessionHandoffBySessionID(ctx, svc.store, sessionID)
}

// FindLiveHandoffBySession returns the live (requested or active) handoff for a
// session, preferring active over requested. Returns nil with no error when the
// session has no live handoff.
func (svc *chatbotSession) FindLiveHandoffBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionHandoff, error) {
	hh, _, err := store.SearchChatbotSessionHandoffs(ctx, svc.store, types.ChatbotSessionHandoffFilter{
		SessionID: sessionID,
		Paging:    filter.Paging{Limit: 1000},
	})
	if err != nil {
		return nil, err
	}
	var requested *types.ChatbotSessionHandoff
	for _, h := range hh {
		if h.ClosedAt != nil {
			continue
		}
		switch h.Status {
		case "active":
			return h, nil
		case "requested":
			if requested == nil {
				requested = h
			}
		}
	}
	return requested, nil
}

// AppendOperatorMessage appends an operator turn to the conversation tied to
// the handoff's step. The handoff must be in "active" state. Returns the
// conversation ID so the caller can emit SSE without re-resolving the step.
// Used by the legacy operator console; the fat SendOperatorMessage below is
// the canonical path for the aggregated session API.
func (svc *chatbotSession) AppendOperatorMessage(ctx context.Context, handoffID uint64, message string) (uint64, error) {
	h, err := svc.FindHandoffByID(ctx, handoffID)
	if err != nil {
		return 0, err
	}
	if h == nil {
		return 0, fmt.Errorf("handoff not found")
	}
	if h.Status != "active" {
		return 0, fmt.Errorf("handoff not active")
	}

	step, err := svc.FindStepByID(ctx, h.StepID)
	if err != nil {
		return 0, fmt.Errorf("step lookup: %w", err)
	}
	if step == nil || step.ConversationID == 0 {
		return 0, fmt.Errorf("step has no conversation")
	}

	conv, err := store.LookupAiConversationByID(ctx, svc.store, step.ConversationID)
	if err != nil {
		return 0, fmt.Errorf("conversation lookup: %w", err)
	}
	if conv == nil {
		return 0, fmt.Errorf("conversation not found")
	}

	conv.Messages = append(conv.Messages, types.AiConversationMessage{
		Role:    "assistant",
		Content: "[operator] " + message,
	})
	if err := store.UpdateAiConversation(ctx, svc.store, conv); err != nil {
		return 0, err
	}
	return conv.ID, nil
}

// FindPendingHandoffs find all pending handoffs (requested or active)
func (svc *chatbotSession) FindPendingHandoffs(ctx context.Context) ([]*types.ChatbotSessionHandoff, error) {
	f := types.ChatbotSessionHandoffFilter{
		Status: "requested",
		Paging: filter.Paging{Limit: 1000},
	}

	hh, _, err := store.SearchChatbotSessionHandoffs(ctx, svc.store, f)
	if err != nil {
		return nil, err
	}

	return hh, nil
}

// invokeAutomationAsync dispatches an automation in the background and logs
// failures via actionlog. The caller does not wait and errors do not surface
// to the step lifecycle.
func (svc *chatbotSession) invokeAutomationAsync(_ context.Context, resourceID string, vars *expr.Vars, phase string) {
	// Snapshot the identity so the goroutine doesn't outlive the request ctx.
	svcCtx := pkgAuth.SetIdentityToContext(context.Background(), pkgAuth.ServiceUser())
	go func() {
		if err := svc.invokeAutomation(svcCtx, resourceID, vars); err != nil {
			svc.actionlog.Record(svcCtx, &actionlog.Action{
				Resource: "chatbot-session-step",
				Action:   "automation-error",
				Error:    fmt.Sprintf("%s (async): %v", phase, err),
			})
		}
	}()
}

// invokeAutomation executes an automation by resource identifier synchronously under service identity.
// Resource format: corteza::automation:ng-automation/{ID}
func (svc *chatbotSession) invokeAutomation(ctx context.Context, resourceID string, vars *expr.Vars) error {
	if automationService.DefaultNgAutomation == nil {
		return fmt.Errorf("automation service not available")
	}

	prefix := automationTypes.NgAutomationResourceType + "/"
	idStr := strings.TrimPrefix(resourceID, prefix)
	automationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || automationID == 0 {
		return fmt.Errorf("invalid automation resource identifier: %s", resourceID)
	}

	if vars == nil {
		vars = &expr.Vars{}
	}

	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())
	_, err = automationService.DefaultNgAutomation.ExecAndWait(svcCtx, automationID, automationTypes.NgAutomationExecParams{
		Input: vars,
	})
	return err
}

// WithDeps wires orchestration deps (bus, runtime, agent svc, conv svc).
// Called once post-init from service.Initialize.
func (svc *chatbotSession) WithDeps(bus *observability.Bus, rt previewAgenticRunner, ag previewAgentLookup, cv previewConvService) *chatbotSession {
	svc.bus = bus
	svc.runtime = rt
	svc.agent = ag
	svc.conv = cv
	return svc
}

// Open validates the chatbot has scenarios, creates a fresh AiConversation
// and ChatbotSession in "active" state. JWT mint stays in REST.
func (svc *chatbotSession) Open(ctx context.Context, cb *types.Chatbot) (*types.ChatbotSession, *types.AiConversation, error) {
	if cb == nil || len(cb.Scenarios) == 0 {
		return nil, nil, ChatbotSessionErrChatbotNoScenarios()
	}
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())

	conv, err := svc.conv.Create(svcCtx, &types.AiConversation{})
	if err != nil {
		return nil, nil, err
	}

	session, err := svc.Create(svcCtx, &types.ChatbotSession{
		ChatbotID: cb.ID,
		Status:    "active",
	})
	if err != nil {
		return nil, nil, err
	}
	return session, conv, nil
}

// Start runs the scenario's before-automation, persists the new step row, and
// broadcasts step_start. If the index is beyond the last scenario, emits
// session_closed instead.
func (svc *chatbotSession) Start(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64, scenarioIndex int, vars *expr.Vars) {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())

	if scenarioIndex >= len(cb.Scenarios) {
		_ = svc.UpdateStatus(svcCtx, sessionID, "closed")
		svc.emit(convID, "session_closed", map[string]any{})
		return
	}

	scenario := &cb.Scenarios[scenarioIndex]
	step, err := svc.StartStep(svcCtx, sessionID, scenario, convID, scenarioIndex, vars)
	if err != nil {
		svc.emit(convID, "step_error", map[string]any{
			"scenarioID":    scenario.ID,
			"scenarioIndex": scenarioIndex,
			"phase":         "before",
			"error":         err.Error(),
		})
		return
	}

	svc.stateInit(svcCtx, sessionID, scenario.ID, scenario.Type, convID)

	svc.emit(convID, "step_start", map[string]any{
		"scenarioID":    scenario.ID,
		"scenarioIndex": scenarioIndex,
		"type":          scenario.Type,
		"config":        ChatbotScenarioConfig(scenario),
		"stepID":        strconv.FormatUint(step.ID, 10),
	})
}

// SubmitMessage handles a user message turn for a persisted widget session.
// Active session → agent runtime; handoff_active/requested → bypass agent.
func (svc *chatbotSession) SubmitMessage(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64, input string) error {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())

	session, err := svc.FindByID(svcCtx, sessionID)
	if err != nil || session == nil {
		return ChatbotSessionErrSessionNotFound()
	}

	switch session.Status {
	case "handoff_requested", "handoff_active":
		msg := types.AiConversationMessage{Role: "user", Content: input}
		svc.appendConversationMessage(svcCtx, convID, msg)
		if sid := svc.scenarioIDForActiveStep(svcCtx, cb, sessionID); sid != "" {
			svc.stateAppendHistory(svcCtx, sessionID, sid, msg)
		}
		svc.emit(convID, "user_message", map[string]any{"content": input})
		return nil
	case "active":
		// fall through
	default:
		return ChatbotSessionErrSessionNotActive()
	}

	step, err := svc.FindStepBySession(svcCtx, sessionID)
	if err != nil || step == nil {
		return ChatbotSessionErrNoActiveStep()
	}
	if step.ScenarioIndex < 0 || step.ScenarioIndex >= len(cb.Scenarios) {
		return ChatbotSessionErrScenarioOutOfRange()
	}
	scenario := &cb.Scenarios[step.ScenarioIndex]
	if scenario.Type != "conversation" || scenario.AgentID == 0 {
		return ChatbotSessionErrScenarioNotMessage()
	}

	agent, err := svc.resolveAgent(ctx, scenario.AgentID)
	if err != nil || agent == nil {
		return ChatbotSessionErrAgentUnavailable()
	}
	if !agent.Invocation.System.Enabled || agent.Invocation.System.ServiceAccount == 0 {
		return ChatbotSessionErrAgentNotConfigured()
	}

	if conv, err := store.LookupAiConversationByID(svcCtx, svc.store, convID); err == nil && conv != nil && conv.AgentID == 0 {
		conv.AgentID = agent.ID
		_ = store.UpdateAiConversation(svcCtx, svc.store, conv)
	}

	svc.stateAppendHistory(svcCtx, sessionID, scenario.ID, types.AiConversationMessage{
		Role:    "user",
		Content: input,
	})

	saCtx := ImpersonateServiceAccount(context.Background(), svc.store, agent.Invocation.System.ServiceAccount)
	agentID := agent.ID
	scenarioID := scenario.ID
	go func() {
		resp, err := svc.runtime.Run(saCtx, &runtime.AgentRequest{
			AgentID:        agentID,
			Input:          input,
			ConversationID: convID,
		})
		if err != nil {
			svc.emit(convID, "agent_error", map[string]any{"error": err.Error()})
		} else if resp != nil && resp.Output != "" {
			svc.stateAppendHistory(saCtx, sessionID, scenarioID, types.AiConversationMessage{
				Role:    "assistant",
				Content: resp.Output,
			})
			svc.emit(convID, "token", map[string]any{"text": resp.Output})
		}
		svc.emit(convID, "done", map[string]any{})
	}()
	return nil
}

// SubmitForm validates and processes a form scenario submission. Non-empty
// returned errors map signals rejected input without state change.
func (svc *chatbotSession) SubmitForm(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64, fields map[string]string) (map[string]string, error) {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())

	session, err := svc.FindByID(svcCtx, sessionID)
	if err != nil || session == nil {
		return nil, ChatbotSessionErrSessionNotFound()
	}
	if session.Status != "active" {
		return nil, ChatbotSessionErrSessionNotActive()
	}

	step, err := svc.FindStepBySession(svcCtx, sessionID)
	if err != nil || step == nil {
		return nil, ChatbotSessionErrNoActiveStep()
	}
	if step.ScenarioIndex < 0 || step.ScenarioIndex >= len(cb.Scenarios) {
		return nil, ChatbotSessionErrScenarioOutOfRange()
	}
	scenario := &cb.Scenarios[step.ScenarioIndex]
	if scenario.Type != "form" {
		return nil, ChatbotSessionErrScenarioNotForm()
	}

	errs := ValidateChatbotFormFields(scenario.Config, fields)
	if len(errs) > 0 {
		svc.emit(convID, "form_error", map[string]any{
			"scenarioID": scenario.ID,
			"errors":     errs,
		})
		return errs, nil
	}

	msg := types.AiConversationMessage{
		Role:    "user",
		Content: EncodeChatbotFormSubmission(fields),
	}
	svc.appendConversationMessage(svcCtx, convID, msg)
	svc.stateRecordFormSubmit(svcCtx, sessionID, scenario.ID, fields)

	vars := ChatbotVarsFromMap(fields)
	if err := svc.FinalizeStep(svcCtx, step.ID, scenario, vars); err != nil {
		return nil, err
	}
	svc.emit(convID, "step_complete", map[string]any{
		"scenarioID":    scenario.ID,
		"scenarioIndex": step.ScenarioIndex,
	})

	svc.Start(svcCtx, cb, sessionID, convID, step.ScenarioIndex+1, vars)
	return nil, nil
}

// AdvanceStep finalizes the current step (after-automation) and starts the
// next scenario. Closes any open handoff first.
func (svc *chatbotSession) AdvanceStep(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64) error {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())

	if h, _ := svc.FindHandoffBySession(svcCtx, sessionID); h != nil && h.ClosedAt == nil {
		_ = svc.CloseHandoff(svcCtx, h.ID)
		svc.emit(convID, "handoff_complete", map[string]any{
			"handoffID": strconv.FormatUint(h.ID, 10),
		})
	}

	step, err := svc.FindStepBySession(svcCtx, sessionID)
	if err != nil || step == nil {
		return ChatbotSessionErrNoActiveStep()
	}
	if step.ScenarioIndex < 0 || step.ScenarioIndex >= len(cb.Scenarios) {
		return ChatbotSessionErrScenarioOutOfRange()
	}
	scenario := &cb.Scenarios[step.ScenarioIndex]

	if err := svc.FinalizeStep(svcCtx, step.ID, scenario, nil); err != nil {
		return err
	}
	svc.emit(convID, "step_complete", map[string]any{
		"scenarioID":    scenario.ID,
		"scenarioIndex": step.ScenarioIndex,
	})

	svc.Start(svcCtx, cb, sessionID, convID, step.ScenarioIndex+1, nil)
	return nil
}

// CloseSession finalizes the current step (if any) and marks the session closed.
func (svc *chatbotSession) CloseSession(ctx context.Context, cb *types.Chatbot, sessionID, convID uint64) error {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())

	if step, err := svc.FindStepBySession(svcCtx, sessionID); err == nil && step != nil {
		if step.ScenarioIndex >= 0 && step.ScenarioIndex < len(cb.Scenarios) {
			scenario := &cb.Scenarios[step.ScenarioIndex]
			_ = svc.FinalizeStep(svcCtx, step.ID, scenario, nil)
			svc.emit(convID, "step_complete", map[string]any{
				"scenarioID":    scenario.ID,
				"scenarioIndex": step.ScenarioIndex,
			})
		}
	}

	if err := svc.UpdateStatus(svcCtx, sessionID, "closed"); err != nil {
		return err
	}
	svc.emit(convID, "session_closed", map[string]any{})
	return nil
}

// AcceptHandoffByID activates a requested handoff and emits handoff_active.
// Idempotent when already active. operatorID/operatorName are recorded on
// the per-scenario session State for the active conversation step.
func (svc *chatbotSession) AcceptHandoffByID(ctx context.Context, cb *types.Chatbot, sessionID, convID, handoffID uint64, operatorID uint64, operatorName string) error {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())

	h, err := svc.FindHandoffByID(svcCtx, handoffID)
	if err != nil || h == nil {
		return ChatbotSessionErrHandoffNotFound()
	}
	if h.Status == "active" {
		return nil
	}
	if h.Status != "requested" {
		return ChatbotSessionErrHandoffNotRequested()
	}
	if err := svc.ActivateHandoff(svcCtx, handoffID); err != nil {
		return err
	}

	if cb != nil {
		if sid := svc.scenarioIDForActiveStep(svcCtx, cb, sessionID); sid != "" {
			svc.stateRecordHandoffOperator(svcCtx, sessionID, sid, operatorID, operatorName)
		}
	}

	svc.emit(convID, "handoff_active", map[string]any{
		"handoffID":    strconv.FormatUint(handoffID, 10),
		"operatorID":   strconv.FormatUint(operatorID, 10),
		"operatorName": operatorName,
	})
	return nil
}

// SendOperatorMessage persists an operator message and emits operator_message.
// Requires the handoff to be active.
func (svc *chatbotSession) SendOperatorMessage(ctx context.Context, cb *types.Chatbot, sessionID, convID, handoffID uint64, msg, operator string) error {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())

	h, err := svc.FindHandoffByID(svcCtx, handoffID)
	if err != nil || h == nil {
		return ChatbotSessionErrHandoffNotFound()
	}
	if h.Status != "active" {
		return ChatbotSessionErrHandoffNotActive()
	}

	m := types.AiConversationMessage{
		Role:    "assistant",
		Content: "[operator] " + msg,
	}
	svc.appendConversationMessage(svcCtx, convID, m)
	if cb != nil {
		if sid := svc.scenarioIDForActiveStep(svcCtx, cb, sessionID); sid != "" {
			svc.stateAppendHistory(svcCtx, sessionID, sid, m)
		}
	}
	svc.emit(convID, "operator_message", map[string]any{
		"content":  msg,
		"operator": operator,
	})
	return nil
}

// RequestHandoffPublic wraps RequestHandoff and emits handoff_requested for
// widget-driven flows.
func (svc *chatbotSession) RequestHandoffPublic(ctx context.Context, sessionID, convID uint64, reason string) (*types.ChatbotSessionHandoff, error) {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())

	step, err := svc.FindStepBySession(svcCtx, sessionID)
	if err != nil || step == nil {
		return nil, ChatbotSessionErrNoActiveStep()
	}

	h, err := svc.RequestHandoff(svcCtx, sessionID, step.ID)
	if err != nil {
		return nil, err
	}

	svc.emit(convID, "handoff_requested", map[string]any{
		"handoffID": strconv.FormatUint(h.ID, 10),
		"reason":    reason,
	})
	return h, nil
}

// CloseHandoffPublic wraps CloseHandoff with the operator_message-style emission.
func (svc *chatbotSession) CloseHandoffPublic(ctx context.Context, convID, handoffID uint64) error {
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())
	if err := svc.CloseHandoff(svcCtx, handoffID); err != nil {
		return err
	}
	svc.emit(convID, "handoff_complete", map[string]any{
		"handoffID": strconv.FormatUint(handoffID, 10),
	})
	return nil
}

// emit broadcasts an SSE event for the conversation. No-op when bus is nil.
func (svc *chatbotSession) emit(convID uint64, name string, details map[string]any) {
	if svc.bus == nil {
		return
	}
	svc.bus.EmitEvent(observability.AgentEvent{
		ConversationID: strconv.FormatUint(convID, 10),
		Event:          name,
		Details:        details,
	})
}

// appendConversationMessage appends a message to the AiConversation. Best-effort.
func (svc *chatbotSession) appendConversationMessage(ctx context.Context, convID uint64, msg types.AiConversationMessage) {
	conv, err := store.LookupAiConversationByID(ctx, svc.store, convID)
	if err != nil || conv == nil {
		return
	}
	conv.Messages = append(conv.Messages, msg)
	_ = store.UpdateAiConversation(ctx, svc.store, conv)
}

// scenarioIDForActiveStep returns the scenario ID of the active step for a
// session, or "" if it cannot be resolved. Used to key State writes.
func (svc *chatbotSession) scenarioIDForActiveStep(ctx context.Context, cb *types.Chatbot, sessionID uint64) string {
	step, err := svc.FindStepBySession(ctx, sessionID)
	if err != nil || step == nil {
		return ""
	}
	if step.ScenarioIndex < 0 || step.ScenarioIndex >= len(cb.Scenarios) {
		return ""
	}
	return cb.Scenarios[step.ScenarioIndex].ID
}

// stateInit ensures a per-scenario state entry exists and seeds the typed
// sub-struct. Idempotent — re-running on an existing entry leaves data intact.
func (svc *chatbotSession) stateInit(ctx context.Context, sessionID uint64, scenarioID, scenarioType string, convID uint64) {
	if sessionID == 0 || scenarioID == "" {
		return
	}
	session, err := svc.FindByID(ctx, sessionID)
	if err != nil || session == nil {
		return
	}
	entry := session.State.Upsert(scenarioID, scenarioType)
	switch scenarioType {
	case "conversation":
		if entry.Conversation == nil {
			entry.Conversation = &types.ChatbotConversationStepState{ConversationID: convID}
		} else if entry.Conversation.ConversationID == 0 {
			entry.Conversation.ConversationID = convID
		}
	case "form":
		if entry.Form == nil {
			entry.Form = &types.ChatbotFormStepState{}
		}
	}
	_ = store.UpdateChatbotSession(ctx, svc.store, session)
}

// stateAppendHistory mirrors a conversation message into the per-scenario
// State.Conversation.History. Best-effort.
func (svc *chatbotSession) stateAppendHistory(ctx context.Context, sessionID uint64, scenarioID string, msg types.AiConversationMessage) {
	if sessionID == 0 || scenarioID == "" {
		return
	}
	session, err := svc.FindByID(ctx, sessionID)
	if err != nil || session == nil {
		return
	}
	entry := session.State.ForScenario(scenarioID)
	if entry == nil {
		return
	}
	if entry.Conversation == nil {
		entry.Conversation = &types.ChatbotConversationStepState{}
	}
	entry.Conversation.History = append(entry.Conversation.History, msg)
	_ = store.UpdateChatbotSession(ctx, svc.store, session)
}

// stateRecordFormSubmit persists submitted form fields into per-scenario state.
func (svc *chatbotSession) stateRecordFormSubmit(ctx context.Context, sessionID uint64, scenarioID string, fields map[string]string) {
	if sessionID == 0 || scenarioID == "" {
		return
	}
	session, err := svc.FindByID(ctx, sessionID)
	if err != nil || session == nil {
		return
	}
	entry := session.State.Upsert(scenarioID, "form")
	if entry.Form == nil {
		entry.Form = &types.ChatbotFormStepState{}
	}
	entry.Form.Fields = fields
	entry.Form.Submitted = true
	_ = store.UpdateChatbotSession(ctx, svc.store, session)
}

// stateRecordHandoffOperator stores operator identity on the active
// conversation step's handoff state. No-op if no active conversation entry.
func (svc *chatbotSession) stateRecordHandoffOperator(ctx context.Context, sessionID uint64, scenarioID string, operatorID uint64, operatorName string) {
	if sessionID == 0 || scenarioID == "" {
		return
	}
	session, err := svc.FindByID(ctx, sessionID)
	if err != nil || session == nil {
		return
	}
	entry := session.State.ForScenario(scenarioID)
	if entry == nil || entry.Conversation == nil {
		return
	}
	entry.Conversation.Handoff = &types.ChatbotConversationHandoffState{
		OperatorID:   operatorID,
		OperatorName: operatorName,
	}
	_ = store.UpdateChatbotSession(ctx, svc.store, session)
}

// resolveAgent loads an Agent by ID under service identity.
func (svc *chatbotSession) resolveAgent(ctx context.Context, agentID uint64) (*types.Agent, error) {
	if agentID == 0 {
		return nil, ChatbotSessionErrAgentUnavailable()
	}
	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())
	return store.LookupAgentByID(svcCtx, svc.store, agentID)
}

// loadChatbotSession loads a ChatbotSession by ID for RBAC checks.
func loadChatbotSession(ctx context.Context, s store.Storer, ID uint64) (*types.ChatbotSession, error) {
	if ID == 0 {
		return nil, fmt.Errorf("invalid chatbot session ID")
	}
	return store.LookupChatbotSessionByID(ctx, s, ID)
}

// loadChatbotSessionStep loads a ChatbotSessionStep by ID for RBAC checks.
func loadChatbotSessionStep(ctx context.Context, s store.Storer, ID uint64) (*types.ChatbotSessionStep, error) {
	if ID == 0 {
		return nil, fmt.Errorf("invalid chatbot session step ID")
	}
	return store.LookupChatbotSessionStepByID(ctx, s, ID)
}

// loadChatbotSessionHandoff loads a ChatbotSessionHandoff by ID for RBAC checks.
func loadChatbotSessionHandoff(ctx context.Context, s store.Storer, ID uint64) (*types.ChatbotSessionHandoff, error) {
	if ID == 0 {
		return nil, fmt.Errorf("invalid chatbot session handoff ID")
	}
	return store.LookupChatbotSessionHandoffByID(ctx, s, ID)
}
