package manager

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"go.uber.org/zap"
)

func init() {
	id.Init(cli.Context())
}

// ---- mocks ----

type mockRegistry struct {
	exec types.Executable
	err  error
}

func (m *mockRegistry) GetExecutable(_ context.Context, _ id.ID, _ int) (types.Executable, error) {
	return m.exec, m.err
}

type mockLedger struct {
	mu              sync.Mutex
	registerCalled  int
	completedCalled int
	failedCalled    int
	lastErr         error
	registerErr     error
	completedErr    error
	failedErr       error
}

func (m *mockLedger) RegisterExecution(_ context.Context, _, _ id.ID, _ int, _ types.ExecutionParams) error {
	m.mu.Lock()
	m.registerCalled++
	m.mu.Unlock()
	return m.registerErr
}
func (m *mockLedger) ExecutionCompleted(_ context.Context, _, _ id.ID, _ int) error {
	m.mu.Lock()
	m.completedCalled++
	m.mu.Unlock()
	return m.completedErr
}
func (m *mockLedger) ExecutionFailed(_ context.Context, _, _ id.ID, _ int, err error) error {
	m.mu.Lock()
	m.failedCalled++
	m.lastErr = err
	m.mu.Unlock()
	return m.failedErr
}
func (m *mockLedger) StepStarted(_ context.Context, _, _, _ id.ID, _ int) error { return nil }
func (m *mockLedger) StepCompleted(_ context.Context, _, _, _ id.ID, _ int, _ any) error {
	return nil
}
func (m *mockLedger) StepFailed(_ context.Context, _, _, _ id.ID, _ int, _ error) error { return nil }
func (m *mockLedger) ExecutionPaused(_ context.Context, _, _, _ id.ID, _, _ int, _ error) error {
	return nil
}
func (m *mockLedger) RecordFrame(_ context.Context, _, _ id.ID, _ int, _ types.StackFrame) error {
	return nil
}
func (m *mockLedger) IsExecutableInUse(_ context.Context, _ id.ID, _ int) (bool, error) {
	return false, nil
}

type mockGovernor struct {
	mu           sync.Mutex
	addCalled    int
	removeCalled int
	addErr       error
}

func (m *mockGovernor) AddExecution(_ id.ID, _ int, _ types.Budget, _ types.RateLimit) error {
	m.mu.Lock()
	m.addCalled++
	m.mu.Unlock()
	return m.addErr
}
func (m *mockGovernor) RemoveExecution(_ id.ID) {
	m.mu.Lock()
	m.removeCalled++
	m.mu.Unlock()
}
func (m *mockGovernor) Request(_ id.ID, _ int) (<-chan struct{}, error) {
	ch := make(chan struct{})
	close(ch)
	return ch, nil
}

// ---- step handlers ----

type noopHandler struct{}

func (h *noopHandler) ExecN(_ context.Context, _ *types.ExecRequest) (types.ExecResponse, error) {
	return nil, nil
}

type alwaysFailHandler struct{ err error }

func (h *alwaysFailHandler) ExecN(_ context.Context, _ *types.ExecRequest) (types.ExecResponse, error) {
	return nil, h.err
}

type blockHandler struct{ release <-chan struct{} }

func (h *blockHandler) ExecN(_ context.Context, _ *types.ExecRequest) (types.ExecResponse, error) {
	<-h.release
	return nil, nil
}

// ---- builder helpers ----

func singleStepExec() types.Executable {
	return types.Executable{
		ID:       id.MustNumID(id.Next()),
		Revision: 1,
		Steps: []types.Step{
			{
				ID:      id.MustNumID(id.Next()),
				Handle:  "step1",
				Kind:    "action",
				Handler: &noopHandler{},
			},
		},
	}
}

func newRM(t *testing.T, reg Registry, led Ledger, gov Governor, cfg Config) *runtimeManager {
	t.Helper()
	rm, err := RuntimeManager(context.Background(), zap.NewNop(), reg, led, gov, cfg)
	if err != nil {
		t.Fatalf("RuntimeManager: %v", err)
	}
	return rm
}

func defaultConfig() Config {
	return Config{
		MaxConcurrent: 10,
		MaxQueued:     0,
		SlotTimeout:   5 * time.Second,
	}
}

func waitDone(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for execution to finish")
	}
}
