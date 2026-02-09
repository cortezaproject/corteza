package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/modern-go/reflect2"
)

var (
	ErrExecutionStopped = errors.New("runtime: execution stopped")
	ErrStepFailed       = errors.New("runtime: step failed")
)

type executionGate interface {
	Request(execID id.ID, ops int) (<-chan struct{}, error)
}

type stateLedger interface {
	StepStarted(ctx context.Context, executableID, executionID, stepID id.ID, rev int) error
	StepCompleted(ctx context.Context, executableID, executionID, stepID id.ID, rev int, out any) error
	StepFailed(ctx context.Context, executableID, executionID, stepID id.ID, rev int, err error) error

	ExecutionCompleted(ctx context.Context, executableID, executionID id.ID, rev int) error
	ExecutionFailed(ctx context.Context, executableID, executionID id.ID, rev int, err error) error

	RecordFrame(ctx context.Context, executableID, executionID id.ID, rev int, frame types.StackFrame) error
}

type StepResult struct {
	StepID      id.ID
	StartedAt   time.Time
	CompletedAt time.Time
	Output      map[string]expr.TypedValue
	Error       error
}

type executionState struct {
	CurrentStep     *id.ID
	CompletedSteps  map[id.ID]StepResult
	InProgressSteps map[id.ID]bool
	mux             sync.RWMutex
}

type runtime struct {
	executionID id.ID
	exec        types.Executable
	scheduler   *scheduler

	globalState *expr.Vars
	entryPoint  string

	gate   executionGate
	ledger stateLedger

	state *executionState

	stopped atomic.Bool
	blocked atomic.Bool

	stopCh   chan struct{}
	resumeCh chan struct{}
}

func Runtime(
	executionID id.ID,
	exec types.Executable,
	gate executionGate,
	ledger stateLedger,
	entryPoint string,
) *runtime {
	state := &executionState{
		CompletedSteps:  make(map[id.ID]StepResult),
		InProgressSteps: make(map[id.ID]bool),
	}

	return &runtime{
		executionID: executionID,
		exec:        exec,
		gate:        gate,
		ledger:      ledger,
		state:       state,
		entryPoint:  entryPoint,
		scheduler:   newScheduler(exec),
		stopCh:      make(chan struct{}),
		resumeCh:    make(chan struct{}, 1),
	}
}

