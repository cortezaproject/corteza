package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"encoding/json"
	"time"
)

type NodeSync struct {
	NodeID       uint64    `json:"nodeID,string"`
	ModuleID     uint64    `json:"moduleID,string"`
	SyncType     string    `json:"syncType"`
	SyncStatus   string    `json:"syncStatus"`
	TimeOfAction time.Time `json:"timeOfAction"`
}

func (r NodeSync) Clone() *NodeSync {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}
