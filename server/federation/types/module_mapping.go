package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ModuleMappingFilter struct {
		NodeID             uint64 `json:"nodeID"`
		ComposeModuleID    uint64 `json:"composeModuleID"`
		ComposeNamespaceID uint64 `json:"composeNamespaceID"`
		FederationModuleID uint64 `json:"federationModuleID"`
		Query              string `json:"query"`

		Check func(*ModuleMapping) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
