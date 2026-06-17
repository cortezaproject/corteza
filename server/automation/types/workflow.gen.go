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

type Workflow struct {
	ID           uint64                           `json:"workflowID,string"`
	TenantID     uint64                           `json:"tenantID,string,omitempty"`
	ProjectID    uint64                           `json:"projectID,string,omitempty"`
	Handle       string                           `json:"handle"`
	Meta         *WorkflowMeta                    `json:"meta,omitempty"`
	Enabled      bool                             `json:"enabled"`
	Trace        bool                             `json:"trace"`
	KeepSessions int                              `json:"keepSessions"`
	Scope        *expr.Vars                       `json:"scope"`
	Steps        WorkflowStepSet                  `json:"steps"`
	Paths        WorkflowPathSet                  `json:"paths"`
	Issues       WorkflowIssueSet                 `json:"issues,omitempty"`
	RunAs        uint64                           `json:"runAs,string"`
	OwnedBy      uint64                           `json:"ownedBy,string"`
	CreatedAt    time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt    *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy    uint64                           `json:"createdBy,string"`
	UpdatedBy    uint64                           `json:"updatedBy,string,omitempty"`
	DeletedBy    uint64                           `json:"deletedBy,string,omitempty"`
	Labels       map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (m *WorkflowMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m WorkflowMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseWorkflowMeta(ss []string) (p *WorkflowMeta, err error) {
	p = &WorkflowMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func (m *WorkflowStepSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m WorkflowStepSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseWorkflowStepSet(ss []string) (p WorkflowStepSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *WorkflowPathSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m WorkflowPathSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseWorkflowPathSet(ss []string) (p WorkflowPathSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *WorkflowIssueSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m WorkflowIssueSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseWorkflowIssueSet(ss []string) (p WorkflowIssueSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
