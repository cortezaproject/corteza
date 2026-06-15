package wfexec

import (
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"time"
)

type (
	ExecResponse interface{}
	partial      struct{}

	errHandler struct {
		handler Step
		results *expr.Vars
	}

	termination struct{}

	delayed struct {
		// when not nil, assuming delayed
		resumeAt time.Time

		// state to be resumed
		state *State
	}

	// when session is resumed from a delay we'll replace
	// delay step on state with the a generic step that will return resumed{}
	resumed struct{}
)

func Delay(until time.Time) *delayed {
	return &delayed{resumeAt: until}
}

func Resume() *resumed {
	return &resumed{}
}

func ErrorHandler(h Step, results *expr.Vars) *errHandler {
	return &errHandler{handler: h, results: results}
}

func Termination() *termination {
	return &termination{}
}

type (
	loopBreak    struct{}
	loopContinue struct{}
)

func LoopBreak() *loopBreak       { return &loopBreak{} }
func LoopContinue() *loopContinue { return &loopContinue{} }

// responseWithWarnings wraps a *expr.Vars scope together with a set of
// non-fatal warnings produced while computing it. The session unwraps
// this and attaches the warnings to the current State so they land in
// the emitted Frame's Warnings slice alongside the scope in Frame.Scope.
//
// Used by the join gateway to surface variable conflicts across
// parallel branches without changing the normal *expr.Vars contract.
type responseWithWarnings struct {
	scope    *expr.Vars
	warnings []string
}

// ResponseWithWarnings constructs a wrapper response. Callers pass the
// scope they would have returned normally plus a list of human-readable
// warning messages that the session should propagate into the frame.
func ResponseWithWarnings(scope *expr.Vars, warnings []string) *responseWithWarnings {
	return &responseWithWarnings{scope: scope, warnings: warnings}
}
