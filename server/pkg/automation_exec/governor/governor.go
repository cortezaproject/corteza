package governor

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
)

type (
	governor struct {
		paused atomic.Bool
		mu     sync.Mutex

		gates   gates
		global  globalPolicy
		exec    map[id.ID]*execPolicy
		metrics Metrics
		config  Config
		now     func() time.Time
	}
)

var (
	closedChan <-chan struct{} = func() <-chan struct{} {
		ch := make(chan struct{})
		close(ch)
		return ch
	}()
)

// The Governor decides whether execution is allowed right now.
//
// Enforces budgets and rate limits
// Applies global and per-execution constraints
// Can pause the entire system
// Hands out permission gates for operations
func Governor(ctx context.Context) (svc *governor) {
	svc = &governor{
		gates: gates{
			globalPause: newGate(),
			exec:        make(map[id.ID]*execGates),
		},
		exec: make(map[id.ID]*execPolicy),
		config: Config{
			WatcherFallbackInterval: 5 * time.Second,
		},
		now: time.Now,
	}

	svc.watch(ctx)

	return
}

// AddExecution prepares the state for the given execution
//
// The function must be called before any other function regarding the execution.
func (g *governor) AddExecution(execID id.ID, maxOpPerRequest int, budget Budget, rate RateLimit) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.exec[execID]; ok {
		return ErrExecutionExists
	}

	if exceeds(maxOpPerRequest, g.global.budget) ||
		exceeds(maxOpPerRequest, g.global.rate) {
		return ErrStepTooExpensive
	}

	execBudgetWindow := newWindow(budget.MaxOps, budget.Window, g.now())
	execRateWindow := newWindow(rate.MaxOps, rate.Window, g.now())

	if exceeds(maxOpPerRequest, execBudgetWindow) ||
		exceeds(maxOpPerRequest, execRateWindow) {
		return ErrExecutableTooExpensive
	}

	g.exec[execID] = &execPolicy{
		budget: execBudgetWindow,
		rate:   execRateWindow,
	}

	return nil
}

// RemoveExecution removes the execution from the state
//
// It is up to the parent services to assure the execution is terminated.
func (g *governor) RemoveExecution(execID id.ID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.exec, execID)

	if eg := g.gates.exec[execID]; eg != nil {
		if eg.budget != nil {
			eg.budget.open()
		}
		if eg.rate != nil {
			eg.rate.open()
		}
		delete(g.gates.exec, execID)
	}
}

// Request evaluates an execution request for the specified cost
//
// This function assumes InitExecution has been called.
// This function intentionally does NOT validate execution existence.
// Violating this contract is a programmer error and will panic via nil dereference.
// This keeps the hot path branch-free and fast.
func (g *governor) Request(execID id.ID, ops int) (<-chan struct{}, error) {
	g.metrics.TotalRequests.Add(1)

	if ops <= 0 {
		return nil, ErrInvalidOps
	}

	now := g.now()

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.paused.Load() {
		g.metrics.BlockedByPause.Add(1)
		return g.gates.globalPause.ch, nil
	}

	ep := g.exec[execID]

	g.refreshLocked(now, ep, execID)

	if wouldExceed(ep.budget, ops) {
		g.metrics.BlockedByExecBudget.Add(1)
		return g.execBudgetGate(execID), nil
	}
	if wouldExceed(ep.rate, ops) {
		g.metrics.BlockedByExecRate.Add(1)
		return g.execRateGate(execID), nil
	}

	if wouldExceed(g.global.budget, ops) {
		g.metrics.BlockedByGlobalBudget.Add(1)
		return g.globalBudgetGate(), nil
	}
	if wouldExceed(g.global.rate, ops) {
		g.metrics.BlockedByGlobalRate.Add(1)
		return g.globalRateGate(), nil
	}

	reserve(&ep.budget, ops)
	reserve(&ep.rate, ops)
	reserve(&g.global.budget, ops)
	reserve(&g.global.rate, ops)

	g.metrics.GrantedImmediately.Add(1)
	return closedChan, nil
}

func exceeds(ops int, w windowCounter) bool {
	return w.max > 0 && ops > w.max
}

func wouldExceed(w windowCounter, ops int) bool {
	return w.max > 0 && w.used+ops > w.max
}

func reserve(w *windowCounter, ops int) {
	if w.max > 0 {
		w.used += ops
	}
}

func (g *governor) execBudgetGate(execID id.ID) <-chan struct{} {
	eg := g.gates.exec[execID]
	if eg == nil {
		eg = &execGates{}
		g.gates.exec[execID] = eg
	}
	if eg.budget == nil {
		eg.budget = newGate()
	}
	return eg.budget.ch
}

func (g *governor) execRateGate(execID id.ID) <-chan struct{} {
	eg := g.gates.exec[execID]
	if eg == nil {
		eg = &execGates{}
		g.gates.exec[execID] = eg
	}
	if eg.rate == nil {
		eg.rate = newGate()
	}
	return eg.rate.ch
}

func (g *governor) globalBudgetGate() <-chan struct{} {
	if g.gates.globalBudgetHold == nil {
		g.gates.globalBudgetHold = newGate()
	}
	return g.gates.globalBudgetHold.ch
}

func (g *governor) globalRateGate() <-chan struct{} {
	if g.gates.globalRateHold == nil {
		g.gates.globalRateHold = newGate()
	}
	return g.gates.globalRateHold.ch
}

func (g *governor) refreshLocked(now time.Time, ep *execPolicy, execID id.ID) {
	var (
		unblockGlobalBudget bool
		unblockGlobalRate   bool
		unblockExecBudget   bool
		unblockExecRate     bool
	)

	reset := func(w *windowCounter) bool {
		if w.win > 0 && !w.reset.After(now) {
			w.used = 0
			w.reset = now.Add(w.win)
			return true
		}
		return false
	}

	unblockGlobalBudget = reset(&g.global.budget)
	unblockGlobalRate = reset(&g.global.rate)
	unblockExecBudget = reset(&ep.budget)
	unblockExecRate = reset(&ep.rate)

	if unblockGlobalBudget && g.gates.globalBudgetHold != nil {
		g.gates.globalBudgetHold.open()
		g.gates.globalBudgetHold = nil
	}
	if unblockGlobalRate && g.gates.globalRateHold != nil {
		g.gates.globalRateHold.open()
		g.gates.globalRateHold = nil
	}

	if eg := g.gates.exec[execID]; eg != nil {
		if unblockExecBudget && eg.budget != nil {
			eg.budget.open()
			eg.budget = nil
		}
		if unblockExecRate && eg.rate != nil {
			eg.rate.open()
			eg.rate = nil
		}
		if eg.budget == nil && eg.rate == nil {
			delete(g.gates.exec, execID)
		}
	}
}
