package types

import (
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

var (
	NodeStatusPending       = "pending"
	NodeStatusPairRequested = "pair_requested"
	NodeStatusPaired        = "paired"
)

type (
	Node struct {
		ID     uint64 `json:"nodeID,string"`
		Name   string `json:"name"`
		Status string `json:"status"`

		// Base URL of the remote server
		BaseURL string `json:"baseURL"`

		Contact string `json:"contact"`

		// Node ID on the remote server that points back to us
		SharedNodeID uint64 `json:"sharedNodeID,string"`

		PairToken string `json:"-"`
		AuthToken string `json:"-"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

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
