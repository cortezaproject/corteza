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

type UserGroup struct {
	ID         uint64                           `json:"userGroupID,string"`
	TenantID   uint64                           `json:"tenantID,string,omitempty"`
	ProjectID  uint64                           `json:"projectID,string,omitempty"`
	Handle     string                           `json:"handle"`
	Meta       *UserGroupMeta                   `json:"meta"`
	Config     *UserGroupConfig                 `json:"config"`
	IsRoot     bool                             `json:"isRoot"`
	ArchivedAt *time.Time                       `json:"archivedAt,omitempty"`
	CreatedAt  time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
	Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r UserGroup) Clone() *UserGroup {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *UserGroupMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m UserGroupMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseUserGroupMeta(ss []string) (p *UserGroupMeta, err error) {
	p = &UserGroupMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func (m *UserGroupConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m UserGroupConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseUserGroupConfig(ss []string) (p *UserGroupConfig, err error) {
	p = &UserGroupConfig{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}
