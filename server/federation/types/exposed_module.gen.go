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

type ExposedModule struct {
	ID                 uint64         `json:"moduleID,string"`
	TenantID           uint64         `json:"tenantID,string,omitempty"`
	ProjectID          uint64         `json:"projectID,string,omitempty"`
	Handle             string         `json:"handle"`
	Name               string         `json:"name"`
	NodeID             uint64         `json:"nodeID,string"`
	ComposeModuleID    uint64         `json:"composeModuleID,string"`
	ComposeNamespaceID uint64         `json:"composeNamespaceID,string"`
	Fields             ModuleFieldSet `json:"fields"`
	CreatedAt          time.Time      `json:"createdAt,omitempty"`
	UpdatedAt          *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt          *time.Time     `json:"deletedAt,omitempty"`
	CreatedBy          uint64         `json:"createdBy,string"`
	UpdatedBy          uint64         `json:"updatedBy,string,omitempty"`
	DeletedBy          uint64         `json:"deletedBy,string,omitempty"`
}

func (m *ModuleFieldSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ModuleFieldSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseModuleFieldSet(ss []string) (p ModuleFieldSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
