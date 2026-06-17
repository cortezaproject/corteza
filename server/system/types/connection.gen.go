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

type Connection struct {
	ID             uint64                           `json:"connectionID,string"`
	Handle         string                           `json:"handle"`
	Revision       int                              `json:"revision"`
	Status         string                           `json:"status"`
	Source         string                           `json:"source,omitempty"`
	Meta           ConnectionMeta                   `json:"meta"`
	Service        ConnectionService                `json:"service"`
	Resources      ConnectionResources              `json:"resources"`
	Operations     ConnectionOperations             `json:"operations"`
	DerivedParams  []ConnectionDerivedParam         `json:"derivedParams,omitempty"`
	CatalogID      string                           `json:"catalogID,omitempty"`
	InstalledCount int                              `json:"installedCount,omitempty"`
	CreatedAt      time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt      *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy      uint64                           `json:"createdBy,string"`
	UpdatedBy      uint64                           `json:"updatedBy,string,omitempty"`
	DeletedBy      uint64                           `json:"deletedBy,string,omitempty"`
	Labels         map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (m *ConnectionMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ConnectionMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseConnectionMeta(ss []string) (p ConnectionMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ConnectionService) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ConnectionService) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseConnectionService(ss []string) (p ConnectionService, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ConnectionResources) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ConnectionResources) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseConnectionResources(ss []string) (p ConnectionResources, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ConnectionOperations) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ConnectionOperations) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseConnectionOperations(ss []string) (p ConnectionOperations, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
