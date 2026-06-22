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

type NgAutomation struct {
	ID        uint64                           `json:"automationID,string"`
	TenantID  uint64                           `json:"tenantID,string,omitempty"`
	ProjectID uint64                           `json:"projectID,string,omitempty"`
	Handle    string                           `json:"handle"`
	Meta      *NgAutomationMeta                `json:"meta,omitempty"`
	Enabled   bool                             `json:"enabled"`
	Scope     *expr.Vars                       `json:"scope"`
	Triggers  NgAutomationTriggerSet           `json:"triggers"`
	Steps     NgAutomationStepSet              `json:"steps"`
	Paths     NgAutomationPathSet              `json:"paths"`
	Issues    NgAutomationIssueSet             `json:"issues,omitempty"`
	RunAs     uint64                           `json:"runAs,string"`
	OwnedBy   uint64                           `json:"ownedBy,string"`
	CreatedAt time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy uint64                           `json:"createdBy,string"`
	UpdatedBy uint64                           `json:"updatedBy,string,omitempty"`
	DeletedBy uint64                           `json:"deletedBy,string,omitempty"`
	Labels    map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r NgAutomation) Clone() *NgAutomation {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *NgAutomationMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NgAutomationMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNgAutomationMeta(ss []string) (p *NgAutomationMeta, err error) {
	p = &NgAutomationMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func (m *NgAutomationTriggerSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NgAutomationTriggerSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNgAutomationTriggerSet(ss []string) (p NgAutomationTriggerSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *NgAutomationStepSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NgAutomationStepSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNgAutomationStepSet(ss []string) (p NgAutomationStepSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *NgAutomationPathSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NgAutomationPathSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNgAutomationPathSet(ss []string) (p NgAutomationPathSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *NgAutomationIssueSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m NgAutomationIssueSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseNgAutomationIssueSet(ss []string) (p NgAutomationIssueSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
