package service

import (
	"testing"

	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"github.com/stretchr/testify/require"
)

func TestSession_Start(t *testing.T) {
	var (
		req = require.New(t)
		ses = &session{}
		g   = wfexec.NewGraph()
		err error
	)

	_, err = ses.Start(g, nil, types.SessionStartParams{})
	req.EqualError(err, "could not find starting step")

	g.AddStep(wfexec.NewGenericStep(nil))
	_, err = ses.Start(g, nil, types.SessionStartParams{StepID: 4321})
	req.EqualError(err, "trigger staring step references non-existing step")

	// Adding another orphaned step and starting session w/o explicitly specifying the starting step
	g.AddStep(wfexec.NewGenericStep(nil))
	_, err = ses.Start(g, nil, types.SessionStartParams{})
	req.EqualError(err, "cannot start workflow session multiple starting steps found")

	// add a generic step with a known ID so we can use it as a starting point
	s := wfexec.NewGenericStep(nil)
	s.SetID(42)
	g.AddStep(s)
	// add parents to the 42 step
	g.AddStep(wfexec.NewGenericStep(nil), s)
	_, err = ses.Start(g, nil, types.SessionStartParams{StepID: 42})
	req.EqualError(err, "cannot start workflow on a step with parents")

}
