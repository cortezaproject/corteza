package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	ProjectFilter struct {
		ProjectID []string      `json:"projectID"`
		TenantID  uint64        `json:"tenantID,string,omitempty"`
		Handle    string        `json:"handle"`
		Status    ProjectStatus `json:"status"`
		Query     string        `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Check func(*Project) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

func (m ProjectMode) Valid() bool {
	return m == ProjectModeFree || m == ProjectModeGated
}
