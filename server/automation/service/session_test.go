package service

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/stretchr/testify/require"
)

func TestSession_Start(t *testing.T) {
	var (
		req = require.New(t)
		ses = &session{}
		g   = wfexec.NewGraph()
		ctx = context.Background()

		err error
	)

	_, _, err = ses.Start(ctx, g, types.SessionStartParams{Invoker: auth.Anonymous()})
	req.EqualError(err, "could not find starting step")

	g.AddStep(wfexec.NewGenericStep(nil))
	_, _, err = ses.Start(ctx, g, types.SessionStartParams{StepID: 4321, Invoker: auth.Anonymous()})
	req.EqualError(err, "trigger staring step references non-existing step")

	// Adding another orphaned step and starting session w/o explicitly specifying the starting step
	g.AddStep(wfexec.NewGenericStep(nil))
	_, _, err = ses.Start(ctx, g, types.SessionStartParams{Invoker: auth.Anonymous()})
	req.EqualError(err, "cannot start workflow session multiple starting steps found")

	// add a generic step with a known ID so we can use it as a starting point
	s := wfexec.NewGenericStep(nil)
	s.SetID(42)
	g.AddStep(s)
	// add parents to the 42 step
	g.AddStep(wfexec.NewGenericStep(nil), s)
	_, _, err = ses.Start(ctx, g, types.SessionStartParams{StepID: 42, Invoker: auth.Anonymous()})
	req.EqualError(err, "cannot start workflow on a step with parents")
}

func TestSessionStackTraces(t *testing.T) {
	stepID := uint64(0)
	getStepID := func() uint64 {
		stepID++
		return uint64(stepID)
	}

	iter1 := &wfexec.Frame{StepID: getStepID(), ParentID: 0, Action: "iterator initialized"}
	iter1stepExpr := &wfexec.Frame{StepID: getStepID(), ParentID: iter1.StepID, Action: ""}
	iter1stepCnt := &wfexec.Frame{StepID: getStepID(), ParentID: iter1stepExpr.StepID, Action: "loop continue"}

	iter2 := &wfexec.Frame{StepID: getStepID(), ParentID: iter1.StepID, Action: "iterator initialized"}
	iter2stepExpr := &wfexec.Frame{StepID: getStepID(), ParentID: iter2.StepID, Action: ""}
	iter2stepCnt := &wfexec.Frame{StepID: getStepID(), ParentID: iter2stepExpr.StepID, Action: "loop continue"}

	rando1 := &wfexec.Frame{StepID: getStepID(), ParentID: iter2.StepID, Action: ""}
	rando2 := &wfexec.Frame{StepID: getStepID(), ParentID: iter2.StepID, Action: ""}

	trace := []*wfexec.Frame{
		// first loop
		iter1,
		iter1stepExpr,
		iter1stepCnt,
		iter1,
		iter1stepExpr,
		iter1stepCnt,
		iter1,

		// second loop
		iter2,
		iter2stepExpr,
		iter2stepCnt,
		iter2,
		iter2stepExpr,
		iter2stepCnt,
		iter2,

		rando1,
		rando2,
	}

	ses := &types.Session{}

	for _, f := range trace {
		ses.AppendRuntimeStacktrace(f)
	}

	require.Len(t, ses.RuntimeStacktrace, 10)

	require.Equal(t, iter1.StepID, ses.RuntimeStacktrace[0].StepID)
	require.Equal(t, iter1stepExpr.StepID, ses.RuntimeStacktrace[1].StepID)
	require.Equal(t, iter1stepCnt.StepID, ses.RuntimeStacktrace[2].StepID)
	require.Equal(t, iter1.StepID, ses.RuntimeStacktrace[3].StepID)
	require.Equal(t, iter2.StepID, ses.RuntimeStacktrace[4].StepID)
	require.Equal(t, iter2stepExpr.StepID, ses.RuntimeStacktrace[5].StepID)
	require.Equal(t, iter2stepCnt.StepID, ses.RuntimeStacktrace[6].StepID)
	require.Equal(t, iter2.StepID, ses.RuntimeStacktrace[7].StepID)
	require.Equal(t, rando1.StepID, ses.RuntimeStacktrace[8].StepID)
	require.Equal(t, rando2.StepID, ses.RuntimeStacktrace[9].StepID)
}

// Before:
// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/automation/service
// BenchmarkSessionStackTraces_1000-12                42758             28594 ns/op           87288 B/op         15 allocs/op
// BenchmarkSessionStackTraces_10000-12                3430            348308 ns/op         1160442 B/op         23 allocs/op
// BenchmarkSessionStackTraces_100000-12                229           5299188 ns/op        13317368 B/op         33 allocs/op
// BenchmarkSessionStackTraces_1000000-12                19          52896750 ns/op        128431352 B/op        43 allocs/op
// BenchmarkSessionStackTraces_10000000-12                3         436311292 ns/op        1202337018 B/op       53 allocs/op

// After:
// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/automation/service
// BenchmarkSessionStackTraces_1000-12                40142             29860 ns/op             120 B/op          4 allocs/op
// BenchmarkSessionStackTraces_10000-12                4387            264930 ns/op             120 B/op          4 allocs/op
// BenchmarkSessionStackTraces_100000-12                447           2735284 ns/op             120 B/op          4 allocs/op
// BenchmarkSessionStackTraces_1000000-12                42          26596636 ns/op             120 B/op          4 allocs/op
// BenchmarkSessionStackTraces_10000000-12                4         265553667 ns/op             120 B/op          4 allocs/op
func benchmarkSessionStackTraces(b *testing.B, iters int) {
	stepID := uint64(0)
	getStepID := func() uint64 {
		stepID++
		return uint64(stepID)
	}

	iter := &wfexec.Frame{
		StepID: getStepID(),
		Action: "iterator initialized",
	}

	loop := []*wfexec.Frame{iter}
	loop = append(loop, &wfexec.Frame{
		StepID:   getStepID(),
		ParentID: iter.ParentID,
		Action:   "",
	})
	loop = append(loop, &wfexec.Frame{
		StepID:   getStepID(),
		ParentID: loop[len(loop)-1].StepID,
		Action:   "loop continue",
	})

	frames := []*wfexec.Frame{}
	for x := 0; x < iters; x++ {
		frames = append(frames, loop...)
	}

	frames = append(frames, iter)
	frames = append(frames, &wfexec.Frame{
		StepID:   getStepID(),
		ParentID: iter.ParentID,
		Action:   "termination",
	})

	b.StopTimer()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ses := &types.Session{}
		b.StartTimer()
		for _, f := range frames {
			ses.AppendRuntimeStacktrace(f)
		}
		b.StopTimer()
	}
}

func BenchmarkSessionStackTraces_1000(b *testing.B) {
	benchmarkSessionStackTraces(b, 1000)
}

func BenchmarkSessionStackTraces_10000(b *testing.B) {
	benchmarkSessionStackTraces(b, 10000)
}

func BenchmarkSessionStackTraces_100000(b *testing.B) {
	benchmarkSessionStackTraces(b, 100000)
}

func BenchmarkSessionStackTraces_1000000(b *testing.B) {
	benchmarkSessionStackTraces(b, 1000000)
}

func BenchmarkSessionStackTraces_10000000(b *testing.B) {
	benchmarkSessionStackTraces(b, 10000000)
}
