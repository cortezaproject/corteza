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
	"go.uber.org/zap"
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
	ExecutionPaused(ctx context.Context, executableID, executionID, stepID id.ID, phaseIndex, rev int, err error) error

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
	log         *zap.Logger
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
	log *zap.Logger,
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
		log:         log,
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
	r.log.Debug("execution started",
		zap.Stringer("executionID", r.executionID),
		zap.String("executableID", r.exec.ID.String()),
		zap.String("entryPoint", r.entryPoint),
	)
	defer r.log.Debug("execution loop exited", zap.Stringer("executionID", r.executionID))

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

		step, frameID, parentID, more, err := r.scheduler.Next(ctx, r.globalState, r.entryPoint)
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
			r.log.Error("step execution failed",
				zap.String("stepID", step.ID.String()),
				zap.String("handle", step.Handle),
				zap.String("kind", step.Kind),
				zap.Error(err),
			)
			// executeStep only propagates errors not already dispatched internally.
			return r.fail(ctx, err)
		}

		// After executeStep returns nil:
		// - A recoverable pause was triggered: r.blocked is true; loop to the wait select.
		// - Mid-phase success: top frame is still this step with more phases remaining.
		// - Catch was pushed: top frame is a different step.
		// Only call OnStepComplete when the step is truly finished.
		if r.blocked.Load() {
			continue
		}

		// If the top frame is still this step and it has more phases, loop again.
		if tf := r.scheduler.topFrame(); tf != nil && tf.stepID.Equal(step.ID) {
			if phased, ok := step.Handler.(types.PhasedStepHandler); ok {
				if tf.phaseIndex < len(phased.Phases()) {
					continue
				}
			}
		}

		// Check if this was a termination step
		if step.Kind == "termination" {
			allTerminated, err := r.scheduler.TerminateBranch(step.ID)
			if err != nil {
				return r.fail(ctx, err)
			}

			if allTerminated {
				return r.complete(ctx)
			}

			continue
		}

		// If catch was pushed, the top frame is a different step — skip OnStepComplete.
		tf := r.scheduler.topFrame()
		if tf == nil || !tf.stepID.Equal(step.ID) {
			continue
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

	// Obtain the current frame so we can read/write phaseIndex and retryCount.
	cf := r.scheduler.topFrame()

	// Only log StepStarted on the very first phase (phaseIndex == 0).
	if cf == nil || cf.phaseIndex == 0 {
		r.log.Debug("step started",
			zap.String("stepID", step.ID.String()),
			zap.String("handle", step.Handle),
			zap.String("kind", step.Kind),
		)
		if err := r.ledger.StepStarted(ctx, r.exec.ID, r.executionID, step.ID, r.exec.Revision); err != nil {
			return fmt.Errorf("step started: %w", err)
		}
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

	// Execute: phased or single-shot.
	var output types.ExecResponse
	if phased, ok := step.Handler.(types.PhasedStepHandler); ok {
		phases := phased.Phases()
		phIdx := 0
		if cf != nil {
			phIdx = cf.phaseIndex
		}
		if phIdx >= len(phases) {
			// All phases already completed — treat as success with no output.
			output = nil
			err = nil
		} else {
			output, err = phases[phIdx](ctx, &types.ExecRequest{Scope: inputVars})
			if err == nil && cf != nil {
				cf.phaseIndex++
				cf.retryCount = 0
				// If more phases remain, return without completing the step.
				if cf.phaseIndex < len(phases) {
					return nil
				}
			}
		}
	} else {
		output, err = step.Handler.ExecN(ctx, &types.ExecRequest{Scope: inputVars})
	}

	// remove "" from inputVars as it's already outlined by "global"
	delete(inputVars, "")

	result.CompletedAt = time.Now()

	if err != nil {
		result.Error = err

		r.log.Warn("step failed",
			zap.String("stepID", step.ID.String()),
			zap.String("handle", step.Handle),
			zap.String("kind", step.Kind),
			zap.Error(err),
		)

		r.state.mux.Lock()
		delete(r.state.InProgressSteps, step.ID)
		r.state.CompletedSteps[step.ID] = result
		r.state.CurrentStep = nil
		r.state.mux.Unlock()

		_ = r.ledger.StepFailed(ctx, r.exec.ID, r.executionID, step.ID, r.exec.Revision, err)

		// Record frame in ledger even on failure
		phIdx := 0
		if cf != nil {
			phIdx = cf.phaseIndex
		}
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

		return r.dispatchError(ctx, step, cf, phIdx, err)
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

			// Preserve internal iterator bindings.
			if iterHandler, ok := vars.GetValue()["_iter_handler"]; ok {
				outputMap["_iter_handler"] = iterHandler
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

	r.log.Debug("step completed",
		zap.String("stepID", step.ID.String()),
		zap.String("handle", step.Handle),
		zap.String("kind", step.Kind),
		zap.Duration("duration", result.CompletedAt.Sub(result.StartedAt)),
	)

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

// dispatchError implements the error dispatch logic from the spec:
//
//	if step.Recoverable && isRecoverable(err) && (MaxRetries==0 || retryCount<MaxRetries): pause
//	else: find catch handler; if none → fail; else push catch frame
func (r *runtime) dispatchError(ctx context.Context, step *types.Step, cf *frame, phIdx int, err error) error {
	if step.Recoverable && types.IsRecoverable(err) {
		retryCount := 0
		if cf != nil {
			retryCount = cf.retryCount
		}
		if step.MaxRetries == 0 || retryCount < step.MaxRetries {
			if cf != nil {
				cf.retryCount++
			}
			_ = r.ledger.ExecutionPaused(ctx, r.exec.ID, r.executionID, step.ID, phIdx, r.exec.Revision, err)
			r.Block()
			return nil
		}
	}

	// Walk stack for a catch handler.
	stackIdx := r.scheduler.topFrameIdx()
	catchID, _ := r.scheduler.FindCatch(stackIdx)
	if catchID.IsZero() {
		return r.fail(ctx, fmt.Errorf("%w: %s: %v", ErrStepFailed, step.ID, err))
	}

	// Build error vars for the catch scope.
	errVars, _ := buildErrVars(step.ID, err)
	if pushErr := r.scheduler.PushCatch(catchID, errVars); pushErr != nil {
		return r.fail(ctx, fmt.Errorf("%w: push catch: %v (original: %v)", ErrStepFailed, pushErr, err))
	}

	return nil
}

// buildErrVars constructs the error context map injected into the catch scope.
func buildErrVars(stepID id.ID, err error) (map[string]expr.TypedValue, error) {
	msg, _ := expr.NewString(err.Error())
	sid, _ := expr.NewString(stepID.String())
	return map[string]expr.TypedValue{
		"errorMessage": msg,
		"errorStepID":  sid,
	}, nil
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
			keys := ""
			for k := range r.scheduler.completedOutputs {
				keys += k + ","
			}
			for i := len(r.scheduler.stack) - 1; i >= 0; i-- {
				f := r.scheduler.stack[i]
				keys += fmt.Sprintf("stack:%v[%v],", f.handle, f.stepID)
			}
			r.log.Debug("scheduler output not found for scope",
				zap.String("scope", arg.Scope),
				zap.String("available_keys", keys),
				zap.Any("arg", arg),
			)
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
