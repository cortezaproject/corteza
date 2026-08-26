package types

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/stretchr/testify/require"
)

// makeScope builds a scope roughly the size of a workflow that keeps a set of
// records cached in a variable
func makeScope(records int) *expr.Vars {
	cache := make(map[string]interface{}, records)
	for i := 0; i < records; i++ {
		cache[fmt.Sprintf("rec%d", i)] = fmt.Sprintf("some reasonably sized record value %d", i)
	}

	vv, err := expr.NewVars(map[string]interface{}{"cache": cache})
	if err != nil {
		panic(err)
	}

	return vv
}

func makeFrame(stepID uint64, at time.Time, scope *expr.Vars) *wfexec.Frame {
	return &wfexec.Frame{StepID: stepID, CreatedAt: at, Scope: scope}
}

// countRetained reports how many frames the underlying array still points at,
// including the ones sliced off the end
func countRetained(s Stacktrace) (n int) {
	full := s[:cap(s)]
	for _, f := range full {
		if f != nil {
			n++
		}
	}

	return
}

func TestSession_runtimeStacktraceIsCapped(t *testing.T) {
	var (
		req = require.New(t)
		ses = &Session{}
		t0  = time.Now()
	)

	for i := 0; i < maxRuntimeStacktraceFrames*10; i++ {
		ses.AppendRuntimeStacktrace(makeFrame(uint64(i), t0.Add(time.Duration(i)*time.Millisecond), nil))
	}

	req.Len(ses.RuntimeStacktrace, maxRuntimeStacktraceFrames)

	// evicted frames must be unreachable, otherwise the underlying array keeps
	// every step's scope copy alive
	req.Equal(maxRuntimeStacktraceFrames, countRetained(ses.RuntimeStacktrace))

	// the cap keeps the most recent frames
	req.Equal(uint64(maxRuntimeStacktraceFrames*10-1), ses.RuntimeStacktrace[maxRuntimeStacktraceFrames-1].StepID)
}

func TestSession_elapsedTimeSurvivesEviction(t *testing.T) {
	var (
		req = require.New(t)
		ses = &Session{}
		t0  = time.Now()
	)

	for i := 0; i < maxRuntimeStacktraceFrames*10; i++ {
		ses.AppendRuntimeStacktrace(makeFrame(uint64(i), t0.Add(time.Duration(i)*time.Millisecond), nil))
	}

	// elapsed time is measured from the first recorded step, not from the
	// oldest frame that happens to still be around
	last := ses.RuntimeStacktrace[len(ses.RuntimeStacktrace)-1]
	req.Equal(uint(maxRuntimeStacktraceFrames*10-1), last.ElapsedTime)
}

func TestSession_tracedSessionKeepsEveryFrame(t *testing.T) {
	var (
		req = require.New(t)
		t0  = time.Now()

		// Apply() sets this when the session is started with Trace: true
		ses = &Session{Stacktrace: Stacktrace{}}
	)

	for i := 0; i < maxRuntimeStacktraceFrames*3; i++ {
		ses.AppendRuntimeStacktrace(makeFrame(uint64(i), t0, nil))
	}

	req.Len(ses.RuntimeStacktrace, maxRuntimeStacktraceFrames*3)
}

func TestSession_fullStacktraceKeepsEveryFrame(t *testing.T) {
	var (
		req = require.New(t)
		t0  = time.Now()
		ses = &Session{runtimeOpts: runtimeOptions{fullStacktrace: true}}
	)

	for i := 0; i < maxRuntimeStacktraceFrames*3; i++ {
		ses.AppendRuntimeStacktrace(makeFrame(uint64(i), t0, nil))
	}

	req.Len(ses.RuntimeStacktrace, maxRuntimeStacktraceFrames*3)
}

func TestSession_disabledStacktraceKeepsNothing(t *testing.T) {
	var (
		req = require.New(t)
		t0  = time.Now()
		ses = &Session{runtimeOpts: runtimeOptions{disableStacktrace: true}}
	)

	for i := 0; i < maxRuntimeStacktraceFrames*3; i++ {
		ses.AppendRuntimeStacktrace(makeFrame(uint64(i), t0, nil))
	}

	req.Len(ses.RuntimeStacktrace, 0)
}

// A loop over the same two steps must not accumulate frames, and the frames it
// slices off must be released
func TestSession_loopTruncationReleasesFrames(t *testing.T) {
	var (
		req = require.New(t)
		t0  = time.Now()
		ses = &Session{}
	)

	ses.AppendRuntimeStacktrace(&wfexec.Frame{StepID: 10, CreatedAt: t0, Action: "iterator initialized"})

	for i := 0; i < 1000; i++ {
		ses.AppendRuntimeStacktrace(makeFrame(10, t0, nil))
		ses.AppendRuntimeStacktrace(makeFrame(11, t0, nil))
	}

	req.LessOrEqual(len(ses.RuntimeStacktrace), maxRuntimeStacktraceFrames)
	req.Equal(len(ses.RuntimeStacktrace), countRetained(ses.RuntimeStacktrace))
}

// BenchmarkSession_runtimeStacktrace reports how much heap a session still
// holds on to after stepping through a workflow that keeps a record cache in
// scope; every frame snapshots that scope, so this is what grows the server
func BenchmarkSession_runtimeStacktrace(b *testing.B) {
	const steps = 2000

	heapAlloc := func() uint64 {
		var ms runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&ms)
		return ms.HeapAlloc
	}

	for _, tc := range []struct {
		name string
		ses  func() *Session
	}{
		{"default", func() *Session { return &Session{} }},
		{"traced", func() *Session { return &Session{Stacktrace: Stacktrace{}} }},
		{"full", func() *Session {
			return &Session{runtimeOpts: runtimeOptions{fullStacktrace: true}}
		}},
		{"disabled", func() *Session {
			return &Session{runtimeOpts: runtimeOptions{disableStacktrace: true}}
		}},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			var retained uint64

			for n := 0; n < b.N; n++ {
				var (
					ses  = tc.ses()
					t0   = time.Now()
					base = heapAlloc()
				)

				for i := 0; i < steps; i++ {
					// every step snapshots the whole scope, same as MakeFrame does
					ses.AppendRuntimeStacktrace(makeFrame(uint64(i), t0, makeScope(200)))
				}

				retained = heapAlloc() - base
				runtime.KeepAlive(ses)
			}

			b.ReportMetric(float64(retained)/(1<<20), "retained-MB")
		})
	}
}
