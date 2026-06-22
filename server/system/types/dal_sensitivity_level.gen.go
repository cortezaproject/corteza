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

type DalSensitivityLevel struct {
	ID        uint64                  `json:"sensitivityLevelID,string"`
	TenantID  uint64                  `json:"tenantID,string,omitempty"`
	ProjectID uint64                  `json:"projectID,string,omitempty"`
	Handle    string                  `json:"handle"`
	Level     int                     `json:"level"`
	Meta      DalSensitivityLevelMeta `json:"meta"`
	Labels    map[string]string       `json:"labels,omitempty"`
	CreatedAt time.Time               `json:"createdAt,omitempty"`
	UpdatedAt *time.Time              `json:"updatedAt,omitempty"`
	DeletedAt *time.Time              `json:"deletedAt,omitempty"`
	CreatedBy uint64                  `json:"createdBy,string"`
	UpdatedBy uint64                  `json:"updatedBy,string,omitempty"`
	DeletedBy uint64                  `json:"deletedBy,string,omitempty"`
}

func (r DalSensitivityLevel) Clone() *DalSensitivityLevel {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *DalSensitivityLevelMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m DalSensitivityLevelMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseDalSensitivityLevelMeta(ss []string) (p DalSensitivityLevelMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
