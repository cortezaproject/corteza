package types

import (
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

var (
	NodeSyncTypeStructure = "sync_structure"
	NodeSyncTypeData      = "sync_data"
	NodeSyncStatusSuccess = "success"
	NodeSyncStatusError   = "error"
)

type (
	NodeSync struct {
		NodeID     uint64 `json:"nodeID,string"`
		ModuleID   uint64 `json:"moduleID,string"`
		SyncStatus string `json:"syncStatus"`
		SyncType   string `json:"syncType"`

		TimeOfAction time.Time `json:"timeOfAction"`
	}

	NodeSyncFilter struct {
		NodeID     uint64 `json:"nodeID"`
		ModuleID   uint64 `json:"moduleID"`
		SyncStatus string `json:"syncStatus"`
		SyncType   string `json:"syncType"`

		Query string `json:"name"`

		Check func(*NodeSync) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)
