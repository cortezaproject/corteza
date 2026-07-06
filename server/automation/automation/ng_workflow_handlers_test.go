package automation

import (
	"context"
	"testing"

	atypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/stretchr/testify/require"
)

type mockWorkflowRunner struct {
	lookupByID     func(context.Context, uint64) (*atypes.Workflow, error)
	lookupByHandle func(context.Context, string) (*atypes.Workflow, error)
	exec           func(context.Context, uint64, atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error)
}

func (m mockWorkflowRunner) LookupByID(ctx context.Context, id uint64) (*atypes.Workflow, error) {
	return m.lookupByID(ctx, id)
}

func (m mockWorkflowRunner) LookupByHandle(ctx context.Context, h string) (*atypes.Workflow, error) {
	return m.lookupByHandle(ctx, h)
}

func (m mockWorkflowRunner) Exec(ctx context.Context, id uint64, p atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error) {
	return m.exec(ctx, id, p)
}

func testWorkflow(id uint64) *atypes.Workflow {
	return &atypes.Workflow{ID: id, Handle: "wf", Enabled: true}
}

func TestNgWorkflowHandler_Exec_ByID(t *testing.T) {
	var (
		req       = require.New(t)
		result, _ = expr.NewVars(map[string]interface{}{"answer": 42})

		gotID uint64
		gotP  atypes.WorkflowExecParams
	)

	h := ngWorkflowHandler{wf: mockWorkflowRunner{
		lookupByID: func(_ context.Context, id uint64) (*atypes.Workflow, error) {
			return testWorkflow(id), nil
		},
		exec: func(_ context.Context, id uint64, p atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error) {
			gotID, gotP = id, p
			return result, 99, nil, nil
		},
	}}

	in := &expr.Vars{}
	req.NoError(in.Set("workflow", expr.Must(expr.NewID(uint64(123)))))

	out, err := h.Exec().Handler(context.Background(), in)
	req.NoError(err)
	req.NotNil(out)
	req.True(out.Has("results"))

	req.Equal(uint64(123), gotID)
	req.False(gotP.Async, "workflow must run synchronously")
	req.True(gotP.Wait, "must wait for the workflow result")
	req.NotNil(gotP.Input, "input must default to empty vars, never nil")
}

func TestNgWorkflowHandler_Exec_ByHandle(t *testing.T) {
	var (
		req       = require.New(t)
		result, _ = expr.NewVars(nil)
		gotHandle string
	)

	h := ngWorkflowHandler{wf: mockWorkflowRunner{
		lookupByHandle: func(_ context.Context, handle string) (*atypes.Workflow, error) {
			gotHandle = handle
			return testWorkflow(7), nil
		},
		exec: func(_ context.Context, _ uint64, _ atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error) {
			return result, 0, nil, nil
		},
	}}

	in := &expr.Vars{}
	req.NoError(in.Set("workflow", expr.Must(expr.NewHandle("my-sub"))))

	_, err := h.Exec().Handler(context.Background(), in)
	req.NoError(err)
	req.Equal("my-sub", gotHandle)
}

func TestNgWorkflowHandler_Exec_NumericStringIsID(t *testing.T) {
	var (
		req       = require.New(t)
		result, _ = expr.NewVars(nil)

		gotID           uint64
		lookupByHandled bool
	)

	h := ngWorkflowHandler{wf: mockWorkflowRunner{
		lookupByID: func(_ context.Context, id uint64) (*atypes.Workflow, error) {
			gotID = id
			return testWorkflow(id), nil
		},
		lookupByHandle: func(_ context.Context, _ string) (*atypes.Workflow, error) {
			lookupByHandled = true
			return nil, nil
		},
		exec: func(_ context.Context, _ uint64, _ atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error) {
			return result, 0, nil, nil
		},
	}}

	// the workflow selector emits the ID as a string; it must resolve by ID, not handle
	in := &expr.Vars{}
	req.NoError(in.Set("workflow", expr.Must(expr.NewString("123456"))))

	_, err := h.Exec().Handler(context.Background(), in)
	req.NoError(err)
	req.Equal(uint64(123456), gotID)
	req.False(lookupByHandled, "a numeric string must resolve by ID, not handle")
}

func TestNgWorkflowHandler_Exec_RunsWorkflowWithTriggers(t *testing.T) {
	var (
		req       = require.New(t)
		result, _ = expr.NewVars(nil)
		gotID     uint64
	)

	// a regular (trigger-ful) workflow runs the same as a manual run — no gating
	h := ngWorkflowHandler{wf: mockWorkflowRunner{
		lookupByID: func(_ context.Context, id uint64) (*atypes.Workflow, error) {
			return &atypes.Workflow{ID: id, Handle: "regular", Enabled: true}, nil
		},
		exec: func(_ context.Context, id uint64, _ atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error) {
			gotID = id
			return result, 0, nil, nil
		},
	}}

	in := &expr.Vars{}
	req.NoError(in.Set("workflow", expr.Must(expr.NewID(uint64(1)))))

	_, err := h.Exec().Handler(context.Background(), in)
	req.NoError(err)
	req.Equal(uint64(1), gotID)
}

func TestNgWorkflowHandler_Exec_RequiresWorkflow(t *testing.T) {
	req := require.New(t)

	h := ngWorkflowHandler{wf: mockWorkflowRunner{}}

	_, err := h.Exec().Handler(context.Background(), &expr.Vars{})
	req.Error(err)
}

func TestNgWorkflowHandler_Exec_PassesInput(t *testing.T) {
	var (
		req       = require.New(t)
		result, _ = expr.NewVars(nil)
		input, _  = expr.NewVars(map[string]interface{}{"foo": "bar"})
		gotInput  *expr.Vars
	)

	h := ngWorkflowHandler{wf: mockWorkflowRunner{
		lookupByID: func(_ context.Context, id uint64) (*atypes.Workflow, error) {
			return testWorkflow(id), nil
		},
		exec: func(_ context.Context, _ uint64, p atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error) {
			gotInput = p.Input
			return result, 0, nil, nil
		},
	}}

	in := &expr.Vars{}
	req.NoError(in.Set("workflow", expr.Must(expr.NewID(uint64(5)))))
	req.NoError(in.Set("input", input))

	_, err := h.Exec().Handler(context.Background(), in)
	req.NoError(err)
	req.NotNil(gotInput)
	req.True(gotInput.Has("foo"))
}
