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
