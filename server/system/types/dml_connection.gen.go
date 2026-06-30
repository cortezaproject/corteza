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

type DmlConnection struct {
	ID        uint64              `json:"connectionID,string"`
	Handle    string              `json:"handle"`
	Label     string              `json:"label"`
	Params    DmlConnectionParams `json:"params"`
	CreatedAt time.Time           `json:"createdAt,omitempty"`
	UpdatedAt *time.Time          `json:"updatedAt,omitempty"`
	DeletedAt *time.Time          `json:"deletedAt,omitempty"`
}

func (r DmlConnection) Clone() *DmlConnection {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *DmlConnectionParams) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m DmlConnectionParams) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseDmlConnectionParams(ss []string) (p DmlConnectionParams, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
