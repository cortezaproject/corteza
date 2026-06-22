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

type DataPrivacyRequestComment struct {
	ID        uint64     `json:"commentID,string"`
	RequestID uint64     `json:"requestID,string"`
	Comment   string     `json:"comment"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	CreatedBy uint64     `json:"createdBy,string"`
	UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
	DeletedBy uint64     `json:"deletedBy,string,omitempty"`
}

func (r DataPrivacyRequestComment) Clone() *DataPrivacyRequestComment {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}
