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
	"strconv"
	"time"
)

type Agent struct {
	ID         uint64                           `json:"agentID,string"`
	TenantID   uint64                           `json:"tenantID,string,omitempty"`
	ProjectID  uint64                           `json:"projectID,string,omitempty"`
	Handle     string                           `json:"handle"`
	Status     string                           `json:"status"`
	Revision   int                              `json:"revision"`
	Meta       AgentMeta                        `json:"meta"`
	Behavior   AgentBehavior                    `json:"behavior"`
	Execution  AgentExecution                   `json:"execution"`
	Access     AgentAccess                      `json:"access"`
	Invocation AgentInvocation                  `json:"invocation"`
	CreatedAt  time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy  uint64                           `json:"createdBy,string"`
	UpdatedBy  uint64                           `json:"updatedBy,string,omitempty"`
	DeletedBy  uint64                           `json:"deletedBy,string,omitempty"`
	Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r Agent) Clone() *Agent {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

// AgentAccessIDList is a []uint64 that serializes each element as a JSON string.
type AgentAccessIDList []uint64

func (ll AgentAccessIDList) MarshalJSON() ([]byte, error) {
	ss := make([]string, len(ll))
	for i, id := range ll {
		ss[i] = strconv.FormatUint(id, 10)
	}
	return json.Marshal(ss)
}

func (ll *AgentAccessIDList) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*ll = make(AgentAccessIDList, 0, len(raw))
	for _, r := range raw {
		s := string(r)
		if len(s) >= 2 && s[0] == '"' {
			s = s[1 : len(s)-1]
		}
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*ll = append(*ll, id)
	}
	return nil
}

func (m *AgentMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAgentMeta(ss []string) (p AgentMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *AgentBehavior) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentBehavior) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAgentBehavior(ss []string) (p AgentBehavior, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *AgentExecution) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentExecution) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAgentExecution(ss []string) (p AgentExecution, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *AgentAccess) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentAccess) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAgentAccess(ss []string) (p AgentAccess, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *AgentInvocation) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentInvocation) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAgentInvocation(ss []string) (p AgentInvocation, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
