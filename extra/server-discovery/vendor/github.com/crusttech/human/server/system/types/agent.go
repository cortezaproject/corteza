package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	AgentFilter struct {
		AgentID   []string `json:"agentID"`
		ProjectID uint64   `json:"projectID,string,omitempty"`
		Handle    string   `json:"handle"`
		Status    string   `json:"status"`
		Query     string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Agent) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
