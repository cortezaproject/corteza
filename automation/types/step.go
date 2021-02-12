package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

type (
	// WorkflowStep describes one workflow step
	WorkflowStep struct {
		ID   uint64           `json:"stepID,string"`
		Kind WorkflowStepKind `json:"kind"`

		// reference to function or subprocess (workflow)
		Ref string `json:"ref"`

		// set of expressions to evaluate, test or pass to function
		// invalid for for kind=~gateway:*
		Arguments []*Expr `json:"arguments"`

		// only valid when kind=function
		Results []*Expr `json:"results"`

		Meta WorkflowStepMeta `json:"meta,omitempty"`

		Labels map[string]string `json:"labels,omitempty"`
	}

	WorkflowStepMeta struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	// WorkflowPath defines connection between two workflow steps
	WorkflowPath struct {
		// Expression to evaluate over the input variables; results will be set to scope under variable Name
		Expr string `json:"expr,omitempty"`

		eval expr.Evaluable

		ParentID uint64           `json:"parentID,string"`
		ChildID  uint64           `json:"childID,string"`
		Meta     WorkflowPathMeta `json:"meta,omitempty"`
	}

	WorkflowPathMeta struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	WorkflowStepKind string
)

const (
	WorkflowStepKindExpressions WorkflowStepKind = "expressions"   // ref
	WorkflowStepKindGateway     WorkflowStepKind = "gateway"       // ref = join|fork|excl|incl
	WorkflowStepKindFunction    WorkflowStepKind = "function"      // ref = <function ref>
	WorkflowStepKindIterator    WorkflowStepKind = "iterator"      // ref = <iterator function ref>
	WorkflowStepKindError       WorkflowStepKind = "error"         // no ref
	WorkflowStepKindTermination WorkflowStepKind = "termination"   // no ref
	WorkflowStepKindPrompt      WorkflowStepKind = "prompt"        // ref = <client function>
	WorkflowStepKindErrHandler  WorkflowStepKind = "error-handler" // no ref
	WorkflowStepKindVisual      WorkflowStepKind = "visual"        // ref = <*>
	WorkflowStepKindDebug       WorkflowStepKind = "debug"         // ref = <*>
)

// Scan on WorkflowStepSet gracefully handles conversion from NULL
func (vv WorkflowStepSet) Value() (driver.Value, error) {
	return json.Marshal(vv)
}

func (vv *WorkflowStepSet) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = WorkflowStepSet{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("can not scan '%v' into WorkflowStepSet: %w", string(b), err)
		}
	}

	return nil
}

func (t WorkflowPath) GetExpr() string              { return t.Expr }
func (t *WorkflowPath) SetEval(eval expr.Evaluable) { t.eval = eval }
func (t WorkflowPath) Eval(ctx context.Context, scope *expr.Vars) (interface{}, error) {
	return t.eval.Eval(ctx, scope)
}
func (t WorkflowPath) Test(ctx context.Context, scope *expr.Vars) (bool, error) {
	return t.eval.Test(ctx, scope)
}

func (vv *WorkflowPathSet) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = WorkflowPathSet{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("can not scan '%v' into WorkflowPathSet: %w", string(b), err)
		}
	}

	return nil
}

// Scan on WorkflowPathSet gracefully handles conversion from NULL
func (vv WorkflowPathSet) Value() (driver.Value, error) {
	return json.Marshal(vv)
}
