package wfexec

import (
	"context"
	"time"
)

type (
	ExecResponse interface{}
	partial      struct{}

	message struct {
		ID   uint64
		Type string
		Body string
		msg  *Expression
	}
)

func NewPartial() *partial {
	return &partial{}
}

func DelayExecution(until time.Time) *suspended {
	return &suspended{resumeAt: &until}
}

func WaitForInput() *suspended {
	return &suspended{input: true}
}

func NewMessage(typ string, msg *Expression) *message {
	return &message{
		ID:   nextID(),
		Type: typ,
		msg:  msg,
	}
}

func (m *message) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	var err error
	if m.Body, err = m.msg.eval.EvalString(ctx, r.Scope); err != nil {
		return nil, err
	}

	return m, nil
}
