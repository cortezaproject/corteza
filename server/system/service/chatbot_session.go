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
	"github.com/crusttech/human/server/system/types"
)

type (
	chatbotSession struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        chatbotSessionAccessController
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
