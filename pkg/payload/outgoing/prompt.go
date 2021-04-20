package outgoing

import (
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"time"
)

type (
	Prompt struct {
		Ref       string     `json:"ref"`
		SessionID uint64     `json:"sessionID,string"`
		CreatedAt time.Time  `json:"createdAt"`
		StateID   uint64     `json:"stateID,string"`
		Payload   *expr.Vars `json:"payload"`
	}

	Prompts []*Prompt
)
