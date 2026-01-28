package automation_exec

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/governor"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/ledger"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/registry"
	manager "github.com/cortezaproject/corteza/server/pkg/automation_exec/runtime_manager"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/supervisor"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

type (
	automationService struct {
		reg registrySvc
		// @todo when needed
		sup interface{}

		rm runtimeManagerAPI
	}

	registrySvc interface {
		Add(ctx context.Context, exec types.Executable) error
		Deprecate(ctx context.Context, executableID id.ID, revision int) error
		Remove(ctx context.Context, executableID id.ID, revision int) error
	}

	runtimeManagerAPI interface {
		Start(ctx context.Context, executionID id.ID, revision int, params *expr.Vars) (id.ID, error)
	}
)

func AutomationService(ctx context.Context, cfg manager.Config) (*automationService, error) {
	led := ledger.Ledger()
	gov := governor.Governor(ctx)
	reg := registry.Registry(led)

	rm, err := manager.RuntimeManager(ctx, reg, led, gov, cfg)
	if err != nil {
		return nil, err
	}

	sup := supervisor.Supervisor(led, rm)

	return &automationService{
		reg: reg,
		rm:  rm,
		sup: sup,
	}, nil
}

// ---- facade API ----

func (s *automationService) RegisterExecutable(ctx context.Context, exe types.Executable) error {
	return s.reg.Add(ctx, exe)
}

func (s *automationService) DeprecateExecutable(ctx context.Context, exeID id.ID, rev int) error {
	return s.reg.Deprecate(ctx, exeID, rev)
}

func (s *automationService) RemoveExecutable(ctx context.Context, exeID id.ID, rev int) error {
	return s.reg.Remove(ctx, exeID, rev)
}

func (s *automationService) Execute(
	ctx context.Context,
	exeID id.ID,
	rev int,
	params *expr.Vars,
) (id.ID, error) {
	return s.rm.Start(ctx, exeID, rev, params)
}
