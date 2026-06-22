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

type Page struct {
	ID             uint64                           `json:"pageID,string"`
	TenantID       uint64                           `json:"tenantID,string,omitempty"`
	ProjectID      uint64                           `json:"projectID,string,omitempty"`
	Title          string                           `json:"title"`
	Handle         string                           `json:"handle"`
	SelfID         uint64                           `json:"selfID,string"`
	ModuleID       uint64                           `json:"moduleID,string"`
	NamespaceID    uint64                           `json:"namespaceID,string"`
	Meta           PageMeta                         `json:"meta"`
	Config         PageConfig                       `json:"config"`
	Blocks         PageBlocks                       `json:"blocks"`
	Children       PageSet                          `json:"children,omitempty"`
	Visible        bool                             `json:"visible"`
	Weight         int                              `json:"weight"`
	Description    string                           `json:"description"`
	CreatedAt      time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt      *time.Time                       `json:"deletedAt,omitempty"`
	CreatedByAgent uint64                           `json:"createdByAgent,string,omitempty"`
	Labels         map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r Page) Clone() *Page {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *PageMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m PageMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParsePageMeta(ss []string) (p PageMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *PageConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m PageConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func ParsePageConfig(ss []string) (p PageConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *PageBlocks) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m PageBlocks) Value() (driver.Value, error) { return json.Marshal(m) }

func ParsePageBlocks(ss []string) (p PageBlocks, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
