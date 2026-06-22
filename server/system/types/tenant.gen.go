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

type Tenant struct {
	ID          uint64                           `json:"tenantID,string"`
	Handle      string                           `json:"handle"`
	Status      TenantStatus                     `json:"status"`
	Config      TenantConfig                     `json:"config"`
	Meta        TenantMeta                       `json:"meta"`
	CreatedAt   time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time                       `json:"updatedAt,omitempty"`
	SuspendedAt *time.Time                       `json:"suspendedAt,omitempty"`
	DeletedAt   *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy   uint64                           `json:"createdBy,string"`
	UpdatedBy   uint64                           `json:"updatedBy,string,omitempty"`
	DeletedBy   uint64                           `json:"deletedBy,string,omitempty"`
	Labels      map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r Tenant) Clone() *Tenant {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *TenantConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m TenantConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseTenantConfig(ss []string) (p TenantConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *TenantMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m TenantMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseTenantMeta(ss []string) (p TenantMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
