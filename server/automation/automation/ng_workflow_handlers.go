package automation

import (
	"context"
	"strconv"

	atypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/expr"
)

type (
	ngWorkflowConstructSvc interface {
		AddFunctions(ff ...atypes.ConstructFunction)
		AddTriggers(tt ...atypes.ConstructTrigger)
	}

	workflowRunner interface {
		LookupByID(ctx context.Context, workflowID uint64) (*atypes.Workflow, error)
		LookupByHandle(ctx context.Context, handle string) (*atypes.Workflow, error)
		Exec(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error)
	}

	ngWorkflowHandler struct {
		reg ngWorkflowConstructSvc
		wf  workflowRunner
	}
)

func NgWorkflowHandler(reg ngWorkflowConstructSvc, wf workflowRunner) *ngWorkflowHandler {
	h := &ngWorkflowHandler{
		reg: reg,
		wf:  wf,
	}

	h.register()
	return h
}

func (h ngWorkflowHandler) register() {
	h.reg.AddFunctions(
		h.Exec(),
	)
}

func (h ngWorkflowHandler) Exec() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "workflowExec",
		Kind:   "function",
		Groups: []string{"Workflows"},
		Labels: map[string]string{"workflow": "step,workflow"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Run Workflow",
			Description: "Run a workflow and return its results",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "share-2"},
		},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "workflow",
				Types:        []string{"ID", "Handle"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Workflow",
					Description: "Workflow to run, referenced by ID or handle.",
				},
			},
			{
				ArgumentName: "input",
				Types:        []string{"Vars"},
				Meta: &atypes.ParamMeta{
					Label:       "Input",
					Description: "Scope passed to the workflow as its input.",
				},
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "results",
				Types:        []string{"Vars"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Meta: atypes.ConstructSegmentMeta{},
			Sections: []atypes.ConstructSection{{
				Meta: atypes.ConstructSectionMeta{},
				Elements: []atypes.SectionElement{
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Workflow", Argument: "workflow", Required: true}},
					{Input: atypes.SectionElementInput{Type: "Expression", Label: "Input", Argument: "input"}},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var wf *atypes.Workflow

			if !in.Has("workflow") {
				return nil, errors.InvalidData("workflow argument is required")
			}

			switch v := expr.Must(expr.Select(in, "workflow")).Get().(type) {
			case uint64:
				wf, err = h.wf.LookupByID(ctx, v)
			case int64:
				wf, err = h.wf.LookupByID(ctx, uint64(v))
			case string:
				if id, perr := strconv.ParseUint(v, 10, 64); perr == nil {
					wf, err = h.wf.LookupByID(ctx, id)
				} else {
					wf, err = h.wf.LookupByHandle(ctx, v)
				}
			default:
				return nil, errors.InvalidData("workflow must be referenced by ID or handle")
			}
			if err != nil {
				return nil, err
			}
			if wf == nil {
				return nil, errors.NotFound("workflow not found")
			}

			p := atypes.WorkflowExecParams{
				// run synchronously and wait for the result, same as a manual run
				Async: false,
				Wait:  true,
				Input: expr.EmptyVars(),
			}

			if in.Has("input") {
				if vars, ok := expr.Must(expr.Select(in, "input")).(*expr.Vars); ok && vars != nil {
					p.Input = vars
				}
			}

			result, _, _, err := h.wf.Exec(ctx, wf.ID, p)
			if err != nil {
				return nil, err
			}

			out = &expr.Vars{}
			if result != nil {
				if err = expr.Assign(out, "results", result); err != nil {
					return nil, err
				}
			}

			return out, nil
		},
	}
}
