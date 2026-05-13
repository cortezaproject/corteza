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

	svc.actionlog.Record(ctx, &actionlog.Action{
		Resource: "chatbot-session-handoff",
		Action:   "request",
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

	// Update session status
	return svc.UpdateStatus(ctx, h.SessionID, "handoff_active")
}

// CloseHandoff close handoff (Variant B: close session; Variant A: advance step)
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
		Paging: filter.Paging{Limit: 1000},
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
	var step *types.ChatbotSessionStep

	existing, _, _ := store.SearchChatbotSessionSteps(ctx, svc.store, types.ChatbotSessionStepFilter{
		SessionID: sessionID,
		Check: func(s *types.ChatbotSessionStep) (bool, error) {
			return s.ScenarioIndex == scenarioIndex, nil
		},
	})
	var existingStep *types.ChatbotSessionStep
	if len(existing) > 0 {
		existingStep = existing[0]
	}

	if existingStep == nil {
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

		// Run before automation if present — blocks; failure aborts step
		if scenario.Automation.Before != "" {
			if err := svc.invokeAutomation(ctx, scenario.Automation.Before); err != nil {
				step.Status = "failed"
				_ = store.UpdateChatbotSessionStep(ctx, svc.store, step)
				return step, fmt.Errorf("before automation: %w", err)
			}
		}
	} else {
		step = existingStep
	}

	// Mark step complete
	if err := svc.CompleteStep(ctx, step.ID); err != nil {
		return step, fmt.Errorf("complete step: %w", err)
	}
	step.Status = "complete"

	// Run after automation — log errors only, don't abort
	if scenario.Automation.After != "" {
		if err := svc.invokeAutomation(ctx, scenario.Automation.After); err != nil {
			svc.actionlog.Record(ctx, &actionlog.Action{
				Resource: "chatbot-session-step",
				Action:   "automation-error",
				Error:    err.Error(),
			})
		}
	}

	return step, nil
}

// invokeAutomation executes an automation by resource identifier synchronously under service identity.
// Resource format: corteza::automation:ng-automation/{ID}
func (svc *chatbotSession) invokeAutomation(ctx context.Context, resourceID string) error {
	if automationService.DefaultNgAutomation == nil {
		return fmt.Errorf("automation service not available")
	}

	prefix := automationTypes.NgAutomationResourceType + "/"
	idStr := strings.TrimPrefix(resourceID, prefix)
	automationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || automationID == 0 {
		return fmt.Errorf("invalid automation resource identifier: %s", resourceID)
	}

	svcCtx := pkgAuth.SetIdentityToContext(ctx, pkgAuth.ServiceUser())
	_, err = automationService.DefaultNgAutomation.ExecAndWait(svcCtx, automationID, automationTypes.NgAutomationExecParams{
		Input: &expr.Vars{},
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