func (r *runtime) Start(ctx context.Context, global *expr.Vars) error {
	r.globalState = global

	for {
		select {
		case <-r.stopCh:
			return ErrExecutionStopped
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if r.blocked.Load() {
			select {
			case <-r.resumeCh:
			case <-r.stopCh:
				return ErrExecutionStopped
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		step, frameID, parentID, more, err := r.scheduler.Next(ctx)
		if err != nil {
			return err
		}
		if !more {
			return r.complete(ctx)
		}

		if err := r.waitForPermission(ctx, r.getStepOps(step)); err != nil {
			return err
		}

		if err := r.executeStep(ctx, step, frameID, parentID); err != nil {
			return r.fail(ctx, err)
		}

		if err := r.scheduler.OnStepComplete(ctx, step.ID); err != nil {
			return err
		}
	}
}

func (r *runtime) Stop() {
	if r.stopped.CompareAndSwap(false, true) {
		close(r.stopCh)
	}
}

func (r *runtime) Block() { r.blocked.Store(true) }
func (r *runtime) Resume() {
	if r.blocked.CompareAndSwap(true, false) {
		select {
		case r.resumeCh <- struct{}{}:
		default:
		}
	}
}

func (r *runtime) IsBlocked() bool { return r.blocked.Load() }
func (r *runtime) IsStopped() bool { return r.stopped.Load() }

func (r *runtime) executeStep(ctx context.Context, step *types.Step, frameID, parentID id.ID) error {
	r.state.mux.Lock()
	r.state.CurrentStep = &step.ID
	r.state.InProgressSteps[step.ID] = true
	r.state.mux.Unlock()

	if err := r.ledger.StepStarted(ctx, r.exec.ID, r.executionID, step.ID, r.exec.Revision); err != nil {
		return fmt.Errorf("step started: %w", err)
	}

	// Resolve inputs from scheduler stack
	inputVars, err := r.resolveInputs(step)
	if err != nil {
		return fmt.Errorf("resolve inputs: %w", err)
	}

	result := StepResult{
		StepID:    step.ID,
		StartedAt: time.Now(),
	}

	output, err := step.Handler.ExecN(ctx, &types.ExecRequest{
		Scope: inputVars,
	})

	// remove "" from inputVars as it's already outlined by "global"
	delete(inputVars, "")

	result.CompletedAt = time.Now()

	if err != nil {
		result.Error = err

		r.state.mux.Lock()
		delete(r.state.InProgressSteps, step.ID)
		r.state.CompletedSteps[step.ID] = result
		r.state.CurrentStep = nil
		r.state.mux.Unlock()

		_ = r.ledger.StepFailed(ctx, r.exec.ID, r.executionID, step.ID, r.exec.Revision, err)

		// Record frame in ledger even on failure
		_ = r.ledger.RecordFrame(ctx, r.exec.ID, r.executionID, r.exec.Revision, types.StackFrame{
			ID:        frameID,
			StepID:    step.ID,
			ParentID:  parentID,
			Handle:    step.Handle,
			Kind:      step.Kind,
			Input:     inputVars,
			Output:    nil,
			StartedAt: result.StartedAt,
			EndedAt:   &result.CompletedAt,
			Error:     err,
		})

		return fmt.Errorf("%w: %s: %v", ErrStepFailed, step.ID, err)
	}

	// Extract results based on step definition
	outputMap := make(map[string]expr.TypedValue)
	if output != nil {
		if vars, ok := output.(*expr.Vars); ok {
			// Map each defined result to a Vars containing that value
			for _, rst := range step.Results {
				v := vars.GetValue()[rst.ArgumentName]

				if !reflect2.IsNil(v) {
					outputMap[rst.ArgumentName] = v
				}
			}
		}
	}

	result.Output = outputMap

	r.state.mux.Lock()
	delete(r.state.InProgressSteps, step.ID)
	r.state.CompletedSteps[step.ID] = result
	r.state.CurrentStep = nil
	r.state.mux.Unlock()

	// Store outputs in scheduler for future step resolution
	if err := r.scheduler.StoreOutputs(step.ID, outputMap); err != nil {
		return fmt.Errorf("store outputs: %w", err)
	}

	if err := r.ledger.StepCompleted(ctx, r.exec.ID, r.executionID, step.ID, r.exec.Revision, outputMap); err != nil {
		return fmt.Errorf("step completed: %w", err)
	}

	// Record frame in ledger
	_ = r.ledger.RecordFrame(ctx, r.exec.ID, r.executionID, r.exec.Revision, types.StackFrame{
		ID:        frameID,
		StepID:    step.ID,
		ParentID:  parentID,
		Handle:    step.Handle,
		Kind:      step.Kind,
		Input:     inputVars,
		Output:    expr.Must(expr.NewVars(outputMap)),
		StartedAt: result.StartedAt,
		EndedAt:   &result.CompletedAt,
	})

	return nil
}

func (r *runtime) resolveInputs(step *types.Step) (map[string]*expr.Vars, error) {
	out := make(map[string]*expr.Vars, 4)

	out[""] = r.globalState
	out["global"] = r.globalState
	out[r.entryPoint] = r.globalState

	if len(step.Arguments) == 0 {
		return out, nil
	}

	for _, arg := range step.Arguments {
		// If no scope, we're using the global one
		if arg.Scope == "" || arg.Scope == r.entryPoint {
			continue
		}

		// Get entire output map from the scope handle
		outputs, err := r.scheduler.FindOutput(arg.Scope)
		if err != nil {
			return nil, fmt.Errorf("resolve %s from scope %s: %w", arg.Target, arg.Scope, err)
		}

		out[arg.Scope] = outputs
	}

	return out, nil
}

func (r *runtime) waitForPermission(ctx context.Context, ops int) error {
	permCh, err := r.gate.Request(r.executionID, ops)
	if err != nil {
		return err
	}

	select {
	case <-permCh:
		return nil
	case <-r.stopCh:
		return ErrExecutionStopped
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *runtime) complete(ctx context.Context) error {
	if err := r.ledger.ExecutionCompleted(ctx, r.exec.ID, r.executionID, r.exec.Revision); err != nil {
		return fmt.Errorf("execution complete: %w", err)
	}
	return nil
}

func (r *runtime) fail(ctx context.Context, err error) error {
	if repErr := r.ledger.ExecutionFailed(ctx, r.exec.ID, r.executionID, r.exec.Revision, err); repErr != nil {
		return fmt.Errorf("execution failed: %w (report error: %v)", err, repErr)
	}
	return err
}

func (r *runtime) getStepOps(step *types.Step) int {
	return 1
}
