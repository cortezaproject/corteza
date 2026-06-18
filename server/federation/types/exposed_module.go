package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
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
