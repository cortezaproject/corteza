package runtime

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/cli"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/id"
)

func init() {
	id.Init(cli.Context())
}

// ---- helpers ----

func stepID(n uint64) id.ID { return id.MustNumID(n) }

func mkExe(steps ...types.Step) types.Executable {
	return types.Executable{
		ID:       stepID(999),
		Revision: 1,
		Steps:    steps,
	}
}

// noopHandler satisfies types.StepHandler.
type noopHandler struct{}

func (h noopHandler) ExecN(_ context.Context, _ *types.ExecRequest) (types.ExecResponse, error) {
	return nil, nil
}

// mkStep creates a step with the given id, kind and children (shallow).
func mkStep(sid uint64, kind string, children ...types.Step) types.Step {
	s := types.Step{
		ID:       stepID(sid),
		Handle:   fmt.Sprintf("step%d", sid),
		Kind:     kind,
		Handler:  noopHandler{},
		Children: children,
	}
	return s
}

// ---- mockIteratorHandler ----

type mockIteratorHandler struct {
	calls   int
	maxIter int
}

func (m *mockIteratorHandler) Start(_ context.Context, _ *expr.Vars) error { return nil }
func (m *mockIteratorHandler) More(_ context.Context, _ *expr.Vars) (bool, error) {
	return m.calls < m.maxIter, nil
}
func (m *mockIteratorHandler) Next(_ context.Context, _ *expr.Vars) (*expr.Vars, error) {
	m.calls++
	v, err := expr.NewVars(nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}
func (m *mockIteratorHandler) ExecN(_ context.Context, _ *types.ExecRequest) (types.ExecResponse, error) {
	return nil, nil
}

// ---- mockGatewayHandler ----

type mockGatewayHandler struct {
	indices []int
}

func (m *mockGatewayHandler) Select(_ context.Context, _ map[string]*expr.Vars) ([]int, error) {
	return m.indices, nil
}

func (m *mockGatewayHandler) ExecN(_ context.Context, _ *types.ExecRequest) (types.ExecResponse, error) {
	return nil, nil
}
