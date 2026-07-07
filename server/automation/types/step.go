package types

import (
	"context"

	"github.com/crusttech/human/server/pkg/expr"
)

// WorkflowPath defines connection between two workflow steps
type WorkflowPath struct {
	// Expression to evaluate over the input variables; results will be set to scope under variable Name
	Expr string `json:"expr,omitempty"`

	eval expr.Evaluable

	ParentID uint64           `json:"parentID,string"`
	ChildID  uint64           `json:"childID,string"`
	Meta     WorkflowPathMeta `json:"meta,omitempty"`
}

// IsDeferred fn returns true if type of step is delay or prompt
func (s WorkflowStep) IsDeferred() (is bool) {
	switch s.Kind {
	case WorkflowStepKindPrompt:
		return true
	case WorkflowStepKindDelay:
		return true
	}
	return false
}

// HasDeferred fn returns true if wf-step is delay or prompt
func (vv WorkflowStepSet) HasDeferred() bool {
	for _, s := range vv {
		if s.IsDeferred() {
			return true
		}
	}

	return false
}

func (t WorkflowPath) GetExpr() string              { return t.Expr }
func (t *WorkflowPath) SetEval(eval expr.Evaluable) { t.eval = eval }
func (t WorkflowPath) Eval(ctx context.Context, scope *expr.Vars) (interface{}, error) {
	return t.eval.Eval(ctx, scope)
}
func (t WorkflowPath) Test(ctx context.Context, scope *expr.Vars) (bool, error) {
	return t.eval.Test(ctx, scope)
}
