package manager

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/runtime"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
)

//
// ===== errors =====
//

var (
	ErrSystemDraining     = errors.New("manager: system is draining")
	ErrSlotTimeout        = errors.New("manager: concurrency slot timeout")
	ErrExecutableNotFound = errors.New("manager: executable not found")
	ErrExecutableInactive = errors.New("manager: executable not active")
	ErrExecutionNotFound  = errors.New("manager: execution not found")
	ErrQueueFull          = errors.New("manager: start queue full")
	ErrInvalidConfig      = errors.New("manager: invalid configuration")
	ErrNotInDebugMode     = errors.New("manager: not in debug mode")
	ErrStepInProgress     = errors.New("manager: step already in progress")
)

//
// ===== config =====
//

type Config struct {
	MaxConcurrent int
	MaxQueued     int
	SlotTimeout   time.Duration
}

//
// ===== external dependencies =====
//

type Registry interface {
	GetExecutable(ctx context.Context, id id.ID, revision int) (types.Executable, error)
}

// type Ledger interface {
// 	RegisterExecution(ctx context.Context, exec types.Execution) error
// 	UpdateExecution(ctx context.Context, execID id.ID, status types.Status, endedAt *time.Time) error
// }

type Ledger interface {
	StepStarted(ctx context.Context, executableID, executionID, stepID id.ID, rev int) error
	StepCompleted(ctx context.Context, executableID, executionID, stepID id.ID, rev int, out any) error
	StepFailed(ctx context.Context, executableID, executionID, stepID id.ID, rev int, err error) error

	ExecutionCompleted(ctx context.Context, executableID, executionID id.ID, rev int) error
	ExecutionFailed(ctx context.Context, executableID, executionID id.ID, rev int, err error) error

	IsExecutableInUse(ctx context.Context, executableID id.ID, rev int) (bool, error)

	RegisterExecution(ctx context.Context, executionID, executableID id.ID, rev int, params *expr.Vars) error
}

type Governor interface {
	AddExecution(execID id.ID, maxOpPerRequest int, budget types.Budget, rate types.RateLimit) error
	RemoveExecution(execID id.ID)
	Request(execID id.ID, ops int) (<-chan struct{}, error)
}

type Runtime interface {
	Start(ctx context.Context, scope *expr.Vars) error
	Stop()
	Block()
	Resume()
	IsBlocked() bool
}

//
// ===== internal types =====
//

type queuedStart struct {
	execID     id.ID
	executable types.Executable
	params     *expr.Vars
}

type runtimeEntry struct {
	execID       id.ID
	executableID id.ID
	revision     int

	runtime Runtime
	cancel  context.CancelFunc

	debug    atomic.Bool
	stepGate chan struct{}
}

//
// ===== runtime manager =====
//

type runtimeManager struct {
	registry Registry
	ledger   Ledger
	governor Governor
	config   Config

	runtimes map[id.ID]*runtimeEntry
	mu       sync.RWMutex

	queue    []queuedStart
	queueMux sync.Mutex
	qSignal  chan struct{}

	slots chan struct{}

	running  atomic.Int32
	draining atomic.Bool
}

//
// ===== constructor =====
//

// The Runtime Manager is the orchestrator.
//
// * Starts executions
// * Queues them
// * Enforces concurrency limits
// * Wires runtime + governor + ledger
// * Cleans up when executions finish
func RuntimeManager(
	ctx context.Context,
	reg Registry,
	led Ledger,
	gov Governor,
	cfg Config,
) (*runtimeManager, error) {
	if cfg.MaxConcurrent <= 0 {
		return nil, ErrInvalidConfig
	}
	if cfg.SlotTimeout <= 0 {
		cfg.SlotTimeout = 30 * time.Second
	}

	rm := &runtimeManager{
		registry: reg,
		ledger:   led,
		governor: gov,
		config:   cfg,
		runtimes: make(map[id.ID]*runtimeEntry),
		queue:    make([]queuedStart, 0),
		qSignal:  make(chan struct{}, 1),
		slots:    make(chan struct{}, cfg.MaxConcurrent),
	}

	go rm.watchQueue(ctx)
	return rm, nil
}

//
// ===== admission =====
//

func (rm *runtimeManager) Start(
	ctx context.Context,
	executableID id.ID,
	revision int,
	params *expr.Vars,
) (id.ID, error) {
	if rm.draining.Load() {
		return id.Zero(), ErrSystemDraining
	}

	executable, err := rm.registry.GetExecutable(ctx, executableID, revision)
	if err != nil {
		return id.Zero(), ErrExecutableNotFound
	}

	active, err := rm.ledger.IsExecutableInUse(ctx, executableID, revision)
	if err != nil || !active {
		return id.Zero(), ErrExecutableInactive
	}

	eid := id.MustNumID(id.Next())

	err = rm.ledger.RegisterExecution(ctx, eid, executableID, revision, params)
	if err != nil {
		return id.Zero(), err
	}

	if err := rm.enqueue(queuedStart{
		execID:     eid,
		executable: executable,
		params:     params,
	}); err != nil {
		return id.Zero(), err
	}

	rm.signalQueue()
	return eid, nil
}

