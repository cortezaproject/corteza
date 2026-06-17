package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"time"
)

type NodeSync struct {
	NodeID       uint64    `json:"nodeID,string"`
	ModuleID     uint64    `json:"moduleID,string"`
	SyncType     string    `json:"syncType"`
	SyncStatus   string    `json:"syncStatus"`
	TimeOfAction time.Time `json:"timeOfAction"`
}
