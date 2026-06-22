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

type ProjectGroup struct {
	ID        uint64           `json:"projectGroupID,string"`
	TenantID  uint64           `json:"tenantID,string,omitempty"`
	ProjectID uint64           `json:"projectID,string"`
	Handle    string           `json:"handle"`
	Meta      ProjectGroupMeta `json:"meta"`
	CreatedAt time.Time        `json:"createdAt,omitempty"`
	UpdatedAt *time.Time       `json:"updatedAt,omitempty"`
	DeletedAt *time.Time       `json:"deletedAt,omitempty"`
}

func (r ProjectGroup) Clone() *ProjectGroup {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *ProjectGroupMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ProjectGroupMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseProjectGroupMeta(ss []string) (p ProjectGroupMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
