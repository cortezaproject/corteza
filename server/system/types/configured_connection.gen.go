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

type ConfiguredConnection struct {
	ID           uint64                           `json:"configurationID,string"`
	TenantID     uint64                           `json:"tenantID,string,omitempty"`
	ProjectID    uint64                           `json:"projectID,string,omitempty"`
	ConnectionID uint64                           `json:"connectionID,string"`
	Name         string                           `json:"name"`
	Status       string                           `json:"status"`
	Connection   Connection                       `json:"connection"`
	Config       ConfiguredConnectionConfig       `json:"config"`
	CreatedAt    time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt    *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy    uint64                           `json:"createdBy,string"`
	UpdatedBy    uint64                           `json:"updatedBy,string,omitempty"`
	DeletedBy    uint64                           `json:"deletedBy,string,omitempty"`
	Labels       map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (m *Connection) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m Connection) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseConnection(ss []string) (p Connection, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ConfiguredConnectionConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ConfiguredConnectionConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseConfiguredConnectionConfig(ss []string) (p ConfiguredConnectionConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
