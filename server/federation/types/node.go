package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

var (
	NodeStatusPending       = "pending"
	NodeStatusPairRequested = "pair_requested"
	NodeStatusPaired        = "paired"
	NodeStatusFailed        = "failed"
)

type (
	NodeFilter struct {
		Query  string `json:"name"`
		Status string `json:"status"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(node *Node) (bool, error) `json:"-"`

		Deleted filter.State `json:"deleted"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)
