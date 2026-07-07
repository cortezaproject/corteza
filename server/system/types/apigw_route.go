package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ApigwRouteFilter struct {
		ApigwRouteID []string `json:"apigwRouteID"`
		Route        string   `json:"route"`
		Endpoint     string   `json:"endpoint"`
		Method       string   `json:"method"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*ApigwRoute) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
