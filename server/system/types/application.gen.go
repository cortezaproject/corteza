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

type Application struct {
	ID        uint64                           `json:"applicationID,string"`
	TenantID  uint64                           `json:"tenantID,string,omitempty"`
	ProjectID uint64                           `json:"projectID,string,omitempty"`
	Name      string                           `json:"name"`
	Enabled   bool                             `json:"enabled"`
	Weight    int                              `json:"weight"`
	Unify     *ApplicationUnify                `json:"unify,omitempty"`
	OwnerID   uint64                           `json:"ownerID"`
	CreatedAt time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt *time.Time                       `json:"deletedAt,omitempty"`
	Labels    map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	Flags     []string                         `json:"flags,omitempty"`
}

func (m *ApplicationUnify) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ApplicationUnify) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseApplicationUnify(ss []string) (p ApplicationUnify, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
