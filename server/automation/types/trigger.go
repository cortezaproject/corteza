package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
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
		TriggerID  []string `json:"triggerID"`
		WorkflowID []string `json:"workflowID"`

		EventType    string `json:"eventType"`
		ResourceType string `json:"resourceType"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

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
