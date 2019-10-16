package service

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	automationScript struct {
		logger        *zap.Logger
		scriptManager automationScriptManager
		ac            automationScriptAccessController
		trg           automationTrigger
	}

	automationScriptManager interface {
		FindScriptByID(context.Context, uint64) (*automation.Script, error)
		FindScripts(context.Context, automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error)
		CreateScript(context.Context, *automation.Script) error
		UpdateScript(context.Context, *automation.Script) error
		DeleteScript(context.Context, *automation.Script) error
	}

	automationScriptAccessController interface {
		CanGrant(context.Context) bool

		CanCreateAutomationScript(context.Context) bool
		CanReadAutomationScript(context.Context, *automation.Script) bool
		CanUpdateAutomationScript(context.Context, *automation.Script) bool
		CanDeleteAutomationScript(context.Context, *automation.Script) bool

		FilterReadableScripts(ctx context.Context) *permissions.ResourceFilter
	}
)

func AutomationScript(sm automationScriptManager) automationScript {
	var svc = automationScript{
		scriptManager: sm,
		logger:        DefaultLogger.Named("automation-script"),
		ac:            DefaultAccessControl,
		trg:           DefaultAutomationTriggerManager,
	}

	return svc
}

func (svc automationScript) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc automationScript) FindByID(ctx context.Context, scriptID uint64) (*automation.Script, error) {
	if s, err := svc.loadCombo(ctx, scriptID); err != nil {
		return nil, err
	} else {
		return s, nil
	}
}

func (svc automationScript) Find(ctx context.Context, f automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error) {
	f.IsReadable = svc.ac.FilterReadableScripts(ctx)
	return svc.scriptManager.FindScripts(ctx, f)
}

func (svc automationScript) Create(ctx context.Context, mod *automation.Script) (err error) {
	if !svc.ac.CanCreateAutomationScript(ctx) {
		return ErrNoCreatePermissions.withStack()
	}

	if mod.RunAs > 0 {
		if !svc.ac.CanGrant(ctx) {
			return ErrNoGrantPermissions
		}
	}

	err = mod.Triggers().Walk(func(t *automation.Trigger) error {
		return svc.trg.isValid(ctx, mod, t)
	})

	if err != nil {
		return
	}

	return svc.scriptManager.CreateScript(ctx, mod)
}

func (svc automationScript) Update(ctx context.Context, mod *automation.Script) (err error) {
	var s *automation.Script

	if s, err = svc.loadCombo(ctx, mod.ID); err != nil {
		return err
	}

	if !svc.ac.CanUpdateAutomationScript(ctx, s) {
		return ErrNoUpdatePermissions.withStack()
	}

	// Users need to have grant privileges to
	// set script runner
	if mod.RunAs != s.RunAs {
		if !svc.ac.CanGrant(ctx) {
			return ErrNoGrantPermissions
		}
	}

	s.Name = mod.Name
	s.SourceRef = mod.SourceRef
	s.Source = mod.Source
	s.Async = mod.Async
	s.RunAs = mod.RunAs
	s.Timeout = mod.Timeout
	s.Critical = mod.Critical
	s.Enabled = mod.Enabled

	// Never run in a UA
	s.RunInUA = false

	err = mod.Triggers().Walk(func(t *automation.Trigger) error {
		return svc.trg.isValid(ctx, mod, t)
	})

	if err != nil {
		return
	}

	s.AddTrigger(automation.STMS_UPDATE, mod.Triggers()...)

	return svc.scriptManager.UpdateScript(ctx, s)
}

func (svc automationScript) Delete(ctx context.Context, scriptID uint64) (err error) {
	if s, err := svc.loadCombo(ctx, scriptID); err != nil {
		return err
	} else if !svc.ac.CanDeleteAutomationScript(ctx, s) {
		return ErrNoDeletePermissions.withStack()
	} else {
		return svc.scriptManager.DeleteScript(ctx, s)
	}
}

func (svc automationScript) loadCombo(ctx context.Context, scriptID uint64) (s *automation.Script, err error) {

	if scriptID > 0 {
		if s, err = svc.scriptManager.FindScriptByID(ctx, scriptID); err != nil {
			return
		} else if !svc.ac.CanReadAutomationScript(ctx, s) {
			err = ErrNoReadPermissions.withStack()
			return
		}
	}

	return
}
