package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/automation"
)

type (
	// Handles automation triggers storing and loading
	automationTrigger struct {
		logger         *zap.Logger
		triggerManager automationTriggerManager
	}

	automationTriggerManager interface {
		FindTriggerByID(context.Context, uint64) (*automation.Trigger, error)
		FindTriggers(context.Context, automation.TriggerFilter) (automation.TriggerSet, automation.TriggerFilter, error)
		CreateTrigger(context.Context, *automation.Script, *automation.Trigger) error
		UpdateTrigger(context.Context, *automation.Script, *automation.Trigger) error
		DeleteTrigger(context.Context, *automation.Trigger) error
	}
)

func AutomationTrigger(tm automationTriggerManager) automationTrigger {
	var svc = automationTrigger{
		triggerManager: tm,
		logger:         DefaultLogger.Named("automation-trigger"),
	}

	return svc
}

func (svc automationTrigger) FindByID(ctx context.Context, triggerID uint64) (*automation.Trigger, error) {
	// @todo security check - can user read this trigger?
	return svc.triggerManager.FindTriggerByID(ctx, triggerID)
}

func (svc automationTrigger) Find(ctx context.Context, f automation.TriggerFilter) (automation.TriggerSet, automation.TriggerFilter, error) {
	// @todo security check - can user read these triggers?
	return svc.triggerManager.FindTriggers(ctx, f)
}

func (svc automationTrigger) Create(ctx context.Context, s *automation.Script, t *automation.Trigger) (err error) {
	// @todo security check - can user create trigger on this specific resource
	return svc.triggerManager.CreateTrigger(ctx, s, t)
}

func (svc automationTrigger) Update(ctx context.Context, s *automation.Script, t *automation.Trigger) (err error) {
	// @todo security check - can user create update triggers on this specific resource
	return svc.triggerManager.UpdateTrigger(ctx, s, t)
}

func (svc automationTrigger) Delete(ctx context.Context, t *automation.Trigger) (err error) {
	// @todo security check - can user create delete triggers on this specific resource
	return svc.triggerManager.DeleteTrigger(ctx, t)
}
