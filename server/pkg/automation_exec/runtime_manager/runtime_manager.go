package manager

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
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
	GetExecutable(ctx context.Context, id types.ExecutableID, revision int) (types.Executable, error)
	IsActive(ctx context.Context, id types.ExecutableID, revision int) (bool, error)
}

type Ledger interface {
	RegisterExecution(ctx context.Context, exec types.Execution) error
	UpdateExecution(ctx context.Context, execID id.ID, status types.Status, endedAt *time.Time) error
}

type Governor interface {
	InitExecution(execID id.ID, limits types.ExecutableLimits) error
	CleanupExecution(execID id.ID)
}

type Runtime interface {
	Execute(ctx context.Context) error
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
	params     map[string]any
}

type runtimeEntry struct {
	execID       id.ID
	executableID types.ExecutableID
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

	queue   []queuedStart
	queueMu sync.Mutex
	qSignal chan struct{}

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
	execID types.ExecutableID,
	revision int,
	params map[string]any,
) (id.ID, error) {
	if rm.draining.Load() {
		return id.Zero(), ErrSystemDraining
	}

	exe, err := rm.registry.GetExecutable(ctx, execID, revision)
	if err != nil {
		return id.Zero(), ErrExecutableNotFound
	}

	active, err := rm.registry.IsActive(ctx, execID, revision)
	if err != nil || !active {
		return id.Zero(), ErrExecutableInactive
	}

	eid := id.MustNumID(id.Next())

	if err := rm.ledger.RegisterExecution(ctx, types.Execution{
		ID:           eid,
		ExecutableID: execID,
		Revision:     uint32(revision),
		Status:       types.StatusCreated,
		CreatedAt:    time.Now(),
		Variables:    params,
	}); err != nil {
		return id.Zero(), err
	}

	if err := rm.enqueue(queuedStart{
		execID:     eid,
		executable: exe,
		params:     params,
	}); err != nil {
		return id.Zero(), err
	}

	rm.signalQueue()
	return eid, nil
}

//
// ===== queue =====
//

func (rm *runtimeManager) enqueue(q queuedStart) error {
	rm.queueMu.Lock()
	defer rm.queueMu.Unlock()

	if rm.config.MaxQueued > 0 && len(rm.queue) >= rm.config.MaxQueued {
		return ErrQueueFull
	}

	rm.queue = append(rm.queue, q)
	return nil
}

func (rm *runtimeManager) dequeue() (queuedStart, bool) {
	rm.queueMu.Lock()
	defer rm.queueMu.Unlock()

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

		if err := rm.startQueued(ctx, item); err != nil {
			_ = rm.enqueue(item)
			rm.signalQueue()
			return
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

	if err := rm.governor.InitExecution(q.execID, q.executable.Limits); err != nil {
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
	go rm.runRuntime(rctx, entry)
	return nil
}

func (rm *runtimeManager) runRuntime(ctx context.Context, e *runtimeEntry) {
	err := e.runtime.Execute(ctx)
	rm.onExit(ctx, e, err)
}

func (rm *runtimeManager) onExit(ctx context.Context, e *runtimeEntry, err error) {
	now := time.Now()
	status := types.StatusCompleted

	if err != nil {
		status = types.StatusFailed
	}

	_ = rm.ledger.UpdateExecution(ctx, e.execID, status, &now)

	rm.governor.CleanupExecution(e.execID)
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

//
// ===== debug =====
//

func (rm *runtimeManager) EnableDebug(execID id.ID) error {
	e, err := rm.get(execID)
	if err != nil {
		return err
	}

	e.runtime.Block()
	if !e.runtime.IsBlocked() {
		return ErrStepInProgress
	}

	e.debug.Store(true)
	rm.running.Add(-1)
	return nil
}

func (rm *runtimeManager) Step(execID id.ID) error {
	e, err := rm.get(execID)
	if err != nil {
		return err
	}

	if !e.debug.Load() {
		return ErrNotInDebugMode
	}

	select {
	case e.stepGate <- struct{}{}:
		return nil
	default:
		return ErrStepInProgress
	}
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
	params map[string]any,
) Runtime {
	panic("runtime factory not implemented")
}
