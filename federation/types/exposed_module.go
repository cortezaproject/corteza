package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	ExposedModule struct {
		ID              uint64                 `json:"moduleID,string"`
		NodeID          uint64                 `json:"nodeID,string"`
		ComposeModuleID uint64                 `json:"composeModuleID,string"`
		Fields          ModuleFieldMappingList `json:"fields"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ExposedModuleFilter struct {
		NodeID          uint64 `json:"node"`
		ComposeModuleID uint64 `json:"composeModuleID"`
		Query           string `json:"query"`

		Check func(*ExposedModule) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
