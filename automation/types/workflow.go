package types

import (
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

		// Collection of issues from the last parse
		Issues WorkflowIssueSet `json:"issues,omitempty"`

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

	WorkflowIssue struct {
		// url encoded location of the error:
		Culprit     map[string]int `json:"culprit"`
		Description string         `json:"description"`
	}
)

// Resource returns a resource ID for this type
func (r *Workflow) RBACResource() rbac.Resource {
	return WorkflowRBACResource.AppendID(r.ID)
}

// CheckDeferred returns true if any of the steps is deferred.
//
// Workflow is considered deferred when delay or prompt step types are used.
// Deferred workflows can not short-circuit triggers or prevent creation/update on before triggers
//
// @todo add flag on workflow to explicitly mark workflow as deferred even when there are no delay or prompt steps
func (r Workflow) CheckDeferred() bool {
	return r.Steps.HasDeferred()
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

// Scan on WorkflowStepSet gracefully handles conversion from NULL
func (set WorkflowIssueSet) Value() (driver.Value, error) {
	return json.Marshal(set)
}

func (set *WorkflowIssueSet) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*set = WorkflowIssueSet{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, set); err != nil {
			return fmt.Errorf("can not scan '%v' into WorkflowIssueSet: %w", string(b), err)
		}
	}

	return nil
}

func (set WorkflowIssueSet) Error() string {
	switch len(set) {
	case 0:
		return fmt.Sprintf("no workflow issue found")
	case 1:
		return fmt.Sprintf("1 workflow issue found")
	default:
		return fmt.Sprintf("%d workflow issues found", len(set))
	}
}

func (set WorkflowIssueSet) Append(err error, culprit map[string]int) WorkflowIssueSet {
	if culprit == nil {
		culprit = make(map[string]int)
	}

	return append(set, &WorkflowIssue{
		Culprit:     culprit,
		Description: err.Error(),
	})
}

func (set WorkflowIssueSet) SetCulprit(name string, pos int) WorkflowIssueSet {
	for i := range set {
		set[i].Culprit[name] = pos
	}

	return set
}
