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

type SharedModule struct {
	ID                         uint64         `json:"moduleID,string"`
	TenantID                   uint64         `json:"tenantID,string,omitempty"`
	ProjectID                  uint64         `json:"projectID,string,omitempty"`
	Handle                     string         `json:"handle"`
	NodeID                     uint64         `json:"nodeID,string"`
	Name                       string         `json:"name"`
	ExternalFederationModuleID uint64         `json:"externalFederationModuleID,string"`
	Fields                     ModuleFieldSet `json:"fields"`
	CreatedAt                  time.Time      `json:"createdAt,omitempty"`
	UpdatedAt                  *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt                  *time.Time     `json:"deletedAt,omitempty"`
	CreatedBy                  uint64         `json:"createdBy,string"`
	UpdatedBy                  uint64         `json:"updatedBy,string,omitempty"`
	DeletedBy                  uint64         `json:"deletedBy,string,omitempty"`
}

func (r SharedModule) Clone() *SharedModule {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}
