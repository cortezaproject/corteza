package manager

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/runtime"
	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/modern-go/reflect2"
	"go.uber.org/zap"
)

//
// ===== errors =====
//

var (
	ErrSystemDraining     = errors.New("manager: system is draining")
	ErrSlotTimeout        = errors.New("manager: concurrency slot timeout")
	ErrExecutableNotFound = errors.New("manager: executable not found")
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

type Ledger interface {
	StepStarted(ctx context.Context, executableID, executionID, stepID id.ID, rev int) error
	StepCompleted(ctx context.Context, executableID, executionID, stepID id.ID, rev int, out any) error
	StepFailed(ctx context.Context, executableID, executionID, stepID id.ID, rev int, err error) error

	ExecutionCompleted(ctx context.Context, executableID, executionID id.ID, rev int) error
	ExecutionFailed(ctx context.Context, executableID, executionID id.ID, rev int, err error) error

	RecordFrame(ctx context.Context, executableID, executionID id.ID, rev int, frame types.StackFrame) error

	IsExecutableInUse(ctx context.Context, executableID id.ID, rev int) (bool, error)

	RegisterExecution(ctx context.Context, executableID, executionID id.ID, rev int) error
}

type Governor interface {
	AddExecution(execID id.ID, maxOpPerRequest int, budget types.Budget, rate types.RateLimit) error
	RemoveExecution(execID id.ID)
	Request(execID id.ID, ops int) (<-chan struct{}, error)
}

type Runtime interface {
	Start(ctx context.Context, global *expr.Vars) error
	Stop()
	Block()
	Resume()
	IsBlocked() bool
}

//
// ===== internal types =====
//

type RuntimeEntry struct {
	execID       id.ID
	executableID id.ID
	revision     int

	executable  types.Executable
	globalState *expr.Vars
	entryPoint  string

	runtime Runtime
	cancel  context.CancelFunc

	Done    chan struct{}
	Started chan struct{}
	Err     error

	invoker auth.Identifiable
	runner  auth.Identifiable
}

//
// ===== runtime manager =====
//

type runtimeManager struct {
	log *zap.Logger

	registry Registry
	ledger   Ledger
	governor Governor
	config   Config

	executions map[id.ID]*RuntimeEntry
	mu         sync.RWMutex

	queue    []*RuntimeEntry
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
	log *zap.Logger,
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
		log:        log,
		registry:   reg,
		ledger:     led,
		governor:   gov,
		config:     cfg,
		executions: make(map[id.ID]*RuntimeEntry),
		queue:      make([]*RuntimeEntry, 0),
		qSignal:    make(chan struct{}, 1),
		slots:      make(chan struct{}, cfg.MaxConcurrent),
	}

	go rm.watchQueue(ctx)
	return rm, nil
}

//
// ===== admission =====
//

func (rm *runtimeManager) Start(ctx context.Context, executableID id.ID, revision int, params types.ExecutionParams) (id.ID, error) {
	if rm.draining.Load() {
		return id.Zero(), ErrSystemDraining
	}

	global, err := rm.validateGlobalState(params.Input)
	if err != nil {
		return id.Zero(), err
	}

	executable, err := rm.registry.GetExecutable(ctx, executableID, revision)
	if err != nil {
		return id.Zero(), ErrExecutableNotFound
	}

	eid := id.MustNumID(id.Next())

	err = rm.ledger.RegisterExecution(ctx, executableID, eid, revision)
	if err != nil {
		return id.Zero(), err
	}

	invoker := auth.GetIdentityFromContext(ctx)
	runner := executable.RunAs

	if reflect2.IsNil(runner) {
		runner = invoker
	}

	entry := &RuntimeEntry{
		execID:       eid,
		executableID: executableID,
		revision:     revision,
		executable:   executable,
		globalState:  global,
		entryPoint:   params.EntryPoint,

		Done:    make(chan struct{}),
		Started: make(chan struct{}),

		runner:  runner,
		invoker: invoker,
	}

	rm.mu.Lock()
	rm.executions[eid] = entry
	rm.mu.Unlock()

	if err := rm.enqueue(entry); err != nil {
		rm.mu.Lock()
		delete(rm.executions, eid)
		rm.mu.Unlock()
		return id.Zero(), err
	}

	rm.signalQueue()
	return eid, nil
}

