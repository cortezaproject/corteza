package wfexec

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"testing"
	"time"
)

type (
	sesTestStep struct {
		StepIdentifier
		name string
		exec func(context.Context, *ExecRequest) (ExecResponse, error)
	}

	sesTestTemporal struct {
		StepIdentifier
		delay time.Duration
		until time.Time
	}
)

var (
	// used for testing to produce lower numbers that are easier to inspect and compare
	testID = atomic.NewUint64(0)
)

func (s *sesTestStep) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	if s.exec != nil {
		return s.exec(ctx, r)
	}

	var (
		args = &struct {
			Path    string
			Counter int64
		}{}
	)

	if err := r.Scope.Decode(args); err != nil {
		return nil, err
	}

	return expr.NewVars(map[string]interface{}{
		"counter": args.Counter + 1,
		"path":    args.Path + "/" + s.name,
		s.name:    "executed",
	})
}

func (s *sesTestTemporal) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	if s.until.IsZero() {
		s.until = now().Add(s.delay)
	}

	if now().Before(s.until) {
		return Delay(s.until), nil
	}

	return expr.NewVars(map[string]interface{}{
		"waitForMoment": "executed",
	})
}

func TestSession_TwoStepWorkflow(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
		wf  = NewGraph()
		ses = NewSession(ctx, wf)

		s1 = &sesTestStep{name: "s1"}
		s2 = &sesTestStep{name: "s2"}

		scope = &expr.Vars{}
	)

	scope.Set("two", 1)
	scope.Set("three", 1)

	wf.AddStep(s1, s2) // 1st execute s1 then s2
	req.NoError(ses.Exec(ctx, s1, scope))
	req.NoError(ses.Wait(ctx))
	req.NoError(ses.Error())
	req.NotNil(ses.Result())
	req.Equal("/s1/s2", expr.Must(expr.Select(ses.Result(), "path")).Get())
}

func TestSession_SplitAndMerge(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
		wf  = NewGraph()
		ses = NewSession(ctx, wf)

		start  = &sesTestStep{name: "start"}
		split1 = &sesTestStep{name: "split1"}
		split2 = &sesTestStep{name: "split2"}
		split3 = &sesTestStep{name: "split3"}

		end = JoinGateway(split1, split2, split3)
	)

	wf.AddStep(start, split1, split2, split3)
	wf.AddStep(split1, end)
	wf.AddStep(split2, end)
	wf.AddStep(split3, end)
	ses.Exec(ctx, start, nil)
	ses.Wait(ctx)
	req.True(ses.Idle())
	req.NoError(ses.Error())
	req.NotNil(ses.Result())
	// split3 only!
	req.Equal("/start/split3", expr.Must(expr.Select(ses.Result(), "path")).Get())
	req.Contains(ses.Result().Dict(), "split1")
	req.Contains(ses.Result().Dict(), "split2")
	req.Contains(ses.Result().Dict(), "split3")
}

func TestSession_Delays(t *testing.T) {
	t.SkipNow()
	var (
		// how fast we want to go (lower = faster)
		//
		unit  = time.Millisecond
		delay = unit * 3

		ctx = context.Background()
		req = require.New(t)
		wf  = NewGraph()
		ses = NewSession(ctx, wf,
			// for testing we need much shorter worker intervals
			SetWorkerInterval(unit),
		)

		start = &sesTestStep{name: "start"}

		waitForMoment = &sesTestTemporal{delay: delay}

		waitForInputStateId atomic.Uint64
		waitForInput        = &sesTestStep{name: "waitForInput", exec: func(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
			if !r.Input.Has("input") {
				waitForInputStateId.Store(r.StateID)
				return Prompt(0, "", nil), nil
			}

			out := &expr.Vars{}
			_ = out.Set("waitForInput", "executed")
			r.Input.Copy(out, "input")

			return out, nil
		}}
	)

	ctx, cancelFn := context.WithTimeout(ctx, time.Second*5)
	defer cancelFn()

	wf.AddStep(start, waitForMoment)
	wf.AddStep(waitForMoment, waitForInput)

	req.NoError(ses.Exec(ctx, start, nil))

	// wait-for-moment step needs to be executed before we can resume wait-for-input
	req.NoError(ses.Wait(ctx))
	time.Sleep(delay + unit)
	req.NotZero(waitForInputStateId.Load())

	// should not be completed yet...
	req.True(ses.Idle())
	req.True(ses.Suspended())

	// push in the input
	input := &expr.Vars{}
	input.Set("inout", "foo")
	_, err := ses.Resume(ctx, waitForInputStateId.Load(), input)
	req.NoError(err)

	req.False(ses.Suspended())
	req.NoError(ses.Wait(ctx))
	time.Sleep(2 * unit)

	// should not be completed yet...
	req.True(ses.Idle())
	req.NoError(ses.Error())
	req.NotNil(ses.Result())
	req.Contains(ses.Result().Dict(), "waitForMoment")
	req.Contains(ses.Result().Dict(), "waitForInput")
	req.Equal("foo", expr.Must(expr.Select(ses.Result(), "input")).Get())

}

func bmSessionSimpleStepSequence(c uint64, b *testing.B) {
	var (
		ctx = context.Background()
		g   = NewGraph()
		err error
	)

	for i := uint64(1); i <= c; i++ {
		s := &sesTestStep{name: "start"}
		s.SetID(i)
		g.AddStep(s)
		if i > 1 {
			g.AddParent(s, g.StepByID(i-1))
		}
	}

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		ses := NewSession(ctx, g)
		if err = ses.Exec(ctx, g.StepByID(1), nil); err != nil {
			b.Fatal(err.Error())
		}

		ses.Wait(ctx)
	}
	b.StopTimer()
}

func BenchmarkSessionSimple1StepSequence(b *testing.B)       { bmSessionSimpleStepSequence(1, b) }
func BenchmarkSessionSimple10StepSequence(b *testing.B)      { bmSessionSimpleStepSequence(10, b) }
func BenchmarkSessionSimple100StepSequence(b *testing.B)     { bmSessionSimpleStepSequence(100, b) }
func BenchmarkSessionSimple1000StepSequence(b *testing.B)    { bmSessionSimpleStepSequence(1000, b) }
func BenchmarkSessionSimple10000StepSequence(b *testing.B)   { bmSessionSimpleStepSequence(10000, b) }
func BenchmarkSessionSimple100000StepSequence(b *testing.B)  { bmSessionSimpleStepSequence(100000, b) }
func BenchmarkSessionSimple1000000StepSequence(b *testing.B) { bmSessionSimpleStepSequence(1000000, b) }
