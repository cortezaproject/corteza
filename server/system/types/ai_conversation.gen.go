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

type AiConversation struct {
	ID         uint64                 `json:"aiConversationID,string"`
	TenantID   uint64                 `json:"tenantID,string,omitempty"`
	ProjectID  uint64                 `json:"projectID,string,omitempty"`
	AgentID    uint64                 `json:"agentID,string"`
	Messages   AiConversationMessages `json:"messages"`
	TokenCount int                    `json:"tokenCount"`
	CreatedAt  time.Time              `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time             `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time             `json:"deletedAt,omitempty"`
	CreatedBy  uint64                 `json:"createdBy,string"`
	UpdatedBy  uint64                 `json:"updatedBy,string,omitempty"`
	DeletedBy  uint64                 `json:"deletedBy,string,omitempty"`
}

func (m *AiConversationMessages) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AiConversationMessages) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAiConversationMessages(ss []string) (p AiConversationMessages, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
