package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	TenantFilter struct {
		TenantID []string     `json:"tenantID"`
		Handle   string       `json:"handle"`
		Status   TenantStatus `json:"status"`
		Query    string       `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Check func(*Tenant) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
