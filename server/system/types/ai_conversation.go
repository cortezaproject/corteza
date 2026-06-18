package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	AiConversationMessage struct {
		Role        string                     `json:"role"`
		Content     string                     `json:"content"`
		Operator    string                     `json:"operator,omitempty"`
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
