package types

import (
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"time"
)

type Trigger struct {
	ID      uint64 `json:"triggerID,string"`
	Enabled bool   `json:"enabled"`

	WorkflowID uint64 `json:"workflowID,string"`
	// Start workflow on this step. If 0, find first (only) orphan
	StepID uint64 `json:"stepID,string"`

	// Resource type that can trigger the workflow
	ResourceType string

	// Event type that can trigger the workflow
	EventType string

	// Trigger constraints
	Constraints []TriggerConstraint

	// Initial input scope,
	// will be merged merged with workflow variables
	Input wfexec.Variables

	Labels map[string]string `json:"labels,omitempty"`

	OwnedBy   uint64     `json:"ownedBy,string"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	CreatedBy uint64     `json:"createdBy,string" `
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	DeletedBy uint64     `json:"deletedBy,string,omitempty"`
}

type TriggerConstraint struct {
	Name   string   `json:"name"`
	Op     string   `json:"op,omitempty"`
	Values []string `json:"values,omitempty"`
}

type TriggerFilter struct {
	TriggerID  []uint64 `json:"triggerID"`
	WorkflowID []uint64 `json:"workflowID"`

	EventType    string `json:"eventType"`
	ResourceType string `json:"resourceType"`

	Deleted filter.State `json:"deleted"`

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

func ParseTriggerConstraintSet(ss []string) (p TriggerConstraintSet, err error) {
	p = TriggerConstraintSet{}
	return p, parseStringsInput(ss, &p)
}

// Resource returns a resource ID for this type
func (r *Trigger) RBACResource() rbac.Resource {
	return TriggerRBACResource.AppendID(r.ID)
}
