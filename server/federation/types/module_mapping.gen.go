package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/sql"
)

type ModuleMapping struct {
	TenantID           uint64                `json:"tenantID,string,omitempty"`
	ProjectID          uint64                `json:"projectID,string,omitempty"`
	NodeID             uint64                `json:"nodeID,string"`
	FederationModuleID uint64                `json:"federationModuleID,string"`
	ComposeModuleID    uint64                `json:"composeModuleID,string"`
	ComposeNamespaceID uint64                `json:"composeNamespaceID,string"`
	FieldMapping       ModuleFieldMappingSet `json:"fields"`
}

func (r ModuleMapping) Clone() *ModuleMapping {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *ModuleFieldMappingSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ModuleFieldMappingSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseModuleFieldMappingSet(ss []string) (p ModuleFieldMappingSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
