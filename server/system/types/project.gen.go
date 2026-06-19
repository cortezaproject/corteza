package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/sql"
	"time"
)

type Project struct {
	ID         uint64                           `json:"projectID,string"`
	TenantID   uint64                           `json:"tenantID,string,omitempty"`
	Handle     string                           `json:"handle"`
	Status     ProjectStatus                    `json:"status"`
	Config     ProjectConfig                    `json:"config"`
	Meta       ProjectMeta                      `json:"meta"`
	Governance ProjectGovernance                `json:"governance"`
	CreatedAt  time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy  uint64                           `json:"createdBy,string"`
	UpdatedBy  uint64                           `json:"updatedBy,string,omitempty"`
	DeletedBy  uint64                           `json:"deletedBy,string,omitempty"`
	Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (m *ProjectConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ProjectConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseProjectConfig(ss []string) (p ProjectConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ProjectMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ProjectMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseProjectMeta(ss []string) (p ProjectMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ProjectGovernance) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ProjectGovernance) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseProjectGovernance(ss []string) (p ProjectGovernance, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
