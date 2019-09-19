package service

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/automation"
)

type (
	// Handles automation triggers storing and loading
	automationTrigger struct {
		logger         *zap.Logger
		triggerManager automationTriggerManager

		mod ModuleService
		ac  automationTriggerAccessController
	}

	automationTriggerManager interface {
		FindTriggerByID(context.Context, uint64) (*automation.Trigger, error)
		FindTriggers(context.Context, automation.TriggerFilter) (automation.TriggerSet, automation.TriggerFilter, error)
		CreateTrigger(context.Context, *automation.Script, *automation.Trigger) error
		UpdateTrigger(context.Context, *automation.Script, *automation.Trigger) error
		DeleteTrigger(context.Context, *automation.Trigger) error
	}

	automationTriggerAccessController interface {
		CanUpdateAutomationScript(context.Context, *automation.Script) bool
		CanManageAutomationTriggersOnModule(context.Context, *types.Module) bool
	}
)

func AutomationTrigger(tm automationTriggerManager) automationTrigger {
	var svc = automationTrigger{
		triggerManager: tm,
		logger:         DefaultLogger.Named("automation-trigger"),

		ac:  DefaultAccessControl,
		mod: DefaultModule,
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
	if err = svc.isValid(ctx, s, t); err != nil {
		return
	}

	if !svc.ac.CanUpdateAutomationScript(ctx, s) {
		return ErrNoTriggerManagementPermissions
	}

	return svc.triggerManager.CreateTrigger(ctx, s, t)
}

func (svc automationTrigger) Update(ctx context.Context, s *automation.Script, t *automation.Trigger) (err error) {
	if err = svc.isValid(ctx, s, t); err != nil {
		return
	}

	if !svc.ac.CanUpdateAutomationScript(ctx, s) {
		return ErrNoTriggerManagementPermissions
	}

	return svc.triggerManager.UpdateTrigger(ctx, s, t)
}

func (svc automationTrigger) Delete(ctx context.Context, s *automation.Script, t *automation.Trigger) (err error) {
	if err = svc.isValid(ctx, s, t); err != nil {
		return
	}

	if !svc.ac.CanUpdateAutomationScript(ctx, s) {
		return ErrNoTriggerManagementPermissions
	}

	return svc.triggerManager.DeleteTrigger(ctx, t)
}

// Validates trigger (in compose context)
func (svc automationTrigger) isValid(ctx context.Context, s *automation.Script, t *automation.Trigger) error {
	if s == nil {
		return errors.WithStack(automation.ErrAutomationScriptInvalid)
	}

	if !t.Enabled {
		return nil
	}

	if t.Resource != "compose:record" {
		// Accepting only compose:record resources
		return errors.WithStack(automation.ErrAutomationTriggerInvalidResource)
	}

	if t.IsDeferred() {
		// @todo validate condition for deferred triggers
		return nil
	}

	switch t.Event {
	case "manual",
		"beforeCreate", "beforeUpdate", "beforeDelete",
		"afterCreate", "afterUpdate", "afterDelete":
		var moduleID = t.Uint64Condition()

		if t.Event != "manual" && moduleID == 0 {
			return errors.WithStack(automation.ErrAutomationTriggerInvalidCondition)
		}

		if moduleID > 0 {
			if m, err := svc.mod.With(ctx).FindByID(s.NamespaceID, moduleID); err != nil {
				return err
			} else if !svc.ac.CanManageAutomationTriggersOnModule(ctx, m) {
				return errors.WithStack(ErrNoTriggerManagementPermissions)
			}
		}

	case "interval",
		"deferred":
		if s.RunAs == 0 {
			return errors.WithStack(automation.ErrAutomationScriptMissingUser)
		}

	default:
		return errors.WithStack(automation.ErrAutomationTriggerInvalidEvent)
	}

	return nil
}
