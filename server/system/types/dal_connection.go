package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	DalConnectionFilter struct {
		DalConnectionID []string `json:"connectionID"`
		Handle          string   `json:"handle"`
		Type            string   `json:"type"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*DalConnection) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Paging
		filter.Sorting
	}
)

var (
	// Used to identify the primary DAL connection instead of an extra flag
	DalPrimaryConnectionResourceType = "corteza::system:primary-dal-connection"
	DalPrimaryConnectionHandle       = "primary-database"
)

func (c DalConnection) HasIssues() bool {
	return len(c.Issues) > 0
}
