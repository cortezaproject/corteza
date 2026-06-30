package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"encoding/json"
)

type DmlImportRun struct {
	ID           uint64            `json:"runID,string"`
	ConnectionID uint64            `json:"connectionID,string"`
	MappingID    uint64            `json:"mappingID,string"`
	Method       DmlImportMethod   `json:"method"`
	Status       string            `json:"status"`
	Processed    uint64            `json:"processed"`
	Failed       uint64            `json:"failed"`
	Error        string            `json:"error,omitempty"`
	Cursor       map[string]string `json:"cursor,omitempty"`
}

func (r DmlImportRun) Clone() *DmlImportRun {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}
