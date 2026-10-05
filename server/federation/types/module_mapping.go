package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ModuleMappingFilter struct {
		NodeID             uint64 `json:"nodeID,string"`
		ComposeModuleID    uint64 `json:"composeModuleID,string"`
		ComposeNamespaceID uint64 `json:"composeNamespaceID,string"`
		FederationModuleID uint64 `json:"federationModuleID,string"`
		Query              string `json:"query"`

		Check func(*ModuleMapping) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
