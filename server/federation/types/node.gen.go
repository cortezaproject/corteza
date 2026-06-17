package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"time"
)

type Node struct {
	ID           uint64     `json:"nodeID,string"`
	TenantID     uint64     `json:"tenantID,string,omitempty"`
	SharedNodeID uint64     `json:"sharedNodeID,string"`
	Name         string     `json:"name"`
	BaseURL      string     `json:"baseURL"`
	Status       string     `json:"status"`
	Contact      string     `json:"contact"`
	PairToken    string     `json:"-"`
	AuthToken    string     `json:"-"`
	CreatedAt    time.Time  `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
	CreatedBy    uint64     `json:"createdBy,string"`
	UpdatedBy    uint64     `json:"updatedBy,string,omitempty"`
	DeletedBy    uint64     `json:"deletedBy,string,omitempty"`
}
