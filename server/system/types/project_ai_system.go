package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ProjectGroupFilter struct {
		ProjectGroupID []uint64     `json:"projectGroupID"`
		TenantID       uint64       `json:"tenantID,string,omitempty"`
		ProjectID      uint64       `json:"projectID,string,omitempty"`
		Handle         string       `json:"handle,omitempty"`
		Query          string       `json:"query,omitempty"`
		Deleted        filter.State `json:"deleted"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*ProjectGroup) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	ProjectGroupEntryFilter struct {
		ProjectGroupID uint64 `json:"projectGroupID,string,omitempty"`
		ResourceRef    string `json:"resourceRef,omitempty"`
		Limit          uint   `json:"-"`
	}
)
