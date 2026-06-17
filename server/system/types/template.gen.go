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

type Template struct {
	ID         uint64                           `json:"templateID,string"`
	TenantID   uint64                           `json:"tenantID,string,omitempty"`
	ProjectID  uint64                           `json:"projectID,string,omitempty"`
	OwnerID    uint64                           `json:"ownerID,string"`
	Handle     string                           `json:"handle"`
	Language   string                           `json:"language"`
	Type       DocumentType                     `json:"type"`
	Partial    bool                             `json:"partial"`
	Meta       TemplateMeta                     `json:"meta"`
	Template   string                           `json:"template"`
	CreatedAt  time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
	LastUsedAt *time.Time                       `json:"lastUsedAt,omitempty"`
	Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (m *TemplateMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m TemplateMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseTemplateMeta(ss []string) (p TemplateMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
