package wfexec

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	wfTestStep struct {
		stepIdentifier
		name string
	}
)

func (wfTestStep) Exec(context.Context, *ExecRequest) (ExecResponse, error) {
	return nil, nil
}

func TestWorkflow(t *testing.T) {
	var (
		req = require.New(t)
		wf  = NewGraph()

		s          = &wfTestStep{name: "s1"}
		c1, c2, c3 = &wfTestStep{name: "c1"}, &wfTestStep{name: "c2"}, &wfTestStep{name: "c3"}
	)

	wf.AddStep(s, c1, c2, c3)
	req.Equal(wf.Children(s), Steps{c1, c2, c3})
	req.Equal(wf.Parents(c1), Steps{s})
	req.Equal(wf.Parents(c2), Steps{s})
	req.Equal(wf.Parents(c3), Steps{s})
}
