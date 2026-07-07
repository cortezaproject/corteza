package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	UserGroupPath struct {
		SelfID uint64 `json:"selfID,string"`
		Name   string `json:"name"`
	}

	UserGroupFilter struct {
		UserGroupID []string `json:"userGroupID"`
		MemberID    uint64   `json:"memberID,string"`

		Query string `json:"query"`

		Handle string `json:"handle"`
		Name   string `json:"name"`

		Deleted  filter.State `json:"deleted"`
		Archived filter.State `json:"archived"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Check func(*UserGroup) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

// FindByHandle finds userGroup by it's handle
func (set UserGroupSet) FindByHandle(handle string) *UserGroup {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}
