package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/crusttech/human/server/pkg/revisions"
	"reflect"
	"time"
)

type (
	NodeSync struct {
		NodeID       uint64    `json:"nodeID,string"`
		ModuleID     uint64    `json:"moduleID,string"`
		SyncType     string    `json:"syncType"`
		SyncStatus   string    `json:"syncStatus"`
		TimeOfAction time.Time `json:"timeOfAction"`
	}
)

func (r NodeSync) Clone() *NodeSync {
	dup := r
	return &dup
}

func (r NodeSync) Diff(cmp *NodeSync) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &NodeSync{}
	}
	if r.NodeID != cmp.NodeID {
		out = append(out, &revisions.Change{Key: "nodeID", Old: []any{cmp.NodeID}, New: []any{r.NodeID}})
	}

	if r.ModuleID != cmp.ModuleID {
		out = append(out, &revisions.Change{Key: "moduleID", Old: []any{cmp.ModuleID}, New: []any{r.ModuleID}})
	}

	if r.SyncType != cmp.SyncType {
		out = append(out, &revisions.Change{Key: "syncType", Old: []any{cmp.SyncType}, New: []any{r.SyncType}})
	}

	if r.SyncStatus != cmp.SyncStatus {
		out = append(out, &revisions.Change{Key: "syncStatus", Old: []any{cmp.SyncStatus}, New: []any{r.SyncStatus}})
	}

	if !reflect.DeepEqual(r.TimeOfAction, cmp.TimeOfAction) {
		out = append(out, &revisions.Change{Key: "timeOfAction", Old: []any{cmp.TimeOfAction}, New: []any{r.TimeOfAction}})
	}

	return out
}
