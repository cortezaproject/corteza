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

type PageLayout struct {
	ID             uint64                           `json:"pageLayoutID,string"`
	TenantID       uint64                           `json:"tenantID,string,omitempty"`
	ProjectID      uint64                           `json:"projectID,string,omitempty"`
	Handle         string                           `json:"handle"`
	Primary        bool                             `json:"primary"`
	PageID         uint64                           `json:"pageID,string"`
	ParentID       uint64                           `json:"parentID,string"`
	NamespaceID    uint64                           `json:"namespaceID,string"`
	Weight         int                              `json:"weight"`
	Meta           PageLayoutMeta                   `json:"meta,omitempty"`
	Config         PageLayoutConfig                 `json:"config"`
	Blocks         PageLayoutBlocks                 `json:"blocks,omitempty"`
	OwnedBy        uint64                           `json:"ownedBy,string"`
	CreatedAt      time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt      *time.Time                       `json:"deletedAt,omitempty"`
	CreatedByAgent uint64                           `json:"createdByAgent,string,omitempty"`
	Labels         map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r PageLayout) Clone() *PageLayout {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *PageLayoutMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m PageLayoutMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParsePageLayoutMeta(ss []string) (p PageLayoutMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *PageLayoutConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m PageLayoutConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func ParsePageLayoutConfig(ss []string) (p PageLayoutConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *PageLayoutBlocks) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m PageLayoutBlocks) Value() (driver.Value, error) { return json.Marshal(m) }

func ParsePageLayoutBlocks(ss []string) (p PageLayoutBlocks, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
