package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	DalSensitivityLevelFilter struct {
		DalSensitivityLevelID []string `json:"sensitivityLevelID"`
		Handle                string   `json:"handle"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*DalSensitivityLevel) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Paging
		filter.Sorting
	}
)

func (ss DalSensitivityLevelSet) Len() int { return len(ss) }
func (ss DalSensitivityLevelSet) Less(i, j int) bool {
	return ss[i].Level < ss[j].Level
}
func (ss DalSensitivityLevelSet) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
