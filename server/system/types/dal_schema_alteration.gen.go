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
	"time"
)

type DalSchemaAlteration struct {
	ID           uint64                     `json:"alterationID,string"`
	BatchID      uint64                     `json:"batchID,string"`
	DependsOn    uint64                     `json:"dependsOn,string,omitempty"`
	Resource     string                     `json:"resource"`
	ResourceType string                     `json:"resourceType"`
	ConnectionID uint64                     `json:"connectionID,string"`
	Kind         string                     `json:"kind"`
	Params       *DalSchemaAlterationParams `json:"params"`
	Error        string                     `json:"error,omitempty"`
	CreatedAt    time.Time                  `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time                 `json:"updatedAt,omitempty"`
	DeletedAt    *time.Time                 `json:"deletedAt,omitempty"`
	CompletedAt  *time.Time                 `json:"completedAt,omitempty"`
	DismissedAt  *time.Time                 `json:"dismissedAt,omitempty"`
	CreatedBy    uint64                     `json:"createdBy,string"`
	UpdatedBy    uint64                     `json:"updatedBy,string,omitempty"`
	DeletedBy    uint64                     `json:"deletedBy,string,omitempty"`
	CompletedBy  uint64                     `json:"completedBy,string,omitempty"`
	DismissedBy  uint64                     `json:"dismissedBy,string,omitempty"`
}

func (r DalSchemaAlteration) Clone() *DalSchemaAlteration {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *DalSchemaAlterationParams) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m DalSchemaAlterationParams) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseDalSchemaAlterationParams(ss []string) (p DalSchemaAlterationParams, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
