package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/sql"
	"time"
)

type DalConnection struct {
	ID        uint64              `json:"connectionID,string"`
	Handle    string              `json:"handle"`
	Type      string              `json:"type"`
	Config    DalConnectionConfig `json:"config"`
	Meta      DalConnectionMeta   `json:"meta"`
	Issues    []dal.Issue         `json:"issues,omitempty"`
	Labels    map[string]string   `json:"labels,omitempty"`
	CreatedAt time.Time           `json:"createdAt,omitempty"`
	UpdatedAt *time.Time          `json:"updatedAt,omitempty"`
	DeletedAt *time.Time          `json:"deletedAt,omitempty"`
	CreatedBy uint64              `json:"createdBy,string"`
	UpdatedBy uint64              `json:"updatedBy,string,omitempty"`
	DeletedBy uint64              `json:"deletedBy,string,omitempty"`
}

func (m *DalConnectionConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m DalConnectionConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseDalConnectionConfig(ss []string) (p DalConnectionConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *DalConnectionMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m DalConnectionMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseDalConnectionMeta(ss []string) (p DalConnectionMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
