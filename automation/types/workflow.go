package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
		Meta    *WorkflowMeta     `json:"meta"`
		Enabled bool              `json:"enabled"`

		Trace bool `json:"trace"`

		// how much time do we keep completed sessions (in sec)
		KeepSessions int `json:"keepSessions"`

		// Initial input scope
		Scope Variables `json:"scope"`

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
		Name        string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	// WorkflowStep describes one workflow step
	WorkflowStep struct {
		ID   uint64           `json:"stepID,string"`
		Kind WorkflowStepKind `json:"kind"`

		// reference to function or subprocess (workflow)
		Ref string `json:"ref"`

		// set of expressions to evaluate or pass to function
		// invalid for for kind=~gateway:*
		Arguments []*WorkflowExpression `json:"arguments,string"`

		// only valid when kind=function
		Results []*WorkflowExpression `json:"results,string"`

		Meta WorkflowStepMeta `json:"meta"`
	}

	WorkflowStepMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	// WorkflowPath defines connection between two workflow steps
	WorkflowPath struct {
		ParentID uint64 `json:"parentID,string"`
		ChildID  uint64 `json:"childID,string"`

		// test expression for gateway paths
		Test string           `json:"test,string"`
		Meta WorkflowPathMeta `json:"meta"`
	}

	WorkflowPathMeta struct {
		Label       string                 `json:"label"`
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	WorkflowStepKind string
)

const (
	WorkflowStepKindExpressions WorkflowStepKind = "expressions" // ref
	WorkflowStepKindGateway     WorkflowStepKind = "gateway"     // ref=join|fork|excl|incl
	WorkflowStepKindFunction    WorkflowStepKind = "function"    // ref=<function ref>
	//WorkflowStepKindLoop        WorkflowStepKind = "loop"
	//WorkflowStepKindSubprocess  WorkflowStepKind = "subprocess"
	//WorkflowStepKindPrompt      WorkflowStepKind = "prompt" // ref = client function
	//WorkflowStepKindNotify      WorkflowStepKind = "notify" // ref = error, warning, info
	//WorkflowStepKindEvent       WorkflowStepKind = "event" // ref = ??
)

// Resource returns a resource ID for this type
func (r *Workflow) RBACResource() rbac.Resource {
	return WorkflowRBACResource.AppendID(r.ID)
}

func ParseWorkflowMeta(ss []string) (p *WorkflowMeta, err error) {
	p = &WorkflowMeta{}
	return p, parseStringsInput(ss, p)
}

func ParseWorkflowStepSet(ss []string) (p WorkflowStepSet, err error) {
	p = WorkflowStepSet{}
	return p, parseStringsInput(ss, &p)
}

func ParseWorkflowPathSet(ss []string) (p WorkflowPathSet, err error) {
	p = WorkflowPathSet{}
	return p, parseStringsInput(ss, &p)
}

func parseStringsInput(ss []string, p interface{}) (err error) {
	if len(ss) == 0 {
		return
	}

	return json.Unmarshal([]byte(ss[0]), &p)
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
