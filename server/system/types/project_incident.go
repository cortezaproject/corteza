package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ProjectIncidentFilter struct {
		IncidentID []uint64 `json:"incidentID"`
		TenantID   uint64   `json:"tenantID,string,omitempty"`
		ProjectID  uint64   `json:"projectID,string,omitempty"`
		RevisionID uint64   `json:"revisionID,string,omitempty"`
		Status     string   `json:"status"`
		Query      string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found; it can
		// modify the resource and return false if store should not return it.
		Check func(*ProjectIncident) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
