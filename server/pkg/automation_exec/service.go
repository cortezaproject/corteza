package automation_exec

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/governor"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/ledger"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/registry"
	manager "github.com/cortezaproject/corteza/server/pkg/automation_exec/runtime_manager"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/supervisor"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"go.uber.org/zap"
)

type (
	automationService struct {
		log *zap.Logger
		reg registrySvc
		// @todo when needed
		sup interface{}

		rm  runtimeManagerAPI
		led ledgerAPI
	}

	registrySvc interface {
		Add(ctx context.Context, exec types.Executable) error
		Deprecate(ctx context.Context, executableID id.ID, revision int) error
		Remove(ctx context.Context, executableID id.ID, revision int) error
	}

	ledgerAPI interface {
		GetExecution(ctx context.Context, executableID, executionID id.ID, revision int) (*types.Execution, error)
		GetTrace(ctx context.Context, executableID, executionID id.ID, revision int) ([]types.StackFrame, error)
	}

	runtimeManagerAPI interface {
		Start(ctx context.Context, executableID id.ID, revision int, params *expr.Vars) (id.ID, error)
		Get(execID id.ID) (*manager.RuntimeEntry, error)
	}
)

func AutomationService(ctx context.Context, log *zap.Logger, cfg manager.Config) (*automationService, error) {
	log.Debug("initializing ledger")
	led := ledger.Ledger(log.Named("ledger"))
	log.Debug("initialized ledger")

	log.Debug("initializing governor")
	gov := governor.Governor(ctx, log.Named("governor"))
	log.Debug("initialized governor")

	log.Debug("initializing registry")
	reg := registry.Registry(log.Named("registry"), led)
	log.Debug("initialized registry")

	log.Debug("initializing runtime-manager")
	rm, err := manager.RuntimeManager(ctx, log.Named("runtime-manager"), reg, led, gov, cfg)
	log.Debug("initialized runtime-manager")

	if err != nil {
		return nil, err
	}

	log.Debug("initializing supervisor")
	sup := supervisor.Supervisor(log.Named("supervisor"), led, rm)
	log.Debug("initialized supervisor")

	return &automationService{
		log: log,
		reg: reg,
		rm:  rm,
		led: led,
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

func (s *automationService) ExecuteAndWait(
	ctx context.Context,
	exeID id.ID,
	rev int,
	params *expr.Vars,
) (*types.ExecutionResult, error) {
	executionID, err := s.rm.Start(ctx, exeID, rev, params)
	if err != nil {
		return nil, err
	}

	entry, err := s.rm.Get(executionID)
	if err != nil {
		return nil, err
	}

	select {
	case <-entry.Done:
		// Fetch execution details from ledger
		exec, err := s.led.GetExecution(ctx, exeID, executionID, rev)
		if err != nil {
			return nil, err
		}

		return s.prepMetaResponse(exec), nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *automationService) GetExecutionTrace(ctx context.Context, exeID id.ID, rev int, executionID id.ID) ([]types.StackFrame, error) {
	return s.led.GetTrace(ctx, exeID, executionID, rev)
}

func (s *automationService) prepMetaResponse(exec *types.Execution) *types.ExecutionResult {
	out := &types.ExecutionResult{
		ExecutionID:  exec.ID,
		ExecutableID: exec.ExecutableID,
		Revision:     exec.Revision,
		Status:       exec.Status,
		StartedAt:    exec.CreatedAt,
		EndedAt:      exec.EndedAt,
	}

	if exec.EndedAt != nil {
		out.Duration = formatDuration(exec.EndedAt.Sub(exec.CreatedAt))
	}

	if exec.Error != nil {
		out.Error = exec.Error.Error()
	}

	return out
}

func formatDuration(d time.Duration) string {
	ms := float64(d.Nanoseconds()) / 1e6

	if ms < 100000 {
		return fmt.Sprintf("%.3fms", ms)
	}

	seconds := d.Seconds()
	if seconds < 1000 {
		return fmt.Sprintf("%.3fs", seconds)
	}

	minutes := d.Minutes()
	if minutes < 1000 {
		return fmt.Sprintf("%.3fmin", minutes)
	}

	return fmt.Sprintf("%.3fh", d.Hours())
}
