package wfexec

import (
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/expr"
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

		// error handler result variable names
		//
		// Holds the error/errorMessage/errorStepID to variable-name mapping
		// declared on the error handling step. Kept apart from results
		// (previous step outputs) so that a normal step running between the
		// handler and the actual error does not overwrite it.
		errHandlerResults *expr.Vars

		// error handled flag, this gets restarted on every new state!
		errHandled bool

		loops []Iterator

		// names of variables written on this path since it branched off
		//
		// Used by the join gateway to merge only what each parallel branch
		// actually changed instead of whole scopes; reset for every path
		// when execution branches.
		dirty map[string]bool

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
		stateId:           nextID(),
		owner:             s.owner,
		sessionId:         s.sessionId,
		parent:            s.step,
		errHandler:        s.errHandler,
		errHandlerResults: s.errHandlerResults,
		results:           s.results,
		loops:             s.loops,
		dirty:             s.dirty,

		step:  current,
		scope: scope,
	}
}

// NextBranch returns next state for one of the parallel paths
//
// Variables written before the branching point are not part of the branch's
// changes, so the set of written variables starts empty.
func (s State) NextBranch(current Step, scope *expr.Vars) *State {
	st := s.Next(current, scope)
	st.dirty = nil
	return st
}

// markDirty records variables written by the current step
func (s *State) markDirty(vv *expr.Vars) {
	if vv.IsEmpty() {
		return
	}

	var nn []string
	_ = vv.Each(func(k string, _ expr.TypedValue) error {
		nn = append(nn, k)
		return nil
	})

	s.markDirtyNames(nn...)
}

// markDirtyNames records variables written directly to the scope
func (s *State) markDirtyNames(nn ...string) {
	if len(nn) == 0 {
		return
	}

	if s.dirty == nil {
		s.dirty = make(map[string]bool)
	} else {
		// state's set is shared with the states it spawned; copy before write
		dirty := make(map[string]bool, len(s.dirty)+1)
		for k := range s.dirty {
			dirty[k] = true
		}
		s.dirty = dirty
	}

	for _, n := range nn {
		s.dirty[n] = true
	}
}

func (s State) MakeRequest() *ExecRequest {
	return &ExecRequest{
		SessionID: s.sessionId,
		StateID:   s.stateId,
		Scope:     s.scope,
		Input:     s.input,
		Results:   s.results,
		Parent:    s.parent,
		dirty:     s.dirty,
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

// MakeLightFrame builds a frame that describes the step without snapshotting
// any of the state's variables
//
// Cloning the scope is what makes stacktrace collection expensive, and a frame
// only needs its variables when someone is going to read them back.
func (s State) MakeLightFrame() *Frame {
	f := &Frame{
		CreatedAt: s.created,
		SessionID: s.sessionId,
		StateID:   s.stateId,
		NextSteps: s.next.IDs(),
		Action:    s.action,
	}

	s.describe(f)
	return f
}

// describe fills in the parts of a frame that do not require cloning
func (s State) describe(f *Frame) {
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
}

// MakeFrame builds a frame with a full snapshot of the state's variables
func (s State) MakeFrame() *Frame {
	var (
		// might not be the most optimal way but we need to
		// un-reference scope, input, result variables
		unref = func(vars *expr.Vars) *expr.Vars {
			aux, err := vars.Clone()
			if err != nil {
				return expr.EmptyVars()
			}

			// Since we're cloning vars, this will always hold
			return aux.(*expr.Vars)
		}
	)

	f := s.MakeLightFrame()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		f.Input = unref(s.input)
		wg.Done()
	}()

	go func() {
		f.Scope = unref(s.scope)
		wg.Done()
	}()

	go func() {
		f.Results = unref(s.results)
		wg.Done()
	}()

	wg.Wait()

	return f
}

func (s *State) Error() string {
	if s.err == nil {
		return ""
	}

	return s.err.Error()
}
