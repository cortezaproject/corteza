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

type ResourceTranslation struct {
	ID        uint64     `json:"translationID,string"`
	TenantID  uint64     `json:"tenantID,string,omitempty"`
	ProjectID uint64     `json:"projectID,string,omitempty"`
	Lang      Lang       `json:"lang"`
	Resource  string     `json:"resource"`
	K         string     `json:"key"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	OwnedBy   uint64     `json:"ownedBy,string"`
	CreatedBy uint64     `json:"createdBy,string"`
	UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
	DeletedBy uint64     `json:"deletedBy,string,omitempty"`
}

func (r ResourceTranslation) Clone() *ResourceTranslation {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}
