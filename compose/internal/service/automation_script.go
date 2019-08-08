package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/pkg/automation"
)

type (
	automationScript struct {
		logger        *zap.Logger
		scriptManager automationScriptManager
		ns            NamespaceService
		ac            automationScriptAccessController
	}

	automationScriptManager interface {
		FindScriptByID(context.Context, uint64) (*automation.Script, error)
		FindScripts(context.Context, automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error)
		CreateScript(context.Context, *automation.Script) error
		UpdateScript(context.Context, *automation.Script) error
		DeleteScript(context.Context, *automation.Script) error
	}

	automationScriptAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool

		CanCreateAutomationScript(context.Context, *types.Namespace) bool
		CanReadAnyAutomationScript(context.Context) bool
		CanReadAutomationScript(context.Context, *automation.Script) bool
		CanUpdateAutomationScript(context.Context, *automation.Script) bool
		CanDeleteAutomationScript(context.Context, *automation.Script) bool
	}

	automationScriptNamespaceFinder interface {
		FindByID(uint64) (*types.Namespace, error)
	}
)

func AutomationScript(sm automationScriptManager) automationScript {
	var svc = automationScript{
		scriptManager: sm,
		logger:        DefaultLogger.Named("automation-script"),
		ac:            DefaultAccessControl,
		ns:            Namespace(),
	}

	return svc
}

func (svc automationScript) FindByID(ctx context.Context, namespaceID, scriptID uint64) (*automation.Script, error) {
	if _, err := svc.loadNamespace(ctx, namespaceID); err != nil {
		return nil, err
	}

	if script, err := svc.scriptManager.FindScriptByID(ctx, scriptID); err != nil {
		return nil, err
	} else if !svc.ac.CanReadAutomationScript(ctx, script) {
		return nil, ErrNoCreatePermissions.withStack()
	} else {
		return script, nil
	}
}

func (svc automationScript) Find(ctx context.Context, namespaceID uint64, f automation.ScriptFilter) (automation.ScriptSet, automation.ScriptFilter, error) {
	if _, err := svc.loadNamespace(ctx, namespaceID); err != nil {
		return nil, f, err
	}

	f.AccessCheck = permissions.InitAccessCheckFilter(
		"read",
		auth.GetIdentityFromContext(ctx).Roles(),
		svc.ac.CanReadAnyAutomationScript(ctx),
	)

	return svc.scriptManager.FindScripts(ctx, f)
}

func (svc automationScript) Create(ctx context.Context, namespaceID uint64, s *automation.Script) (err error) {
	if ns, err := svc.loadNamespace(ctx, namespaceID); err != nil {
		return err
	} else if !svc.ac.CanCreateAutomationScript(ctx, ns) {
		return ErrNoCreatePermissions.withStack()
	}

	return svc.scriptManager.CreateScript(ctx, s)
}

func (svc automationScript) Update(ctx context.Context, namespaceID uint64, s *automation.Script) (err error) {
	if _, err := svc.loadNamespace(ctx, namespaceID); err != nil {
		return err
	} else if !svc.ac.CanUpdateAutomationScript(ctx, s) {
		return ErrNoCreatePermissions.withStack()
	}

	return svc.scriptManager.UpdateScript(ctx, s)
}

func (svc automationScript) Delete(ctx context.Context, namespaceID uint64, s *automation.Script) (err error) {
	if _, err := svc.loadNamespace(ctx, namespaceID); err != nil {
		return err
	} else if !svc.ac.CanDeleteAutomationScript(ctx, s) {
		return ErrNoCreatePermissions.withStack()
	}

	return svc.scriptManager.DeleteScript(ctx, s)
}

func (svc automationScript) loadNamespace(ctx context.Context, namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}

	if ns, err = svc.ns.With(ctx).FindByID(namespaceID); err != nil {
		return
	}

	if !svc.ac.CanReadNamespace(ctx, ns) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}
