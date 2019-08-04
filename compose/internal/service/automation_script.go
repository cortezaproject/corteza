package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/automation"
)

type (
	automationScript struct {
		logger        *zap.Logger
		scriptManager automationScriptManager
	}

	automationScriptManager interface {
		FindScriptByID(context.Context, uint64) (*automation.Script, error)
		FindScripts(context.Context, automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error)
		CreateScript(context.Context, *automation.Script) error
		UpdateScript(context.Context, *automation.Script) error
		DeleteScript(context.Context, *automation.Script) error
	}
)

func AutomationScript(sm automationScriptManager) automationScript {
	var svc = automationScript{
		scriptManager: sm,
		logger:        DefaultLogger.Named("automation-script"),
	}

	return svc
}

func (svc automationScript) FindByID(ctx context.Context, scriptID uint64) (*automation.Script, error) {
	// @todo security check - can user read this script?
	return svc.scriptManager.FindScriptByID(ctx, scriptID)
}

func (svc automationScript) Find(ctx context.Context, f automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error) {
	// @todo security check - can user read these scripts?
	return svc.scriptManager.FindScripts(ctx, f)
}

func (svc automationScript) Create(ctx context.Context, s *automation.Script) (err error) {
	// @todo security check - can user create scripts?
	// @todo security check - can make scripts with security-definer?
	return svc.scriptManager.CreateScript(ctx, s)
}

func (svc automationScript) Update(ctx context.Context, s *automation.Script) (err error) {
	// @todo security check - can user update this script?
	// @todo security check - can make scripts with security-definer?
	return svc.scriptManager.UpdateScript(ctx, s)
}

func (svc automationScript) Delete(ctx context.Context, s *automation.Script) (err error) {
	// @todo security check - can user delete this script?
	return svc.scriptManager.DeleteScript(ctx, s)
}
