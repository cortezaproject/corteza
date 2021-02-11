package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"time"
)

type (
	// Workflow represents entire workflow definition
	Workflow struct {
		ID      uint64            `json:"workflowID,string"`
		Handle  string            `json:"handle"`
		Labels  map[string]string `json:"labels,omitempty"`
		Meta    *WorkflowMeta     `json:"meta,omitempty"`
		Enabled bool              `json:"enabled"`

		Trace bool `json:"trace"`

		// how much time do we keep completed sessions (in sec)
		KeepSessions int `json:"keepSessions"`

		// Initial input scope
		Scope *expr.Vars `json:"scope"`

		Steps WorkflowStepSet `json:"steps"`
		Paths WorkflowPathSet `json:"paths"`

		RunAs uint64 `json:"runAs,string"`

		OwnedBy   uint64     `json:"ownedBy,string"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	WorkflowFilter struct {
		WorkflowID []uint64 `json:"workflowID"`

		Query string `json:"query"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Workflow) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	WorkflowMeta struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

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
	WorkflowStepKindMessage     WorkflowStepKind = "message"       // ref = error|warning|info, ...
	WorkflowStepKindPrompt      WorkflowStepKind = "prompt"        // ref = <client function>
	WorkflowStepKindErrHandler  WorkflowStepKind = "error-handler" // no ref
	WorkflowStepKindVisual      WorkflowStepKind = "visual"        // ref = <*>
	WorkflowStepKindDebug       WorkflowStepKind = "debug"         // ref = <*>
	//WorkflowStepKindSubprocess  WorkflowStepKind = "subprocess"
	//WorkflowStepKindEvent       WorkflowStepKind = "event" // ref = ??
)

// Resource returns a resource ID for this type
func (r *Workflow) RBACResource() rbac.Resource {
	return WorkflowRBACResource.AppendID(r.ID)
}

func (vv *WorkflowMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = WorkflowMeta{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("can not scan '%v' into WorkflowMeta: %w", string(b), err)
		}
	}

	return nil
}

// Scan on WorkflowMeta gracefully handles conversion from NULL
func (vv *WorkflowMeta) Value() (driver.Value, error) {
	if vv == nil {
		return []byte("null"), nil
	}

	return json.Marshal(vv)
}

func (t WorkflowPath) GetExpr() string              { return t.Expr }
func (t *WorkflowPath) SetEval(eval expr.Evaluable) { t.eval = eval }
func (t WorkflowPath) Eval(ctx context.Context, scope *expr.Vars) (interface{}, error) {
	return t.eval.Eval(ctx, scope)
}
func (t WorkflowPath) Test(ctx context.Context, scope *expr.Vars) (bool, error) {
	return t.eval.Test(ctx, scope)
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

// Scan on WorkflowStepSet gracefully handles conversion from NULL
func (vv WorkflowStepSet) Value() (driver.Value, error) {
	return json.Marshal(vv)
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
