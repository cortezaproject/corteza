package wfexec

import (
	"fmt"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/stretchr/testify/require"
)

// stateWithScope builds a state whose scope holds a record cache, the shape of
// a workflow that keeps looked-up records in a variable
func stateWithScope(records int) State {
	cache := make(map[string]interface{}, records)
	for i := 0; i < records; i++ {
		cache[fmt.Sprintf("rec%d", i)] = fmt.Sprintf("some reasonably sized record value %d", i)
	}

	vv, err := expr.NewVars(map[string]interface{}{"cache": cache})
	if err != nil {
		panic(err)
	}

	return State{
		created:   time.Now(),
		sessionId: 42,
		stateId:   43,
		action:    "step executed",
		scope:     vv,
		input:     expr.EmptyVars(),
		results:   expr.EmptyVars(),
	}
}

func TestState_MakeLightFrameSkipsVariables(t *testing.T) {
	var (
		req = require.New(t)
		s   = stateWithScope(50)

		light = s.MakeLightFrame()
		full  = s.MakeFrame()
	)

	// the light frame still describes the step
	req.Equal(s.sessionId, light.SessionID)
	req.Equal(s.stateId, light.StateID)
	req.Equal(s.action, light.Action)
	req.Equal(full.SessionID, light.SessionID)
	req.Equal(full.StateID, light.StateID)
	req.Equal(full.Action, light.Action)

	// ...but carries none of the variables
	req.Nil(light.Scope)
	req.Nil(light.Input)
	req.Nil(light.Results)

	// while the full frame snapshots them
	req.NotNil(full.Scope)
	req.Len(full.Scope.GetValue(), 1)
}

// BenchmarkState_MakeFrame contrasts the cost of snapshotting a step against
// only describing it
func BenchmarkState_MakeFrame(b *testing.B) {
	for _, records := range []int{100, 1000} {
		s := stateWithScope(records)

		b.Run(fmt.Sprintf("full/%d-records", records), func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_ = s.MakeFrame()
			}
		})

		b.Run(fmt.Sprintf("light/%d-records", records), func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				_ = s.MakeLightFrame()
			}
		})
	}
}
