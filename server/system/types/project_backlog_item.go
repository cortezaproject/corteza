package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ProjectBacklogItemFilter struct {
		BacklogItemID []uint64 `json:"backlogItemID"`
		TenantID      uint64   `json:"tenantID,string,omitempty"`
		ProjectID     uint64   `json:"projectID,string,omitempty"`
		EventID       uint64   `json:"eventID,string,omitempty"`
		Category      string   `json:"category"`
		Status        string   `json:"status"`
		Query         string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found; it can
		// modify the resource and return false if store should not return it.
		Check func(*ProjectBacklogItem) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
