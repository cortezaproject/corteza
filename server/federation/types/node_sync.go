package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

var (
	NodeSyncTypeStructure = "sync_structure"
	NodeSyncTypeData      = "sync_data"
	NodeSyncStatusSuccess = "success"
	NodeSyncStatusError   = "error"
)

type (
	NodeSyncFilter struct {
		NodeID     uint64 `json:"nodeID"`
		RelNodeID  uint64 `json:"relNodeID"`
		ModuleID   uint64 `json:"moduleID"`
		SyncStatus string `json:"syncStatus"`
		SyncType   string `json:"syncType"`

		Query string `json:"name"`

		Check func(*NodeSync) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
