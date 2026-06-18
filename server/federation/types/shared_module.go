package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	SharedModuleFilter struct {
		NodeID                     uint64 `json:"nodeID,string"`
		ExternalFederationModuleID uint64 `json:"externalFederationModuleID,string"`

		Handle string `json:"handle"`
		Name   string `json:"name"`
		Query  string `json:"query"`

		Check func(*SharedModule) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