func (rm *runtimeManager) Stop(execID id.ID) error {
	rm.mu.RLock()
	e, ok := rm.executions[execID]
	rm.mu.RUnlock()

	if ok {
		if e.cancel != nil {
			e.cancel()
		}
		if e.runtime != nil {
			e.runtime.Stop()
		}
		return nil
	}

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

func (rm *runtimeManager) enqueue(entry *RuntimeEntry) error {
	rm.queueMux.Lock()
	defer rm.queueMux.Unlock()

	if rm.config.MaxQueued > 0 && len(rm.queue) >= rm.config.MaxQueued {
		return ErrQueueFull
	}

	rm.queue = append(rm.queue, entry)
	return nil
}

func (rm *runtimeManager) dequeue() (*RuntimeEntry, bool) {
	rm.queueMux.Lock()
	defer rm.queueMux.Unlock()

	if len(rm.queue) == 0 {
		return nil, false
	}

	entry := rm.queue[0]
	rm.queue = rm.queue[1:]
	return entry, true
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
		entry, ok := rm.dequeue()
		if !ok {
			return
		}

		// @todo do we want to do something more with the context
		execCtx := context.Background()

		// Encode runner into execution context
		// runner is used as identity and for access control
		execCtx = auth.SetIdentityToContext(execCtx, entry.runner)

		// Encode invoker into execution context
		// invoker is used
		// @todo :)
		// execCtx = context.WithValue(execCtx, workflowInvokerCtxKey{}, ent)

		err := rm.startQueued(execCtx, entry)
		if err != nil {
			panic(err)
		}
	}
}

//
// ===== runtime start =====
//

func (rm *runtimeManager) startQueued(ctx context.Context, entry *RuntimeEntry) error {
	select {
	case rm.slots <- struct{}{}:
	case <-time.After(rm.config.SlotTimeout):
		return ErrSlotTimeout
	case <-ctx.Done():
		return ctx.Err()
	}

	if err := rm.governor.AddExecution(entry.execID, 0, types.Budget{}, types.RateLimit{}); err != nil {
		rm.releaseSlot()
		return err
	}

	rctx, cancel := context.WithCancel(ctx)
	rt := rm.createRuntime(entry.execID, entry.executable, entry.entryPoint)

	rm.mu.Lock()
	entry.runtime = rt
	entry.cancel = cancel
	rm.mu.Unlock()

	close(entry.Started)

	rm.running.Add(1)
	go rm.runRuntime(rctx, entry, entry.globalState)
	return nil
}

func (rm *runtimeManager) runRuntime(ctx context.Context, e *RuntimeEntry, global *expr.Vars) {
	err := e.runtime.Start(ctx, global)
	rm.onExit(ctx, e, err)
}

func (rm *runtimeManager) onExit(ctx context.Context, e *RuntimeEntry, err error) {
	e.Err = err

	status := types.StatusCompleted
	if err != nil {
		status = types.StatusFailed
	}

	switch status {
	case types.StatusCompleted:
		_ = rm.ledger.ExecutionCompleted(ctx, e.executableID, e.execID, e.revision)
	case types.StatusFailed:
		_ = rm.ledger.ExecutionFailed(ctx, e.executableID, e.execID, e.revision, err)
	}

	rm.governor.RemoveExecution(e.execID)
	rm.releaseSlot()

	rm.mu.Lock()
	close(e.Done)
	delete(rm.executions, e.execID)
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

func (rm *runtimeManager) Get(execID id.ID) (*RuntimeEntry, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	e, ok := rm.executions[execID]
	if !ok {
		return nil, ErrExecutionNotFound
	}
	return e, nil
}

func (rm *runtimeManager) createRuntime(executionID id.ID, exe types.Executable, entryPoint string) Runtime {
	return runtime.Runtime(
		executionID,
		exe,
		rm.governor,
		rm.ledger,
		entryPoint,
	)
}

func (rm *runtimeManager) validateGlobalState(global *expr.Vars) (out *expr.Vars, err error) {
	out = global

	if global == nil {
		return expr.EmptyVars(), nil
	}

	if out.Type() == (expr.Unresolved{}).Type() {
		return nil, fmt.Errorf("global state is not resolved")
	}

	return
}
