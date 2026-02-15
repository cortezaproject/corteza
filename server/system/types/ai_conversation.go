package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/sql"
)

type (
	AiConversation struct {
		ID         uint64                 `json:"aiConversationID,string"`
		AgentID    uint64                 `json:"agentID,string"`
		Messages   AiConversationMessages `json:"messages"`
		TokenCount int                    `json:"tokenCount"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	AiConversationMessage struct {
		Role        string                     `json:"role"`
		Content     string                     `json:"content"`
		ToolCalls   []AiConversationToolCall   `json:"toolCalls,omitempty"`
		ToolResults []AiConversationToolResult `json:"toolResults,omitempty"`
	}

	AiConversationToolCall struct {
		CallID string `json:"callID"`
		Name   string `json:"name"`
		Data   string `json:"data"`
	}

	AiConversationToolResult struct {
		CallID string `json:"callID"`
		Data   string `json:"data"`
		Error  string `json:"error,omitempty"`
	}

	AiConversationMessages []AiConversationMessage

	AiConversationFilter struct {
		AiConversationID []uint64     `json:"aiConversationID"`
		AgentID          uint64       `json:"agentID,string"`
		Deleted          filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*AiConversation) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

func (m *AiConversationMessages) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AiConversationMessages) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAiConversationMessages(ss []string) (p AiConversationMessages, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
