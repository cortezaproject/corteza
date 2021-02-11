package wfexec

import (
	"time"
)

type (
	ExecResponse interface{}
	partial      struct{}

	errHandler struct {
		handler Step
	}

	termination struct{}
)

func DelayExecution(until time.Time) *suspended {
	return &suspended{resumeAt: &until}
}

func WaitForInput() *suspended {
	return &suspended{input: true}
}

func ErrorHandler(h Step) *errHandler {
	return &errHandler{handler: h}
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
