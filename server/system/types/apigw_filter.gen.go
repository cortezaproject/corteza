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

type ApigwFilter struct {
	ID        uint64            `json:"filterID,string"`
	TenantID  uint64            `json:"tenantID,string,omitempty"`
	ProjectID uint64            `json:"projectID,string,omitempty"`
	Route     uint64            `json:"routeID,string"`
	Weight    uint64            `json:"weight,string"`
	Kind      string            `json:"kind,omitempty"`
	Ref       string            `json:"ref,omitempty"`
	Enabled   bool              `json:"enabled,omitempty"`
	Params    ApigwFilterParams `json:"params"`
	CreatedAt time.Time         `json:"createdAt,omitempty"`
	UpdatedAt *time.Time        `json:"updatedAt,omitempty"`
	DeletedAt *time.Time        `json:"deletedAt,omitempty"`
	CreatedBy uint64            `json:"createdBy,string"`
	UpdatedBy uint64            `json:"updatedBy,string,omitempty"`
	DeletedBy uint64            `json:"deletedBy,string,omitempty"`
}

func (m *ApigwFilterParams) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ApigwFilterParams) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseApigwFilterParams(ss []string) (p ApigwFilterParams, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
