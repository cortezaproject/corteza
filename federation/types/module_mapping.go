package types

import (
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	ModuleMapping struct {
		FederationModuleID uint64                 `json:"federationModuleID,string"`
		ComposeModuleID    uint64                 `json:"composeModuleID,string"`
		FieldMapping       ModuleFieldMappingList `json:"fields"`
	}

	ModuleMappingFilter struct {
		ComposeModuleID    uint64 `json:"composeModuleID"`
		FederationModuleID uint64 `json:"federationModuleID"`
		Query              string `json:"query"`

		Check func(*ModuleMapping) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
