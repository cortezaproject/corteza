package wfexec

import "time"

type (
	ExecResponse interface{}
	Joined       struct{}
)

func DelayExecution(until time.Time) *suspended {
	return &suspended{resumeAt: &until}
}

func WaitForInput() *suspended {
	return &suspended{input: true}
}
