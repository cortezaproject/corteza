package wfexec

import (
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"time"
)

type (

	// state holds information about Session ID
	State struct {
		created   time.Time
		completed *time.Time

		// state identifier
		stateId uint64

		// who's running this?
		owner auth.Identifiable

		// Session identifier
		sessionId uint64

		// parent, parent step
		parent Step

		// current step
		step Step

		// next steps
		next Steps

		// step error (if any)
		err error

		// input variables that were sent to resume the session
		input *expr.Vars

		// scope
		scope *expr.Vars

		// step execution results
		results *expr.Vars

		// error handling step
		errHandler Step
		handled    bool

		loops []Iterator

		action string
	}
)

func NewState(ses *Session, owner auth.Identifiable, caller, current Step, scope *expr.Vars) *State {
	return &State{
		stateId:   nextID(),
		owner:     owner,
		sessionId: ses.id,
		created:   *now(),
		parent:    caller,
		step:      current,
		scope:     scope,

		loops: make([]Iterator, 0, 4),
	}
}

func FinalState(ses *Session, scope *expr.Vars) *State {
	return &State{
		stateId:   nextID(),
		sessionId: ses.id,
		created:   *now(),
		completed: now(),
		scope:     scope,
	}
}

func (s State) Next(current Step, scope *expr.Vars) *State {
	return &State{
		stateId:    nextID(),
		owner:      s.owner,
		sessionId:  s.sessionId,
		parent:     s.step,
		errHandler: s.errHandler,
		loops:      s.loops,

		step:  current,
		scope: scope,
	}
}

func (s State) MakeRequest() *ExecRequest {
	return &ExecRequest{
		SessionID: s.sessionId,
		StateID:   s.stateId,
		Scope:     s.scope,
		Input:     s.input,
		Parent:    s.parent,
	}
}

func (s *State) newLoop(i Iterator) {
	s.loops = append(s.loops, i)
}

// ends loop and returns step that leads out of the loop
func (s *State) loopEnd() (out Steps) {
	l := len(s.loops) - 1
	if l < 0 {
		panic("not inside a loop")
	}

	out = Steps{s.loops[l].Break()}
	s.loops = s.loops[:l]
	return
}

func (s State) loopCurr() Iterator {
	l := len(s.loops)
	if l > 0 {
		return s.loops[l-1]
	}

	return nil
}

func (s State) MakeFrame() *Frame {
	var (
		// might not be the most optimal way but we need to
		// un-reference scope, input, result variables
		unref = func(vars *expr.Vars) *expr.Vars {
			out := &expr.Vars{}
			tmp, _ := json.Marshal(vars)
			_ = json.Unmarshal(tmp, out)
			return out
		}
	)

	f := &Frame{
		CreatedAt: s.created,
		SessionID: s.sessionId,
		StateID:   s.stateId,
		Input:     unref(s.input),
		Scope:     unref(s.scope),
		Results:   unref(s.results),
		NextSteps: s.next.IDs(),
		Action:    s.action,
	}

	if s.err != nil {
		f.Error = s.err.Error()
	}

	if s.step != nil {
		f.StepID = s.step.ID()
	}

	if s.parent != nil {
		f.ParentID = s.parent.ID()
	}

	if s.completed != nil {
		f.StepTime = uint(s.completed.Sub(s.created) / time.Millisecond)
	}

	return f
}

func (s *State) Error() string {
	if s.err == nil {
		return ""
	}

	return s.err.Error()
}
