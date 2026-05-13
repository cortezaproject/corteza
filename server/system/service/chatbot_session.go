package service

import (
	"context"
	"fmt"
	"time"

	automationService "github.com/crusttech/human/server/automation/service"
	automationTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/expr"
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
			return errors.New("not allowed to create chatbot session")
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

		svc.actionlog.Record(ctx, actionlog.New().
			Resource("chatbot-session").
			Action("create").
			ID(new.ID))

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

// FindStepBySession find current step in session
func (svc *chatbotSession) FindStepBySession(ctx context.Context, sessionID uint64) (*types.ChatbotSessionStep, error) {
	return store.LookupChatbotSessionStepBySessionID(ctx, svc.store, sessionID)
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

// RequestHandoff request handoff from AI to human
func (svc *chatbotSession) RequestHandoff(ctx context.Context, sessionID, stepID uint64) (*types.ChatbotSessionHandoff, error) {
	// Update session status
	if err := svc.UpdateStatus(ctx, sessionID, "handoff_requested"); err != nil {
		return nil, err
	}

	// Create handoff record
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

	svc.actionlog.Record(ctx, actionlog.New().
		Resource("chatbot-session-handoff").
		Action("request").
		ID(h.ID))

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

	// Update session status
	return svc.UpdateStatus(ctx, h.SessionID, "handoff_active")
}

// CloseHandoff close handoff (Variant B: close session; Variant A: advance step)
func (svc *chatbotSession) CloseHandoff(ctx context.Context, handoffID uint64) error {
	h, err := svc.FindHandoffByID(ctx, handoffID)
	if err != nil {
		return err
	}

	now := now()
	h.Status = "closed"
	h.ClosedAt = now
	h.UpdatedAt = now

	if err := store.UpdateChatbotSessionHandoff(ctx, svc.store, h); err != nil {
		return err
	}

	// Variant B: close session
	return svc.UpdateStatus(ctx, h.SessionID, "closed")
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
		Paging: types.Paging{Limit: 1000},
	}

	hh, _, err := store.SearchChatbotSessionHandoffs(ctx, svc.store, f)
	if err != nil {
		return nil, err
	}

	return hh, nil
}

// ExecuteStep executes a scenario step with before/after automation hooks.
// Full flow: beforeAutomation → mark complete → afterAutomation
// Caller is responsible for executing the scenario itself.
func (svc *chatbotSession) ExecuteStep(ctx context.Context, sessionID uint64, scenario *types.ChatbotScenario, conversationID uint64, scenarioIndex int, input string) (*types.ChatbotSessionStep, error) {
	// Get or create step
	var step *types.ChatbotSessionStep
	existingStep, _ := svc.FindStepBySession(ctx, sessionID)

	if existingStep == nil {
		// Create new step
		var err error
		step, err = svc.CreateStep(ctx, &types.ChatbotSessionStep{
			SessionID:      sessionID,
			ScenarioIndex:  scenarioIndex,
			ConversationID: conversationID,
			Status:         "active",
		})
		if err != nil {
			return nil, fmt.Errorf("create step: %w", err)
		}

		// 2. Run before automation if present
		if scenario.BeforeAutomationID != nil && *scenario.BeforeAutomationID != 0 {
			if err := svc.invokeAutomation(ctx, *scenario.BeforeAutomationID); err != nil {
				// Before automation failure = mark step failed and abort
				step.Status = "failed"
				_ = store.UpdateChatbotSessionStep(ctx, svc.store, step)
				return step, fmt.Errorf("before automation: %w", err)
			}
		}
	} else {
		step = existingStep
	}

	// 3. Mark step complete
	step.Status = "complete"
	if err := svc.CompleteStep(ctx, step.ID); err != nil {
		return step, fmt.Errorf("complete step: %w", err)
	}

	// 4. Run after automation if present (log errors, don't abort)
	if scenario.AfterAutomationID != nil && *scenario.AfterAutomationID != 0 {
		if err := svc.invokeAutomation(ctx, *scenario.AfterAutomationID); err != nil {
			svc.actionlog.Record(ctx, &actionlog.Action{
				Resource:    "chatbot-session-step",
				Action:      "automation-error",
				Error:       err.Error(),
			})
		}
	}

	// 5. SSE emission and next message handled by caller (widget)

	return step, nil
}

// invokeAutomation executes an automation by ID synchronously.
// Runs as service account under system context.
func (svc *chatbotSession) invokeAutomation(ctx context.Context, automationID uint64) error {
	if automationService.DefaultNgAutomation == nil {
		return fmt.Errorf("automation service not available")
	}

	svcCtx := auth.SetIdentityToContext(ctx, auth.ServiceUser())
	_, err := automationService.DefaultNgAutomation.ExecAndWait(svcCtx, automationID, automationTypes.NgAutomationExecParams{
		Input: &expr.Vars{},
	})
	return err
}
