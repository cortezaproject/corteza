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

type Role struct {
	ID         uint64                           `json:"roleID,string"`
	TenantID   uint64                           `json:"tenantID,string,omitempty"`
	ProjectID  uint64                           `json:"projectID,string,omitempty"`
	Name       string                           `json:"name"`
	Handle     string                           `json:"handle"`
	Meta       *RoleMeta                        `json:"meta"`
	ArchivedAt *time.Time                       `json:"archivedAt,omitempty"`
	CreatedAt  time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
	Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r Role) Clone() *Role {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *RoleMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m RoleMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseRoleMeta(ss []string) (p *RoleMeta, err error) {
	p = &RoleMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}
