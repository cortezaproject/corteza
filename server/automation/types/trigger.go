package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/sql"
	"time"
)

type (
	Trigger struct {
		ID      uint64 `json:"triggerID,string"`
		Enabled bool   `json:"enabled"`

		WorkflowID uint64 `json:"workflowID,string"`
		// Start workflow on this step. If 0, find first (only) orphan
		StepID uint64 `json:"stepID,string"`

		// Resource type that can trigger the workflow
		ResourceType string `json:"resourceType"`

		// Event type that can trigger the workflow
		EventType string `json:"eventType"`

		// Trigger constraints
		Constraints TriggerConstraintSet `json:"constraints"`

		// Initial input scope,
		// will be merged merged with workflow variables
		Input *expr.Vars `json:"input"`

		Labels map[string]string `json:"labels,omitempty"`
		Meta   *TriggerMeta      `json:"meta,omitempty"`

		OwnedBy   uint64     `json:"ownedBy,string"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	TriggerConstraint struct {
		Name   string   `json:"name"`
		Op     string   `json:"op,omitempty"`
		Values []string `json:"values,omitempty"`
	}

	TriggerMeta struct {
		Description string                 `json:"description"`
		Visual      map[string]interface{} `json:"visual"`
	}

	TriggerFilter struct {
		TriggerID  []uint64 `json:"triggerID"`
		WorkflowID []uint64 `json:"workflowID"`

		EventType    string `json:"eventType"`
		ResourceType string `json:"resourceType"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Trigger) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

func ParseTriggerMeta(ss []string) (p *TriggerMeta, err error) {
	p = &TriggerMeta{}
	return p, parseStringsInput(ss, p)
}

func ParseTriggerConstraintSet(ss []string) (p TriggerConstraintSet, err error) {
	p = TriggerConstraintSet{}
	return p, parseStringsInput(ss, &p)
}

func (vv *TriggerConstraintSet) Scan(src any) error          { return sql.ParseJSON(src, vv) }
func (vv TriggerConstraintSet) Value() (driver.Value, error) { return json.Marshal(vv) }

func (vv *TriggerMeta) Scan(src any) error           { return sql.ParseJSON(src, vv) }
func (vv *TriggerMeta) Value() (driver.Value, error) { return json.Marshal(vv) }

func (set TriggerSet) FilterByWorkflowID(workflowID uint64) (vv TriggerSet) {
	// Make sure we never return nil
	vv = TriggerSet{}

	for i := range set {
		if set[i].WorkflowID == workflowID {
			vv = append(vv, set[i])
		}
	}

	return
}
