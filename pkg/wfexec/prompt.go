package wfexec

import (
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"time"
)

type (
	prompted struct {
		// when !nul, assuming waiting for input we're waiting for input
		payload *expr.Vars

		// user to be prompted
		ownerId uint64

		// state to be resumed
		state *State

		// prompt reference; something client can use
		// for orientation, what kind of prompt is expected
		ref string
	}

	PendingPrompt struct {
		Ref       string     `json:"ref"`
		SessionID uint64     `json:"sessionID,string"`
		CreatedAt time.Time  `json:"created"`
		StateID   uint64     `json:"stateID,string"`
		Payload   *expr.Vars `json:"payload"`
	}
)

func Prompt(ownerId uint64, ref string, payload *expr.Vars) *prompted {
	return &prompted{payload: payload, ref: ref, ownerId: ownerId}
}

func (p *prompted) toPending() *PendingPrompt {
	return &PendingPrompt{
		Ref:       p.ref,
		CreatedAt: p.state.created,
		StateID:   p.state.stateId,
		Payload:   p.payload,
	}
}