func (rm *runtimeManager) Stop(execID id.ID) error {
	// 1) If running, stop runtime + cancel its context.
	rm.mu.RLock()
	e, ok := rm.runtimes[execID]
	rm.mu.RUnlock()

	if ok {
		// Cancel context first to unblock waits, then hard-stop the runtime.
		if e.cancel != nil {
			e.cancel()
		}
		e.runtime.Stop()
		return nil
	}

	// 2) If queued (not started yet), drop from queue.
	rm.queueMux.Lock()
	defer rm.queueMux.Unlock()

	for i := range rm.queue {
		if rm.queue[i].execID.Equal(execID) {
			rm.queue = append(rm.queue[:i], rm.queue[i+1:]...)
			return nil
		}
	}

	return ErrExecutionNotFound
}

//
// ===== queue =====
//

func (rm *runtimeManager) enqueue(q queuedStart) error {
	rm.queueMux.Lock()
	defer rm.queueMux.Unlock()

	if rm.config.MaxQueued > 0 && len(rm.queue) >= rm.config.MaxQueued {
		return ErrQueueFull
	}

	rm.queue = append(rm.queue, q)
	return nil
}

func (rm *runtimeManager) dequeue() (queuedStart, bool) {
	rm.queueMux.Lock()
	defer rm.queueMux.Unlock()

	if len(rm.queue) == 0 {
		return queuedStart{}, false
	}

	q := rm.queue[0]
	rm.queue = rm.queue[1:]
	return q, true
}

func (rm *runtimeManager) signalQueue() {
	select {
	case rm.qSignal <- struct{}{}:
	default:
	}
}

func (rm *runtimeManager) watchQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-rm.qSignal:
			rm.processQueue(ctx)
		}
	}
}

func (rm *runtimeManager) processQueue(ctx context.Context) {
	for {
		item, ok := rm.dequeue()
		if !ok {
			return
		}

		err := rm.startQueued(ctx, item)
		if err != nil {
			// @todo consider re-enqueuing on fail
			panic(err)
		}
	}
}

//
// ===== runtime start =====
//

func (rm *runtimeManager) startQueued(ctx context.Context, q queuedStart) error {
	select {
	case rm.slots <- struct{}{}:
	case <-time.After(rm.config.SlotTimeout):
		return ErrSlotTimeout
	case <-ctx.Done():
		return ctx.Err()
	}

	// @todo limits
	if err := rm.governor.AddExecution(q.execID, 0, types.Budget{}, types.RateLimit{}); err != nil {
		rm.releaseSlot()
		return err
	}

	rctx, cancel := context.WithCancel(ctx)
	rt := rm.createRuntime(q.execID, q.executable, q.params)

	entry := &runtimeEntry{
		execID:       q.execID,
		executableID: q.executable.ID,
		revision:     q.executable.Revision,
		runtime:      rt,
		cancel:       cancel,
		stepGate:     make(chan struct{}, 1),
	}

	rm.mu.Lock()
	rm.runtimes[q.execID] = entry
	rm.mu.Unlock()

	rm.running.Add(1)
	go rm.runRuntime(rctx, entry, q.params)
	return nil
}

func (rm *runtimeManager) runRuntime(ctx context.Context, e *runtimeEntry, scope *expr.Vars) {
	err := e.runtime.Start(ctx, scope)
	rm.onExit(ctx, e, err)
}

func (rm *runtimeManager) onExit(ctx context.Context, e *runtimeEntry, err error) {
	status := types.StatusCompleted

	if err != nil {
		status = types.StatusFailed
	}

	switch status {
	case types.StatusCompleted:
		// @todo error handling?
		_ = rm.ledger.ExecutionCompleted(ctx, e.executableID, e.execID, e.revision)
	case types.StatusFailed:
		_ = rm.ledger.ExecutionFailed(ctx, e.executableID, e.execID, e.revision, err)

	default:
		// @todo
	}

	rm.governor.RemoveExecution(e.execID)
	rm.releaseSlot()

	rm.mu.Lock()
	delete(rm.runtimes, e.execID)
	rm.mu.Unlock()

	rm.running.Add(-1)
	rm.signalQueue()
}

func (rm *runtimeManager) releaseSlot() {
	select {
	case <-rm.slots:
	default:
	}
}

func (rm *runtimeManager) SetDraining(draining bool) {
	rm.draining.Store(draining)
}

func (rm *runtimeManager) IsDraining() bool {
	return rm.draining.Load()
}

//
// ===== helpers =====
//

func (rm *runtimeManager) get(execID id.ID) (*runtimeEntry, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	e, ok := rm.runtimes[execID]
	if !ok {
		return nil, ErrExecutionNotFound
	}
	return e, nil
}

func (rm *runtimeManager) createRuntime(
	execID id.ID,
	exe types.Executable,
	params *expr.Vars,
) Runtime {
	return runtime.Runtime(
		exe.ID,
		exe,
		rm.governor,
		rm.ledger,
	)

	panic("runtime factory not implemented")
}
