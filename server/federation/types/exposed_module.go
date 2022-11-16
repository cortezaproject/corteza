package types

import (
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	ExposedModule struct {
		ID                 uint64         `json:"moduleID,string"`
		NodeID             uint64         `json:"nodeID,string"`
		ComposeModuleID    uint64         `json:"composeModuleID,string"`
		ComposeNamespaceID uint64         `json:"composeNamespaceID,string"`
		Handle             string         `json:"handle"`
		Name               string         `json:"name"`
		Fields             ModuleFieldSet `json:"fields"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	ExposedModuleFilter struct {
		NodeID             uint64 `json:"nodeID,string"`
		ComposeModuleID    uint64 `json:"composeModuleID,string"`
		ComposeNamespaceID uint64 `json:"composeNamespaceID,string"`

		LastSync uint64 `json:"lastSync"`
		Handle   string `json:"handle"`
		Name     string `json:"name"`
		Query    string `json:"query"`

		Check func(*ExposedModule) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
