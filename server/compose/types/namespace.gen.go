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

type Namespace struct {
	ID             uint64                           `json:"namespaceID,string"`
	TenantID       uint64                           `json:"tenantID,string,omitempty"`
	ProjectID      uint64                           `json:"projectID,string,omitempty"`
	Slug           string                           `json:"slug"`
	Enabled        bool                             `json:"enabled"`
	Meta           NamespaceMeta                    `json:"meta"`
	Name           string                           `json:"name"`
	CreatedAt      time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt      *time.Time                       `json:"deletedAt,omitempty"`
	CreatedByAgent uint64                           `json:"createdByAgent,string,omitempty"`
	Labels         map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (m *NamespaceMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NamespaceMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNamespaceMeta(ss []string) (p NamespaceMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
