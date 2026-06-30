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

type DmlMapping struct {
	ID              uint64          `json:"mappingID,string"`
	ConnectionID    uint64          `json:"connectionID,string"`
	NamespaceHandle string          `json:"namespaceHandle"`
	SourceIdent     string          `json:"sourceIdent"`
	ModuleHandle    string          `json:"moduleHandle"`
	ModuleName      string          `json:"moduleName"`
	Skip            bool            `json:"skip"`
	Identifier      string          `json:"identifier"`
	Columns         DmlColumnMapSet `json:"columns"`
	CreatedAt       time.Time       `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time      `json:"updatedAt,omitempty"`
	DeletedAt       *time.Time      `json:"deletedAt,omitempty"`
}

func (r DmlMapping) Clone() *DmlMapping {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *DmlColumnMapSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m DmlColumnMapSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseDmlColumnMapSet(ss []string) (p DmlColumnMapSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
