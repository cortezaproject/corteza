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

type Queue struct {
	ID        uint64     `json:"queueID,string"`
	TenantID  uint64     `json:"tenantID,string,omitempty"`
	ProjectID uint64     `json:"projectID,string,omitempty"`
	Consumer  string     `json:"consumer"`
	Queue     string     `json:"queue"`
	Meta      QueueMeta  `json:"meta"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	CreatedBy uint64     `json:"createdBy,string"`
	UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
	DeletedBy uint64     `json:"deletedBy,string,omitempty"`
}

func (m *QueueMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m QueueMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseQueueMeta(ss []string) (p QueueMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
