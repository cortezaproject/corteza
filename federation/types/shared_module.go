package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	SharedModule struct {
		ID                         uint64                 `json:"moduleID,string"`
		NodeID                     uint64                 `json:"nodeID,string"`
		Handle                     string                 `json:"handle"`
		Name                       string                 `json:"name"`
		ExternalFederationModuleID uint64                 `json:"externalFederationModuleID,string"`
		Fields                     ModuleFieldMappingList `json:"fields"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	SharedModuleFilter struct {
		NodeID uint64 `json:"node"`
		Query  string `json:"query"`

		Handle string `json:"handle"`
		Name   string `json:"name"`

		Check func(*SharedModule) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
