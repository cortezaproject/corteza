package payload

import (
	"github.com/cortezaproject/corteza-server/pkg/payload/outgoing"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

func Prompt(p *wfexec.PendingPrompt) *outgoing.Prompt {
	return &outgoing.Prompt{
		Ref:       p.Ref,
		SessionID: p.SessionID,
		CreatedAt: p.CreatedAt,
		StateID:   p.StateID,
		Payload:   p.Payload,
	}
}

func Prompts(prompts []*wfexec.PendingPrompt) *outgoing.Prompts {
	ps := make([]*outgoing.Prompt, len(prompts))
	for k, p := range prompts {
		ps[k] = Prompt(p)
	}
	retval := outgoing.Prompts(ps)
	return &retval
}
