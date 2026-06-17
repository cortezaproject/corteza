package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/expr"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/sql"
	"time"
)

type Trigger struct {
	ID           uint64                           `json:"triggerID,string"`
	TenantID     uint64                           `json:"tenantID,string,omitempty"`
	ProjectID    uint64                           `json:"projectID,string,omitempty"`
	WorkflowID   uint64                           `json:"workflowID,string"`
	StepID       uint64                           `json:"stepID,string"`
	Enabled      bool                             `json:"enabled"`
	Meta         *TriggerMeta                     `json:"meta,omitempty"`
	ResourceType string                           `json:"resourceType"`
	EventType    string                           `json:"eventType"`
	Constraints  TriggerConstraintSet             `json:"constraints"`
	Input        *expr.Vars                       `json:"input"`
	OwnedBy      uint64                           `json:"ownedBy,string"`
	CreatedAt    time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt    *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy    uint64                           `json:"createdBy,string"`
	UpdatedBy    uint64                           `json:"updatedBy,string,omitempty"`
	DeletedBy    uint64                           `json:"deletedBy,string,omitempty"`
	Labels       map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (m *TriggerMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m TriggerMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseTriggerMeta(ss []string) (p *TriggerMeta, err error) {
	p = &TriggerMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func (m *TriggerConstraintSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m TriggerConstraintSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseTriggerConstraintSet(ss []string) (p TriggerConstraintSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
