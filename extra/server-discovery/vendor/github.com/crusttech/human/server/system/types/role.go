package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	RoleFilter struct {
		RoleID      []string `json:"roleID"`
		TenantID    uint64   `json:"tenantID,string,omitempty"`
		ProjectID   uint64   `json:"projectID,string,omitempty"`
		MemberID    uint64   `json:"memberID,string"`
		UserGroupID uint64   `json:"userGroupID,string"`

		// @todo will migrate from MemberID/UserGroupID in a later releae
		// For internal use only
		Resource string `json:"-"`

		Query string `json:"query"`

		Handle string `json:"handle"`
		Name   string `json:"name"`

		Deleted  filter.State `json:"deleted"`
		Archived filter.State `json:"archived"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Check func(*Role) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

// FindByHandle finds role by it's handle
func (set RoleSet) FindByHandle(handle string) *Role {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}
