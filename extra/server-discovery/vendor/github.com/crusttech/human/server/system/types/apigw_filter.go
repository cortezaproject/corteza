package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ApigwFilterParams map[string]interface{}

	ApigwFilterFilter struct {
		ApigwFilterID []string `json:"apigwFilterID"`
		RouteID       uint64   `json:"routeID,string"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		Kind string `json:"kind"`
		Ref  string `json:"ref"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*ApigwFilter) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
